package main

import (
	"easyimage_go/biz/handler/image"
	"easyimage_go/biz/mw"
	genrouter "easyimage_go/biz/router"
	"easyimage_go/docs"
	"easyimage_go/utils/config"
	"easyimage_go/utils/logger"
	"embed"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
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

//go:embed internal/version/version.txt
var version string

//	@contact.name	buyfakett
//	@contact.url	https://github.com/buyfakett

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						authorization
// 读取图片数据（支持本地文件和远程URL）
func readImageData(path string) ([]byte, string, error) {
	// 判断是远程URL还是本地文件
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// 远程URL
		resp, err := http.Get(path)
		if err != nil {
			return nil, "", fmt.Errorf("无法下载远程图片: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("远程图片下载失败，状态码: %d", resp.StatusCode)
		}

		// 读取响应体
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("读取远程图片失败: %w", err)
		}

		// 从URL中提取文件名
		urlParts := strings.Split(path, "/")
		filename := urlParts[len(urlParts)-1]
		return data, filename, nil
	} else {
		// 本地文件
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, "", fmt.Errorf("本地文件不存在: %s", path)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, "", fmt.Errorf("读取本地文件失败: %w", err)
		}

		filename := path
		return data, filename, nil
	}
}

func main() {
	config.InitConfig(defaultConfigContent)
	// 如果显示版本信息，直接退出
	if config.CliCfg.ShowVersion {
		config.ShowVersionAndExit(version)
	}

	// 检查是否提供了图片路径（命令行模式）
	if config.CliCfg.ImagePath != "" {
		// 命令行模式：处理图片并退出
		fmt.Printf("正在处理图片: %s\n", config.CliCfg.ImagePath)

		// 读取图片数据
		data, filename, err := readImageData(config.CliCfg.ImagePath)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}

		// 处理图片
		url, err := image.ProcessImage(data, filename)
		if err != nil {
			fmt.Printf("图片处理失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("图片处理成功！\n")
		fmt.Printf("访问URL: %s\n", url)
		return
	}

	// 服务模式：启动Web服务器
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

	docs.SwaggerInfo.Version = version
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
