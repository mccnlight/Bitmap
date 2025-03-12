package bmp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ApplyCrop(pixels []Pixel, width, height int, options string) ([]Pixel, int, int, error) {
	optionSlice := strings.Split(options, "-")

	optionSliceInt := []int{}

	for _, option := range optionSlice {
		num, err := strconv.Atoi(option)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("%v", err)
		}

		optionSliceInt = append(optionSliceInt, num)
	}

	offsetX, offsetY, cropWidth, cropHeight, err := parseAndValidateOptions(optionSliceInt, height, width)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%v", err)
	}

	// Create new slice for cropped pixels
	croppedPixels := make([]Pixel, cropWidth*cropHeight)

	// Copy the cropped part
	for y := 0; y < cropHeight; y++ {
		rowStart := (offsetY+y)*width + offsetX
		rowEnd := rowStart + cropWidth
		if rowEnd > (offsetY+y+1)*width { // Prevents out-of-bounds panic
			rowEnd = (offsetY + y + 1) * width
		}
		copy(croppedPixels[y*cropWidth:(y+1)*cropWidth], pixels[rowStart:rowEnd])
	}

	return croppedPixels, cropWidth, cropHeight, nil
}

func parseAndValidateOptions(options []int, height, width int) (int, int, int, int, error) {
	var offsetX, offsetY, cropWidth, cropHeight int

	if len(options) != 2 && len(options) != 4 {
		return 0, 0, 0, 0, errors.New("invalid number of options")
	}

	if len(options) == 2 {
		offsetX, offsetY = options[0], options[1]
		cropWidth, cropHeight = width-offsetX, height-offsetY
	} else {
		offsetX, offsetY, cropWidth, cropHeight = options[0], options[1], options[2], options[3]
	}

	// Validate crop dimensions
	if offsetX < 0 || offsetY < 0 || offsetX >= width || offsetY >= height {
		return 0, 0, 0, 0, errors.New("crop offset is out of bounds")
	}

	// Ensure crop area is within bounds
	if offsetX+cropWidth > width || offsetY+cropHeight > height {
		return 0, 0, 0, 0, errors.New("crop dimensions exceed image bounds")
	}

	return offsetX, offsetY, cropWidth, cropHeight, nil
}
