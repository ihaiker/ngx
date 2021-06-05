package hooks

import (
	"bytes"
	"github.com/ihaiker/ngx/v2/config"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var compile = regexp.MustCompile(`\$\{.*\}`)

type ParametersHooker struct {
	params map[string]interface{}
}

func Parameter() *ParametersHooker {
	ps := &ParametersHooker{
		params: map[string]interface{}{},
	}

	envs := map[string]string{}
	for _, env := range os.Environ() {
		nameAndValue := strings.SplitN(env, "=", 2)
		envs[nameAndValue[0]] = nameAndValue[1]
	}
	ps.params["env"] = envs
	ps.params["__func__"] = template.FuncMap{}
	return ps
}

func (this *ParametersHooker) Execute(conf *config.Configuration) (err error) {
	for _, item := range conf.Body {
		for idx, arg := range item.Args {
			if compile.MatchString(arg) {
				if item.Args[idx], err = this.template(arg); err != nil {
					return
				}
			}
		}
		if len(item.Body) > 0 {
			subConf := &config.Configuration{Body: item.Body}
			if err = this.Execute(subConf); err != nil {
				return
			}
		}
	}
	return
}

func (this *ParametersHooker) Func(name string, fn interface{}) *ParametersHooker {
	(this.params["__func__"].(template.FuncMap))[name] = fn
	return this
}

func (this *ParametersHooker) Add(name string, obj interface{}) *ParametersHooker {
	this.params[name] = obj
	return this
}

func (this *ParametersHooker) FuncMap() template.FuncMap {
	return this.params["__func__"].(template.FuncMap)
}

func (this *ParametersHooker) template(arg string) (string, error) {
	temp, err := template.New("").Funcs(this.FuncMap()).
		Delims("${", "}").Parse(arg)
	if err != nil {
		return arg, err
	}
	out := bytes.NewBufferString("")
	err = temp.Execute(out, this.params)
	arg = out.String()
	return arg, err
}
