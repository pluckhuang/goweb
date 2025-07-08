package job

import (
	"fmt"
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type TwoJob struct {
	*LockableJob
}

func NewTwoJob(
	l logger.LoggerV1,
	client *rlock.Client,
	timeout time.Duration) *TwoJob {
	return &TwoJob{
		LockableJob: NewLockableJob(l, client, "two:job", timeout),
	}
}

func (r *TwoJob) Run() error {
	if !r.EnsureLock() {
		return nil
	}
	fmt.Println("TwoJob is running over")
	return nil
}

func (r *TwoJob) Name() string {
	return "TwoJob"
}

func (r *TwoJob) Close() error {
	return r.Unlock()
}
