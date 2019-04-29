package compiler

import (
	goErrors "errors"
	"github.com/alecthomas/participle"
	"github.com/byxor/NeverScript/compiler/grammar"
	"github.com/byxor/NeverScript/shared/checksums"
	"github.com/byxor/NeverScript/shared/tokens"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

const (
	bytecodeSize = 10 * 1000 * 1000
	junkByte     = 0xFF
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
		log.Printf("%+v\n", declaration)

		if declaration.EndOfLine != "" {
			pushBytes(tokens.EndOfLine)
			continue
		}

		if declaration.BooleanAssignment != nil {
			name := []byte{junkByte, junkByte, junkByte, junkByte}
			value := convertBooleanTextToByte(declaration.BooleanAssignment.Value)

			pushBytes(tokens.Name)
			pushBytes(name...)
			pushBytes(tokens.Equals)

			// The QB format has no Boolean type.
			// Instead, we use Ints with a value of 0 or 1.
			pushBytes(tokens.Int, value, 0, 0, 0)
			continue
		}

		if declaration.IntegerAssignment != nil {
			assignment := declaration.IntegerAssignment

			name := []byte{junkByte, junkByte, junkByte, junkByte}

			value, err := convertIntegerNodeToUint32(*assignment.Value)
			if err != nil {
				return []byte{}, errors.Wrap(err, "Failed to convert grammar.Integer to uint32")
			}

			valueBytes := checksums.LittleEndian(value)

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
		participle.Lexer(grammar.NsLexer),
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

func convertIntegerNodeToUint32(node grammar.Integer) (uint32, error) {
	text, base, nodeIsEmpty := (func() (text string, base int, nodeIsEmpty bool) {
		nodeIsEmpty = false

		if node.Base16 != "" {
			base = 16
			text = node.Base16[2:]
			return
		}

		if node.Base10 != "" {
			base = 10
			text = node.Base10
			return
		}

		if node.Base8 != "" {
			base = 8
			text = node.Base8[2:]
			return
		}

		if node.Base2 != "" {
			base = 2
			text = node.Base2[2:]
			return
		}

		nodeIsEmpty = true
		return
	})()

	if nodeIsEmpty {
		return 0, goErrors.New("Integer node is empty")
	}

	value, err := strconv.ParseUint(text, base, 32)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}
