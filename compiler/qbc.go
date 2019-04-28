package qbc

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/repr"
	"github.com/byxor/NeverScript/compiler/checksums"
	"github.com/byxor/NeverScript/compiler/grammar"
)

func Compile(code string) ([]byte, error) {

	qbFile, err := parse(code)
	if err != nil {
		return []byte{}, err
	}

	bytecode := make([]byte, 500)
	marker := 0

	push := func(bytes ...byte) {
		for i, b := range bytes {
			bytecode[marker+i] = b
		}
		marker += len(bytes)
	}

	for _, statement := range qbFile.Statements {
		if statement.EndOfLine != nil {
			push(0x01)
		} else if statement.Assignment != nil {
			name := []byte{0x66, 0x6F, 0x6F, 0x00}
			checksum := checksums.Generate(string(name))
			checksumBytes := checksums.LittleEndian(checksum)

			push(0x016)
			push(checksumBytes...)
			push(0x07, 0x0A, 0x00, 0x00, 0x00)
			push(0x01)
			push(0x2B)
			push(checksumBytes...)
			push(name...)
		}

	}

	push(0x00)
	return bytecode[0:marker], nil
}

func parse(code string) (*grammar.QbFile, error) {
	parser := participle.MustBuild(
		&grammar.QbFile{},
		participle.UseLookahead(2),
	)

	result := &grammar.QbFile{}

	err := parser.ParseString(code, result)
	if err != nil {
		return nil, err
	}

	repr.Println(result, repr.Indent("  "), repr.OmitEmpty(true))

	return result, nil
}
