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
		{[]byte{0x00}, token.EndOfFile},
		{[]byte{0x01}, token.EndOfLine},
	}

	for _, entry := range entries {
		tokens := make(chan token.Token)
		go token.GetTokens(tokens, entry.input)
		assert.Equal(t, entry.output, <-tokens)
	}
}
