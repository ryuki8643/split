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

func (s LineSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Set the maximum memory limit to 1GB (in bytes).
	const maxMemoryLimit = 1 * 1024 * 1024 * 1024

	// Initialize a buffer to store file data temporarily.
	buffer := make([]byte, 0)

	// Line counter to keep track of lines read from the input file.
	var lineCounter int64= 0

	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Increase the line counter.
		lineCounter++

		// Add the line to the buffer.
		buffer = append(buffer, line...)
		buffer = append(buffer, '\n') // Add a newline character after each line.

		// If we have read 1000 lines, write to the output file.
		if lineCounter == s.separateLineNumber {
			// Create the output file.
			outputFilePath, err := fileNameCreater.Create(outputCounter)
			if err != nil {
				return err
			}
			outFile, err := os.Create(outputFilePath)
			if err != nil {
				return fmt.Errorf(createFileErrorMsg, err)
			}

			// Write the buffer data to the output file.
			_, err = outFile.Write(buffer)
			if err != nil {
				return fmt.Errorf(fileWriteErrorMsg, err)
			}

			// Close the output file.
			outFile.Close()

			// Reset the buffer and line counter for the next output file.
			buffer = buffer[:0]
			lineCounter = 0

			// Increment the output file counter.
			outputCounter++

		}

		if len(buffer) > maxMemoryLimit {
			return fmt.Errorf(maxMemoryLimitExceededErrorMsg)
		}
	}

	if len(buffer) > 0 {
		// Create the output file.
		outputFilePath, err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg, err)
		}

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer)
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Close the output file.
		outFile.Close()

		// Increment the output file counter.
		outputCounter++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(fileReadErrorMsg, err)
	}

	return nil
}

type ByteSplitter struct {
	separateByteStr string
}

func (s ByteSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	separateByte, err := separateByteStrToInt(s.separateByteStr)
	if err != nil {
		return err
	}

	// 1GB memory limit (in bytes).
	const maxMemoryLimit = 1 * 1024 * 1024 * 1024

	// Buffer to store file data temporarily.
	buffer := make([]byte, 0, maxMemoryLimit)

	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file and write to the output files.
	for {
		// Read separateByte of data from the input file.
		n, err := file.Read(buffer[:separateByte])

		if err != nil && err != io.EOF {
			return fmt.Errorf(fileReadErrorMsg, err)
		}

		if n == 0 {
			// Reached the end of the file, exit the loop.
			break
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

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Reset the buffer for the next output file.
		buffer = buffer[:0]

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
		}

		separateByte := number * factor
		return separateByte, nil
	}
	return 0, fmt.Errorf(separateByteInvalidErrorMsg)

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

	// 1GB memory limit (in bytes).
	const maxMemoryLimit = 1 * 1024 * 1024 * 1024

	// Buffer to store file data temporarily.
	buffer := make([]byte, 0, maxMemoryLimit)

	// Output file counter to keep track of split files.
	outputCounter := 0

	// Read the input file and write to the output files.
	for {
		// Read separateByte of data from the input file.
		n, err := file.Read(buffer[:splitSize])

		if err != nil && err != io.EOF {
			return fmt.Errorf(fileReadErrorMsg, err)
		}

		if n == 0 {
			// Reached the end of the file, exit the loop.
			break
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

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Reset the buffer for the next output file.
		buffer = buffer[:0]

		// Increment the output file counter.
		outputCounter++
	}

	return nil
}

type PieceLineSplitter struct {
	separatePieceNumber int64
}

func (s PieceLineSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	// Set the maximum memory limit to 1GB (in bytes).
	const maxMemoryLimit = 1 * 1024 * 1024 * 1024

	// Initialize a buffer to store file data temporarily.
	buffer := make([]byte, 0)

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

		// Increase the line counter.
		lineCounter++

		// Add the line to the buffer.
		buffer = append(buffer, line...)
		buffer = append(buffer, '\n') // Add a newline character after each line.

		// If we have read 1000 lines, write to the output file.
		if lineCounter == fileLinesPerPiece {
			// Create the output file.
			outputFilePath, err := fileNameCreater.Create(outputCounter)
			if err != nil {
				return err
			}
			outFile, err := os.Create(outputFilePath)
			if err != nil {
				return fmt.Errorf(createFileErrorMsg, err)
			}

			// Write the buffer data to the output file.
			_, err = outFile.Write(buffer)
			if err != nil {
				return fmt.Errorf(fileWriteErrorMsg, err)
			}

			// Close the output file.
			outFile.Close()

			// Reset the buffer and line counter for the next output file.
			buffer = buffer[:0]
			lineCounter = 0

			// Increment the output file counter.
			outputCounter++

		}

		if len(buffer) > maxMemoryLimit {
			return fmt.Errorf(maxMemoryLimitExceededErrorMsg)
		}
	}

	if len(buffer) > 0 {
		// Create the output file.
		outputFilePath, err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg, err)
		}

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer)
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Close the output file.
		outFile.Close()

		// Increment the output file counter.
		outputCounter++
	}

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
	buffer := make([]byte, 0)

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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Increase the line counter.
		lineCounter++

		// Add the line to the buffer.
		buffer = append(buffer, line...)
		buffer = append(buffer, '\n') // Add a newline character after each line.
		// Create the output file.
		outputFilePath, err := fileNameCreater.Create(outputCounter % int(s.separatePieceNumber))
		if err != nil {
			return err
		}
		outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg, err)
		}

		// Write the buffer data to the output file.
		_, err = outFile.Write(buffer)
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Close the output file.
		err = outFile.Close()
		if err != nil {
			return fmt.Errorf(fileWriteErrorMsg, err)
		}

		// Increment the output file counter.
		outputCounter++

		// Reset the buffer and line counter for the next output file.
		buffer = buffer[:0]

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(fileReadErrorMsg, err)
	}

	return nil
}

func countLinesByFile(file *os.File) (int64, error) {
	fileForCountLine, err := os.Open(file.Name())
	defer fileForCountLine.Close()
	if err != nil {
		return 0, fmt.Errorf(fileReadErrorMsg, err)

	}
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
			if result.K == 0 {
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
		if result.K == 0 {
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
	if result.N == 0 {
		return chunk{}, fmt.Errorf(chunkFormatInvalidErrorMsg)
	}
	return result, nil
}
