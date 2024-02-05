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

	if strings.HasSuffix(arg, ".vm") {
		out, err := os.Create(strings.TrimSuffix(filepath.Base(arg), ".vm") + ".asm")
		if err != nil {
			panic(err)
		}

		codeWriter := NewCodeWriter(out)
		defer codeWriter.Close()
		codeWriter.WriteInit()

		in, err := os.Open(arg)
		if err != nil {
			panic(err)
		}

		parse(in, codeWriter, arg)
	} else {
		vmFound := false
		var codeWriter *CodeWriter

		err := filepath.WalkDir(arg, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".vm") {
				if !vmFound {
					dir := filepath.Dir(path)
					out, err := os.Create(dir + string(filepath.Separator) + filepath.Base(dir) + ".asm")
					if err != nil {
						panic(err)
					}
					codeWriter = NewCodeWriter(out)
					codeWriter.WriteInit()
					vmFound = true
				}

				in, err := os.Open(path)
				if err != nil {
					return err
				}
				parse(in, codeWriter, path)
			}

			if d.IsDir() {
				if codeWriter != nil {
					codeWriter.Close()
				}
				vmFound = false
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func parse(in io.Reader, codeWriter *CodeWriter, path string) {
	codeWriter.SetFileName(strings.TrimSuffix(filepath.Base(path), ".vm"))

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
		case C_LABEL:
			codeWriter.WriteLabel(parser.Arg1())
		case C_GOTO:
			codeWriter.WriteGoto(parser.Arg1())
		case C_IF:
			codeWriter.WriteIf(parser.Arg1())
		case C_FUNCTION:
			codeWriter.WriteFunction(parser.Arg1(), parser.Arg2())
		case C_RETURN:
			codeWriter.WriteReturn()
		case C_CALL:
			codeWriter.WriteCall(parser.Arg1(), parser.Arg2())
		default:
			panic(fmt.Sprintf("unknown command type: %v", cmdType))
		}
	}
}
