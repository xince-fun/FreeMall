package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xince-fun/FreeMall/app/leaf/service"
	leaf "github.com/xince-fun/FreeMall/kitex_gen/leaf"
	"github.com/xince-fun/FreeMall/pkg/times"
	"strconv"
)

// LeafServiceImpl implements the last service interface defined in the IDL.
type LeafServiceImpl struct {
	segmentIdGenUseCase   *service.SegmentIdGenUseCase
	snowflakeIdGenUseCase *service.SnowflakeIdGenUseCase
}

// GenSegmentId implements the LeafServiceImpl interface.
func (s *LeafServiceImpl) GenSegmentId(ctx context.Context, request *leaf.IdRequest) (resp *leaf.IdResponse, err error) {
	resp = new(leaf.IdResponse)
	id, err := s.segmentIdGenUseCase.GetSegID(ctx, request.Tag)
	if err != nil {
		klog.Errorf("get segment id failed, err: %v", err)
		return resp, err
	}
	resp.Id = strconv.FormatInt(id, 10)
	return resp, nil
}

// GenSnowflakeId implements the LeafServiceImpl interface.
func (s *LeafServiceImpl) GenSnowflakeId(ctx context.Context, request *leaf.IdRequest) (resp *leaf.IdResponse, err error) {
	resp = new(leaf.IdResponse)
	id, err := s.snowflakeIdGenUseCase.GetSnowflakeID(ctx)
	if err != nil {
		klog.Errorf("get snowflake id failed, err: %v", err)
		return resp, err
	}
	resp.Id = strconv.FormatInt(id, 10)
	return resp, err
}

// DecodeSnowflakeId implements the LeafServiceImpl interface.
func (s *LeafServiceImpl) DecodeSnowflakeId(ctx context.Context, request *leaf.DecodeSnokflakeRequest) (resp *leaf.DecodeSnokflakeResponse, err error) {
	resp = new(leaf.DecodeSnokflakeResponse)
	snowflakeId, err := strconv.ParseInt(request.Id, 10, 64)
	if err != nil {
		klog.Errorf("parse snowflake id failed, err: %v", err)
		return resp, err
	}
	originTimestamp := (snowflakeId >> 22) + 1288834974657
	timeStr := times.GetDateTimeStr(times.UnixToMS(originTimestamp))
	resp.Timestamp = strconv.FormatInt(originTimestamp, 10) + "(" + timeStr + ")"

	workerId := (snowflakeId >> 12) ^ (snowflakeId >> 22 << 10)
	resp.WorkerId = strconv.FormatInt(workerId, 10)

	sequence := snowflakeId ^ (snowflakeId >> 12 << 12)
	resp.SequenceId = strconv.FormatInt(sequence, 10)
	return resp, err
}
