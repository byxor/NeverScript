package compiler

type syntaxTree struct {
	Declarations []*declaration `@@*`
}

type declaration struct {
	EndOfLine         string             `  @Semicolon`
	BooleanAssignment *booleanAssignment `| @@`
	IntegerAssignment *integerAssignment `| @@`
	StringAssignment  *stringAssignment  `| @@`
	FunctionDeclaration *functionDeclaration `| @@`
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

type functionDeclaration struct {
	Function string `@FunctionKeyword`
	Name string `@Identifier`
	OpenCurlyBrace string `@OpenCurlyBrace`
	Statements []*declaration `@@*`
	CloseCurlyBrace string `@CloseCurlyBrace`
}

type integer struct {
	Base2  string `  @Integer_Base2`
	Base8  string `| @Integer_Base8`
	Base16 string `| @Integer_Base16`
	Base10 string `| @Integer_Base10`
}
