package router

import (
	"easyimage_go/biz/handler"

	"github.com/gin-gonic/gin"
)

func diyRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/ping", handler.Ping)
	apiGroup.GET("/server_info", handler.ServerInfo)
	apiGroup.GET("/metrics", handler.Metrics)
}
