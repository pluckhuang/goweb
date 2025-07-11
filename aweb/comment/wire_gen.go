// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pluckhuang/goweb/aweb/comment/grpc"
	"github.com/pluckhuang/goweb/aweb/comment/ioc"
	"github.com/pluckhuang/goweb/aweb/comment/repository"
	"github.com/pluckhuang/goweb/aweb/comment/repository/dao"
	"github.com/pluckhuang/goweb/aweb/comment/service"
)

// Injectors from wire.go:

func Init() *App {
	loggerV1 := ioc.InitLogger()
	client := ioc.InitEtcdClient()
	db := ioc.InitDB(loggerV1)
	commentDAO := dao.NewCommentDAO(db)
	commentRepository := repository.NewCommentRepo(commentDAO, loggerV1)
	commentService := service.NewCommentSvc(commentRepository)
	commentServiceServer := grpc.NewGrpcServer(commentService)
	server := ioc.InitGRPCxServer(loggerV1, client, commentServiceServer)
	app := &App{
		server: server,
	}
	return app
}

// wire.go:

var serviceProviderSet = wire.NewSet(dao.NewCommentDAO, repository.NewCommentRepo, service.NewCommentSvc, grpc.NewGrpcServer)

var thirdProvider = wire.NewSet(ioc.InitEtcdClient, ioc.InitLogger, ioc.InitDB)
