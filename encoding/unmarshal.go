package encoding

import (
	"fmt"
	"github.com/ihaiker/ngx/config"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Unmarshaler interface {
	UnmarshalNgx(item *config.Configuration) error
}

func Unmarshal(data []byte, v interface{}) error {
	return UnmarshalWithOptions(data, v, *defaults)
}

func UnmarshalWithOptions(data []byte, v interface{}, options Options) error {
	if v == nil {
		return fmt.Errorf("the value is empty")
	}
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%s not be a interface", reflect.TypeOf(v))
	}

	if conf, err := config.ParseWith(data, encodingOptions()); err != nil {
		return err
	} else {
		return UnmarshalDirectives(v, conf, options)
	}
}

func UnmarshalDirectives(v interface{}, cfg *config.Configuration, options Options) error {
	if us, match := v.(Unmarshaler); match {
		return us.UnmarshalNgx(cfg)
	}

	if cfg == nil || len(cfg.Body) == 0 {
		return nil
	}

	value := reflectValue(v)

	if val, handled, err := options.TypeHandlers.UnmarshalNgx(value.Type(), cfg); err != nil {
		return err
	} else if handled {
		value.Set(val)
		return nil
	}

	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("ngx")

		fieldTagName, format := split2(tag, ",")
		if fieldTagName == "-" {
			continue
		}
		if fieldTagName == "" {
			fieldTagName = field.Name
		}

		items := cfg.Body.Gets(fieldTagName)
		if len(items) == 0 {
			continue
		}

		//当用户指定了 TypeHandler 或者实现了 Unmarshaler
		itemConfig := directive2Conf(items[0])

		//实现了 Unmarshaler, 这种情况下只能使用
		if has, err := unmarshalNgx(value, i, fieldTagName, itemConfig); err != nil {
			return err
		} else if has {
			continue
		}

		if val, has, err := options.TypeHandlers.UnmarshalNgx(field.Type, itemConfig); err != nil {
			return err
		} else if has {
			value.Field(i).Set(val)
			continue
		}

		if isBase(field.Type) {
			if d := cfg.Body.Get(fieldTagName); d != nil {
				if val, err := baseValue(field.Type, d, format); err != nil {
					return err
				} else {
					value.Field(i).Set(val)
				}
			}
			continue
		}

		fieldValue := value.Field(i)
		if v, err := assemblyValue(field.Type, fieldValue, items, options); err == nil {
			if isPtr(field.Type) {
				fieldValue.Set(v)
			} else {
				fieldValue.Set(v.Elem())
			}
			continue
		} else {
			return err
		}
	}
	return nil
}

//使用返回返回 value 值，
func reflectValue(obj interface{}) reflect.Value {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
func isPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}

func isString(fieldType reflect.Type) bool {
	if isPtr(fieldType) {
		fieldType = fieldType.Elem()
	}
	return fieldType.Kind() == reflect.String
}

//是否是简单类型
func isBase(fieldType reflect.Type) bool {
	if isPtr(fieldType) {
		fieldType = fieldType.Elem()
	}
	switch fieldType.Kind() {
	default:
		return false
	case reflect.Struct:
		return fieldType.String() == "time.Time"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool, reflect.String:
	}
	return true
}

func baseValue(fieldType reflect.Type, item *config.Directive, format string) (reflect.Value, error) {
	if isPtr(fieldType) {
		out, err := baseValue(fieldType.Elem(), item, format)
		if err == nil {
			v := reflect.New(fieldType.Elem())
			v.Elem().Set(out)
			return v, err
		}
		return out, err
	}
	sv := index(item.Args, 0)
	v := reflect.New(fieldType)
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(sv, 10, 64); err != nil {
			return v, err
		} else {
			v.Elem().SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(sv, 10, 64); err != nil {
			return v, err
		} else {
			v.Elem().SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(sv, 64); err != nil {
			return v, err
		} else {
			v.Elem().SetFloat(f)
		}
	case reflect.Bool:
		v.Elem().SetBool(sv == "" || sv == "true")

	case reflect.String:

		v.Elem().SetString(sv)
	case reflect.Struct:
		if fieldType.String() != "time.Time" {
			return v, os.ErrInvalid
		}
		if format == "" {
			format = time.RFC3339
		}
		if t, err := time.Parse(format, sv); err != nil {
			return v, err
		} else {
			v.Elem().Set(reflect.ValueOf(t))
		}
	default:
		return v, os.ErrInvalid
	}
	return v.Elem(), nil
}

func structValue(fieldType reflect.Type, fieldValue reflect.Value, items config.Directives, opt Options) (reflect.Value, error) {
	value := reflect.New(fieldType)

	body := conf()
	for _, item := range items {
		for _, arg := range item.Args {
			fieldName, fieldValue := compileSplit2(arg, ":|=")
			body.Body = append(body.Body, &config.Directive{
				Name: fieldName, Args: []string{fieldValue}, Line: item.Line,
			})
		}
		body.Body = append(body.Body, item.Body...)
	}

	err := UnmarshalDirectives(value.Interface(), body, opt)
	return value, err
}

func mapValue(keyType, valueType reflect.Type, items config.Directives, opt Options) (reflect.Value, error) {
	if !isBase(keyType) {
		return reflect.Value{}, fmt.Errorf("%s not support %s", items[0].Name, keyType)
	}

	m := reflect.MakeMap(reflect.MapOf(keyType, valueType))

	for _, item := range items {
		//基本类型
		if isBase(valueType) {
			if len(item.Args) == 2 {
				//直接设置值 mapField mapKey mapValue;
				key, err := baseValue(keyType, config.New("key", item.Args[0]), opt.DateFormat)
				if err != nil {
					return reflect.Value{}, err
				}
				value, err := baseValue(valueType, config.New("value", item.Args[1]), opt.DateFormat)
				if err != nil {
					return reflect.Value{}, err
				}
				m.SetMapIndex(key, value)
			} else {
				//直接设置值 mapField { mapKey: mapValue }
				for _, d := range item.Body {
					key, err := baseValue(keyType, config.New("key", d.Name), opt.DateFormat)
					if err != nil {
						return reflect.Value{}, err
					}
					value, err := baseValue(valueType, config.New("value", d.Args...), opt.DateFormat)
					if err != nil {
						return reflect.Value{}, err
					}
					m.SetMapIndex(key, value)
				}
			}
		} else { //非基本类型
			setMap := func(key string, cfg *config.Configuration) error {
				keyValue, err := baseValue(keyType, config.New("key", key), opt.DateFormat)
				if err != nil {
					return err
				}

				vs := reflect.New(valueType)
				if isPtr(valueType) {
					vs = reflect.New(valueType.Elem())
				}

				if err := UnmarshalDirectives(vs.Interface(), cfg, opt); err != nil {
					return err
				}
				if isPtr(valueType) {
					m.SetMapIndex(keyValue, vs)
				} else {
					m.SetMapIndex(keyValue, vs.Elem())
				}
				return nil
			}

			if len(item.Args) == 1 {
				if err := setMap(item.Args[0], conf(item.Body...)); err != nil {
					return reflect.Value{}, err
				}
			} else {
				for _, sub := range item.Body {
					if err := setMap(sub.Name, conf(sub.Body...)); err != nil {
						return reflect.Value{}, err
					}
				}
			}
		}
	}
	return m, nil
}

func sliceValue(sliceType reflect.Type, items config.Directives, opt Options) (reflect.Value, error) {
	if isBase(sliceType) {
		length := 0
		for _, item := range items {
			length += len(item.Args)
		}
		values := reflect.MakeSlice(reflect.SliceOf(sliceType), length, length)

		idx := 0
		for _, item := range items {
			for _, arg := range item.Args {
				v, err := baseValue(sliceType, config.New("key", arg), opt.DateFormat)
				if err != nil {
					return reflect.Value{}, err
				}
				values.Index(idx).Set(v)
				idx++
			}
		}
		return values, nil
	} else {
		length := len(items)
		slice := reflect.MakeSlice(reflect.SliceOf(sliceType), length, length)
		for i, item := range items {
			vs := reflect.New(sliceType)
			if isPtr(sliceType) {
				vs = reflect.New(sliceType.Elem())
			}
			if err := UnmarshalDirectives(vs.Interface(), conf(item.Body...), opt); err != nil {
				return reflect.Value{}, nil
			} else {
				if isPtr(sliceType) {
					slice.Index(i).Set(vs)
				} else {
					slice.Index(i).Set(vs.Elem())
				}
			}
		}
		return slice, nil
	}
}

func assemblyValue(fieldType reflect.Type, fieldValue reflect.Value, items config.Directives, opt Options) (reflect.Value, error) {
	//所有处理按照interface处理
	if fieldType.Kind() == reflect.Ptr {
		if out, err := assemblyValue(fieldType.Elem(), fieldValue, items, opt); err == nil {
			v := reflect.New(fieldType.Elem())
			v.Elem().Set(reflect.Indirect(out))
			return v, nil
		} else {
			return out, err
		}
	}

	switch fieldType.Kind() {
	case reflect.Array:
		//return reflect.Value{}, fmt.Errorf("Invalid %s", items.Name)

	case reflect.Map:
		val := reflect.New(fieldType)
		m, err := mapValue(fieldType.Key(), fieldType.Elem(), items, opt)
		if err != nil {
			return val, err
		}
		if fieldValue.IsValid() {
			for mr := fieldValue.MapRange(); mr.Next(); {
				m.SetMapIndex(mr.Key(), mr.Value())
			}
		}
		val.Elem().Set(m)
		return val, nil

	case reflect.Slice:
		val := reflect.New(fieldType)
		slice, err := sliceValue(fieldType.Elem(), items, opt)
		if err != nil {
			return val, err
		}
		val.Elem().Set(slice)
		if fieldValue.IsValid() {
			val.Elem().Set(reflect.AppendSlice(fieldValue, slice))
		}
		return val, nil

	case reflect.Struct:
		val, err := structValue(fieldType, fieldValue, items, opt)
		return val, err
	}

	return reflect.Value{}, fmt.Errorf("不支持: %s", fieldType.String())
}

func unmarshalNgx(value reflect.Value, idx int, fieldTagName string, config *config.Configuration) (bool, error) {
	field := value.Type().Field(idx)

	var fieldValue reflect.Value
	if isPtr(field.Type) {
		fieldValue = reflect.New(field.Type.Elem())
	} else {
		fieldValue = reflect.New(field.Type)
	}

	if us, match := fieldValue.Interface().(Unmarshaler); match {
		if err := us.UnmarshalNgx(config); err != nil {
			return false, err
		}
		if isPtr(field.Type) {
			value.Field(idx).Set(fieldValue)
		} else {
			value.Field(idx).Set(fieldValue.Elem())
		}
		return true, nil
	}
	return false, nil
}
