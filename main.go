package main

import (
	"fmt"
	"io"
	"os"
)

const (
	tooBigFileNumberErrorMsg       = "fileNumber is too big"
	negativeDigitErrorMsg          = "digit is negative"
	negativeFileNumberErrorMsg     = "fileNumber is negative"
	maxMemoryLimitExceededErrorMsg = "memory limit exceeded"
	createFileErrorMsg             = "failed to create the output file:%w"
	fileWriteErrorMsg              = "failed to write to the output file:%w"
	fileReadErrorMsg               = "failed to read from the input file:%w"
	fileCloseErrorMsg              = "failed to close the output file:%w"
	fileOpenErrorMsg               = "failed to open the input file:%w"
	separateByteInvalidErrorMsg    = "separate byte is invalid"
	chunkFormatInvalidErrorMsg     = "chunk format is invalid"
	invalidArgumentErrorMsg        = "invalid argument:%d"
	tooManyFlagErrorMsg            = "only one of -l, -n, -b can be used"
)

var writer io.Writer

func init() {
	writer = os.Stdout
}

func main() {
	fileName, splitter, fileNameCreater, err := ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	err = splitter.Split(file, fileNameCreater)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
