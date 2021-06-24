package hooks

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type _demoSub struct {
	Sub string
}
type _demo struct {
	Name string
	Desc string
	Sub  *_demoSub
	Sub2 _demoSub
}

type TestVariablesSuite struct {
	suite.Suite
	vars *Variables
}

func (p *TestVariablesSuite) SetupTest() {
	p.vars = NewVariables()
}

func (p TestVariablesSuite) TestMap() {
	p.vars.Parameter("args", map[string]interface{}{
		"name": "ngx",
		"desc": "similar to nginx configuration",
		"demo": &_demo{
			Name: "demo",
			Desc: "",
		},
	})
	val, err := p.vars.Get(".args.name")
	p.Nil(err)
	p.Equal("ngx", val)

	val, err = p.vars.Get(".args.demo.Name")
	p.Nil(err)
	p.Equal("demo", val)
}

func (p TestVariablesSuite) TestStruct() {
	p.vars.Parameter("demo", &_demo{
		Name: "demo",
		Desc: "demo_desc",
	})

	val, err := p.vars.Get(".demo.Name")
	p.Nil(err)
	p.Equal("demo", val)

	val, err = p.vars.Get(".demo.Desc")
	p.Nil(err)
	p.Equal("demo_desc", val)

	val, err = p.vars.Get(".demo.Sub")
	p.T().Log(val, err)

	val, err = p.vars.Get(".demo.Sub2")
	p.T().Log(val, err)

	val, err = p.vars.Get(".demo.Desc1")
	p.T().Log(val, err)
}

func (p *TestVariablesSuite) TestGet() {
	p.vars.Parameter("val", "val")
	val, err := p.vars.Get(".val")
	p.Nil(err)
	p.Equal("val", val)
}

func TestVariables(t *testing.T) {
	suite.Run(t, new(TestVariablesSuite))
}
