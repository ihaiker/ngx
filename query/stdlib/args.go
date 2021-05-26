package stdlib

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

func init() {
	_methods["args"] = args
}

//获取参数的值，参数item 仅仅可以传入 config.Directive 和 config.Directives
func args(item interface{}) ([]string, error) {
	if c, m := item.(*config.Directive); m {
		return c.Args, nil
	}

	if c, m := item.(config.Directives); m {
		if o := c.Index(0); o != nil {
			return o.Args, nil
		}
	}
	return nil, fmt.Errorf("`args` method parameters only need to allow *config.Driective and config.Directive")
}
