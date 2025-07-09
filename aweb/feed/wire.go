//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/feed/events"
	"github.com/pluckhuang/goweb/aweb/feed/grpc"
	"github.com/pluckhuang/goweb/aweb/feed/ioc"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
	"github.com/pluckhuang/goweb/aweb/feed/repository/dao"
	"github.com/pluckhuang/goweb/aweb/feed/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewFeedPushEventDAO,
	dao.NewFeedPullEventDAO,
	repository.NewFeedEventRepo,
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitDB,
	ioc.InitFollowClient,
	ioc.InitGlobalVal,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.RegisterHandler,
		service.NewFeedService,
		grpc.NewFeedEventGrpcSvc,
		events.NewFeedEventConsumer,
		ioc.InitConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
