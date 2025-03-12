package bmp

import "fmt"

// Applies various filters like blue, red, green, grayscale, negative, pixelate or blur
func ApplyFilter(pixels []Pixel, width, height int, filterType string) ([]Pixel, error) {
	switch filterType {
	case "pixelate":
		return applyPixelation(pixels, width, height, 20), nil
	case "blur":
		return applyBlur(pixels, width, height, 25), nil
	case "blue":
		for i := range pixels {
			pixels[i].Red = 0
			pixels[i].Green = 0
		}
	case "red":
		for i := range pixels {
			pixels[i].Green = 0
			pixels[i].Blue = 0
		}
	case "green":
		for i := range pixels {
			pixels[i].Red = 0
			pixels[i].Blue = 0
		}
	case "grayscale":
		for i := range pixels {
			gray := uint8(0.299*float64(pixels[i].Red) + 0.587*float64(pixels[i].Green) + 0.114*float64(pixels[i].Blue))
			pixels[i].Red, pixels[i].Green, pixels[i].Blue = gray, gray, gray
		}
	case "negative":
		for i := range pixels {
			pixels[i].Red = 255 - pixels[i].Red
			pixels[i].Green = 255 - pixels[i].Green
			pixels[i].Blue = 255 - pixels[i].Blue
		}
	default:
		return nil, fmt.Errorf("invalid filter type - '%s'", filterType)
	}
	return pixels, nil
}

// Applies pixelation to the image
func applyPixelation(pixels []Pixel, width, height, blockSize int) []Pixel {
	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			var sumR, sumG, sumB, count int

			// Collect block colors
			for dy := 0; dy < blockSize && (y+dy) < height; dy++ {
				for dx := 0; dx < blockSize && (x+dx) < width; dx++ {
					idx := (y+dy)*width + (x + dx)
					sumR += int(pixels[idx].Red)
					sumG += int(pixels[idx].Green)
					sumB += int(pixels[idx].Blue)
					count++
				}
			}

			// Calculate the average color of the block
			avgR := uint8(sumR / count)
			avgG := uint8(sumG / count)
			avgB := uint8(sumB / count)

			// Apply the averaged color to all pixels in the block
			for dy := 0; dy < blockSize && (y+dy) < height; dy++ {
				for dx := 0; dx < blockSize && (x+dx) < width; dx++ {
					idx := (y+dy)*width + (x + dx)
					pixels[idx] = Pixel{Red: avgR, Green: avgG, Blue: avgB}
				}
			}
		}
	}
	return pixels
}

// Applies blur to the image
func applyBlur(pixels []Pixel, width, height int, kernelSize int) []Pixel {
	blurredPixels := make([]Pixel, len(pixels))
	radius := kernelSize / 2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sumR, sumG, sumB, count int

			// Iterate over the surrounding pixels in the kernel
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					neighborX := x + dx
					neighborY := y + dy

					// Ensure the neighboring pixel is within bounds
					if neighborX >= 0 && neighborX < width && neighborY >= 0 && neighborY < height {
						index := neighborY*width + neighborX
						pixel := pixels[index]

						// Sum the RGB values of the surrounding pixels
						sumR += int(pixel.Red)
						sumG += int(pixel.Green)
						sumB += int(pixel.Blue)
						count++
					}
				}
			}

			// Apply the averaged color to the current pixel
			newIndex := y*width + x
			blurredPixels[newIndex] = Pixel{
				Red:   byte(sumR / count),
				Green: byte(sumG / count),
				Blue:  byte(sumB / count),
			}
		}
	}

	return blurredPixels
}
