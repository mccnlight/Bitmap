package utils

import "fmt"

// Displays general usage instructions
func DisplayGeneralHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bitmap <command> [arguments]")
	fmt.Println()
	fmt.Println("The commands are:")
	fmt.Println("  header    prints bitmap file header information")
	fmt.Println("  apply     applies processing to the image and saves it to the file")
}

// Displays usage instructions for header command
func DisplayHeaderHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bitmap header <source_file>")
	fmt.Println()
	fmt.Println("Description:")
	fmt.Println("  Prints bitmap file header information")
}

// Displays usage instructions for apply command
func DisplayApplyHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bitmap apply [options] <source_file> <output_file>")
	fmt.Println()
	fmt.Println("The options are:")
	fmt.Println("  -h, --help                                                      prints program usage information")
	fmt.Println("  --mirror=<horizontal|vertical>                                  mirrors the image along the specified axis")
	fmt.Println("  --filter=<blue|red|green|grayscale|negative|pixelate|blur>      applies a specified filter to the image")
	fmt.Println("  --rotate=<right|left|90|-90|180|-180|270|-270>                  rotates the image by the specified angle")
	fmt.Println("  --crop=<offsetX-offsetY-width-height>                           crops the image based on the specified offset and dimensions")
	fmt.Println()
	fmt.Println("Note:")
	fmt.Println("  Multiple options can be combined and applied sequentially")
}
