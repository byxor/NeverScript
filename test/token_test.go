package test

import (
	"github.com/byxor/qbd/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGettingTokens(t *testing.T) {
	type Entry struct {
		input  []byte
		output token.Token
	}

	entries := []Entry{
		{[]byte{}, token.Invalid},

		{[]byte{0x00}, token.EndOfFile},

		{[]byte{0x01}, token.EndOfLine},

		{[]byte{0x16, 0x00, 0x00, 0x00, 0x00}, token.Name},
		{[]byte{0x16, 0xBB, 0xEE, 0xEE, 0xFF}, token.Name},

		{[]byte{0x16, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x16, 0x11, 0x22}, token.Invalid},
		{[]byte{0x16, 0x33}, token.Invalid},
		{[]byte{0x16}, token.Invalid},
	}

	for _, entry := range entries {
		tokens := make(chan token.Token)
		go token.GetTokens(tokens, entry.input)
		assert.Equal(t, entry.output, <-tokens)
	}
}
