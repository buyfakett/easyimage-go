package router

import (
	hImage "easyimage_go/biz/handler/image"
	"easyimage_go/biz/mw"

	"github.com/gin-gonic/gin"
)

func imageRoutes(r *gin.RouterGroup) {
	imageGroup := r.Group("/image")
	{
		imageGroup.PUT("/upload", mw.TokenAuthMiddleware(), hImage.UploadImage)
		imageGroup.GET("/list", mw.TokenAuthMiddleware(), hImage.ListFiles)
		imageGroup.DELETE("/delete", mw.TokenAuthMiddleware(), hImage.DeleteFile)
	}
}
