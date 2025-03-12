package bmp

import (
	"fmt"
)

// Applies horizontal or vertical mirroring
func ApplyMirror(pixels []Pixel, width, height int, mode string) ([]Pixel, error) {
	mirror := make([]Pixel, len(pixels))

	switch mode {
	case "horizontal", "h", "horizontally", "hor":
		// Iterate over each pixel in the bitmap
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				srcIndex := y*width + x               // Original pixel position
				dstIndex := y*width + (width - 1 - x) // Mirrored position in the same row

				mirror[dstIndex] = pixels[srcIndex]
			}
		}
	case "vertical", "v", "vertically", "ver":
		// Iterate over each pixel in the bitmap
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				srcIndex := y*width + x            // Original pixel position
				dstIndex := (height-1-y)*width + x // Mirrored position in the opposite row

				mirror[dstIndex] = pixels[srcIndex]
			}
		}
	default:
		return nil, fmt.Errorf("invalid mirror mode - '%s'", mode)
	}
	return mirror, nil
}
