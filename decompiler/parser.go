package decompiler

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"strings"
)

type ParseResult struct {
	WasSuccessful bool
	Reason        string
	BytesRead     int
	Node          compiler.AstNode
}

func ParserFailure(reason string) ParseResult {
	return ParseResult{
		WasSuccessful: false,
		Reason:        reason,
	}
}

func ParserSuccess(bytesRead int, node compiler.AstNode) ParseResult {
	return ParseResult{
		WasSuccessful: true,
		BytesRead:     bytesRead,
		Node:          node,
	}
}

/*
 * The decompiler isn't a priority, but I've written a tiny bit of it just for fun.
 */
func ParseByteCode(arguments *Arguments) error {
	type ParserFunction func(index int) ParseResult
	var ParseRoot ParserFunction
	var ParseRootBodyNode ParserFunction
	var ParseEndOfFile ParserFunction
	var ParseNewLine ParserFunction
	var ParseComma ParserFunction
	var ParseAssignment ParserFunction
	var ParseScript ParserFunction
	var ParseIfStatement ParserFunction
	var ParseExpression ParserFunction
	var ParseInvocation ParserFunction
	var ParseReturn ParserFunction
	var ParseChecksum ParserFunction
	var ParseFloat ParserFunction
	var ParseInteger ParserFunction
	var ParsePair ParserFunction
	var ParseVector ParserFunction
	var ParseString ParserFunction
	var ParseStruct ParserFunction
	var ParseArray ParserFunction
	var ParseAllArgumentsSymbol ParserFunction
	var ParseNameTableEntry ParserFunction
	var ParseBodyOfCode func(index *int) []compiler.AstNode

	bytes := arguments.ByteCode
	numBytes := len(bytes)

	ParseRoot = func(index int) ParseResult {
		start := index
		var bodyNodes []compiler.AstNode
		hasReadAllBytes := false
		for {
			if index >= numBytes {
				hasReadAllBytes = true
				break
			}
			bodyNodeParseResult := ParseRootBodyNode(index)
			if !bodyNodeParseResult.WasSuccessful {
				return ParserFailure(WrapLine(WrapIndex(index, "Failed to parse root"), bodyNodeParseResult.Reason))
			}
			bodyNodes = append(bodyNodes, bodyNodeParseResult.Node)
			index += bodyNodeParseResult.BytesRead
		}
		if !hasReadAllBytes {
			return ParserFailure(WrapIndex(index, "Some bytes left unread"))
		}
		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Root,
			Data: compiler.AstData_Root{
				BodyNodes: bodyNodes,
			},
		})
	}

	ParseRootBodyNode = func(index int) ParseResult {
		parserFunctions := []ParserFunction{
			ParseEndOfFile,
			ParseNewLine,
			ParseAssignment,
			ParseScript,
			ParseNameTableEntry,
		}
		parseResults := make([]ParseResult, len(parserFunctions))
		for i, parserFunction := range parserFunctions {
			parseResult := parserFunction(index)
			parseResults[i] = parseResult
			if parseResult.WasSuccessful {
				return parseResult
			}
		}
		{
			var message strings.Builder
			message.WriteString(WrapIndex(index, "Bytes not recognised as root body node.\n"))
			message.WriteString(hex.Dump(bytes[index : index+32]))
			message.WriteString("\nPOTENTIAL CAUSES.\n")
			message.WriteString("-------------------------------------------------\n")
			for _, parseResult := range parseResults {
				message.WriteString("\n")
				message.WriteString(parseResult.Reason)
				message.WriteString("\n")
			}
			message.WriteString("-------------------------------------------------\n")
			return ParserFailure(message.String())
		}
	}

	ParseEndOfFile = func(index int) ParseResult {
		if bytes[index] != 0 {
			return ParserFailure(WrapIndex(index, fmt.Sprintf("Not an EOF byte '%#X'", bytes[index])))
		}
		return ParserSuccess(1, compiler.AstNode{
			Kind: compiler.AstKind_EndOfFile,
			Data: compiler.AstData_Empty{},
		})
	}

	ParseNewLine = func(index int) ParseResult {
		if bytes[index] != 1 {
			return ParserFailure(WrapIndex(index, fmt.Sprintf("Not a new-line byte '%#X'", bytes[index])))
		}
		return ParserSuccess(1, compiler.AstNode{
			Kind: compiler.AstKind_NewLine,
			Data: compiler.AstData_Empty{},
		})
	}

	ParseComma = func(index int) ParseResult {
		if bytes[index] != 9 {
			return ParserFailure(WrapIndex(index, fmt.Sprintf("Not a comma byte '%#X'", bytes[index])))
		}
		return ParserSuccess(1, compiler.AstNode{
			Kind: compiler.AstKind_Comma,
			Data: compiler.AstData_Empty{},
		})
	}

	ParseAssignment = func(index int) ParseResult {
		start := index

		checksumParseResult := ParseChecksum(index)
		if !checksumParseResult.WasSuccessful {
			return ParserFailure(WrapLine(WrapIndex(index, "Expected checksum at start of assignment"), checksumParseResult.Reason))
		}
		index += checksumParseResult.BytesRead

		if bytes[index] != 7 {
			return ParserFailure(WrapIndex(index, fmt.Sprintf("Expected 0x7 ('=') in assignment, got %#X", bytes[index])))
		}
		index++

		for {
			if parseResult := ParseNewLine(index); parseResult.WasSuccessful {
				index += parseResult.BytesRead
			} else {
				break
			}
		}

		expressionParseResult := ParseExpression(index)
		if !expressionParseResult.WasSuccessful {
			return ParserFailure(WrapLine(WrapIndex(index, "Expected expression at end of assignment"), expressionParseResult.Reason))
		}
		index += expressionParseResult.BytesRead

		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Assignment,
			Data: compiler.AstData_Assignment{
				NameNode:  checksumParseResult.Node,
				ValueNode: expressionParseResult.Node,
			},
		})
	}

	ParseScript = func(index int) ParseResult {
		start := index

		if bytes[index] != 0x23 {
			return ParserFailure(WrapIndex(index, "Script doesn't start with 0x23"))
		}
		index++

		checksumParseResult := ParseChecksum(index)
		if !checksumParseResult.WasSuccessful {
			return ParserFailure(WrapLine(WrapIndex(index, "Expected checksum for name of script"), checksumParseResult.Reason))
		}
		index += checksumParseResult.BytesRead

		bodyNodes := ParseBodyOfCode(&index)

		if bytes[index] != 0x24 {
			return ParserFailure(WrapLine(WrapIndex(index, "Script doesn't end with 0x24"), hex.Dump(bytes[index:index+64])))
		}
		index++

		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Script,
			Data: compiler.AstData_Script{
				NameNode:              checksumParseResult.Node,
				DefaultParameterNodes: []compiler.AstNode{},
				BodyNodes:             bodyNodes,
			},
		})
	}

	ParseIfStatement = func(index int) ParseResult {
		start := index

		if bytes[index] != 0x47 {
			return ParserFailure(WrapIndex(index, "If statement doesn't start with 0x47"))
		}
		index++
		if index+2 >= numBytes {
			return ParserFailure(WrapIndex(index+2, "Reached EOF when scanning the jump offset"))
		}
		//offsetBytes := bytes[index:index+2]
		//offsetSize := binary.LittleEndian.Uint16(offsetBytes)
		index += 2

		conditionParseResult := ParseExpression(index)
		if !conditionParseResult.WasSuccessful {
			return ParserFailure(WrapIndex(index, WrapLine("Failed to parse if statement condition", conditionParseResult.Reason)))
		}
		index += conditionParseResult.BytesRead

		bodyNodes := ParseBodyOfCode(&index)
		bodies := [][]compiler.AstNode{bodyNodes}

		if bytes[index] == 0x48 { // has 'else'
			index++
			index += 2
			bodyNodes := ParseBodyOfCode(&index)
			bodies = append(bodies, bodyNodes)
		}

		index++
		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_IfStatement,
			Data: compiler.AstData_IfStatement{
				BooleanInvocationData: make([]bool, 1), // FIXME(brandon): Arbitrary data for now, will crash if number of branches exceeds the capacity
				Conditions:            []compiler.AstNode{conditionParseResult.Node},
				Bodies:                bodies,
			},
		})
	}

	ParseExpression = func(index int) ParseResult {
		ParseExpressionInner := func(index int) ParseResult {
			parserFunctions := []ParserFunction{
				ParseAllArgumentsSymbol,
				ParseFloat,
				ParseInteger,
				ParsePair,
				ParseVector,
				ParseStruct,
				ParseArray,
				ParseInvocation,
				ParseChecksum,
				ParseString,
			}
			for _, parserFunction := range parserFunctions {
				parseResult := parserFunction(index)
				if parseResult.WasSuccessful {
					return parseResult
				}
			}
			return ParserFailure(WrapLine(WrapIndex(index, "Bytes not recognised as an expression"), hex.Dump(bytes[index:index+64])))
		}
		if bytes[index] == 0x39 {
			index++
			innerExpressionParseResult := ParseExpressionInner(index)
			if !innerExpressionParseResult.WasSuccessful {
				return ParserFailure(WrapIndex(index, "No expression after '!' (0x39)"))
			}
			return ParserSuccess(1+innerExpressionParseResult.BytesRead, compiler.AstNode{
				Kind: compiler.AstKind_LogicalNot,
				Data: compiler.AstData_UnaryExpression{
					Node: innerExpressionParseResult.Node,
				},
			})
		} else {
			return ParseExpressionInner(index)
		}
	}

	ParseInvocation = func(index int) ParseResult {
		start := index
		checksumParseResult := ParseChecksum(index)
		if !checksumParseResult.WasSuccessful {
			return ParserFailure(WrapLine(WrapIndex(index, "Invocation does not start with checksum"), checksumParseResult.Reason))
		}
		index += checksumParseResult.BytesRead
		var parameterNodes []compiler.AstNode
		parserFunctions := []ParserFunction{
			ParseAssignment,
			ParseExpression,
		}
		for {
			foundParameter := false
			for _, parserFunction := range parserFunctions {
				parseResult := parserFunction(index)
				if parseResult.WasSuccessful {
					foundParameter = true
					parameterNodes = append(parameterNodes, parseResult.Node)
					index += parseResult.BytesRead
					break
				}
			}
			if !foundParameter {
				break
			}
		}
		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Invocation,
			Data: compiler.AstData_Invocation{
				ScriptIdentifierNode: checksumParseResult.Node,
				ParameterNodes:       parameterNodes,
			},
		})
	}

	ParseReturn = func(index int) ParseResult {
		start := index
		if bytes[index] != 0x29 {
			return ParserFailure(WrapIndex(index, "Not a 'return' statement (expected 0x29)"))
		}
		index++
		var parameterNodes []compiler.AstNode
		parserFunctions := []ParserFunction{
			ParseAssignment,
			ParseExpression,
		}
		for {
			foundParameter := false
			for _, parserFunction := range parserFunctions {
				parseResult := parserFunction(index)
				if parseResult.WasSuccessful {
					foundParameter = true
					parameterNodes = append(parameterNodes, parseResult.Node)
					index += parseResult.BytesRead
					break
				}
			}
			if !foundParameter {
				break
			}
		}
		// FIXME(brandon): I'm returning an "invocation" to save time because it uses the same syntax as 'return'.
		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Invocation,
			Data: compiler.AstData_Invocation{
				ScriptIdentifierNode: compiler.AstNode{
					Kind: compiler.AstKind_Checksum,
					Data: compiler.AstData_Checksum{
						ChecksumToken: compiler.Token{
							Kind: compiler.TokenKind_Identifier,
							Data: "return",
						},
					},
				},
				ParameterNodes:       parameterNodes,
			},
		})
	}

	ParseChecksum = func(index int) ParseResult {
		start := index
		isLocalReference := false
		if bytes[index] == 0x2D {
			isLocalReference = true
			index++
		}
		if bytes[index] != 0x16 {
			return ParserFailure(WrapLine(WrapIndex(index, "Checksum doesn't have 0x16"), hex.Dump(bytes[index:index+32])))
		}
		index++
		if index+4 >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when scanning checksum bytes"))
		}
		checksumBytes := bytes[index : index+4]
		index += 4
		node := compiler.AstNode{
			Kind: compiler.AstKind_Checksum,
			Data: compiler.AstData_Checksum{
				ChecksumBytes: checksumBytes,
			},
		}
		if isLocalReference {
			node = compiler.AstNode{
				Kind: compiler.AstKind_LocalReference,
				Data: compiler.AstData_UnaryExpression{
					Node: node,
				},
			}
		}
		return ParserSuccess(index-start, node)
	}

	ParseFloat = func(index int) ParseResult {
		if bytes[index] != 0x1A {
			return ParserFailure(WrapLine(WrapIndex(index, "Float doesn't start with 0x1A"), hex.Dump(bytes[index:index+32])))
		}
		index++
		if index+4 >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when scanning the float's bytes"))
		}
		return ParserSuccess(5, compiler.AstNode{
			Kind: compiler.AstKind_Float,
			Data: compiler.AstData_Float{
				FloatBytes: bytes[index : index+4],
			},
		})
	}

	ParseInteger = func(index int) ParseResult {
		if bytes[index] != 0x17 {
			return ParserFailure(WrapLine(WrapIndex(index, "Integer doesn't start with 0x17"), hex.Dump(bytes[index:index+32])))
		}
		index++
		if index+4 >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when scanning the integer's bytes"))
		}
		return ParserSuccess(5, compiler.AstNode{
			Kind: compiler.AstKind_Integer,
			Data: compiler.AstData_Integer{
				IntegerBytes: bytes[index : index+4],
			},
		})
	}

	ParsePair = func(index int) ParseResult {
		if bytes[index] != 0x1F {
			return ParserFailure(WrapLine(WrapIndex(index, "Pair doesn't start with 0x1F"), hex.Dump(bytes[index:index+32])))
		}
		index++
		if index+8 >= numBytes {
			return ParserFailure(WrapIndex(index+8, "Reached EOF when scanning the pair's float values"))
		}
		return ParserSuccess(9, compiler.AstNode{
			Kind: compiler.AstKind_Pair,
			Data: compiler.AstData_Pair{
				FloatNodeA: compiler.AstNode{
					Kind: compiler.AstKind_Float,
					Data: compiler.AstData_Float{
						FloatBytes: bytes[index : index+4],
					},
				},
				FloatNodeB: compiler.AstNode{
					Kind: compiler.AstKind_Float,
					Data: compiler.AstData_Float{
						FloatBytes: bytes[index+4 : index+8],
					},
				},
			},
		})
	}

	ParseVector = func(index int) ParseResult {
		if bytes[index] != 0x1E {
			return ParserFailure(WrapLine(WrapIndex(index, "Vector doesn't start with 0x1E"), hex.Dump(bytes[index:index+32])))
		}
		index++
		if index+12 >= numBytes {
			return ParserFailure(WrapIndex(index+8, "Reached EOF when scanning the vector's float values"))
		}
		return ParserSuccess(13, compiler.AstNode{
			Kind: compiler.AstKind_Vector,
			Data: compiler.AstData_Vector{
				FloatNodeA: compiler.AstNode{
					Kind: compiler.AstKind_Float,
					Data: compiler.AstData_Float{
						FloatBytes: bytes[index : index+4],
					},
				},
				FloatNodeB: compiler.AstNode{
					Kind: compiler.AstKind_Float,
					Data: compiler.AstData_Float{
						FloatBytes: bytes[index+4 : index+8],
					},
				},
				FloatNodeC: compiler.AstNode{
					Kind: compiler.AstKind_Float,
					Data: compiler.AstData_Float{
						FloatBytes: bytes[index+8 : index+12],
					},
				},
			},
		})
	}

	ParseString = func(index int) ParseResult {
		if bytes[index] != 0x1B {
			return ParserFailure(WrapIndex(index, "String doesn't start with 0x1B"))
		}
		index++
		if index+4 >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when scanning string size"))
		}
		stringSize := binary.LittleEndian.Uint32(bytes[index : index+4])
		index += 4
		if index+int(stringSize) >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when scanning string contents"))
		}
		return ParserSuccess(5+int(stringSize), compiler.AstNode{
			Kind: compiler.AstKind_String,
			Data: compiler.AstData_String{
				StringBytes: bytes[index : index+int(stringSize)],
			},
		})
	}

	ParseStruct = func(index int) ParseResult {
		start := index
		if bytes[index] != 3 {
			return ParserFailure(WrapLine(WrapIndex(index, "Struct doesn't start with 0x3"), hex.Dump(bytes[index:index+32])))
		}
		index++
		var structElementNodes []compiler.AstNode
		structElementParserFunctions := []ParserFunction{
			ParseNewLine,
			ParseComma,
			ParseAssignment,
			ParseExpression,
		}
		for {
			if bytes[index] == 4 {
				index++
				break
			}
			foundElement := false
			for _, parserFunction := range structElementParserFunctions {
				parseResult := parserFunction(index)
				if parseResult.WasSuccessful {
					foundElement = true
					structElementNodes = append(structElementNodes, parseResult.Node)
					index += parseResult.BytesRead
					break
				}
			}
			if !foundElement {
				return ParserFailure(WrapLine(WrapIndex(index, "Bytes not recognised as a struct element"), hex.Dump(bytes[index:index+32])))
			}
		}

		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Struct,
			Data: compiler.AstData_Struct{
				ElementNodes: structElementNodes,
			},
		})
	}

	ParseArray = func(index int) ParseResult {
		start := index
		if bytes[index] != 5 {
			return ParserFailure(WrapLine(WrapIndex(index, "Array doesn't start with 0x5"), hex.Dump(bytes[index:index+32])))
		}
		index++
		var elements []compiler.AstNode
		elementParserFunctions := []ParserFunction{
			ParseNewLine,
			ParseComma,
			ParseExpression,
		}
		for {
			if bytes[index] == 6 {
				index++
				break
			}
			foundElement := false
			for _, parserFunction := range elementParserFunctions {
				parseResult := parserFunction(index)
				if parseResult.WasSuccessful {
					foundElement = true
					elements = append(elements, parseResult.Node)
					index += parseResult.BytesRead
					break
				}
			}
			if !foundElement {
				return ParserFailure(WrapLine(WrapIndex(index, "Bytes not recognised as an array element"), hex.Dump(bytes[index:index+32])))
			}
		}

		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_Array,
			Data: compiler.AstData_Array{
				ElementNodes: elements,
			},
		})
	}

	ParseAllArgumentsSymbol = func(index int) ParseResult {
		if bytes[index] != 0x2C {
			return ParserFailure("Not an AllArgumentsSymbol (<...>) (0x2C)")
		}
		return ParserSuccess(1, compiler.AstNode{
			Kind: compiler.AstKind_AllArguments,
			Data: compiler.AstData_Empty{},
		})
	}

	ParseNameTableEntry = func(index int) ParseResult {
		start := index

		if bytes[index] != 0x2B {
			return ParserFailure(WrapLine(WrapIndex(index, "Name table entry doesn't start with 0x2B"), hex.Dump(bytes[index:index+32])))
		}
		index++

		if index+4 >= numBytes {
			return ParserFailure(WrapIndex(index+4, "Reached EOF when reading checksum for name table entry"))
		}
		checksumBytes := bytes[index : index+4]
		index += 4

		nameStart := index
		var name string
		for {
			if bytes[index] == 0 {
				index++
				name = string(bytes[nameStart : index-1])
				break
			}
			index++
		}

		return ParserSuccess(index-start, compiler.AstNode{
			Kind: compiler.AstKind_NameTableEntry,
			Data: compiler.AstData_NameTableEntry{
				ChecksumBytes: checksumBytes,
				Name:          name,
			},
		})
	}

	ParseBodyOfCode = func(index *int) []compiler.AstNode {
		var bodyNodes []compiler.AstNode
		parserFunctions := []ParserFunction{
			ParseNewLine,
			ParseAssignment,
			ParseIfStatement,
			ParseExpression,
			ParseReturn,
		}
		for {
			foundSomething := false
			for _, parserFunction := range parserFunctions {
				parseResult := parserFunction(*index)
				if parseResult.WasSuccessful {
					foundSomething = true
					bodyNodes = append(bodyNodes, parseResult.Node)
					*index += parseResult.BytesRead
					break
				}
			}
			if !foundSomething {
				break
			}
		}
		return bodyNodes
	}

	parseResult := ParseRoot(0)
	if !parseResult.WasSuccessful {
		return errors.New(parseResult.Reason)
	}

	arguments.RootNode = parseResult.Node
	return nil
}

func WrapLine(outer, inner string) string {
	return fmt.Sprintf("%s:\n%s", outer, inner)
}

func WrapIndex(index int, message string) string {
	return fmt.Sprintf("[%#X] %s", index, message)
}
