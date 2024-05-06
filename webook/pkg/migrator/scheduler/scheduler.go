package scheduler

import (
	"context"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/gormx/connpool"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
	"time"
)

// Scheduler
// 统一管理整个迁移过程
// 为了方便用户操作和理解而引入的
type Scheduler[T migrator.Entity] struct {
	lock sync.Mutex
	src *gorm.DB
	dst *gorm.DB
	pool *connpool.DoubleWritePool
	l logger.LoggerV1
	pattern string
	cancelFull func()
	cancelIncr func()
	producer events.Producer

	// 允许多个全量校验同时运行
	fulls map[string]func()
}

func NewScheduler[T migrator.Entity](
	l logger.LoggerV1,
	src *gorm.DB,
	dst *gorm.DB,
	pool *connpool.DoubleWritePool,
	producer events.Producer) *Scheduler[T] {
	return &Scheduler[T]{
		l: l,
		src: src,
		dst: dst,
		pattern: connpool.PatternSrcOnly,
		cancelFull: func() {

		},
		cancelIncr: func() {

		},
		pool: pool,
		producer: producer,
	}
}

// 这一个也不是必须的，就是你可以考虑利用配置中心，监听配置中心的变化
// 把全量校验，增量校验做成分布式任务，利用分布式任务调度平台来调度
func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	// 将这个暴露为 HTTP 接口
	// 你可以配上对应的 UI
	server.POST("/src_only", ginx.Wrap(s.SrcOnly))
	server.POST("/src_first", ginx.Wrap(s.SrcFirst))
	server.POST("/dst_first", ginx.Wrap(s.DstFirst))
	server.POST("/dst_only", ginx.Wrap(s.DstOnly))
	server.POST("/full/start", ginx.Wrap(s.StartFullValidation))
	server.POST("/full/stop", ginx.Wrap(s.StopFullValidation))
	server.POST("/incr/stop", ginx.Wrap(s.StopIncrementValidation))
	server.POST("/incr/start", ginx.WrapBodyV1[StartIncrRequest](s.StartIncrementValidation))
}

// ---- 下面是四个阶段 ---- //

func (s *Scheduler[T]) SrcOnly(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcOnly
	s.pool.UpdatePattern(connpool.PatternSrcOnly)
	return ginx.Result{
		Msg: "Src only, OK",
	}, nil
}

func (s *Scheduler[T]) SrcFirst(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcFirst
	s.pool.UpdatePattern(connpool.PatternSrcFirst)
	return ginx.Result{
		Msg: "Src first, OK",
	}, nil
}

func (s *Scheduler[T]) DstOnly(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstOnly
	s.pool.UpdatePattern(connpool.PatternDstOnly)
	return ginx.Result{
		Msg: "Dst only, OK",
	}, nil
}

func (s *Scheduler[T]) DstFirst(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstFirst
	s.pool.UpdatePattern(connpool.PatternDstFirst)
	return ginx.Result{
		Msg: "Dst first, OK",
	}, nil
}

// StartIncrementValidation
// 开启增量校验
func (s *Scheduler[T]) StartIncrementValidation(c *gin.Context,
	req StartIncrRequest) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 开启新的时候，关掉旧的
	cancel := s.cancelIncr
	v, err := s.newValidator()
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg: "增量校验 系统异常",
		}, nil
	}
	v.Incr().Utime(req.Utime).
		SleepInterval(time.Duration(req.Interval) * time.Millisecond)

	go func() {
		var ctx context.Context
		ctx, s.cancelIncr = context.WithCancel(context.Background())
		cancel()
		err := v.Validate(ctx)
		s.l.Warn("退出增量校验", logger.Error(err))
	}()
	return ginx.Result{
		Msg: "启动增量校验成功",
	}, nil
}

// StopIncrementValidation
// 停止增量校验
func (s *Scheduler[T]) StopIncrementValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelIncr()
	return ginx.Result{
		Msg: "停止增量校验，OK",
	}, nil
}

// StartFullValidation
// 开启全量校验
func (s *Scheduler[T]) StartFullValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cancel := s.cancelFull
	v, err := s.newValidator()
	if err != nil {
		return ginx.Result{}, err
	}
	var ctx context.Context
	ctx, s.cancelFull = context.WithCancel(context.Background())

	go func() {
		// 开启新的时候，关掉旧的
		cancel()
		err := v.Validate(ctx)
		if err != nil {
			s.l.Warn("退出增量校验", logger.Error(err))
		}
	}()
	return ginx.Result{
		Msg: "启动全量校验成功",
	}, nil
}

// StopFullValidation
// 停止全量校验
func (s *Scheduler[T]) StopFullValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelIncr()
	return ginx.Result{
		Msg: "停止全量校验，OK",
	}, nil
}


func (s *Scheduler[T]) newValidator() (*validator.Validator[T], error) {
	switch s.pattern {
	case connpool.PatternSrcOnly, connpool.PatternSrcFirst:
		return validator.NewValidator[T](s.src, s.dst, "SRC", s.l, s.producer), nil
	case connpool.PatternDstOnly, connpool.PatternDstFirst:
		return validator.NewValidator[T](s.dst, s.src, "DST", s.l, s.producer), nil
	}
	return nil, fmt.Errorf("未知的 pattern %s", s.pattern)
}

type StartIncrRequest struct {
	Utime int64 `json:"utime"`
	// 毫秒数
	// json 不能正确处理 time.Duration 类型
	Interval int64 `json:"interval"`
}