package main

import (
	"bytes"
	"fmt"
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
		buffer.WriteString("M=M+D\n")
	case "sub":
		popBinary(&buffer)
		buffer.WriteString("M=M-D\n")
	case "neg":
		popUnary(&buffer)
		buffer.WriteString("M=-M\n")
	case "eq":
		c.jump = writeCompare(&buffer, c.jump, "JEQ")
	case "gt":
		c.jump = writeCompare(&buffer, c.jump, "JGT")
	case "lt":
		c.jump = writeCompare(&buffer, c.jump, "JLT")
	case "and":
		popBinary(&buffer)
		buffer.WriteString("M=M&D\n")
	case "or":
		popBinary(&buffer)
		buffer.WriteString("M=M|D\n")
	case "not":
		popUnary(&buffer)
		buffer.WriteString("M=!M\n")
	}

	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WritePushPop(command CommandType, segment string, index int) {
	buffer := bytes.Buffer{}
	switch command {
	case C_PUSH:
		switch segment {
		case "constant":
			buffer.WriteString("@" + strconv.Itoa(index) + "\n")
			buffer.WriteString("D=A\n")
		case "local", "argument", "this", "that":
			switch segment {
			case "local":
				buffer.WriteString("@LCL\n")
			case "argument":
				buffer.WriteString("@ARG\n")
			case "this":
				buffer.WriteString("@THIS\n")
			case "that":
				buffer.WriteString("@THAT\n")
			}
			writeOffset(&buffer, index)
		case "pointer", "temp", "static":
			switch segment {
			case "pointer":
				if index == 0 {
					buffer.WriteString("@THIS\n")
				} else if index == 1 {
					buffer.WriteString("@THAT\n")
				} else {
					panic("invalid pointer index")
				}
			case "temp":
				buffer.WriteString("@R" + strconv.Itoa(5+index) + "\n")
			case "static":
				buffer.WriteString("@" + c.fileName + "." + strconv.Itoa(index) + "\n")
			}
			buffer.WriteString("D=M\n")
		}

		push(&buffer)
	case C_POP:
		switch segment {
		case "local", "argument", "this", "that":
			switch segment {
			case "local":
				buffer.WriteString("@LCL\n")
			case "argument":
				buffer.WriteString("@ARG\n")
			case "this":
				buffer.WriteString("@THIS\n")
			case "that":
				buffer.WriteString("@THAT\n")
			}

			buffer.WriteString("D=M\n")
			buffer.WriteString("@" + strconv.Itoa(index) + "\n")
			buffer.WriteString("D=D+A\n")
		case "pointer", "temp", "static":
			switch segment {
			case "pointer":
				if index == 0 {
					buffer.WriteString("@THIS\n")
				}
				if index == 1 {
					buffer.WriteString("@THAT\n")
				}
			case "temp":
				buffer.WriteString("@R" + strconv.Itoa(5+index) + "\n")
			case "static":
				buffer.WriteString("@" + c.fileName + "." + strconv.Itoa(index) + "\n")
			}

			buffer.WriteString("D=A\n")
		}
		buffer.WriteString("@R13\n")
		buffer.WriteString("M=D\n")

		popUnary(&buffer)

		buffer.WriteString("D=M\n")
		buffer.WriteString("@R13\n")
		buffer.WriteString("A=M\n")
		buffer.WriteString("M=D\n")
	default:
		panic(fmt.Sprintf("unknown command type: %v", command))
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
	buf.WriteString("@SP\n")
	buf.WriteString("A=M-1\n")
	buf.WriteString("M=0\n")
}

func writeTrue(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("A=M-1\n")
	buf.WriteString("M=-1\n")
}

func writeCompare(buffer *bytes.Buffer, jump int, compare string) int {
	popBinary(buffer)
	buffer.WriteString("D=M-D\n")
	buffer.WriteString("@JUMP_" + strconv.Itoa(jump) + "\n")
	buffer.WriteString("D;" + compare + "\n")
	writeFalse(buffer)
	buffer.WriteString("@END_" + strconv.Itoa(jump) + "\n")
	buffer.WriteString("0;JMP\n")
	buffer.WriteString("(JUMP_" + strconv.Itoa(jump) + ")\n")
	writeTrue(buffer)
	buffer.WriteString("(END_" + strconv.Itoa(jump) + ")\n")
	return jump + 1
}

// get value from segment[index] and put it in D
func writeOffset(buf *bytes.Buffer, index int) {
	buf.WriteString("D=M\n")
	buf.WriteString("@" + strconv.Itoa(index) + "\n")
	buf.WriteString("A=D+A\n")
	buf.WriteString("D=M\n")
}

// one value in D, one value in M
func popBinary(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("AM=M-1\n")
	buf.WriteString("D=M\n")
	buf.WriteString("A=A-1\n")
}

// one value in M
func popUnary(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("AM=M-1\n")
}

func push(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("A=M\n")
	buf.WriteString("M=D\n")
	buf.WriteString("@SP\n")
	buf.WriteString("M=M+1\n")
}
