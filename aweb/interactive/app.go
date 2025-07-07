package main

import (
	"github.com/pluckhuang/goweb/aweb/pkg/grpcx"
	"github.com/pluckhuang/goweb/aweb/pkg/saramax"
)

type App struct {
	consumers []saramax.Consumer
	server    *grpcx.Server
}
