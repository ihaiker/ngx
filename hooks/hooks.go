package hooks

import (
	"github.com/ihaiker/ngx/v2/config"
	"os"
	"strings"
)

type (
	Hook interface {
		Execute(item *config.Directive) (current config.Directives, children config.Directives, err error)
	}
	Hooks struct {
		hooks map[string]Hook
		*Variables
	}
)

func New() *Hooks {
	hooks := &Hooks{
		hooks:     map[string]Hook{},
		Variables: NewVariables(),
	}
	root, _ := os.Getwd()
	hooks.Hook(&IncludeHooker{Merge: false, Search: WalkFiles(root)}, "include")
	hooks.Hook(new(SwitchHooker), "@switch")
	hooks.Hook(new(RepeatHooker), "@repeat")
	hooks.Hook(new(TemplateHooker), "@template", "@merge", "@include")
	return hooks
}

var Defaults = New()

// 注册hook, 可以注册多个名字
func (this *Hooks) Hook(hook Hook, names ...string) *Hooks {
	for _, name := range names {
		this.hooks[name] = hook
	}
	return this
}

func (this *Hooks) Execute(conf *config.Configuration) (err error) {

	var current config.Directives
	var children config.Directives

	for idx := 0; ; idx++ {
		if idx == len(conf.Body) {
			break
		}
		item := conf.Body[idx]
		for i, arg := range item.Args {
			if item.Args[i], err = this.ExecutorArgs("${", "}", arg); err != nil {
				return
			}
		}

		if hook, has := this.hooks[item.Name]; has {
			if variable, match := hook.(VariableAdapter); match {
				variable.SetVariables(this.Variables)
			}
			if current, children, err = hook.Execute(item); err != nil {
				return
			}
			if len(current) > 0 {
				conf.Body = append(conf.Body[:idx], append(current, conf.Body[idx+1:]...)...)
				//这里处理include命令的特除性
				if item.Name != current[0].Name || strings.Join(item.Args, ",") != strings.Join(current[0].Args, ",") {
					idx-- //由hook处理完成后，需要再次检查，
					continue
				}
			} else if len(children) > 0 {
				item.Body = append(item.Body, children...)
			} else {
				//处理仅仅只是删除当前节点
				conf.Body = append(conf.Body[:idx], conf.Body[idx+1:]...)
				idx--
				continue
			}
		}

		item = conf.Body[idx]
		if len(item.Body) > 0 {
			sub := &config.Configuration{Body: item.Body}
			if err = this.Execute(sub); err != nil {
				return
			}
			item.Body = sub.Body
		}
	}
	return
}
