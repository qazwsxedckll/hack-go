package main

import (
	"bytes"
	"io"
	"strconv"
)

type CodeWriter struct {
	fileName string
	w        io.WriteCloser
	jump     int
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
	buffer := bytes.Buffer{}

	switch command {
	case "add":
		popBinary(&buffer)
		buffer.WriteString(`M=M+D`)
	case "sub":
		popBinary(&buffer)
		buffer.WriteString(`M=M-D`)
	case "neg":
		popBinary(&buffer)
		buffer.WriteString(`M=-M`)
	case "eq":
		c.jump = writeCompare(&buffer, c.jump, "JEQ")
	case "gt":
		c.jump = writeCompare(&buffer, c.jump, "JGT")
	case "lt":
		c.jump = writeCompare(&buffer, c.jump, "JLT")
	case "and":
		popBinary(&buffer)
		buffer.WriteString(`M=M&D`)
	case "or":
		popBinary(&buffer)
		buffer.WriteString(`M=M|D`)
	case "not":
		popUnary(&buffer)
		buffer.WriteString(`M=!M`)
	}

	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WritePushPop(command CommandType, segment string, index int) {
	buffer := bytes.Buffer{}
	switch command {
	case "push":
		switch segment {
		case "constant":
			buffer.WriteString(`@` + strconv.Itoa(index))
			buffer.WriteString(`D=A`)
			push(&buffer)
		case "local":
		case "argument":
		case "this":
		case "that":
		case "temp":
		case "pointer":
		case "static":
		}
	case "pop":
	}

	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) Close() {
	err := c.w.Close()
	if err != nil {
		panic(err)
	}
}

func writeFalse(buf *bytes.Buffer) {
	buf.WriteString(`@SP`)
	buf.WriteString(`A=M-1`)
	buf.WriteString(`M=0`)
}

func writeTrue(buf *bytes.Buffer) {
	buf.WriteString(`@SP`)
	buf.WriteString(`A=M-1`)
	buf.WriteString(`M=-1`)
}

func writeCompare(buffer *bytes.Buffer, jump int, compare string) int {
	popBinary(buffer)
	buffer.WriteString(`D=M-D`)
	buffer.WriteString(`@JUMP_` + strconv.Itoa(jump))
	buffer.WriteString(`D;` + compare)
	writeFalse(buffer)
	buffer.WriteString(`@END_` + strconv.Itoa(jump))
	buffer.WriteString(`0;JMP`)
	buffer.WriteString(`(JUMP_` + strconv.Itoa(jump) + `)`)
	writeTrue(buffer)
	buffer.WriteString(`(END_` + strconv.Itoa(jump) + `)`)
	return jump + 1
}

// one value in D, one value in M
func popBinary(buf *bytes.Buffer) {
	buf.WriteString(`@SP`)
	buf.WriteString(`AM=M-1`)
	buf.WriteString(`D=M`)
	buf.WriteString(`A=A-1`)
}

// one value in M
func popUnary(buf *bytes.Buffer) {
	buf.WriteString(`@SP`)
	buf.WriteString(`A=M-1`)
}

func push(buf *bytes.Buffer) {
	buf.WriteString(`@SP`)
	buf.WriteString(`A=M`)
	buf.WriteString(`M=D`)
	buf.WriteString(`@SP`)
	buf.WriteString(`M=M+1`)
}
