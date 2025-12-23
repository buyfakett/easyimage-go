package image

import (
	"easyimage_go/biz/response"
	"easyimage_go/utils/config"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DeleteReq struct {
	Path string `form:"path"`
}

// DeleteFile 删除文件
//
//	@Tags			图片
//	@Summary		删除文件
//	@Description	删除相对图片目录下的文件
//	@Accept			application/json
//	@Produce		application/json
//	@Param			path	query		string	true	"文件路径"
//	@Success		200		{object}	response.CommonResp
//	@Router			/api/image/delete [delete]
func DeleteFile(c *gin.Context) {
	var req DeleteReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, &response.CommonResp{
			Code: response.Code_Err,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}
	resp := new(response.CommonResp)

	real := filepath.Join("."+config.Cfg.Image.Uri, req.Path)

	if err := os.Remove(real); err != nil {
		c.JSON(http.StatusOK, &response.CommonResp{
			Code: response.Code_Err,
			Msg:  err.Error(),
		})
		return
	}

	resp.Code = response.Code_Success
	resp.Msg = "文件" + req.Path + "删除成功"
	c.JSON(http.StatusOK, resp)
}
