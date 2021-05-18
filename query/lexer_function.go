package query

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ihaiker/ngx/v2/config"
)

// args(''`)
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
	Function  *Function  `|@@)`
}

//方法执行
type ExecFunction func(items config.Directives, args ...FuncArg) config.Directives
