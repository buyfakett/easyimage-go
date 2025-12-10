package user

import (
	"easyimage_go/biz/response"
	"easyimage_go/utils/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestTokenReq struct {
	Token     string `json:"token" binding:"required,min=1,max=255"`
	CaptchaID string `json:"captcha_id" binding:"required,min=1,max=255"`
	Captcha   string `json:"captcha" binding:"required,min=1,max=10"`
}

// TestToken 测试token权限
// @Tags 测试
// @Summary 测试token权限
// @Description 测试token权限
// @Accept application/json
// @Produce application/json
// @Param req body TestTokenReq true "测试token权限"
// @Success 200 {object} response.CommonResp
// @Router /api/user/test_token [post]
func TestToken(c *gin.Context) {
	req := new(TestTokenReq)
	if err := c.ShouldBind(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// 验证验证码
	if !captchaStore.Verify(req.CaptchaID, req.Captcha, true) {
		c.JSON(http.StatusOK, &response.CommonResp{
			Code: response.Code_CaptchaErr,
			Msg:  "验证码错误或已过期",
		})
		return
	}
	if req.Token != config.Cfg.Server.Token {
		c.JSON(http.StatusOK, &response.CommonResp{
			Code: response.Code_Err,
			Msg:  "token错误",
		})
		return
	}
	c.JSON(http.StatusOK, &response.CommonResp{
		Code: response.Code_Success,
		Msg:  "权限验证成功",
	})
	return
}
