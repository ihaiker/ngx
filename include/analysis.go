package include

import (
	"github.com/ihaiker/ngx/config"
)

type (
	WalkSearch func(args ...string) (files []string, err error)
	IsWalk     func(directive *config.Directive) bool
)

func Walk(cfg *config.Configuration, isWalk IsWalk, search WalkSearch, opt *config.Options) error {
	if isWalk(cfg) {
		if err := includes(search, cfg, opt); err != nil {
			return err
		}
	}
	if cfg.Body != nil {
		for _, d := range cfg.Body {
			if err := Walk(d, isWalk, search, opt); err != nil {
				return err
			}
		}
	}
	return nil
}

// 处理 include 文件
func includes(search WalkSearch, node *config.Directive, opt *config.Options) error {
	files, err := search(node.Args...)
	if err != nil {
		return err
	}
	for _, file := range files {
		fileDirective := &config.Directive{
			Virtual: config.Include, Name: "file", Args: []string{file},
		}
		if doc, err := config.Parse(file, opt); err != nil {
			return err
		} else {
			fileDirective.Body = doc.Body
		}
		node.Body = append(node.Body, fileDirective)
	}
	return nil
}
