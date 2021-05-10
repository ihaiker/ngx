package config

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type charIteratorSuite struct {
	suite.Suite
	it *charIterator
}

func (p *charIteratorSuite) SetupTest() {
	p.it = newCharIteratorWithBytes([]byte(`
# user nginx;
user nobody;
`))
}

func (p *charIteratorSuite) nextAssertIt(assertC string, assertLine int, assertHas bool) {
	c, line, has := p.it.next()
	p.Suite.Equal(has, assertHas)
	p.Suite.Equal(c, assertC)
	p.Suite.Equal(line, assertLine)
}

func (p *charIteratorSuite) TestNext() {
	p.nextAssertIt("\n", 1, true)

	p.nextAssertIt("#", 2, true)
	p.nextAssertIt(" ", 2, true)
	p.nextAssertIt("u", 2, true)
	p.nextAssertIt("s", 2, true)
	p.nextAssertIt("e", 2, true)
	p.nextAssertIt("r", 2, true)
	p.nextAssertIt(" ", 2, true)
	p.nextAssertIt("n", 2, true)
	p.nextAssertIt("g", 2, true)
	p.nextAssertIt("i", 2, true)
	p.nextAssertIt("n", 2, true)
	p.nextAssertIt("x", 2, true)
	p.nextAssertIt(";", 2, true)
	p.nextAssertIt("\n", 2, true)

	p.nextAssertIt("u", 3, true)
	p.nextAssertIt("s", 3, true)
}

func (p *charIteratorSuite) TestFilter() {
	c, line, has := p.it.nextFilter(ValidChars)
	p.Suite.Equal(has, true)
	p.Suite.Equal(line, 2)
	p.Suite.Equal(c, "#")

	c, line, has = p.it.nextFilter(ValidChars)
	p.Suite.Equal(has, true)
	p.Suite.Equal(line, 2)
	p.Suite.Equal(c, "u")

	c, line, has = p.it.nextTo(Not(ValidChars), false)
	p.Suite.Equal(has, true)
	p.Suite.Equal(line, 2)
	p.Suite.Equal(c, "ser")

	p.nextAssertIt(" ", 2, true)

	c, line, has = p.it.nextTo(In(";"), false)
	p.Suite.Equal(has, true)
	p.Suite.Equal(line, 2)
	p.Suite.Equal(c, "nginx")

	p.nextAssertIt(";", 2, true)
}

func TestCharIterator(t *testing.T) {
	suite.Run(t, new(charIteratorSuite))
}
