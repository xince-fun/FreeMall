package main

import (
	"context"
	"fmt"

	"github.com/xince-fun/FreeMall/app/auth/service"
	auth "github.com/xince-fun/FreeMall/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct {
	AuthServiceUseCase *service.AuthServiceUseCase
}

// GetByUserIdAndType implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetByUserIdAndType(ctx context.Context, req *auth.GetByUserIdAndTypeReq) (resp *auth.GetByUserIdAndTypeResp, err error) {
	// TODO: Your code here...
	return
}

// GetByUid implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetByUid(ctx context.Context, req *auth.GetByUidReq) (resp *auth.GetByUidResp, err error) {
	// TODO: Your code here...
	fmt.Println("test")
	return
}

// UpdatePassword implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) UpdatePassword(ctx context.Context, req *auth.UpdatePasswordReq) (resp *auth.UpdatePasswordResp, err error) {
	// TODO: Your code here...
	return
}

// GetAccountByUserName implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetAccountByUserName(ctx context.Context, req *auth.GetAccountByUserNameReq) (resp *auth.GetAccountByUserNameResp, err error) {
	// TODO: Your code here...
	return
}
