package image

import (
	"easyimage_go/biz/response"
	"easyimage_go/utils/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ListFilesReq struct {
	Dir string `form:"dir"` // 查询参数 ?dir=
}

type FileItem struct {
	Name string `json:"name"` // 文件或目录名
	Path string `json:"path"` // 相对路径
	Url  string `json:"url"`  // 文件访问 URL，仅 file 有
	Type string `json:"type"` // file | dir
}

type ListFilesResp struct {
	Code response.Code `json:"code"`
	Msg  string        `json:"msg"`
	Data []*FileItem   `json:"data"`
}

// ListFiles 获取文件/目录列表
// @Tags 图片
// @Summary 获取文件列表
// @Description 获取指定目录下的文件与子目录（不递归）
// @Accept application/json
// @Produce application/json
// @Param dir query string false "相对 图片 的目录，不传表示根目录"
// @Success 200 {object} ListFilesResp
// @Router /api/image/list [get]
func ListFiles(c *gin.Context) {
	var req ListFilesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, ListFilesResp{
			Code: response.Code_Err,
			Msg:  "参数错误: " + err.Error(),
			Data: nil,
		})
		return
	}

	realPath := filepath.Join("."+config.Cfg.Image.Uri, req.Dir)
	entries, err := os.ReadDir(realPath)
	if err != nil {
		c.JSON(http.StatusOK, ListFilesResp{
			Code: response.Code_Err,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	var list []*FileItem

	for _, entry := range entries {
		name := entry.Name()

		// 排除以.开头的文件/目录
		if strings.HasPrefix(name, ".") {
			continue
		}

		rel := filepath.Join(req.Dir, name)
		rel = filepath.ToSlash(rel)

		if entry.IsDir() {
			list = append(list, &FileItem{
				Name: name,
				Path: rel,
				Type: "dir",
			})
		} else {
			list = append(list, &FileItem{
				Name: name,
				Path: rel,
				Url:  config.Cfg.Server.Domain + config.Cfg.Image.Uri + "/" + rel,
				Type: "file",
			})
		}
	}

	c.JSON(http.StatusOK, ListFilesResp{
		Code: response.Code_Success,
		Msg:  "success",
		Data: list,
	})
}
