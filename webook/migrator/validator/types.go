package validator

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/migrator"
	"github.com/Gnoloayoul/JGEBCamp/webook/migrator/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	l         logger.LoggerV1
	p         events.Producer
	direction string
	batchSize int
	highLoad  *atomicx.Value[bool]
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, l logger.LoggerV1, p events.Producer, direction string) *Validator[T] {
	highLoad := atomicx.NewValueOf[bool](false)
	go func() {
		// TODO: 性能判断，优先看数据库，再结合 CPU 与内存
	}()
	return &Validator[T]{
		base:      base,
		target:    target,
		l:         l,
		p:         p,
		direction: direction,
		highLoad:  highLoad,
	}
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		v.validateBaseToTarget(ctx)
		return nil
	})
	eg.Go(func() error {
		v.validateTargetToBase(ctx)
		return nil
	})
	return eg.Wait()
}

// validateBaseToTarget 校验
// 一条条校验
func (v *Validator[T]) validateBaseToTarget(ctx context.Context) {
	offset := -1
	for {
		if v.highLoad.Load() {
			// ToDo: 校验挂起
		}
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		offset++
		var src T
		// ToDo: 改成批量处理
		err := v.base.WithContext(dbCtx).Offset(offset).Order("id").First(&src).Error
		cancel()
		switch err {
		case nil:
			var dst T
			err = v.target.Where("id = ?", src.ID()).First(&dst).Error
			switch err {
			case nil:
				if !src.CompareTo(dst) {
					v.notify(ctx, src.ID(),
						events.InconsistentEventTypeNEQ)
				}
			case gorm.ErrRecordNotFound:
				v.notify(ctx, src.ID(),
					events.InconsistentEventTypeTargetMissing)
			default:
				v.l.Error("查询 target 数据失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
			return
		default:
			v.l.Error("校验数据，查询 base 出错", logger.Error(err))
			continue
		}
	}
}

// validateTargetToBase 反向校验
// 根据 v.batchSize 的宽度返回一组的查询结果
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := -v.batchSize
	for {
		offset = offset + v.batchSize
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)

		var dstTs []T
		err := v.target.WithContext(dbCtx).Select("id").
			Offset(offset).Limit(v.batchSize).Order("id").Find(&dstTs).Error
		cancel()

		if len(dstTs) == 0 {
			return
		}

		switch err {
		case gorm.ErrRecordNotFound:
			return
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err = v.base.Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			default:
				continue
			}
		default:
			continue
		}
		if len(dstTs) < v.batchSize {
			return
		}
	}
}

func (v *Validator[T]) notifyBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.notify(ctx, id, events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(ctx context.Context, id int64, typ string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.p.ProduceInconsistentEvent(ctx,
		events.InconsistentEvent{
			ID:        id,
			Direction: v.direction,
			Type:      typ,
		})
	cancel()
	if err != nil {
		// TODO: 重试发送（或者直接忽略，等下一轮的修复与校验）
		v.l.Error("发送数据不一致的消息失败", logger.Error(err))
	}
}
