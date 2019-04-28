package compiler

import (
	"github.com/alecthomas/participle"
	// "github.com/alecthomas/repr"
	"github.com/byxor/NeverScript/compiler/grammar"
)

func Compile(code string) ([]byte, error) {

	program, err := parse(code)
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

	push(0x01)

	for _, declaration := range program.Declarations {
		if declaration.EndOfLine != nil {
			push(0x01)
		} else if declaration.BooleanAssignment != nil {
			dontCare := byte(0xFF)

			nameChecksum := []byte{dontCare, dontCare, dontCare, dontCare}

			push(0x16)
			push(nameChecksum...)
			push(0x07, 0x17, 0x01, 0x00, 0x00, 0x00)
		}
	}

	push(0x00)
	return bytecode[0:marker], nil
}

func parse(code string) (*grammar.Program, error) {
	parser := participle.MustBuild(
		&grammar.Program{},
		participle.UseLookahead(2),
	)

	result := &grammar.Program{}

	err := parser.ParseString(code, result)
	if err != nil {
		return nil, err
	}

	// repr.Println(result, repr.Indent("  "), repr.OmitEmpty(true))

	return result, nil
}
