package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/encoding"
	"github.com/ihaiker/ngx/v2/query/methods"
	"reflect"
)

// args('')
type function struct {
	Pos  lexer.Position
	Name string    `@Ident`
	Args []funcArg `"("@@ [("," Space @@)+] ")"`
}

type funcArg struct {
	Pos       lexer.Position
	Directive directives `(@@+`
	Index     *int       `|@Number`
	Value     *string    `|@String`
	Boolean   *bool      `|@("true"|"false"|"True"|"False")`
	Arrays    []funcArg  `|("[" @@ ("," Space @@)* "]")`
	Function  *function  `|@@)`
}

func (this *function) callForExpression(items config.Directives, manager *methods.FunctionManager) (config.Directives, error) {
	matched := config.Directives{}
	for _, item := range items {
		if val, err := this.call(config.Directives{item}, manager); err != nil {
			return nil, err
		} else {
			if val.Type().String() == "bool" {
				if val.Bool() {
					matched = append(matched, item)
				}
			} else {
				if val.Kind() != reflect.Ptr || !val.IsNil() {
					if val.Type().String() == "config.Directives" {
						matched = append(matched, val.Interface().(config.Directives)...)
					} else if val.Type().String() == "*config.Directive" {
						matched = append(matched, val.Interface().(*config.Directive))
					} else {
						if conf, err := encoding.MarshalOptions(val.Interface(), *encoding.DefaultOptions()); err == nil {
							matched = append(matched, conf.Body...)
						}
					}
				}
			}
		}
	}
	return matched, nil
}

func (this *function) call(items config.Directives, manager *methods.FunctionManager) (reflect.Value, error) {
	function, has := manager.Get(this.Name)
	if !has {
		return reflect.Value{}, fmt.Errorf("the function %s not found", this.Name)
	}

	fn := reflect.ValueOf(function)
	inParams := make([]reflect.Value, fn.Type().NumIn())
	var err error
	for idx := 0; idx < fn.Type().NumIn(); idx++ {
		arg := this.Args[idx]
		if inParams[idx], err = arg.call(items, manager); err != nil {
			return reflect.Value{}, err
		}
		argType := fn.Type().In(idx)
		if argType.Kind() == reflect.Ptr {
			inParams[idx] = inParams[idx].Addr()
		}
	}

	values := fn.Call(inParams)
	if len(values) == 2 {
		if !values[1].IsNil() {
			err = values[1].Interface().(error)
		}
	}
	return values[0], err
}

func (arg *funcArg) call(items config.Directives, manager *methods.FunctionManager) (reflect.Value, error) {
	if arg.Directive != nil {
		itemValues, err := arg.Directive.call(items)
		return reflect.ValueOf(itemValues), err
	}
	if arg.Index != nil {
		return reflect.ValueOf(*arg.Index), nil
	}
	if arg.Boolean != nil {
		return reflect.ValueOf(*arg.Boolean), nil
	}
	if arg.Value != nil {
		return reflect.ValueOf(*arg.Value), nil
	}

	if arg.Arrays != nil {
		var err error
		length := len(arg.Arrays)
		//获取参数值
		values := make([]reflect.Value, length)
		for i, array := range arg.Arrays {
			if values[i], err = array.call(items, manager); err != nil {
				return reflect.Value{}, err
			}
		}
		//获取切片的类型
		var aryType reflect.Type
		for _, value := range values {
			if aryType == nil {
				aryType = value.Type()
			} else if aryType != value.Type() {
				// []interface{}
				aryType = reflect.TypeOf(func(n interface{}) {}).In(0)
			}
		}
		//设置切片值
		value := reflect.MakeSlice(aryType, length, length)
		for i, v := range values {
			value.Index(i).Set(v)
		}
		return value, nil
	}

	return arg.Function.call(items, manager)
}
