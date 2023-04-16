package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cast"
	"github.com/xince-fun/FreeMall/app/leaf/global"
	"github.com/xince-fun/FreeMall/app/leaf/model"
	"github.com/xince-fun/FreeMall/pkg/dir"
	"github.com/xince-fun/FreeMall/pkg/errno"
	"github.com/xince-fun/FreeMall/pkg/ip"
	"github.com/xince-fun/FreeMall/pkg/times"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type SnowflakeIdGenRepo interface {
	GetPrefixKey(ctx context.Context, prefix string) (*clientv3.GetResponse, error)
	CreateKeyWithOptLock(ctx context.Context, key string, val string) bool
	CreateOrUpdateKey(ctx context.Context, key string, val string) bool
	GetKey(ctx context.Context, key string) (*clientv3.GetResponse, error)
}

var (
	// 起始的时间戳，用于用当前时间戳减去这个时间戳，算出偏移量
	twepoch            = int64(1288834974657)
	workerIdBits       = 10
	maxWorkerId        = -1 ^ (-1 << workerIdBits) // 最多1023个节点
	sequenceBits       = 12
	workerIdShift      = sequenceBits
	timestampLeftShift = sequenceBits + workerIdBits
	sequenceMask       = int64(-1 ^ (-1 << sequenceBits))
	PrefixEtcdPath     = "/leaf_snowflake/"
	PropPath           string
	LeafForever        = PrefixEtcdPath + "/forever"   // 保持所有数据持久的节点
	LeafTemporaryKey   = PrefixEtcdPath + "/temporary" // 临时节点
)

type SnowflakeIdGenUseCase struct {
	repo SnowflakeIdGenRepo

	twepoch       int64
	workerId      int64
	sequence      int64
	lastTimestamp int64
	snowflakeLock sync.Mutex

	model.SnowflakeEtcdHolder
}

// NewSnowflakeIdGenUseCase new a snowflake usecase.
func NewSnowflakeIdGenUseCase(repo SnowflakeIdGenRepo) *SnowflakeIdGenUseCase {
	return &SnowflakeIdGenUseCase{
		repo: repo,
	}
}

// GetSnowflakeID creates a Snowflake ID  运行时，leaf允许最多5ms的回拨；重启时，允许最多3s的回拨
func (uc *SnowflakeIdGenUseCase) GetSnowflakeID(ctx context.Context) (int64, error) {
	if global.GlobalServerConfig.EtcdConfig.SnowflakeEnable {
		// 多个共享变量会被并发访问，同一时间内，应只让一个goroutine访问共享变量
		uc.snowflakeLock.Lock()
		defer uc.snowflakeLock.Unlock()

		var ts = time.Now().UnixMilli()

		if ts < uc.lastTimestamp {
			offset := uc.lastTimestamp - ts
			if offset <= 5 {
				// 等待 2*offset ms就可以唤醒重新尝试获取锁继续执行
				time.Sleep(time.Duration(offset<<1) * time.Millisecond)
				// 重新获取当前时间戳，理论上这次应该比上次记录的时间戳迟了
				ts := time.Now().UnixMilli()
				if ts < uc.lastTimestamp {
					return 0, errno.ErrSnowflakeTimeException
				}
			} else {
				return 0, errno.ErrSnowflakeTimeException
			}
		}

		// 如果从上一个逻辑分支产生的timestamp仍然和lastTimestamp相等
		if uc.lastTimestamp == ts {
			// 自增序列+1然后取后12位的值
			uc.sequence = (uc.sequence + 1) & sequenceMask
			// 当 seq == 0 时，说明当前毫秒已经用完了，需要等待下一毫秒
			if uc.sequence == 0 {
				// 对自增序列做随机起始
				uc.sequence = int64(rand.Intn(100))
				// 生成戳lastTimestamp滞后的时间戳，这里不进行sleep，因为很快就能获得滞后的毫秒数
				ts = times.UtilNextMillis(uc.lastTimestamp)
			}
		} else {
			// 如果是新的ms开始
			uc.sequence = int64(rand.Intn(100))
		}
		uc.lastTimestamp = ts
		id := ((ts - uc.twepoch) << timestampLeftShift) | (uc.workerId << workerIdShift) | uc.sequence

		return id, nil
	} else {
		return 0, nil
	}
}

func initSnowflake(s *SnowflakeIdGenUseCase) {
	s.twepoch = twepoch
	if !(time.Now().UnixMilli() > twepoch) {
		panic("Snowflake not support twepoch gt CurrentTime")
	}

	s.SnowflakeEtcdHolder.Ip = ip.GetOutboundIP()
	// TODO: fix this
	s.SnowflakeEtcdHolder.Port = "9001"

	PrefixEtcdPath = "/snowflake/" + "LEAF_SNOWFLAKE"
	PropPath = filepath.Join(dir.GetCurrentAbPath(), "LEAF_SNOWFLAKE") + "/leafconf/" +
		s.SnowflakeEtcdHolder.Port + "/workerID.toml"
	LeafForever = PrefixEtcdPath + "/forever"
	klog.Info("workerId local cache path: ", PropPath)

	s.SnowflakeEtcdHolder.ListenAddress = s.SnowflakeEtcdHolder.Ip + ":" + s.SnowflakeEtcdHolder.Port
	if !s.initSnowflakeWorkId() {
		klog.Error("Snowflake init workId failed")
		global.GlobalServerConfig.EtcdConfig.SnowflakeEnable = false
	} else {
		s.workerId = int64(s.SnowflakeEtcdHolder.WorkerId)
	}
	if !(s.workerId >= 0 && s.workerId <= int64(maxWorkerId)) {
		panic("Snowflake worker Id can't be greater than " + string(maxWorkerId) + " or less than 0")
	}
}

func (uc *SnowflakeIdGenUseCase) initSnowflakeWorkId() bool {
	var retryCount = 0
RETRY:
	prefixKeyResps, err := uc.repo.GetPrefixKey(newTimeoutCtx(time.Second*2), LeafForever)
	if err == nil {
		// 还没有实例化过
		if prefixKeyResps.Count == 0 {
			uc.SnowflakeEtcdHolder.EtcdAddressNode = LeafForever + "/" + uc.ListenAddress + "-0"
			if success := uc.repo.CreateKeyWithOptLock(newTimeoutCtx(time.Second),
				uc.SnowflakeEtcdHolder.EtcdAddressNode,
				string(uc.buildData())); !success {
				// 其他示例已经创建
				if retryCount > 3 {
					return false
				}
				retryCount++
				goto RETRY
			}
			uc.updateLocalWorkerID(uc.WorkerId)
			go uc.scheduledUploadData(uc.SnowflakeEtcdHolder.EtcdAddressNode)
		} else {
			// 存在的话，说明不是第一次启动leaf应用，etcd存在以前的
			// 自身节点ip:port->1
			nodeMap := make(map[string]int, 0)
			// 自身节点ip:port -> path/ip:port-1
			realNodeMap := make(map[string]string, 0)

			for _, node := range prefixKeyResps.Kvs {
				nodeKey := strings.Split(filepath.Base(string(node.Key)), "-")
				realNodeMap[nodeKey[0]] = string(node.Key)
				nodeMap[nodeKey[0]] = cast.ToInt(nodeKey[1])
			}

			if workId, ok := nodeMap[uc.SnowflakeEtcdHolder.ListenAddress]; ok {
				uc.SnowflakeEtcdHolder.EtcdAddressNode = realNodeMap[uc.SnowflakeEtcdHolder.ListenAddress]
				uc.WorkerId = workId
				if !uc.checkInitTimeStamp(uc.SnowflakeEtcdHolder.EtcdAddressNode) {
					klog.Error("init timestamp check failed, forever node timestamp gt this node time")
					return false
				}

				go uc.scheduledUploadData(uc.SnowflakeEtcdHolder.EtcdAddressNode)
				uc.updateLocalWorkerID(uc.WorkerId)
			} else {
				// 不存在自己的节点则表示是一个新启动的节点，则创建持久节点，不需要check时间
				workId := 0

				// 找到最大的ID
				for _, id := range nodeMap {
					if workId < id {
						workId = id
					}
				}
				uc.SnowflakeEtcdHolder.WorkerId = workId + 1
				uc.SnowflakeEtcdHolder.EtcdAddressNode = LeafForever + "/" + uc.ListenAddress +
					fmt.Sprintf("-%d", uc.SnowflakeEtcdHolder.WorkerId)
				if success := uc.repo.CreateKeyWithOptLock(newTimeoutCtx(time.Second),
					uc.SnowflakeEtcdHolder.EtcdAddressNode,
					string(uc.buildData())); !success {
					// 其他示例已经创建
					if retryCount > 3 {
						return false
					}
					retryCount++
					goto RETRY
				}
				go uc.scheduledUploadData(uc.SnowflakeEtcdHolder.EtcdAddressNode)
				uc.updateLocalWorkerID(uc.WorkerId)
			}
		}
	} else {
		klog.Error("start node ERROR: ", err)
		// 读不到etcd就尝试从本地读取
		if _, err := os.Stat(PropPath); err == nil {
			readFile, err := os.ReadFile(PropPath)
			if err != nil {
				klog.Error("Snowflake read local workerId failed, err: ", err)
				return false
			}
			split := strings.Split(string(readFile), "=")
			uc.WorkerId = cast.ToInt(split[1])
			klog.Warnf("START FAILED ,use local node file properties workerID-{%d}", uc.WorkerId)
		} else {
			klog.Error("workerID file not exist...")
			return false
		}
	}

	return true
}

func (uc *SnowflakeIdGenUseCase) updateLocalWorkerID(workId int) {
	if _, err := os.Stat(PropPath); err != nil {
		os.MkdirAll(filepath.Dir(PropPath), os.ModePerm)
	}

	err := os.WriteFile(PropPath, []byte(fmt.Sprintf("workerID=%d", workId)), os.ModePerm)
	if err != nil {
		klog.Error("Snowflake update local workerId failed, err: ", err)
		return
	}
	return
}

func (uc *SnowflakeIdGenUseCase) scheduledUploadData(addressNode string) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			uc.updateNewData(addressNode)
		}
	}
}

func (uc *SnowflakeIdGenUseCase) updateNewData(path string) {
	if time.Now().UnixMilli() < uc.SnowflakeEtcdHolder.LastUpdateTime {
		return
	}
	success := uc.repo.CreateOrUpdateKey(context.TODO(), path, string(uc.buildData()))
	if !success {
		return
	}
	uc.SnowflakeEtcdHolder.LastUpdateTime = time.Now().UnixMilli()
	return
}

func (uc *SnowflakeIdGenUseCase) buildData() []byte {
	endPoint := new(model.Endpoint)
	endPoint.IP = uc.SnowflakeEtcdHolder.Ip
	endPoint.Port = uc.SnowflakeEtcdHolder.Port
	endPoint.Timestamp = time.Now().UnixMilli()
	encodeAddr, _ := json.Marshal(endPoint)
	return encodeAddr
}

func (uc *SnowflakeIdGenUseCase) deBuildData(val []byte) *model.Endpoint {
	endPoint := new(model.Endpoint)
	_ = json.Unmarshal(val, endPoint)
	return endPoint
}

func (uc *SnowflakeIdGenUseCase) checkInitTimeStamp(addressNode string) bool {
	getKey, err := uc.repo.GetKey(context.TODO(), addressNode)
	if err != nil {
		return false
	}
	endpoint := uc.deBuildData(getKey.Kvs[0].Value)
	return !(endpoint.Timestamp > time.Now().UnixMilli())
}

func newTimeoutCtx(duration time.Duration) context.Context {
	timeoutCtx, _ := context.WithTimeout(context.TODO(), duration)
	return timeoutCtx
}
