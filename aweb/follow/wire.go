//go:build wireinject

package main

import (
	"github.com/google/wire"
	grpc2 "github.com/pluckhuang/goweb/aweb/follow/grpc"
	"github.com/pluckhuang/goweb/aweb/follow/ioc"
	"github.com/pluckhuang/goweb/aweb/follow/repository"
	"github.com/pluckhuang/goweb/aweb/follow/repository/cache"
	"github.com/pluckhuang/goweb/aweb/follow/repository/dao"
	"github.com/pluckhuang/goweb/aweb/follow/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMFollowRelationDAO,
	cache.NewRedisFollowCache,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc2.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.InitLogger,
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
