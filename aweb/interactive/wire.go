//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/interactive/events"
	"github.com/pluckhuang/goweb/aweb/interactive/grpc"
	ioc "github.com/pluckhuang/goweb/aweb/interactive/ioc"
	repository2 "github.com/pluckhuang/goweb/aweb/interactive/repository"
	cache2 "github.com/pluckhuang/goweb/aweb/interactive/repository/cache"
	dao2 "github.com/pluckhuang/goweb/aweb/interactive/repository/dao"
	service2 "github.com/pluckhuang/goweb/aweb/interactive/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.InitSaramaClient,
	ioc.InitSaramaSyncProducer,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	events.NewInteractiveProducer,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		events.NewInteractiveSyncEventConsumer,
		ioc.InitConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
