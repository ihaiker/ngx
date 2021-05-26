package stdlib

import (
	"github.com/ihaiker/ngx/v2/config"
	"regexp"
	"strings"
)

//select(.http.server.server_name, 'equal', 'domain')
func fnSelect(name, operator, value interface{}) (bool, error) {
	if items, match := name.(config.Directives); match {
		return fnSelectItems(items, operator.(string), value.(string))
	}
	if item, match := name.(*config.Directive); match {
		return fnSelectItem(item, operator.(string), value.(string))
	}
	if operator == "equal" {
		return name == value, nil
	} else {
		return name != value, nil
	}
}

func fnSelectItems(items config.Directives, operator, value string) (bool, error) {
	for _, item := range items {
		if b, err := fnSelectItem(item, operator, value); err != nil || b {
			return b, err
		}
	}
	return false, nil
}

func fnSelectItem(item *config.Directive, operator, value string) (bool, error) {
	if operator == "not" {
		matched, err := fnSelectItem(item, "equal", value)
		return !matched, err
	}

	for _, arg := range item.Args {
		switch operator {
		case "equal":
			if arg == value {
				return true, nil
			}
		case "startWith", "hasPrefix":
			if strings.HasPrefix(arg, value) {
				return true, nil
			}
		case "endWith", "hasSuffix":
			if strings.HasSuffix(arg, value) {
				return true, nil
			}
		case "regex":
			if matched, _ := regexp.MatchString(value, arg); matched {
				return true, nil
			}
		case "contains":
			if strings.Contains(arg, value) {
				return true, nil
			}
		}
	}
	return false, nil
}
