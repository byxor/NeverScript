package compiler

import (
	"fmt"
	"strings"
)

type ParseResult struct {
	WasSuccessful  bool
	Reason         string
	Node           AstNode
	TokensConsumed int
}

type Parser struct {
	Tokens []Token
	Result ParseResult
}

func BuildAbstractSyntaxTree(parser *Parser) {
	var ParseRoot func( /*index is always 0*/) ParseResult
	var ParseRootBodyNode func(index int) ParseResult
	var ParseExpression func(index int, allowInvocations bool) ParseResult
	var ParseExpressionBeginningWithLeftParenthesis func(index int) ParseResult
	var ParseAssignment func(index int, allowInvocations bool) ParseResult
	var ParseScript func(index int) ParseResult
	var ParseBodyOfCode func(index int) (ParseResult, []AstNode)
	var ParseWhileLoop func(index int) ParseResult
	var ParseLogicalNot func(index int) ParseResult
	var ParseIfStatement func(index int) ParseResult
	var ParseChecksumOrInvocation func(index int, allowInvocations bool) ParseResult
	var ParseInvocation func(index int) ParseResult
	var ParseInvocationParameter func(index int) ParseResult
	var ParseRandom func(index int) ParseResult
	var ParseComment func(index int) ParseResult
	var ParseChecksum func(index int) ParseResult
	var ParseFloat func(index int) ParseResult
	var ParseInteger func(index int) ParseResult
	var ParseString func(index int) ParseResult
	var ParseArray func(index int) ParseResult
	var ParseStruct func(index int) ParseResult
	var ParseNewLine func(index int) ParseResult
	var ParseBreak func(index int) ParseResult
	var ParseReturn func(index int) ParseResult
	var ParseComma func(index int) ParseResult
	// TODO(brandon): var SkipOverCommentsAndEscapedNewlines func(index int) int
	var GetKind func(index int) TokenKind
	var GetToken func(index int) Token

	ParseRoot = func() ParseResult {
		var bodyNodes AstNodeBuffer

		// Ensure root starts with new-line
		bodyNodes.MaybeSave(ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_NewLine,
				Data: AstData_Empty{},
			},
			TokensConsumed: 1,
		})

		// Parse root body nodes until you can't anymore
		index := 0
		var bodyNodeParseResult ParseResult
		for {
			bodyNodeParseResult = ParseRootBodyNode(index)
			if bodyNodeParseResult.WasSuccessful {
				bodyNodes.MaybeSave(bodyNodeParseResult)
				index += bodyNodeParseResult.TokensConsumed
			} else {
				break
			}
		}

		// Display
		if numOfTokens := len(parser.Tokens); bodyNodes.TokensConsumed < numOfTokens {
			var messageBuilder strings.Builder
			messageBuilder.WriteString("\n\nFinished parsing but didn't read all tokens.\n")
			messageBuilder.WriteString(fmt.Sprintf("Read %d/%d (%d left unread).\n", bodyNodes.TokensConsumed, numOfTokens, numOfTokens-bodyNodes.TokensConsumed))
			for _, unreadToken := range parser.Tokens[bodyNodes.TokensConsumed:numOfTokens] {
				messageBuilder.WriteString(fmt.Sprintf("  %+v,\n", unreadToken))
			}
			messageBuilder.WriteString(fmt.Sprintf("\nPotential cause: %s\n", bodyNodeParseResult.Reason))
			return ParseResult{
				WasSuccessful: false,
				Reason:        messageBuilder.String(),
			}
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Root,
				Data: AstData_Root{
					BodyNodes: bodyNodes.Nodes,
				},
			},
			TokensConsumed: bodyNodes.TokensConsumed,
		}
	}

	ParseRootBodyNode = func(index int) ParseResult {
		if parseResult := ParseNewLine(index); parseResult.WasSuccessful {
			return parseResult
		}

		if parseResult := ParseComment(index); parseResult.WasSuccessful {
			return parseResult
		}

		if parseResult := ParseScript(index); parseResult.WasSuccessful {
			return parseResult
		}

		if parseResult := ParseAssignment(index, true); parseResult.WasSuccessful {
			return parseResult
		}

		if parseResult := ParseExpression(index, true); parseResult.WasSuccessful {
			return parseResult
		}

		return ParseResult{
			WasSuccessful: false,
			Reason:        TokensNotRecognisedError(parser.Tokens[index:], "a root body node"),
		}
	}

	ParseExpression = func(index int, allowInvocations bool) ParseResult {
		parseExpressionInner := func(index int, allowInvocations bool) ParseResult {
			if GetKind(index) == TokenKind_If {
				return ParseResult{
					WasSuccessful: false,
					Reason:        "Failed to parse expression, found 'if'",
				}
			}
			if GetKind(index) == TokenKind_LeftAngleBracket &&
				GetKind(index+1) == TokenKind_Dot &&
				GetKind(index+2) == TokenKind_Dot &&
				GetKind(index+3) == TokenKind_Dot &&
				GetKind(index+4) == TokenKind_RightAngleBracket {
				return ParseResult{
					WasSuccessful: true,
					Node: AstNode{
						Kind: AstKind_AllArguments,
						Data: AstData_Empty{},
					},
					TokensConsumed: 5,
				}
			}
			if GetKind(index) == TokenKind_LeftAngleBracket {
				futureIndex := index + 1
				// parseResult := ParseChecksumOrInvocation(futureIndex, allowInvocations)
				// BUG: ^ invocation comes after the local reference, not inside
				parseResult := ParseChecksum(futureIndex)
				if !parseResult.WasSuccessful {
					return ParseResult{
						WasSuccessful: false,
						Reason:        WrapStr("Failed to parse expression for local reference", parseResult.Reason),
					}
				}
				futureIndex++
				if GetKind(futureIndex) != TokenKind_RightAngleBracket {
					return ParseResult{
						WasSuccessful: false,
						Reason:        WrapStr("Failed to parse local reference, no '>'", parseResult.Reason),
					}
				}
				return ParseResult{
					WasSuccessful: true,
					Node: AstNode{
						Kind: AstKind_LocalReference,
						Data: AstData_LocalReference{
							Node: parseResult.Node,
						},
					},
					TokensConsumed: 2 + parseResult.TokensConsumed,
				}
			}
			if GetKind(index) == TokenKind_Minus {
				nextToken := GetToken(index + 1)
				if nextToken.Kind == TokenKind_Integer {
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_Integer,
							Data: AstData_Integer{
								IntegerToken: Token{
									Kind:       TokenKind_Integer,
									Data:       "-" + nextToken.Data,
									LineNumber: nextToken.LineNumber,
								},
							},
						},
						TokensConsumed: 2,
					}
				} else if nextToken.Kind == TokenKind_Float {
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_Float,
							Data: AstData_Float{
								FloatToken: Token{
									Kind:       TokenKind_Float,
									Data:       "-" + nextToken.Data,
									LineNumber: nextToken.LineNumber,
								},
							},
						},
						TokensConsumed: 2,
					}
				}
			}

			if parseResult := ParseRandom(index); parseResult.WasSuccessful {
				return parseResult
			}
			if GetKind(index) == TokenKind_Identifier || GetKind(index) == TokenKind_RawChecksum {
				return ParseChecksumOrInvocation(index, allowInvocations)
			}
			if GetKind(index) == TokenKind_Bang {
				return ParseLogicalNot(index)
			}
			if GetKind(index) == TokenKind_Integer {
				return ParseInteger(index)
			}
			if GetKind(index) == TokenKind_Float {
				return ParseFloat(index)
			}
			if GetKind(index) == TokenKind_String {
				return ParseString(index)
			}
			if GetKind(index) == TokenKind_LeftParenthesis {
				return ParseExpressionBeginningWithLeftParenthesis(index)
			}
			if GetKind(index) == TokenKind_LeftSquareBracket {
				return ParseArray(index)
			}
			if GetKind(index) == TokenKind_LeftCurlyBrace {
				return ParseStruct(index)
			}
			GetToken(index)
			return ParseResult{
				WasSuccessful: false,
				Reason:        TokensNotRecognisedError(parser.Tokens[index:], "an expression"),
			}
		}

		inPlaceMathOperationParseResult := func(index int, leftParseResult ParseResult, rightParseResult ParseResult, expressionKind AstKind) ParseResult {
			return ParseResult{
				WasSuccessful: true,
				Node: AstNode{
					Kind: AstKind_Assignment,
					Data: AstData_Assignment{
						NameNode: leftParseResult.Node,
						ValueNode: AstNode{
							Kind: expressionKind,
							Data: AstData_BinaryExpression{
								LeftNode:  leftParseResult.Node,
								RightNode: rightParseResult.Node,
							},
						},
					},
				},
				TokensConsumed: leftParseResult.TokensConsumed + 2 + rightParseResult.TokensConsumed,
			}
		}

		expressionParseResult := parseExpressionInner(index, allowInvocations)
		if expressionParseResult.WasSuccessful {
			index += expressionParseResult.TokensConsumed
			if GetKind(index) == TokenKind_Dot {
				index++
				secondExpressionParseResult := ParseExpression(index, true)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_DotExpression,
							Data: AstData_BinaryExpression{
								LeftNode:  expressionParseResult.Node,
								RightNode: secondExpressionParseResult.Node,
							},
						},
						TokensConsumed: expressionParseResult.TokensConsumed + 1 + secondExpressionParseResult.TokensConsumed,
					}
				}
			} else if GetKind(index) == TokenKind_Colon {
				index++
				secondExpressionParseResult := ParseExpression(index, true)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_ColonExpression,
							Data: AstData_BinaryExpression{
								LeftNode:  expressionParseResult.Node,
								RightNode: secondExpressionParseResult.Node,
							},
						},
						TokensConsumed: expressionParseResult.TokensConsumed + 1 + secondExpressionParseResult.TokensConsumed,
					}
				}
			} else if GetKind(index) == TokenKind_And {
				index++
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_LogicalAnd,
							Data: AstData_BinaryExpression{
								LeftNode:  expressionParseResult.Node,
								RightNode: secondExpressionParseResult.Node,
							},
						},
						TokensConsumed: expressionParseResult.TokensConsumed + 1 + secondExpressionParseResult.TokensConsumed,
					}
				}
			} else if GetKind(index) == TokenKind_Or {
				index++
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return ParseResult{
						WasSuccessful: true,
						Node: AstNode{
							Kind: AstKind_LogicalOr,
							Data: AstData_BinaryExpression{
								LeftNode:  expressionParseResult.Node,
								RightNode: secondExpressionParseResult.Node,
							},
						},
						TokensConsumed: expressionParseResult.TokensConsumed + 1 + secondExpressionParseResult.TokensConsumed,
					}
				}
			} else if GetKind(index) == TokenKind_Plus && GetKind(index+1) == TokenKind_Equals {
				index += 2
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return inPlaceMathOperationParseResult(index, expressionParseResult, secondExpressionParseResult, AstKind_AdditionExpression)
				}
			} else if GetKind(index) == TokenKind_Minus && GetKind(index+1) == TokenKind_Equals {
				index += 2
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return inPlaceMathOperationParseResult(index, expressionParseResult, secondExpressionParseResult, AstKind_SubtractionExpression)
				}
			} else if GetKind(index) == TokenKind_Asterisk && GetKind(index+1) == TokenKind_Equals {
				index += 2
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return inPlaceMathOperationParseResult(index, expressionParseResult, secondExpressionParseResult, AstKind_MultiplicationExpression)
				}
			} else if GetKind(index) == TokenKind_ForwardSlash && GetKind(index+1) == TokenKind_Equals {
				index += 2
				secondExpressionParseResult := ParseExpression(index, false)
				if secondExpressionParseResult.WasSuccessful {
					index += secondExpressionParseResult.TokensConsumed
					return inPlaceMathOperationParseResult(index, expressionParseResult, secondExpressionParseResult, AstKind_DivisionExpression)
				}
			}
		}

		return expressionParseResult
	}

	ParseExpressionBeginningWithLeftParenthesis = func(index int) ParseResult {
		oldIndex := index
		index++

		firstParseResult := ParseExpression(index, true)
		if firstParseResult.WasSuccessful {
			index += firstParseResult.TokensConsumed
			if GetKind(index) == TokenKind_RightParenthesis {
				return ParseResult{
					WasSuccessful: true,
					Node: AstNode{
						Kind: AstKind_UnaryExpression,
						Data: AstData_UnaryExpression{
							Node: firstParseResult.Node,
						},
					},
					TokensConsumed: 2 + firstParseResult.TokensConsumed,
				}
			}

			if GetKind(index) == TokenKind_Comma {
				index++
				if GetKind(index) == TokenKind_Float {
					secondParseResult := ParseExpression(index, true)
					if secondParseResult.WasSuccessful {
						index += secondParseResult.TokensConsumed
						if GetKind(index) == TokenKind_RightParenthesis {
							return ParseResult{
								WasSuccessful: true,
								Node: AstNode{
									Kind: AstKind_Pair,
									Data: AstData_Pair{
										FloatNodeA: firstParseResult.Node,
										FloatNodeB: secondParseResult.Node,
									},
								},
								TokensConsumed: 3 + firstParseResult.TokensConsumed + secondParseResult.TokensConsumed,
							}
						}
						if GetKind(index) == TokenKind_Comma {
							index++
							if GetKind(index) == TokenKind_Float {
								thirdParseResult := ParseExpression(index, true)
								if thirdParseResult.WasSuccessful {
									index += thirdParseResult.TokensConsumed
									if GetKind(index) == TokenKind_RightParenthesis {
										return ParseResult{
											WasSuccessful: true,
											Node: AstNode{
												Kind: AstKind_Vector,
												Data: AstData_Vector{
													FloatNodeA: firstParseResult.Node,
													FloatNodeB: secondParseResult.Node,
													FloatNodeC: thirdParseResult.Node,
												},
											},
											TokensConsumed: 4 + firstParseResult.TokensConsumed + secondParseResult.TokensConsumed + thirdParseResult.TokensConsumed,
										}
									}
								}
							}
						}
					}
				}
			}

			handleBinaryOperator := func(astKind AstKind, size int) ParseResult {
				index += size
				secondInnerExpressionParseResult := ParseExpression(index, true)
				if secondInnerExpressionParseResult.WasSuccessful {
					index += secondInnerExpressionParseResult.TokensConsumed
					if GetKind(index) == TokenKind_RightParenthesis {
						return ParseResult{
							WasSuccessful: true,
							Node: AstNode{
								Kind: astKind,
								Data: AstData_BinaryExpression{
									LeftNode:  firstParseResult.Node,
									RightNode: secondInnerExpressionParseResult.Node,
								},
							},
							TokensConsumed: 2 + firstParseResult.TokensConsumed + size + secondInnerExpressionParseResult.TokensConsumed,
						}
					}
				}
				return ParseResult{
					WasSuccessful: false,
					Reason:        "Couldn't parse binary operator expression",
				}
			}
			if GetKind(index) == TokenKind_Plus {
				return handleBinaryOperator(AstKind_AdditionExpression, 1)
			}
			if GetKind(index) == TokenKind_Minus {
				return handleBinaryOperator(AstKind_SubtractionExpression, 1)
			}
			if GetKind(index) == TokenKind_Asterisk {
				return handleBinaryOperator(AstKind_MultiplicationExpression, 1)
			}
			if GetKind(index) == TokenKind_ForwardSlash {
				return handleBinaryOperator(AstKind_DivisionExpression, 1)
			}
			if GetKind(index) == TokenKind_Bang && GetKind(index+1) == TokenKind_Equals {
				return handleBinaryOperator(AstKind_NotEqualExpression, 2)
			}
			if GetKind(index) == TokenKind_LeftAngleBracket && GetKind(index+1) == TokenKind_Equals {
				return handleBinaryOperator(AstKind_LessThanEqualsExpression, 2)
			}
			if GetKind(index) == TokenKind_RightAngleBracket && GetKind(index+1) == TokenKind_Equals {
				return handleBinaryOperator(AstKind_GreaterThanEqualsExpression, 2)
			}
			if GetKind(index) == TokenKind_RightAngleBracket {
				return handleBinaryOperator(AstKind_GreaterThanExpression, 1)
			}
			if GetKind(index) == TokenKind_LeftAngleBracket {
				return handleBinaryOperator(AstKind_LessThanExpression, 1)
			}
			if GetKind(index) == TokenKind_Equals {
				return handleBinaryOperator(AstKind_EqualsExpression, 1)
			}
		}

		return ParseResult{
			WasSuccessful: false,
			Reason:        TokensNotRecognisedError(parser.Tokens[oldIndex:], "an expression beginning with a left parenthesis"),
		}
	}

	ParseChecksumOrInvocation = func(index int, allowInvocations bool) ParseResult {
		if allowInvocations {
			extraTokens := 0
			for {
				if GetKind(index+extraTokens+1) == TokenKind_BackwardSlash && GetKind(index+extraTokens+2) == TokenKind_NewLine {
					extraTokens += 2
				} else {
					break
				}
			}
			if expressionParseResult := ParseExpression(index+extraTokens+1, allowInvocations); expressionParseResult.WasSuccessful {
				invocation := ParseInvocation(index)
				return invocation
			}
		}
		return ParseChecksum(index)
		// TODO(brandon): Check for array access... e.g. identifier[int]
		// edit ^ might not be necessary since it's just an array beside an identifier in the bytecode
	}

	ParseChecksum = func(index int) ParseResult {
		if GetKind(index) == TokenKind_LeftAngleBracket &&
			GetKind(index+2) == TokenKind_RightAngleBracket {
			return ParseResult{
				WasSuccessful: true,
				Node: AstNode{
					Kind: AstKind_LocalReference,
					Data: AstData_LocalReference{
						Node: AstNode{
							Kind: AstKind_Checksum,
							Data: AstData_Checksum{
								ChecksumToken: GetToken(index + 1),
								IsRawChecksum: GetKind(index+1) == TokenKind_RawChecksum,
							},
						},
					},
				},
				TokensConsumed: 3,
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Checksum,
				Data: AstData_Checksum{
					ChecksumToken: GetToken(index),
					IsRawChecksum: GetKind(index) == TokenKind_RawChecksum,
				},
			},
			TokensConsumed: 1,
		}
	}

	ParseFloat = func(index int) ParseResult {
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Float,
				Data: AstData_Float{
					FloatToken: GetToken(index),
				},
			},
			TokensConsumed: 1,
		}
	}

	ParseInteger = func(index int) ParseResult {
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Integer,
				Data: AstData_Integer{
					IntegerToken: GetToken(index),
				},
			},
			TokensConsumed: 1,
		}
	}

	ParseString = func(index int) ParseResult {
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_String,
				Data: AstData_String{
					StringToken: GetToken(index),
				},
			},
			TokensConsumed: 1,
		}
	}

	ParseArray = func(index int) ParseResult {
		index++

		// gather array elements
		var elementNodes AstNodeBuffer
		for {
			if GetKind(index) == TokenKind_BackwardSlash && GetKind(index+1) == TokenKind_NewLine {
				index += 2
				elementNodes.TokensConsumed += 2
			}
			if GetKind(index) == TokenKind_RightSquareBracket {
				break
			}
			if GetKind(index) == TokenKind_SingleLineComment || GetKind(index) == TokenKind_MultiLineComment {
				index++
				elementNodes.TokensConsumed++
			}
			if newLineParseResult := ParseNewLine(index); newLineParseResult.WasSuccessful {
				elementNodes.MaybeSave(newLineParseResult)
				index += newLineParseResult.TokensConsumed
			}
			if commaParseResult := ParseComma(index); commaParseResult.WasSuccessful {
				elementNodes.MaybeSave(commaParseResult)
				index += commaParseResult.TokensConsumed
			}
			if expressionParseResult := ParseExpression(index, true); expressionParseResult.WasSuccessful {
				elementNodes.MaybeSave(expressionParseResult)
				index += expressionParseResult.TokensConsumed
			}
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Array,
				Data: AstData_Array{
					ElementNodes: elementNodes.Nodes,
				},
			},
			TokensConsumed: 2 + elementNodes.TokensConsumed,
		}
	}

	ParseStruct = func(index int) ParseResult {
		var elementNodes AstNodeBuffer

		if GetKind(index) != TokenKind_LeftCurlyBrace {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Couldn't parse struct, no '{' found",
			}
		}

		index++
		indexAfterLastIteration := index
		for { // gather struct elements
			if GetKind(index) == TokenKind_BackwardSlash && GetKind(index+1) == TokenKind_NewLine {
				elementNodes.TokensConsumed += 2
				index += 2
			} else if parseResult := ParseNewLine(index); parseResult.WasSuccessful {
				elementNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if GetKind(index) == TokenKind_If {
				return ParseResult{
					WasSuccessful: false,
					Reason:        "Failed to parse struct elements, found 'if'",
				}
			} else if GetKind(index) == TokenKind_While {
				return ParseResult{
					WasSuccessful: false,
					Reason:        "Failed to parse struct elements, found 'while'",
				}
			} else if GetKind(index) == TokenKind_RightCurlyBrace {
				break
			} else if parseResult := ParseComma(index); parseResult.WasSuccessful {
				elementNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseComment(index); parseResult.WasSuccessful {
				elementNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseAssignment(index, true); parseResult.WasSuccessful {
				elementNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseExpression(index, true); parseResult.WasSuccessful {
				elementNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if index == indexAfterLastIteration {
				return ParseResult{
					WasSuccessful: false,
					Reason:        TokensNotRecognisedError(parser.Tokens[index:], "a struct element"),
				}
			}
			indexAfterLastIteration = index
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Struct,
				Data: AstData_Struct{
					ElementNodes: elementNodes.Nodes,
				},
			},
			TokensConsumed: 2 + elementNodes.TokensConsumed,
		}
	}

	ParseAssignment = func(index int, allowInvocations bool) ParseResult {
		start := index
		nameParseResult := ParseChecksum(index)
		if !nameParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in assignment wasn't a checksum",
			}
		}
		index += nameParseResult.TokensConsumed

		if GetKind(index) != TokenKind_Equals {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Second token in assignment wasn't an 'equals'",
			}
		}
		index++

		valueParseResult := ParseExpression(index, allowInvocations)
		if !valueParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        WrapStr("Couldn't parse expression for value of assignment", valueParseResult.Reason),
			}
		}
		index += valueParseResult.TokensConsumed

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Assignment,
				Data: AstData_Assignment{
					NameNode:  nameParseResult.Node,
					ValueNode: valueParseResult.Node,
				},
			},
			TokensConsumed: index - start,
		}
	}

	pruneStructIfInvoked := func(parseResult *ParseResult, index *int) {
		// TODO(brandon): semantically compress this code, it does basically the same thing twice but for 2 cases
		if parseResult.Node.Kind == AstKind_Invocation {
			invocationData := parseResult.Node.Data.(AstData_Invocation)
			lastParameterIndex := len(invocationData.ParameterNodes) - 1
			lastParameterNode := invocationData.ParameterNodes[lastParameterIndex]
			if lastParameterNode.Kind == AstKind_Struct {
				// remove struct from params
				parseResult.Node.Data = AstData_Invocation{
					ScriptIdentifierNode:              invocationData.ScriptIdentifierNode,
					ParameterNodes:                    invocationData.ParameterNodes[:lastParameterIndex],
					TokensConsumedByEachParameterNode: invocationData.TokensConsumedByEachParameterNode[:lastParameterIndex],
				}

				// skip backwards over the struct so it will be read as the if-statement body
				*index -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]

				// reduce the number of tokens consumed by the condition
				parseResult.TokensConsumed -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]
			}
		} else if parseResult.Node.Kind == AstKind_LogicalNot {
			logicalNotData := parseResult.Node.Data.(AstData_UnaryExpression)
			if logicalNotData.Node.Kind == AstKind_Invocation {
				invocationData := logicalNotData.Node.Data.(AstData_Invocation)
				lastParameterIndex := len(invocationData.ParameterNodes) - 1
				lastParameterNode := invocationData.ParameterNodes[lastParameterIndex]
				if lastParameterNode.Kind == AstKind_Struct {
					// remove struct from params
					parseResult.Node.Data = AstData_UnaryExpression{
						Node: AstNode{
							Kind: AstKind_Invocation,
							Data: AstData_Invocation{
								ScriptIdentifierNode:              invocationData.ScriptIdentifierNode,
								ParameterNodes:                    invocationData.ParameterNodes[:lastParameterIndex],
								TokensConsumedByEachParameterNode: invocationData.TokensConsumedByEachParameterNode[:lastParameterIndex],
							},
						},
					}
					// skip backwards over the struct so it will be read as the if-statement body
					*index -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]
					// reduce the number of tokens consumed by the condition
					parseResult.TokensConsumed -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]
				}
			}
		} else if parseResult.Node.Kind == AstKind_Assignment {
			assignmentData := parseResult.Node.Data.(AstData_Assignment)
			if assignmentData.ValueNode.Kind == AstKind_Invocation {
				invocationData := assignmentData.ValueNode.Data.(AstData_Invocation)
				lastParameterIndex := len(invocationData.ParameterNodes) - 1
				lastParameterNode := invocationData.ParameterNodes[lastParameterIndex]
				if lastParameterNode.Kind == AstKind_Struct {
					// remove struct from params
					parseResult.Node.Data = AstData_Assignment{
						NameNode: assignmentData.NameNode,
						ValueNode: AstNode{
							Kind: AstKind_Invocation,
							Data: AstData_Invocation{
								ScriptIdentifierNode:              invocationData.ScriptIdentifierNode,
								ParameterNodes:                    invocationData.ParameterNodes[:lastParameterIndex],
								TokensConsumedByEachParameterNode: invocationData.TokensConsumedByEachParameterNode[:lastParameterIndex],
							},
						},
					}
					// skip backwards over the struct so it will be read as the if-statement body
					*index -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]
					// reduce the number of tokens consumed by the condition
					parseResult.TokensConsumed -= invocationData.TokensConsumedByEachParameterNode[lastParameterIndex]
				}
			}
		}
	}

	ParseScript = func(index int) ParseResult {
		if GetKind(index) != TokenKind_Script {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in script wasn't 'script'",
			}
		}
		index++

		var defaultParameters AstNodeBuffer

		if GetKind(index) != TokenKind_Identifier && GetKind(index) != TokenKind_RawChecksum {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Second token in script wasn't an identifier or a checksum",
			}
		}
		nameToken := GetToken(index)
		index++

		for {
			if GetKind(index) == TokenKind_NewLine {
				index++
			} else if parseResult := ParseAssignment(index, true); parseResult.WasSuccessful {
				index += parseResult.TokensConsumed
				pruneStructIfInvoked(&parseResult, &index)
				defaultParameters.MaybeSave(parseResult)
			} else {
				break
			}
		}

		bodyParseResult, bodyNodes := ParseBodyOfCode(index)
		if !bodyParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        WrapStr("Couldn't parse script body", bodyParseResult.Reason),
			}
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Script,
				Data: AstData_Script{
					NameNode: AstNode{
						Kind: AstKind_Checksum,
						Data: AstData_Checksum{
							ChecksumToken: nameToken,
						},
					},
					DefaultParameterNodes: defaultParameters.Nodes,
					BodyNodes:             bodyNodes,
				},
			},
			TokensConsumed: 1 + defaultParameters.TokensConsumed + bodyParseResult.TokensConsumed + 1,
		}
	}

	ParseWhileLoop = func(index int) ParseResult {
		if GetKind(index) != TokenKind_While {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in while loop wasn't 'while'",
			}
		}
		index++

		bodyParseResult, bodyNodes := ParseBodyOfCode(index)
		if bodyParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: true,
				Node: AstNode{
					Kind: AstKind_WhileLoop,
					Data: AstData_WhileLoop{
						BodyNodes: bodyNodes,
					},
				},
				TokensConsumed: 1 + bodyParseResult.TokensConsumed,
			}
		}

		return ParseResult{
			WasSuccessful: false,
			Reason:        WrapStr("Couldn't parse while loop body", bodyParseResult.Reason),
		}
	}

	ParseLogicalNot = func(index int) ParseResult {
		if GetKind(index) != TokenKind_Bang {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Logical not doesn't start with '!'",
			}
		}
		index++
		expressionParseResult := ParseExpression(index, true)
		if !expressionParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Failed to parse expression after logical-not",
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_LogicalNot,
				Data: AstData_UnaryExpression{
					Node: expressionParseResult.Node,
				},
			},
			TokensConsumed: 1 + expressionParseResult.TokensConsumed,
		}
	}

	ParseIfStatement = func(index int) ParseResult {
		oldIndex := index
		if GetKind(index) != TokenKind_If {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in if-statement wasn't 'if'",
			}
		}

		index++

		booleanInvocationData := make([]bool, 6500)
		conditions := make([]AstNode, 6500)
		numConditions := 0
		saveCondition := func(conditionParseResult ParseResult, isBooleanInvocation bool) {
			booleanInvocationData[numConditions] = isBooleanInvocation
			conditions[numConditions] = conditionParseResult.Node
			numConditions++
		}

		bodies := make([][]AstNode, 6500)
		numBodies := 0
		saveBody := func(bodyParseResult ParseResult, bodyNodes []AstNode) {
			bodies[numBodies] = bodyNodes
			numBodies++
		}

		conditionIsBooleanInvocation := false
		if GetKind(index) == TokenKind_AtSymbol {
			conditionIsBooleanInvocation = true
			index++
		}

		conditionParseResult := ParseExpression(index, true)
		if !conditionParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        WrapStr("Couldn't parse condition in if-statement", conditionParseResult.Reason),
			}
		}
		index += conditionParseResult.TokensConsumed
		pruneStructIfInvoked(&conditionParseResult, &index)
		saveCondition(conditionParseResult, conditionIsBooleanInvocation)

		bodyParseResult, bodyNodes := ParseBodyOfCode(index)
		if !bodyParseResult.WasSuccessful {
			return ParseResult{
				WasSuccessful: false,
				Reason:        WrapStr("Couldn't parse body of if-statement", bodyParseResult.Reason),
			}
		}
		index += bodyParseResult.TokensConsumed
		saveBody(bodyParseResult, bodyNodes)

		for {
			if GetKind(index) == TokenKind_Else {
				index++
				if GetKind(index) == TokenKind_If {
					index++

					nextConditionIsBooleanInvocation := false
					if GetKind(index) == TokenKind_AtSymbol {
						index++
						nextConditionIsBooleanInvocation = true
					}

					anotherConditionParseResult := ParseExpression(index, true)
					if anotherConditionParseResult.WasSuccessful {
						index += anotherConditionParseResult.TokensConsumed
						pruneStructIfInvoked(&anotherConditionParseResult, &index)
						saveCondition(anotherConditionParseResult, nextConditionIsBooleanInvocation)
						if GetKind(index) != TokenKind_LeftCurlyBrace {
							return ParseResult{
								WasSuccessful: false,
								Reason:        WrapStr("Couldn't find '{' after else-if condition", anotherConditionParseResult.Reason),
							}
						}
					} else {
						return ParseResult{
							WasSuccessful: false,
							Reason:        WrapStr("Failed to parse condition for else-if", anotherConditionParseResult.Reason),
						}
					}
				}

				if GetKind(index) == TokenKind_LeftCurlyBrace {
					anotherBodyParseResult, bodyNodes := ParseBodyOfCode(index)
					if bodyParseResult.WasSuccessful {
						index += anotherBodyParseResult.TokensConsumed
						saveBody(anotherBodyParseResult, bodyNodes)
					} else {
						return ParseResult{
							WasSuccessful: false,
							Reason:        WrapStr("Failed to parse an if-statement body", anotherBodyParseResult.Reason),
						}
					}
				} else {
					return ParseResult{
						WasSuccessful: false,
						Reason:        fmt.Sprintf("Unexpected token after 'else' keyword, '%s'. Expected '{' or 'if'.", GetToken(index).Data),
					}
				}
			} else {
				break
			}
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_IfStatement,
				Data: AstData_IfStatement{
					BooleanInvocationData: booleanInvocationData[:numConditions],
					Conditions:            conditions[:numConditions],
					Bodies:                bodies[:numBodies],
				},
			},
			TokensConsumed: index - oldIndex,
		}
	}

	ParseBodyOfCode = func(index int) (ParseResult, []AstNode) {
		startIndex := index

		if GetKind(index) != TokenKind_LeftCurlyBrace {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in body of code wasn't '{'",
			}, []AstNode{}
		}
		index++

		var bodyNodes AstNodeBuffer
		for {
			if GetKind(index) == TokenKind_OutOfRange {
				break
			} else if GetKind(index) == TokenKind_RightCurlyBrace {
				index++
				break
			} else if parseResult := ParseNewLine(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseBreak(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseReturn(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseIfStatement(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseWhileLoop(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseComment(index); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseAssignment(index, true); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else if parseResult := ParseExpression(index, true); parseResult.WasSuccessful {
				bodyNodes.MaybeSave(parseResult)
				index += parseResult.TokensConsumed
			} else {
				TokensNotRecognisedError(parser.Tokens[index:], "a script body node")
			}
		}

		if (index - startIndex) < 2 {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Didn't parse enough tokens to be a body of code",
			}, []AstNode{}
		}

		return ParseResult{
			WasSuccessful:  true,
			TokensConsumed: 2 + bodyNodes.TokensConsumed,
		}, bodyNodes.Nodes
	}

	ParseRandom = func(index int) ParseResult {
		oldIndex := index

		if GetKind(index) != TokenKind_Random {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in 'random' wasn't 'random'",
			}
		}
		index++

		if GetKind(index) != TokenKind_LeftCurlyBrace {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "Second token in 'random' wasn't '{'",
			}
		}
		index++

		var branchWeights []AstNode
		var branches [][]AstNode
		numBranches := 0
		saveBranch := func(weight AstNode, branch []AstNode) {
			branchWeights = append(branchWeights, weight)
			branches = append(branches, branch)
			numBranches++
		}

		for {
			for {
				if GetKind(index) == TokenKind_NewLine ||
					GetKind(index) == TokenKind_SingleLineComment ||
					GetKind(index) == TokenKind_MultiLineComment {
					index++
				} else {
					break
				}
			}

			if GetKind(index) == TokenKind_RightCurlyBrace {
				index++
				break
			}

			integerParseResult := ParseInteger(index)
			if !integerParseResult.WasSuccessful {
				return ParseResult{
					WasSuccessful: false,
					Reason:        WrapStr("Failed to parse weight for branch", integerParseResult.Reason),
				}
			}
			index += integerParseResult.TokensConsumed

			bodyParseResult, bodyNodes := ParseBodyOfCode(index)
			if !bodyParseResult.WasSuccessful {
				return ParseResult{
					WasSuccessful: false,
					Reason:        WrapStr("Failed to parse body for branch", integerParseResult.Reason),
				}
			}
			index += bodyParseResult.TokensConsumed

			saveBranch(integerParseResult.Node, bodyNodes)
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Random,
				Data: AstData_Random{
					BranchWeights: branchWeights[:numBranches],
					Branches:      branches[:numBranches],
				},
			},
			TokensConsumed: index - oldIndex,
		}
	}

	ParseInvocation = func(index int) ParseResult {
		oldIndex := index
		if GetKind(index) != TokenKind_Identifier &&
			GetKind(index) != TokenKind_RawChecksum &&
			GetKind(index) != TokenKind_Return /* TODO(brandon): remove this hack. (ParseReturn() just calls ParseInvocation() because the syntax so similar) */ {
			return ParseResult{
				WasSuccessful: false,
				Reason:        "First token in invocation wasn't an identifier or checksum",
			}
		}
		scriptIdentifierToken := GetToken(index)

		isRawChecksum := false
		if GetKind(index) == TokenKind_RawChecksum {
			isRawChecksum = true
		}

		var parameterNodes AstNodeBuffer
		var tokensConsumedByEachParameterNode []int

		index++
		for { // gather parameters
			if GetKind(index) == TokenKind_BackwardSlash && GetKind(index+1) == TokenKind_NewLine {
				index += 2
			}
			parameterParseResult := ParseInvocationParameter(index)
			if parameterParseResult.WasSuccessful {
				parameterNodes.MaybeSave(parameterParseResult)
				tokensConsumedByEachParameterNode = append(tokensConsumedByEachParameterNode, parameterParseResult.TokensConsumed)
				index += parameterParseResult.TokensConsumed
			} else {
				break
			}
		}

		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Invocation,
				Data: AstData_Invocation{
					ScriptIdentifierNode: AstNode{
						Kind: AstKind_Checksum,
						Data: AstData_Checksum{
							ChecksumToken: scriptIdentifierToken,
							IsRawChecksum: isRawChecksum,
						},
					},
					ParameterNodes:                    parameterNodes.Nodes,
					TokensConsumedByEachParameterNode: tokensConsumedByEachParameterNode,
				},
			},
			TokensConsumed: index - oldIndex,
		}
	}

	ParseInvocationParameter = func(index int) ParseResult {
		if assignmentParseResult := ParseAssignment(index, false); assignmentParseResult.WasSuccessful {
			return assignmentParseResult
		}
		if expressionParseResult := ParseExpression(index, false); expressionParseResult.WasSuccessful {
			return expressionParseResult
		}
		return ParseResult{
			WasSuccessful: false,
			Reason:        TokensNotRecognisedError(parser.Tokens[index:], "a parameter"),
		}
	}

	ParseComment = func(index int) ParseResult {
		if GetKind(index) != TokenKind_SingleLineComment && GetKind(index) != TokenKind_MultiLineComment {
			return ParseResult{
				WasSuccessful: false,
				Reason:        fmt.Sprintf("Not a comment (%+v)", GetToken(index)),
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Comment,
				Data: AstData_Comment{
					CommentToken: GetToken(index),
				},
			},
			TokensConsumed: 1,
		}
	}

	ParseNewLine = func(index int) ParseResult {
		if GetKind(index) != TokenKind_NewLine {
			return ParseResult{
				WasSuccessful: false,
				Reason:        fmt.Sprintf("Not a new-line token (%+v)", GetToken(index)),
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_NewLine,
				Data: AstData_Empty{},
			},
			TokensConsumed: 1,
		}
	}

	ParseBreak = func(index int) ParseResult {
		if GetKind(index) != TokenKind_Break {
			return ParseResult{
				WasSuccessful: false,
				Reason:        fmt.Sprintf("Not a break token (%+v)", GetToken(index)),
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Break,
				Data: AstData_Empty{},
			},
			TokensConsumed: 1,
		}
	}

	ParseReturn = func(index int) ParseResult {
		if GetKind(index) == TokenKind_Return {
			if invocationParseResult := ParseChecksumOrInvocation(index, true); invocationParseResult.WasSuccessful {
				return ParseResult{
					WasSuccessful: true,
					Node: AstNode{
						Kind: AstKind_Return,
						Data: AstData_UnaryExpression{
							Node: invocationParseResult.Node,
						},
					},
					TokensConsumed: invocationParseResult.TokensConsumed,
				}
			}
		}
		return ParseResult{
			WasSuccessful: false,
			Reason:        "Not a return statement",
		}
	}

	ParseComma = func(index int) ParseResult {
		if GetKind(index) != TokenKind_Comma {
			return ParseResult{
				WasSuccessful: false,
				Reason:        fmt.Sprintf("Not a comma token (%+v)", GetToken(index)),
			}
		}
		return ParseResult{
			WasSuccessful: true,
			Node: AstNode{
				Kind: AstKind_Comma,
				Data: AstData_Empty{},
			},
			TokensConsumed: 1,
		}
	}

	GetKind = func(index int) TokenKind {
		return GetToken(index).Kind
	}

	GetToken = func(index int) Token {
		if numOfTokens := len(parser.Tokens); index >= numOfTokens {
			return Token{
				Kind:       TokenKind_OutOfRange,
				Data:       fmt.Sprintf("<index=%d,numOfTokens=%d>", index, numOfTokens),
				LineNumber: -1,
			}
		}
		return parser.Tokens[index]
	}

	parser.Result = ParseRoot()
}

type AstNodeBuffer struct {
	Nodes          []AstNode
	NumNodes       int
	TokensConsumed int
}

func (this *AstNodeBuffer) MaybeSave(parseResult ParseResult) {
	// Even if the node isn't appended, we need to count the number of tokens it consumed
	this.TokensConsumed += parseResult.TokensConsumed

	// Don't store consecutive newlines; they will break the roq decompiler.
	if parseResult.Node.Kind == AstKind_NewLine {
		i := this.NumNodes - 1
		for {
			if i < 0 {
				break
			}
			earlierNode := this.Nodes[i]
			if earlierNode.Kind == AstKind_Comment {
				i--
			} else if earlierNode.Kind == AstKind_NewLine {
				return
			} else {
				break
			}
		}
	}

	this.Nodes = append(this.Nodes, parseResult.Node)
	this.NumNodes++
}

func TokensNotRecognisedError(tokens []Token, notRecognisedAs string) string {
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("Token stream not recognised as %s: [\n", notRecognisedAs))
	for _, token := range tokens {
		messageBuilder.WriteString(fmt.Sprintf("  %+v,\n", token))
	}
	messageBuilder.WriteString("]")
	return messageBuilder.String()
}

func WrapStr(outer, inner string) string {
	return fmt.Sprintf("<%s: %s>", outer, inner)
}
