package newcompiler

import (
    "errors"
    "fmt"
    "strconv"
)

type NodeKind int

const (
    NodeKind_Case NodeKind = iota
    NodeKind_Switch
    NodeKind_Return
    NodeKind_Break
    NodeKind_Loop
    NodeKind_Else
    NodeKind_ElseIf
    NodeKind_If
    NodeKind_IfStatement
    NodeKind_RandomEntry
    NodeKind_Random
    NodeKind_Bytes
    NodeKind_Byte
    NodeKind_ScriptHeader
    NodeKind_Script
    NodeKind_Struct
    NodeKind_Array
    NodeKind_Pair
    NodeKind_Vector
    NodeKind_SubExpression
    NodeKind_AllArguments
    NodeKind_Comma
    NodeKind_Int
    NodeKind_Float
    NodeKind_String
    NodeKind_LocalQbKey
    NodeKind_QbKey
    NodeKind_RawQbKey
    NodeKind_Operation
    NodeKind_NotOperation
    NodeKind_UnaryMinusOperation
    NodeKind_ParenthesisOperation
    NodeKind_PlusOperation
    NodeKind_MinusOperation
    NodeKind_DivideOperation
    NodeKind_MultiplyOperation
    NodeKind_AssignmentOperation
    NodeKind_EqualityOperation
    NodeKind_InequalityOperation
    NodeKind_PlusEqualOperation
    NodeKind_MinusEqualOperation
    NodeKind_DivideEqualOperation
    NodeKind_MultiplyEqualOperation
    NodeKind_GreaterThanOperation
    NodeKind_LessThanOperation
    NodeKind_GreaterThanEqualOperation
    NodeKind_LessThanEqualOperation
    NodeKind_AndOperation
    NodeKind_OrOperation
    NodeKind_ColonOperation
    NodeKind_DotOperation
    NodeKind_ArrayAccessOperation
    NodeKind_ChunkOfCode
    NodeKind_Expression
    NodeKind_SuperExpression
    NodeKind_LineBreak
    NodeKind_Program
)

type Node interface {
    Kind() NodeKind
    TokensConsumed() uint
    LineNumber() uint
}

func Parse(tokens []Token) (Node, error) {
    var parser parser
    parser.tokens = tokens
    parser.compressAST = true
    parser.expressionIndexReference.indices = make(map[uint]bool)
    parser.superExpressionCache = make(map[uint]Node)
    parser.operationCache = make(map[uint]Node)
    parser.expressionCache = make(map[uint]Node)
    parser.subExpressionCache = make(map[uint]Node)
    return parser.parse()
}

type basicNode struct {
    kind           NodeKind
    data           string
    tokensConsumed uint
    lineNumber     uint
}

func (this basicNode) Kind() NodeKind {
    return this.kind
}

func (this basicNode) TokensConsumed() uint {
    return this.tokensConsumed
}

func (this basicNode) LineNumber() uint {
    return this.lineNumber
}

type wrappedNode struct {
    kind                NodeKind
    node                Node
    extraTokensConsumed uint
}

func (this wrappedNode) Kind() NodeKind {
    return this.kind
}

func (this wrappedNode) TokensConsumed() uint {
    return this.node.TokensConsumed() + this.extraTokensConsumed
}

func (this wrappedNode) LineNumber() uint {
    return this.node.LineNumber()
}

type wrappedNodes struct {
    kind                NodeKind
    nodes               []Node
    extraTokensConsumed uint
}

func (this wrappedNodes) Kind() NodeKind {
    return this.kind
}

func (this wrappedNodes) TokensConsumed() uint {
    tokensConsumed := uint(0)
    for _, node := range this.nodes {
        if node != nil {
            tokensConsumed += node.TokensConsumed()
        }
    }
    return tokensConsumed + this.extraTokensConsumed
}

func (this wrappedNodes) LineNumber() uint {
    return this.nodes[0].LineNumber()
}

type manyWrappedNodes struct {
    kind                NodeKind
    nodeLists           [][]Node
    extraTokensConsumed uint
}

func (this manyWrappedNodes) Kind() NodeKind {
    return this.kind
}

func (this manyWrappedNodes) TokensConsumed() uint {
    tokensConsumed := uint(0)
    for _, nodeList := range this.nodeLists {
        if nodeList != nil {
            for _, node := range nodeList {
                if node != nil {
                    tokensConsumed += node.TokensConsumed()
                }
            }
        }
    }
    return tokensConsumed + this.extraTokensConsumed
}

func (this manyWrappedNodes) LineNumber() uint {
    return this.nodeLists[0][0].TokensConsumed()
}

type fixedSizeWrappedNode struct {
    node           manyWrappedNodes
    tokensConsumed uint
}

func (this fixedSizeWrappedNode) Kind() NodeKind {
    return this.node.Kind()
}

func (this fixedSizeWrappedNode) TokensConsumed() uint {
    return this.tokensConsumed
}

func (this fixedSizeWrappedNode) LineNumber() uint {
    return this.node.TokensConsumed()
}

type rawQbKeyNode struct {
    key        []byte
    lineNumber uint
}

func (this rawQbKeyNode) Kind() NodeKind {
    return NodeKind_RawQbKey
}

func (this rawQbKeyNode) TokensConsumed() uint {
    return 1
}

func (this rawQbKeyNode) LineNumber() uint {
    return this.lineNumber
}

type nodeArray struct {
    nodes          []Node
    tokensConsumed uint
}

func (this *nodeArray) save(node Node) {
    this.nodes = append(this.nodes, node)
    this.tokensConsumed += node.TokensConsumed()
}

type indexReference struct {
    indices map[uint]bool
}

func (this indexReference) contains(index uint) bool {
    if b, found := this.indices[index]; found {
        return b
    }
    return false
}

func (this *indexReference) push(index uint) error {
    if this.contains(index) {
        return errors.New(fmt.Sprintf("index reference already contains index '%d'", index))
    }
    this.indices[index] = true
    return nil
}

func (this *indexReference) pop(index uint) error {
    if !this.contains(index) {
        return errors.New(fmt.Sprintf("index reference does not contain index '%d'", index))
    }
    this.indices[index] = false
    return nil
}

type parser struct {
    tokens                   []Token
    compressAST              bool
    expressionIndexReference indexReference
    superExpressionCache     map[uint]Node
    operationCache           map[uint]Node
    expressionCache          map[uint]Node
    subExpressionCache       map[uint]Node
}

func (this *parser) isOutOfRangeAt(index uint) bool {
    return index >= uint(len(this.tokens))
}

func (this *parser) parse() (Node, error) {
    return this.tryParseProgram()
}

// ChunkOfCode
func (this *parser) tryParseProgram() (Node, error) {
    chunkOfCode, err := this.tryParseChunkOfCodeAt(0)

    if err != nil {
        return nil, err
    } else if chunkOfCode == nil {
        return nil, errors.New("couldn't parse program")
    }

    chunkOfCode_ := chunkOfCode.(wrappedNodes)

    if chunkOfCode.TokensConsumed() != uint(len(this.tokens)) {
        return nil, errors.New("parser didn't consume all tokens")
    }

    return wrappedNodes{
        kind:                NodeKind_Program,
        nodes:               notNilNodes(chunkOfCode_.nodes),
        extraTokensConsumed: 0,
    }, nil
}

// (SuperExpression | LineBreak)*
func (this *parser) tryParseChunkOfCodeAt(index uint) (Node, error) {
    startIndex := index
    var nodes nodeArray
    for {
        index = startIndex + nodes.tokensConsumed

        if this.isOutOfRangeAt(index) {
            break
        }

        superExpression, err := this.tryParseSuperExpressionAt(index)
        if err != nil {
            return nil, err
        } else if superExpression != nil {
            nodes.save(superExpression)
            continue
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            nodes.save(lineBreak)
            continue
        }

        break
    }

    return wrappedNodes{
        kind:                NodeKind_ChunkOfCode,
        nodes:               notNilNodes(nodes.nodes),
        extraTokensConsumed: 0,
    }, nil
}

// IfStatement | Loop | Switch | "break" | Return | Expression
func (this *parser) tryParseSuperExpressionAt(index uint) (Node, error) {
    cachedNode, found := this.superExpressionCache[index]
    if found {
        return cachedNode, nil
    }
    node, err := this._tryParseSuperExpressionAt(index)
    if err != nil {
        return nil, err
    }
    this.superExpressionCache[index] = node
    return node, nil
}

func (this *parser) _tryParseSuperExpressionAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    var ifStatement, loop, switch_, break_, return_, expression Node
    var err error

    node := wrappedNode{
        kind:                NodeKind_SuperExpression,
        node:                nil,
        extraTokensConsumed: 0,
    }

    ifStatement, err = this.tryParseIfStatementAt(index)
    if err != nil {
        return nil, err
    } else if ifStatement != nil {
        node.node = ifStatement
        goto foundNode
    }

    loop, err = this.tryParseLoopAt(index)
    if err != nil {
        return nil, err
    } else if loop != nil {
        node.node = loop
        goto foundNode
    }

    switch_, err = this.tryParseSwitchAt(index)
    if err != nil {
        return nil, err
    } else if switch_ != nil {
        node.node = switch_
        goto foundNode
    }

    break_, err = this.tryParseBreakAt(index)
    if err != nil {
        return nil, err
    } else if break_ != nil {
        node.node = break_
        goto foundNode
    }

    return_, err = this.tryParseReturnAt(index)
    if err != nil {
        return nil, err
    } else if return_ != nil {
        node.node = return_
        goto foundNode
    }

    expression, err = this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expression != nil {
        node.node = expression
        goto foundNode
    }

    return nil, nil

foundNode:
    if this.compressAST {
        return node.node, nil
    } else {
        return node, nil
    }
}

// "\r" | "\n"
func (this *parser) tryParseLineBreakAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]
    if token.Kind() != TokenKind_CarriageReturn && token.Kind() != TokenKind_NewLine {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_LineBreak,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// Operation | SubExpression
func (this *parser) tryParseExpressionAt(index uint) (Node, error) {
    cachedNode, found := this.expressionCache[index]
    if found {
        return cachedNode, nil
    }
    node, err := this._tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    }
    this.expressionCache[index] = node
    return node, nil
}

func (this *parser) _tryParseExpressionAt(index uint) (Node, error) {
    // WARNING: INFINITE LEFT RECURSION?
    if this.expressionIndexReference.contains(index) {
        //log.Printf("WARNING: prevented left-recursion in `tryParseExpressionAt` at index %d\n", index)
        return nil, nil
    } else {
        err := this.expressionIndexReference.push(index)
        if err != nil {
            return nil, err
        }
    }
    // -------------------------------------------------

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    var operation, subExpression Node
    var err error

    var node = wrappedNode{
        kind:                NodeKind_Expression,
        node:                nil,
        extraTokensConsumed: 0,
    }

    operation, err = this.tryParseOperationAt(index)
    if err != nil {
        return nil, err
    } else if operation != nil {
        node.node = operation
        goto foundNode
    }

    subExpression, err = this.tryParseSubExpressionAt(index)
    if err != nil {
        return nil, err
    } else if subExpression != nil {
        node.node = subExpression
        goto foundNode
    }

    return nil, nil

foundNode:
    this.expressionIndexReference.pop(index)
    if this.compressAST {
        return node.node, nil
    } else {
        return node, nil
    }
}

/*
	"!" Expression |
	"(" Expression ")" |

	SubExpression "+=" Expression |
	SubExpression "-=" Expression |
	SubExpression "/=" Expression |
	SubExpression "*=" Expression |

	SubExpression "+" Expression |
	SubExpression "-" Expression |
	SubExpression "/" Expression |
	SubExpression "*" Expression |

	SubExpression "==" Expression |
	SubExpression "=" Expression | // <-- assignment hack?
	SubExpression "!=" Expression |

	SubExpression ">" Expression |
	SubExpression "<" Expression |
	SubExpression "<=" Expression |
	SubExpression ">=" Expression |

	SubExpression "and" Expression |
	SubExpression "or" Expression |

	SubExpression ":" Expression |
	SubExpression "." QbKey |
	SubExpression "[" Expression "]"
*/
func (this *parser) tryParseOperationAt(index uint) (Node, error) {
    cachedNode, found := this.operationCache[index]
    if found {
        return cachedNode, nil
    }
    node, err := this._tryParseOperationAt(index)
    if err != nil {
        return nil, err
    }
    this.operationCache[index] = node
    return node, nil
}

func (this *parser) _tryParseOperationAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    var notOperation Node
    var unaryMinusOperation Node
    var parenthesisOperation Node
    var plusOperation Node
    var minusOperation Node
    var divideOperation Node
    var multiplyOperation Node
    var assignmentOperation Node
    var equalityOperation Node
    var inequalityOperation Node
    var plusEqualOperation Node
    var minusEqualOperation Node
    var divideEqualOperation Node
    var multiplyEqualOperation Node
    var greaterThanOperation Node
    var lessThanOperation Node
    var lessThanEqualOperation Node
    var greaterThanEqualOperation Node
    var andOperation Node
    var orOperation Node
    var colonOperation Node
    var dotOperation Node
    var arrayAccessOperation Node
    var err error

    node := wrappedNode{
        kind:                NodeKind_Operation,
        node:                nil,
        extraTokensConsumed: 0,
    }

    notOperation, err = this.tryParseNotOperationAt(index)
    if err != nil {
        return nil, err
    } else if notOperation != nil {
        node.node = notOperation
        goto foundNode
    }

    unaryMinusOperation, err = this.tryParseUnaryMinusOperationAt(index)
    if err != nil {
        return nil, err
    } else if unaryMinusOperation != nil {
        node.node = unaryMinusOperation
        goto foundNode
    }

    parenthesisOperation, err = this.tryParseParenthesisOperationAt(index)
    if err != nil {
        return nil, err
    } else if parenthesisOperation != nil {
        node.node = parenthesisOperation
        goto foundNode
    }

    plusEqualOperation, err = this.tryParsePlusEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if plusEqualOperation != nil {
        node.node = plusEqualOperation
        goto foundNode
    }

    minusEqualOperation, err = this.tryParseMinusEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if minusEqualOperation != nil {
        node.node = minusEqualOperation
        goto foundNode
    }

    divideEqualOperation, err = this.tryParseDivideEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if divideEqualOperation != nil {
        node.node = divideEqualOperation
        goto foundNode
    }

    multiplyEqualOperation, err = this.tryParseMultiplyEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if multiplyEqualOperation != nil {
        node.node = multiplyEqualOperation
        goto foundNode
    }

    plusOperation, err = this.tryParsePlusOperationAt(index)
    if err != nil {
        return nil, err
    } else if plusOperation != nil {
        node.node = plusOperation
        goto foundNode
    }

    minusOperation, err = this.tryParseMinusOperationAt(index)
    if err != nil {
        return nil, err
    } else if minusOperation != nil {
        node.node = minusOperation
        goto foundNode
    }

    divideOperation, err = this.tryParseDivideOperationAt(index)
    if err != nil {
        return nil, err
    } else if divideOperation != nil {
        node.node = divideOperation
        goto foundNode
    }

    multiplyOperation, err = this.tryParseMultiplyOperationAt(index)
    if err != nil {
        return nil, err
    } else if multiplyOperation != nil {
        node.node = multiplyOperation
        goto foundNode
    }

    equalityOperation, err = this.tryParseEqualityOperationAt(index)
    if err != nil {
        return nil, err
    } else if equalityOperation != nil {
        node.node = equalityOperation
        goto foundNode
    }

    assignmentOperation, err = this.tryParseAssignmentOperationAt(index)
    if err != nil {
        return nil, err
    } else if assignmentOperation != nil {
        node.node = assignmentOperation
        goto foundNode
    }

    inequalityOperation, err = this.tryParseInequalityOperationAt(index)
    if err != nil {
        return nil, err
    } else if inequalityOperation != nil {
        node.node = inequalityOperation
        goto foundNode
    }

    greaterThanOperation, err = this.tryParseGreaterThanOperationAt(index)
    if err != nil {
        return nil, err
    } else if greaterThanOperation != nil {
        node.node = greaterThanOperation
        goto foundNode
    }

    lessThanOperation, err = this.tryParseLessThanOperationAt(index)
    if err != nil {
        return nil, err
    } else if lessThanOperation != nil {
        node.node = lessThanOperation
        goto foundNode
    }

    lessThanEqualOperation, err = this.tryParseLessThanEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if lessThanEqualOperation != nil {
        node.node = lessThanEqualOperation
        goto foundNode
    }

    greaterThanEqualOperation, err = this.tryParseGreaterThanEqualOperationAt(index)
    if err != nil {
        return nil, err
    } else if greaterThanEqualOperation != nil {
        node.node = greaterThanEqualOperation
        goto foundNode
    }

    andOperation, err = this.tryParseAndOperationAt(index)
    if err != nil {
        return nil, err
    } else if andOperation != nil {
        node.node = andOperation
        goto foundNode
    }

    orOperation, err = this.tryParseOrOperationAt(index)
    if err != nil {
        return nil, err
    } else if orOperation != nil {
        node.node = orOperation
        goto foundNode
    }

    colonOperation, err = this.tryParseColonOperationAt(index)
    if err != nil {
        return nil, err
    } else if colonOperation != nil {
        node.node = colonOperation
        goto foundNode
    }

    dotOperation, err = this.tryParseDotOperationAt(index)
    if err != nil {
        return nil, err
    } else if dotOperation != nil {
        node.node = dotOperation
        goto foundNode
    }

    arrayAccessOperation, err = this.tryParseArrayAccessOperationAt(index)
    if err != nil {
        return nil, err
    } else if arrayAccessOperation != nil {
        node.node = arrayAccessOperation
        goto foundNode
    }

    return nil, nil

foundNode:
    if this.compressAST {
        return node.node, nil
    } else {
        return node, nil
    }
}

// "!" Expression
func (this *parser) tryParseNotOperationAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Exclamation {
        return nil, nil
    }
    index++

    expression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expression == nil {
        return nil, nil
    }

    return wrappedNode{
        kind:                NodeKind_NotOperation,
        node:                expression,
        extraTokensConsumed: 1,
    }, nil
}

// "-" Expression
func (this *parser) tryParseUnaryMinusOperationAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Minus {
        return nil, nil
    }
    index++

    expression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expression == nil {
        return nil, nil
    }

    if expression.Kind() == NodeKind_Int || expression.Kind() == NodeKind_Float {
        return basicNode{
            kind:           expression.Kind(),
            data:           "-" + expression.(basicNode).data,
            tokensConsumed: expression.TokensConsumed() + 1,
            lineNumber:     expression.LineNumber(),
        }, nil
    }

    return wrappedNode{
        kind:                NodeKind_UnaryMinusOperation,
        node:                expression,
        extraTokensConsumed: 1,
    }, nil
}

// "(" Expression ")"
func (this *parser) tryParseParenthesisOperationAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    expression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expression == nil {
        return nil, nil
    }
    index += expression.TokensConsumed()

    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    return wrappedNode{
        kind:                NodeKind_ParenthesisOperation,
        node:                expression,
        extraTokensConsumed: 2,
    }, nil
}

// Expression "+" Expression
func (this *parser) tryParsePlusOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_PlusOperation, TokenKind_Plus)
}

// Expression "-" Expression
func (this *parser) tryParseMinusOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_MinusOperation, TokenKind_Minus)
}

// Expression "/" Expression
func (this *parser) tryParseDivideOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_DivideOperation, TokenKind_ForwardSlash)
}

// Expression "*" Expression
func (this *parser) tryParseMultiplyOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_MultiplyOperation, TokenKind_Asterisk)
}

// Expression "=" Expression
func (this *parser) tryParseAssignmentOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationWithLineBreaksAt(index, NodeKind_AssignmentOperation, TokenKind_Equals)
}

// Expression "==" Expression
func (this *parser) tryParseEqualityOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_EqualityOperation, TokenKind_Equals, TokenKind_Equals)
}

// Expression "!=" Expression
func (this *parser) tryParseInequalityOperationAt(index uint) (Node, error) {
    inequalityOperation, err := this._tryParseBinaryOperationAt(index, NodeKind_InequalityOperation, TokenKind_Exclamation, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if inequalityOperation == nil {
        return nil, nil
    }
    return simplifyInequalityOperation(inequalityOperation)
}

// Expression "+=" Expression
func (this *parser) tryParsePlusEqualOperationAt(index uint) (Node, error) {
    plusEqualOperation, err := this._tryParseBinaryOperationAt(index, NodeKind_PlusEqualOperation, TokenKind_Plus, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if plusEqualOperation == nil {
        return nil, nil
    }
    return simplifyInPlaceOperation(plusEqualOperation, NodeKind_PlusOperation)
}

// Expression "-=" Expression
func (this *parser) tryParseMinusEqualOperationAt(index uint) (Node, error) {
    minusEqualOperation, err := this._tryParseBinaryOperationAt(index, NodeKind_MinusEqualOperation, TokenKind_Minus, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if minusEqualOperation == nil {
        return nil, nil
    }
    return simplifyInPlaceOperation(minusEqualOperation, NodeKind_MinusOperation)
}

// Expression "/=" Expression
func (this *parser) tryParseDivideEqualOperationAt(index uint) (Node, error) {
    divideEqualOperation, err := this._tryParseBinaryOperationAt(index, NodeKind_DivideEqualOperation, TokenKind_ForwardSlash, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if divideEqualOperation == nil {
        return nil, nil
    }
    return simplifyInPlaceOperation(divideEqualOperation, NodeKind_DivideOperation)
}

// Expression "*=" Expression
func (this *parser) tryParseMultiplyEqualOperationAt(index uint) (Node, error) {
    multiplyEqualOperation, err := this._tryParseBinaryOperationAt(index, NodeKind_MultiplyEqualOperation, TokenKind_Asterisk, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if multiplyEqualOperation == nil {
        return nil, nil
    }
    return simplifyInPlaceOperation(multiplyEqualOperation, NodeKind_MultiplyOperation)
}

// Expression ">" Expression
func (this *parser) tryParseGreaterThanOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_GreaterThanOperation, TokenKind_RightAngleBracket)
}

// Expression "<" Expression
func (this *parser) tryParseLessThanOperationAt(index uint) (Node, error) {
    lessThan, err := this._tryParseBinaryOperationAt(index, NodeKind_LessThanOperation, TokenKind_LeftAngleBracket)
    if err != nil {
        return nil, err
    } else if lessThan == nil {
        return nil, nil
    }

    // Lookahead to avoid incorrectly parsing `MyFunc <x>` as `MyFunc < x`
    index += lessThan.TokensConsumed()
    if !this.isOutOfRangeAt(index) && this.tokens[index].Kind() == TokenKind_RightAngleBracket {
        return nil, nil
    }

    return lessThan, err
}

// Expression "<=" Expression
func (this *parser) tryParseLessThanEqualOperationAt(index uint) (Node, error) {
    lessThanEqual, err := this._tryParseBinaryOperationAt(index, NodeKind_LessThanEqualOperation, TokenKind_LeftAngleBracket, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if lessThanEqual == nil {
        return nil, nil
    }
    return simplifyLessThanEqualOperation(lessThanEqual)
}

// Expression ">=" Expression
func (this *parser) tryParseGreaterThanEqualOperationAt(index uint) (Node, error) {
    greaterThanEqual, err := this._tryParseBinaryOperationAt(index, NodeKind_GreaterThanEqualOperation, TokenKind_RightAngleBracket, TokenKind_Equals)
    if err != nil {
        return nil, err
    } else if greaterThanEqual == nil {
        return nil, nil
    }
    return simplifyGreaterThanEqualOperation(greaterThanEqual)
}

// Expression "and" Expression
func (this *parser) tryParseAndOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_AndOperation, TokenKind_And)
}

// Expression "or" Expression
func (this *parser) tryParseOrOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_OrOperation, TokenKind_Or)
}

// Expression ":" Expression
func (this *parser) tryParseColonOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_ColonOperation, TokenKind_Colon)
}

// Expression "." Expression
func (this *parser) tryParseDotOperationAt(index uint) (Node, error) {
    return this._tryParseBinaryOperationAt(index, NodeKind_DotOperation, TokenKind_Dot)
}

// Expression "[" Expression "]"
func (this *parser) tryParseArrayAccessOperationAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    leftExpression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if leftExpression == nil {
        return nil, nil
    }
    index += leftExpression.TokensConsumed()

    if this.tokens[index].Kind() != TokenKind_LeftSquareBracket {
        return nil, nil
    }
    index++

    rightExpression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if rightExpression == nil {
        return nil, nil
    }
    index++

    if this.tokens[index].Kind() != TokenKind_RightSquareBracket {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_ArrayAccessOperation,
        nodes:               []Node{leftExpression, rightExpression},
        extraTokensConsumed: 2,
    }, nil
}

// "<...>" | Array | Struct | Int | Float | String | "<" QbKey ">" | QbKey | Pair | Vector | Script | Bytes | Random // Needs RandomRange?
func (this *parser) tryParseSubExpressionAt(index uint) (Node, error) {
    cachedNode, found := this.subExpressionCache[index]
    if found {
        return cachedNode, nil
    }
    node, err := this._tryParseSubExpressionAt(index)
    if err != nil {
        return nil, err
    }
    this.subExpressionCache[index] = node
    return node, nil
}

func (this *parser) _tryParseSubExpressionAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    var allArguments, array, struct_, int_, float_, string_, localQbKey, qbKey, pair, vector, script, bytes, random Node
    var err error

    node := wrappedNode{
        kind:                NodeKind_SubExpression,
        node:                nil,
        extraTokensConsumed: 0,
    }

    allArguments, err = this.tryParseAllArgumentsAt(index)
    if err != nil {
        return nil, err
    } else if allArguments != nil {
        node.node = allArguments
        goto foundNode
    }

    array, err = this.tryParseArrayAt(index)
    if err != nil {
        return nil, err
    } else if array != nil {
        node.node = array
        goto foundNode
    }

    struct_, err = this.tryParseStructAt(index)
    if err != nil {
        return nil, err
    } else if struct_ != nil {
        node.node = struct_
        goto foundNode
    }

    float_, err = this.tryParseFloatAt(index)
    if err != nil {
        return nil, err
    } else if float_ != nil {
        node.node = float_
        goto foundNode
    }

    int_, err = this.tryParseIntAt(index)
    if err != nil {
        return nil, err
    } else if int_ != nil {
        node.node = int_
        goto foundNode
    }

    string_, err = this.tryParseStringAt(index)
    if err != nil {
        return nil, err
    } else if string_ != nil {
        node.node = string_
        goto foundNode
    }

    localQbKey, err = this.tryParseLocalQbKeyAt(index)
    if err != nil {
        return nil, err
    } else if localQbKey != nil {
        node.node = localQbKey
        goto foundNode
    }

    qbKey, err = this.tryParseQbKeyAt(index)
    if err != nil {
        return nil, err
    } else if qbKey != nil {
        node.node = qbKey
        goto foundNode
    }

    pair, err = this.tryParsePairAt(index)
    if err != nil {
        return nil, err
    } else if pair != nil {
        node.node = pair
        goto foundNode
    }

    vector, err = this.tryParseVectorAt(index)
    if err != nil {
        return nil, err
    } else if vector != nil {
        node.node = vector
        goto foundNode
    }

    script, err = this.tryParseScriptAt(index)
    if err != nil {
        return nil, err
    } else if script != nil {
        node.node = script
        goto foundNode
    }

    bytes, err = this.tryParseBytesAt(index)
    if err != nil {
        return nil, err
    } else if bytes != nil {
        node.node = bytes
        goto foundNode
    }

    random, err = this.tryParseRandomAt(index)
    if err != nil {
        return nil, err
    } else if random != nil {
        node.node = random
        goto foundNode
    }

    return nil, nil

foundNode:
    if this.compressAST {
        return node.node, nil
    } else {
        return node, nil
    }
}

// "<...>"
func (this *parser) tryParseAllArgumentsAt(index uint) (Node, error) {
    lastTokenIndex := index + 4
    if this.isOutOfRangeAt(lastTokenIndex) {
        return nil, nil
    }

    if this.tokens[index].Kind() == TokenKind_LeftAngleBracket &&
        this.tokens[index+1].Kind() == TokenKind_Dot &&
        this.tokens[index+2].Kind() == TokenKind_Dot &&
        this.tokens[index+3].Kind() == TokenKind_Dot &&
        this.tokens[index+4].Kind() == TokenKind_RightAngleBracket {
        return basicNode{
            kind:           NodeKind_AllArguments,
            data:           "<...>",
            tokensConsumed: 5,
            lineNumber:     this.tokens[index].LineNumber(),
        }, nil
    }

    return nil, nil
}

// "[" (Expression | "," | LineBreak)* "]"
func (this *parser) tryParseArrayAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftSquareBracket {
        return nil, nil
    }
    index++

    var nodes nodeArray
    nodes.nodes = []Node{}
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }

        expression, err := this.tryParseExpressionAt(index)
        if err != nil {
            return nil, err
        } else if expression != nil {
            nodes.save(expression)
            index += expression.TokensConsumed()
            continue
        }

        comma, err := this.tryParseCommaAt(index)
        if err != nil {
            return nil, err
        } else if comma != nil {
            nodes.save(comma)
            index += comma.TokensConsumed()
            continue
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            nodes.save(lineBreak)
            index += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightSquareBracket {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Array,
        nodes:               nodes.nodes,
        extraTokensConsumed: 2,
    }, nil
}

// "{" (Expression | "," | LineBreak)* "}"
func (this *parser) tryParseStructAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++

    var nodes nodeArray
    nodes.nodes = []Node{}
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }

        expression, err := this.tryParseExpressionAt(index)
        if err != nil {
            return nil, err
        } else if expression != nil {
            nodes.save(expression)
            index += expression.TokensConsumed()
            continue
        }

        comma, err := this.tryParseCommaAt(index)
        if err != nil {
            return nil, err
        } else if comma != nil {
            nodes.save(comma)
            index += comma.TokensConsumed()
            continue
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            nodes.save(lineBreak)
            index += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Struct,
        nodes:               nodes.nodes,
        extraTokensConsumed: 2,
    }, nil
}

// ","
func (this *parser) tryParseCommaAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_Comma {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_Comma,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// Float
func (this *parser) tryParseFloatAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_Float {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_Float,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// Int
func (this *parser) tryParseIntAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_Int {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_Int,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// String
func (this *parser) tryParseStringAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_String {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_String,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// "<" QbKey ">"               -
func (this *parser) tryParseLocalQbKeyAt(index uint) (Node, error) {
    lastIndex := index + 2
    if this.isOutOfRangeAt(lastIndex) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftAngleBracket {
        return nil, nil
    }
    index++

    qbKey, err := this.tryParseQbKeyAt(index)
    if err != nil {
        return nil, err
    } else if qbKey == nil {
        return nil, nil
    }
    index += qbKey.TokensConsumed()

    if this.tokens[index].Kind() != TokenKind_RightAngleBracket {
        return nil, nil
    }
    index++

    return wrappedNode{
        kind:                NodeKind_LocalQbKey,
        node:                qbKey,
        extraTokensConsumed: 2,
    }, nil
}

// QbKey
func (this *parser) tryParseQbKeyAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_Identifier &&
        token.Kind() != TokenKind_RawQbKey {
        return nil, nil
    }

    if token.Kind() == TokenKind_RawQbKey {
        var key []byte

        for i := 0; i < 8; i += 2 {
            s := token.Data()[i : i+2]
            int_, err := strconv.ParseInt(s, 16, 32)
            if err != nil {
                return nil, err
            }
            key = append(key, byte(int_))
        }

        return rawQbKeyNode{
            key:        key,
            lineNumber: token.LineNumber(),
        }, nil
    }

    return basicNode{
        kind:           NodeKind_QbKey,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// "(" Expression "," Expression ")"
func (this *parser) tryParsePairAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    expressionOne, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expressionOne == nil {
        return nil, nil
    }
    index += expressionOne.TokensConsumed()

    comma, err := this.tryParseCommaAt(index)
    if err != nil {
        return nil, err
    } else if comma == nil {
        return nil, nil
    }
    index += comma.TokensConsumed()

    expressionTwo, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expressionTwo == nil {
        return nil, nil
    }
    index += expressionTwo.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Pair,
        nodes:               []Node{expressionOne, expressionTwo},
        extraTokensConsumed: 3,
    }, nil
}

// "(" Expression "," Expression "," Expression ")"
func (this *parser) tryParseVectorAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    expressionOne, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expressionOne == nil {
        return nil, nil
    }
    index += expressionOne.TokensConsumed()

    commaOne, err := this.tryParseCommaAt(index)
    if err != nil {
        return nil, err
    } else if commaOne == nil {
        return nil, nil
    }
    index += commaOne.TokensConsumed()

    expressionTwo, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expressionTwo == nil {
        return nil, nil
    }
    index += expressionTwo.TokensConsumed()

    commaTwo, err := this.tryParseCommaAt(index)
    if err != nil {
        return nil, err
    } else if commaTwo == nil {
        return nil, nil
    }
    index += commaTwo.TokensConsumed()

    expressionThree, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if expressionThree == nil {
        return nil, nil
    }
    index += expressionThree.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Vector,
        nodes:               []Node{expressionOne, expressionTwo, expressionThree},
        extraTokensConsumed: 4,
    }, nil
}

// "{" Expression* "}"
func (this *parser) tryParseScriptHeaderAt(index uint) (Node, error) {

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    var expressions nodeArray
    expressions.nodes = []Node{}
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }
        expression, err := this.tryParseExpressionAt(index)
        if err != nil {
            return nil, err
        } else if expression != nil {
            expressions.save(expression)
            index += expression.TokensConsumed()
            continue
        }
        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    // lookahead past newlines to find next { for script body
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            index++
            continue
        }

        if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
            return nil, nil
        }

        break
    }

    return wrappedNodes{
        kind:                NodeKind_ScriptHeader,
        nodes:               expressions.nodes,
        extraTokensConsumed: 2,
    }, nil
}

// "script" QbKey ScriptHeader? LineBreak* "{" ChunkOfCode "}"
func (this *parser) tryParseScriptAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    extraTokensConsumed := uint(0)

    if this.tokens[index].Kind() != TokenKind_Script {
        return nil, nil
    }
    index++
    extraTokensConsumed++

    qbKey, err := this.tryParseQbKeyAt(index)
    if err != nil {
        return nil, err
    } else if qbKey == nil {
        return nil, nil
    }
    index += qbKey.TokensConsumed()

    scriptHeader, err := this.tryParseScriptHeaderAt(index)
    if err != nil {
        return nil, err
    } else if scriptHeader != nil {
        index += scriptHeader.TokensConsumed()
        extraTokensConsumed += 2
    }

    var lineBreaks nodeArray
    lineBreaks.nodes = []Node{}
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }
        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            lineBreaks.save(lineBreak)
            index += lineBreak.TokensConsumed()
            continue
        }
        break
    }
    extraTokensConsumed += lineBreaks.tokensConsumed

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++
    extraTokensConsumed++

    bodyChunk, err := this.tryParseChunkOfCodeAt(index)
    if err != nil {
        return nil, err
    }
    index += bodyChunk.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++
    extraTokensConsumed++

    headerNodes := []Node{}
    if scriptHeader != nil {
        headerNodes = scriptHeader.(wrappedNodes).nodes
    }

    return manyWrappedNodes{
        kind: NodeKind_Script,
        nodeLists: [][]Node{
            {qbKey},
            headerNodes,
            notNilNodes(bodyChunk.(wrappedNodes).nodes),
        },
        extraTokensConsumed: extraTokensConsumed,
    }, nil
}

// "bytes" "(" (Byte | LineBreak)* ")"
func (this *parser) tryParseBytesAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Bytes {
        return nil, nil
    }
    index++

    extraTokensConsumed := uint(0)

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    var bytes_ nodeArray
    for {
        byte_, err := this.tryParseByteAt(index)
        if err != nil {
            return nil, err
        } else if byte_ != nil {
            bytes_.save(byte_)
            index += byte_.TokensConsumed()
            continue
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            index += lineBreak.TokensConsumed()
            extraTokensConsumed += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Bytes,
        nodes:               notNilNodes(bytes_.nodes),
        extraTokensConsumed: 3 + extraTokensConsumed,
    }, nil
}

// HexDigit{2}
func (this *parser) tryParseByteAt(index uint) (Node, error) {
    startIndex := index
    endIndex := startIndex

    if this.isOutOfRangeAt(endIndex) {
        return nil, nil
    }

    digitAsText := ""
    for {
        if len(digitAsText) == 2 {
            break
        }

        if this.isOutOfRangeAt(endIndex) {
            break
        }

        token := this.tokens[endIndex]

        if token.Kind() == TokenKind_Identifier || token.Kind() == TokenKind_Int {
            digitAsText += token.Data()
            endIndex++
            continue
        }

        break
    }

    if len(digitAsText) != 2 {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_Byte,
        data:           digitAsText,
        tokensConsumed: endIndex - startIndex,
        lineNumber:     this.tokens[startIndex].LineNumber(),
    }, nil
}

// "random" "(" (RandomEntry|LineBreak)* ")"
func (this *parser) tryParseRandomAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Random {
        return nil, nil
    }
    index++

    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, nil
    }
    index++

    extraTokensConsumed := uint(0)
    var randomEntries nodeArray
    for {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }

        randomEntry, err := this.tryParseRandomEntryAt(index)
        if err != nil {
            return nil, err
        } else if randomEntry != nil {
            randomEntries.save(randomEntry)
            index += randomEntry.TokensConsumed()
            continue
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            extraTokensConsumed += lineBreak.TokensConsumed()
            index += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Random,
        nodes:               notNilNodes(randomEntries.nodes),
        extraTokensConsumed: 3 + extraTokensConsumed,
    }, nil
}

// "{" ChunkOfCode "}"
func (this *parser) tryParseRandomEntryAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++

    bodyChunk, err := this.tryParseChunkOfCodeAt(index)
    if err != nil {
        return nil, err
    }
    index += bodyChunk.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_RandomEntry,
        nodes:               notNilNodes(bodyChunk.(wrappedNodes).nodes),
        extraTokensConsumed: 0,
    }, nil
}

//  If ElseIf* Else?
func (this *parser) tryParseIfStatementAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if_, err := this.tryParseIfAt(index)
    if err != nil {
        return nil, err
    } else if if_ == nil {
        return nil, nil
    }
    index += if_.TokensConsumed()

    var elseIfs nodeArray
    for {
        elseIf, err := this.tryParseElseIfAt(index)
        if err != nil {
            return nil, err
        } else if elseIf != nil {
            elseIfs.save(elseIf)
            index += elseIf.TokensConsumed()
            continue
        }

        break
    }

    else_, err := this.tryParseElseAt(index)
    if err != nil {
        return nil, err
    }

    ifStatement, err := simplifyElseIfChain(if_, notNilNodes(elseIfs.nodes), else_)
    if err != nil {
        return nil, err
    }

    return ifStatement, nil
}

// "if" "(" Expression* ")" "{" ChunkOfCode "}"
func (this *parser) tryParseIfAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_If {
        return nil, nil
    }
    index++

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftParenthesis {
        return nil, errors.New("if (condition)\n   ^\n   \\____ condition missing")
    }
    index++

    var conditions nodeArray
    for {
        if this.isOutOfRangeAt(index) {
            break
        }

        expression, err := this.tryParseExpressionAt(index)
        if err != nil {
            return nil, err
        } else if expression != nil {
            conditions.save(expression)
            index += expression.TokensConsumed()
            continue
        }

        break
    }

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightParenthesis {
        return nil, nil
    }
    index++

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++

    bodyChunk, err := this.tryParseChunkOfCodeAt(index)
    if err != nil {
        return nil, err
    }
    index += bodyChunk.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++

    return manyWrappedNodes{
        kind: NodeKind_If,
        nodeLists: [][]Node{
            notNilNodes(conditions.nodes),
            notNilNodes(bodyChunk.(wrappedNodes).nodes),
        },
        extraTokensConsumed: 5,
    }, nil
}

// "else" If
func (this *parser) tryParseElseIfAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() != TokenKind_Else {
        return nil, nil
    }
    index++

    if_, err := this.tryParseIfAt(index)
    if err != nil {
        return nil, err
    } else if if_ == nil {
        return nil, nil
    }
    index += if_.TokensConsumed()

    return wrappedNode{
        kind:                NodeKind_ElseIf,
        node:                if_,
        extraTokensConsumed: 1,
    }, nil
}

// "else" "{" ChunkOfCode "}"
func (this *parser) tryParseElseAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Else {
        return nil, nil
    }
    index++

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++

    bodyChunk, err := this.tryParseChunkOfCodeAt(index)
    if err != nil {
        return nil, err
    }
    index += bodyChunk.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++

    return wrappedNodes{
        kind:                NodeKind_Else,
        nodes:               bodyChunk.(wrappedNodes).nodes,
        extraTokensConsumed: 3,
    }, nil
}

// "loop" "{" ChunkOfCode "}" Expression?
func (this *parser) tryParseLoopAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Loop {
        return nil, nil
    }
    index++

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_LeftCurlyBrace {
        return nil, nil
    }
    index++

    bodyChunk, err := this.tryParseChunkOfCodeAt(index)
    if err != nil {
        return nil, err
    }
    index += bodyChunk.TokensConsumed()

    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_RightCurlyBrace {
        return nil, nil
    }
    index++

    expression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    }

    return manyWrappedNodes{
        kind: NodeKind_Loop,
        nodeLists: [][]Node{
            bodyChunk.(wrappedNodes).nodes,
            {expression},
        },
        extraTokensConsumed: 3,
    }, nil
}

// "switch" Expression "{" (Case ChunkOfCode)* "}"
func (this *parser) tryParseSwitchAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }
    if this.tokens[index].Kind() == TokenKind_Switch {
        return nil, errors.New("`switch` syntax not implemented yet")
    }
    return nil, nil
}

// "break"
func (this *parser) tryParseBreakAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    token := this.tokens[index]

    if token.Kind() != TokenKind_Break {
        return nil, nil
    }

    return basicNode{
        kind:           NodeKind_Break,
        data:           token.Data(),
        tokensConsumed: 1,
        lineNumber:     token.LineNumber(),
    }, nil
}

// "return" Expression*
func (this *parser) tryParseReturnAt(index uint) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    if this.tokens[index].Kind() != TokenKind_Return {
        return nil, nil
    }
    index++

    var expressions nodeArray
    expressions.nodes = []Node{}
    for {
        if this.isOutOfRangeAt(index) {
            break
        }

        expression, err := this.tryParseExpressionAt(index)
        if err != nil {
            return nil, err
        } else if expression != nil {
            expressions.save(expression)
            index += expression.TokensConsumed()
            continue
        }

        break
    }

    return wrappedNodes{
        kind:                NodeKind_Return,
        nodes:               expressions.nodes,
        extraTokensConsumed: 1,
    }, nil
}

// SubExpression LineBreak* operatorTokens LineBreak* Expression?
func (this *parser) _tryParseBinaryOperationWithLineBreaksAt(index uint, nodeKind NodeKind, operatorTokens ...TokenKind) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    leftSubExpression, err := this.tryParseSubExpressionAt(index)
    if err != nil {
        return nil, err
    } else if leftSubExpression == nil {
        return nil, nil
    }
    index += leftSubExpression.TokensConsumed()

    var lineBreaks1 nodeArray
    for {
        if this.isOutOfRangeAt(index) {
            break
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            lineBreaks1.save(lineBreak)
            index += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    for _, operatorToken := range operatorTokens {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }
        if this.tokens[index].Kind() != operatorToken {
            return nil, nil
        }
        index++
    }

    var lineBreaks2 nodeArray
    for {
        if this.isOutOfRangeAt(index) {
            break
        }

        lineBreak, err := this.tryParseLineBreakAt(index)
        if err != nil {
            return nil, err
        } else if lineBreak != nil {
            lineBreaks2.save(lineBreak)
            index += lineBreak.TokensConsumed()
            continue
        }

        break
    }

    rightExpression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if rightExpression != nil {
        index += rightExpression.TokensConsumed()
    }

    return manyWrappedNodes{
        kind: nodeKind,
        nodeLists: [][]Node{
            {leftSubExpression, rightExpression},
            notNilNodes(lineBreaks1.nodes),
            notNilNodes(lineBreaks2.nodes),
        },
        extraTokensConsumed: uint(len(operatorTokens)),
    }, nil
}

// SubExpression operatorTokens Expression
func (this *parser) _tryParseBinaryOperationAt(index uint, nodeKind NodeKind, operatorTokens ...TokenKind) (Node, error) {
    if this.isOutOfRangeAt(index) {
        return nil, nil
    }

    leftSubExpression, err := this.tryParseSubExpressionAt(index)
    if err != nil {
        return nil, err
    } else if leftSubExpression == nil {
        return nil, nil
    }
    index += leftSubExpression.TokensConsumed()

    for _, operatorToken := range operatorTokens {
        if this.isOutOfRangeAt(index) {
            return nil, nil
        }
        if this.tokens[index].Kind() != operatorToken {
            return nil, nil
        }
        index++
    }

    rightExpression, err := this.tryParseExpressionAt(index)
    if err != nil {
        return nil, err
    } else if rightExpression == nil {
        return nil, nil
    }
    index += rightExpression.TokensConsumed()

    return manyWrappedNodes{
        kind: nodeKind,
        nodeLists: [][]Node{
            {leftSubExpression, rightExpression},
            {},
            {},
        },
        extraTokensConsumed: uint(len(operatorTokens)),
    }, nil
}

func (this *parser) lineNumberAt(index uint) uint {
    return this.tokens[index].LineNumber()
}

func notNilNodes(nodes []Node) []Node {
    if nodes == nil {
        return []Node{}
    }
    return nodes
}
