package hooks

import (
	"github.com/ihaiker/ngx/v2/config"
	"os"
	"strings"
)

type (
	Next func(filter func(*config.Directive) bool) *config.Directive
	Hook interface {
		Execute(item *config.Directive, next Next) (current config.Directives, children config.Directives, err error)
	}
	Hooks interface {
		Execute(conf *config.Configuration) error
		RegHook(hook Hook, names ...string) Hooks
		Vars() *Variables
	}
	HooksAdapter interface {
		SetHooks(hooks Hooks)
	}

	_hooks struct {
		hooks map[string]Hook
		*Variables
	}
)

func New() Hooks {
	hooks := &_hooks{
		hooks:     map[string]Hook{},
		Variables: NewVariables(),
	}
	root, _ := os.Getwd()
	hooks.RegHook(&IncludeHooker{Merge: false, Search: WalkFiles(root)}, "include")
	hooks.RegHook(new(SwitchHooker), "@switch")
	hooks.RegHook(new(RepeatHooker), "@repeat")
	hooks.RegHook(new(TemplateHooker), "@template", "@merge", "@include")
	hooks.RegHook(new(IfElseHooker), "@if")
	return hooks
}

var Defaults = New()

// 注册hook, 可以注册多个名字
func (this *_hooks) RegHook(hook Hook, names ...string) Hooks {
	for _, name := range names {
		this.hooks[name] = hook
	}
	return this
}

func (this *_hooks) Vars() *Variables {
	return this.Variables
}

func (this *_hooks) Execute(conf *config.Configuration) (err error) {

	var current config.Directives
	var children config.Directives

	for idx := 0; ; idx++ {
		if idx == len(conf.Body) {
			break
		}

		next := func(filter func(*config.Directive) bool) *config.Directive {
			if idx+1 < len(conf.Body) && filter(conf.Body[idx+1]) {
				item := conf.Body[idx+1]
				conf.Body = append(conf.Body[:idx+1], conf.Body[idx+2:]...)
				return item
			}
			return nil
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
			if hooks, match := hook.(HooksAdapter); match {
				hooks.SetHooks(this)
			}

			if current, children, err = hook.Execute(item, next); err != nil {
				return
			}
			if len(current) > 0 {
				conf.Body = append(conf.Body[:idx], append(current, conf.Body[idx+1:]...)...)
				//这里处理include命令的特除性，因为include使用merge模式，当前指令变化，需要再次检查
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
