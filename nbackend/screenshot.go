package backend

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"log"
	"os"
)

// CaptureScreenshot captures all active displays and returns temp file path
func CaptureScreenshot() (string, error) {
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return "", fmt.Errorf("no active displays found")
	}

	var images []*image.RGBA
	var bounds []image.Rectangle

	// Capture all displays
	for i := 0; i < n; i++ {
		img, err := screenshot.CaptureDisplay(i)
		if err != nil {
			log.Printf("Failed to capture display %d: %v", i, err)
			continue
		}
		images = append(images, img)
		bounds = append(bounds, screenshot.GetDisplayBounds(i))
	}

	if len(images) == 0 {
		return "", fmt.Errorf("failed to capture any displays")
	}

	// Create combined image
	combined, err := combineScreenshots(images, bounds)
	if err != nil {
		return "", err
	}

	// Save to temp file
	tmpFile, err := os.CreateTemp("", "gpt-snapper-*.png")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if err := png.Encode(tmpFile, combined); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// combineScreenshots stitches multiple displays into a single image
func combineScreenshots(images []*image.RGBA, bounds []image.Rectangle) (*image.RGBA, error) {
	minX, minY, maxX, maxY := getCombinedBounds(bounds)
	width := maxX - minX
	height := maxY - minY

	combined := image.NewRGBA(image.Rect(0, 0, width, height))

	for i, img := range images {
		b := bounds[i]
		xOffset := b.Min.X - minX
		yOffset := b.Min.Y - minY

		// Draw each pixel into combined image
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				combined.Set(x+xOffset, y+yOffset, img.At(x, y))
			}
		}
	}

	return combined, nil
}

// getCombinedBounds calculates total area covering all displays
func getCombinedBounds(bounds []image.Rectangle) (minX, minY, maxX, maxY int) {
	minX = bounds[0].Min.X
	minY = bounds[0].Min.Y
	maxX = bounds[0].Max.X
	maxY = bounds[0].Max.Y

	for _, b := range bounds[1:] {
		if b.Min.X < minX {
			minX = b.Min.X
		}
		if b.Min.Y < minY {
			minY = b.Min.Y
		}
		if b.Max.X > maxX {
			maxX = b.Max.X
		}
		if b.Max.Y > maxY {
			maxY = b.Max.Y
		}
	}
	return
}
