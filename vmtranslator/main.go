package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: vmtranslator <file.vm>|<directory>")
		return
	}

	arg := os.Args[1]

	if info, err := os.Stat(arg); err == nil && info.IsDir() {
		files, err := os.ReadDir(arg)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".vm") {
				in, err := os.Open(arg + "/" + file.Name())
				if err != nil {
					panic(err)
				}
				out, err := os.Create(arg + "/" + file.Name()[:len(file.Name())-3] + ".asm")
				if err != nil {
					panic(err)
				}
				parse(in, out, out.Name())
			}
		}
	} else if strings.HasSuffix(arg, ".vm") {
		in, err := os.Open(arg)
		if err != nil {
			panic(err)
		}
		out, err := os.OpenFile(arg[:len(arg)-3]+".asm", os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic(err)
		}
		parse(in, out, out.Name())
	} else {
		panic("Invalid argument. Please provide a .vm file or a directory.")
	}
}

func parse(in io.Reader, out io.WriteCloser, filename string) {
	codeWriter := NewCodeWriter(out)
	codeWriter.SetFileNmae(filename)
	defer codeWriter.w.Close()

	parser, err := NewParser(in)
	if err != nil {
		panic(err)
	}

	for parser.HasMoreCommands() {
		parser.Advance()

		cmdType := parser.CommandType()
		switch cmdType {
		case C_ARITHMETIC:
			codeWriter.WriteArithmetic(parser.Arg1())
		case C_PUSH, C_POP:
			codeWriter.WritePushPop(cmdType, parser.Arg1(), parser.Arg2())
		default:
			panic(fmt.Sprintf("unknown command type: %v", cmdType))
		}
	}
}
