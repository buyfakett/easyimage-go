//go:build cgo
// +build cgo

package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
)

// convertToWebPImpl converts an image to WebP format
// This implementation uses CGO and external C libraries
func convertToWebPImpl(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
	if err != nil {
		return nil, fmt.Errorf("failed to encode WebP: %v", err)
	}
	return buf.Bytes(), nil
}

// encodeAsJPEG encodes an image as JPEG
func encodeAsJPEG(img image.Image, quality int) ([]byte, string, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/jpeg", nil
}

// encodeAsPNG encodes an image as PNG
func encodeAsPNG(img image.Image) ([]byte, string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/png", nil
}
