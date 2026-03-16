package densecode

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func (c *Configuration) RenderPNG(filename string) error {
	matrix := c.ToMatrix()
	imgSize := c.size * c.ModuleSize
	palette := colorPaletteForBits(c.bitsPerModule())

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			colorIdx := matrix[i][j]
			if colorIdx >= len(palette) {
				colorIdx = 0
			}
			color := palette[colorIdx]

			for px := 0; px < c.ModuleSize; px++ {
				for py := 0; py < c.ModuleSize; py++ {
					img.Set(j*c.ModuleSize+px, i*c.ModuleSize+py, color)
				}
			}
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	return png.Encode(file, img)
}

func (c *Configuration) RenderTerminal() {
	matrix := c.ToMatrix()
	palette := colorPaletteForBits(c.bitsPerModule())
	reset := "\033[0m"

	for i := 0; i < c.size; i++ {
		for j := 0; j < c.size; j++ {
			colorIdx := matrix[i][j]
			if colorIdx >= len(palette) {
				colorIdx = 0
			}
			c := palette[colorIdx]
			fmt.Printf("\033[48;2;%d;%d;%dm  %s", c.R, c.G, c.B, reset)
		}
		fmt.Println()
	}
	fmt.Println()
}
