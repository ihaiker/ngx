package config

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type tokenIteratorSite struct {
	suite.Suite
	it *tokenIterator
}

func (p *tokenIteratorSite) SetupTest() {
	p.it = newTokenIteratorWithBytes([]byte(`
http  {
    include mime.types;
    default_type application/octet-stream;
}
`))
}

func (p *tokenIteratorSite) assertIt(aToken string, aLine int, aHas bool) {
	token, line, has := p.it.next()
	p.Equal(token, aToken)
	p.Equal(line, aLine)
	p.Equal(has, aHas)
}

func (p *tokenIteratorSite) TestNext() {
	p.assertIt("http", 2, true)
	p.assertIt("{", 2, true)

	p.assertIt("include", 3, true)
	p.assertIt("mime.types", 3, true)
	p.assertIt(";", 3, true)

	p.assertIt("default_type", 4, true)
	p.assertIt("application/octet-stream", 4, true)
	p.assertIt(";", 4, true)

	p.assertIt("}", 5, true)

	p.assertIt("", 0, false)
}

func (p *tokenIteratorSite) TestExport() {
	words, token, err := p.it.expectNext(In(";", "{"))
	p.Nil(err)
	p.Equal([]string{"http"}, words)
	p.Equal("{", token)

	words, token, err = p.it.expectNext(In(";", "{"))
	p.Nil(err)
	p.Equal([]string{"include", "mime.types"}, words)
	p.Equal(";", token)

	words, token, err = p.it.expectNext(In(";", "{"))
	p.Nil(err)
	p.Equal([]string{"default_type", "application/octet-stream"}, words)
	p.Equal(";", token)

	p.assertIt("}", 5, true)
	p.assertIt("", 0, false)
}

func TestTokenIterator(t *testing.T) {
	suite.Run(t, new(tokenIteratorSite))
}
