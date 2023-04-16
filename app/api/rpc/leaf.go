package rpc

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/xince-fun/FreeMall/kitex_gen/leaf"
	"github.com/xince-fun/FreeMall/kitex_gen/leaf/leafservice"
)

var leafClient leafservice.Client

func initLeaf() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	c, err := leafservice.NewClient(
		"LeafService",
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	hlog.Infof("client: %+v", c)
	leafClient = c
}

func GetSegmentId(ctx context.Context, req *leaf.IdRequest) (string, error) {
	resp, err := leafClient.GenSegmentId(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func GetSnowflakeId(ctx context.Context, req *leaf.IdRequest) (string, error) {
	resp, err := leafClient.GenSnowflakeId(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func DecodeSnowflakeId(ctx context.Context, req *leaf.DecodeSnokflakeRequest) (*leaf.DecodeSnokflakeResponse, error) {
	resp, err := leafClient.DecodeSnowflakeId(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
