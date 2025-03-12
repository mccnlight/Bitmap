package bmp

import (
	"errors"
	"fmt"
	"strings"
)

// Parses command-line arguments while maintaining order
func ParseArgs(args []string) (command string, filename string, outputFilename string, orderedOptions []Option, err error) {
	if len(args) < 2 {
		return "", "", "", nil, errors.New("invalid number of arguments")
	}

	command = args[0] // "header" or "apply"

	// Handle "header" command (only requires filename)
	if command == "header" {
		if len(args) != 2 {
			return "", "", "", nil, errors.New("usage: ./bitmap header <bmp_file>")
		}
		filename = args[1]
		return command, filename, "", nil, nil
	}

	// Handle "apply" command (requires at least one option, input file, and output file)
	if command == "apply" {
		if len(args) < 4 {
			return "", "", "", nil, errors.New("usage: ./bitmap apply [options] <source_file> <output_file>")
		}

		filename = args[len(args)-2]       // Second-to-last argument is the source file
		outputFilename = args[len(args)-1] // Last argument is the output file

		for i := 1; i < len(args)-2; i++ { // Ignore the last two arguments (file names)
			if strings.HasPrefix(args[i], "--") {
				// Break down the option into the option name and its associated value
				parts := strings.SplitN(args[i], "=", 2)
				if len(parts) != 2 {
					return "", "", "", nil, fmt.Errorf("invalid option format: %s", args[i])
				}
				name, value := parts[0], parts[1]

				// Slice of struct preserves the insertion order of the applied options
				orderedOptions = append(orderedOptions, Option{Name: name, Value: value})
			} else {
				return "", "", "", nil, fmt.Errorf("unexpected argument: %s", args[i])
			}
		}

		return command, filename, outputFilename, orderedOptions, nil
	}

	// If command is neither "header" nor "apply", then return an error
	return "", "", "", nil, fmt.Errorf("unknown command: %s", command)
}
