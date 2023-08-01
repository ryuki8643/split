package main

import (
	"testing"
)

func TestCreateAlphabetFileName(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{2, 25, "xaz"},
		{3, 10, "xaak"},
		{4, 456, "xaaro"},
		{1, 0, "xa"},
		{2, 26, "xba"},
		{2, 27, "xbb"},
		{2, 53, "xcb"},
		{2, 675, "xzz"},
		{3, 100, "xadw"},
		{3, 702, "xbba"},
		{3, 703, "xbbb"},
		{4, 18278, "xbbba"},
		{4, 18279, "xbbbb"},
		{5, 456789, "xazzsv"},
		{6, 4567890, "xajzxgc"},
		{6, 999999, "xacexhn"},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,_ := createAlphabetFileName(tc.digit, tc.fileNumber)
		if got != tc.expectedFileName {
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, got)
		}
	}
}

func TestCreateAlphabetFileNameNegativeDigitError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{0, 25, ""},
		{-2, 10, ""},
		{-3, 456, ""},
		{-9999, 0, ""},
		{-2414124, 26, ""},
		{-1321, 27, ""},
		{-42, 53, ""},
		{-2, 675, ""},
		{-22, 100, ""},
		{-233, 702, ""},
		{-223, 703, ""},
		{-2414, 18278, ""},
		{-924148, 18279, ""},
		{-2, 456789, ""},
		{0, 4567890, ""},
		{0, 999999, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createAlphabetFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, negativeDigitErrorMsg,got)
		} else if err.Error() != negativeDigitErrorMsg {
			// Error should be negativeDigitErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateAlphabetFileNameNegativeFileNumberError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{2, -22, ""},
		{3, -1, ""},
		{4, -2, ""},
		{1, -21414, ""},
		{2, -2, ""},
		{2, -444, ""},
		{2, -53, ""},
		{2, -675, ""},
		{3, -100, ""},
		{3, -24944, ""},
		{3, -2494, ""},
		{4, -20494, ""},
		{4, -23004, ""},
		{5, -2094, ""},
		{6, -2003, ""},
		{6, -203, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createAlphabetFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, negativeFileNumberErrorMsg ,got)
		} else if err.Error() != negativeFileNumberErrorMsg {
			// Error should be negativeFileNumberErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateAlphabetFileNameTooBigFileNumberErrorMsg(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{1, 26, ""},
		{1, 27, ""},
		{2, 676, ""},
		{2, 677, ""},
		{3, 17576, ""},
		{3, 17577, ""},
		{4, 456976, ""},
		{4, 456977, ""},
		{5, 11881376, ""},
		{5, 11881377, ""},
		{6, 308915776, ""},
		{6, 308915777, ""},

	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createAlphabetFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, tooBigFileNumberErrorMsg ,got)
		} else if err.Error() != tooBigFileNumberErrorMsg {
			// Error should be tooBigFileNumberErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileName(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{2, 99, "x99"},
		{3, 10, "x010"},
		{4, 456, "x0456"},
		{1, 0, "x0"},
		{2, 9, "x09"},
		{2, 27, "x27"},
		{2, 53, "x53"},
		{3, 675, "x675"},
		{3, 100, "x100"},
		{3, 702, "x702"},
		{3, 703, "x703"},
		{4, 1278, "x1278"},
		{4, 1879, "x1879"},
		{5, 45789, "x45789"},
		{7, 4567890, "x4567890"},
		{6, 999999, "x999999"},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,_ := createNumericFileName(tc.digit, tc.fileNumber)
		if got != tc.expectedFileName {
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, got)
		}
	}
}

func TestCreateNumericFileNameNegativeDigitError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{0, 25, ""},
		{-2, 10, ""},
		{-3, 456, ""},
		{-9999, 0, ""},
		{-2414124, 26, ""},
		{-1321, 27, ""},
		{-42, 53, ""},
		{-2, 675, ""},
		{-22, 100, ""},
		{-233, 702, ""},
		{-223, 703, ""},
		{-2414, 18278, ""},
		{-924148, 18279, ""},
		{-2, 456789, ""},
		{0, 4567890, ""},
		{0, 999999, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createNumericFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, negativeDigitErrorMsg,got)
		} else if err.Error() != negativeDigitErrorMsg {
			// Error should be negativeDigitErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileNameNegativeFileNumberError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{2, -22, ""},
		{3, -1, ""},
		{4, -2, ""},
		{1, -21414, ""},
		{2, -2, ""},
		{2, -444, ""},
		{2, -53, ""},
		{2, -675, ""},
		{3, -100, ""},
		{3, -24944, ""},
		{3, -2494, ""},
		{4, -20494, ""},
		{4, -23004, ""},
		{5, -2094, ""},
		{6, -2003, ""},
		{6, -203, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createNumericFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, negativeFileNumberErrorMsg ,got)
		} else if err.Error() != negativeFileNumberErrorMsg {
			// Error should be negativeFileNumberErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileNameTooBigFileNumberErrorMsg(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		fileNumber int
		expectedFileName   string
	}{
		{1, 10, ""},
		{1, 11, ""},
		{2, 100, ""},
		{2, 101, ""},
		{3, 1000, ""},
		{3, 1001, ""},
		{4, 10000, ""},
		{4, 10001, ""},
		{5, 100000, ""},
		{5, 100001, ""},
		{6, 1000000, ""},
		{6, 1000001, ""},

	}

	// Execute each test case
	for _, tc := range testCases {
		got,err := createNumericFileName(tc.digit, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.fileNumber, tooBigFileNumberErrorMsg ,got)
		} else if err.Error() != tooBigFileNumberErrorMsg {
			// Error should be tooBigFileNumberErrorMsg
			t.Errorf("Input: %d, %d, Expected: %s, Got: %s", tc.digit, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}