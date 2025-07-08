package job

import (
	"context"
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/cronjob/service"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type RankingJob struct {
	*LockableJob
	rankingSvc service.RankingService
	timeout    time.Duration
}

func NewRankingJob(
	l logger.LoggerV1,
	client *rlock.Client,
	rankingSvc service.RankingService,
	timeout time.Duration) *RankingJob {
	return &RankingJob{
		LockableJob: NewLockableJob(l, client, "ranking:job", timeout),
		rankingSvc:  rankingSvc,
		timeout:     timeout,
	}
}

func (r *RankingJob) Run() error {
	if !r.EnsureLock() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.rankingSvc.TopN(ctx)
}

func (r *RankingJob) Name() string {
	return "RankingJob"
}

func (r *RankingJob) Close() error {
	return r.Unlock()
}
