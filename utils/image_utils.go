package utils

import (
	"image"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

// ConvertToWebP 将图像转换为WebP格式
// 实际实现在webp_utils.go或webp_utils_cgo.go中，取决于是否启用了CGO
func ConvertToWebP(img image.Image, quality int) ([]byte, error) {
	// 实际实现会根据构建标签自动选择CGO或非CGO版本
	// 这里只是一个占位符，实际实现在webp_utils*.go文件中
	return convertToWebPImpl(img, quality)
}
