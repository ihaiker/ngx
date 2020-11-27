package encoding

import (
	"fmt"
	"github.com/ihaiker/ngx/config"
	"reflect"
	"strconv"
	"time"
)

type Marshaler interface {
	MarshalNgx() (config.Directives, error)
}

func Marshal(v interface{}) ([]byte, error) {
	cfg, err := MarshalConfig(v, time.RFC3339)
	if err != nil {
		return nil, err
	}
	c := config.Body("content", cfg...)
	return c.BodyBytes(), nil
}

func MarshalConfig(v interface{}, format string) (config.Directives, error) {
	if v == nil {
		return nil, nil
	}
	if mg, match := v.(Marshaler); match {
		return mg.MarshalNgx()
	}

	value := reflect.ValueOf(v)
	valueType := value.Type()
	if value.Kind() == reflect.Ptr {
		valueType = value.Elem().Type()
		value = value.Elem()
	}

	items := config.Directives{}
	if valueType.String() == "time.Time" {
		t := value.Interface().(time.Time)
		return config.Directives{config.New("key", t.Format(format))}, nil
	}

	switch valueType.Kind() {
	case reflect.String:
		return config.Directives{config.New("key", strconv.Quote(value.String()))}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return config.Directives{config.New("key", fmt.Sprintf("%d", value.Int()))}, nil
	case reflect.Float32, reflect.Float64:
		return config.Directives{config.New("key", fmt.Sprintf("%f", value.Float()))}, nil
	case reflect.Bool:
		return config.Directives{config.New("key", strconv.FormatBool(value.Bool()))}, nil

	case reflect.Map:
		for mr := value.MapRange(); mr.Next(); {
			item := config.New(mr.Key().String())
			if isBase(valueType.Elem()) {
				item.AddArgs(mr.Value().String())
			} else {
				if d, err := MarshalConfig(mr.Value().Interface(), format); err != nil {
					return nil, err
				} else {
					item.AddBodyDirective(d...)
				}
			}
			items = append(items, item)
		}
	case reflect.Slice:
		if isBase(valueType.Elem()) {
			ary := config.New("array")
			for i := 0; i < value.Len(); i++ {
				v := value.Index(i).Interface()
				if vItem, err := MarshalConfig(v, format); err != nil {
					return nil, err
				} else {
					for _, item := range vItem {
						ary.AddArgs(item.Args...)
					}
				}
			}
			items = append(items, ary)
		} else {
			for i := 0; i < value.Len(); i++ {
				ary := config.New("array")
				v := value.Index(i).Interface()
				if vItem, err := MarshalConfig(v, format); err != nil {
					return nil, err
				} else {
					ary.AddBodyDirective(vItem...)
				}
				items = append(items, ary)
			}
		}

	case reflect.Struct:
		for i := 0; i < value.Type().NumField(); i++ {
			field := value.Type().Field(i)
			fieldValue := value.Field(i)

			if field.Type.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					continue
				}
				fieldValue = fieldValue.Elem()
			} else if fieldValue.IsZero() {
				continue
			}

			fieldName, f := split2(field.Tag.Get("ngx"), ",")
			if fieldName == "" {
				fieldName = field.Name
			}

			if confItems, err := MarshalConfig(fieldValue.Interface(), f); err != nil {
				return nil, err
			} else {
				if isBase(field.Type) || field.Type.Kind() == reflect.Slice {
					for _, item := range confItems {
						item.Name = fieldName
						items = append(items, item)
					}
				} else {
					items = append(items, config.Body(fieldName, confItems...))
				}
			}
		}
	}
	return items, nil
}
