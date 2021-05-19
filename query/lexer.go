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
	Directive []Directive `[@@+`
	Function  *Function   `|@@]`
}

func (expr *Expression) Select(items config.Directives) (config.Directives, error) {
	if expr.Directive != nil || len(expr.Directive) != 0 {
		return expr.directive(items)
	}
	return expr.function(items)
}

//指令检索
func (expr *Expression) directive(items config.Directives) (matched config.Directives, err error) {
	dir := expr.Directive[0]
	if matched, err = dir.Select(items); err != nil {
		return
	}

	if len(expr.Directive[1:]) == 0 {
		return
	}
	expr.Directive = expr.Directive[1:]

	subItems := config.Directives{}
	for _, item := range matched {
		if item.Body != nil {
			subItems = append(subItems, item.Body...)
		}
	}
	return expr.directive(subItems)
}

func (expr *Expression) function(items config.Directives) (config.Directives, error) {
	return items, nil
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
