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

func parse(in io.ReadSeeker, out io.Writer) error {
	symbolTable := NewSymbolTable()

	parser, err := NewParser(in)
	if err != nil {
		return err
	}

	rom := 0
	for parser.HasMoreCommands() {
		parser.Advance()

		switch parser.CommandType() {
		case ACommand, CCommand:
			rom++
		case LCommand:
			symbolTable.AddEntry(parser.Symbol(), rom)
		}
	}

	_, err = in.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	parser, err = NewParser(in)
	if err != nil {
		return err
	}

	ram := 16
	for parser.HasMoreCommands() {
		parser.Advance()

		switch symbol := parser.Symbol(); parser.CommandType() {
		case ACommand:
			i, err := strconv.Atoi(symbol)
			if err != nil {
				if !symbolTable.Contains(symbol) {
					symbolTable.AddEntry(symbol, ram)
					ram++
				}
				i = symbolTable.GetAddress(symbol)
			}
			fmt.Fprintf(out, "0%015b\n", i)
		case CCommand:
			fmt.Fprintf(out, "111%s%s%s\n", Comp(parser.Comp()), Dest(parser.Dest()), Jump(parser.Jump()))
		}
	}
	return nil
}
