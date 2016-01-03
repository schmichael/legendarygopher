package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type p struct {
	rc   io.ReadCloser
	sz   int64
	cur  int64
	last time.Time
}

func progger(rc io.ReadCloser, sz int64) io.ReadCloser {
	return &p{rc: rc, sz: sz, last: time.Now()}
}

func (p *p) Read(buf []byte) (int, error) {
	n, err := p.rc.Read(buf)
	p.cur += int64(n)
	now := time.Now()
	if now.After(p.last.Add(3 * time.Second)) {
		fmt.Fprintf(os.Stderr, "%d/%d (%d%%) done\n", p.cur, p.sz, int(float64(p.cur)/float64(p.sz)*100))
		p.last = now
	}
	return n, err
}

func (p *p) Close() error { return p.rc.Close() }
