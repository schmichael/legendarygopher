package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/schmichael/legendarygopher/lg"
)

const invalidUTFError = `invalid UTF-8 characters
try running:
    iconv -f utf-8 -t utf-8 -c %q > out.xml
    %s out.xml
`

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		usageExit()
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file %q: %v\n", flag.Arg(0), err)
		os.Exit(11)
	}
	progf, err := progger(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting file size: %v\n", err)
		os.Exit(11)
	}

	lgr, err := lg.New(progf)
	if err != nil {
		if xmlerr, ok := err.(*xml.SyntaxError); ok && xmlerr.Msg == "invalid UTF-8" {
			fmt.Fprintf(os.Stderr, invalidUTFError, flag.Arg(0), os.Args[0])
			os.Exit(12)
		}
		fmt.Fprintf(os.Stderr, "error reading legends file %q: %v\n", flag.Arg(0), err, err)
		os.Exit(12)
	}
	f.Close()

	fmt.Println(lgr)
}

func usageExit() {
	fmt.Fprintf(os.Stderr, "incorrect usage, expected: %s dump.xml\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(10)
}
