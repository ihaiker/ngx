package methods

import (
	"fmt"
	"reflect"
)

type FunctionManager struct {
	fns map[string]interface{}
}

//get a function named `fnName`
func (this *FunctionManager) Get(fnName string) (fn interface{}, has bool) {
	fn, has = this.fns[fnName]
	return
}

//add function to manager. the function will check the parameters type.
//allow the following types: string、int、config.Directives、[]string 、[]int、[]interface{}
func (this *FunctionManager) Add(name string, execFn interface{}) error {
	fn := reflect.ValueOf(execFn)
	if fn.Kind() != reflect.Func {
		return fmt.Errorf("not a function: %s", name)
	}
	for idx := 0; idx < fn.Type().NumIn(); idx++ {
		parameterType := fn.Type().In(idx)
		if !this.checkParameterType(parameterType) {
			return fmt.Errorf("the %d parameter type `%s` of function named %s not allow",
				idx, parameterType.String(), name)
		}
	}

	if fn.Type().NumOut() > 2 {
		return fmt.Errorf("the %s function only allow two out value.", name)
	}

	outType := fn.Type().Out(0)
	if !this.checkParameterType(outType) {
		return fmt.Errorf("the frist out value type `%s` of function named %s not allow",
			outType.String(), name)
	}
	if fn.Type().NumOut() == 2 {
		outType = fn.Type().Out(1)
		if outType.String() != "error" {
			return fmt.Errorf("the second out value type `%s` of function named %s not allow",
				outType.String(), name)
		}
	}
	this.fns[name] = execFn
	return nil
}

// Check the parameter type
func (this *FunctionManager) checkParameterType(parameterType reflect.Type) bool {
	if parameterType.Kind() == reflect.Ptr || parameterType.Kind() == reflect.Array ||
		parameterType.Kind() == reflect.Slice {
		return this.checkParameterType(parameterType.Elem())
	}

	switch parameterType.String() {
	case "error", "string", "int", "config.Directive", "bool", "interface {}":
		return true
	default:
		return false
	}
}

func New() *FunctionManager {
	return &FunctionManager{
		fns: map[string]interface{}{},
	}
}
