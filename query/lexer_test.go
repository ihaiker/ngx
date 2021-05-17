package query

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TestLexerSuite struct {
	suite.Suite
}

func (p TestLexerSuite) lexer(str string) *Expression {
	expr, err := Lexer(str)
	p.Require().Nil(err)
	return expr
}

func (p *TestLexerSuite) TestBase() {
	base := func(nql string) {
		p.T().Log("test:", nql)
		expr := p.lexer(nql)
		items := strings.Split(nql[1:], ".")
		for i, item := range items {
			if strings.HasPrefix(item, "[") {
				name := item[2 : len(item)-2]
				p.Equal(name, expr.Directive[i].Name)
			} else {
				p.Equal(item, expr.Directive[i].Name)
			}
		}
	}
	base(".http")
	base(".http.server")
	base(".http.server.server_name")
	base(".http.server.listen.set_proxy_header")
	base(".['test_name']")
	base(".http.server.['server_name']")
}

func (p TestLexerSuite) TestBaseIndex() {
	base := func(nql string) {
		p.T().Log("test:", nql)
		expr := p.lexer(nql)
		items := strings.Split(nql[1:], ".")
		for i, item := range items {
			start := strings.Index(item, "[")
			name := item[:start]
			p.Equal(name, expr.Directive[i].Name)
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
	p.Equal("http", expr.Directive[0].Name)
	p.Equal("server", expr.Directive[1].Name)
	p.Equal("location", expr.Directive[2].Name)
	p.Equal("set_proxy_header", expr.Directive[3].Name)
	p.Equal(1, *expr.Directive[1].Index.Start)
	p.Equal(2, *expr.Directive[3].Index.Start)
	p.Equal(3, *expr.Directive[3].Index.End)
}

func TestLexer(t *testing.T) {
	suite.Run(t, new(TestLexerSuite))
}
