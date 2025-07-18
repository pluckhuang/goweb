//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/internal/repository"
	"github.com/pluckhuang/goweb/aweb/internal/repository/cache"
	"github.com/pluckhuang/goweb/aweb/internal/repository/dao"
	"github.com/pluckhuang/goweb/aweb/internal/service"
	"github.com/pluckhuang/goweb/aweb/internal/web"
	ijwt "github.com/pluckhuang/goweb/aweb/internal/web/jwt"
	"github.com/pluckhuang/goweb/aweb/ioc"
)

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,

		// DAO 部分
		dao.NewUserDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
