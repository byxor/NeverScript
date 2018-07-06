package test

import (
	"github.com/byxor/qbd/token"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChannelIsClosedWhenInputIsEmpty(t *testing.T) {
	tokens := make(chan token.Token)

	go token.GetTokens(tokens, []byte{})

	_, more := read(tokens)

	assert.Equal(t, false, more)
}

func TestChannelIsClosedWhenFinished(t *testing.T) {
	tokens := make(chan token.Token)
	input := []byte{0x00, 0x00, 0x00}

	go token.GetTokens(tokens, input)

	for i := 0; i <= len(input); i++ {
		expectingMore := i < len(input)
		_, more := read(tokens)
		assert.Equal(t, expectingMore, more)
	}
}

func TestExtractingTokens(t *testing.T) {
	entries := []struct {
		input    []byte
		expected token.Token
	}{
		{[]byte{0x00}, token.EndOfFile},
		{[]byte{0x01}, token.EndOfLine},
		{[]byte{0x03}, token.StartOfStruct},
		{[]byte{0x04}, token.EndOfStruct},
		{[]byte{0x05}, token.StartOfArray},
		{[]byte{0x06}, token.EndOfArray},

		{[]byte{0x23}, token.StartOfFunction},
		{[]byte{0x24}, token.EndOfFunction},
		{[]byte{0x29}, token.Return},

		{[]byte{0x22}, token.Break},

		{[]byte{0x25}, token.StartOfIf},
		{[]byte{0x26}, token.Else},
		{[]byte{0x27}, token.ElseIf},
		{[]byte{0x28}, token.EndOfIf},

		{[]byte{0x07}, token.Assignment},

		{[]byte{0x0A}, token.Subtraction},
		{[]byte{0x0B}, token.Addition},
		{[]byte{0x0C}, token.Division},
		{[]byte{0x0D}, token.Multiplication},

		{[]byte{0x11}, token.EqualityCheck},
		{[]byte{0x12}, token.LessThanCheck},
		{[]byte{0x13}, token.LessThanOrEqualCheck},
		{[]byte{0x14}, token.GreaterThanCheck},
		{[]byte{0x15}, token.GreaterThanOrEqualCheck},

		{[]byte{0x16, 0x00, 0x00, 0x00, 0x00}, token.Name},
		{[]byte{0x16, 0xBB, 0xEE, 0xEE, 0xFF}, token.Name},

		{[]byte{0x17, 0x00, 0x00, 0x00, 0x00}, token.Integer},
		{[]byte{0x17, 0xBA, 0x5E, 0xBA, 0x11}, token.Integer},

		{[]byte{0x1A, 0x00, 0x00, 0x00, 0x00}, token.Float},
		{[]byte{0x1A, 0x12, 0x34, 0x56, 0x78}, token.Float},

		{[]byte{0x2B, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00}, token.ChecksumTableEntry},
		{[]byte{0x2B, 0x11, 0x22, 0x33, 0x44, 0x43, 0x6F, 0x63, 0x6B, 0x00}, token.ChecksumTableEntry},

		// Invalid names (not enough bytes)
		{[]byte{0x16, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x16, 0x11, 0x22}, token.Invalid},
		{[]byte{0x16, 0x33}, token.Invalid},
		{[]byte{0x16}, token.Invalid},

		// Invalid floats (not enough bytes)
		{[]byte{0x1A, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x1A, 0x11, 0x22}, token.Invalid},
		{[]byte{0x1A, 0x33}, token.Invalid},
		{[]byte{0x1A}, token.Invalid},

		// Invalid integers (not enough bytes)
		{[]byte{0x17, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x17, 0x11, 0x22}, token.Invalid},
		{[]byte{0x17, 0x33}, token.Invalid},
		{[]byte{0x17}, token.Invalid},

		// Invalid checksum table entries (not enough bytes)
		{[]byte{0x2B, 0x00, 0x00, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x2B, 0xAB, 0xCD, 0xEF, 0x00}, token.Invalid},
		{[]byte{0x2B, 0x12, 0x34, 0x56}, token.Invalid},
		{[]byte{0x2B, 0xFF, 0xDE}, token.Invalid},
		{[]byte{0x2B, 0xE2}, token.Invalid},
		{[]byte{0x2B}, token.Invalid},

		// Invalid checksum table entries (not null-terminated)
		{[]byte{0x2B, 0x11, 0x22, 0x33, 0x44, 0x43, 0x6F, 0x63, 0x6B}, token.Invalid},
		{[]byte{0x2B, 0x00, 0xAA, 0xEE, 0xFF, 0x01, 0x02, 0x03}, token.Invalid},
	}

	for _, entry := range entries {
		tokens := make(chan token.Token)
		go token.GetTokens(tokens, entry.input)

		token, _ := read(tokens)
		assert.Equal(t, entry.expected, token)
	}
}

func TestExtractingMultipleTokens(t *testing.T) {
	entries := []struct {
		input  []byte
		output []token.Token
	}{
		{
			[]byte{0x01, 0x01},
			[]token.Token{token.EndOfLine, token.EndOfLine},
		},
		{
			[]byte{0x01, 0x00},
			[]token.Token{token.EndOfLine, token.EndOfFile},
		},
		{
			[]byte{
				0x17, 0x00, 0x00, 0x00, 0x00,
				0x17, 0x01, 0x00, 0x00, 0x00,
			},
			[]token.Token{token.Integer, token.Integer},
		},
		{
			[]byte{
				0x16, 0xFF, 0x00, 0x00, 0xDD,
				0x2B, 0x11, 0x11, 0x11, 0x11, 0x68, 0x69, 0x00,
			},
			[]token.Token{token.Name, token.ChecksumTableEntry},
		},
		{
			[]byte{
				0x01,
				0x23,
				0x16, 0x93, 0x4D, 0xCD, 0xA1,
			},
			[]token.Token{token.EndOfLine, token.StartOfFunction, token.Name},
		},
	}

	for _, entry := range entries {
		tokens := make(chan token.Token)
		go token.GetTokens(tokens, entry.input)

		for _, expected := range entry.output {
			token, _ := read(tokens)
			assert.Equal(t, expected, token)
		}
	}
}

func read(tokens chan token.Token) (token token.Token, more bool) {
	select {
	case token, more = <-tokens:
		return
	case <-time.After(3 * time.Second):
		panic("Timed out while reading token...")
	}
}
