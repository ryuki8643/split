package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type FileSplitter interface {
	Split(file *os.File, fileNameCreater FileNameCreater) error
}


type LineFileSplitter struct {
	separateLineNumber int
}

func (s LineFileSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

    // Set the maximum memory limit to 1GB (in bytes).
    const maxMemoryLimit = 1 * 1024 * 1024 * 1024

    // Initialize a buffer to store file data temporarily.
    buffer := make([]byte, 0)

    // Line counter to keep track of lines read from the input file.
    lineCounter := 0

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
            outputFilePath,err := fileNameCreater.Create(outputCounter)
            if err != nil {
                return err
            }
            outFile, err := os.Create(outputFilePath)
            if err != nil {
                return fmt.Errorf(createFileErrorMsg,err)
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
            outputFilePath,err := fileNameCreater.Create(outputCounter)
            if err != nil {
                return err
            }
            outFile, err := os.Create(outputFilePath)
            if err != nil {
                return fmt.Errorf(createFileErrorMsg,err)
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

    fmt.Println("File splitting is complete.")
    return nil
}

type ByteSplitter struct {
	separateByteStr string
}

func (s ByteSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {

	separateByte,err := separateByteStrToInt(s.separateByteStr)
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
		outputFilePath,err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg,err)
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

	fmt.Println("File splitting is complete.")
	return nil
}


func separateByteStrToInt(separateByteStr string) (int, error) {
	numberPattern := `^\d+$`
	kiloBytePattern := `^\d+k$`
	megaBytePattern := `^\d+m$`

	// Compile the regular expressions
	numberRegex := regexp.MustCompile(numberPattern)
	kiloByteRegex := regexp.MustCompile(kiloBytePattern)
	megaByteRegex := regexp.MustCompile(megaBytePattern)
	var separateByte int
	if  numberRegex.MatchString(separateByteStr) {
		s,err:=strconv.Atoi(separateByteStr)
		if err != nil {
			return 0,fmt.Errorf(separateByteInvalidErrorMsg)
		}
		separateByte = s
	} else if kiloByteRegex.MatchString(separateByteStr) {
		s,err := strconv.Atoi(separateByteStr[:len(separateByteStr)-1])
		if err != nil {
			return 0,fmt.Errorf(separateByteInvalidErrorMsg)
		}
		separateByte =s * 1024
	} else if megaByteRegex.MatchString(separateByteStr) {
		s,err := strconv.Atoi(separateByteStr[:len(separateByteStr)-1])
		if err != nil {
			return 0,fmt.Errorf(separateByteInvalidErrorMsg)
		}
		separateByte = s * 1024 * 1024
	} else {
		return 0,fmt.Errorf(separateByteInvalidErrorMsg)
	}
	return separateByte,nil
}


type PieceSplitter struct {
	separatePieceNumberStr int64
}

func (s PieceSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {
	
	// Get the file information.
	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	// Calculate the size of each piece.
	splitSize:=info.Size()/s.separatePieceNumberStr

	if info.Size()%s.separatePieceNumberStr != 0 {
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
		outputFilePath,err := fileNameCreater.Create(outputCounter)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf(createFileErrorMsg,err)
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

	fmt.Println("File splitting is complete.")
	return nil
}