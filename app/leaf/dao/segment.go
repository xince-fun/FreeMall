package dao

import (
	"context"
	"github.com/xince-fun/FreeMall/app/leaf/model"
	"gorm.io/gorm"
)

type SegmentIdGenRepoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewSegmentIdGenRepo(db *gorm.DB, tableName string) *SegmentIdGenRepoImpl {
	return &SegmentIdGenRepoImpl{
		db:        db,
		tableName: tableName,
	}
}

func (s *SegmentIdGenRepoImpl) GetAllLeafAllocs(ctx context.Context) (leafs []*model.LeafAlloc, err error) {
	if err := s.db.Table(s.tableName).WithContext(ctx).Find(&leafs).Error; err != nil {

		return nil, err
	}

	return
}

func (s *SegmentIdGenRepoImpl) UpdateMaxIdAndGetLeafAlloc(ctx context.Context, tag string) (leaf *model.LeafAlloc, err error) {
	// Begin
	// UPDATE table SET max_id=max_id+step WHERE biz_tag=xxx
	// SELECT tag, max_id, step FROM table WHERE biz_tag=xxx
	// Commit
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Table(s.tableName).Where("biz_tag = ?", tag).
			Update("max_id", gorm.Expr("max_id + step")).Error; err != nil {

			return err
		}
		if err = tx.Table(s.tableName).Where("biz_tag = ?", tag).First(&leaf).Error; err != nil {

			return err
		}

		return nil
	})

	return
}

func (s *SegmentIdGenRepoImpl) UpdateMaxIdByCustomStepAndGetLeafAlloc(ctx context.Context, tag string, step int) (leaf *model.LeafAlloc, err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Table(s.tableName).Where("biz_tag = ?", tag).
			Update("max_id", gorm.Expr("max_id + ?", step)).Error; err != nil {

			return err
		}

		if err = tx.Table(s.tableName).Where("biz_tag = ?", tag).First(&leaf).Error; err != nil {

			return err
		}

		return nil
	})

	return
}

func (s *SegmentIdGenRepoImpl) GetAllTags(ctx context.Context) (tags []string, err error) {
	if err := s.db.Table(s.tableName).WithContext(ctx).Pluck("biz_tag", &tags).Error; err != nil {

		return nil, err
	}

	return
}

func (s *SegmentIdGenRepoImpl) GetLeafAlloc(ctx context.Context, tag string) (leaf *model.LeafAlloc, err error) {
	if err := s.db.Table(s.tableName).WithContext(ctx).Where("biz_tag = ?", tag).First(&leaf).Error; err != nil {

		return nil, err
	}

	return
}
