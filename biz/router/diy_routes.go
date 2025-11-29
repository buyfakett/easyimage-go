package router

import (
	"easyimage_go/biz/handler"
	"easyimage_go/biz/mw"

	"github.com/gin-gonic/gin"
)

func diyRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/ping", handler.Ping)
	apiGroup.GET("/server_info", handler.ServerInfo)
	apiGroup.GET("/metrics", handler.Metrics)
	apiGroup.GET("/test_token", mw.TokenAuthMiddleware(), handler.TestToken)
}
