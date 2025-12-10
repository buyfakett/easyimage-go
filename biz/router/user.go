package router

import (
	hUser "easyimage_go/biz/handler/user"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/test_token", hUser.TestToken)
		userGroup.GET("/captcha", hUser.GenerateCaptcha)
	}
}
