package main

import (
	"fmt"
	"math"
	"strconv"
)

type FileNameCreater interface {
	Create(fileNumber int) (string, error)
}

type AlphabetFileNameCreater struct {
	digit int
	prefix string
}

func (fileNameCreater AlphabetFileNameCreater)Create(fileNumber int) (string, error) {
	if fileNameCreater.digit < 1 {
		return "", fmt.Errorf(negativeDigitErrorMsg)
	}
	if fileNumber < 0 {
		return "", fmt.Errorf(negativeFileNumberErrorMsg)
	}
	var fileName string
	if math.Pow(26, float64(fileNameCreater.digit)) <= float64(fileNumber) {
		return "", fmt.Errorf(tooBigFileNumberErrorMsg)
	}
	for i := 0; i < fileNameCreater.digit; i++ {
		remainder := fileNumber % 26
		fileName = string(rune(remainder+97)) + fileName
		fileNumber = fileNumber / 26
	}
	if fileNameCreater.prefix != "" {
		fileName = fileNameCreater.prefix + fileName
	} else {
		fileName = "x" + fileName
	}
	return fileName, nil
}

type NumericFileNameCreater struct {
	digit int
	prefix string
}

func (fileNameCreater NumericFileNameCreater) Create(fileNumber int) (string, error) {
	if fileNameCreater.digit < 1 {
		return "", fmt.Errorf(negativeDigitErrorMsg)
	}
	if fileNumber < 0 {
		return "", fmt.Errorf(negativeFileNumberErrorMsg)
	}
	if math.Pow(10, float64(fileNameCreater.digit)) <= float64(fileNumber) {
		return "", fmt.Errorf(tooBigFileNumberErrorMsg)
	}
	var fileName string
	for i := 0; i < fileNameCreater.digit; i++ {
		fileName = strconv.Itoa(fileNumber%10) + fileName
		fileNumber = fileNumber / 10
	}
	if fileNameCreater.prefix != "" {
		fileName = fileNameCreater.prefix + fileName
	} else {
		fileName = "x" + fileName
	}
	
	return fileName, nil
}
