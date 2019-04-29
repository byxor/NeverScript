package grammar

type SyntaxTree struct {
	Declarations []*Declaration `@@*`
}

type Declaration struct {
	EndOfLine         *EndOfLine         `  @@`
	BooleanAssignment *BooleanAssignment `| @@`
	IntegerAssignment *IntegerAssignment `| @@`
}

type EndOfLine struct {
	Value *string `";"`
}

type BooleanAssignment struct {
	Name    string   `@Ident`
	Equals  string   `"="`
	Boolean *Boolean `@@`
}

type Boolean struct {
	Value string `@"true"|@"false"`
}

type IntegerAssignment struct {
	Name    string   `@Ident`
	Equals  string   `"="`
	Integer *Integer `@@`
}

type Integer struct {
	Decimal *float64 `@Int`
}
