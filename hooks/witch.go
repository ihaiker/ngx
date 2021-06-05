package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

func Switch(params Parameters) Hook {
	return func(conf *config.Configuration) error {
		return setSwitch(conf, params)
	}
}

func setSwitch(conf *config.Configuration, params Parameters) (err error) {
	var items config.Directives
	for idx := 0; ; idx++ {
		if len(conf.Body) == idx {
			return
		}
		item := conf.Body[idx]
		//不是switch指令，继续搜索下级
		if item.Name != "@switch" {
			if len(item.Body) > 0 {
				subConf := &config.Configuration{
					Body: item.Body,
				}
				if err = setSwitch(subConf, params); err != nil {
					return
				}
				item.Body = subConf.Body
			}
			continue
		}

		if len(item.Args) == 0 {
			err = fmt.Errorf("not found switch value at line %d", item.Line)
			return
		}

		value := getSwitchValue(item.Args[0], params)
		if items, err = getSwitchItems(item, value, params); err != nil {
			return
		}

		item.Name = getArys(item.Args, 1)
		item.Args = sliceArgs(item.Args, 2)

		if item.Name == "" {
			conf.Body = append(conf.Body[:idx], append(items, conf.Body[idx+1:]...)...)
			idx += len(items) - 1
		} else {
			item.Body = items
		}
	}
}

func getSwitchValue(name string, params Parameters) string {
	value, err := params.template("${" + name + "}")
	if err != nil {
		return ""
	}
	return value
}

func getSwitchItems(conf *config.Directive, value string, params Parameters) (outs config.Directives, err error) {
	outs = config.Directives{}

	subConf := &config.Configuration{Body: conf.Body}
	if err = setSwitch(subConf, params); err != nil {
		return
	}
	conf.Body = subConf.Body

	appendItem := func(d *config.Directive) {
		if d.Name == "" {
			outs = append(outs, d.Body...)
		} else {
			outs = append(outs, d)
		}
	}

	matched := false
	for _, it := range conf.Body {
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
