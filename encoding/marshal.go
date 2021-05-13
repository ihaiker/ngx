package encoding

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
	"strconv"
	"time"
)

type Marshaler interface {
	MarshalNgx() (*config.Configuration, error)
}

func Marshal(v interface{}) ([]byte, error) {
	return MarshalWithOptions(v, *defaults)
}

func MarshalWithOptions(v interface{}, options Options) ([]byte, error) {
	cfg, err := MarshalOptions(v, options)
	if err != nil {
		return nil, err
	}
	return []byte(cfg.Pretty()), nil
}

func conf(items ...*config.Directive) *config.Configuration {
	return &config.Configuration{
		Source: "code", Options: config.Encoding(), Body: items,
	}
}
func directive2Conf(item *config.Directive) *config.Configuration {
	return &config.Configuration{
		Source: "code", Options: config.Encoding(), Body: item.Body,
	}
}

func MarshalOptions(v interface{}, options Options) (*config.Configuration, error) {
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

	if items, handlered, err := options.TypeHandlers.MarshalNgx(v); err != nil || handlered {
		return items, err
	}

	items := config.Directives{}
	if valueType.String() == "time.Time" {
		t := value.Interface().(time.Time)
		value := strconv.Quote(t.Format(options.DateFormat))
		return conf(config.New("key", value)), nil
	}

	switch valueType.Kind() {
	case reflect.String:
		return conf(config.New("key", value.String())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return conf(config.New("key", fmt.Sprintf("%d", value.Int()))), nil
	case reflect.Float32, reflect.Float64:
		return conf(config.New("key", fmt.Sprintf("%f", value.Float()))), nil
	case reflect.Bool:
		return conf(config.New("key", strconv.FormatBool(value.Bool()))), nil

	case reflect.Map:
		for mr := value.MapRange(); mr.Next(); {
			item := config.New(mr.Key().String())
			if isBase(valueType.Elem()) {
				item.AddArgs(mr.Value().String())
			} else {
				if d, err := MarshalOptions(mr.Value().Interface(), options); err != nil {
					return nil, err
				} else {
					item.AddBodyDirective(d.Body...)
				}
			}
			items = append(items, item)
		}
	case reflect.Slice:
		if isBase(valueType.Elem()) {
			ary := config.New("array")
			for i := 0; i < value.Len(); i++ {
				val := value.Index(i).Interface()
				if vItem, err := MarshalOptions(val, options); err != nil {
					return nil, err
				} else {
					for _, item := range vItem.Body {
						ary.AddArgs(item.Args...)
					}
				}
			}
			items = append(items, ary)
		} else {
			for i := 0; i < value.Len(); i++ {
				ary := config.New("array")
				val := value.Index(i).Interface()
				if vItem, err := MarshalOptions(val, options); err != nil {
					return nil, err
				} else {
					ary.AddBodyDirective(vItem.Body...)
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

			fieldName, format := split2(field.Tag.Get("ngx"), ",")
			if fieldName == "" {
				fieldName = field.Name
			}

			if fieldValue.Kind().String() == "time.Time" {
				val := fieldValue.Interface().(time.Time).Format(format)
				items = append(items, config.New(fieldName, val))
			} else {
				if confItems, err := MarshalOptions(fieldValue.Interface(), options); err != nil {
					return nil, err
				} else {
					if isBase(field.Type) || field.Type.Kind() == reflect.Slice {
						for _, item := range confItems.Body {
							item.Name = fieldName
							items = append(items, item)
						}
					} else {
						items = append(items, config.Body(fieldName, confItems.Body...))
					}
				}
			}
		}
	}
	return conf(items...), nil
}
