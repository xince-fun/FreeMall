package service

import (
	"context"

	"github.com/xince-fun/FreeMall/app/auth/model"
)

type AuthServiceUseCase struct {
	repo AuthServiceRepo
}

type AuthServiceRepo interface {
	GetByUserIdAndType(ctx context.Context, userId int64, sysType int8) (account *model.AuthAccount, err error)
	GetByUid(ctx context.Context, uid string) (account *model.AuthAccount, err error)
	UpdatePassword(ctx context.Context, userId int64, sysType int8, password string) (err error)
	UpdateAccountInfo(ctx context.Context, account *model.AuthAccount) (err error)
	DeleteByUserIdAndType(ctx context.Context, userId int64, sysType int8) (err error)
	UpdateUserInfoByUserId(ctx context.Context, userId int64, sysType int8, account *model.AuthAccount) (err error)
	GetMerchantInfoByTenantId(ctx context.Context, tenantId int64) (merchantInfo *model.AuthAccount, err error)
}
