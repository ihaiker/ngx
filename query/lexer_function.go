package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ihaiker/ngx/v2/config"
)

// args('')
type Function struct {
	Pos  lexer.Position
	Name string    `@Ident`
	Args []FuncArg `"(" [Whitespace] @@ [( [Whitespace] "," [Whitespace] @@ )+] [Whitespace] ")"`
}

type FuncArg struct {
	Pos       lexer.Position
	Directive Directives `(@@+`
	Index     *int       `|@Number`
	Value     *string    `|@String`
	Function  *Function  `|@@`
	Arrays    []FuncArg  `|("[" [Whitespace] @@ ([Whitespace]","[Whitespace] @@)* [Whitespace] "]") )`
}

type ExecutorFunction interface {
	Select(args []FuncArg, items config.Directives, fns map[string]ExecutorFunction) (config.Directives, error)
}

func (this *Function) Select(items config.Directives, fns map[string]ExecutorFunction) (config.Directives, error) {
	fn, has := fns[this.Name]
	if !has {
		return nil, fmt.Errorf("function not found: %s", this.Name)
	}
	matched, err := fn.Select(this.Args, items, fns)
	if err != nil {
		err = fmt.Errorf("error function %s: %s", this.Name, err.Error())
	}
	return matched, err
}

func ExecutorFunctions(fns map[string]ExecutorFunction) map[string]ExecutorFunction {
	defaults := map[string]ExecutorFunction{
		"select": new(selectFunc),
	}
	for name, value := range fns {
		defaults[name] = value
	}
	return defaults
}
