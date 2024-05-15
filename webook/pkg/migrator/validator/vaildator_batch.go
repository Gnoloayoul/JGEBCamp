package validator

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

func (v *Validator[T]) ValidateV1(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		v.validateBaseToTargetBatch(ctx)
		return nil
	})
	eg.Go(func() error {
		v.validateTargetToBase(ctx)
		return nil
	})
	return eg.Wait()
}

func (v *Validator[T]) fullFromBaseBatch(ctx context.Context, offset int, batchSize int) ([]T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src []T
	err := v.base.WithContext(dbCtx).
		// 最好不要取等号
		Limit(batchSize).
		Offset(offset).
		Order("id").Find(&src).Error
	return src, err
}

// validateBaseToTargetBatch 校验
// 批量提取、校验
func (v *Validator[T]) validateBaseToTargetBatch(ctx context.Context) {
	v.batchSize = 100
	v.fromBaseBatch = v.fullFromBaseBatch
	offset := 0
	for {
		if v.highLoad.Load() {
			// ToDo: 校验挂起
		}
		// ToDo: 改成批量处理
		src, err := v.fromBaseBatch(ctx, offset, v.batchSize)
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			// 超时或者取消了
			return
		case nil:
			for _, k := range src {
				var dst T
				err = v.target.Where("id = ?", k.ID()).First(&dst).Error
				switch err {
				case context.Canceled, context.DeadlineExceeded:
					// 超时或者取消了
					return
				case nil:
					if !k.CompareTo(dst) {
						v.notify(ctx, k.ID(),
							events.InconsistentEventTypeNEQ)
					}
				case gorm.ErrRecordNotFound:
					v.notify(ctx, k.ID(),
						events.InconsistentEventTypeTargetMissing)
				default:
					v.l.Error("查询 target 数据失败", logger.Error(err))
				}
			}

		case gorm.ErrRecordNotFound:
			offset += v.batchSize
			var data T
			for {
				offset++
				err := v.base.WithContext(ctx).
					Offset(offset).
					Order("id").First(&data).Error
				src = append(src, data)
				if err != nil {
					// 到这里，数据比完了，全量校验结束
					// 想要同时支持全量校验 + 增量校验，这里不能直接退出
					// 同时，要考虑到用户的需要：有些情况希望退出，有些不想退出
					// 当用户希望继续时， sleep 一下
					if v.sleepInterval <= 0 {
						return
					}
					time.Sleep(v.sleepInterval)
					continue
				}
			}



		default:
			v.l.Error("校验数据，查询 base 出错", logger.Error(err))

			// 课堂演示方便，你可以删掉
			time.Sleep(time.Second)

		}
		offset += v.batchSize
	}
}