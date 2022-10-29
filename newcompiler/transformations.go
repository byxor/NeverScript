package newcompiler

// Examples:
// --------------------------------------------------
//
// `x += 5`                    -->  `x = (x + 5)`
//
// `<speed> *= 0.5`            -->  `<speed> = (<speed> * 0.5)`
//
// `Change mrs_ramos /= 2022`  -->  `Change mrs_ramos = (mrs_ramos / 2022)`
//
func simplifyInPlaceOperation(node Node, nodeKind NodeKind) (Node, error) {
    operation := node.(manyWrappedNodes)

    leftHandSide := operation.nodeLists[0][0]
    rightHandSide := operation.nodeLists[0][1]

    newRightHandSide := wrappedNode{
        kind: NodeKind_ParenthesisOperation,
        node: manyWrappedNodes{
            kind: nodeKind,
            nodeLists: [][]Node{
                {leftHandSide, rightHandSide},
                {},
                {},
            },
            extraTokensConsumed: 0,
        },
        extraTokensConsumed: 0,
    }

    simplifiedOperation := manyWrappedNodes{
        kind: NodeKind_AssignmentOperation,
        nodeLists: [][]Node{
            {
                operation.nodeLists[0][0],
                newRightHandSide,
            },
            {},
            {},
        },
        extraTokensConsumed: 0,
    }

    return fixedSizeWrappedNode{
        node:           simplifiedOperation,
        tokensConsumed: operation.TokensConsumed(),
    }, nil
}

/*
Removes all "else-if" nodes from an if statement so they only use "if" and "else".

```
if (firstCondition) {
	doFirst
	x++
} else if (secondCondition) {
	doSecond
	y++
} else {
	doSomethingElse
}
```

becomes

```
if (firstCondition) {
	doFirst
	x++
} else {
	if (secondCondition) {
		doSecond
		y++
	} else {
		doSomethingElse
	}
}

```
*/
func simplifyElseIfChain(if_ Node, elseIfs []Node, else_ Node) (Node, error) {
    if len(elseIfs) == 0 {
        return wrappedNodes{
            kind:                NodeKind_IfStatement,
            nodes:               []Node{if_, else_},
            extraTokensConsumed: 0,
        }, nil
    }
    newElseBody, err := simplifyElseIfChain(elseIfs[0].(wrappedNode).node, elseIfs[1:], else_)
    if err != nil {
        return nil, err
    }
    newElse := wrappedNodes{
        kind:                NodeKind_Else,
        nodes:               newElseBody.(wrappedNodes).nodes,
        extraTokensConsumed: 1,
    }
    ifStatement := wrappedNodes{
        kind: NodeKind_IfStatement,
        nodes: []Node{
            if_,
            newElse,
        },
        extraTokensConsumed: 0,
    }
    return ifStatement, nil
}

// `x != y`
// becomes
// `!(x == y)`
func simplifyInequalityOperation(inequalityOperation Node) (Node, error) {
    return wrappedNode{
        kind: NodeKind_NotOperation,
        node: wrappedNode{
            kind: NodeKind_ParenthesisOperation,
            node: manyWrappedNodes{
                kind: NodeKind_EqualityOperation,
                nodeLists: [][]Node{
                    {
                        inequalityOperation.(manyWrappedNodes).nodeLists[0][0],
                        inequalityOperation.(manyWrappedNodes).nodeLists[0][1],
                    },
                    {},
                    {},
                },
                extraTokensConsumed: 2,
            },
            extraTokensConsumed: 0,
        },
        extraTokensConsumed: 0,
    }, nil
}

// `x <= y`
// becomes
// `!(x > y)`
func simplifyLessThanEqualOperation(lessThanEqualOperation Node) (Node, error) {
    return wrappedNode{
        kind: NodeKind_NotOperation,
        node: wrappedNode{
            kind: NodeKind_ParenthesisOperation,
            node: manyWrappedNodes{
                kind: NodeKind_GreaterThanOperation,
                nodeLists: [][]Node{
                    {
                        lessThanEqualOperation.(manyWrappedNodes).nodeLists[0][0],
                        lessThanEqualOperation.(manyWrappedNodes).nodeLists[0][1],
                    },
                    {},
                    {},
                },
                extraTokensConsumed: 2,
            },
            extraTokensConsumed: 0,
        },
        extraTokensConsumed: 0,
    }, nil
}

// `x >= y`
// becomes
// `!(x < y)`
func simplifyGreaterThanEqualOperation(greaterThanEqual Node) (Node, error) {
    return wrappedNode{
        kind: NodeKind_NotOperation,
        node: wrappedNode{
            kind: NodeKind_ParenthesisOperation,
            node: manyWrappedNodes{
                kind: NodeKind_LessThanOperation,
                nodeLists: [][]Node{
                    {
                        greaterThanEqual.(manyWrappedNodes).nodeLists[0][0],
                        greaterThanEqual.(manyWrappedNodes).nodeLists[0][1],
                    },
                    {},
                    {},
                },
                extraTokensConsumed: 2,
            },
            extraTokensConsumed: 0,
        },
        extraTokensConsumed: 0,
    }, nil
}
