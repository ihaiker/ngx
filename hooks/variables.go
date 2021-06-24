package hooks

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type (
	Variables struct {
		vars map[string]interface{}
	}
	VariableAdapter interface {
		SetVariables(*Variables)
	}
)

func NewVariables() *Variables {
	v := new(Variables)
	v.vars = map[string]interface{}{}
	envs := map[string]string{}
	for _, env := range os.Environ() {
		nameAndValue := strings.SplitN(env, "=", 2)
		envs[nameAndValue[0]] = nameAndValue[1]
	}
	v.vars["env"] = envs
	v.vars["__func__"] = template.FuncMap{}
	return v
}

func (this *Variables) Func(name string, fn interface{}) *Variables {
	(this.vars["__func__"].(template.FuncMap))[name] = fn
	return this
}

func (this *Variables) Parameter(name string, obj interface{}) *Variables {
	this.vars[name] = obj
	return this
}

func (this *Variables) ExecutorArgs(leftDelims, rightDelims string, arg string) (string, error) {
	if !strings.HasPrefix(arg, leftDelims) {
		return arg, nil
	}

	if strings.HasPrefix(arg, "${.env.") {
		name := arg[7 : len(arg)-len(rightDelims)]
		return os.Getenv(name), nil
	}

	funcs := this.vars["__func__"].(template.FuncMap)
	temp, err := template.New("").Funcs(funcs).
		Delims(leftDelims, rightDelims).Parse(arg)
	if err != nil {
		return "", err
	}
	out := bytes.NewBufferString("")
	err = temp.Execute(out, this.vars)
	return out.String(), err
}

func _is(t reflect.Type, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if t.Kind() == kind ||
			(t.Kind() == reflect.Ptr && t.Elem().Kind() == kind) {
			return true
		}
	}
	return false
}

func (this *Variables) Get(name string) (interface{}, error) {
	names := strings.Split(name[1:], ".")
	value := this.vars[names[0]]
	for _, fieldName := range names[1:] {
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		index := -1
		if strings.Contains(fieldName, "[") {
			idx := strings.Index(fieldName, "[")
			fieldName = fieldName[0:idx]
			index, _ = strconv.Atoi(fieldName[idx : len(fieldName)-1])
		}

		if rv.Kind() == reflect.Map { //map key 必须为string 类型
			fieldValue := rv.MapIndex(reflect.ValueOf(fieldName))
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("not found %s", name)
			}
			value = fieldValue.Interface()
		} else if rv.Kind() == reflect.Struct {
			fieldValue := rv.FieldByName(fieldName)
			if !fieldValue.IsValid() {
				return nil, fmt.Errorf("not found %s", name)
			}
			value = fieldValue.Interface()
		} else {
			return nil, fmt.Errorf("not found %s", name)
		}

		if index != -1 && _is(reflect.TypeOf(value), reflect.Slice) {
			sliceValue := reflect.ValueOf(value)
			if sliceValue.Kind() == reflect.Ptr {
				sliceValue = sliceValue.Elem()
			}
			if index >= sliceValue.Len() {
				return nil, fmt.Errorf("index outof %s", name)
			}
			value = sliceValue.Index(index).Interface()
		}
	}
	return value, nil
}
