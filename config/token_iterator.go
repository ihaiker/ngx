package config

import (
	"fmt"
	"io"
)

type tokenIterator struct {
	it *charIterator
}

func newTokenIterator(filename string) (*tokenIterator, error) {
	chatIt, err := newCharIterator(filename)
	if err != nil {
		return nil, err
	}
	tokenIt := &tokenIterator{it: chatIt}
	return tokenIt, nil
}

func newTokenIteratorWithBytes(bs []byte) *tokenIterator {
	chatIt := newCharIteratorWithBytes(bs)
	tokenIt := &tokenIterator{it: chatIt}
	return tokenIt
}

func (self *tokenIterator) next() (token string, tokenLine int, tokenHas bool, err error) {
	for {
		char, line, has := self.it.nextFilter(ValidChars)
		if !has {
			return
		}
		switch char {
		case ";", "{", "}":
			{
				token = char
				tokenLine = line
				tokenHas = true
				return
			}
		case "#":
			{
				word, _, _ := self.it.nextTo(In("\n"), false)
				token = char + word
				tokenLine = line
				tokenHas = true
				return
			}
		case "'", `"`, "`":
			{
				word, _, wordHas := self.it.nextTo(In(char), true)
				if !wordHas {
					err = fmt.Errorf("not found the string end quote %s at line : %d", char, line)
					return
				}
				token = word[:len(word)-1] //remove quota
				tokenLine = line
				tokenHas = true
				return
			}
		default:
			word, _, wordHas := self.it.nextTo(Not(ValidChars).Or(In(";", "{")), false)
			if !wordHas {
				err = fmt.Errorf("not found the directive end `;` or body block start `{` sat line : %d", line)
				return
			}
			token = char + word
			tokenLine = line
			tokenHas = true
			return
		}
	}
}

func (self *tokenIterator) expectNext(filter filter) (tokens []string, lastToken string, err error) {
	tokens = make([]string, 0)
	var token string
	var has bool
	for {
		if token, _, has, err = self.next(); err != nil {
			return
		} else if has {
			if filter(token, "") {
				lastToken = token
				return
			}
			tokens = append(tokens, token)
		} else {
			return tokens, "", io.EOF
		}
	}
}
