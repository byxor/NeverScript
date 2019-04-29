package compiler

import (
	"github.com/alecthomas/participle"
	"github.com/byxor/NeverScript/compiler/grammar"
	"github.com/byxor/NeverScript/shared/tokens"
	"github.com/pkg/errors"
	"encoding/binary"
)

const (
	bytecodeSize = 10 * 1000 * 1000
	junkByte = 0xFF
)

func Compile(code string) ([]byte, error) {
	syntaxTree, err := parseCodeIntoSyntaxTree(code)

	if err != nil {
		return []byte{}, errors.Wrap(err, "Failed to get syntax tree")
	}

	bytecode := make([]byte, bytecodeSize)
	numberOfUsedBytes := 0

	pushBytes := func(bytes ...byte) {
		for i, b := range bytes {
			bytecode[numberOfUsedBytes+i] = b
		}
		numberOfUsedBytes += len(bytes)
	}

	pushBytes(tokens.EndOfLine)

	for _, declaration := range syntaxTree.Declarations {
		if declaration.EndOfLine != nil {
			pushBytes(tokens.EndOfLine)
			continue
		}

		if declaration.BooleanAssignment != nil {
			name := []byte{junkByte, junkByte, junkByte, junkByte}
			value := convertBooleanTextToByte(declaration.BooleanAssignment.Boolean.Value)

			pushBytes(tokens.Name)
			pushBytes(name...)
			pushBytes(tokens.Equals)

			// The QB format has no Boolean type.
			// Instead, we must use Ints with a value of 0 or 1.
			//
			// Floats are also legal, but Ints are simpler.
			pushBytes(tokens.Int, value, 0, 0, 0)
			continue
		}

		if declaration.IntegerAssignment != nil {
			name := []byte{junkByte, junkByte, junkByte, junkByte}

			value := *declaration.IntegerAssignment.Integer.Decimal
			valueBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(valueBytes, value)

			pushBytes(tokens.Name)
			pushBytes(name...)
			pushBytes(tokens.Equals)
			pushBytes(tokens.Int)
			pushBytes(valueBytes...)
			continue
		}
	}

	pushBytes(tokens.EndOfFile)

	return bytecode[:numberOfUsedBytes], nil
}

func parseCodeIntoSyntaxTree(code string) (*grammar.SyntaxTree, error) {
	parser := participle.MustBuild(
		&grammar.SyntaxTree{},
		participle.UseLookahead(2),
	)

	syntaxTree := &grammar.SyntaxTree{}

	err := parser.ParseString(code, syntaxTree)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run participle")
	}

	return syntaxTree, nil
}

func convertBooleanTextToByte(text string) byte {
	if text == "true" {
		return 0x01
	} else {
		return 0x00
	}
}
