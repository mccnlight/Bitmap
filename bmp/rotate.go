package bmp

import (
	"fmt"
)

// Applies rotation to the image. Supports rotatiob by 90, 180, and 270 degrees
func ApplyRotate(pixels []Pixel, width int, height int, angle int) ([]Pixel, int, int, error) {
	var newWidth, newHeight int
	var rotatedPixels []Pixel

	switch angle {
	case 90:
		newWidth, newHeight = height, width
		rotatedPixels = make([]Pixel, newWidth*newHeight)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				newX := height - y - 1
				newY := x
				rotatedPixels[newY*newWidth+newX] = pixels[y*width+x]
			}
		}
	case 180:
		newWidth, newHeight = width, height
		rotatedPixels = make([]Pixel, newWidth*newHeight)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				newX := width - x - 1
				newY := height - y - 1
				rotatedPixels[newY*newWidth+newX] = pixels[y*width+x]
			}
		}
	case 270:
		newWidth, newHeight = height, width
		rotatedPixels = make([]Pixel, newWidth*newHeight)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				newX := y
				newY := width - x - 1
				rotatedPixels[newY*newWidth+newX] = pixels[y*width+x]
			}
		}
	default:
		return nil, 0, 0, fmt.Errorf("'%v' is not a valid angle value", angle)
	}

	return rotatedPixels, newWidth, newHeight, nil
}

// Converts the given angle to an int value or returns an error in case of wrong angle value
func ParseRotationValue(rotation string) (int, error) {
	switch rotation {
	case "right", "90", "-270":
		return 90, nil
	case "left", "270", "-90":
		return 270, nil
	case "180", "-180":
		return 180, nil
	default:
		return 0, fmt.Errorf("'%s' is not a valid angle value", rotation)
	}
}
