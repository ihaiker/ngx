package query

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

type Expression struct {
	Directive []Directive `("."@@)+`
}

type Directive struct {
	Name  string `(@Ident|("[" @String "]"))`
	Index *Index `["["@@"]"]`
}

type Index struct {
	Start *int `[@Number]`
	End   *int `[":"][@Number]`
}

func Lexer(str string) (expr *Expression, err error) {
	var def *stateful.Definition
	def, err = stateful.NewSimple([]stateful.Rule{
		{"String", `("(\\"|[^"])*")|('(\\'|[^'])*')`, nil},
		{"Number", `[-]?(\d*\.)?\d+`, nil},
		{"Ident", `[a-zA-Z_]\w*`, nil},
		{"Whitespace", `\s+`, nil},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
	})
	options := []participle.Option{
		participle.Lexer(def),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	}
	expr = &Expression{}
	var parser *participle.Parser
	if parser, err = participle.Build(expr, options...); err != nil {
		return
	}
	//fmt.Println(parser.String())
	err = parser.ParseString(str, str, expr)
	return
}
