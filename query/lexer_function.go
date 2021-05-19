package query

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// args('')
type Function struct {
	Pos  lexer.Position
	Name string     `@Ident`
	Args []*FuncArg `"(" [Whitespace] @@ [( [Whitespace] "," [Whitespace] @@ )+] [Whitespace] ")"`
}

type FuncArg struct {
	Pos       lexer.Position
	Directive *Directive `(@@`
	Index     *int       `|@Number`
	Value     *string    `|@String`
	Function  *Function  `|@@`
	Arrays    []FuncArg  `|("[" @@ (","@@)* "]") )`
}
