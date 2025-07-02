//go:build wireinject

package main

import (
	"github.com/google/wire"
	grpc2 "github.com/pluckhuang/goweb/aweb/comment/grpc"
	ioc "github.com/pluckhuang/goweb/aweb/comment/ioc"
	"github.com/pluckhuang/goweb/aweb/comment/repository"
	"github.com/pluckhuang/goweb/aweb/comment/repository/dao"
	"github.com/pluckhuang/goweb/aweb/comment/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitDB,
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
