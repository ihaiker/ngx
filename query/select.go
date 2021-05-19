package query

import (
	"github.com/ihaiker/ngx/v2/config"
)

//Selects 查询配置内容。
func Selects(conf *config.Configuration, queries ...string) (items []*config.Directive, err error) {
	items = conf.Body
	var expr *Expression

	for _, query := range queries {
		if expr, err = Lexer(query); err != nil {
			return
		}
		if items, err = expr.Select(items); err != nil {
			return
		}
		if len(items) == 0 {
			return
		}
	}
	return
}
