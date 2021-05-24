package query

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

type selectFunc struct{}

func (this *selectFunc) Select(args []FuncArg, items config.Directives, fns map[string]ExecutorFunction) (matched config.Directives, err error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("must be more then three paramters")
	}

	var match bool
	for _, item := range items {
		if match, err = this.checkItem(args, item, fns); err != nil {
			return
		} else if match {
			matched = append(matched, item)
		}
	}
	return
}

func (this *selectFunc) checkItem(args []FuncArg, item *config.Directive, fns map[string]ExecutorFunction) (matched bool, err error) {
	var selectItems config.Directives
	if args[0].Directive != nil {
		if selectItems, err = args[0].Directive.Select(item.Body); err != nil {
			return
		}
	} else if args[0].Function != nil {
		if selectItems, err = args[0].Function.Select(item.Body, fns); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("invoid first paramter type")
		return
	}
	operator := *args[1].Value
	value := *args[2].Value

	for _, selectItem := range selectItems {
		switch operator {
		case "equal":
			matched = len(selectItem.Args) == 1 && value == selectItem.Args[0]
			return
		}
	}
	return
}
