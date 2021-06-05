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

type Parameters map[string]interface{}

func Parameter() Parameters {
	ps := map[string]interface{}{}

	envs := map[string]string{}
	for _, env := range os.Environ() {
		nameAndValue := strings.SplitN(env, "=", 2)
		envs[nameAndValue[0]] = nameAndValue[1]
	}
	(ps)["env"] = envs
	(ps)["__func__"] = template.FuncMap{}
	return ps
}

func (this *Parameters) Func(name string, fn interface{}) *Parameters {
	((*this)["__func__"].(template.FuncMap))[name] = fn
	return this
}

func (this *Parameters) Add(name string, obj interface{}) *Parameters {
	(*this)[name] = obj
	return this
}

func (this *Parameters) FuncMap() template.FuncMap {
	return (*this)["__func__"].(template.FuncMap)
}

// 参数替换，使用text/template方式来替换，且设置delims为 `${` `}`
func (this *Parameters) Hook() Hook {
	return func(conf *config.Configuration) error {
		return this.replace(conf.Body)
	}
}

func (this *Parameters) template(arg string) (string, error) {
	temp, err := template.New("").Funcs(this.FuncMap()).
		Delims("${", "}").Parse(arg)
	if err != nil {
		return arg, err
	}
	out := bytes.NewBufferString("")
	err = temp.Execute(out, this)
	arg = out.String()
	return arg, err
}

func (this *Parameters) replace(items config.Directives) (err error) {
	for _, item := range items {
		for idx, arg := range item.Args {
			if compile.MatchString(arg) {
				if item.Args[idx], err = this.template(arg); err != nil {
					return
				}
			}
		}
		if len(item.Body) > 0 {
			if err = this.replace(item.Body); err != nil {
				return
			}
		}
	}
	return nil
}

func (this *Parameters) Exec(conf *config.Configuration) error {
	return this.Hook()(conf)
}
