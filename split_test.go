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
