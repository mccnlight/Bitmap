package bmp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Represents BMP header structure (first 14 bytes)
type BMPHeader struct {
	Signature  [2]byte // "BM"
	FileSize   uint32  // File size in bytes
	Reserved   uint32  // Reserved (always 0)
	DataOffset uint32  // Offset to image data
}

// Represents DIBHeader structure. It supports BMP v3 (40-byte), BMP v4 (108-byte), and BMP v5 (124-byte)
type DIBHeader struct {
	DibHeaderSize uint32 // Size of the DIB header (40, 108, or 124 bytes)
	Width         int32  // Image width in pixels
	Height        int32  // Image height in pixels
	Planes        uint16 // Number of color planes (must be 1)
	BitCount      uint16 // Bits per pixel (e.g., 24 for true color)
	Compression   uint32 // Compression type (0 for uncompressed)
	ImageSize     uint32 // Size of raw image data (can be 0 for uncompressed)
	XPixelsPerM   int32  // Horizontal resolution (pixels per meter)
	YPixelsPerM   int32  // Vertical resolution (pixels per meter)
	ColorsUsed    uint32 // Number of colors in the palette (0 for true color)
	ColorsImp     uint32 // Number of important colors (0 means all)

	// BMP v4+ Fields (Present if DibHeaderSize >= 108)
	RedMask    uint32   // Bit mask for the red channel
	GreenMask  uint32   // Bit mask for the green channel
	BlueMask   uint32   // Bit mask for the blue channel
	AlphaMask  uint32   // Bit mask for the alpha channel
	ColorSpace uint32   // Color space type
	Endpoints  [9]int32 // CIE XYZ color space endpoints
	GammaRed   uint32   // Gamma correction for red channel
	GammaGreen uint32   // Gamma correction for green channel
	GammaBlue  uint32   // Gamma correction for blue channel

	// BMP v5 Fields (Present if DibHeaderSize == 124)
	Intent      uint32 // Rendering intent (e.g., perceptual, colorimetric)
	ProfileData uint32 // Offset to ICC color profile data
	ProfileSize uint32 // Size of the ICC profile
	Reserved    uint32 // Always 0
}

// Represents a command-line option that consists of name and its value
type Option struct {
	Name  string // The option name (e.g., "--mirror", "--filter", "--rotate", etc)
	Value string // The associated value (e.g., "horizontal", "90", "negative", etc)
}

// Reads the BMP and DIB headers from a file
func ReadHeaders(filename string) (*BMPHeader, *DIBHeader, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file - %v", err)
	}
	defer file.Close()

	// Check if file is at least large enough to contain a BMP header (14 bytes) and a DIB header (40 bytes)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting file info - %v", err)
	}
	if fileInfo.Size() < 54 { // 14 bytes (BMP header) + 40 bytes (DIB header)
		return nil, nil, errors.New("not a valid BMP file")
	}

	bmpHeader, dibHeader, err := readHeaders(file)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	// Check the validity of the provided file
	if err := validateFile(*bmpHeader, *dibHeader, filename, fileInfo); err != nil {
		return nil, nil, err
	}

	return bmpHeader, dibHeader, nil
}

func validateFile(bmpHeader BMPHeader, dibHeader DIBHeader, filename string, fileInfo fs.FileInfo) error {
	// Ensure that it is a valid BMP file
	if string(bmpHeader.Signature[:]) != "BM" {
		return fmt.Errorf("%s is not a valid BMP file", filename)
	}

	// Ensure that the bit count is 24 (since we only support uncompressed 24-bit BMP files)
	if dibHeader.BitCount != 24 {
		return fmt.Errorf("%s is not a valid 24-bit BMP file (BitCount = %d)", filename, dibHeader.BitCount)
	}

	// Ensure compression is set to 0 (uncompressed)
	if dibHeader.Compression != 0 {
		return fmt.Errorf("%s is a compressed BMP file (Compression = %d), which is not supported", filename, dibHeader.Compression)
	}

	// Validate the correctness of data offset
	expectedDataOffset := int64(14 + dibHeader.DibHeaderSize)
	if int64(bmpHeader.DataOffset) != expectedDataOffset {
		return fmt.Errorf("unexpected pixel data offset: got %d, expected %d", bmpHeader.DataOffset, expectedDataOffset)
	}

	// Validate file size consistency
	expectedSize := int64(bmpHeader.DataOffset) + int64(dibHeader.Width)*int64(dibHeader.Height)*3
	if fileInfo.Size() < expectedSize {
		return fmt.Errorf("%s is corrupted or incomplete (file size too small)", filename)
	}

	return nil
}

// Prints the BMP and DIB header information
func PrintHeader(bmp *BMPHeader, dib *DIBHeader) {
	fmt.Println("BMP Header:")
	fmt.Printf("- Signature %s\n", string(bmp.Signature[:]))
	fmt.Printf("- FileSizeInBytes %d\n", bmp.FileSize)
	fmt.Printf("- HeaderSize %d\n", bmp.DataOffset)

	fmt.Println("DIB Header:")
	fmt.Printf("- DibHeaderSize %d\n", dib.DibHeaderSize)
	fmt.Printf("- WidthInPixels %d\n", dib.Width)
	fmt.Printf("- HeightInPixels %d\n", dib.Height)
	fmt.Printf("- PixelSizeInBits %d\n", dib.BitCount)
	fmt.Printf("- ImageSizeInBytes %d\n", dib.ImageSize)
}

func readHeaders(file *os.File) (*BMPHeader, *DIBHeader, error) {
	var bmpHeader BMPHeader
	if err := binary.Read(file, binary.LittleEndian, &bmpHeader); err != nil {
		return nil, nil, fmt.Errorf("error reading BMP header - %v", err)
	}

	// Read DIB header size (first 4 bytes of DIB header)
	var dibHeaderSize uint32
	if err := binary.Read(file, binary.LittleEndian, &dibHeaderSize); err != nil {
		return nil, nil, fmt.Errorf("error reading DIB header size - %v", err)
	}
	dibHeader := DIBHeader{DibHeaderSize: dibHeaderSize}

	// Validate supported header sizes
	if dibHeaderSize != 40 && dibHeaderSize != 108 && dibHeaderSize != 124 {
		return nil, nil, fmt.Errorf("%s has an unsupported BMP format (DIB Header Size = %d)", file.Name(), dibHeaderSize)
	}

	// Read the full DIB header as bytes
	dibHeaderBytes := make([]byte, dibHeaderSize-4) // Already read first 4 bytes
	if _, err := file.Read(dibHeaderBytes); err != nil {
		return nil, nil, fmt.Errorf("error reading full DIB header - %v", err)
	}
	buffer := bytes.NewReader(dibHeaderBytes)

	// Read mandatory BMP v3 fields
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Width); err != nil {
		return nil, nil, fmt.Errorf("error reading Width - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Height); err != nil {
		return nil, nil, fmt.Errorf("error reading Height - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Planes); err != nil {
		return nil, nil, fmt.Errorf("error reading Planes - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.BitCount); err != nil {
		return nil, nil, fmt.Errorf("error reading BitCount - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Compression); err != nil {
		return nil, nil, fmt.Errorf("error reading Compression - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ImageSize); err != nil {
		return nil, nil, fmt.Errorf("error reading ImageSize - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.XPixelsPerM); err != nil {
		return nil, nil, fmt.Errorf("error reading XPixelsPerM - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.YPixelsPerM); err != nil {
		return nil, nil, fmt.Errorf("error reading YPixelsPerM - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ColorsUsed); err != nil {
		return nil, nil, fmt.Errorf("error reading ColorsUsed - %v", err)
	}
	if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ColorsImp); err != nil {
		return nil, nil, fmt.Errorf("error reading ColorsImp - %v", err)
	}

	// Stop reading extra fields if DibHeaderSize == 40 (BMP v3)
	if dibHeaderSize == 40 {
		return &bmpHeader, &dibHeader, nil
	}

	// Read BMP v4 fields (if present)
	if dibHeaderSize >= 108 {
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.RedMask); err != nil {
			return nil, nil, fmt.Errorf("error reading RedMask - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.GreenMask); err != nil {
			return nil, nil, fmt.Errorf("error reading GreenMask - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.BlueMask); err != nil {
			return nil, nil, fmt.Errorf("error reading BlueMask - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.AlphaMask); err != nil {
			return nil, nil, fmt.Errorf("error reading AlphaMask - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ColorSpace); err != nil {
			return nil, nil, fmt.Errorf("error reading ColorSpace - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Endpoints); err != nil {
			return nil, nil, fmt.Errorf("error reading ColorSpaceEndpoints - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.GammaRed); err != nil {
			return nil, nil, fmt.Errorf("error reading GammaRed - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.GammaGreen); err != nil {
			return nil, nil, fmt.Errorf("error reading GammaGreen - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.GammaBlue); err != nil {
			return nil, nil, fmt.Errorf("error reading GammaBlue - %v", err)
		}
	}

	// Stop reading extra fields if DibHeaderSize == 108 (BMP v4)
	if dibHeaderSize == 108 {
		return &bmpHeader, &dibHeader, nil
	}

	// Read BMP v5 fields (if present)
	if dibHeaderSize == 124 {
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Intent); err != nil {
			return nil, nil, fmt.Errorf("error reading Intent - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ProfileData); err != nil {
			return nil, nil, fmt.Errorf("error reading ProfileData - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.ProfileSize); err != nil {
			return nil, nil, fmt.Errorf("error reading ProfileSize - %v", err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &dibHeader.Reserved); err != nil {
			return nil, nil, fmt.Errorf("error reading Reserved - %v", err)
		}
	}

	return &bmpHeader, &dibHeader, nil
}
