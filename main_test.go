package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {

	defer deleteOutputPrefixXFiles()
	testing.Init()

	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"testfile.txt"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
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
	main()

	// Check that the output files were created.
	if _, err := os.Stat("xaa"); os.IsNotExist(err) {
		t.Fatal("xaa file was not created.")
	}
	if _, err := os.Stat("xab"); os.IsNotExist(err) {
		t.Fatal("xab was not created.")
	}
	if _, err := os.Stat("xac"); os.IsNotExist(err) {
		t.Fatal("xac was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("xaa")
	if err != nil {
		t.Fatal(err)
	}

	if countLinesByByte(output1) != 1000 {
		t.Fatal("xaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("xab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 1000 {
		t.Fatal("xab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("xac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 7 {
		t.Fatal("xac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}

	if countFilesPrefixX() != 3 {
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

func TestMainLineFileSplitterSplit1000Lines(t *testing.T) {

	defer deleteOutputFiles()
	testing.Init()

	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-l", "1000", "testfile.txt", "output"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
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
	main()

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

func TestMainByteFileSplitterSplit2500Byte(t *testing.T) {

	defer deleteOutputFiles()
	testing.Init()

	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-b", "1k", "testfile.txt", "output"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
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
	main()

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

func TestMainPieceSplitterSelectPieceByteSplitter(t *testing.T) {
	defer deleteOutputFiles()
	testing.Init()

	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-n", "3", "testfile.txt", "output"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
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
	main()

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

func TestMainPieceSplitterSelectPieceLineSplitter(t *testing.T) {

	defer deleteOutputFiles()
	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-n", "l/3", "testfile.txt", "output"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 3007; i++ {
		fmt.Fprintln(testFile, "line", i)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	main()

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

	if countLinesByByte(output1) != 1003 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("outputab")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 1003 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("outputac")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 1001 {
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

func TestMainPieceSplitterSelectPieceRoundRobinLineSplitter(t *testing.T) {

	defer deleteOutputFiles()
	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-n", "r/3", "testfile.txt", "output"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 2003; i++ {
		fmt.Fprintln(testFile, "line", i)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	main()

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

	if countLinesByByte(output1) != 668 {
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
	if string(output1)+string(output2)+string(output3) == testFileString {
		t.Fatal("output files have correct content. should be different.")
	}

}

func TestMainNoPrefixAndNumericFileName(t *testing.T) {

	defer deleteOutputPrefixXFiles()
	testing.Init()

	// reset flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	args := []string{"-l", "700", "-d", "-a", "3", "testfile.txt"}

	args = append([]string{"-test.v"}, args...)
	// set test args
	os.Args = args

	// Create a test file with 2000 lines.
	testFile, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()
	for i := 1; i <= 2002; i++ {
		fmt.Fprintln(testFile, "line", i)
	}

	// Reset the file pointer to the beginning.
	if _, err := testFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// Call the Split function with the mock fileNameCreater.
	main()

	// Check that the output files were created.
	if _, err := os.Stat("x000"); os.IsNotExist(err) {
		t.Fatal("x000 file was not created.")
	}
	if _, err := os.Stat("x001"); os.IsNotExist(err) {
		t.Fatal("x001 was not created.")
	}
	if _, err := os.Stat("x002"); os.IsNotExist(err) {
		t.Fatal("x002 was not created.")
	}

	// Check that the output files have the correct lines
	output1, err := os.ReadFile("x000")
	if err != nil {
		t.Fatal(err)
	}

	if countLinesByByte(output1) != 700 {
		t.Fatal("outputaa file has incorrect number of lines.Expected 1000, got ", countLinesByByte(output1))
	}

	output2, err := os.ReadFile("x001")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output2) != 700 {
		t.Fatal("outputab file has incorrect number of lines. Expected 1000, got ", countLinesByByte(output2))
	}
	output3, err := os.ReadFile("x002")
	if err != nil {
		t.Fatal(err)
	}
	if countLinesByByte(output3) != 602 {
		t.Fatal("outputac file has incorrect number of lines. Expected 7, got ", countLinesByByte(output3))
	}

	if countFilesPrefixX() != 3 {
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

func countFilesPrefixX() int {
	outputPrefix := "x" // Prefix of files to be counted

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

func deleteOutputPrefixXFiles() {
	outputPrefix := "x" // Prefix of files to be deleted

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
