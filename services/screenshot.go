package services

import (
	"fmt"
	"github.com/jasonlovesdoggo/displayindex"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

type ScreenshotService struct {
	tempDir string
}

func NewScreenshotService(tempDir string) *ScreenshotService {
	return &ScreenshotService{
		tempDir: tempDir,
	}
}

func (s *ScreenshotService) CaptureCurrentScreen() (string, error) {
	index, err := displayindex.CurrentDisplayIndex()
	if err != nil {
		return "", fmt.Errorf("failed to get current display: %w", err)
	}

	img, err := screenshot.CaptureDisplay(index)
	if err != nil {
		return "", fmt.Errorf("failed to capture screen: %w", err)
	}

	return s.saveImage(img)
}

func (s *ScreenshotService) saveImage(img *image.RGBA) (string, error) {
	filename := filepath.Join(s.tempDir, fmt.Sprintf("screenshot_%d.png", time.Now().UnixNano()))

	err := os.MkdirAll(s.tempDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}

func (s *ScreenshotService) CleanupTempFiles() error {
	return filepath.Walk(s.tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && time.Since(info.ModTime()) > 24*time.Hour {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}
