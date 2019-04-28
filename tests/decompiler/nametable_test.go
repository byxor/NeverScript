package decompiler

import (
	"github.com/byxor/NeverScript/decompiler/nametable"
	. "github.com/byxor/NeverScript/decompiler/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableGeneration(t *testing.T) {
	tokens := []Token{

		Token{NameTableEntry,
			[]byte{0x2B, 0x00, 0x00, 0x00, 0x00,
				0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00}},

		Token{NameTableEntry,
			[]byte{0x2B, 0x00, 0x00, 0x00, 0x01,
				0x48, 0x69, 0x00}},

		Token{NameTableEntry,
			[]byte{0x2B, 0x00, 0x00, 0x00, 0xFC,
				0x54, 0x48, 0x55, 0x47, 0x50, 0x72, 0x6F, 0x00}},

		Token{NameTableEntry,
			[]byte{0x2B, 0xFF, 0xFF, 0xFF, 0xFF,
				0x41, 0x69, 0x72, 0x53, 0x74, 0x61, 0x74, 0x73, 0x00}},
	}

	entries := []struct {
		checksum string
		expected string
	}{
		{"00000000", "Hello"},
		{"00000001", "Hi"},
		{"000000fc", "THUGPro"},
		{"ffffffff", "AirStats"},
	}

	nameTable := nametable.BuildFrom(tokens)

	for _, entry := range entries {
		name := nameTable.Get(entry.checksum)
		assert.Equal(t, entry.expected, name)
	}
}

func TestUnrecognisedChecksums(t *testing.T) {
	entries := []struct {
		checksum string
		expected string
	}{
		{"00000000", "%00000000%"},
		{"11223344", "%11223344%"},
		{"abcdef12", "%abcdef12%"},
		{"ffeeffee", "%ffeeffee%"},
	}

	nameTable := nametable.BuildFrom([]Token{})

	for _, entry := range entries {
		assert.Equal(t, entry.expected, nameTable.Get(entry.checksum))
	}
}
