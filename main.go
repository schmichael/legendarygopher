package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/schmichael/legendarygopher/lg"
)

func main() {
	bind := "localhost:6060"
	flag.StringVar(&bind, "http", bind, "start web server")
	flag.Parse()
	if len(flag.Args()) < 1 {
		usageExit()
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file %q: %v\n", flag.Arg(0), err)
		os.Exit(11)
	}

	rc, err := progger(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting file size: %v\n", err)
		os.Exit(11)
	}

	if strings.HasSuffix(flag.Arg(0), ".gz") {
		if rc, err = gzip.NewReader(rc); err != nil {
			fmt.Fprintf(os.Stderr, "error decompressing %q: %v", flag.Arg(0), err)
			os.Exit(11)
		}
	}

	// Convert from cp437 to utf8
	reader := charmap.CodePage437.NewDecoder().Reader(rc)

	// Let's see how much memory it takes
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	alloc := m.Alloc

	world, err := lg.New(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading legends file %q: %v\n", flag.Arg(0), err, err)
		os.Exit(12)
	}
	rc.Close()
	f.Close()

	runtime.ReadMemStats(&m)
	fmt.Fprintf(os.Stderr, "took about %d MB\n", (m.Alloc-alloc)/1024/1024)

	if bind == "" {
		// Don't start web server; just exit
		fmt.Println(world)
		return
	}

	fmt.Printf("Open http://%s\n", bind)
	runserver(bind, world)
}

func usageExit() {
	fmt.Fprintf(os.Stderr, "incorrect usage, expected: %s dump.xml\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(10)
}
