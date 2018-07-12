package test

import (
	"github.com/byxor/qbd/code"
	. "github.com/byxor/qbd/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

const any = 0x00

func TestSyntax(t *testing.T) {
	entries := []struct {
		tokens   []Token
		expected string
	}{
		{[]Token{}, ""},

		{[]Token{{EndOfFile, nil}}, ""},

		// Ends of lines
		{[]Token{{EndOfLine, nil}}, ";"},
		{[]Token{{EndOfLine, nil}, {EndOfLine, nil}}, ";"},
		{[]Token{{EndOfLine, nil}, {EndOfLine, nil}, {EndOfLine, nil}}, ";"},

		// Integers
		{[]Token{{Integer, []byte{any, 0x00, 0x00, 0x00, 0x00}}}, "0"},
		{[]Token{{Integer, []byte{any, 0x01, 0x00, 0x00, 0x00}}}, "1"},
		{[]Token{{Integer, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}}}, "-1"},
		{[]Token{{Integer, []byte{any, 0x2A, 0x43, 0x0F, 0x00}}}, "1000234"},

		// Unknown Names
		{[]Token{{Name, []byte{any, 0x11, 0x22, 0x33, 0x44}}}, "%11223344%"},
		{[]Token{{Name, []byte{any, 0x12, 0x00, 0x00, 0x99}}}, "%12000099%"},
		{[]Token{{Name, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}}}, "%ffffffff%"},

		// Known Names
		{[]Token{
			{Name, []byte{any, 0x00, 0x00, 0x00, 0x00}},
			{NameTableEntry, []byte{
				any, 0x00, 0x00, 0x00, 0x00,
				0x47, 0x45, 0x54, 0x44, 0x4F, 0x57, 0x4E, 0x00},
			},
		}, "GETDOWN"},
	}
	for _, entry := range entries {
		code := code.GenerateUsing(entry.tokens)
		assert.Equal(t, entry.expected, code)
	}
}