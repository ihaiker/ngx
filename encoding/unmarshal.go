package encoding

import (
	"fmt"
	"github.com/xhaiker/ngx/config"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Unmarshaler interface {
	UnmarshalNgx(item config.Directives) error
}

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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool, reflect.String:
	case reflect.Struct:
		return fieldType.String() == "time.Time"
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

func structValue(fieldType reflect.Type, items config.Directives) (reflect.Value, error) {
	value := reflect.New(fieldType)

	body := config.Directives{}

	for _, item := range items {
		for _, arg := range item.Args {
			fieldName, fieldValue := compileSplit2(arg, ":|=")
			body = append(body, &config.Directive{
				Name: fieldName, Args: []string{fieldValue}, Line: item.Line,
			})
		}
		body = append(body, item.Body...)
	}
	err := UnmarshalWith(value.Interface(), body)
	return value.Elem(), err
}

func mapValue(keyType, valueType reflect.Type, items config.Directives, format string) (reflect.Value, error) {
	if !isBase(keyType) {
		return reflect.Value{}, fmt.Errorf("%s not support %s", items[0].Name, keyType)
	}

	m := reflect.MakeMap(reflect.MapOf(keyType, valueType))

	for _, item := range items {
		//基本类型
		if isBase(valueType) {
			if len(item.Args) == 2 {
				//直接设置值 mapField mapKey mapValue;
				key, err := baseValue(keyType, config.New("key", item.Args[0]), format)
				if err != nil {
					return reflect.Value{}, err
				}
				value, err := baseValue(valueType, config.New("value", item.Args[1]), format)
				if err != nil {
					return reflect.Value{}, err
				}
				m.SetMapIndex(key, value)
			} else {
				//直接设置值 mapField { mapKey: mapValue }
				for _, d := range item.Body {
					key, err := baseValue(keyType, config.New("key", d.Name), format)
					if err != nil {
						return reflect.Value{}, err
					}
					value, err := baseValue(valueType, config.New("value", d.Args...), format)
					if err != nil {
						return reflect.Value{}, err
					}
					m.SetMapIndex(key, value)
				}
			}
		} else { //非基本类型
			setMap := func(key string, bodyItems config.Directives) error {
				keyValue, err := baseValue(keyType, config.New("key", key), format)
				if err != nil {
					return err
				}

				vs := reflect.New(valueType)
				if isPtr(valueType) {
					vs = reflect.New(valueType.Elem())
				}

				d := config.Directives{{Name: key, Body: bodyItems}}
				if err := UnmarshalWith(vs.Interface(), d); err != nil {
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
				if err := setMap(item.Args[0], item.Body); err != nil {
					return reflect.Value{}, err
				}
			} else {
				for _, sub := range item.Body {
					if err := setMap(sub.Name, sub.Body); err != nil {
						return reflect.Value{}, err
					}
				}
			}
		}
	}
	return m, nil
}

func sliceValue(sliceType reflect.Type, items config.Directives, format string) (reflect.Value, error) {
	if isBase(sliceType) {
		length := 0
		for _, item := range items {
			length += len(item.Args)
		}
		values := reflect.MakeSlice(reflect.SliceOf(sliceType), length, length)

		idx := 0
		for _, item := range items {
			for _, arg := range item.Args {
				v, err := baseValue(sliceType, config.New("key", arg), format)
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
			if err := UnmarshalWith(vs.Interface(), config.Directives{item}); err != nil {
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

func Unmarshal(data []byte, v interface{}, opt *config.Options) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%s not be a interface", reflect.TypeOf(v))
	}
	if conf, err := config.ParseWith(data, opt); err != nil {
		return err
	} else {
		return UnmarshalWith(v, conf.Body)
	}
}

func UnmarshalWith(v interface{}, item config.Directives) error {
	if us, match := v.(Unmarshaler); match {
		return us.UnmarshalNgx(item)
	}

	value := reflectValue(v)
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

		//实现了 Unmarshaler
		if has, err := unmarshalNgx(value, i, fieldTagName, item); err != nil {
			return err
		} else if has {
			continue
		}

		if v, has := Defaults.dealWith(field.Type, item.Gets(fieldTagName)); has {
			value.Field(i).Set(v)
		} else if isBase(field.Type) {
			if d := item.Get(fieldTagName); d != nil {
				if v, err := baseValue(field.Type, d, format); err != nil {
					return err
				} else {
					value.Field(i).Set(v)
					continue
				}
			}
		}

		fieldValue := value.Field(i)
		if v, err := assemblyValue(field.Type, fieldValue, item.Gets(fieldTagName), format); err == nil {
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

func assemblyValue(fieldType reflect.Type, value reflect.Value, item config.Directives, format string) (reflect.Value, error) {
	//所有处理按照interface处理
	if fieldType.Kind() == reflect.Ptr {
		if out, err := assemblyValue(fieldType.Elem(), value, item, format); err == nil {
			v := reflect.New(fieldType.Elem())
			v.Elem().Set(out)
			return v, nil
		} else {
			return out, err
		}
	}

	switch fieldType.Kind() {
	case reflect.Array:
		//return reflect.Value{}, fmt.Errorf("Invalid %s", item.Name)

	case reflect.Map:
		v := reflect.New(fieldType)
		m, err := mapValue(fieldType.Key(), fieldType.Elem(), item, format)
		if err != nil {
			return v, err
		}
		if value.IsValid() {
			for mr := value.MapRange(); mr.Next(); {
				m.SetMapIndex(mr.Key(), mr.Value())
			}
		}
		v.Elem().Set(m)
		return v, nil

	case reflect.Slice:
		v := reflect.New(fieldType)
		slice, err := sliceValue(fieldType.Elem(), item, format)
		if err != nil {
			return v, err
		}
		v.Elem().Set(slice)
		if value.IsValid() {
			v.Elem().Set(reflect.AppendSlice(value, slice))
		}
		return v, nil

	case reflect.Struct:
		return structValue(fieldType, item)
	}

	return reflect.Value{}, os.ErrInvalid
}

func unmarshalNgx(value reflect.Value, idx int, fieldTagName string, items config.Directives) (bool, error) {
	field := value.Type().Field(idx)

	var fieldValue reflect.Value
	if isPtr(field.Type) {
		fieldValue = reflect.New(field.Type.Elem())
	} else {
		fieldValue = reflect.New(field.Type)
	}

	if us, match := fieldValue.Interface().(Unmarshaler); match {
		if err := us.UnmarshalNgx(items.Gets(fieldTagName)); err != nil {
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
