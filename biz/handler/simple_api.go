package handler

import (
	"easyimage_go/utils/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping 测试网络接口
// @Tags 测试
// @Summary 测试网络接口
// @Description 测试网络接口
// @Accept application/json
// @Produce application/json
// @Router /api/ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "pong",
	})
}

// ServerInfo 服务信息
// @Tags 测试
// @Summary 服务信息
// @Description 服务信息
// @Accept application/json
// @Produce application/json
// @Router /api/server_info [get]
func ServerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"name":    config.Cfg.Server.Name,
			"version": config.Cfg.Server.Version,
		},
	})
}

// TestToken 测试token权限
// @Tags 测试
// @Summary 测试token权限
// @Description 测试token权限
// @Accept application/json
// @Produce application/json
// @Router /api/test_token [get]
func TestToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "权限验证成功",
	})
}
