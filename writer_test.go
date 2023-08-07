package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLineFileSplitterSplit1000Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := LineFileSplitter{1000}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())
	for i := 1; i <= 2007; i++ {
		fmt.Fprintln(testFile, "line", i)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files were created.
	if _, err := os.Stat("outputaa"); os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	if _, err := os.Stat("outputab"); os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if _, err := os.Stat("outputac"); os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	if countLines(output1) != 1000 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLines(output1))
	}
	
	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLines(output2) != 1000 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLines(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLines(output3) != 7 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLines(output3))
	}

	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}


	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) != testFileString {
		t.Fatal("Incorrect output file content.")
	}

}

func TestLineFileSplitterSplit500Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := LineFileSplitter{500}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())
	for i := 1; i <= 716; i++ {
		fmt.Fprintln(testFile, "line", i)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files were created.
	if _, err := os.Stat("outputaa"); os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	if _, err := os.Stat("outputab"); os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	if countLines(output1) != 500 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLines(output1))
	}
	
	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLines(output2) != 216 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLines(output2))
	}

	if countFiles() != 2 {
		t.Fatal("Incorrect number of output files.")
	}


	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2) != testFileString {
		t.Fatal("Incorrect output file content.")
	}

}



func TestLineFileSplitterSplitWithEmptyFile(t *testing.T) {
	
	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := LineFileSplitter{1000}

	// Create a test file with 0 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files were created.
	if _, err := os.Stat("outputaa"); os.IsExist(err) {
		t.Fatal("null file was created.")
	}
}





func countLines(data []byte) int {
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	return count
}

func countFiles() int {
	outputPrefix := "output" // Prefix of files to be counted

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get the current working directory:", err)
		return 0
	}

	// Get the list of files
	files, err := filepath.Glob(filepath.Join(currentDir, outputPrefix+"*"))
	if err != nil {
		fmt.Println("Failed to get the list of files:", err)
		return 0
	}

	return len(files)
}

func deleteOutputFiles() {
	outputPrefix := "output" // Prefix of files to be deleted

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get the current working directory:", err)
		return
	}

	// Get the list of files
	files, err := filepath.Glob(filepath.Join(currentDir, outputPrefix+"*"))
	if err != nil {
		fmt.Println("Failed to get the list of files:", err)
		return
	}

	// Delete the files
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Failed to delete the file:", err)
		} else {
			fmt.Println("File deleted:", file)
		}
	}
}

func TestByteFileSplitterSplit1kByte(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := ByteSplitter{"1k"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())
	
	// Write 2500 bytes to the file.
	byteCount := 2500
	data := make([]byte,byteCount)
	for i := 0; i < byteCount; i++ {
		data[i] = byte(i % 256)
	}
	_, err = testFile.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files were created.
	out1Stat, err := os.Stat("outputaa");
	if  os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1024 {
		t.Fatal("outputaa file has incorrect size. Expected 1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab");
	if  os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1024 {
		t.Fatal("outputab file has incorrect size. Expected 1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac");
	if  os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 452 {
		t.Fatal("outputac file has incorrect size. Expected 452, got ", out3Stat.Size())
	}	

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	
	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}

	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output count is correct
	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}


	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files have the correct content
	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) != testFileString {
		t.Fatal("Incorrect output file content.")
	}

}


func TestByteFileSplitterSplit1mByte(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := ByteSplitter{"1m"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())
	
	// Write 2500000 bytes to the file.
	byteCount := 2500000
	data := make([]byte,byteCount)
	for i := 0; i < byteCount; i++ {
		data[i] = byte(i % 256)
	}
	_, err = testFile.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files were created.
	out1Stat, err := os.Stat("outputaa");
	if  os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1024*1024 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab");
	if  os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1024*1024 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac");
	if  os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 402848 {
		t.Fatal("outputac file has incorrect size. Expected 402848, got ", out3Stat.Size())
	}	

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	
	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}

	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output count is correct
	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}


	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check that the output files have the correct content
	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) != testFileString {
		t.Fatal("Incorrect output file content.")
	}

}

func TestByteFileSplitterCannotSplit0ByteFile(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := ByteSplitter{"1m"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())
	
	// Write 0 bytes to the file.
	byteCount := 0
	data := make([]byte,byteCount)
	for i := 0; i < byteCount; i++ {
		data[i] = byte(i % 256)
	}
	_, err = testFile.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	

	// Check that the output count is correct
	if countFiles() != 0 {
		t.Fatal("Incorrect number of output files.")
	}

}



func TestSeparateByteStrToInt(t *testing.T) {
    tests := []struct {
        input string
        expected int
        err error
    }{
        {"100", 100, nil},
        {"1k", 1024, nil},
        {"2m", 2097152, nil},
        {"abc", 0, fmt.Errorf(separateByteInvalidErrorMsg)},
    }

    for _, test := range tests {
        result, err := separateByteStrToInt(test.input)
        if result != test.expected || !reflect.DeepEqual(err, test.err) {
            t.Errorf("separateByteStrToInt(%q) = (%v, %v), expected (%v, %v)", test.input, result, err, test.expected, test.err)
        }
    }
}