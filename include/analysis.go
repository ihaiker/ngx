package include

import (
	"github.com/ihaiker/ngx/config"
)

type (
	WalkSearch func(args ...string) (files []string, err error)
	IsWalk     func(directive *config.Directive) bool
)

func Walk(cfg *config.Configuration, walk IsWalk, search WalkSearch, opt *config.Options) error {
	if cfg.Body != nil {

		for i := 0; ; i++ {
			if len(cfg.Body) == i {
				break
			}
			item := cfg.Body[i]

			if walk(item) {
				ds, err := includes(search, item, opt)
				if err != nil {
					return err
				}
				if opt.MergeInclude {
					cfg.Body = append(cfg.Body[:i], append(ds, cfg.Body[i+1:]...)...)
					i += len(ds) - 1
				} else {
					item.Body = append(item.Body, ds...)
				}
			} else {
				if err := Walk(item, walk, search, opt); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//includes  处理include文件
func includes(search WalkSearch, node *config.Directive, opt *config.Options) (config.Directives, error) {
	files, err := search(node.Args...)
	if err != nil {
		return nil, err
	}
	includeFiles := config.Directives{}
	for _, file := range files {
		doc, err := config.Parse(file, opt)
		if err != nil {
			return nil, err
		}

		if opt.MergeInclude {
			includeFiles = append(includeFiles, doc.Body...)
		} else {
			doc.Virtual = config.Include
			doc.Name = "file"
			doc.Args = []string{file}
			includeFiles = append(includeFiles, doc)
		}
	}
	return includeFiles, nil
}
