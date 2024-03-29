package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	CommandType       string
	ArithmeticCommand string
)

var cmdMap = map[string]CommandType{
	"add":      C_ARITHMETIC,
	"sub":      C_ARITHMETIC,
	"neg":      C_ARITHMETIC,
	"eq":       C_ARITHMETIC,
	"gt":       C_ARITHMETIC,
	"lt":       C_ARITHMETIC,
	"and":      C_ARITHMETIC,
	"or":       C_ARITHMETIC,
	"not":      C_ARITHMETIC,
	"push":     C_PUSH,
	"pop":      C_POP,
	"label":    C_LABEL,
	"goto":     C_GOTO,
	"if-goto":  C_IF,
	"function": C_FUNCTION,
	"return":   C_RETURN,
	"call":     C_CALL,
}

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
	cmd := strings.Fields(p.commands[p.current])[0]
	if c, ok := cmdMap[cmd]; ok {
		return c
	}

	panic(fmt.Sprintf("unknown command: %s", p.commands[p.current]))
}

func (p *Parser) Arg1() string {
	if p.CommandType() == C_RETURN {
		panic("no argument for return command")
	}

	if p.CommandType() == C_ARITHMETIC {
		return p.commands[p.current]
	}

	return strings.Fields(p.commands[p.current])[1]
}

func (p *Parser) Arg2() int {
	if p.CommandType() == C_PUSH || p.CommandType() == C_POP || p.CommandType() == C_FUNCTION || p.CommandType() == C_CALL {
		i, err := strconv.Atoi(strings.Fields(p.commands[p.current])[2])
		if err != nil {
			panic(err)
		}

		return i
	}

	panic("no second argument for command")
}
