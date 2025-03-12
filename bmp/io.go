package bmp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// Represents a single pixel in the image (for 24-bit BMP files)
type Pixel struct {
	Blue  byte
	Green byte
	Red   byte
}

// Extracts pixel data from a BMP file
func ReadPixels(filename string, bmpHeader *BMPHeader, dibHeader *DIBHeader) ([]Pixel, error) {
	// Ensure the image is within reasonable size limits
	if dibHeader.Width > 65536 || dibHeader.Height > 65536 {
		return nil, fmt.Errorf("image is too large to process")
	}

	// Open the BMP file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file - %v", err)
	}
	defer file.Close()

	// Seek to the start of the pixel data
	_, err = file.Seek(int64(bmpHeader.DataOffset), 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking to pixel data - %v", err)
	}

	// Get image dimensions
	width, height := int(dibHeader.Width), int(dibHeader.Height)
	pixels := make([]Pixel, width*height)

	// Calculate row size and padding
	rowSize := ((width * 3) + 3) / 4 * 4 // BMP rows are padded to 4-byte alignment
	buf := make([]byte, rowSize)         // Buffer for reading full rows

	// Read pixel data (BMP stores pixels bottom-up)
	for y := height - 1; y >= 0; y-- {
		_, err := file.Read(buf) // Read the entire row into buffer
		if err != nil {
			return nil, fmt.Errorf("error reading pixel data: %v", err)
		}

		// Extract RGB values from buf and assign to pixels slice
		for x := 0; x < width; x++ {
			bufIndex := x * 3
			pixelIndex := y*width + x
			pixels[pixelIndex] = Pixel{
				Blue:  buf[bufIndex],
				Green: buf[bufIndex+1],
				Red:   buf[bufIndex+2],
			}
		}
	}

	return pixels, nil
}

// Writes the modified pixel data to an output BMP file
func WritePixels(filename string, bmpHeader *BMPHeader, dibHeader *DIBHeader, pixels []Pixel) error {
	// Create the output BMP file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file - %v", err)
	}
	defer file.Close()

	err = writeHeaders(file, *bmpHeader, *dibHeader)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	writer := bufio.NewWriter(file)
	width, height := int(dibHeader.Width), int(dibHeader.Height)
	rowSize := ((width * 3) + 3) &^ 3 // Align to 4-byte boundary
	padding := rowSize - (width * 3)
	paddingBytes := make([]byte, padding)

	rowBuffer := make([]byte, rowSize)

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			idx := y*width + x
			offset := x * 3
			rowBuffer[offset] = pixels[idx].Blue
			rowBuffer[offset+1] = pixels[idx].Green
			rowBuffer[offset+2] = pixels[idx].Red
		}
		copy(rowBuffer[width*3:], paddingBytes)

		_, err = writer.Write(rowBuffer)
		if err != nil {
			return fmt.Errorf("error writing pixel row: %v", err)
		}
	}

	return writer.Flush()
}

func writeHeaders(file *os.File, bmpHeader BMPHeader, dibHeader DIBHeader) error {
	// Write BMP Header (Only first 14 bytes)
	err := binary.Write(file, binary.LittleEndian, bmpHeader.Signature)
	if err != nil {
		return fmt.Errorf("error writing BMP Signature - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, bmpHeader.FileSize)
	if err != nil {
		return fmt.Errorf("error writing BMP FileSize - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, bmpHeader.Reserved)
	if err != nil {
		return fmt.Errorf("error writing BMP Reserved - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, bmpHeader.DataOffset)
	if err != nil {
		return fmt.Errorf("error writing BMP DataOffset - %v", err)
	}

	// Write DIB Header (Make sure to write only valid fields)
	err = binary.Write(file, binary.LittleEndian, dibHeader.DibHeaderSize)
	if err != nil {
		return fmt.Errorf("error writing DIB Header Size - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.Width)
	if err != nil {
		return fmt.Errorf("error writing DIB Width - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.Height)
	if err != nil {
		return fmt.Errorf("error writing DIB Height - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.Planes)
	if err != nil {
		return fmt.Errorf("error writing DIB Planes - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.BitCount)
	if err != nil {
		return fmt.Errorf("error writing DIB BitCount - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.Compression)
	if err != nil {
		return fmt.Errorf("error writing DIB Compression - %v", err)
	}

	// Handle ImageSize correctly for BMP v3
	if dibHeader.DibHeaderSize == 40 && dibHeader.Compression == 0 {
		dibHeader.ImageSize = 0 // BMP v3 does not require an explicit ImageSize
	}

	err = binary.Write(file, binary.LittleEndian, dibHeader.ImageSize)
	if err != nil {
		return fmt.Errorf("error writing DIB ImageSize - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.XPixelsPerM)
	if err != nil {
		return fmt.Errorf("error writing DIB XPixelsPerM - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.YPixelsPerM)
	if err != nil {
		return fmt.Errorf("error writing DIB YPixelsPerM - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.ColorsUsed)
	if err != nil {
		return fmt.Errorf("error writing DIB ColorsUsed - %v", err)
	}
	err = binary.Write(file, binary.LittleEndian, dibHeader.ColorsImp)
	if err != nil {
		return fmt.Errorf("error writing DIB ColorsImp - %v", err)
	}

	// Read BMP v4 fields (if present)
	if dibHeader.DibHeaderSize >= 108 {
		if err := binary.Write(file, binary.LittleEndian, dibHeader.RedMask); err != nil {
			return fmt.Errorf("error writing RedMask - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.GreenMask); err != nil {
			return fmt.Errorf("error writing GreenMask - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.BlueMask); err != nil {
			return fmt.Errorf("error writing BlueMask - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.AlphaMask); err != nil {
			return fmt.Errorf("error writing AlphaMask - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.ColorSpace); err != nil {
			return fmt.Errorf("error writing ColorSpace - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.Endpoints); err != nil {
			return fmt.Errorf("error writing ColorSpaceEndpoints - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.GammaRed); err != nil {
			return fmt.Errorf("error writing GammaRed - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.GammaGreen); err != nil {
			return fmt.Errorf("error writing GammaGreen - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.GammaBlue); err != nil {
			return fmt.Errorf("error writing GammaBlue - %v", err)
		}
	}

	// Write BMP v5 fields (if present)
	if dibHeader.DibHeaderSize == 124 {
		if err := binary.Write(file, binary.LittleEndian, dibHeader.Intent); err != nil {
			return fmt.Errorf("error writing Intent - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.ProfileData); err != nil {
			return fmt.Errorf("error writing ProfileData - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.ProfileSize); err != nil {
			return fmt.Errorf("error writing ProfileSize - %v", err)
		}
		if err := binary.Write(file, binary.LittleEndian, dibHeader.Reserved); err != nil {
			return fmt.Errorf("error writing Reserved - %v", err)
		}
	}

	return nil
}
