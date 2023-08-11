package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	tooBigFileNumberErrorMsg   = "fileNumber is too big"
	negativeDigitErrorMsg      = "digit is negative"
	negativeFileNumberErrorMsg = "fileNumber is negative"

	maxMemoryLimitExceededErrorMsg = "memory limit exceeded"
	createFileErrorMsg             = "failed to create the output file:%w"
	fileWriteErrorMsg              = "failed to write to the output file:%w"
	fileReadErrorMsg               = "failed to read from the input file:%w"
	fileCloseErrorMsg              = "failed to close the output file:%w"
	fileOpenErrorMsg               = "failed to open the input file:%w"
	separateByteInvalidErrorMsg    = "separate byte is invalid"
)

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := flag.Args()[0]

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
