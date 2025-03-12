package utils

import (
	"fmt"
	"os"
)

// Checks for error and return an error message if there is some
func HandleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
