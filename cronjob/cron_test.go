package cronjob

import (
	cron "github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCronJob(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	// 通过 AddJob/AddFunc 来添加任务，该方法是线程安全的
	expr.AddJob("@every 1s", myJob{})
	expr.AddFunc("@every 3s", func() {
		t.Log("开始")
		time.Sleep(time.Second * 12)
		t.Log("结束")
	})
	expr.Start()
	time.Sleep(time.Second * 10)
	// 发出停止信号，不会调度新任务，但也不同中断已经运行的任务
	stop := expr.Stop()
	t.Log("发出停止信号")
	<-stop.Done()
	t.Log("彻底停止")
	//TODO: 试一试封装：保证一次只运行一次任务
}

type myJob struct {
}

func (m myJob) Run() {

}
