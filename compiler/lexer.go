package compiler

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	SourceCode     string
	SourceCodeSize int
	Index          int
	LineNumber     int

	Tokens    []Token
	NumTokens int

	StartOfIdentifier int
	StartOfInteger    int
	StartOfString     int
}

func LexSourceCode(lexer *Lexer) { // do lexical analysis (build an array of Tokens)

	CanFindKeywordAtIndex := func(keyword string, index int) bool {
		return strings.HasPrefix(lexer.SourceCode[index:], keyword)
	}

	CanFindKeyword := func(keyword string) bool {
		return CanFindKeywordAtIndex(keyword, lexer.Index)
	}

	CanFindSingleLineComment := func() (string, bool) {
		if CanFindKeyword("//") {
			start := lexer.Index
			end := start + 2
			for {
				if end >= len(lexer.SourceCode) {
					return lexer.SourceCode[start:end], true
				} else if lexer.SourceCode[end] == '\n' {
					return lexer.SourceCode[start:end], true
				}
				end++
			}
		}
		return "", false
	}

	CanFindMultiLineComment := func() (string, bool) {
		if CanFindKeyword("/*") {
			start := lexer.Index
			end := start + 2
			nesting := 1
			for {
				if end >= len(lexer.SourceCode) {
					return lexer.SourceCode[start:end], true
				} else if CanFindKeywordAtIndex("*/", end) {
					nesting--
					if nesting <= 0 {
						end += 2
						return lexer.SourceCode[start:end], true
					}
				} else if CanFindKeywordAtIndex("/*", end) {
					nesting++
				}
				if lexer.SourceCode[end] == '\n' {
					lexer.LineNumber++
				}
				end++
			}
		}
		return "", false
	}

	CanFindRawChecksum := func() (string, bool) {
		start := lexer.Index
		end := start
		if lexer.SourceCode[end] != '#' {
			return "", false
		}
		end++

		for i := 0; i < 8; i++ {
			switch lexer.SourceCode[end] {
			case '0':
				fallthrough
			case '1':
				fallthrough
			case '2':
				fallthrough
			case '3':
				fallthrough
			case '4':
				fallthrough
			case '5':
				fallthrough
			case '6':
				fallthrough
			case '7':
				fallthrough
			case '8':
				fallthrough
			case '9':
				fallthrough
			case 'a':
				fallthrough
			case 'b':
				fallthrough
			case 'c':
				fallthrough
			case 'd':
				fallthrough
			case 'e':
				fallthrough
			case 'f':
				fallthrough
			case 'A':
				fallthrough
			case 'B':
				fallthrough
			case 'C':
				fallthrough
			case 'D':
				fallthrough
			case 'E':
				fallthrough
			case 'F':
				end++
			default:
				return "", false
			}
		}
		return lexer.SourceCode[start:end], true
	}

	CanFindString := func() (string, bool) {
		start := lexer.Index
		end := start

		oldLineNumber := lexer.LineNumber

		stage := 0
		// 0 = looking for (")
		// 1 = found ("), scanning string, checking for next (")
		// 2 = found next (")
		for {
			if end < len(lexer.SourceCode) {
				if lexer.SourceCode[end] == '\n' {
					lexer.LineNumber++
				}
			}
			switch stage {
			case 0:
				if lexer.SourceCode[end] == '"' {
					stage = 1
				} else {
					lexer.LineNumber = oldLineNumber
					return "", false
				}
			case 1:
				if end >= len(lexer.SourceCode) {
					return lexer.SourceCode[start:], true
				} else if lexer.SourceCode[end] == '"' {
					return lexer.SourceCode[start : end+1], true
				}
			}
			end++
		}
	}

	CanFindFloat := func() (string, bool) {
		start := lexer.Index
		end := start

		stage := 0
		// 0 = waiting to scan first digit
		// 1 = first digit found, scanning digits but checking for '.'
		// 2 = '.' found, scanning for first digit after '.'
		// 3 = first digit after '.' found, scanning for more digits
		for {
			switch stage {
			case 0:
				if unicode.IsDigit(rune(lexer.SourceCode[end])) {
					stage = 1
				} else {
					return "", false
				}
			case 1:
				if end >= len(lexer.SourceCode) {
					return lexer.SourceCode[start:end], false
				} else if lexer.SourceCode[end] == '.' {
					stage = 2
				} else if !unicode.IsDigit(rune(lexer.SourceCode[end])) {
					return "", false
				}
			case 2:
				if end >= len(lexer.SourceCode) {
					return lexer.SourceCode[start:end], false
				} else if unicode.IsDigit(rune(lexer.SourceCode[end])) {
					stage = 3
				} else {
					return "", false
				}
			case 3:
				if end >= len(lexer.SourceCode) ||
					!unicode.IsDigit(rune(lexer.SourceCode[end])) {
					return lexer.SourceCode[start:end], true
				}
			}
			end++
		}
	}

	CanFindInteger := func() (string, bool) {
		start := lexer.Index
		end := start
		for {
			if end >= len(lexer.SourceCode) {
				return lexer.SourceCode[start:end], start != end
			}
			if !unicode.IsDigit(rune(lexer.SourceCode[end])) {
				return lexer.SourceCode[start:end], start != end
			}
			end++
		}
	}

	CanFindIdentifier := func() (string, bool) {
		start := lexer.Index
		end := start
		for {
			if end >= len(lexer.SourceCode) {
				break
			}
			if !(unicode.IsLetter(rune(lexer.SourceCode[end])) ||
				unicode.IsDigit(rune(lexer.SourceCode[end])) ||
				lexer.SourceCode[end] == '_') {
				break
			}
			end++
		}

		return lexer.SourceCode[start:end], start != end
	}

	SaveToken := func(lexer *Lexer, kind TokenKind, data string) {
		lexer.Tokens = append(lexer.Tokens, Token{
			Kind: kind,
			Data: data,
			LineNumber: lexer.LineNumber,
		})
		lexer.NumTokens++
	}

	lexer.LineNumber = 1
	for {
		if lexer.Index >= lexer.SourceCodeSize {
			break
		}

		if data, found := CanFindFloat(); found {
			SaveToken(lexer, TokenKind_Float, data)
			lexer.Index += len(data)
		} else if data, found := CanFindInteger(); found {
			SaveToken(lexer, TokenKind_Integer, data)
			lexer.Index += len(data)
		} else if data, found := CanFindString(); found {
			SaveToken(lexer, TokenKind_String, data)
			lexer.Index += len(data)
		} else if data, found := CanFindSingleLineComment(); found {
			SaveToken(lexer, TokenKind_SingleLineComment, data)
			lexer.Index += len(data)
		} else if data, found := CanFindMultiLineComment(); found {
			SaveToken(lexer, TokenKind_MultiLineComment, data)
			lexer.Index += len(data)
		} else if data, found := CanFindRawChecksum(); found {
			SaveToken(lexer, TokenKind_RawChecksum, data)
			lexer.Index += len(data)
		} else {
			// Check for single-character tokens
			switch lexer.SourceCode[lexer.Index] {
			case '\t':
				fallthrough
			case ' ':
				lexer.Index++
			case '\r':
				fallthrough
			case '\n':
				lexer.LineNumber++
				SaveToken(lexer, TokenKind_NewLine, "\\n")
				lexer.Index++
			case '=':
				SaveToken(lexer, TokenKind_Equals, "=")
				lexer.Index++
			case '@':
				SaveToken(lexer, TokenKind_AtSymbol, "@")
				lexer.Index++
			case '[':
				SaveToken(lexer, TokenKind_LeftSquareBracket, "[")
				lexer.Index++
			case ']':
				SaveToken(lexer, TokenKind_RightSquareBracket, "]")
				lexer.Index++
			case '{':
				SaveToken(lexer, TokenKind_LeftCurlyBrace, "{")
				lexer.Index++
			case '}':
				SaveToken(lexer, TokenKind_RightCurlyBrace, "}")
				lexer.Index++
			case '(':
				SaveToken(lexer, TokenKind_LeftParenthesis, "(")
				lexer.Index++
			case ')':
				SaveToken(lexer, TokenKind_RightParenthesis, ")")
				lexer.Index++
			case '<':
				SaveToken(lexer, TokenKind_LeftAngleBracket, "<")
				lexer.Index++
			case '>':
				SaveToken(lexer, TokenKind_RightAngleBracket, ">")
				lexer.Index++
			case '+':
				SaveToken(lexer, TokenKind_Plus, "+")
				lexer.Index++
			case '-':
				SaveToken(lexer, TokenKind_Minus, "-")
				lexer.Index++
			case '*':
				SaveToken(lexer, TokenKind_Asterisk, "*")
				lexer.Index++
			case '/':
				SaveToken(lexer, TokenKind_ForwardSlash, "/")
				lexer.Index++
			case '\\':
				SaveToken(lexer, TokenKind_BackwardSlash, "\\")
				lexer.Index++
			case ',':
				SaveToken(lexer, TokenKind_Comma, ",")
				lexer.Index++
			case '.':
				SaveToken(lexer, TokenKind_Dot, ".")
				lexer.Index++
			case '!':
				SaveToken(lexer, TokenKind_Bang, "!")
				lexer.Index++
			case ':':
				SaveToken(lexer, TokenKind_Colon, ":")
				lexer.Index++
			default:
				// Check for multi-character tokens
				if CanFindKeyword("or") {
					SaveToken(lexer, TokenKind_Or, "or")
					lexer.Index += 2
				} else if CanFindKeyword("if") {
					SaveToken(lexer, TokenKind_If, "if")
					lexer.Index += 2
				} else if CanFindKeyword("and") {
					SaveToken(lexer, TokenKind_And, "and")
					lexer.Index += 3
				} else if CanFindKeyword("else") {
					SaveToken(lexer, TokenKind_Else, "else")
					lexer.Index += 4
				} else if CanFindKeyword("while") {
					SaveToken(lexer, TokenKind_While, "while")
					lexer.Index += 5
				} else if CanFindKeyword("break") {
					SaveToken(lexer, TokenKind_Break, "break")
					lexer.Index += 5
				} else if CanFindKeyword("script") {
					SaveToken(lexer, TokenKind_Script, "script")
					lexer.Index += 6
				} else if CanFindKeyword("random") {
					SaveToken(lexer, TokenKind_Random, "random")
					lexer.Index += 6
				} else if CanFindKeyword("return") {
					SaveToken(lexer, TokenKind_Return, "return")
					lexer.Index += 6
				} else if identifier, found := CanFindIdentifier(); found {
					SaveToken(lexer, TokenKind_Identifier, identifier)
					lexer.Index += len(identifier)
				} else {
					character := lexer.SourceCode[lexer.Index]
					fmt.Printf("\nLexer failed at character: '%c' (%#x)...\n", character, character)
					fmt.Println("\nRegistered tokens:")
					for i, token := range lexer.Tokens {
						if i >= lexer.NumTokens {
							break
						}
						fmt.Printf("%+v\n", token)
					}
					fmt.Println()
					return
				}
			}
		}
	}
	lexer.Tokens = lexer.Tokens[:lexer.NumTokens]
}
