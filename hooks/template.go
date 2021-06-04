package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"regexp"
)

var (
	templateCompile = regexp.MustCompile(`@template\(([a-zA-Z][a-zA-Z0-9]*)\)`)
	iormCompile     = regexp.MustCompile(`@(merge|include)\(([a-zA-Z][a-zA-Z0-9]*)\)`)
)

func Template(conf *config.Configuration) (err error) {
	//查找模板文件
	templates := searchTemplate(conf)

	//先行处理模板中的合并指令
	if err = searchMerge(templates, &config.Configuration{Body: templates}); err != nil {
		return
	}

	//执行合并模板
	err = searchMerge(templates, conf)
	return
}

func getTemplate(templates config.Directives, name string) *config.Directive {
	for _, directive := range templates {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func isTemplate(value string) (matched bool, name string) {
	if matched = templateCompile.MatchString(value); !matched {
		return
	}
	matches := templateCompile.FindStringSubmatch(value)
	name = matches[1]
	return
}

func isIncludeOrMerge(value string) (include bool, merge bool, name string) {
	if !iormCompile.MatchString(value) {
		return
	}
	matches := iormCompile.FindStringSubmatch(value)
	include = matches[1] == "include"
	merge = !include
	name = matches[2]
	return
}

//搜索模板
func searchTemplate(conf *config.Configuration) config.Directives {
	var templates config.Directives
	for idx := 0; ; idx++ {
		if idx == len(conf.Body) {
			break
		}

		item := conf.Body[idx]
		if len(item.Body) > 0 {
			subConf := &config.Configuration{Body: item.Body}
			if temps := searchTemplate(subConf); temps != nil {
				item.Body = subConf.Body
				templates = append(templates, temps...)
			}
		}

		if matched, name := isTemplate(item.Name); matched {
			//remove template directive
			conf.Body = append(conf.Body[:idx], conf.Body[idx+1:]...)
			item.Name = name
			templates = append(templates, item)
		}
	}
	return templates
}

func searchMerge(templates config.Directives, conf *config.Configuration) error {
	for idx := 0; ; idx++ {
		if idx == len(conf.Body) {
			break
		}

		//先把下层的合并了
		item := conf.Body[idx]
		if len(item.Body) > 0 {
			subConf := &config.Configuration{Body: item.Body}
			if err := searchMerge(templates, subConf); err != nil {
				return err
			}
			item.Body = subConf.Body
		}

		item = conf.Body[idx]
		if isInclude, isMerge, name := isIncludeOrMerge(item.Name); isInclude || isMerge {
			temp := getTemplate(templates, name)
			if temp == nil {
				return fmt.Errorf("The template %s not found at line %d", item.Args[0], item.Line)
			}

			if isMerge {
				if items, single := merge(item, temp); single {
					item.Name = items[0].Name
					item.Args = items[0].Args
					item.Body = items[0].Body
				} else { //下级指令合并
					conf.Body = append(conf.Body[:idx], append(items, conf.Body[idx+1:]...)...)
					idx += len(items) - 1
				}
			} else if isInclude {
				conf.Body = append(conf.Body[:idx], append(temp.Body, conf.Body[idx+1:]...)...)
				idx += len(temp.Body) - 1
			}
		}
	}
	return nil
}

func merge(item *config.Directive, temp *config.Directive) (merged config.Directives, single bool) {
	if len(item.Args) > 0 || len(temp.Args) > 0 {
		mergeSingle(item, temp)
		return config.Directives{item}, true
	}
	return mergeBody(item.Body, temp.Body), false
}

func getArys(args []string, index int) string {
	if len(args) > index {
		return args[index]
	}
	return ""
}
func sliceArgs(args []string, start int) []string {
	if len(args) > start {
		return args[start:]
	}
	return nil
}

func mergeSingle(item *config.Directive, temp *config.Directive) {
	item.Name = getArys(item.Args, 0)
	if item.Name == "" {
		item.Name = getArys(temp.Args, 0)
	}

	item.Args = sliceArgs(item.Args, 1)
	if len(item.Args) == 0 && len(temp.Args) > 0 {
		item.Args = sliceArgs(temp.Args, 1)
	}
	item.Body = mergeBody(item.Body, temp.Body)
}

func mergeBody(writes config.Directives, temps config.Directives) config.Directives {
	outs := temps.Clone()
	names := writes.Names()

	for _, name := range names {
		writeItems := writes.Gets(name)
		outItems := outs.Gets(name)

		//模板中不包含，直接添加
		if len(outItems) == 0 {
			outs = append(outs, writeItems...)
			continue
		}

		//模板中包含，且都唯一，直接替换
		if len(outItems) == 1 && len(writeItems) == 1 {
			outItems[0].Line = writeItems[0].Line
			outItems[0].Name = writeItems[0].Name
			outItems[0].Virtual = writeItems[0].Virtual
			outItems[0].Args = writeItems[0].Args
			outItems[0].Body = writeItems[0].Body
			continue
		}

		//模板中包含，都不唯一，单个寻找，匹配规则为 directive.arg1 相等
		for _, writeItem := range writeItems {

			var searchOutItem *config.Directive
			for _, outItem := range outItems {
				if writeItem.Name == outItem.Name &&
					getArys(writeItem.Args, 0) == getArys(outItem.Args, 0) {
					searchOutItem = outItem
					break
				}
			}

			//没有找到
			if searchOutItem == nil {
				outs = append(outs, writeItem)
				continue
			}

			//找到了，
			searchOutItem.Args = writeItem.Args
			searchOutItem.Body = writeItem.Body
		}
	}
	return outs
}
