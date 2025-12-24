package image

import (
	"bytes"
	"easyimage_go/biz/response"
	"easyimage_go/utils"
	"easyimage_go/utils/config"
	"fmt"
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

// ProcessImage 处理图片（转换格式和存储）
// 参数：
//
//	fileData: 图片文件数据
//	filename: 原始文件名
//
// 返回：
//
//	存储后的文件URL
//	错误信息
func ProcessImage(fileData []byte, filename string) (string, error) {
	// 检测文件类型
	contentType := http.DetectContentType(fileData)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	// 初始化文件扩展名
	ext := filepath.Ext(filename)

	// 检查是否为HEIC格式
	isHEIC := false
	if contentType == "application/octet-stream" {
		// HEIC文件可能被识别为octet-stream，进一步检查文件扩展名
		if ext == ".heic" || ext == ".HEIC" {
			isHEIC = true
		}
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
			return "", convertErr
		}

		// 更新文件数据和扩展名
		fileData = jpegData
		contentType = "image/jpeg"
		ext = ".jpg"
	}

	if !allowedTypes[contentType] {
		return "", fmt.Errorf("不支持的文件类型: %s", contentType)
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
			return "", decodeErr
		}

		// 转换为WebP
		webpData, convertErr := utils.ConvertToWebP(img, config.Cfg.Image.WebPQuality)
		if convertErr != nil {
			return "", convertErr
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
		return "", err
	}

	finalFilePath := filepath.Join(imgDir, finalFileName)

	// 直接将内存中的数据写入最终目标文件
	if err := os.WriteFile(finalFilePath, finalFileData, 0644); err != nil {
		return "", err
	}

	// 生成返回URL
	url := config.Cfg.Server.Domain + config.Cfg.Image.Uri + "/" + year + "/" + month + "/" + day + "/" + finalFileName
	return url, nil
}

// UploadImage 上传图片
//
//	@Tags			图片
//	@Summary		上传图片
//	@Description	上传图片
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Param			file	formData	file	true	"图片文件"
//	@Success		200		{object}	UploadResp
//	@Security		ApiKeyAuth
//	@router			/api/image/upload [PUT]
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

	// 读取文件内容
	fileData, err := io.ReadAll(fileHeader)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取文件失败: "+err.Error())
		return
	}

	// 处理图片
	url, err := ProcessImage(fileData, req.File.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp.Code = response.Code_Success
	resp.Msg = "上传图片成功"
	resp.Url = url

	c.JSON(http.StatusOK, resp)
}

// UploadImageForm 上传图片
//
//	@Tags			图片
//	@Summary		表单上传图片
//	@Description	使用multipart/form-data格式上传图片，支持Apple自动化调用
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Param			image	formData	file	true	"图片文件（支持curl -F \"image=@/path/to/file\"）"
//	@Success		200		{object}	UploadResp
//	@Security		ApiKeyAuth
//	@router			/api/image/upload [POST]
func UploadImageForm(c *gin.Context) {
	resp := new(UploadResp)

	// 先尝试获取image字段（支持curl -F "image=@/path/to/file"）
	file, err := c.FormFile("image")
	if err != nil {
		// 如果没有image字段，尝试获取file字段（兼容原有的multipart/form-data）
		file, err = c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "获取图片文件失败: "+err.Error())
			return
		}
	}

	// 打开文件
	fileHeader, err := file.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "打开文件失败: "+err.Error())
		return
	}
	defer fileHeader.Close()

	// 读取文件内容
	fileData, err := io.ReadAll(fileHeader)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取文件失败: "+err.Error())
		return
	}

	// 处理图片
	url, err := ProcessImage(fileData, file.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp.Code = response.Code_Success
	resp.Msg = "上传图片成功"
	resp.Url = url

	c.JSON(http.StatusOK, resp)
}
