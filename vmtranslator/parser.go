package main

import (
	"bufio"
	"io"
	"strings"
)

type CommandType string

const (
	C_ARITHMETIC CommandType = "C_ARITHMETIC"
	C_PUSH       CommandType = "C_PUSH"
	C_POP        CommandType = "C_POP"
	C_LABEL      CommandType = "C_LABEL"
	C_GOTO       CommandType = "C_GOTO"
	C_IF         CommandType = "C_IF"
	C_FUNCTION   CommandType = "C_FUNCTION"
	C_RETURN     CommandType = "C_RETURN"
	C_CALL       CommandType = "C_CALL"
)

type Parser struct {
	commands []string
	current  int
}

func NewParser(r io.Reader) (*Parser, error) {
	s := bufio.NewScanner(r)

	var commands []string
	for s.Scan() {
		text := s.Text()

		i := strings.Index(text, "//")
		if i != -1 {
			text = text[:i]
		}

		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		commands = append(commands, text)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return &Parser{
		commands: commands,
		current:  -1,
	}, nil
}

func (p *Parser) HasMoreCommands() bool {
	return p.current+1 < len(p.commands)
}

func (p *Parser) Advance() {
	p.current++
}

func (p *Parser) CommandType() CommandType {
	// if strings.HasPrefix(p.commands[p.current], "@") {
	// 	return ACommand
	// } else if strings.HasPrefix(p.commands[p.current], "(") {
	// 	return LCommand
	// } else {
	// 	return CCommand
	// }
}

func (p *Parser) Arg1() string {
}

func (p *Parser) Arg2() int {
}
