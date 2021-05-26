package methods_test

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/query/methods"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type TestManagerSuite struct {
	suite.Suite
	manager *methods.FunctionManager
}

type testItem struct {
	Name     string
	Function interface{}
	Error    string
}

func (p *TestManagerSuite) SetupTest() {
	p.manager = methods.New()
}

func (p TestManagerSuite) TestArgsType() {
	functions := []interface{}{
		func() bool {
			return false
		},
		func(items config.Directives, name string, idx int, has bool) (string, error) {
			return strconv.QuoteRune(rune(name[idx])), nil
		},
		func(name, operator, value interface{}) bool {
			return false
		},
		func(all ...string) bool {
			return false
		},
	}

	var err error
	for _, function := range functions {
		err = p.manager.Add("test_demo", function)
		p.Nil(err)
	}
}

func TestManager(t *testing.T) {
	suite.Run(t, new(TestManagerSuite))
}
