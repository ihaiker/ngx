package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/query/methods"
	"strings"
)

type expression struct {
	Pos       lexer.Position
	Directive directives `[@@+`
	Function  *function  `|@@]`
}

func (expr *expression) call(items config.Directives, fnm *methods.FunctionManager) (config.Directives, error) {
	if expr.Directive != nil || len(expr.Directive) != 0 {
		return expr.Directive.call(items)
	}
	return expr.Function.callForExpression(items, fnm)
}

func parseLexer(str string) (expr *expression, err error) {
	var def *stateful.Definition
	def, err = stateful.NewSimple([]stateful.Rule{
		{"String", `("(\\"|[^"])*")|('(\\'|[^'])*')`, nil},
		{"Regex", `/(\\/|[^/])*/`, nil},
		{"Number", `[-]?(\d*\.)?\d+`, nil},
		{"Ident", `[a-zA-Z_]\w*`, nil},
		{"Space", `[ ]{1}`, nil},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
		//{"whitespace", `\s+`, nil}, //此处用于控制空格，如果这样写可以允许不用语法之间使用多个空格
	})
	options := []participle.Option{
		participle.Lexer(def),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.UseLookahead(2),
	}

	expr = &expression{}
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
