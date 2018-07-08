package test

import (
	"github.com/byxor/qbd/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChannelIsClosedWhenThereAreNoBytes(t *testing.T) {
	tokenChannel := make(chan tokens.Token)
	go tokens.Extract(tokenChannel, []byte{})
	_, more := readOneFrom(tokenChannel)
	assert.False(t, more)
}

func TestChannelIsClosedWhenFinished(t *testing.T) {
	tokenChannel := make(chan tokens.Token)
	bytes := []byte{0x00, 0x00, 0x00}

	go tokens.Extract(tokenChannel, bytes)

	for i := 0; i <= len(bytes); i++ {
		expectingMore := i < len(bytes)
		_, more := readOneFrom(tokenChannel)
		assert.Equal(t, expectingMore, more)
	}
}

func TestChannelIsClosedUponReceivingInvalidToken(t *testing.T) {
	tokenChannel := make(chan tokens.Token)
	bytes := []byte{0x16, 0x00}
	go tokens.Extract(tokenChannel, bytes)
	readOneFrom(tokenChannel)
	_, more := readOneFrom(tokenChannel)
	assert.False(t, more)
}

func TestExtractingTokens(t *testing.T) {
	entries := []struct {
		bytes    []byte
		expected tokens.Token
	}{
		{[]byte{0x00}, tokens.EndOfFile},
		{[]byte{0x01}, tokens.EndOfLine},
		{[]byte{0x03}, tokens.StartOfStruct},
		{[]byte{0x04}, tokens.EndOfStruct},
		{[]byte{0x05}, tokens.StartOfArray},
		{[]byte{0x06}, tokens.EndOfArray},

		{[]byte{0x07}, tokens.Assignment},

		{[]byte{0x09}, tokens.Comma},

		{[]byte{0x0A}, tokens.Subtraction},
		{[]byte{0x0B}, tokens.Addition},
		{[]byte{0x0C}, tokens.Division},
		{[]byte{0x0D}, tokens.Multiplication},

		{[]byte{0x0E}, tokens.StartOfExpression},
		{[]byte{0x0F}, tokens.EndOfExpression},

		{[]byte{0x11}, tokens.EqualityCheck},
		{[]byte{0x12}, tokens.LessThanCheck},
		{[]byte{0x13}, tokens.LessThanOrEqualCheck},
		{[]byte{0x14}, tokens.GreaterThanCheck},
		{[]byte{0x15}, tokens.GreaterThanOrEqualCheck},

		{[]byte{0x16, 0x00, 0x00, 0x00, 0x00}, tokens.Name},
		{[]byte{0x16, 0xBB, 0xEE, 0xEE, 0xFF}, tokens.Name},

		{[]byte{0x17, 0x00, 0x00, 0x00, 0x00}, tokens.Integer},
		{[]byte{0x17, 0xBA, 0x5E, 0xBA, 0x11}, tokens.Integer},

		{[]byte{0x1A, 0x00, 0x00, 0x00, 0x00}, tokens.Float},
		{[]byte{0x1A, 0x12, 0x34, 0x56, 0x78}, tokens.Float},

		{[]byte{0x23}, tokens.StartOfFunction},
		{[]byte{0x24}, tokens.EndOfFunction},
		{[]byte{0x29}, tokens.Return},

		{[]byte{0x22}, tokens.Break},

		{[]byte{0x25}, tokens.StartOfIf},
		{[]byte{0x26}, tokens.Else},
		{[]byte{0x27}, tokens.ElseIf},
		{[]byte{0x28}, tokens.EndOfIf},

		{[]byte{0x2B, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00}, tokens.ChecksumTableEntry},
		{[]byte{0x2B, 0x11, 0x22, 0x33, 0x44, 0x43, 0x6F, 0x63, 0x6B, 0x00}, tokens.ChecksumTableEntry},

		{[]byte{0x2C}, tokens.AllLocalReferences},
		{[]byte{0x2D}, tokens.LocalReference},

		{[]byte{0x39}, tokens.Not},

		{[]byte{0x3C}, tokens.StartOfSwitch},
		{[]byte{0x3D}, tokens.EndOfSwitch},
		{[]byte{0x3E}, tokens.SwitchCase},
		{[]byte{0x3F}, tokens.DefaultSwitchCase},

		{[]byte{0x42}, tokens.NamespaceAccess},

		{[]byte{0x47, 0x00, 0x00}, tokens.OptimisedIf},
		{[]byte{0x47, 0xDE, 0xAD}, tokens.OptimisedIf},

		{[]byte{0x48, 0x00, 0x00}, tokens.OptimisedElse},
		{[]byte{0x48, 0xDE, 0xAD}, tokens.OptimisedElse},

		{[]byte{0x49, 0x00, 0x00}, tokens.ShortJump},
		{[]byte{0x49, 0xBE, 0xEF}, tokens.ShortJump},

		// Invalid names (not enough bytes)
		{[]byte{0x16, 0x00, 0x00, 0x00}, tokens.Invalid},
		{[]byte{0x16, 0x11, 0x22}, tokens.Invalid},
		{[]byte{0x16, 0x33}, tokens.Invalid},
		{[]byte{0x16}, tokens.Invalid},

		// Invalid floats (not enough bytes)
		{[]byte{0x1A, 0x00, 0x00, 0x00}, tokens.Invalid},
		{[]byte{0x1A, 0x11, 0x22}, tokens.Invalid},
		{[]byte{0x1A, 0x33}, tokens.Invalid},
		{[]byte{0x1A}, tokens.Invalid},

		// Invalid integers (not enough bytes)
		{[]byte{0x17, 0x00, 0x00, 0x00}, tokens.Invalid},
		{[]byte{0x17, 0x11, 0x22}, tokens.Invalid},
		{[]byte{0x17, 0x33}, tokens.Invalid},
		{[]byte{0x17}, tokens.Invalid},

		// Invalid checksum table entries (not enough bytes)
		{[]byte{0x2B, 0x00, 0x00, 0x00, 0x00, 0x00}, tokens.Invalid},
		{[]byte{0x2B, 0xAB, 0xCD, 0xEF, 0x00}, tokens.Invalid},
		{[]byte{0x2B, 0x12, 0x34, 0x56}, tokens.Invalid},
		{[]byte{0x2B, 0xFF, 0xDE}, tokens.Invalid},
		{[]byte{0x2B, 0xE2}, tokens.Invalid},
		{[]byte{0x2B}, tokens.Invalid},

		// Invalid checksum table entries (not null-terminated)
		{[]byte{0x2B, 0x11, 0x22, 0x33, 0x44, 0x43, 0x6F, 0x63, 0x6B}, tokens.Invalid},
		{[]byte{0x2B, 0x00, 0xAA, 0xEE, 0xFF, 0x01, 0x02, 0x03}, tokens.Invalid},

		// Invalid shortjump (not enough bytes)
		{[]byte{0x49, 0x00}, tokens.Invalid},
		{[]byte{0x49}, tokens.Invalid},

		// Invalid optimised if (not enough bytes)
		{[]byte{0x47, 0x11}, tokens.Invalid},
		{[]byte{0x47}, tokens.Invalid},

		// Invalid optimised else (not enough bytes)
		{[]byte{0x48, 0xFE}, tokens.Invalid},
		{[]byte{0x48}, tokens.Invalid},
	}

	for _, entry := range entries {
		tokenChannel := make(chan tokens.Token)
		go tokens.Extract(tokenChannel, entry.bytes)
		token, _ := readOneFrom(tokenChannel)
		assert.Equal(t, entry.expected.String(), token.String())
	}
}

func TestExtractingMultipleTokens(t *testing.T) {
	entries := []struct {
		bytes  []byte
		output []tokens.Token
	}{
		{[]byte{0x01, 0x01}, []tokens.Token{tokens.EndOfLine, tokens.EndOfLine}},
		{[]byte{0x01, 0x00}, []tokens.Token{tokens.EndOfLine, tokens.EndOfFile}},

		{[]byte{0x01, 0x01, 0x01},
			[]tokens.Token{tokens.EndOfLine, tokens.EndOfLine, tokens.EndOfLine}},

		{[]byte{
			0x17, 0x00, 0x00, 0x00, 0x00,
			0x17, 0x01, 0x00, 0x00, 0x00,
		}, []tokens.Token{tokens.Integer, tokens.Integer}},

		{[]byte{
			0x16, 0xFF, 0x00, 0x00, 0xDD,
			0x2B, 0x11, 0x11, 0x11, 0x11, 0x68, 0x69, 0x00,
		}, []tokens.Token{tokens.Name, tokens.ChecksumTableEntry}},

		{[]byte{
			0x01,
			0x23,
			0x16, 0x93, 0x4D, 0xCD, 0xA1,
		}, []tokens.Token{tokens.EndOfLine, tokens.StartOfFunction, tokens.Name}},
	}

	for _, entry := range entries {
		tokenChannel := make(chan tokens.Token)
		go tokens.Extract(tokenChannel, entry.bytes)

		for _, expected := range entry.output {
			token, _ := readOneFrom(tokenChannel)
			assert.Equal(t, expected.String(), token.String())
		}
	}
}

func readOneFrom(tokens chan tokens.Token) (token tokens.Token, more bool) {
	select {
	case token, more = <-tokens:
		return
	case <-time.After(3 * time.Second):
		panic("Timed out while reading tokens...")
	}
}
