package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
)

type TemplateHooker struct {
	vars      *Variables
	templates map[string]*config.Directive
}

func (this *TemplateHooker) SetVariables(variables *Variables) {
	this.vars = variables
}

func (this *TemplateHooker) Execute(item *config.Directive, _ Next) (config.Directives, config.Directives, error) {
	if item.Name == "@template" {
		return this.template(item)
	} else if item.Name == "@merge" {
		return this.merge(item)
	} else {
		return this.include(item)
	}
}

func (this *TemplateHooker) template(item *config.Directive) (current config.Directives, children config.Directives, err error) {
	if this.templates == nil {
		this.templates = map[string]*config.Directive{}
	}
	if len(item.Args) == 0 {
		err = fmt.Errorf("not found tempalte name at line %d", item.Line)
		return
	}
	name := item.Args[0]
	//当前不是个单节点
	if len(item.Args) > 1 {
		item.Name = getArys(item.Args, 1)
		item.Args = sliceArgs(item.Args, 2)
		this.templates[name] = item
	} else {
		this.templates[name] = &config.Directive{
			Body: item.Body,
		}
	}
	return
}

func (this *TemplateHooker) include(item *config.Directive) (current config.Directives, children config.Directives, err error) {
	name := getArys(item.Args, 0)
	if template, has := this.templates[name]; has {
		if template.Name == "" {
			current = template.Body.Clone()
		} else {
			current = config.Directives{template.Clone()}
		}
	} else {
		err = fmt.Errorf("not found tempalte name at line %d", item.Line)
	}
	return
}

func (this *TemplateHooker) merge(item *config.Directive) (current config.Directives, children config.Directives, err error) {
	name := getArys(item.Args, 0)
	if template, has := this.templates[name]; !has {
		err = fmt.Errorf("not found tempalte name at line %d", item.Line)
	} else {
		current = this.mergeIt(item.Clone(), template.Clone())
	}
	return
}

func (t *TemplateHooker) mergeIt(item *config.Directive, temp *config.Directive) (merged config.Directives) {
	if len(item.Args) > 1 || len(temp.Args) > 0 {
		t.mergeSingle(item, temp)
		return config.Directives{item}
	}
	return t.mergeBody(item.Body, temp.Body)
}

func (t *TemplateHooker) mergeSingle(item *config.Directive, temp *config.Directive) {
	//改变本身的指令名
	item.Name = getArys(item.Args, 1)
	if item.Name == "" {
		item.Name = temp.Name
	}

	item.Args = sliceArgs(item.Args, 2)
	if len(item.Args) == 0 && len(temp.Args) > 0 {
		item.Args = sliceArgs(temp.Args, 0)
	}
	item.Body = t.mergeBody(item.Body, temp.Body)
}

func (t *TemplateHooker) mergeBody(writes config.Directives, temps config.Directives) config.Directives {
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
