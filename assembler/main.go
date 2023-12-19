package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: assembler <file.asm>")
		return
	}

	err := parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
}

func parse(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s.hack", strings.Split(file.Name(), ".")[0]))
	if err != nil {
		return err
	}

	fmt.Printf("out: %v\n", out)

	return nil
}
