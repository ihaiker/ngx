package query

import "github.com/alecthomas/participle/v2/lexer"

type Directive struct {
	Pos lexer.Position

	Name  string `"."(@Ident|("[" @String "]"))`
	Args  *Args  `[("("@@")")]`
	Index *Index `["["@@"]"]`
}

type Index struct {
	Pos   lexer.Position
	Start *int `[@Number]`
	End   *int `[":"][@Number]`
}

type Args struct {
	Pos      lexer.Position
	Left     *Arg   `@@`
	Operator string `[(([Whitespace] @("&""&") [Whitespace]) | ([Whitespace] @("|""|") [Whitespace]))`
	Right    *Arg   `@@]`
}

type Arg struct {
	Pos        lexer.Position
	Comparison string  `(([@("!" | "@" | "^" | "$" | "&")]`
	Value      *string `@String)`
	Group      *Args   `|("(" @@ ")"))`
}
