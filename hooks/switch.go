package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"strings"
)

type SwitchHooker struct {
	*Variables
}

func (s *SwitchHooker) SetVariables(variables *Variables) {
	s.Variables = variables
}

func (s *SwitchHooker) Execute(item *config.Directive, _ Next) (current config.Directives, children config.Directives, err error) {
	if len(item.Args) == 0 {
		err = fmt.Errorf("not found switch value at line %d", item.Line)
		return
	}

	var value string
	var items config.Directives

	if value, err = s.getSwitchValue(item.Args[0]); err != nil {
		return
	}
	if items, err = s.getSwitchItems(item, value); err != nil {
		return
	}

	if len(item.Args) > 1 {
		current = config.Directives{{
			Line:    item.Line,
			Virtual: item.Virtual,
			Name:    getArys(item.Args, 1),
			Args:    sliceArgs(item.Args, 2),
			Body:    items,
		}}
	} else {
		current = items
	}
	return
}

func (s *SwitchHooker) getSwitchValue(name string) (string, error) {
	if !strings.HasPrefix(name, ".") {
		return name, nil
	}
	return s.ExecutorArgs("${", "}", fmt.Sprintf("${%s}", name))
}

func (s *SwitchHooker) getSwitchItems(item *config.Directive, value string) (outs config.Directives, err error) {
	outs = config.Directives{}
	appendItem := func(d *config.Directive) {
		if d.Name == "" {
			outs = append(outs, d.Body...)
		} else {
			outs = append(outs, d)
		}
	}
	matched := false
	for _, it := range item.Body {
		switch it.Name {
		case "@case":
			if it.Args[0] == value {
				it.Name = getArys(it.Args, 1)
				it.Args = sliceArgs(it.Args, 2)
				appendItem(it)
				matched = true
			}
		case "@default":
			it.Name = getArys(it.Args, 0)
			it.Args = sliceArgs(it.Args, 1)
			if !matched {
				appendItem(it)
				matched = true
			}
		default:
			appendItem(it)
		}
	}
	return
}
