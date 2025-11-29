//go:build !cgo
// +build !cgo

package utils

import (
	"fmt"
	"image"
)

// convertToWebPImpl converts an image to WebP format
// This is a stub implementation for builds that don't support WebP
func convertToWebPImpl(img image.Image, quality int) ([]byte, error) {
	// Without CGO, we can't do WebP conversion, so we return an error
	// The caller should handle this and potentially use a different format
	return nil, fmt.Errorf("WebP conversion not supported in this build (CGO disabled)")
}
