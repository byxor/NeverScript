package tokens

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Assignment
	LocalReference
	Subtraction
	Addition
	Division
	Multiplication
	EqualityCheck
	LessThanCheck
	LessThanOrEqualCheck
	GreaterThanCheck
	GreaterThanOrEqualCheck
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
