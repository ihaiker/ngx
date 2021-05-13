package query

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"strings"
)

//按照条件查询,
func selectOne(directives []*config.Directive, q string) ([]*config.Directive, error) {
	expr, err := Lexer(q)
	if err != nil {
		return nil, fmt.Errorf("search condition error：[%s]", q)
	}
	matched := make([]*config.Directive, 0)
	for _, directive := range directives {
		for _, body := range directive.Body {
			if expr.Match(body) {
				matched = append(matched, body)
			}
		}
	}
	return matched, nil
}

//Selects 查询配置内容。
// server.[server_name('_') & listen('8081' | '8080')]
//
func Selects(conf *config.Configuration, queries ...string) ([]*config.Directive, error) {
	current := conf.Body //[]*config.Directive{conf}
	for _, q := range queries {
		directives, err := selectOne(current, q)
		if err != nil {
			return nil, err
		}
		if directives == nil || len(directives) == 0 {
			return nil, fmt.Errorf("not found: %s", strings.Join(queries, " "))
		}
		current = directives
	}
	return current, nil
}
