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
	ChecksumTableEntry
	Invalid
)
