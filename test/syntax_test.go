package test

import (
	"fmt"
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
		{
			[]Token{
				{Name, []byte{any, 0x00, 0x00, 0x00, 0x00}},
				{NameTableEntry, []byte{
					any, 0x00, 0x00, 0x00, 0x00,
					0x47, 0x45, 0x54, 0x44, 0x4F, 0x57, 0x4E, 0x00},
				},
			},
			"GETDOWN",
		},
		{
			[]Token{
				{Name, []byte{any, 0xFF, 0xEE, 0x00, 0xCC}},
				{NameTableEntry, []byte{
					any, 0xFF, 0xEE, 0x00, 0xCC,
					0x54, 0x75, 0x72, 0x62, 0x6F, 0x54, 0x69, 0x6D, 0x65, 0x00},
				},
			},
			"TurboTime",
		},

		// Addition
		{
			[]Token{
				{Integer, []byte{any, 0x09, 0x00, 0x00, 0x00}},
				{Addition, []byte{any}},
				{Integer, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
			},
			"9 + 10",
		},
		{
			[]Token{
				{Name, []byte{any, 0x09, 0x00, 0x00, 0x00}},
				{Addition, []byte{any}},
				{Name, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
				{NameTableEntry, []byte{any, 0x0A, 0x00, 0x00, 0x00, 0x66, 0x6F, 0x6F, 0x00}},
			},
			"%09000000% + foo",
		},
	}
	for _, entry := range entries {
		code := code.GenerateUsing(entry.tokens)
		fmt.Println(entry.tokens)
		fmt.Println(code)
		assert.Equal(t, entry.expected, code)
	}
}
