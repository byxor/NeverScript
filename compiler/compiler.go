package compiler

import (
	"github.com/alecthomas/participle"
	"github.com/byxor/NeverScript/compiler/grammar"
	"github.com/byxor/NeverScript/shared/tokens"
)

func Compile(code string) ([]byte, error) {
	syntaxTree, err := parseCodeIntoSyntaxTree(code)
	if err != nil {
		return []byte{}, err
	}

	bytecode := make([]byte, 500)
	index := 0
	push := func(bytes ...byte) {
		for i, b := range bytes {
			bytecode[index+i] = b
		}
		index += len(bytes)
	}

	push(tokens.EndOfLine)

	for _, declaration := range syntaxTree.Declarations {
		if declaration.EndOfLine != nil {
			push(tokens.EndOfLine)
		} else if declaration.BooleanAssignment != nil {
			dontCare := byte(0xFF)
			nameChecksum := []byte{dontCare, dontCare, dontCare, dontCare}
			trueOrFalse := declaration.BooleanAssignment.Boolean.Value

			var assignmentValue byte
			if trueOrFalse == "true" {
				assignmentValue = 0x01
			} else {
				assignmentValue = 0x00
			}

			push(0x16)
			push(nameChecksum...)
			push(0x07)
			push(0x17, assignmentValue, 0x00, 0x00, 0x00)
		}
	}

	push(tokens.EndOfFile)

	return bytecode[0:index], nil
}

func parseCodeIntoSyntaxTree(code string) (*grammar.SyntaxTree, error) {
	parser := participle.MustBuild(
		&grammar.SyntaxTree{},
		participle.UseLookahead(2),
	)

	syntaxTree := &grammar.SyntaxTree{}

	err := parser.ParseString(code, syntaxTree)
	if err != nil {
		return nil, err
	}

	return syntaxTree, nil
}
