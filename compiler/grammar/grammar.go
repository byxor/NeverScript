package grammar

type Program struct {
	Declarations []*Declaration `@@*`
}

type Declaration struct {
	EndOfLine         *EndOfLine         `  @@`
    BooleanAssignment *BooleanAssignment `| @@`
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
