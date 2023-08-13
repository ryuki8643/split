package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var buffer *bytes.Buffer

func init() {
	buffer = &bytes.Buffer{}
	writer = buffer
}

func TestLineFileSplitterSplit1000Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := LineSplitter{1000}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 1000 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 1000 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 7 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
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
	splitter := LineSplitter{500}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 500 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 216 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
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
	splitter := LineSplitter{1000}

	// Create a test file with 0 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

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

func countLinesByByte(data []byte) int {
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

func TestByteFileSplitterSplit3072Byte(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := ByteSplitter{"1k"}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 2500 bytes to the file.
	byteCount := 3072
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1024 {
		t.Fatal("outputaa file has incorrect size. Expected 1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1024 {
		t.Fatal("outputab file has incorrect size. Expected 1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 1024 {
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

func TestByteFileSplitterSplit2500Byte(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := ByteSplitter{"1k"}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 2500 bytes to the file.
	byteCount := 2500
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1024 {
		t.Fatal("outputaa file has incorrect size. Expected 1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1024 {
		t.Fatal("outputab file has incorrect size. Expected 1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
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

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 2500000 bytes to the file.
	byteCount := 2500000
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1024*1024 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1024*1024 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
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

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 0 bytes to the file.
	byteCount := 0
	data := make([]byte, byteCount)
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
		input    string
		expected int
		err      error
	}{
		{"100", 100, nil},
		{"1k", 1024, nil},
		{"1K", 1024, nil},
		{"1kb", 1024, nil},
		{"1KB", 1024, nil},
		{"2m", 2097152, nil},
		{"2M", 2097152, nil},
		{"2mb", 2097152, nil},
		{"2MB", 2097152, nil},
		{"1g", 1073741824, nil},
		{"1G", 1073741824, nil},
		{"1gb", 1073741824, nil},
		{"1GB", 1073741824, nil},
		{"1t", 1099511627776, nil},
		{"1T", 1099511627776, nil},
		{"1tb", 1099511627776, nil},
		{"1TB", 1099511627776, nil},
		{"1p", 1125899906842624, nil},
		{"1P", 1125899906842624, nil},
		{"1pb", 1125899906842624, nil},
		{"1PB", 1125899906842624, nil},
		{"1e", 1152921504606846976, nil},
		{"1E", 1152921504606846976, nil},
		{"1eb", 1152921504606846976, nil},
		{"1EB", 1152921504606846976, nil},
		{"abc", 0, fmt.Errorf(separateByteInvalidErrorMsg)},
	}

	for _, test := range tests {
		result, err := separateByteStrToInt(test.input)
		if result != test.expected || !reflect.DeepEqual(err, test.err) {
			t.Errorf("separateByteStrToInt(%q) = (%v, %v), expected (%v, %v)", test.input, result, err, test.expected, test.err)
		}
	}
}

func TestBytePieceSplitter5000ByteFileto3file(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceByteSplitter{3}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 5000 bytes to the file.
	byteCount := 5000
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1667 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1667 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 1666 {
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

func TestBytePieceSplitter3000ByteFileto3file(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceByteSplitter{3}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 5000 bytes to the file.
	byteCount := 3000
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1000 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1000 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 1000 {
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

func TestBytePieceSplitter0ByteFileto0file(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceByteSplitter{3}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 5000 bytes to the file.
	byteCount := 0
	data := make([]byte, byteCount)
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

	if countFiles() != 0 {
		t.Fatal("Incorrect number of output files.")
	}

}

func TestPieceLineFileSplitterSplit2007Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineSplitter{3}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 669 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 669 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
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

func TestPieceLineFileSplitterSplit713Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineSplitter{4}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 713; i++ {
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
	if _, err := os.Stat("outputad"); os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	if countLinesByByte(output1) != 179 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 179 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 179 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}
	output4, err := os.ReadFile("outputad")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output4) != 176 {
		t.Fatal("outputad file has incorrect number of lines. Expected 7, got ", countLinesByByte(output4))
	}

	if countFiles() != 4 {
		t.Fatal("Incorrect number of output files.")
	}

	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3)+string(output4) != testFileString {
		t.Fatal("Incorrect output file content.")
	}

}

func TestPieceLineFileSplitterSplitWithEmptyFile(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineSplitter{1000}

	// Create a test file with 0 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

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

func TestPieceLineRoundRobinFileSplitterSplit2007Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineRoundRobinSplitter{3}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 669 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 669 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}

	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}

	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) == testFileString {
		t.Fatal("output files have correct content. should be different.")
	}

}

func TestPieceLineRoundRobinFileSplitterSplit713Lines(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineRoundRobinSplitter{4}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 713; i++ {
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
	if _, err := os.Stat("outputad"); os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("outputaa")
	if err != nil {
		t.Fatal(err)
	}

	if countLinesByByte(output1) != 179 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 178 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 178 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}
	output4, err := os.ReadFile("outputad")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output4) != 178 {
		t.Fatal("outputad file has incorrect number of lines. Expected 7, got ", countLinesByByte(output4))
	}

	if countFiles() != 4 {
		t.Fatal("Incorrect number of output files.")
	}

	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3)+string(output4) == testFileString {
		t.Fatal("output files have correct content. should be different.")
	}

}

func TestPieceLineRoundRobinFileSplitterSplitWithEmptyFile(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceLineRoundRobinSplitter{1000}

	// Create a test file with 0 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	err = splitter.Split(testFile, fileNameCreater)
	if err != nil {
		t.Fatal(err)
	}

	if countFiles() != 0 {
		t.Fatal("Incorrect number of output files.")
	}

}

func TestParseCHUNK(t *testing.T) {
	tests := []struct {
		input string
		want  chunk
		err   error
	}{
		{
			input: "10",
			want: chunk{
				N: 10,
			},
			err: nil,
		},
		{
			input: "5/10",
			want: chunk{
				K: 5,
				N: 10,
			},
			err: nil,
		},
		{
			input: "l/10",
			want: chunk{
				L: true,
				N: 10,
			},
			err: nil,
		},
		{
			input: "r/10",
			want: chunk{
				R: true,
				N: 10,
			},
			err: nil,
		},
		{
			input: "l/5/10",
			want: chunk{
				L: true,
				K: 5,
				N: 10,
			},
			err: nil,
		},
		{
			input: "r/5/10",
			want: chunk{
				R: true,
				K: 5,
				N: 10,
			},
			err: nil,
		},
		{
			input: "l",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "17/10",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "20ss/ll",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "20/ll",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "5/10/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "l/30/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "r/a/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "r/10/a",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "l/5/10/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "0/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "r/0/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "r/-2/20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "r/4/-20",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
		{
			input: "0",
			want:  chunk{},
			err:   fmt.Errorf(chunkFormatInvalidErrorMsg),
		},
	}

	for _, test := range tests {
		got, err := parseCHUNK(test.input)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("parseCHUNK(%q) error = %v, wantErr %v", test.input, err, test.err)
		}
		if got != test.want {
			t.Errorf("parseCHUNK(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestPieceSplitterSelectPieceByteSplitter(t *testing.T) {
	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	chunkStr := "3"
	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{chunkStr}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 5000 bytes to the file.
	byteCount := 5000
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1667 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1667 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 1666 {
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

func TestPieceSplitterSelectPieceByteSplitterAndStdout(t *testing.T) {
	defer deleteOutputFiles()
	buffer.Reset()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	chunkStr := "1/3"
	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{chunkStr}

	// Create a test file
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// Write 5000 bytes to the file.
	byteCount := 4999
	data := make([]byte, byteCount)
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
	out1Stat, err := os.Stat("outputaa")
	if os.IsNotExist(err) {
		t.Fatal("outputaa file was not created.")
	}
	// Check that the output files have the correct size
	if out1Stat.Size() != 1667 {
		t.Fatal("outputaa file has incorrect size. Expected 1024*1024, got ", out1Stat.Size())
	}
	out2Stat, err := os.Stat("outputab")
	if os.IsNotExist(err) {
		t.Fatal("outputab was not created.")
	}
	if out2Stat.Size() != 1667 {
		t.Fatal("outputab file has incorrect size. Expected 1024*1024, got ", out2Stat.Size())
	}
	out3Stat, err := os.Stat("outputac")
	if os.IsNotExist(err) {
		t.Fatal("outputac was not created.")
	}
	if out3Stat.Size() != 1665 {
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

	// Check that the output files stdout correct content
	if buffer.String() == string(output1) {
		t.Fatal("Incorrect output file content.")
	}
}

func TestPieceSplitterSelectPieceLineSplitter(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{"l/3"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 669 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 669 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
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

func TestPieceSplitterSelectPieceLineSplitterStdout3(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{"l/3/3"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 2005; i++ {
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 669 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 667 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
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

	// Check that the output files stdout correct content
	if buffer.String() == string(output3) {
		t.Fatal("Incorrect output file content.")
	}

}

func TestPieceSplitterSelectPieceRoundRobinLineSplitter(t *testing.T) {

	defer deleteOutputFiles()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{"r/3"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 669 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 669 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}

	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}

	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) == testFileString {
		t.Fatal("output files have correct content. should be different.")
	}

}

func TestPieceSplitterSelectPieceRoundRobinLineSplitterStdout2(t *testing.T) {

	defer deleteOutputFiles()
	buffer.Reset()
	// Create a mock fileNameCreater.
	fileNameCreater := AlphabetFileNameCreater{digit: 2, prefix: "output"}

	// Create a LineFileSplitter instance.
	splitter := PieceSplitter{"r/2/3"}

	// Create a test file with 2000 lines.
	testFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 2005; i++ {
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

	if countLinesByByte(output1) != 669 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 668 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 668 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}

	if countFiles() != 3 {
		t.Fatal("Incorrect number of output files.")
	}

	testFileContent, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	testFileString := string(testFileContent)
	if string(output1)+string(output2)+string(output3) == testFileString {
		t.Fatal("output files have correct content. should be different.")
	}

	// Check that the output files stdout correct content
	if buffer.String() == string(output2) {
		t.Fatal("Incorrect output file content.")
	}

}
