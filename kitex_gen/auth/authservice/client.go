// Code generated by Kitex v0.5.1. DO NOT EDIT.

package authservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	auth "github.com/xince-fun/FreeMall/kitex_gen/auth"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetByUserIdAndType(ctx context.Context, req *auth.GetByUserIdAndTypeReq, callOptions ...callopt.Option) (r *auth.GetByUserIdAndTypeResp, err error)
	GetByUid(ctx context.Context, req *auth.GetByUidReq, callOptions ...callopt.Option) (r *auth.GetByUidResp, err error)
	UpdatePassword(ctx context.Context, req *auth.UpdatePasswordReq, callOptions ...callopt.Option) (r *auth.UpdatePasswordResp, err error)
	GetAccountByUserName(ctx context.Context, req *auth.GetAccountByUserNameReq, callOptions ...callopt.Option) (r *auth.GetAccountByUserNameResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kAuthServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kAuthServiceClient struct {
	*kClient
}

func (p *kAuthServiceClient) GetByUserIdAndType(ctx context.Context, req *auth.GetByUserIdAndTypeReq, callOptions ...callopt.Option) (r *auth.GetByUserIdAndTypeResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetByUserIdAndType(ctx, req)
}

func (p *kAuthServiceClient) GetByUid(ctx context.Context, req *auth.GetByUidReq, callOptions ...callopt.Option) (r *auth.GetByUidResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetByUid(ctx, req)
}

func (p *kAuthServiceClient) UpdatePassword(ctx context.Context, req *auth.UpdatePasswordReq, callOptions ...callopt.Option) (r *auth.UpdatePasswordResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdatePassword(ctx, req)
}

func (p *kAuthServiceClient) GetAccountByUserName(ctx context.Context, req *auth.GetAccountByUserNameReq, callOptions ...callopt.Option) (r *auth.GetAccountByUserNameResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetAccountByUserName(ctx, req)
}
