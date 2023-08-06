package main

import (
	"testing"
)

func TestCreateAlphabetFileName(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{2,"", 25, "xaz"},
		{3,"a", 10, "aaak"},
		{4,"b", 456, "baaro"},
		{1,"c", 0, "ca"},
		{2, "d",26, "dba"},
		{2,"dada", 27, "dadabb"},
		{2,"hshara.fhnfaf", 53, "hshara.fhnfafcb"},
		{2,"sa", 675, "sazz"},
		{3, "sfaf",100, "sfafadw"},
		{3,"sad", 702, "sadbba"},
		{3, "fhhtrnvxcx",703, "fhhtrnvxcxbbb"},
		{4, "123345",18278, "123345bbba"},
		{4, "x",18279, "xbbbb"},
		{5, "ogjsg",456789, "ogjsgazzsv"},
		{6, "",4567890, "xajzxgc"},
		{6, "geslgposejpgjspeojgjhpgjghjpag",999999, "geslgposejpgjspeojgjhpgjghjpagacexhn"},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = AlphabetFileName{}
		got,_ := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if got != tc.expectedFileName {
			t.Errorf("Input: %d, %s,%d, Expected: %s, Got: %s", tc.digit, tc.prefix,tc.fileNumber, tc.expectedFileName, got)
		}
	}
}

func TestCreateAlphabetFileNameNegativeDigitError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{0, "", 25, ""},
		{-2, "a", 10, ""},
		{-3, "b", 456, ""},
		{-9999, "c", 0, ""},
		{-2414124, "d", 26, ""},
		{-1321, "dada", 27, ""},
		{-42, "hshara.fhnfaf", 53, ""},
		{-2, "sa", 675, ""},
		{-22, "sfaf", 100, ""},
		{-233, "sad", 702, ""},
		{-223, "fhhtrnvxcx", 703, ""},
		{-2414, "123345", 18278, ""},
		{-924148, "x", 18279, ""},
		{-2, "ogjsg", 456789, ""},
		{0, "", 4567890, ""},
		{0, "geslgposejpgjspeojgjhpgjghjpag", 999999, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = AlphabetFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.prefix,tc.fileNumber, negativeDigitErrorMsg,got)
		} else if err.Error() != negativeDigitErrorMsg {
			// Error should be negativeDigitErrorMsg
			t.Errorf("Input: %d, %s %d, Expected: %s, Got: %s", tc.digit,tc.prefix, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateAlphabetFileNameNegativeFileNumberError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{2, "", -22, ""},
		{3, "a", -1, ""},
		{4, "b", -2, ""},
		{1, "c", -21414, ""},
		{2, "d", -2, ""},
		{2, "dada", -444, ""},
		{2, "hshara.fhnfaf", -53, ""},
		{2, "sa", -675, ""},
		{3, "sfaf", -100, ""},
		{3, "sad", -24944, ""},
		{3, "fhhtrnvxcx", -2494, ""},
		{4, "123345", -20494, ""},
		{4, "x", -23004, ""},
		{5, "ogjsg", -2094, ""},
		{6, "", -2003, ""},
		{6, "geslgposejpgjspeojgjhpgjghjpag", -203, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = AlphabetFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
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
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{1, "", 26, ""},
		{1, "a", 27, ""},
		{2, "b", 676, ""},
		{2, "c", 677, ""},
		{3, "d", 17576, ""},
		{3, "dada", 17577, ""},
		{4, "hshara.fhnfaf", 456976, ""},
		{4, "sa", 456977, ""},
		{5, "sfaf", 11881376, ""},
		{5, "sad", 11881377, ""},
		{6, "fhhtrnvxcx", 308915776, ""},
		{6, "123345", 308915777, ""},

	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = AlphabetFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got no error and return: %s", tc.digit, tc.prefix,tc.fileNumber, tooBigFileNumberErrorMsg ,got)
		} else if err.Error() != tooBigFileNumberErrorMsg {
			// Error should be tooBigFileNumberErrorMsg
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got: %s", tc.digit, tc.prefix,tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileName(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{2, "", 99, "x99"},
		{3, "a", 10, "a010"},
		{4, "b", 456, "b0456"},
		{1, "c", 0, "c0"},
		{2, "d", 9, "d09"},
		{2, "dada", 27, "dada27"},
		{2, "hshara.fhnfaf", 53, "hshara.fhnfaf53"},
		{3, "sa", 675, "sa675"},
		{3, "sfaf", 100, "sfaf100"},
		{3, "sad", 702, "sad702"},
		{3, "fhhtrnvxcx", 703, "fhhtrnvxcx703"},
		{4, "123345", 1278, "1233451278"},
		{4, "x", 1879, "x1879"},
		{5, "ogjsg", 45789, "ogjsg45789"},
		{7, "", 4567890, "x4567890"},
		{6, "geslgposejpgjspeojgjhpgjghjpag", 999999, "geslgposejpgjspeojgjhpgjghjpag999999"},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = NumericFileName{}
		got,_ := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if got != tc.expectedFileName {
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got: %s", tc.digit, tc.prefix,tc.fileNumber, tc.expectedFileName, got)
		}
	}
}

func TestCreateNumericFileNameNegativeDigitError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{0, "", 25, ""},
		{-2, "a", 10, ""},
		{-3, "b", 456, ""},
		{-9999, "c", 0, ""},
		{-2414124, "d", 26, ""},
		{-1321, "dada", 27, ""},
		{-42, "hshara.fhnfaf", 53, ""},
		{-2, "sa", 675, ""},
		{-22, "sfaf", 100, ""},
		{-233, "sad", 702, ""},
		{-223, "fhhtrnvxcx", 703, ""},
		{-2414, "123345", 18278, ""},
		{-924148, "x", 18279, ""},
		{-2, "ogjsg", 456789, ""},
		{0, "", 4567890, ""},
		{0, "geslgposejpgjspeojgjhpgjghjpag", 999999, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = NumericFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got no error and return: %s", tc.digit,tc.prefix, tc.fileNumber, negativeDigitErrorMsg,got)
		} else if err.Error() != negativeDigitErrorMsg {
			// Error should be negativeDigitErrorMsg
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got: %s", tc.digit,tc.prefix, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileNameNegativeFileNumberError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{2, "", -22, ""},
		{3, "a", -1, ""},
		{4, "b", -2, ""},
		{1, "c", -21414, ""},
		{2, "d", -2, ""},
		{2, "dada", -444, ""},
		{2, "hshara.fhnfaf", -53, ""},
		{2, "sa", -675, ""},
		{3, "sfaf", -100, ""},
		{3, "sad", -24944, ""},
		{3, "fhhtrnvxcx", -2494, ""},
		{4, "123345", -20494, ""},
		{4, "x", -23004, ""},
		{5, "ogjsg", -2094, ""},
		{6, "", -2003, ""},
		{6, "geslgposejpgjspeojgjhpgjghjpag", -203, ""},
	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = NumericFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got no error and return: %s", tc.digit,tc.prefix, tc.fileNumber, negativeFileNumberErrorMsg ,got)
		} else if err.Error() != negativeFileNumberErrorMsg {
			// Error should be negativeFileNumberErrorMsg
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got: %s", tc.digit, tc.prefix, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}


func TestCreateNumericFileNameTooBigFileNumberErrorMsg(t *testing.T) {
	// Define test cases
	testCases := []struct {
		digit int
		prefix string
		fileNumber int
		expectedFileName   string
	}{
		{1, "", 10, ""},
		{1, "a", 11, ""},
		{2, "b", 100, ""},
		{2, "c", 101, ""},
		{3, "d", 1000, ""},
		{3, "dada", 1001, ""},
		{4, "hshara.fhnfaf", 10000, ""},
		{4, "sa", 10001, ""},
		{5, "sfaf", 100000, ""},
		{5, "sad", 100001, ""},
		{6, "fhhtrnvxcx", 1000000, ""},
		{6, "123345", 1000001, ""},

	}

	// Execute each test case
	for _, tc := range testCases {
		var fileNameCreater FileNameCreater = NumericFileName{}
		got,err := fileNameCreater.Create(tc.digit,tc.prefix, tc.fileNumber)
		if err == nil {
			// Error should be thrown
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got no error and return: %s", tc.digit,tc.prefix, tc.fileNumber, tooBigFileNumberErrorMsg ,got)
		} else if err.Error() != tooBigFileNumberErrorMsg {
			// Error should be tooBigFileNumberErrorMsg
			t.Errorf("Input: %d, %s, %d, Expected: %s, Got: %s", tc.digit,tc.prefix, tc.fileNumber, tc.expectedFileName, err.Error())
		}
	}
}