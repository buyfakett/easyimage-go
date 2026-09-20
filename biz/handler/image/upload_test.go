package image

import (
	"bytes"
	"easyimage_go/utils/config"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessImagePreservesPNG(t *testing.T) {
	var source bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&source, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}

	tempDir := t.TempDir()
	oldWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWorkingDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	oldConfig := config.Cfg
	config.Cfg.Image.Uri = "/i"
	config.Cfg.Server.Domain = "http://example.test"
	t.Cleanup(func() {
		config.Cfg = oldConfig
	})

	url, err := ProcessImage(source.Bytes(), "sample.png")
	if err != nil {
		t.Fatalf("process PNG: %v", err)
	}
	if !strings.HasSuffix(url, ".png") {
		t.Fatalf("expected PNG URL, got %q", url)
	}

	relativePath := strings.TrimPrefix(url, config.Cfg.Server.Domain)
	stored, err := os.ReadFile(filepath.Join(tempDir, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatalf("read stored PNG: %v", err)
	}
	if !bytes.Equal(stored, source.Bytes()) {
		t.Fatal("stored PNG differs from source")
	}
}
