package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {

}

func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	// 这里有一个问题。你要不要检测 status 或者 version?
	// WHERE version = ?
	// 要。你们的作业记得修改
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime": time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
			"utime": time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
			"next_time": next.UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Where("id = ?", id).Updates(map[string]any{
			"status": jobStatusPaused,
			"utime": time.Now().UnixMilli(),
	}).Error
}

type GORMJobDAO struct {
	db *gorm.DB
}

type Job struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`

	// 第一个问题：哪些任务可以抢？哪些任务已经被人占着？哪些任务永远不会被运行
	// 用状态来标记
	Status int

	// 另外一个问题，定时任务，我怎么知道，已经到时间了呢？
	// NextTime 下一次被调度的时间
	// next_time <= now 这样一个查询条件
	// and status = 0
	// 要建立索引
	// 更加好的应该是 next_time 和 status 的联合索引
	NextTime int64 `gorm:"index"`
	// cron 表达式
	Cron string

	Version int

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}

const (
	jobStatusWaiting = iota
	// 已经被抢占
	jobStatusRunning
	// 还可以有别的取值

	// 暂停调度
	jobStatusPaused
)