package main

import "io"

type closer struct {
	r io.Reader
	c io.Closer
}

func (c *closer) Read(p []byte) (int, error) { return c.r.Read(p) }
func (c *closer) Close() error               { return c.c.Close() }
