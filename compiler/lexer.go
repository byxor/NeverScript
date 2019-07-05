package compiler

import (
	"github.com/alecthomas/participle/lexer"
)

var nsLexer = lexer.Must(lexer.Regexp(
	`(?m)` +
		`(\s+)` +
		`` +
		`|(?P<FunctionKeyword>func)` +
		`` +
		`|(?P<Boolean>true|false)` +
		`|(?P<String>"(?:\\.|[^"])*")` +
		`` +
		`|(?P<Integer_Base16>0x[0-9a-fA-F]+)` +
		`|(?P<Integer_Base8>0o[0-7]+)` +
		`|(?P<Integer_Base2>0b[0-1]+)` +
		`|(?P<Integer_Base10>[0-9]+)` +
		`` +
		`|(?P<Identifier>[a-zA-Z][a-zA-Z_\d]*)` +
		`` +
		`|(?P<OpenCurlyBrace>\{)` +
		`|(?P<CloseCurlyBrace>\})` +
		`|(?P<Semicolon>;)` +
		`|(?P<Equals>=)` +
		``,
))
