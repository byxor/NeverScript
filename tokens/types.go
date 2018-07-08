package tokens

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Assignment
	LocalReference
	AllLocalReferences
	Subtraction
	Addition
	Division
	Multiplication
	Not
	EqualityCheck
	LessThanCheck
	LessThanOrEqualCheck
	GreaterThanCheck
	GreaterThanOrEqualCheck
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
	Break
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
	ChecksumTableEntry
	NamespaceAccess
	Invalid
)

//go:generate stringer -type=Token
