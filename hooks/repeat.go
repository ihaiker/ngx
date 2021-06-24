package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
	"strings"
)

type RepeatHooker struct {
	vars  *Variables
	hooks Hooks
}

func (this *RepeatHooker) SetHooks(hooks Hooks) {
	this.hooks = hooks
}

type repeatValueInfo struct {
	Length, Index   int
	IsFirst, IsLast bool
}

func (this *RepeatHooker) SetVariables(variables *Variables) {
	this.vars = variables
}

func (this *RepeatHooker) Execute(item *config.Directive) (current config.Directives, children config.Directives, err error) {
	var values []interface{}
	var key string //循环变量的名字
	if key, values, err = this.getRepeatArray(item); err != nil {
		return
	}
	if key == "" {
		key = "item"
	}

	current = make([]*config.Directive, 0)
	for idx, value := range values {

		this.vars.Parameter(key, value)
		this.vars.Parameter(key+"_info", repeatValueInfo{
			Length: len(values), Index: idx,
			IsFirst: idx == 0, IsLast: idx == len(values)-1,
		})

		items := item.Body.Clone()
		if err = this.executeValue(items, value); err != nil {
			return
		}
		current = append(current, items...)
	}
	return
}

func (this *RepeatHooker) executeValue(items config.Directives, value interface{}) (err error) {
	for _, item := range items {
		for i, arg := range item.Args {
			if item.Args[i], err = this.vars.ExecutorArgs("${", "}", arg); err != nil {
				return
			}
		}
		conf := &config.Configuration{Body: item.Body}
		if err = this.hooks.Execute(conf); err != nil {
			return
		}
	}
	return
}

func (this *RepeatHooker) getRepeatArg(name string) (values []interface{}, err error) {
	var obj interface{}
	values = make([]interface{}, 0)
	if obj, err = this.vars.Get(name); err != nil || obj == nil {
		return
	}

	vs := reflect.ValueOf(obj)
	if vs.Kind() == reflect.Ptr { //指针对象转换
		vs = vs.Elem()
	}
	if vs.Kind() == reflect.Slice {
		for i := 0; i < vs.Len(); i++ {
			values = append(values, vs.Index(i).Interface())
		}
	} else {
		err = fmt.Errorf("%s not a slice object", name)
	}
	return
}

func (this *RepeatHooker) getRepeatArray(item *config.Directive) (key string, values []interface{}, err error) {
	values = make([]interface{}, 0)
	//查找给出 vars.args 参数作为repeat
	if len(item.Args) == 1 {
		name := item.Args[0]
		//此参数必须 . 开头
		if !strings.HasPrefix(name, ".") {
			err = fmt.Errorf("invalid arguments at %s: %d", item.Name, item.Line)
		} else {
			values, err = this.getRepeatArg(name)
		}
		return
	} else if len(item.Args) == 3 && item.Args[1] == "in" { // key in array
		values, err = this.getRepeatArg(item.Args[2])
		key = item.Args[0]
	} else if len(item.Args) == 0 {
		//查找@args参数，作为repeat循环
		for idx := 0; ; idx++ {
			if len(item.Body) == idx {
				break
			}
			subItem := item.Body[idx].Clone()
			if subItem.Name == "@args" {
				subItem.Name = getArys(subItem.Args, 0)
				subItem.Args = sliceArgs(subItem.Args, 1)
				values = append(values, subItem)
				item.Body = append(item.Body[:idx], item.Body[idx+1:]...)
				idx--
			}
		}
	} else {
		err = fmt.Errorf("invalid grammar at (%s %s):%d", item.Name, strings.Join(item.Args, ","), item.Line)
	}
	return
}
