package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoSymbol(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{
			name: "Add.asm",
			file: "testdata/Add/Add.asm",
		},
		{
			name: "MaxL.asm",
			file: "testdata/Max/MaxL.asm",
		},
		{
			name: "PongL.asm",
			file: "testdata/Pong/PongL.asm",
		},
		{
			name: "RectL.asm",
			file: "testdata/Rect/RectL.asm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			require.NoError(t, err)

			buf := bytes.Buffer{}

			err = parse(file, &buf)
			require.NoError(t, err)

			want, err := os.ReadFile(tt.file[:len(tt.file)-4] + ".hack")
			require.NoError(t, err)
			require.Equal(t, string(want), buf.String())
		})
	}
}

func TestWithSymbol(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{
			name: "Max.asm",
			file: "testdata/Max/Max.asm",
		},
		{
			name: "Pong.asm",
			file: "testdata/Pong/Pong.asm",
		},
		{
			name: "Rect.asm",
			file: "testdata/Rect/Rect.asm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			require.NoError(t, err)

			buf := bytes.Buffer{}

			err = parse(file, &buf)
			require.NoError(t, err)

			want, err := os.ReadFile(tt.file[:len(tt.file)-4] + ".hack")
			require.NoError(t, err)
			require.Equal(t, string(want), buf.String())
		})
	}
}
