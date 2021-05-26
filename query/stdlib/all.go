package stdlib

import (
	"github.com/ihaiker/ngx/v2/query/methods"
)

var _methods = map[string]interface{}{
	"select": fnSelect,
}

func Defaults() *methods.FunctionManager {
	defs := methods.New()
	for name, function := range _methods {
		_ = defs.Add(name, function)
	}
	return defs
}
