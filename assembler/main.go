package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: assembler <file.asm>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}

	out, err := os.Create(fmt.Sprintf("%s.hack", strings.Split(file.Name(), ".")[0]))
	if err != nil {
		fmt.Printf("error creating output file: %v\n", err)
		return
	}

	err = parse(file, out)
	if err != nil {
		fmt.Printf("error parsing file: %v\n", err)
		return
	}
}

func parse(in io.Reader, out io.Writer) error {
	parser, err := NewParser(in)
	if err != nil {
		return err
	}

	for parser.HasMoreCommands() {
		parser.Advance()

		switch parser.CommandType() {
		case ACommand:
			i, err := strconv.Atoi(parser.Symbol())
			if err != nil {
				return fmt.Errorf("error converting symbol to int: %v", err)
			}
			fmt.Fprintf(out, "0%015b\n", i)
		case CCommand:
			fmt.Fprintf(out, "111%s%s%s\n", Comp(parser.Comp()), Dest(parser.Dest()), Jump(parser.Jump()))
		}
	}

	return nil
}
