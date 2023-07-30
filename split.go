package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Get file size
	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fileSize := info.Size()

	// Calculate split sizes
	splitSize := fileSize / 2
	if fileSize%2 != 0 {
		splitSize++
	}

	// Create output files
	out1, err := os.Create(filename + "aa")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer out1.Close()

	out2, err := os.Create(filename + "ab")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer out2.Close()

	// Copy data to output files
	_, err = io.CopyN(out1, file, splitSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	_, err = io.Copy(out2, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}