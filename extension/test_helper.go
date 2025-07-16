package extension

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func createTestImage(path string) error {
	return createTestImageWithSize(path, 128, 128)
}

func createTestImageWithSize(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a simple pattern with fewer colors for better compression
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Simple checkerboard pattern with just 2 colors
			if (x/4+y/4)%2 == 0 {
				img.Set(x, y, color.RGBA{220, 220, 220, 255}) // Light gray
			} else {
				img.Set(x, y, color.RGBA{180, 180, 180, 255}) // Slightly darker gray
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	return encoder.Encode(f, img)
}
