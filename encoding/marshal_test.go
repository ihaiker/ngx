package encoding

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestMarshalSuite struct {
	suite.Suite
}

func (p TestMarshalSuite) Marshal(v interface{}, options ...Options) string {
	var bs []byte
	var err error

	if len(options) == 0 {
		bs, err = Marshal(v)
	} else {
		bs, err = MarshalWithOptions(v, options[0])
	}
	p.Nil(err)
	p.T().Log("\n---", string(bs))
	return string(bs)
}

func (p *TestMarshalSuite) TestBase() {
	out := p.Marshal(&Test{
		Name: "name",
	})
	p.Equal(`name: name;`, out)
}

func TestMarshal(t *testing.T) {
	suite.Run(t, new(TestMarshalSuite))
}
