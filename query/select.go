package query

import (
	"github.com/ihaiker/ngx/v2/config"
)

//Selects 查询配置内容。
func Selects(conf *config.Configuration, queries ...string) ([]*config.Directive, error) {
	current := []*config.Directive{{
		Line: 1, Name: conf.Source,
		Body: conf.Body,
	}}
	return current, nil
}
