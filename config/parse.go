package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func subDirectives(it *tokenIterator, opt *Options) ([]*Directive, error) {
	directives := make([]*Directive, 0)
	for {
		token, line, has := it.next()
		if !has {
			break
		}
		if token == ";" || token == "}" {
			break
		} else if token[0] == '#' { //注释
			if !opt.RemoveCommits {
				directives = append(directives, &Directive{
					Line: line, Name: "#", Args: []string{strings.Trim(token[1:], " ")},
				})
			}
		} else {

			if opt.Delimiter {
				if strings.HasSuffix(token, ":") {
					token = token[0 : len(token)-1]
				}
			}

			if args, lastToken, err := it.expectNext(In(";", "{")); err != nil {
				return nil, fmt.Errorf("not found end (%s) [;{] in %d", token, line)
			} else if lastToken == ";" {
				directives = append(directives, &Directive{
					Line: line, Name: token, Args: args,
				})
			} else {
				directive := &Directive{
					Line: line, Name: token, Args: args,
				}
				if subDirs, err := subDirectives(it, opt); err != nil {
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

func MustParse(filename string, opt *Options) *Configuration {
	cfg, err := Parse(filename, opt)
	if err != nil {
		panic(err)
	}
	return cfg
}

func MustParseWith(bs []byte, opt *Options) *Configuration {
	cfg, err := ParseWith(bs, opt)
	if err != nil {
		panic(err)
	}
	return cfg
}

func MustParseIO(reader io.Reader, options *Options) *Configuration {
	cfg, err := ParseIO(reader, options)
	if err != nil {
		panic(err)
	}
	return cfg
}

func parse(bs []byte, cfg *Configuration) (err error) {
	if cfg.Options == nil {
		cfg.Options = &Options{
			Delimiter:     true,
			RemoveQuote:   false,
			RemoveCommits: false,
			MergeInclude:  false,
		}
	}
	it := newTokenIteratorWithBytes(bs, cfg.Options)
	cfg.Body, err = subDirectives(it, cfg.Options)
	return
}

func Parse(filename string, options *Options) (cfg *Configuration, err error) {
	cfg = &Configuration{Source: fmt.Sprintf("file://%s", filename)}
	cfg.Options = options
	var bs []byte
	if bs, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	err = parse(bs, cfg)
	return
}

func ParseIO(reader io.Reader, options *Options) (cfg *Configuration, err error) {
	var bs []byte
	if bs, err = ioutil.ReadAll(reader); err != nil {
		return
	}
	cfg = &Configuration{Source: "io", Options: options}
	err = parse(bs, cfg)
	return
}

func ParseWith(bs []byte, options *Options) (cfg *Configuration, err error) {
	cfg = &Configuration{Source: "bytes", Options: options}
	err = parse(bs, cfg)
	return
}
