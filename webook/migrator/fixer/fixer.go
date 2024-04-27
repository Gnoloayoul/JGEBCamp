package fixer

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/migrator"
	"github.com/Gnoloayoul/JGEBCamp/webook/migrator/events"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Fixer[T migrator.Entity] struct {
	base *gorm.DB
	target *gorm.DB
	columns []string
}

// Fix
// 数据修复（最直接版本）
// ToDo: 改成批量版本
func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	var t T
	err := f.base.WithContext(ctx).
		Where("id = ?", evt.ID).First(&t).Error
	switch err {
	case nil:
		return f.target.WithContext(ctx).
			Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&t).Error
	case gorm.ErrRecordNotFound:
		return f.target.WithContext(ctx).
			Where("id = ?", evt.ID).Delete(&t).Error
	default:
		return err
	}
}

// FixV1
// 根据通知的 type 来处理
func (f *Fixer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing,
		events.InconsistentEventTypeNEQ:
		var t T
		err := f.base.WithContext(ctx).
			Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).
				Where("id = ?", evt.ID).Delete(new(T)).Error
		case nil:
			return f.target.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.WithContext(ctx).
			Where("id = ?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}

// FixV2
//
func (f *Fixer[T]) FixV2(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing,
		events.InconsistentEventTypeNEQ:
		var t T
		err := f.base.WithContext(ctx).
			Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).
				Where("id = ?", evt.ID).Delete(new(T)).Error
		case nil:
			return f.target.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.WithContext(ctx).
			Where("id = ?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}