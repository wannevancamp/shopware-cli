package extension

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

func createTestImage(path string) error {
	return createTestImageWithSize(path, 128, 128)
}

func createTestImageWithSize(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	return png.Encode(f, img)
}
