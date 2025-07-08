//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/cronjob/ioc"
	"github.com/pluckhuang/goweb/aweb/cronjob/repository"
	"github.com/pluckhuang/goweb/aweb/cronjob/service"
)

func Init() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitEtcd,
		ioc.InitRlockClient,
		ioc.InitRankingRedisCache,
		ioc.InitRankingLocalCache,
		repository.NewCachedRankingRepository,
		ioc.InitInteractiveClient,
		ioc.InitArticleClient,
		service.NewBatchRankingService,
		ioc.InitRankingJob,
		ioc.InitCronJob,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
