package router

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	diyRoutes(apiGroup)
	imageRoutes(apiGroup)
	userRoutes(apiGroup)
}
