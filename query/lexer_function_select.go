package query

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

type selectFunc struct{}

func (s selectFunc) Select(args []FuncArg, items config.Directives, fns map[string]ExecutorFunction) (config.Directives, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("must be 3 paramters")
	}
	//args[0].Directive.
	//fnItems := items
	//args[0].
	//operator := *args[1].Value
	//value := *args[2].Value
	return nil, nil
}
