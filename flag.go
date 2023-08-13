package main

import (
	"flag"
	"fmt"
)

type FlagType int

const (
	UnknownFlag FlagType = iota
	LFlag
	NFlag
	BFlag
)

func ParseFlags() (string, FileSplitter, FileNameCreater, error) {
	var (
		lFlag              int64
		nFlag              string
		bFlag              string
		flagSet            int
		NumberFileNameFlag bool
		suffixLength       int
		flagType           = UnknownFlag
	)

	flag.Int64Var(&lFlag, "l", 0, "Line number for split file")
	flag.StringVar(&nFlag, "n", "", "CHUNKS for split file")
	flag.StringVar(&bFlag, "b", "", "Byte for split file")
	flag.BoolVar(&NumberFileNameFlag, "d", false, "Use numeric file name")
	flag.IntVar(&suffixLength, "a", 0, "Use numeric file name")
	flag.Parse()

	if lFlag != 0 {
		flagSet++
		flagType = LFlag
	}
	if nFlag != "" {
		flagSet++
		flagType = NFlag
	}
	if bFlag != "" {
		flagSet++
		flagType = BFlag
	}

	if flagSet > 1 {
		return "", nil, nil, fmt.Errorf(tooManyFlagErrorMsg)
	}
	var fileName string
	var prefix string
	if flag.NArg() == 1 {
		fileName = flag.Args()[0]
	} else if flag.NArg() == 2 {
		fileName = flag.Args()[0]
		prefix = flag.Args()[1]
	} else {
		return "", nil, nil, fmt.Errorf(invalidArgumentErrorMsg, flag.NArg())
	}

	var splitter FileSplitter
	if flagType == LFlag {
		splitter = LineSplitter{lFlag}
	} else if flagType == NFlag {
		splitter = PieceSplitter{nFlag}
	} else if flagType == BFlag {
		splitter = ByteSplitter{bFlag}
	} else {
		splitter = LineSplitter{1000}
	}

	var fileNameCreater FileNameCreater
	if suffixLength == 0 {
		suffixLength = 2
	}
	if NumberFileNameFlag {

		fileNameCreater = NumericFileNameCreater{suffixLength, prefix}
	} else {
		fileNameCreater = AlphabetFileNameCreater{suffixLength, prefix}
	}

	return fileName, splitter, fileNameCreater, nil

}
