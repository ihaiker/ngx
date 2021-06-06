package hooks

import (
	"github.com/ihaiker/ngx/v2/config"
)

type RepeatHooker struct {
	vars *Variables
}

func (this *RepeatHooker) SetVariables(variables *Variables) {
	this.vars = variables
}

func (this *RepeatHooker) Execute(item *config.Directive) (current config.Directives, children config.Directives, err error) {
	var values []interface{}
	if values, err = this.getRepeatArray(item); err != nil {
		return
	}

	current = make([]*config.Directive, 0)
	for idx, value := range values {
		this.vars.Parameter("item", value)
		this.vars.Parameter("item_length", len(values))
		this.vars.Parameter("item_index", idx)
		this.vars.Parameter("item_first", idx == 0)
		this.vars.Parameter("item_last", idx == len(values)-1)

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
		if len(item.Body) > 0 {
			if err = this.executeValue(item.Body, value); err != nil {
				return
			}
		}
	}
	return
}

func (this *RepeatHooker) getRepeatArray(item *config.Directive) (values []interface{}, err error) {
	values = make([]interface{}, 0)
	if len(item.Args) > 0 {
		return
	}

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
	return
}
