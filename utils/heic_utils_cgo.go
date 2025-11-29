//go:build cgo
// +build cgo

package utils

import (
	"bytes"
	"fmt"
	"image/jpeg"

	"github.com/jdeng/goheif"
)

// ConvertHEICtoJPEG converts HEIC image data to JPEG format
// This implementation uses CGO and external C libraries
func ConvertHEICtoJPEG(data []byte) ([]byte, error) {
	// Decode HEIC image
	img, err := goheif.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode HEIC image: %w", err)
	}

	// Encode as JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JPEG image: %w", err)
	}

	return buf.Bytes(), nil
}

// IsHEIC checks if the given data is HEIC format
// This implementation checks the file signature
func IsHEIC(data []byte) bool {
	// Basic check for HEIC signature
	if len(data) < 12 {
		return false
	}

	// Check for 'ftyp' box
	if string(data[4:8]) == "ftyp" {
		// Check for HEIC major brand
		if string(data[8:12]) == "heic" || string(data[8:12]) == "heix" {
			return true
		}

		// Check for compatible brands (at position 16+)
		for i := 16; i < len(data)-4 && i < 64; i += 4 {
			if string(data[i:i+4]) == "heic" || string(data[i:i+4]) == "heix" {
				return true
			}
		}
	}

	return false
}
