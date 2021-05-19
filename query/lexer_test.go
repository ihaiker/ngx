package query

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TestLexerSuite struct {
	suite.Suite
}

func (p TestLexerSuite) lexer(str string) *Expression {
	p.T().Log("test:", str)
	expr, err := Lexer(str)
	p.Require().Nil(err)
	return expr
}

func (p *TestLexerSuite) TestBase() {
	base := func(nql string) {
		expr := p.lexer(nql)
		items := strings.Split(nql[1:], ".")
		for i, item := range items {
			if strings.HasPrefix(item, "[") {
				name := item[2 : len(item)-2]
				p.Equal(name, *expr.Directive[i].Name)
			} else {
				p.Equal(item, *expr.Directive[i].Name)
			}
		}
	}
	base(".http")
	base(".http.server")
	base(".http.server.server_name")
	base(".http.server.listen.set_proxy_header")
	base(".['test_name']")
	base(".http.server.['server_name']")

	expr := p.lexer(".")
	p.Len(expr.Directive, 1)
	p.Nil(expr.Directive[0].Name)
	p.Nil(expr.Directive[0].Args)
	p.Nil(expr.Directive[0].Index)

	expr = p.lexer(".[0]")
	p.Len(expr.Directive, 1)
	p.Nil(expr.Directive[0].Name)
	p.Equal(0, *expr.Directive[0].Index.Start)
}

func (p TestLexerSuite) TestNameOther() {
	expr := p.lexer(".*[0]")
	p.Nil(expr.Directive[0].Name)
	p.Nil(expr.Directive[0].Regex)

	expr = p.lexer("./^abc/")
	p.Nil(expr.Directive[0].Name)
	p.Equal(`/^abc/`, *expr.Directive[0].Regex)

	expr = p.lexer("./.*/")
	p.Nil(expr.Directive[0].Name)
	p.Equal(`/.*/`, *expr.Directive[0].Regex)

}

func (p TestLexerSuite) TestBaseIndex() {
	base := func(nql string) {
		expr := p.lexer(nql)
		items := strings.Split(nql[1:], ".")
		for i, item := range items {
			start := strings.Index(item, "[")
			name := item[:start]
			p.Equal(name, *expr.Directive[i].Name)
			p.Equal(i, *expr.Directive[i].Index.Start)
			if expr.Directive[i].Index.End != nil {
				p.Equal(i+1, *expr.Directive[i].Index.End)
			}
		}
	}
	base(".http[0]")
	base(".http[0].server[1]")
	base(".http[0].server[1].server_name[2:3]")
	base(".http[0:1].server[1].listen[2:3].set_proxy_header[3:4]")

	nql := ".['http'].server[1].location.set_proxy_header[2:3]"
	p.T().Log(nql)
	expr := p.lexer(nql)
	p.Len(expr.Directive, 4)
	p.Equal("http", *expr.Directive[0].Name)
	p.Equal("server", *expr.Directive[1].Name)
	p.Equal("location", *expr.Directive[2].Name)
	p.Equal("set_proxy_header", *expr.Directive[3].Name)
	p.Equal(1, *expr.Directive[1].Index.Start)
	p.Equal(2, *expr.Directive[3].Index.Start)
	p.Equal(3, *expr.Directive[3].Index.End)
}

func (p TestLexerSuite) TestArgs() {
	expr := p.lexer(".http('name')")
	p.Equal("name", *expr.Directive[0].Args.Left.Value)

	expr = p.lexer(".http.server.server_name('api.aginx.io')")
	p.Equal("api.aginx.io", *expr.Directive[2].Args.Left.Value)

	for _, opt := range []string{"", "!", "@", "^", "$", "&"} {
		expr = p.lexer(fmt.Sprintf(".server_name(%s'api.aginx.io')", opt))
		p.Equal(opt, expr.Directive[0].Args.Left.Comparison)
	}

	expr = p.lexer(".http.server[1].server_name('api.aginx.io')")
	p.Equal("api.aginx.io", *expr.Directive[2].Args.Left.Value)
	p.Equal(1, *expr.Directive[1].Index.Start)

	expr = p.lexer(".http.server[1].server_name('ssl' && 'http')")
	p.Equal("ssl", *expr.Directive[2].Args.Left.Value)
	p.Equal("&&", expr.Directive[2].Args.Operator)
	p.Equal("http", *expr.Directive[2].Args.Right.Value)

	expr = p.lexer(".http.server[1].server_name('ssl' && (('api.aginx.io' || 'test') || 'v2.aginx.io'))")
	p.Equal(1, *expr.Directive[1].Index.Start)
	p.Equal("ssl", *expr.Directive[2].Args.Left.Value)
	p.Equal("&&", expr.Directive[2].Args.Operator)
	p.Equal("api.aginx.io", *expr.Directive[2].Args.Right.Group.Left.Group.Left.Value)
	p.Equal("||", expr.Directive[2].Args.Right.Group.Operator)
	p.Equal("v2.aginx.io", *expr.Directive[2].Args.Right.Group.Right.Value)

	expr = p.lexer(`.server_name(('api.aginx.io' || 'v2.aginx.io') && 'ssl')`)
	p.Equal("api.aginx.io", *expr.Directive[0].Args.Left.Group.Left.Value)
	p.Equal("v2.aginx.io", *expr.Directive[0].Args.Left.Group.Right.Value)
	p.Equal("ssl", *expr.Directive[0].Args.Right.Value)
}

func (p TestLexerSuite) TestFunction() {
	expr := p.lexer("args('name', .http)")
	p.Equal("args", expr.Function.Name)
	p.Equal("name", *expr.Function.Args[0].Value)
	p.Equal("http", *expr.Function.Args[1].Directive.Name)

	expr = p.lexer("arg( 1, 'name', 'SELECT' )")
	p.Equal("arg", expr.Function.Name)
	p.Equal(1, *expr.Function.Args[0].Index)
	p.Equal("name", *expr.Function.Args[1].Value)

	expr = p.lexer("arg1( arg2( arg3( 'name' ) ) , 'http' )")
	p.Equal("arg1", expr.Function.Name)
	p.Equal("arg2", expr.Function.Args[0].Function.Name)
	p.Equal("http", *expr.Function.Args[1].Value)
	p.Equal("name", *expr.Function.Args[0].Function.Args[0].Function.Args[0].Value)
}

func TestLexer(t *testing.T) {
	suite.Run(t, new(TestLexerSuite))
}
