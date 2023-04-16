package model

import "go.uber.org/atomic"

type Segment struct {
	SegmentParam
	Buffer *SegmentBuffer
}

type SegmentParam struct {
	Value *atomic.Int64 // 内存生成的每一个id号
	MaxId int64         // 当前号段允许的最大id值
	Step  int           // 步长
}

func NewSegment(segmentBuffer *SegmentBuffer) *Segment {
	return &Segment{
		SegmentParam: SegmentParam{
			Value: atomic.NewInt64(0),
		},
		Buffer: segmentBuffer,
	}
}

func (s *Segment) GetValue() *atomic.Int64 {
	return s.Value
}

func (s *Segment) SetValue(value *atomic.Int64) {
	s.Value = value
}

func (s *Segment) GetMax() int64 {
	return s.MaxId
}

func (s *Segment) SetMax(max int64) {
	s.MaxId = max
}

func (s *Segment) GetStep() int {
	return s.Step
}

func (s *Segment) SetStep(step int) {
	s.Step = step
}

func (s *Segment) GetIdle() int64 {
	value := s.GetValue().Load()
	return s.GetMax() - value
}

// GetBuffer 获取当前号段所属的SegmentBuffer
func (s *Segment) GetBuffer() *SegmentBuffer {
	return s.Buffer
}
