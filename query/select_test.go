package query

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TestSelectSuite struct {
	suite.Suite
}

func (p *TestSelectSuite) sel(queries ...string) config.Directives {
	p.T().Log("test: ", strings.Join(queries, " | "))
	conf, err := config.Parse(
		"_testdata/nginx.conf")
	p.Nil(err)

	items, err := Selects(conf, queries...)
	p.Nil(err)
	return items
}

func (p *TestSelectSuite) TestSelectBase() {
	items := p.sel(".")
	p.Len(items, 4)
	p.Equal("user", items[0].Name)
	p.Equal("http", items[3].Name)

	items = p.sel(".user")
	p.Len(items, 1)
	p.Equal("nginx", items[0].Args[0])

	items = p.sel(".events.worker_connections")
	p.Len(items, 1)
	p.Equal("1024", items[0].Args[0])

	items = p.sel(".http.server.server_name")
	p.Len(items, 3)
	p.Equal("_", items[0].Args[0])
	p.Equal("aginx.x.do", items[1].Args[0])
	p.Equal("test.renzhen.la", items[2].Args[0])
}

func (p TestSelectSuite) TestSelectRegex() {
	items := p.sel("./(user)|events/")
	p.Len(items, 2)
	p.Equal("user", items[0].Name)
	p.Equal("events", items[1].Name)
}

func (p TestSelectSuite) TestSelectIndex() {
	domains := []string{"_", "aginx.x.do", "test.renzhen.la"}
	for i, domain := range domains {
		items := p.sel(fmt.Sprintf(".http.server[%d].server_name", i))
		p.Len(items, 1)
		p.Equal(domain, items[0].Args[0])
	}

	for i, _ := range domains {
		items := p.sel(fmt.Sprintf(".http.server[%d:].server_name", i))
		p.Len(items, 3-i)
		for j, domain := range domains[i:] {
			p.Equal(domain, items[j].Args[0])
		}
	}

	for i, _ := range domains {
		items := p.sel(fmt.Sprintf(".http.server[:%d].server_name", 3-i))
		p.Len(items, 3-i)
		for j, domain := range domains[:3-i] {
			p.Equal(domain, items[j].Args[0])
		}
	}
}

func (p TestSelectSuite) TestArgs() {
	items := p.sel(".http.server.location('/health').return")
	p.Len(items, 1)
	p.Equal([]string{"200", "ok"}, items[0].Args)
}

func (p TestSelectSuite) TestFN() {
	items := p.sel(".http.server", "select(.server_name,'equal','aginx.x.do')")
	p.Len(items, 1)
}

func TestSelect(t *testing.T) {
	suite.Run(t, new(TestSelectSuite))
}
