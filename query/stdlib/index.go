package stdlib

import (
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
)

func init() {
	_methods["index"] = index
}

func index(obj interface{}, index int) interface{} {
	if items, match := obj.(config.Directives); match {
		if index < 0 {
			index = len(items) + index
		}
		if index < 0 || index >= len(items) {
			return nil
		}
		return items[index]
	} else if v := reflect.ValueOf(obj); v.Kind() == reflect.Slice {
		if index < 0 {
			index = v.Len() + index
		}
		if index < 0 || index >= v.Len() {
			return nil
		}
		return v.Index(index).Interface()
	} else {
		return nil
	}
}
