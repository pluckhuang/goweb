//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/article/events"
	grpc2 "github.com/pluckhuang/goweb/aweb/article/grpc"
	ioc "github.com/pluckhuang/goweb/aweb/article/ioc"
	"github.com/pluckhuang/goweb/aweb/article/repository"
	"github.com/pluckhuang/goweb/aweb/article/repository/cache"
	"github.com/pluckhuang/goweb/aweb/article/repository/dao"
	"github.com/pluckhuang/goweb/aweb/article/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewArticleGORMDAO,
	cache.NewArticleRedisCache,
	repository.NewCachedArticleRepository,
	service.NewArticleService,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitKafka,
	ioc.InitSyncProducer,
	events.NewSaramaSyncProducer,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
