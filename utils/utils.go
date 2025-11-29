package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandomFilename 生成随机文件名（时间戳+随机数）
func GenerateRandomFilename(ext string) string {
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := fmt.Sprintf("%x", randomBytes)
	return randomStr + ext
}

// GetExtensionByType 根据Content-Type获取文件扩展名
func GetExtensionByType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
