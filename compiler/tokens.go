package compiler

type Token struct {
	Kind       TokenKind
	Data       string
	LineNumber int
}

type TokenKind int

const (
	TokenKind_Identifier TokenKind = iota
	TokenKind_Equals
	TokenKind_Integer
	TokenKind_String
	TokenKind_LeftSquareBracket
	TokenKind_RightSquareBracket
	TokenKind_Comma
	TokenKind_NewLine
	TokenKind_Colon
	TokenKind_LeftCurlyBrace
	TokenKind_RightCurlyBrace
	TokenKind_LeftParenthesis
	TokenKind_RightParenthesis
	TokenKind_LeftAngleBracket
	TokenKind_RightAngleBracket
	TokenKind_RawChecksum
	TokenKind_Plus
	TokenKind_Minus
	TokenKind_Asterisk
	TokenKind_ForwardSlash
	TokenKind_BackwardSlash
	TokenKind_Bang
	TokenKind_Random
	TokenKind_If
	TokenKind_Else
	TokenKind_While
	TokenKind_Break
	TokenKind_Return
	TokenKind_Script
	TokenKind_Float
	TokenKind_SingleLineComment
	TokenKind_MultiLineComment
	TokenKind_Dot
	TokenKind_And
	TokenKind_Or
	TokenKind_OutOfRange
)

func (tokenKind TokenKind) String() string {
	return [...]string{
		"TokenKind_Identifier",
		"TokenKind_Equals",
		"TokenKind_Integer",
		"TokenKind_String",
		"TokenKind_LeftSquareBracket",
		"TokenKind_RightSquareBracket",
		"TokenKind_Comma",
		"TokenKind_NewLine",
		"TokenKind_Colon",
		"TokenKind_LeftCurlyBrace",
		"TokenKind_RightCurlyBrace",
		"TokenKind_LeftParenthesis",
		"TokenKind_RightParenthesis",
		"TokenKind_LeftAngleBracket",
		"TokenKind_RightAngleBracket",
		"TokenKind_RawChecksum",
		"TokenKind_Plus",
		"TokenKind_Minus",
		"TokenKind_Asterisk",
		"TokenKind_ForwardSlash",
		"TokenKind_BackwardSlash",
		"TokenKind_Bang",
		"TokenKind_Random",
		"TokenKind_If",
		"TokenKind_Else",
		"TokenKind_While",
		"TokenKind_Break",
		"TokenKind_Return",
		"TokenKind_Script",
		"TokenKind_Float",
		"TokenKind_SingleLineComment",
		"TokenKind_MultiLineComment",
		"TokenKind_Dot",
		"TokenKind_And",
		"TokenKind_Or",
		"TokenKind_OutOfRange",
	}[tokenKind]
}
