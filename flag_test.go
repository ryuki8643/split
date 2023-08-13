package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {

	var (
		testlogfile string
		testpaniconexit0 bool
		testtimeout string
		testrun string
	)
	// set up test cases
	testCases := []struct {
		args []string
		err  error
	}{
		{
			args: []string{"input.txt"},
			err:  nil,
		},
		{
			args: []string{"input.txt", "output_"},
			err:  nil,
		},
		{
			args: []string{"-l", "100", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-n", "10", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-b", "100K", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-d", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-a", "3", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-l", "100", "-n", "10", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-l", "100", "-b", "100K", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-n", "10", "-b", "100K", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-l", "100", "-n", "10", "-b", "100K", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-l", "100", "-d", "-a", "3", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-n", "10", "-d", "-a", "3", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-b", "100K", "-d", "-a", "3", "input.txt"},
			err:  nil,
		},
		{
			args: []string{"-l", "100", "-n", "10", "-d", "-a", "3", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-l", "100", "-b", "100K", "-d", "-a", "3", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-n", "10", "-b", "100K", "-d", "-a", "3", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"-l", "100", "-n", "10", "-b", "100K", "-d", "-a", "3", "input.txt"},
			err:  fmt.Errorf(tooManyFlagErrorMsg),
		},
		{
			args: []string{"input.txt","prefix","l","100"},
			err:  fmt.Errorf(invalidArgumentErrorMsg,4),
		},
		
	}

	// run test cases
	for _, tc := range testCases {
		testing.Init()

		// reset flag set
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		tc.args = append([]string{"-test.v"}, tc.args...)
		// set test args
		os.Args = tc.args

		flag.StringVar(&testlogfile, "test.testlogfile", "", "set to write test log to this file")
		flag.BoolVar(&testpaniconexit0, "test.paniconexit0", false, "set to cause a panic instead of an exit (for testing)")
		flag.StringVar(&testtimeout, "test.timeout", "10m", "if positive, sets an aggregate time limit for all tests")
		flag.StringVar(&testrun, "test.run", "", "regular expression to select tests to run")

		// call function
		_, _, _, err := ParseFlags()
		// check error
		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected error %v but got %v for args %v", tc.err, err, tc.args[1:])
		}
	}
}


