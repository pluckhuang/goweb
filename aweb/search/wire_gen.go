// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func Init() *App {
	client := ioc.InitESClient()
	anyDAO := dao.NewAnyESDAO(client)
	anyRepository := repository.NewAnyRepository(anyDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	collectDAO := dao.NewCollectDAO(client)
	likeDAO := dao.NewLikeDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO, collectDAO, likeDAO)
	syncService := service.NewSyncService(anyRepository, articleRepository)
	syncServiceServer := grpc.NewSyncServiceServer(syncService)
	searchService := service.NewSearchService(articleRepository)
	searchServiceServer := grpc.NewSearchService(searchService)
	clientv3Client := ioc.InitEtcdClient()
	loggerV1 := ioc.InitLogger()
	server := ioc.InitGRPCxServer(syncServiceServer, searchServiceServer, clientv3Client, loggerV1)
	saramaClient := ioc.InitKafka()
	articleConsumer := events.NewArticleConsumer(saramaClient, loggerV1, syncService)
	interactiveConsumer := events.NewInteractiveConsumer(saramaClient, loggerV1, syncService)
	v := ioc.NewConsumers(articleConsumer, interactiveConsumer)
	app := &App{
		server:    server,
		consumers: v,
	}
	return app
}

// wire.go:

var serviceProviderSet = wire.NewSet(dao.NewArticleElasticDAO, dao.NewAnyESDAO, dao.NewCollectDAO, dao.NewLikeDAO, repository.NewArticleRepository, repository.NewAnyRepository, service.NewSyncService, service.NewSearchService)

var thirdProvider = wire.NewSet(ioc.InitESClient, ioc.InitEtcdClient, ioc.InitLogger, ioc.InitKafka)
