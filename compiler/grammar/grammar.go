package grammar

type SyntaxTree struct {
	Declarations []*Declaration `@@*`
}

type Declaration struct {
	EndOfLine         string             `  @Semicolon`
	BooleanAssignment *BooleanAssignment `| @@`
	IntegerAssignment *IntegerAssignment `| @@`
}

type BooleanAssignment struct {
	Name    string `@Identifier`
	Equals  string `@Equals`
	Value   string `@Boolean`
}

type IntegerAssignment struct {
	Name    string   `@Identifier`
	Equals  string   `@Equals`
	Value   *Integer `@@`
}

type Integer struct {
	Base10  string `  @Integer_Base10`
	Base16  string `| @Integer_Base16`
}
