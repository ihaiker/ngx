# 语法介绍

```shell script
cat nginx.conf | ngc '.http.server | 
  select(.server_name equal "api.aginx.io" and equal "v2.aginx.io") |
  args() | args_length() | length() | 

  or(select(.server_name, equal, "v2.aginx.io"), select(.server_name, equal, "v2.aginx.io"))'

cat nginx.conf | ngc '.http.server | select(args_length(.listen), "equal", args_length(.server_name))'


https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/
https://github.com/antlr/antlr4
```


```golang 
func (this *Function) executorFn(function interface{}, item *config.Directive, manager *FunctionManager) (config.Directives, error) {
	outValues, err := this.realExecutor(function, item, manager)
	if err != nil {
		return nil, err
	}

	switch len(outValues) {
	case 2:
		if err, match := outValues[1].Interface().(error); !match {
			return nil, fmt.Errorf("invalid out value type of the function %s", this.Name)
		} else if err != nil {
			return nil, err
		}
		fallthrough
	case 1:
		if outValues[0].Type().Kind() == reflect.Bool {
			if outValues[0].Bool() {
				return config.Directives{item}, nil
			}
		} else if outValues[0].Type().String() == reflect.ValueOf(config.Directives{}).String() {
			return outValues[0].Interface().(config.Directives), nil
		} else {
			if conf, err := encoding.MarshalOptions(outValues[0].Interface(), *encoding.DefaultOptions()); err != nil {
				return nil, err
			} else {
				return conf.Body, nil
			}
		}
		return nil, nil
	}
	return nil, fmt.Errorf("invalid %s function return type: %s", this.Name, this.Pos.GoString())
}
```
