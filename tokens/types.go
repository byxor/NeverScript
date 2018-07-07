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
	StartOfFunction
	EndOfFunction
	Return
	Break
	StartOfIf
	Else
	ElseIf
	EndOfIf
	Integer
	Float
	Name
	WeirdThing
	ChecksumTableEntry
	Invalid
)

//go:generate stringer -type=Token
