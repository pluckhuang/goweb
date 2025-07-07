package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/pluckhuang/goweb/aweb/job/job"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 初始化日志
func InitLogger() logger.LoggerV1 {
	cfg := zap.NewDevelopmentConfig()
	if err := viper.UnmarshalKey("log", &cfg); err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

// 初始化 Redis 客户端
func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// 初始化分布式锁客户端
func InitRLock(client *redis.Client) *rlock.Client {
	return rlock.NewClient(client)
}

// 初始化所有 Job
func InitJobs(l logger.LoggerV1, client *rlock.Client) (job.Job, job.Job) {
	rjob := job.NewTestJob(l, client, 30*time.Second)
	twoJob := job.NewTwoJob(l, client, 30*time.Second)
	return rjob, twoJob
}

func initCronJobBuilder(l logger.LoggerV1) *job.CronJobBuilder {
	return job.NewCronJobBuilder(l, prometheus.SummaryOpts{
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
}

// 初始化 Cron 并注册任务
func InitCron(l logger.LoggerV1, rjob job.Job, twoJob job.Job) *cron.Cron {
	builder := initCronJobBuilder(l)
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddJob("@every 1s", builder.Build(rjob)); err != nil {
		panic(err)
	}
	if _, err := c.AddJob("@every 5s", builder.Build(twoJob)); err != nil {
		panic(err)
	}
	return c
}

func main() {
	// 初始化依赖
	logger := InitLogger()
	redisClient := InitRedis()
	lockClient := InitRLock(redisClient)
	rjob, twoJob := InitJobs(logger, lockClient)
	cron := InitCron(logger, rjob, twoJob)

	// 启动定时任务
	cron.Start()

	// 监听系统信号，实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 收到信号后停止定时任务
	<-cron.Stop().Done()
}
