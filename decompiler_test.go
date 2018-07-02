package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecompiler(t *testing.T) {
	type TestEntry struct {
		bytes    []byte
		expected string
	}

	var entries = []TestEntry{
		{[]byte{}, ""},

		// File endings
		{[]byte{0x00}, ""},

		// Line endings
		{[]byte{0x01}, ";"},

		// Integers
		{[]byte{0x17, 0x00, 0x00, 0x00, 0x00}, "0"},
		{[]byte{0x17, 0x0A, 0x00, 0x00, 0x00}, "10"},
		{[]byte{0x17, 0xCC, 0xDD, 0xEE, 0xFF}, "-1122868"},
	}

	for _, entry := range entries {
		assert.Equal(t, entry.expected, Decompile(entry.bytes))
	}
}
