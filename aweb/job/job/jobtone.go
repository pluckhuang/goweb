package job

import (
	"fmt"
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type TestJob struct {
	*LockableJob
}

func NewTestJob(
	l logger.LoggerV1,
	client *rlock.Client,
	timeout time.Duration) *TestJob {
	return &TestJob{
		LockableJob: NewLockableJob(l, client, "one:job", timeout),
	}
}

func (r *TestJob) Run() error {
	if !r.EnsureLock() {
		return nil
	}
	fmt.Println("TestJob is running over")
	return nil
}

func (r *TestJob) Name() string {
	return "TestJob"
}

func (r *TestJob) Close() error {
	return r.Unlock()
}
