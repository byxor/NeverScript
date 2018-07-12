package test

import (
	"github.com/byxor/qbd/table"
	. "github.com/byxor/qbd/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableGeneration(t *testing.T) {
	tokens := []Token{

		Token{ChecksumTableEntry,
			[]byte{0x2B, 0x00, 0x00, 0x00, 0x00,
				0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00}},

		Token{ChecksumTableEntry,
			[]byte{0x2B, 0x01, 0x00, 0x00, 0x00,
				0x48, 0x69, 0x00}},
	}

	entries := []struct {
		checksum int
		expected string
	}{
		{0, "Hello"},
		{1, "Hi"},
	}

	nameTable := table.GenerateUsing(tokens)

	for _, entry := range entries {
		name := nameTable.Get(entry.checksum)
		assert.Equal(t, entry.expected, name)
	}
}

func TestUnrecognisedChecksums(t *testing.T) {
	entries := []struct {
		checksum int
		expected string
	}{
		{0x00000000, "&00000000"},
		{0x11223344, "&44332211"},
		{0xABCDEF12, "&12efcdab"},
	}

	nameTable := table.GenerateUsing([]Token{})

	for _, entry := range entries {
		assert.Equal(t, entry.expected, nameTable.Get(entry.checksum))
	}
}
