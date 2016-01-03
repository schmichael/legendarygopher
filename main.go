package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/text/encoding/charmap"

	"github.com/schmichael/legendarygopher/lg"
)

func main() {
	bind := ""
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

	progf, err := progger(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting file size: %v\n", err)
		os.Exit(11)
	}

	// Convert from cp437 to utf8
	decr := charmap.CodePage437.NewDecoder().Reader(progf)

	// Let's see how much memory it takes
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	alloc := m.Alloc

	world, err := lg.New(decr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading legends file %q: %v\n", flag.Arg(0), err, err)
		os.Exit(12)
	}
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
