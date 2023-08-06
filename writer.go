package main

import (
	"bufio"
	"fmt"
	"os"
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

