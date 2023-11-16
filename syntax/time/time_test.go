package time

import (
	"context"
	"testing"
	"time"
)

// 等间隔循环执行
func TestTick(t *testing.T) {
	tm := time.NewTicker(time.Second)
	// 这一句防gorutine泄露
	defer tm.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			t.Log("over time, or cancel")
		case now := <-tm.C:
			t.Log(now.Unix())
		}
	}
	t.Log("退出了循环")
}

// 等时间间隔执行
func TestTimer(t *testing.T) {
	tm := time.NewTimer(time.Second)
	defer tm.Stop()
	go func() {
		for now := range tm.C {
			t.Log(now.Unix())
		}
	}()

	time.Sleep(time.Second * 10)
}
