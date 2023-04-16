package model

import (
	"go.uber.org/atomic"
	"sync"
)

type SegmentBuffer struct {
	SegmentBufferParam
	Segments []*Segment // 双buffer
}

type SegmentBufferParam struct {
	BizKey        string
	RWMutex       *sync.RWMutex
	CurrentPos    int          // 当前的使用的segment的index
	NextReady     bool         // 下一个segment是否处于可切换状态
	InitOk        bool         // 是否初始化完成
	ThreadRunning *atomic.Bool // 线程是否在运行中
	Step          int          // 步长
	MinStep       int
	UpdatedTime   int64
}

func NewSegmentBuffer() *SegmentBuffer {
	s := new(SegmentBuffer)
	s.Segments = make([]*Segment, 0, 0)
	segment1 := NewSegment(s)
	segment2 := NewSegment(s)
	s.Segments = append(s.Segments, segment1, segment2)
	s.CurrentPos = 0
	s.NextReady = false
	s.InitOk = false
	s.ThreadRunning = atomic.NewBool(false)
	s.RWMutex = &sync.RWMutex{}
	return s
}

func (b *SegmentBuffer) GetKey() string {
	return b.BizKey
}

func (b *SegmentBuffer) SetKey(key string) {
	b.BizKey = key
}

func (b *SegmentBuffer) GetSegments() []*Segment {
	return b.Segments
}

func (b *SegmentBuffer) GetAnotherSegment() *Segment {
	return b.Segments[b.NextPos()]
}

func (b *SegmentBuffer) GetCurrent() *Segment {
	return b.Segments[b.CurrentPos]
}

func (b *SegmentBuffer) NextPos() int {
	return (b.CurrentPos + 1) & 1
}

func (b *SegmentBuffer) SwitchPos() {
	b.CurrentPos = b.NextPos()
}

func (b *SegmentBuffer) IsInitOk() bool {
	return b.InitOk
}

func (b *SegmentBuffer) SetInitOk(initOk bool) {
	b.InitOk = initOk
}

func (b *SegmentBuffer) IsNextReady() bool {
	return b.NextReady
}

func (b *SegmentBuffer) SetNextReady(nextReady bool) {
	b.NextReady = nextReady
}

func (b *SegmentBuffer) GetThreadRunning() *atomic.Bool {
	return b.ThreadRunning
}

func (b *SegmentBuffer) RLock() {
	b.RWMutex.RLock()
}

func (b *SegmentBuffer) RUnLock() {
	b.RWMutex.RUnlock()
}

func (b *SegmentBuffer) WLock() {
	b.RWMutex.Lock()
}

func (b *SegmentBuffer) WUnLock() {
	b.RWMutex.Unlock()
}

func (b *SegmentBuffer) GetStep() int {
	return b.Step
}

func (b *SegmentBuffer) SetStep(step int) {
	b.Step = step
}

func (b *SegmentBuffer) GetMinStep() int {
	return b.MinStep
}

func (b *SegmentBuffer) SetMinStep(minStep int) {
	b.MinStep = minStep
}

func (b *SegmentBuffer) GetUpdateTimeStamp() int64 {
	return b.UpdatedTime
}

func (b *SegmentBuffer) SetUpdateTimeStamp(ts int64) {
	b.UpdatedTime = ts
}
