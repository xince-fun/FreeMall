package dao

import (
	"context"
	"github.com/xince-fun/FreeMall/app/auth/model"
	"gorm.io/gorm"
)

type AuthServiceRepoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewAuthServiceRepo(db *gorm.DB, tableName string) *AuthServiceRepoImpl {
	return &AuthServiceRepoImpl{
		db:        db,
		tableName: tableName,
	}
}

func (s *AuthServiceRepoImpl) GetByUid(ctx context.Context, uid string) (account *model.AuthAccount, err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Where("uid = ?", uid).First(&account).Error; err != nil {
		return nil, err
	}
	if account == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return account, nil
}

func (s *AuthServiceRepoImpl) UpdatePassword(ctx context.Context, userId int64, sysType int8, password string) (err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Where("user_id = ? and sys_type = ?", userId, sysType).
		Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceRepoImpl) UpdateAccountInfo(ctx context.Context, account *model.AuthAccount) (err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(s.tableName).Where("user_id = ? and sys_type = ?", account.UserId, account.SysType).
			Updates(map[string]interface{}{
				"username": account.UserName,
				"password": account.Password,
				"status":   account.Status,
			}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceRepoImpl) UpdateUserInfoByUserId(ctx context.Context, userId int64, sysType int8, account *model.AuthAccount) (err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(s.tableName).Where("user_id = ? and sys_type = ?", userId, sysType).
			Updates(map[string]interface{}{
				"tenant_id": account.TenantId,
			}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceRepoImpl) DeleteByUserIdAndType(ctx context.Context, userId int64, sysType int8) (err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Where("user_id = ? and sys_type = ?", userId, sysType).Update("status", -1).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceRepoImpl) GetMerchantInfoByTenantId(ctx context.Context, tenantId int64) (merchantInfo *model.AuthAccount, err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Where("tenant_id = ? and sys_type = ?", tenantId, 2).First(&merchantInfo).Error; err != nil {
		return nil, err
	}
	return merchantInfo, nil
}

func (s *AuthServiceRepoImpl) GetByUserIdAndType(ctx context.Context, userId int64, sysType int8) (account *model.AuthAccount, err error) {
	if err := s.db.WithContext(ctx).Table(s.tableName).Where("user_id = ? and sys_type = ?", userId, sysType).First(&account).Error; err != nil {
		return nil, err
	}
	if account == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return account, nil
}
