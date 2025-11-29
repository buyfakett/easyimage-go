//go:build !cgo
// +build !cgo

package utils

import (
	"fmt"
)

// ConvertHEICtoJPEG converts HEIC image data to JPEG format
// This is a stub implementation for builds that don't support HEIC
func ConvertHEICtoJPEG(data []byte) ([]byte, error) {
	return nil, fmt.Errorf("HEIC format not supported in this build")
}

// IsHEIC checks if the given data is HEIC format
// This is a stub implementation for builds that don't support HEIC
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
