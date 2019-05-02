package compiler

import (
	"github.com/alecthomas/participle"
	"github.com/byxor/NeverScript"
	"github.com/pkg/errors"
)

func buildSyntaxTreeFrom(sourceCode NeverScript.SourceCode) (syntaxTree, error) {
	syntaxTree := syntaxTree{}

	parser := participle.MustBuild(
		&syntaxTree,
		participle.Lexer(nsLexer),
		participle.UseLookahead(2),
	)

	err := parser.ParseString(sourceCode.ToString(), &syntaxTree)
	if err != nil {
		return syntaxTree, errors.Wrap(err, "Failed to run participle")
	}

	return syntaxTree, nil
}
