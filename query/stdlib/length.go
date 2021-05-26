package stdlib

import (
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
)

func init() {
	_methods["length"] = length
}

func length(obj interface{}) int {
	if items, m := obj.(config.Directives); m {
		return len(items)
	} else if ary := reflect.ValueOf(obj); ary.Kind() == reflect.Slice {
		return ary.Len()
	} else if obj != nil {
		return 1
	}
	return 0
}
