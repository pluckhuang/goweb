package ioc

import "github.com/gin-gonic/gin"

func InitWebServer() *gin.Engine {
	server := gin.Default()
	return server
}
