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

			separateByte:= number * factor
			return separateByte,nil
		}
		return 0, fmt.Errorf(separateByteInvalidErrorMsg)
			
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

type PieceLineSplitter struct {
	separatePieceNumberStr int64
}

func (s PieceLineSplitter) Split(file *os.File, fileNameCreater FileNameCreater) error {
	
    // Set the maximum memory limit to 1GB (in bytes).
    const maxMemoryLimit = 1 * 1024 * 1024 * 1024

    // Initialize a buffer to store file data temporarily.
    buffer := make([]byte, 0)

	// count file line number
	fileLineNum := countLinesByFile(file)


	fileSizePerPiece := fileLineNum / s.separatePieceNumberStr

	if fileLineNum%s.separatePieceNumberStr != 0 {
		fileSizePerPiece++
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
        if lineCounter == fileSizePerPiece {
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

func countLinesByFile(file *os.File) int64 {
	var count int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count
}