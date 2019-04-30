package compiler

type syntaxTree struct {
	Declarations []*declaration `@@*`
}

type Declaration struct {
	EndOfLine         string             `  @Semicolon`
	BooleanAssignment *booleanAssignment `| @@`
	IntegerAssignment *integerAssignment `| @@`
	StringAssignment  *stringAssignment  `| @@`
}

type booleanAssignment struct {
	Name   string `@Identifier`
	Equals string `@Equals`
	Value  string `@Boolean`
}

type integerAssignment struct {
	Name   string   `@Identifier`
	Equals string   `@Equals`
	Value  *integer `@@`
}

type stringAssignment struct {
	Name   string `@Identifier`
	Equals string `@Equals`
	Value  string `@String`
}

type integer struct {
	Base2  string `  @Integer_Base2`
	Base8  string `| @Integer_Base8`
	Base16 string `| @Integer_Base16`
	Base10 string `| @Integer_Base10`
}
