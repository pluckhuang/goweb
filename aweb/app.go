package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pluckhuang/goweb/aweb/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
