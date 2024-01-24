package main

import (
	"fmt"
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
				// TODO: Process .vm file and write to .asm file
			}
		}
	} else if strings.HasSuffix(arg, ".vm") {
		// TODO: Process .vm file and write to .asm file
	} else {
		panic("Invalid argument. Please provide a .vm file or a directory.")
	}
}
