package code

import (
	"fmt"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strings"
)

func GenerateUsing(tokens []Token) string {
	return strings.TrimSpace(
		generateUsing(clean(tokens), nametable.BuildFrom(tokens)),
	)
}

func generateUsing(tokens []Token, nameTable nametable.NameTable) string {
	if len(tokens) == 0 {
		return ""
	}
	if evaluator, ok := evaluators[tokens[0].Type]; ok {
		result := evaluator(tokens[0].Chunk, nameTable)
		return result + generateUsing(tokens[1:], nameTable)
	} else {
		panic(fmt.Sprintf("No evaluator found for %s tokens!", tokens[0].Type.String()))
	}
}
