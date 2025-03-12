package main

import (
	"fmt"
	"os"

	"git.platform.alem.school/amibragim/bitmap/bmp"
	"git.platform.alem.school/amibragim/bitmap/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.DisplayGeneralHelp()
		os.Exit(1)
	}

	if len(os.Args) == 3 && os.Args[1] == "header" && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		utils.DisplayHeaderHelp()
		os.Exit(0)
	}

	if len(os.Args) == 3 && os.Args[1] == "apply" && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		utils.DisplayApplyHelp()
		os.Exit(0)
	}

	command, filename, outputFilename, orderedOptions, err := bmp.ParseArgs(os.Args[1:])
	utils.HandleError(err)

	bmpHeader, dibHeader, err := bmp.ReadHeaders(filename)
	utils.HandleError(err)

	var croppedWidth, croppedHeight int

	switch command {
	case "header":
		bmp.PrintHeader(bmpHeader, dibHeader)

	case "apply":
		pixels, err := bmp.ReadPixels(filename, bmpHeader, dibHeader)
		utils.HandleError(err)

		// Process options sequentially
		for _, opt := range orderedOptions {
			switch opt.Name {
			case "--mirror":
				pixels, err = bmp.ApplyMirror(pixels, int(dibHeader.Width), int(dibHeader.Height), opt.Value)

			case "--filter":
				pixels, err = bmp.ApplyFilter(pixels, int(dibHeader.Width), int(dibHeader.Height), opt.Value)

			case "--rotate":
				angle, err := bmp.ParseRotationValue(opt.Value)
				utils.HandleError(err)

				var newWidth, newHeight int
				pixels, newWidth, newHeight, err = bmp.ApplyRotate(pixels, int(dibHeader.Width), int(dibHeader.Height), angle)

				// Update image properties after rotating
				dibHeader.Width = int32(newWidth)
				dibHeader.Height = int32(newHeight)

			case "--crop":
				pixels, croppedWidth, croppedHeight, err = bmp.ApplyCrop(pixels, int(dibHeader.Width), int(dibHeader.Height), opt.Value)

				// Update image properties after cropping
				dibHeader.Width = int32(croppedWidth)
				dibHeader.Height = int32(croppedHeight)
				dibHeader.ImageSize = uint32(croppedWidth * croppedHeight * 3)
			default:
				utils.HandleError(fmt.Errorf("undefined option - %s", opt.Name))
			}
			utils.HandleError(err)
		}

		err = bmp.WritePixels(outputFilename, bmpHeader, dibHeader, pixels)
		utils.HandleError(err)

	default:
		utils.DisplayGeneralHelp()
		os.Exit(1)
	}
}
