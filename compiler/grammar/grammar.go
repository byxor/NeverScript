package grammar

type QbFile struct {
	Statements []*Statement `@@*`
}

type Statement struct {
	EndOfLine  *EndOfLine  `@@`
	Assignment *Assignment `|@@`
}

type EndOfLine struct {
	Character string `@";"`
}

type Assignment struct {
	Name  string      `@Ident "="`
	Value *Expression `@@`
}

type Expression struct {
	Addition    *Addition    `@@`
	Subtraction *Subtraction `|@@`
	Value       *Value       `|@@`
}

type Addition struct {
	Left  *Expression `"(" @@ "+"`
	Right *Expression `@@ ")"`
}

type Subtraction struct {
	Left  *Expression `"(" @@ "-"`
	Right *Expression `@@ ")"`
}

type Value struct {
	Int    string `@Int`
	String string `|@String`
	Float  string `|@Float`
	Name   string `|@Ident`
}
