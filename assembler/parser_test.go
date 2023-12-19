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
			name: "Max.asm",
			file: "testdata/Max/Max.asm",
			want: []string{
				`@R0`,
				`D=M`,
				`@R1`,
				`D=D-M`,
				`@OUTPUT_FIRST`,
				`D;JGT`,
				`@R1`,
				`D=M`,
				`@OUTPUT_D`,
				`0;JMP`,
				`(OUTPUT_FIRST)`,
				`@R0`,
				`D=M`,
				`(OUTPUT_D)`,
				`@R2`,
				`M=D`,
				`(INFINITE_LOOP)`,
				`@INFINITE_LOOP`,
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

// func TestCommandType(t *testing.T) {
// 	file, err := os.Open("testdata/Add.asm")
// 	require.NoError(t, err)

// 	p := NewParser(file)

// 	p.Advance()
// 	require.Equal(t, ACommand, p.CommandType())

// 	p.Advance()
// 	require.Equal(t, CCommand, p.CommandType())

// 	p.Advance()
// 	require.Equal(t, LCommand, p.CommandType())

// 	p.Advance()
// 	require.Equal(t, CCommand, p.CommandType())
// }
