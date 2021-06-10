package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func subDirectives(it *tokenIterator) ([]*Directive, error) {
	directives := make([]*Directive, 0)
	for {
		token, line, has, err := it.next()
		if err != nil {
			return nil, err
		}
		if !has {
			break
		}
		if token == ";" || token == "}" {
			break
		} else if token[0] == '#' { //注释
			directives = append(directives, &Directive{
				Line: line, Name: "#", Args: []string{strings.Trim(token[1:], " ")},
			})
		} else {
			//@Unsupported， 放弃配置文件制定是否需要添加分隔符。
			//改为默认支持，向下兼容
			if strings.HasSuffix(token, ":") {
				token = token[0 : len(token)-1]
			}

			if args, lastToken, err := it.expectNext(In(";", "{")); err != nil {
				return nil, fmt.Errorf("not found directive end `;` or block directive start `{` at `%s` in %d", token, line)
			} else if lastToken == ";" {
				directives = append(directives, &Directive{
					Line: line, Name: token, Args: args,
				})
			} else {
				directive := &Directive{
					Line: line, Name: token, Args: args,
				}
				if subDirs, err := subDirectives(it); err != nil {
					return nil, err
				} else {
					directive.Body = subDirs
				}
				directives = append(directives, directive)
			}
		}
	}
	return directives, nil
}

func MustParse(filename string) *Configuration {
	cfg, err := Parse(filename)
	if err != nil {
		panic(err)
	}
	return cfg
}

func MustParseBytes(bs []byte) *Configuration {
	cfg, err := ParseBytes(bs)
	if err != nil {
		panic(err)
	}
	return cfg
}

func MustParseIO(reader io.Reader) *Configuration {
	cfg, err := ParseIO(reader)
	if err != nil {
		panic(err)
	}
	return cfg
}

func parse(bs []byte, cfg *Configuration) (err error) {
	it := newTokenIteratorWithBytes(bs)
	cfg.Body, err = subDirectives(it)
	return
}

func Parse(filename string) (cfg *Configuration, err error) {
	cfg = &Configuration{Source: fmt.Sprintf("file://%s", filename)}

	if _, err = os.Stat(filename); !(err == nil || os.IsExist(err)) {
		err = fmt.Errorf("file not found: %s", filename)
		return
	}

	var bs []byte
	if bs, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	err = parse(bs, cfg)
	return
}

func ParseIO(reader io.Reader) (cfg *Configuration, err error) {
	var bs []byte
	if bs, err = ioutil.ReadAll(reader); err != nil {
		return
	}
	cfg = &Configuration{Source: "io"}
	err = parse(bs, cfg)
	return
}

func ParseBytes(bs []byte) (cfg *Configuration, err error) {
	cfg = &Configuration{Source: "bytes"}
	err = parse(bs, cfg)
	return
}
