package query

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type testSelectFunctionSuite struct {
	suite.Suite
}

func (p *testSelectFunctionSuite) sel(queries ...string) config.Directives {
	p.T().Log("test: ", strings.Join(queries, " | "))
	conf, err := config.Parse("_testdata/nginx.conf")
	p.Nil(err)

	items, err := Select(conf, queries...)
	p.Nil(err)
	return items
}

func (p *testSelectFunctionSuite) sels(query string) config.Directives {
	p.T().Log("test: ", query)
	conf, err := config.Parse("_testdata/nginx.conf")
	p.Nil(err)

	items, err := Selects(conf, query)
	p.Nil(err)
	return items
}

func (p *testSelectFunctionSuite) TestSelect() {
	items := p.sel(".http.server", "select(.server_name, 'equal', 'aginx.x.do')", ".location.proxy_pass")
	p.Len(items, 1)
	p.Equal("http://127.0.0.1:8012", items.Index(0).Args[0])

	items = p.sel(".http.upstream", "or(select(., 'equal', 't1'), select(., 'equal', 't2'))")
	p.Len(items, 2)

}

func (p *testSelectFunctionSuite) TestArgs() {
	items := p.sel(".http", "args(.log_format)")
	p.Len(items, 1)
	p.Equal("main", items.Index(0).Args[0])

	items = p.sel(".http.log_format", "args(.)")
	p.Len(items, 1)
	p.Equal("main", items.Index(0).Args[0])

	items = p.sels(".http.log_format | args(.)")
	p.Len(items, 1)
	p.Equal("main", items.Index(0).Args[0])

	items = p.sel(".http.server.location.proxy_pass", "args(.)")
	p.Len(items, 2)
	p.Equal("http://127.0.0.1:8012", items.Index(0).Args[0])
	p.Equal("http://127.0.0.1:8011", items.Index(1).Args[0])
}

func (p testSelectFunctionSuite) TestHas() {
	items := p.sel(`.http.server`, `has(.location, "gzip")`, `.server_name`)
	p.Len(items, 1)
	p.Equal("aginx.x.do", items.Index(0).Args[0])
}

func (p testSelectFunctionSuite) TestLength() {
	items := p.sel(`length(.http.server)`)
	p.Len(items, 1)
	p.Equal("3", items.Index(0).Args[0])
}

func (p testSelectFunctionSuite) TestIndex() {
	args := []string{
		`main`,
		`$remote_addr - $remote_user [$time_local] "$request" `,
		`$status $body_bytes_sent "$http_referer" `,
		`"$http_user_agent" "$http_x_forwarded_for"`,
	}
	for i, arg := range args {
		items := p.sel(".http", fmt.Sprintf("index(args(.log_format), %d)", i))
		p.Len(items, 1)
		p.Equal("key", items.Index(0).Name)
		p.Equal(arg, items.Index(0).Args[0])
	}
	for i := -1; i > -2; i-- {
		items := p.sel(".http", fmt.Sprintf(`index(args(.log_format), %d)`, i))
		p.Len(items, 1)
		p.Equal(args[len(args)+i], items.Index(0).Args[0])
	}
}

func TestSelectFunction(t *testing.T) {
	suite.Run(t, new(testSelectFunctionSuite))
}
