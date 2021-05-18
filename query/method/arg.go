package method

import "github.com/ihaiker/ngx/v2/config"

func arg(items config.Directives, index int, value string) config.Directives {
	newItems := config.Directives{}
	for _, item := range items {
		if len(item.Args) > index {
			if item.Args[index] == value {
				newItems = append(newItems, item)
			}
		}
	}
	return newItems
}
