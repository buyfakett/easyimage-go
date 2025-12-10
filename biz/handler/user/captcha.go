package user

import (
	"easyimage_go/biz/response"
	"easyimage_go/utils/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// 验证码驱动
var captchaDriver = base64Captcha.NewDriverDigit(
	60,  // 高度
	240, // 宽度
	4,   // 长度
	0.7, // 最大倾斜度
	80,  // 点数
)

// 验证码存储 - 设置5分钟过期
var captchaStore = base64Captcha.NewMemoryStore(10240, time.Duration(config.Cfg.Server.CaptchaExpireTime)*time.Minute)

// CaptchaResp 验证码响应
type CaptchaResp struct {
	Code response.Code `json:"code"`
	Msg  string        `json:"msg"`
	Data *CaptchaData  `json:"data"`
}

// CaptchaData 验证码数据
type CaptchaData struct {
	ID          string `json:"id"`
	Base64Image string `json:"base64_image"`
}

// GenerateCaptcha 生成验证码
// @Tags 用户
// @Summary 生成验证码
// @Description 生成登录验证码
// @Produce application/json
// @Success 200 {object} CaptchaResp
// @router /api/user/captcha [GET]
func GenerateCaptcha(c *gin.Context) {
	// 创建验证码
	captcha := base64Captcha.NewCaptcha(captchaDriver, captchaStore)
	id, base64Image, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &CaptchaResp{
			Code: response.Code_Err,
			Msg:  "生成验证码失败",
		})
		return
	}

	c.JSON(http.StatusOK, &CaptchaResp{
		Code: response.Code_Success,
		Msg:  "生成验证码成功",
		Data: &CaptchaData{
			ID:          id,
			Base64Image: base64Image,
		},
	})
}
