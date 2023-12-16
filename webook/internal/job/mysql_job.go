package job

import (
	"context"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"golang.org/x/sync/semaphore"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未知任务，你是否注册了？ %s", j.Name)
	}
	return fn(ctx, j)
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{funcs: make(map[string]func(ctx context.Context, j domain.Job) error)}
}

// Scheduler
// 调度器
type Scheduler struct {
	execs   map[string]Executor
	svc     service.JobService
	l       logger.LoggerV1
	limiter *semaphore.Weighted
}

func NewScheduler(svc service.JobService, l logger.LoggerV1) *Scheduler {
	return &Scheduler{svc: svc, l: l,
		limiter: semaphore.NewWeighted(200),
		execs:   make(map[string]Executor)}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}
