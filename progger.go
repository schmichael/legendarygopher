package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type p struct {
	f    *os.File
	sz   int64
	cur  int64
	last time.Time
}

func progger(f *os.File) (io.ReadCloser, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &p{f: f, sz: fi.Size(), last: time.Now()}, nil
}

func (p *p) Read(buf []byte) (int, error) {
	n, err := p.f.Read(buf)
	p.cur += int64(n)
	now := time.Now()
	if now.After(p.last.Add(3 * time.Second)) {
		fmt.Fprintf(os.Stderr, "%d/%d (%d%%) done\n", p.cur, p.sz, int(float64(p.cur)/float64(p.sz)*100))
		p.last = now
	}
	return n, err
}

func (p *p) Close() error { return p.f.Close() }
