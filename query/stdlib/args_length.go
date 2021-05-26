package stdlib

import (
	"github.com/ihaiker/ngx/v2/config"
)

func init() {
	_methods["args_length"] = args_length
}

func args_length(item *config.Directive) (int, error) {
	if argsValue, err := args(item); err != nil {
		return 0, err
	} else {
		return length(argsValue), nil
	}
}
