package ioc

import (
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/cronjob/job"
	"github.com/pluckhuang/goweb/aweb/cronjob/service"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(l logger.LoggerV1, client *rlock.Client, svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(l, client, svc, time.Second*30)
}

func InitCronJob(l logger.LoggerV1, rjob *job.RankingJob) *cron.Cron {
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "goweb",
		Subsystem: "aweb",
		Name:      "cron_job",
		Help:      "定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 10s", builder.Build(rjob))
	if err != nil {
		panic(err)
	}
	return expr
}
