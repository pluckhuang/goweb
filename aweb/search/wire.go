//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/search/events"
	"github.com/pluckhuang/goweb/aweb/search/grpc"
	"github.com/pluckhuang/goweb/aweb/search/ioc"
	"github.com/pluckhuang/goweb/aweb/search/repository"
	"github.com/pluckhuang/goweb/aweb/search/repository/dao"
	"github.com/pluckhuang/goweb/aweb/search/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewCollectDAO,
	dao.NewLikeDAO,
	repository.NewArticleRepository,
	repository.NewAnyRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	ioc.InitESClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitKafka)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
		grpc.NewSearchService,
		events.NewArticleConsumer,
		events.NewInteractiveConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
