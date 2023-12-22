package main

import (
	"bufio"
	"io"
	"strings"
)

type CommandType string

const (
	ACommand CommandType = "A_COMMAND"
	CCommand CommandType = "C_COMMAND"
	LCommand CommandType = "L_COMMAND"
)

type Parser struct {
	commands []string
	current  int
}

// TODO: Constants must be non-negative and are written in decimal
// notation. A user-defined symbol can be any sequence of letters, digits, underscore (_),
// dot (.), dollar sign ($), and colon (:) that does not begin with a digit

// TODO: Command validation
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
	if strings.HasPrefix(p.commands[p.current], "@") {
		return ACommand
	} else if strings.HasPrefix(p.commands[p.current], "(") {
		return LCommand
	} else {
		return CCommand
	}
}

func (p *Parser) Symbol() string {
	if p.commands[p.current][len(p.commands[p.current])-1] == ')' {
		return p.commands[p.current][1 : len(p.commands[p.current])-1]
	} else {
		return p.commands[p.current][1:]
	}
}

func (p *Parser) Dest() string {
	s := strings.Split(p.commands[p.current], "=")
	if len(s) == 2 {
		return s[0]
	}

	return ""
}

func (p *Parser) Comp() string {
	equal := strings.Index(p.commands[p.current], "=")
	semicolon := strings.Index(p.commands[p.current], ";")

	// dest=comp;jump
	if equal != -1 && semicolon != -1 {
		return p.commands[p.current][equal+1 : semicolon]
	}

	// comp;jump
	if equal == -1 && semicolon != -1 {
		return p.commands[p.current][:semicolon]
	}

	// dest=comp
	if equal != -1 && semicolon == -1 {
		return p.commands[p.current][equal+1:]
	}

	return ""
}

func (p *Parser) Jump() string {
	s := strings.Split(p.commands[p.current], ";")
	if len(s) == 2 {
		return s[1]
	}

	return ""
}
