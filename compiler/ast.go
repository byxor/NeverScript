package compiler

type AstNode struct {
	Kind AstKind
	Data AstData
}

type AstKind int

const (
	AstKind_Root = iota
	AstKind_Assignment
	AstKind_Invocation
	AstKind_Comment
	AstKind_NewLine
	AstKind_Script
	AstKind_WhileLoop
	AstKind_Break
	AstKind_Return
	AstKind_IfStatement
	AstKind_LogicalNot
	AstKind_LogicalAnd
	AstKind_LogicalOr
	AstKind_LocalReference
	AstKind_AllArguments
	AstKind_Checksum
	AstKind_Float
	AstKind_Integer
	AstKind_String
	AstKind_AdditionExpression
	AstKind_SubtractionExpression
	AstKind_MultiplicationExpression
	AstKind_DivisionExpression
	AstKind_GreaterThanExpression
	AstKind_GreaterThanEqualsExpression
	AstKind_LessThanExpression
	AstKind_LessThanEqualsExpression
	AstKind_EqualsExpression
	AstKind_DotExpression
	AstKind_UnaryExpression
	AstKind_Pair
	AstKind_Vector
	AstKind_Struct
	AstKind_Array
	AstKind_Comma
	AstKind_Random
)

func (astKind AstKind) String() string {
	return [...]string{
		"AstKind_Root",
		"AstKind_Assignment",
		"AstKind_Invocation",
		"AstKind_Comment",
		"AstKind_NewLine",
		"AstKind_Script",
		"AstKind_WhileLoop",
		"AstKind_Break",
		"AstKind_Return",
		"AstKind_IfStatement",
		"AstKind_LogicalNot",
		"AstKind_LogicalAnd",
		"AstKind_LogicalOr",
		"AstKind_LocalReference",
		"AstKind_AllArguments",
		"AstKind_Checksum",
		"AstKind_Float",
		"AstKind_Integer",
		"AstKind_String",
		"AstKind_AdditionExpression",
		"AstKind_SubtractionExpression",
		"AstKind_MultiplicationExpression",
		"AstKind_DivisionExpression",
		"AstKind_GreaterThanExpression",
		"AstKind_GreaterThanEqualsExpression",
		"AstKind_LessThanExpression",
		"AstKind_LessThanEqualsExpression",
		"AstKind_EqualsExpression",
		"AstKind_DotExpression",
		"AstKind_UnaryExpression",
		"AstKind_Pair",
		"AstKind_Vector",
		"AstKind_Struct",
		"AstKind_Array",
		"AstKind_Comma",
		"AstKind_Random",
	}[astKind]
}

type AstData interface {
	astData()
}

type AstData_Root struct {
	BodyNodes []AstNode
}
func (astData AstData_Root) astData() {}

type AstData_Assignment struct {
	NameNode  AstNode
	ValueNode AstNode
}
func (astData AstData_Assignment) astData() {}

type AstData_Invocation struct {
	ScriptIdentifierNode              AstNode
	ParameterNodes                    []AstNode
	TokensConsumedByEachParameterNode []int
}
func (astData AstData_Invocation) astData() {}

type AstData_Empty struct{}
func (astData AstData_Empty) astData() {}

type AstData_Script struct {
	NameNode              AstNode
	DefaultParameterNodes []AstNode
	BodyNodes             []AstNode
}
func (astData AstData_Script) astData() {}

type AstData_WhileLoop struct {
	BodyNodes    []AstNode
}
func (astData AstData_WhileLoop) astData() {}

type AstData_IfStatement struct {
	Conditions  []AstNode
	Bodies      [][]AstNode
}
func (astData AstData_IfStatement) astData() {}

type AstData_Comment struct {
	CommentToken Token
}
func (astData AstData_Comment) astData() {}

type AstData_LocalReference struct {
	Node AstNode
}
func (astData AstData_LocalReference) astData() {}

type AstData_Checksum struct {
	IsRawChecksum bool
	ChecksumToken Token

}
func (astData AstData_Checksum) astData() {}

type AstData_Float struct {
	FloatToken Token
}
func (astData AstData_Float) astData() {}

type AstData_Integer struct {
	IntegerToken Token
}
func (astData AstData_Integer) astData() {}

type AstData_String struct {
	StringToken Token
}
func (astData AstData_String) astData() {}

type AstData_BinaryExpression struct {
	LeftNode  AstNode
	RightNode AstNode
}
func (astData AstData_BinaryExpression) astData() {}

type AstData_Pair struct {
	FloatNodeA AstNode
	FloatNodeB AstNode
}
func (astData AstData_Pair) astData() {}

type AstData_Vector struct {
	FloatNodeA AstNode
	FloatNodeB AstNode
	FloatNodeC AstNode
}
func (astData AstData_Vector) astData() {}

type AstData_UnaryExpression struct {
	Node AstNode
}
func (astData AstData_UnaryExpression) astData() {}

type AstData_Struct struct {
	ElementNodes []AstNode
}
func (astData AstData_Struct) astData() {}

type AstData_Array struct {
	ElementNodes []AstNode
}
func (astData AstData_Array) astData() {}

type AstData_Random struct {
	BranchWeights []AstNode
	Branches [][]AstNode
}
func (astData AstData_Random) astData() {}