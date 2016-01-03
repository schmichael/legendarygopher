package main

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

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

	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting file size: %v\n", err)
		os.Exit(11)
	}

	rc := progger(f, fi.Size())

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

	var dec lg.Decoder = xml.NewDecoder(reader)
	if strings.HasSuffix(flag.Arg(0), ".json") || strings.HasSuffix(flag.Arg(0), ".json.gz") {
		dec = json.NewDecoder(reader)
	}

	start := time.Now()
	world, err := lg.New(dec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading legends file %q: %v\n", flag.Arg(0), err, err)
		os.Exit(12)
	}
	dur := time.Now().Sub(start)
	rc.Close()
	f.Close()

	runtime.ReadMemStats(&m)
	fmt.Fprintf(os.Stderr, "took %s (%d KBps) and approximately %d MB of memory\n",
		dur, (fi.Size()/1024)/int64(math.Max(1, float64(dur/time.Second))), (m.Alloc-alloc)/1024/1024)

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
