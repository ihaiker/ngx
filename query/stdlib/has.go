package stdlib

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

func init() {
	_methods["has"] = has
}

func has(input interface{}, name string) (bool, error) {
	if item, m := input.(*config.Directive); m {
		for _, directive := range item.Body {
			if directive.Name == name {
				return true, nil
			}
		}
		return false, nil
	}

	if items, m := input.(config.Directives); m {
		for _, item := range items {
			if has, _ := has(item, name); has {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("`has` method parameters only need to allow *config.Driective and config.Directive")
}
