package image

import (
	"bytes"
	"easyimage_go/biz/response"
	"easyimage_go/utils"
	"easyimage_go/utils/config"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadReq struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadResp struct {
	response.CommonResp
	Url string `json:"url"`
}

// UploadImage 上传图片
// @Tags 图片
// @Summary 上传图片
// @Description 上传图片
// @Accept multipart/form-data
// @Produce multipart/form-data
// @Param file formData file true "图片文件"
// @Success 200 {object} UploadResp
// @Security ApiKeyAuth
// @router /api/image/upload [PUT]
func UploadImage(c *gin.Context) {
	req := new(UploadReq)
	if err := c.ShouldBind(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp := new(UploadResp)

	// 验证文件类型
	fileHeader, err := req.File.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "打开文件失败: "+err.Error())
		return
	}
	defer fileHeader.Close()

	// 检测文件类型
	buffer := make([]byte, 512)
	if _, err := fileHeader.Read(buffer); err != nil {
		c.String(http.StatusBadRequest, "读取文件失败: "+err.Error())
		return
	}

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	// 重置文件指针
	if _, err := fileHeader.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "文件处理失败: "+err.Error())
		return
	}

	// 初始化文件扩展名
	ext := filepath.Ext(req.File.Filename)

	// 检查是否为HEIC格式
	isHEIC := false
	if contentType == "application/octet-stream" {
		// HEIC文件可能被识别为octet-stream，进一步检查文件扩展名
		if ext == ".heic" || ext == ".HEIC" {
			isHEIC = true
		}
	}

	// 读取文件内容用于可能的HEIC转换
	fileData, err := io.ReadAll(fileHeader)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取文件失败: "+err.Error())
		return
	}

	// 检查文件签名是否为HEIC
	checkLen := 64
	if len(fileData) < 64 {
		checkLen = len(fileData)
	}
	if utils.IsHEIC(fileData[:checkLen]) {
		isHEIC = true
	}

	// 如果是HEIC格式，尝试转换为JPEG
	if isHEIC {
		jpegData, convertErr := utils.ConvertHEICtoJPEG(fileData)
		if convertErr != nil {
			c.String(http.StatusBadRequest, "HEIC格式转换失败，请转换为JPEG或PNG格式后再上传: "+convertErr.Error())
			return
		}

		// 更新文件数据和扩展名
		fileData = jpegData
		contentType = "image/jpeg"
		ext = ".jpg"
	}

	if !allowedTypes[contentType] {
		c.String(http.StatusBadRequest, "不支持的文件类型: "+contentType)
		return
	}

	// 获取文件扩展名
	ext = utils.GetExtensionByType(contentType)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	var finalFileName string
	var finalFileData []byte

	if ext != ".webp" {
		// 解码图像
		img, _, decodeErr := image.Decode(bytes.NewReader(fileData))
		if decodeErr != nil {
			c.String(http.StatusInternalServerError, "解码图像失败: "+decodeErr.Error())
			return
		}

		// 转换为WebP
		webpData, convertErr := utils.ConvertToWebP(img, config.Cfg.Image.WebPQuality)
		if convertErr != nil {
			c.String(http.StatusInternalServerError, "转换WebP失败: "+convertErr.Error())
			return
		}

		finalFileName = utils.GenerateRandomFilename(".webp")
		finalFileData = webpData
	} else {
		// 不转换，使用原始文件
		finalFileName = utils.GenerateRandomFilename(ext)
		finalFileData = fileData
	}

	// 创建 i/年/月/日 目录结构
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	// 创建年/月/日目录结构
	imgDir := "." + config.Cfg.Image.Uri
	imgDir = filepath.Join(imgDir, year, month, day)
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		c.String(http.StatusInternalServerError, "创建目标目录失败: "+err.Error())
		return
	}

	finalFilePath := filepath.Join(imgDir, finalFileName)

	// 直接将内存中的数据写入最终目标文件
	if err := os.WriteFile(finalFilePath, finalFileData, 0644); err != nil {
		c.String(http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}

	resp.Code = response.Code_Success
	resp.Msg = "上传图片成功"
	resp.Url = config.Cfg.Server.Domain + config.Cfg.Image.Uri + "/" + year + "/" + month + "/" + day + "/" + finalFileName

	c.JSON(http.StatusOK, resp)
}
