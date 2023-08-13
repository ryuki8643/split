package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type FileSplitter interface {
	Split(file *os.File, fileNameCreater FileNameCreater) error
}

type LineSplitter struct {
	separateLineNumber int64
}

const bufferSize = 1024 * 1024

func (s LineSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Initialize a buffer to store file data temporarily.
	buffer := make([]byte, 0, bufferSize)

	// Line counter to keep track of lines read from the input file.
	var lineCounter int64 = 0

	var outFile *os.File
	var outputFilePath string
	var err error
	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// If we have read 1000 lines, write to the output file.
		if lineCounter == 0 || lineCounter%s.separateLineNumber == 0 {
			// Increment the output file counter.

			// Create the output file.
			outputFilePath, err = fileNameCreater.Create(outputCounter)
			if err != nil {
				return err
			}
			outFile, err = os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf(createFileErrorMsg, err)
			}
			defer outFile.Close()
			outputCounter++
		}

		// Increase the line counter.
		lineCounter++

		// Add the line to the buffer.
		buffer, err = writeFileBy1Line(outFile, line, buffer)
		if err != nil {
			return err
		}

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(fileReadErrorMsg, err)
	}

	return nil
}

func writeFileBy1Line(outFile *os.File, line string, buffer []byte) ([]byte, error) {
	buffer = append(buffer, line...)
	buffer = append(buffer, '\n') // Add a newline character after each line.

	// Write the buffer data to the output file.
	_, err := outFile.Write(buffer)
	if err != nil {
		return nil, fmt.Errorf(fileWriteErrorMsg, err)
	}

	// Reset the buffer and line counter for the next output file.
	buffer = buffer[:0]
	return buffer, nil
}

type ByteSplitter struct {
	separateByteStr string
}

func (s ByteSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	separateByte, err := separateByteStrToInt(s.separateByteStr)
	if err != nil {
		return err
	}

	// Output file counter to keep track of split files.
	outputCounter := 0

	info, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := info.Size()
	if fileSize == 0 {
		return nil
	}

	// Read the input file and write to the output files.
	for {

		if err != nil && err != io.EOF {
			return fmt.Errorf(fileReadErrorMsg, err)
		}

		// Create the output file.
		outputFilePath, err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg, err)
		}
		defer outFile.Close()

		fileEnd, err := writeFileBy1KSize(file, outFile, outputFilePath, separateByte)
		if err != nil {
			return err
		}

		// If the output file is empty, delete it.
		outFileInfo, err := outFile.Stat()
		if err != nil {
			return err
		}
		size := outFileInfo.Size()
		if size == 0 {
			defer func() {
				outFile.Close()
				os.Remove(outputFilePath)
			}()
		}

		if fileEnd {
			break
		}
		// Increment the output file counter.
		outputCounter++
	}

	return nil
}

func separateByteStrToInt(separateByteStr string) (int, error) {
	re := regexp.MustCompile(`^(\d+)([kmgtpesy]*)b?$`)

	match := re.FindStringSubmatch(strings.ToLower(separateByteStr))
	if len(match) >= 3 {
		numberStr := match[1]
		unit := match[2]

		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return 0, fmt.Errorf(separateByteInvalidErrorMsg)
		}

		factor := 1
		switch unit {
		case "k":
			factor = 1024
		case "m":
			factor = 1024 * 1024
		case "g":
			factor = 1024 * 1024 * 1024
		case "t":
			factor = 1024 * 1024 * 1024 * 1024
		case "p":
			factor = 1024 * 1024 * 1024 * 1024 * 1024
		case "e":
			factor = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		}

		separateByte := number * factor
		return separateByte, nil
	}
	return 0, fmt.Errorf(separateByteInvalidErrorMsg)

}

func writeFileBy1KSize(file, outFile *os.File, outputputFilePath string, size int) (bool, error) {

	// Create the buffer for reading the input file.
	buffer := make([]byte, 1024)

	// Read 1KB of data from the input file.
	for size > 0 {
		if size < 1024 {
			buffer = make([]byte, size)
		}
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return false, fmt.Errorf(fileReadErrorMsg, err)
		}

		if n == 0 {
			// Reached the end of the file, exit the loop.
			return true, nil
		}

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer[:n])
		if err != nil {
			return false, fmt.Errorf(fileWriteErrorMsg, err)
		}
		size -= n

		// Reset the buffer for the next output file.
		buffer = make([]byte, 1024)
	}

	return false, nil
}

type PieceSplitter struct {
	chunkStr string
}

func (s PieceSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {
	chunk, err := parseCHUNK(s.chunkStr)
	if err != nil {
		return err
	}
	var splitter FileSplitter
	if chunk.R {
		splitter = PieceLineRoundRobinSplitter{chunk.N}
	} else if chunk.L {
		splitter = PieceLineSplitter{chunk.N}
	} else {
		splitter = PieceByteSplitter{chunk.N}
	}

	err = splitter.Split(file, fileNameCreater)
	if err != nil {
		return err
	}
	if chunk.K > 0 {
		fileName, err := fileNameCreater.Create(int(chunk.K) - 1)
		if err != nil {
			return err
		}
		file, err = os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
	return nil

}

type PieceByteSplitter struct {
	separatePieceNumber int64
}

func (s PieceByteSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Get the file information.
	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	// Calculate the size of each piece.
	splitSize := info.Size() / s.separatePieceNumber

	if info.Size()%s.separatePieceNumber != 0 {
		splitSize++
	}
	if info.Size() == 0 {
		return nil
	}

	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file and write to the output files.
	for {
		if err != nil && err != io.EOF {
			return fmt.Errorf(fileReadErrorMsg, err)
		}

		// Create the output file.
		outputFilePath, err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg, err)
		}
		defer outFile.Close()

		fileEnd, err := writeFileBy1KSize(file, outFile, outputFilePath, int(splitSize))
		if err != nil {
			return err
		}

		// If the output file is empty, delete it.
		outFileInfo, err := outFile.Stat()
		if err != nil {
			return err
		}
		if outFileInfo.Size() == 0 {
			defer func() {
				outFile.Close()
				os.Remove(outputFilePath)
			}()
		}

		if fileEnd {
			break
		}

		// Increment the output file counter.
		outputCounter++
	}

	return nil
}

type PieceLineSplitter struct {
	separatePieceNumber int64
}

func (s PieceLineSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Initialize a buffer to store file data temporarily.
	buffer := make([]byte, 0, bufferSize)
	var outFile *os.File

	// count file line number
	fileLineNum, err := countLinesByFile(file)
	if err != nil {
		return err
	}

	fileLinesPerPiece := fileLineNum / s.separatePieceNumber

	if fileLineNum%s.separatePieceNumber != 0 {
		fileLinesPerPiece++
	}

	// Line counter to keep track of lines read from the input file.
	var lineCounter int64 = 0

	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// If we have read 1000 lines, write to the output file.
		if lineCounter == 0 || lineCounter%fileLinesPerPiece == 0 {
			// Create the output file.
			outputFilePath, err := fileNameCreater.Create(outputCounter)
			if err != nil {
				return err
			}
			outFile, err = os.Create(outputFilePath)
			if err != nil {
				return fmt.Errorf(createFileErrorMsg, err)
			}

			// Close the output file.
			defer outFile.Close()

			// Increment the output file counter.
			outputCounter++

		}

		// Increase the line counter.
		lineCounter++

		buffer, err = writeFileBy1Line(outFile, line, buffer)
		if err != nil {
			return err
		}

	}

	// Check for errors.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf(fileReadErrorMsg, err)
	}

	return nil
}

type PieceLineRoundRobinSplitter struct {
	separatePieceNumber int64
}

func (s PieceLineRoundRobinSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Initialize a buffer to store file data temporarily.
	buffer := make([]byte, 0, bufferSize)

	// count file line number
	fileLineNum, err := countLinesByFile(file)
	if err != nil {
		return err
	}

	fileLinesPerPiece := fileLineNum / s.separatePieceNumber

	if fileLineNum%s.separatePieceNumber != 0 {
		fileLinesPerPiece++
	}

	// Line counter to keep track of lines read from the input file.
	lineCounter := 0

	scanner := bufio.NewScanner(file)
	outFiles := make([]*os.File, 0, s.separatePieceNumber)

	for scanner.Scan() {
		if len(outFiles) == (lineCounter % int(s.separatePieceNumber)) {
			outputFilePath, err := fileNameCreater.Create(lineCounter % int(s.separatePieceNumber))
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf(createFileErrorMsg, err)
			}
			outFiles = append(outFiles, outFile)
			defer outFile.Close()
		}
		line := scanner.Text()
		// Create the output file.
		outFile := outFiles[lineCounter%int(s.separatePieceNumber)]

		buffer, err = writeFileBy1Line(outFile, line, buffer)
		if err != nil {
			return err
		}

		// Increment the output file counter.
		lineCounter++

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(fileReadErrorMsg, err)
	}

	return nil
}

func countLinesByFile(file *os.File) (int64, error) {
	fileForCountLine, err := os.Open(file.Name())
	if err != nil {
		return 0, fmt.Errorf(fileReadErrorMsg, err)
	}
	defer fileForCountLine.Close()
	var count int64
	scanner := bufio.NewScanner(fileForCountLine)
	for scanner.Scan() {
		count++
	}
	return count, nil
}

type chunk struct {
	R bool
	L bool
	K int64
	N int64
}

func parseCHUNK(chunkStr string) (chunk, error) {
	result := chunk{}
	parts := strings.Split(chunkStr, "/")
	var err error
	if len(parts) == 1 {
		result.N, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}
	} else if len(parts) == 2 {
		if parts[0] == "l" || parts[0] == "r" {
			result.L = parts[0] == "l"
			result.R = parts[0] == "r"
		} else {
			result.K, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
			}
			if result.K <= 0 {
				return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
			}
		}
		result.N, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}
	} else if len(parts) == 3 {
		if parts[0] == "l" || parts[0] == "r" {
			result.L = parts[0] == "l"
			result.R = parts[0] == "r"
		} else {
			return result, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}
		result.K, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}
		if result.K <= 0 {
			return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}

		result.N, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
		}
	} else {
		return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
	}
	if result.N < result.K {
		return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
	}
	if result.N <= 0 {
		return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
	}
	return result, nil
}
