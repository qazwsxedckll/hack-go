package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"strconv"
)

var validLabel = `^[A-Za-z_:.][\w_:.]*$`

type CodeWriter struct {
	fileName string
	w        io.WriteCloser
	jump     int
	retIndex int
	function string
}

func NewCodeWriter(w io.WriteCloser) *CodeWriter {
	return &CodeWriter{
		w: w,
	}
}

func (c *CodeWriter) SetFileName(fileName string) {
	c.fileName = fileName
}

func (c *CodeWriter) WriteArithmetic(command string) {
	buffer := bytes.Buffer{}

	switch command {
	case "add":
		popUnaryAndGetTop(&buffer)
		buffer.WriteString("M=M+D\n")
	case "sub":
		popUnaryAndGetTop(&buffer)
		buffer.WriteString("M=M-D\n")
	case "neg":
		buffer.WriteString("@SP\n")
		buffer.WriteString("A=M-1\n")
		buffer.WriteString("M=-M\n")
	case "eq":
		c.jump = writeCompare(&buffer, c.jump, "JEQ")
	case "gt":
		c.jump = writeCompare(&buffer, c.jump, "JGT")
	case "lt":
		c.jump = writeCompare(&buffer, c.jump, "JLT")
	case "and":
		popUnaryAndGetTop(&buffer)
		buffer.WriteString("M=M&D\n")
	case "or":
		popUnaryAndGetTop(&buffer)
		buffer.WriteString("M=M|D\n")
	case "not":
		buffer.WriteString("@SP\n")
		buffer.WriteString("A=M-1\n")
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

func (c *CodeWriter) WriteInit() {
	buffer := bytes.Buffer{}
	buffer.WriteString("@256\n")
	buffer.WriteString("D=A\n")
	buffer.WriteString("@SP\n")
	buffer.WriteString("M=D\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	c.WriteCall("Sys.init", 0)
}

func (c *CodeWriter) WriteLabel(label string) {
	reg := regexp.MustCompile(validLabel)
	if !reg.MatchString(label) {
		panic("invalid label: " + label)
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("(" + c.function + "$" + label + ")\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WriteGoto(label string) {
	buffer := bytes.Buffer{}
	buffer.WriteString("@" + c.function + "$" + label + "\n")
	buffer.WriteString("0;JMP\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WriteIf(label string) {
	buffer := bytes.Buffer{}
	popUnary(&buffer)
	buffer.WriteString("@" + c.function + "$" + label + "\n")
	buffer.WriteString("D;JNE\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WriteCall(functionName string, numArgs int) {
	buffer := bytes.Buffer{}
	returnAddress := "_" + functionName + "_RETURN_" + strconv.Itoa(c.retIndex)
	c.retIndex++
	// push return address
	buffer.WriteString("@" + returnAddress + "\n")
	buffer.WriteString("D=A\n")
	push(&buffer)
	// push LCL, ARG, THIS, THAT
	buffer.WriteString("@LCL\n")
	buffer.WriteString("D=M\n")
	push(&buffer)
	buffer.WriteString("@ARG\n")
	buffer.WriteString("D=M\n")
	push(&buffer)
	buffer.WriteString("@THIS\n")
	buffer.WriteString("D=M\n")
	push(&buffer)
	buffer.WriteString("@THAT\n")
	buffer.WriteString("D=M\n")
	push(&buffer)
	// ARG = SP - (n + 5)
	buffer.WriteString("@" + strconv.Itoa(numArgs+5) + "\n")
	buffer.WriteString("D=D-A\n")
	buffer.WriteString("@ARG\n")
	buffer.WriteString("M=D\n")
	// LCL = SP
	buffer.WriteString("@SP\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@LCL\n")
	buffer.WriteString("M=D\n")
	// goto function
	buffer.WriteString("@" + functionName + "\n")
	buffer.WriteString("0;JMP\n")
	// return address
	buffer.WriteString("(" + returnAddress + ")\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WriteReturn() {
	buffer := bytes.Buffer{}
	// FRAME = LCL
	buffer.WriteString("@LCL\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@R13\n")
	buffer.WriteString("M=D\n")
	// RET = *(FRAME - 5)
	// if arg num is 0, pop() will overwrite RET, so we need to save it
	buffer.WriteString("@5\n")
	buffer.WriteString("A=D-A\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@R14\n")
	buffer.WriteString("M=D\n")
	// *ARG = pop()
	popUnary(&buffer)
	buffer.WriteString("@ARG\n")
	buffer.WriteString("A=M\n")
	buffer.WriteString("M=D\n")
	// SP = ARG + 1
	buffer.WriteString("@ARG\n")
	buffer.WriteString("D=M+1\n")
	buffer.WriteString("@SP\n")
	buffer.WriteString("M=D\n")
	// THAT = *(FRAME - 1)
	buffer.WriteString("@R13\n")
	buffer.WriteString("AM=M-1\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@THAT\n")
	buffer.WriteString("M=D\n")
	// THIS = *(FRAME - 2)
	buffer.WriteString("@R13\n")
	buffer.WriteString("AM=M-1\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@THIS\n")
	buffer.WriteString("M=D\n")
	// ARG = *(FRAME - 3)
	buffer.WriteString("@R13\n")
	buffer.WriteString("AM=M-1\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@ARG\n")
	buffer.WriteString("M=D\n")
	// LCL = *(FRAME - 4)
	buffer.WriteString("@R13\n")
	buffer.WriteString("AM=M-1\n")
	buffer.WriteString("D=M\n")
	buffer.WriteString("@LCL\n")
	buffer.WriteString("M=D\n")
	// goto RET
	buffer.WriteString("@R14\n")
	buffer.WriteString("A=M\n")
	buffer.WriteString("0;JMP\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (c *CodeWriter) WriteFunction(functionName string, numLocals int) {
	reg := regexp.MustCompile(validLabel)
	if !reg.MatchString(functionName) {
		panic("invalid function name: " + functionName)
	}

	buffer := bytes.Buffer{}
	c.function = functionName
	buffer.WriteString("(" + functionName + ")\n")
	_, err := c.w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	for i := 0; i < numLocals; i++ {
		c.WritePushPop(C_PUSH, "constant", 0)
	}
}

func (c *CodeWriter) Close() {
	err := c.w.Close()
	if err != nil {
		if !errors.Is(err, fs.ErrClosed) {
			panic(err)
		}
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
	popUnaryAndGetTop(buffer)
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

// pop value and put it in D, then point to the last value
func popUnaryAndGetTop(buf *bytes.Buffer) {
	popUnary(buf)
	buf.WriteString("A=A-1\n")
}

// pop value and put it in D
func popUnary(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("AM=M-1\n")
	buf.WriteString("D=M\n")
}

func push(buf *bytes.Buffer) {
	buf.WriteString("@SP\n")
	buf.WriteString("A=M\n")
	buf.WriteString("M=D\n")
	buf.WriteString("@SP\n")
	buf.WriteString("M=M+1\n")
}
