package encoding

import (
	"encoding/json"
	"github.com/ihaiker/ngx/v2/config"
)

type _json map[string]interface{}

func Json(conf *config.Configuration) ([]byte, error) {
	values := _loop_json(conf.Body)
	values["@source"] = conf.Source
	return json.Marshal(values)
}
func JsonIndent(conf *config.Configuration, prefix, indent string) ([]byte, error) {
	values := _loop_json(conf.Body)
	values["@source"] = conf.Source
	return json.MarshalIndent(values, prefix, indent)
}

func _item_value(item *config.Directive) interface{} {
	if len(item.Body) == 0 {
		if len(item.Args) == 0 {
			return true
		} else if len(item.Args) == 1 {
			return item.Args[0]
		} else {
			return item.Args
		}
	} else {
		value := _json{
			//"@line": item.Line,
		}
		if len(item.Args) > 0 {
			value["@args"] = item.Args
		}
		//if item.Virtual != "" {
		//	value["@virtual"] = item.Virtual
		//}
		if len(item.Body) > 0 {
			childData := _loop_json(item.Body)
			for name, child := range childData {
				value[name] = child
			}
		}
		return value
	}
}

func _loop_json(items config.Directives) _json {
	values := _json{}
	names := items.Names()
	for _, name := range names {
		if name == "#" {
			continue
		}
		subs := items.Gets(name)
		if len(subs) == 1 {
			item := subs[0]
			value := _item_value(item)
			values[item.Name] = value
		} else {
			jsons := make([]interface{}, len(subs))
			for idx, sub := range subs {
				jsons[idx] = _item_value(sub)
			}
			values[name] = jsons
		}
	}
	return values
}
