//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/oauth2/grpc"
	ioc "github.com/pluckhuang/goweb/oauth2/ioc"
	"github.com/pluckhuang/goweb/oauth2/service"
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitRedis,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitService,
		service.NewOauth2Service,
		grpc.NewOauth2ServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
