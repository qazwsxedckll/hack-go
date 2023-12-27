package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name string
		file string
		want []string
	}{
		{
			name: "Add.asm",
			file: "testdata/Add/Add.asm",
			want: []string{
				`@2`,
				`D=A`,
				`@3`,
				`D=D+A`,
				`@0`,
				`M=D`,
			},
		},
		{
			name: "MaxL.asm",
			file: "testdata/Max/MaxL.asm",
			want: []string{
				`@0`,
				`D=M`,
				`@1`,
				`D=D-M`,
				`@12`,
				`D;JGT`,
				`@1`,
				`D=M`,
				`@2`,
				`M=D`,
				`@16`,
				`0;JMP`,
				`@0`,
				`D=M`,
				`@2`,
				`M=D`,
				`@16`,
				`0;JMP`,
			},
		},
		{
			name: "Max.asm",
			file: "testdata/Max/Max.asm",
			want: []string{
				`@R0`,
				`D=M`,
				`@R1`,
				`D=D-M`,
				`@ITSR0`,
				`D;JGT`,
				`@R1`,
				`D=M`,
				`@R2`,
				`M=D`,
				`@END`,
				`0;JMP`,
				`(ITSR0)`,
				`@R0`,
				`D=M`,
				`@R2`,
				`M=D`,
				`(END)`,
				`@END`,
				`0;JMP`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			require.NoError(t, err)

			p, err := NewParser(file)
			require.NoError(t, err)

			require.Equal(t, tt.want, p.commands)
		})
	}
}
