package hooks_test

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/hooks"
	"github.com/ihaiker/ngx/v2/query"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type testHookSuite struct {
	suite.Suite
	conf *config.Configuration
}

func (p *testHookSuite) SetupTest() {
	var err error
	p.conf, err = config.Parse("./_testdata/nginx.conf")
	p.Nil(err)
}

func (p testHookSuite) TestIncludeMerge() {
	include := &hooks.IncludeHooker{
		Merge:  true,
		Search: hooks.WalkFiles("./_testdata"),
	}
	err := include.Execute(p.conf)
	p.Nil(err)
	items, err := query.Selects(p.conf, ".http.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])
}

func (p testHookSuite) TestIncludeNotMerge() {
	include := &hooks.IncludeHooker{
		Merge:  false,
		Search: hooks.WalkFiles("./_testdata"),
	}
	err := include.Execute(p.conf)
	p.Nil(err)

	items, err := query.Selects(p.conf, ".http.include.file.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])
}

func (p testHookSuite) TestTemplate() {
	template := &hooks.TemplateHooker{}
	err := template.Execute(p.conf)
	p.Nil(err)

	services, err := query.Selects(p.conf, ".http.server")
	p.Nil(err)
	p.Len(services, 2)

	p.Equal([]string{"arg1", "arg2"}, services[0].Args)
	p.Equal([]string{"arg0"}, services[1].Args)
}

func (p testHookSuite) TestParameterHook() {
	tmpdir := os.Getenv("TMPDIR")
	_ = os.Setenv("WORKER_CONNECTIONS", "24")

	parameters := hooks.Parameter()
	parameters.Add("access_log", tmpdir)

	err := parameters.Execute(p.conf)
	p.Nil(err)

	items, err := query.Selects(p.conf, ".events.worker_connections")
	p.Nil(err)
	p.Len(items, 1)
	p.Equal("24", items[0].Args[0])

	items, err = query.Selects(p.conf, ".http.access_log")
	p.Nil(err)
	p.Len(items, 1)
	p.Equal(tmpdir, items[0].Args[0])
}

func (p *testHookSuite) TestSwitch() {
	set := func(env, field string, serverName, listen string) {
		_ = os.Setenv("SERVER_TYPE", env)
		params := hooks.Parameter()
		params.Add("serverType", field)

		conf := &config.Configuration{
			Body: p.conf.Body.Clone(),
		}
		hooker := &hooks.SwitchHooker{
			Parameters: params,
		}
		err := hooker.Execute(conf)
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

func (p *testHookSuite) TestAll() {
	include := &hooks.IncludeHooker{
		Merge:  true,
		Search: hooks.WalkFiles("./_testdata"),
	}
	_ = os.Setenv("SERVER_TYPE", "https")
	_ = os.Setenv("WORKER_CONNECTIONS", "24")

	parameters := hooks.Parameter()
	parameters.Add("serverType", "http")
	parameters.Add("access_log", os.Getenv("TMPDIR"))

	templateHooker := &hooks.TemplateHooker{}
	switchHooker := &hooks.SwitchHooker{
		Parameters: parameters,
	}
	hookers := hooks.New(parameters, include, templateHooker, switchHooker)

	err := hookers.Execute(p.conf)
	p.Nil(err)

	p.T().Log(p.conf.Pretty())
}

func TestAfter(t *testing.T) {
	suite.Run(t, new(testHookSuite))
}
