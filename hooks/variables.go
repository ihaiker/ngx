package hooks

import (
	"bytes"
	"os"
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
