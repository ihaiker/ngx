package encoding

import (
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
	"time"
)

type (
	TypeHandler interface {
		MarshalNgx(v interface{}) (*config.Configuration, error)
		UnmarshalNgx(item *config.Configuration) (v interface{}, err error)
	}
	typeManager map[reflect.Type]TypeHandler
	Options     struct {
		DateFormat   string
		TypeHandlers typeManager
	}
)

func (h typeManager) UnmarshalNgx(fieldType reflect.Type, item *config.Configuration) (value reflect.Value, handlered bool, err error) {
	if handler, has := h[fieldType]; has {
		handlered = true
		var val interface{}
		val, err = handler.UnmarshalNgx(item)
		value = reflect.ValueOf(val)
		if fieldType.Kind() != reflect.Ptr {
			value = value.Elem()
		}
	}
	return
}

func (h typeManager) MarshalNgx(val interface{}) (items *config.Configuration, handlered bool, err error) {
	fieldType := reflect.ValueOf(val).Type()
	if handler, has := h[fieldType]; has {
		handlered = true
		items, err = handler.MarshalNgx(val)
	}
	return
}

func (h *typeManager) Reg(v interface{}, handler TypeHandler) *typeManager {
	if v == nil {
		return h
	}

	fieldType := reflect.ValueOf(v).Type()
	(*h)[fieldType] = handler
	if fieldType.Kind() == reflect.Ptr {
		(*h)[fieldType.Elem()] = handler
	} else {
		(*h)[reflect.PtrTo(fieldType)] = handler
	}
	return h
}

var defaults = DefaultOptions()

func DefaultOptions() *Options {
	return &Options{
		DateFormat:   time.RFC3339,
		TypeHandlers: map[reflect.Type]TypeHandler{},
	}
}

func RegTypeHandler(v interface{}, handler TypeHandler) {
	defaults.TypeHandlers.Reg(v, handler)
}
