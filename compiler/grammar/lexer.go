package grammar

import (
	"github.com/alecthomas/participle/lexer"
)

// Some of this lexer was copied and pasted from the participle repository.
// See: _examples/ini/main.go
//
// Some pieces may be irrelevant; feel free to remove them.
//
var NsLexer = lexer.Must(lexer.Regexp(
	`(?m)` +
		`(\s+)` +
		`` +
		`|(?P<Boolean>true|false)` +
		`|(?P<Integer_Base16>0x[0-9a-fA-F]+)` +
		`|(?P<Integer_Base10>[0-9]+)` +
	    `` + 
		`|(?P<Identifier>[a-zA-Z][a-zA-Z_\d]*)` +
		// `|(?P<String>"(?:\\.|[^"])*")` +
		// `|(?P<Float>\d+(?:\.\d+)?)` +
		// `|(?P<Punct>[][=])` +
		`` +
		`|(?P<Semicolon>;)` +
		`|(?P<Equals>=)` +
		``,
))
