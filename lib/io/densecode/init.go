package densecode

import "image/color"

var ColorPalette = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 255, 0, 255},     // Green
	{0, 0, 255, 255},     // Blue
	{255, 255, 0, 255},   // Yellow
	{255, 0, 255, 255},   // Magenta
	{0, 255, 255, 255},   // Cyan
}

var colorPalette1Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
}

var colorPalette2Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 0, 255, 255},     // Blue
}

var colorPalette4Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 255, 0, 255},     // Lime
	{0, 0, 255, 255},     // Blue
	{255, 255, 0, 255},   // Yellow
	{255, 0, 255, 255},   // Magenta
	{0, 255, 255, 255},   // Cyan
	{128, 0, 0, 255},     // Maroon
	{0, 128, 0, 255},     // Green
	{0, 0, 128, 255},     // Navy
	{128, 128, 0, 255},   // Olive
	{128, 0, 128, 255},   // Purple
	{0, 128, 128, 255},   // Teal
	{192, 192, 192, 255}, // Silver
	{128, 128, 128, 255}, // Gray
}

func colorPaletteForBits(bits int) []color.RGBA {
	switch bits {
	case 1:
		return colorPalette1Bit
	case 2:
		return colorPalette2Bit
	case 4:
		return colorPalette4Bit
	default:
		return ColorPalette
	}
}
