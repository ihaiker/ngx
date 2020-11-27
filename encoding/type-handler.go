package encoding

import (
	"github.com/ihaiker/ngx/config"
	"reflect"
)

type (
	TypeHandler func(fieldType reflect.Type, item config.Directives) interface{}
	typeManager map[reflect.Type]TypeHandler
)

func (h *typeManager) dealWith(fieldType reflect.Type, item config.Directives) (reflect.Value, bool) {
	if handler, has := (*h)[fieldType]; has {
		v := handler(fieldType, item)
		return reflect.ValueOf(v), true
	}
	return reflect.Value{}, false
}

func (h *typeManager) With(fieldType reflect.Type, handler TypeHandler) *typeManager {
	(*h)[fieldType] = handler
	if fieldType.Kind() == reflect.Ptr {
		(*h)[fieldType.Elem()] = handler
	} else {
		(*h)[reflect.PtrTo(fieldType)] = handler
	}
	return h
}

var Defaults = new(typeManager)
