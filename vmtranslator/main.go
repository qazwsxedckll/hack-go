package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: vmtranslator <file.vm>|<directory>")
		return
	}

	arg := os.Args[1]
	err := filepath.WalkDir(arg, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".vm") {
			in, err := os.Open(path)
			if err != nil {
				return err
			}
			parse(in, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func parse(in io.Reader, path string) {
	out, err := os.Create(path[:len(path)-3] + ".asm")
	if err != nil {
		panic(err)
	}

	codeWriter := NewCodeWriter(out)
	codeWriter.SetFileNmae(strings.TrimSuffix(filepath.Base(path), ".vm"))
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
