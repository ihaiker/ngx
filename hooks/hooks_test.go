package hooks_test

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/hooks"
	"github.com/ihaiker/ngx/v2/query"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testAfterSuite struct {
	suite.Suite
	conf *config.Configuration
}

func (p *testAfterSuite) SetupTest() {
	var err error
	p.conf, err = config.Parse("./_testdata/nginx.conf")
	p.Nil(err)
}

func (p testAfterSuite) TestIncludeMerge() {
	include := hooks.Include(true, hooks.WalkFiles("./_testdata"))
	err := include(p.conf)
	p.Nil(err)
	items, err := query.Selects(p.conf, ".http.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])
}

func (p testAfterSuite) TestIncludeNotMerge() {
	include := hooks.Include(false, hooks.WalkFiles("./_testdata"))
	err := include(p.conf)
	p.Nil(err)

	items, err := query.Selects(p.conf, ".http.include.file.server.server_name")
	p.Len(items, 2)
	p.Equal("a", items[0].Args[0])
	p.Equal("b", items[1].Args[0])
}

func (p testAfterSuite) TestTemplate() {
	err := hooks.Template(p.conf)
	p.Nil(err)
	p.T().Log(p.conf.Pretty())
}

func TestAfter(t *testing.T) {
	suite.Run(t, new(testAfterSuite))
}
