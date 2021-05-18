package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"strings"
)

type Expression struct {
	Pos       lexer.Position
	Directive []Directive `[@@+`
	Function  *Function   `|@@]`
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
	if err = parser.ParseString(str, str, expr); err != nil {
		if utp, match := err.(participle.UnexpectedTokenError); match {
			err = fmt.Errorf("Expected: [%s]\n%s\n%s^\n", utp.Expected,
				str, strings.Repeat(" ", utp.Position().Offset))
		}
	}
	return
}
