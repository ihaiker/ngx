package query

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/query/methods"
	"github.com/ihaiker/ngx/v2/query/stdlib"
)

var manager = stdlib.Defaults() //系统标准方法库

//Selects 查询配置内容。
func Selects(conf *config.Configuration, queries ...string) ([]*config.Directive, error) {
	return SelectWithFn(conf, manager, queries...)
}

//查询匹配配内容并附带方法管理器
func SelectWithFn(conf *config.Configuration, manager *methods.FunctionManager, queries ...string) (items []*config.Directive, err error) {
	items = config.Directives{{Name: "", Body: conf.Body}}
	var expr *expression
	for _, query := range queries {
		if expr, err = parseLexer(query); err != nil {
			return
		}
		if items, err = expr.call(items, manager); err != nil {
			return
		}
		if len(items) == 0 {
			return
		}
	}
	return
}
