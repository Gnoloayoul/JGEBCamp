package validator

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base          *gorm.DB
	target        *gorm.DB
	l             logger.LoggerV1
	p             events.Producer
	direction     string
	batchSize     int
	highLoad      *atomicx.Value[bool]
	utime         int64
	sleepInterval time.Duration
	fromBase      func(ctx context.Context, offset int) (T, error)
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	l logger.LoggerV1,
	p events.Producer,
	direction string) *Validator[T] {
	highLoad := atomicx.NewValueOf[bool](false)
	go func() {
		// TODO: 性能判断，优先看数据库，再结合 CPU 与内存
	}()
	res := &Validator[T]{base: base, target: target,
		l: l, p: p, direction: direction,
		highLoad: highLoad}
	res.fromBase = res.fullFromBase
	return res
}

func (v *Validator[T]) SleepInterval(i time.Duration) *Validator[T] {
	v.sleepInterval = i
	return v
}

func (v *Validator[T]) Utime(utime int64) *Validator[T] {
	v.utime = utime
	return v
}

// Incr
// 发动增量校验
func (v *Validator[T]) Incr() *Validator[T] {
	v.fromBase = v.intrFromBase
	return v
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
	offset := 0
	for {
		if v.highLoad.Load() {
			// ToDo: 校验挂起
		}
		// ToDo: 改成批量处理
		src, err := v.fromBase(ctx, offset)
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			// 超时或者取消了
			return
		case nil:
			var dst T
			err = v.target.Where("id = ?", src.ID()).First(&dst).Error
			switch err {
			case context.Canceled, context.DeadlineExceeded:
				// 超时或者取消了
				return
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
			}
		case gorm.ErrRecordNotFound:
			// 到这里，数据比完了，全量校验结束
			// 想要同时支持全量校验 + 增量校验，这里不能直接退出
			// 同时，要考虑到用户的需要：有些情况希望退出，有些不想退出
			// 当用户希望继续时， sleep 一下
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		default:
			v.l.Error("校验数据，查询 base 出错", logger.Error(err))

			// 课堂演示方便，你可以删掉
			time.Sleep(time.Second)

		}
		offset++
	}
}

func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	// 找到了 base 中的数据
	// 例如 .Order("id DESC")，每次插入数据，就会导致你的 offset 不准了
	// 如果我的表没有 id 这个列怎么办？
	// 找一个类似的列，比如说 ctime (创建时间）
	// 作业。你改成批量，性能要好很多
	err := v.base.WithContext(dbCtx).
		// 最好不要取等号
		Offset(offset).
		Order("id").First(&src).Error
	return src, err
}

func (v *Validator[T]) intrFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	// 找到了 base 中的数据
	// 例如 .Order("id DESC")，每次插入数据，就会导致你的 offset 不准了
	// 如果我的表没有 id 这个列怎么办？
	// 找一个类似的列，比如说 ctime (创建时间）
	// 作业。你改成批量，性能要好很多
	err := v.base.WithContext(dbCtx).
		// 最好不要取等号
		Where("utime > ?", v.utime).
		Offset(offset).
		Order("utime ASC, id ASC").First(&src).Error
	return src, err
}

// validateTargetToBase 反向校验
// 根据 v.batchSize 的宽度返回一组的查询结果
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := 0
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)

		var dstTs []T
		err := v.target.WithContext(dbCtx).
			Where("utime > ?", v.utime).
			Select("id").
			Offset(offset).Limit(v.batchSize).
			Order("utime").Find(&dstTs).Error
		cancel()

		if len(dstTs) == 0 {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		}

		switch err {
		case context.Canceled, context.DeadlineExceeded:
			// 超时或者被人取消了
			return
		case gorm.ErrRecordNotFound:
			// 没数据了。直接返回
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err = v.base.Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case context.Canceled, context.DeadlineExceeded:
				// 超时或者被人取消了
				return
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			default:
				// 记录日志
			}
		default:
			// 记录日志，continue 掉
			v.l.Error("查询target 失败", logger.Error(err))
		}
		offset += len(dstTs)
		if len(dstTs) < v.batchSize {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
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
