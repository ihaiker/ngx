package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/ihaiker/ngx/v2/config"
	"strings"
)

type Expression struct {
	Pos       lexer.Position
	Directive Directives `[@@+`
	Function  *Function  `|@@]`
}

func (expr *Expression) Select(items config.Directives, fns map[string]ExecutorFunction) (config.Directives, error) {
	if expr.Directive != nil || len(expr.Directive) != 0 {
		return expr.Directive.Select(items)
	}
	return expr.Function.Select(items, fns)
}

func Lexer(str string) (expr *Expression, err error) {
	var def *stateful.Definition
	def, err = stateful.NewSimple([]stateful.Rule{
		{"String", `("(\\"|[^"])*")|('(\\'|[^'])*')`, nil},
		{"Regex", `/(\\/|[^/])*/`, nil},
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
