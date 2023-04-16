package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	"github.com/xince-fun/FreeMall/app/leaf/model"
	"github.com/xince-fun/FreeMall/pkg/errno"
	"go.uber.org/atomic"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type SegmentIdGenRepo interface {
	GetAllLeafAllocs(ctx context.Context) (leafs []*model.LeafAlloc, err error)
	UpdateMaxIdAndGetLeafAlloc(ctx context.Context, tag string) (leaf *model.LeafAlloc, err error)
	UpdateMaxIdByCustomStepAndGetLeafAlloc(ctx context.Context, tag string, step int) (leaf *model.LeafAlloc, err error)
	GetAllTags(ctx context.Context) (tags []string, err error)
	GetLeafAlloc(ctx context.Context, tag string) (leaf *model.LeafAlloc, err error)
}

const (
	SegmentDuration = time.Minute * 15
	MaxStep         = 1000000
)

type SegmentIdGenUseCase struct {
	repo            SegmentIdGenRepo
	singleGroup     singleflight.Group
	segmentDuration int64    // 号码消耗时间
	cache           sync.Map // k biz-tag : v model.SegmentBuffer

	twepoch       int64
	workerID      int64
	sequence      int64
	lastTimestamp int64
	snowFlakeLock sync.Mutex
}

func NewSegmentIdGenUseCase(repo SegmentIdGenRepo) *SegmentIdGenUseCase {
	s := &SegmentIdGenUseCase{
		repo: repo,
	}
	// 号段模式启用
	if global.GlobalServerConfig.MysqlConfig.SegmentEnable {
		_ = s.updateCacheFromDb()
		go s.updateCacheFromDbAtEveryMinute()
	}

	return s
}

func (uc *SegmentIdGenUseCase) GetAllLeafs(ctx context.Context) ([]*model.LeafAlloc, error) {
	return nil, nil
}

// GetSegID creates a Segment, and returns the new Segment.
func (uc *SegmentIdGenUseCase) GetSegID(ctx context.Context, tag string) (int64, error) {
	if global.GlobalServerConfig.MysqlConfig.SegmentEnable {
		value, ok := uc.cache.Load(tag)
		if !ok {
			return 0, errno.ErrTagNotFound
		}
		segmentBuffer := value.(*model.SegmentBuffer)
		if !segmentBuffer.IsInitOk() {
			//
			_, err, _ := uc.singleGroup.Do(tag, func() (res interface{}, err error) {
				if !segmentBuffer.IsInitOk() {
					err := uc.updateSegmentFromDb(ctx, tag, segmentBuffer.GetCurrent())
					if err != nil {
						segmentBuffer.SetInitOk(false)
						return 0, err
					}
					segmentBuffer.SetInitOk(true)
				}
				return
			})
			if err != nil {
				return 0, nil
			}
		}

		return uc.getIdFromSegmentBuffer(ctx, segmentBuffer)
	} else {
		return 0, nil
	}
}

func (uc *SegmentIdGenUseCase) updateSegmentFromDb(ctx context.Context, bizTag string, segment *model.Segment) (err error) {

	var leafAlloc *model.LeafAlloc

	segmentBuffer := segment.GetBuffer()

	// 如果buffer没有DB数据初始化(也就是第一次进行DB数据初始化)
	if !segmentBuffer.IsInitOk() {
		leafAlloc, err = uc.repo.UpdateMaxIdAndGetLeafAlloc(ctx, bizTag)
		if err != nil {
			klog.Error("db error : ", err)
			return fmt.Errorf("db error : %s %w", err, errno.ErrDBFailed)
		}
		segmentBuffer.SetStep(leafAlloc.Step)
		segmentBuffer.SetMinStep(leafAlloc.Step)
	} else if segmentBuffer.GetUpdateTimeStamp() == 0 {
		// 如果buffer的更新时间是0（初始是0，也就是第二次调用updateSegmentFromDb()）
		leafAlloc, err = uc.repo.UpdateMaxIdAndGetLeafAlloc(ctx, bizTag)
		if err != nil {
			klog.Error("db error : ", err)
			return fmt.Errorf("db error : %s %w", err, errno.ErrDBFailed)
		}
		segmentBuffer.SetUpdateTimeStamp(time.Now().Unix())
		segmentBuffer.SetMinStep(leafAlloc.Step)
	} else {
		// 第三次以及之后的进来 动态设置nextStep
		// 计算当前更新操作和上一次更新时间差
		duration := time.Now().Unix() - segmentBuffer.GetUpdateTimeStamp()
		nextStep := segmentBuffer.GetStep()
		/**
		 *  动态调整step
		 *  1) duration < 15 分钟 : step 变为原来的2倍， 最大为 MAX_STEP
		 *  2) 15分钟 <= duration < 30分钟 : nothing
		 *  3) duration >= 30 分钟 : 缩小step, 最小为DB中配置的step
		 *
		 *  这样做的原因是认为15min一个号段大致满足需求
		 *  如果updateSegmentFromDb()速度频繁(15min多次)，也就是
		 *  如果15min这个时间就把step号段用完，为了降低数据库访问频率，
		 *  我们可以扩大step大小，相反如果将近30min才把号段内的id用完，则可以缩小step
		 */
		// duration < 15 分钟 : step 变为原来的2倍. 最大为 MAX_STEP
		if duration < int64(SegmentDuration) {
			if nextStep*2 > MaxStep {
				// do nothing
			} else {
				nextStep = nextStep * 2
			}
		} else if duration < int64(SegmentDuration)*2 {
			// do nothing
		} else {
			if nextStep/2 >= segmentBuffer.GetMinStep() {
				nextStep = nextStep / 2
			}
		}
		leafAlloc, err = uc.repo.UpdateMaxIdByCustomStepAndGetLeafAlloc(ctx, bizTag, nextStep)
		if err != nil {
			klog.Error("db error : ", err)
			return fmt.Errorf("db error : %s %w", err, errno.ErrDBFailed)
		}
		segmentBuffer.SetUpdateTimeStamp(time.Now().Unix())
		segmentBuffer.SetStep(nextStep)
		segmentBuffer.SetMinStep(leafAlloc.Step)
	}

	value := leafAlloc.MaxId - int64(segmentBuffer.GetStep())
	segment.GetValue().Store(value)
	segment.SetMax(leafAlloc.MaxId)
	segment.SetStep(segmentBuffer.GetStep())

	return
}

func (uc *SegmentIdGenUseCase) loadNextSegmentBuffer(ctx context.Context, cacheSegmentBuffer *model.SegmentBuffer) {
	segment := cacheSegmentBuffer.GetSegments()[cacheSegmentBuffer.NextPos()]
	err := uc.updateSegmentFromDb(ctx, cacheSegmentBuffer.GetKey(), segment)
	if err != nil {
		cacheSegmentBuffer.GetThreadRunning().Store(false)
		return
	}

	cacheSegmentBuffer.WLock()
	defer cacheSegmentBuffer.WUnLock()
	cacheSegmentBuffer.SetNextReady(true)
	cacheSegmentBuffer.GetThreadRunning().Store(false)

	return
}

func waitAndSleep(segmentBuffer *model.SegmentBuffer) {
	roll := 0
	for segmentBuffer.GetThreadRunning().Load() {
		roll++
		if roll > 10000 {
			time.Sleep(time.Millisecond * time.Duration(10))
			break
		}
	}
}

func (uc *SegmentIdGenUseCase) getIdFromSegmentBuffer(ctx context.Context, cacheSegmentBuffer *model.SegmentBuffer) (int64, error) {

	var (
		segment *model.Segment
		value   int64
		err     error
	)

	for {
		if value := func() int64 {
			cacheSegmentBuffer.RLock()
			defer cacheSegmentBuffer.RUnLock()

			segment = cacheSegmentBuffer.GetCurrent()
			if !cacheSegmentBuffer.IsNextReady() &&
				(segment.GetIdle() < int64(0.9*float64(segment.GetStep()))) &&
				cacheSegmentBuffer.GetThreadRunning().CAS(false, true) {

				// 协程中传入空ctx，防止主体执行完成后其调用cancel取消上下文
				go uc.loadNextSegmentBuffer(context.TODO(), cacheSegmentBuffer)
			}

			value := segment.GetValue().Load()
			segment.GetValue().Inc()
			if value < segment.GetMax() {
				return value
			}
			return 0
		}(); value != 0 {
			return value, nil
		}

		// 等待协程异步准备号段完毕
		waitAndSleep(cacheSegmentBuffer)

		value, err = func() (int64, error) {
			// 执行到这里，说明当前号段已经用完，应该切换另一个Segment号段使用
			cacheSegmentBuffer.WLock()
			defer cacheSegmentBuffer.WUnLock()

			// 重复获取value, 并发执行时，Segment可能已经被其他协程切换。再次判断, 防止重复切换Segment
			segment = cacheSegmentBuffer.GetCurrent()
			value = segment.GetValue().Load()
			segment.GetValue().Inc()
			if value < segment.GetMax() {
				return value, nil
			}

			// 执行到这里, 说明其他的协程没有进行Segment切换，
			// 并且当前号段所有号码用完，需要进行切换Segment
			// 如果准备好另一个Segment，直接切换
			if cacheSegmentBuffer.IsNextReady() {
				cacheSegmentBuffer.SwitchPos()
				cacheSegmentBuffer.SetNextReady(false)
			} else { // 如果另一个Segment没有准备好，则返回异常双buffer全部用完
				return 0, errno.ErrIDTwoSegmentsAreNull
			}
			return 0, nil
		}()
		if value != 0 || err != nil {
			return value, err
		}
	}
}

func (uc *SegmentIdGenUseCase) updateCacheFromDbAtEveryMinute() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			_ = uc.updateCacheFromDb()
		}
	}
}

func (uc *SegmentIdGenUseCase) updateCacheFromDb() (err error) {

	bizTags, err := uc.repo.GetAllTags(context.TODO())
	if err != nil {
		klog.Error("load tags error : ", err)
		return err
	}
	if len(bizTags) == 0 {
		klog.Error("no tags found")
		return nil
	}

	// 数据库中的tag
	insertTags := make([]string, 0)
	removeTags := make([]string, 0)
	// 当前cache中的所有tag
	cacheTags := map[string]struct{}{}
	uc.cache.Range(func(key, value interface{}) bool {
		cacheTags[key.(string)] = struct{}{}
		return true
	})

	// 保证cache和数据库tags同步
	// 1.db中新进的tags灌进cache，并实例化对应的SegmentBuffer
	for _, k := range bizTags {
		if _, ok := cacheTags[k]; !ok {
			insertTags = append(insertTags, k)
		}
	}
	for _, k := range insertTags {
		segmentBuffer := model.NewSegmentBuffer()
		segmentBuffer.SetKey(k)
		segment := segmentBuffer.GetCurrent()
		segment.SetValue(atomic.NewInt64(0))
		segment.SetMax(0)
		segment.SetStep(0)
		uc.cache.Store(k, segmentBuffer)
		cacheTags[k] = struct{}{}
		klog.Infof("insert tag %s into cache", k)
	}

	// 2.cache中已失效的tags从cache删除
	for _, k := range bizTags {
		if _, ok := cacheTags[k]; !ok {
			removeTags = append(removeTags, k)
		}
	}
	if len(removeTags) > 0 && len(cacheTags) > 0 {
		for _, tag := range removeTags {
			uc.cache.Delete(tag)
			klog.Infof("remove tag %s from cache", tag)
		}
	}

	return nil
}
