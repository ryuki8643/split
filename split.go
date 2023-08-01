package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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

func createAlphabetFileName(digit int, fileNumber int) (string,error) {
	if (digit<0 || fileNumber<0){
		return "",fmt.Errorf("digit or fileNumber is negative")
	}
	var fileName string
	if (math.Pow(26, float64(digit))<=float64(fileNumber)){
		return "",fmt.Errorf("fileNumber is too big")
	}
	for i := 0; i < digit; i++ {
		remainder := fileNumber % 26
		fileName = string(rune(remainder + 97)) + fileName
		fileNumber = fileNumber / 26
	}
	return fileName,nil
}

func createNumericFileName(digit int, fileNumber int) (string,error) {
	if (digit<0 || fileNumber<0){
		return "",fmt.Errorf("digit or fileNumber is negative")
	}
	if (math.Pow(10, float64(digit))<=float64(fileNumber)){
		return "",fmt.Errorf("fileNumber is too big")
	}
	var fileName string
	for i := 0; i < digit; i++ {
		fileName = strconv.Itoa(fileNumber%10) + fileName
		fileNumber = fileNumber / 10
	}
	return fileName,nil
}

