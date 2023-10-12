package failover

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	// 超时次数
	cnt int32

	// 阈值
	// 当 cnt > threshold， 触发切换
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 先读出来
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 成功换至下一位
			// cnt 归零
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.cnt)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		// 正常了，连续超时计数自然中断（归零）
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		// 不明错误
		return err
	}
	return err
}
