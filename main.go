package main

import (
	"easyimage_go/biz/mw"
	genrouter "easyimage_go/biz/router"
	"easyimage_go/docs"
	"easyimage_go/utils/config"
	"easyimage_go/utils/logger"
	"embed"
	_ "embed"
	"fmt"
	"strings"

	"github.com/gookit/slog"
	swaggerfiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed config/default.yaml
var defaultConfigContent []byte

//go:embed static/*
var staticFS embed.FS

// @contact.name buyfakett
// @contact.url https://github.com/buyfakett

// @securityDefinitions.apikey	ApiKeyAuth
// @in	header
// @name authorization
func main() {
	config.InitConfig(defaultConfigContent)
	logger.InitLog(config.Cfg.Server.LogLevel)
	if config.Cfg.Server.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	r := gin.Default()
	r.Use(mw.StaticFileMiddleware(staticFS))
	r.Static(config.Cfg.Image.Uri, "."+config.Cfg.Image.Uri)

	// 注册路由
	genrouter.RegisterRoutes(r)

	docs.SwaggerInfo.Version = config.Cfg.Server.Version
	docs.SwaggerInfo.Title = config.Cfg.Server.Name
	docs.SwaggerInfo.Description = fmt.Sprintf("%s by [%s](https://github.com/%s).",
		config.Cfg.Server.Name, config.Cfg.Server.Author, config.Cfg.Server.Author)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 注册swagger文档
	if config.Cfg.Server.EnableSwagger {
		slog.Info("Swagger文档已启用")
		r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	if config.Cfg.Server.LogLevel == "debug" {
		slog.Infof("服务启动成功，地址为 http://localhost:%d", config.Cfg.Server.Port)
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, config.Cfg.Image.Uri) {
			c.String(404, "404 - 文件不存在")
			return
		}
		// 其他404可以保持你原来的处理方式
		c.JSON(404, gin.H{"code": 404, "msg": "你访问的页面不存在"})
	})

	// 启动服务
	port := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	if err := r.Run(port); err != nil {
		panic(err)
	}
}
