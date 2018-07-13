package code

import (
	"fmt"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strings"
)

func GenerateUsing(tokens []Token) string {
	cleanTokens := clean(tokens)
	state := stateHolder{
		tokens:     cleanTokens,
		token:      Token{Invalid, nil},
		index:      0,
		names:      nametable.BuildFrom(cleanTokens),
		arrayDepth: 0,
	}
	return strings.TrimSpace(generateUsing(&state))
}

func generateUsing(state *stateHolder) string {
	if len(state.tokens) == 0 {
		return ""
	}

	if state.index >= len(state.tokens) {
		return ""
	}

	state.token = state.tokens[state.index]

	evaluator, ok := evaluators[state.token.Type]
	if !ok {
		panic(fmt.Sprintf("No evaluator found for %s tokens!", state.token.Type.String()))
	}

	atom := evaluator(state)
	tweaked := makeStatefulAdjustments(atom, state)
	state.index++
	return tweaked + generateUsing(state)
}

func makeStatefulAdjustments(atom string, state *stateHolder) string {
	output := atom
	if state.arrayDepth > 0 {
		notTooHigh := state.index < len(state.tokens)-1
		notTooLow := state.index > 0

		if notTooHigh && notTooLow {
			notFirstElement := state.token.Type != StartOfArray
			notLastElement := state.tokens[state.index+1].Type != EndOfArray

			if notFirstElement && notLastElement {
				output += " "
			}
		}
	}
	return output
}

type stateHolder struct {
	tokens     []Token
	token      Token
	index      int
	names      nametable.NameTable
	arrayDepth int
}
