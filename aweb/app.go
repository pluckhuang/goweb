package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pluckhuang/goweb/aweb/internal/events"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
