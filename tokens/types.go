package tokens

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Comma
	Dot
	Assignment
	LocalReference
	AllLocalReferences
	String
	LocalString
	Pair
	Vector
	While
	Repeat
	Break
	Subtraction
	Addition
	Division
	Multiplication
	Not
	Or
	And
	EqualityCheck
	LessThanCheck
	LessThanOrEqualCheck
	GreaterThanCheck
	GreaterThanOrEqualCheck
	ExecuteRandomBlock
	StartOfExpression
	EndOfExpression
	StartOfStruct
	EndOfStruct
	StartOfArray
	EndOfArray
	StartOfSwitch
	EndOfSwitch
	SwitchCase
	DefaultSwitchCase
	StartOfFunction
	EndOfFunction
	Return
	StartOfIf
	Else
	ElseIf
	EndOfIf
	OptimisedIf
	OptimisedElse
	Integer
	Float
	Name
	ShortJump
	LongJump
	ChecksumTableEntry
	NamespaceAccess
	Invalid
)

//go:generate stringer -type=Token
