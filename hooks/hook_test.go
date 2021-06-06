package hooks_test

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/hooks"
	"github.com/ihaiker/ngx/v2/query"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type testHookSuite struct {
	suite.Suite
}

func (p testHookSuite) SetupTest() {
	hooks.Defaults.Hook(&hooks.IncludeHooker{Merge: true,
		Search: hooks.WalkFiles("./_testdata")}, "include")
}

func (p testHookSuite) selects(fileName string, queries string) config.Directives {
	conf, err := config.Parse(fmt.Sprintf("./_testdata/%s.ngx.conf", fileName))
	p.Nil(err)

	err = hooks.Defaults.Execute(conf)
	p.Nil(err)

	items, err := query.Selects(conf, queries)
	p.Nil(err)
	return items
}

func (p testHookSuite) TestParameter() {
	serverName := "v2.aginx.io"
	listen := "80"
	hooks.Defaults.Parameter("serverName", serverName)
	hooks.Defaults.Parameter("listen", listen)

	items := p.selects("params", ".http.server")
	p.Len(items, 1)
	p.Equal(serverName, items[0].Body.Get("server_name").Args[0])
	p.Equal(listen, items[0].Body.Get("listen").Args[0])
	p.Equal(os.Getenv("HOME"), items[0].Body.
		Get("location").Body.Get("root").Args[0])
}

func (p testHookSuite) TestFunc() {
	hooks.Defaults.Func("test_fn", func() string {
		return "test_value"
	})
	items := p.selects("params", ".user")
	p.Len(items, 1)
	p.Equal("test_value", items[0].Args[0])
}

func (p testHookSuite) TestInclude() {
	items := p.selects("include", ".http.include")
	p.Len(items, 0)
	items = p.selects("include", ".http.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])

	//include merge
	hooks.Defaults.Hook(&hooks.IncludeHooker{Merge: false,
		Search: hooks.WalkFiles("./_testdata")}, "include")

	items = p.selects("include", ".http.include")
	p.Len(items, 1)

	items = p.selects("include", ".http.include.file.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])
}

func (p *testHookSuite) TestSwitch() {
	set := func(env, field string, serverName, listen string) {
		_ = os.Setenv("SERVER_TYPE", env)
		conf, err := config.Parse("./_testdata/switch.ngx.conf")
		p.Nil(err)

		hooks.Defaults.Parameter("serverType", field)
		err = hooks.Defaults.Execute(conf)
		p.Nil(err)

		items, err := query.Selects(conf, ".http.server")
		p.Nil(err)
		p.Len(items, 2)

		p.Equal(serverName, items[0].Body.Get("server_name").Args[0])
		p.Equal(listen, items[1].Body.Get("listen").Args[0])
	}

	set("http", "http", "switch_http", "80")
	set("http", "https", "switch_http", "443")
	set("http", "", "switch_http", "8080")

	set("https", "http", "switch_https", "80")
	set("https", "https", "switch_https", "443")
	set("https", "", "switch_https", "8080")

	set("", "http", "switch_8080", "80")
	set("", "https", "switch_8080", "443")
	set("", "", "switch_8080", "8080")
}

func (p *testHookSuite) TestRepeat() {
	items := p.selects("repeat", ".http.server.server_name")
	p.Len(items, 3)
	p.Equal("a0.aginx.io", items[0].Args[0])
	p.Equal("a1.aginx.io", items[1].Args[0])
	p.Equal("a2.aginx.io", items[2].Args[0])

	items = p.selects("repeat", ".http.server.listen")
	p.Len(items, 3)
	p.Equal("80", items[0].Args[0])
	p.Equal("81", items[1].Args[0])
	p.Equal("82", items[2].Args[0])
}

func (p *testHookSuite) TestTemplate() {
	items := p.selects("template", ".http.server.server_name")
	p.Len(items, 2)
}

func TestHooks(t *testing.T) {
	suite.Run(t, new(testHookSuite))
}
