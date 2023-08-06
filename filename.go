package main

import (
	"fmt"
	"math"
	"strconv"
)

func createAlphabetFileName(digit int, prefix string,fileNumber int) (string, error) {
	if digit < 1 {
		return "", fmt.Errorf(negativeDigitErrorMsg)
	}
	if fileNumber < 0 {
		return "", fmt.Errorf(negativeFileNumberErrorMsg)
	}
	var fileName string
	if math.Pow(26, float64(digit)) <= float64(fileNumber) {
		return "", fmt.Errorf(tooBigFileNumberErrorMsg)
	}
	for i := 0; i < digit; i++ {
		remainder := fileNumber % 26
		fileName = string(rune(remainder+97)) + fileName
		fileNumber = fileNumber / 26
	}
	if prefix != "" {
		fileName = prefix + fileName
	} else {
		fileName = "x" + fileName
	}
	return fileName, nil
}

func createNumericFileName(digit int,prefix string, fileNumber int) (string, error) {
	if digit < 1 {
		return "", fmt.Errorf(negativeDigitErrorMsg)
	}
	if fileNumber < 0 {
		return "", fmt.Errorf(negativeFileNumberErrorMsg)
	}
	if math.Pow(10, float64(digit)) <= float64(fileNumber) {
		return "", fmt.Errorf(tooBigFileNumberErrorMsg)
	}
	var fileName string
	for i := 0; i < digit; i++ {
		fileName = strconv.Itoa(fileNumber%10) + fileName
		fileNumber = fileNumber / 10
	}
	if prefix != "" {
		fileName = prefix + fileName
	} else {
		fileName = "x" + fileName
	}
	
	return fileName, nil
}
