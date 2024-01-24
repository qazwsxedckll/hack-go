package main

import "io"

type CodeWriter struct {
	fileName string
	w        io.WriteCloser
}

func NewCodeWriter(w io.WriteCloser) *CodeWriter {
	return &CodeWriter{
		w: w,
	}
}

func (c *CodeWriter) SetFileNmae(fileName string) {
	c.fileName = fileName
}

func (c *CodeWriter) WriteArithmetic(command string) {
}

func (c *CodeWriter) WritePushPop(command string, segment string, index int) {
}

func (c *CodeWriter) Close() {
	err := c.w.Close()
	if err != nil {
		panic(err)
	}
}
