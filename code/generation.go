package code

import (
	"fmt"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strings"
)

func GenerateUsing(tokens []Token) string {
	cleanTokens := clean(tokens)
	state := stateMap{
		"tokens":  cleanTokens,
		"token":   nil,
		"names":   nametable.BuildFrom(cleanTokens),
		"index":   0,
		"inArray": 0,
	}
	return strings.TrimSpace(generateUsing(state))
}

func generateUsing(state stateMap) string {
	tokens := state["tokens"].([]Token)

	if len(tokens) == 0 {
		return ""
	}

	index := state["index"].(int)

	if index >= len(tokens) {
		return ""
	}

	token := tokens[index]
	state["token"] = token

	if evaluator, ok := evaluators[token.Type]; ok {
		result := evaluator(state)
		tweaked := makeStatefulAdjustments(result, state)
		state["index"] = index + 1
		return tweaked + generateUsing(state)
	} else {
		panic(fmt.Sprintf("No evaluator found for %s tokens!", tokens[0].Type.String()))
	}
}

func makeStatefulAdjustments(result string, state stateMap) string {
	tokens := state["tokens"].([]Token)
	token := state["token"].(Token)
	index := state["index"].(int)
	if state["inArray"].(int) > 0 {
		inRange := index < len(tokens)-1 && index > 0
		if inRange {
			notFirstElement := token.Type != StartOfArray
			notLastElement := tokens[index+1].Type != EndOfArray
			if notFirstElement && notLastElement {
				result += " "
			}
		}
	}
	return result
}

type stateMap map[string]interface{}
