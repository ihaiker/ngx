package hooks

import (
	"fmt"
	"github.com/ihaiker/ngx/v2/config"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	File struct {
		Name   string
		Stream io.ReadCloser
	}
	Walk func(args ...string) (files []*File, err error)

	// include后置处理
	IncludeHooker struct {
		Merge  bool
		Search Walk
	}
)

// 本地文件搜索，root为搜索的根路径
func WalkFiles(root string) Walk {
	return func(args ...string) (files []*File, err error) {
		files = make([]*File, 0)
		var matches []string
		for _, arg := range args {
			if matches, err = filepath.Glob(filepath.Join(root, arg)); err != nil {
				return
			}
			for _, match := range matches {
				if f, e := os.Open(match); e != nil {
					err = e
					return
				} else {
					files = append(files, &File{
						Name:   match,
						Stream: f,
					})
				}
			}
		}
		if len(files) == 0 {
			for _, arg := range args {
				if strings.Contains(arg, "*") {
					return
				}
			}
			err = fmt.Errorf("the include file not found: %s", strings.Join(args, ","))
		}
		return
	}
}

func (this *IncludeHooker) Execute(conf *config.Directive, _ Next) (config.Directives, config.Directives, error) {
	if this.Search == nil {
		dir, _ := os.Getwd()
		this.Search = WalkFiles(dir)
	}

	//搜索到的文件
	files, err := this.Search(conf.Args...)
	if err != nil {
		return nil, nil, err
	}

	var doc *config.Configuration
	items := config.Directives{}
	for _, file := range files {
		if doc, err = config.ParseIO(file.Stream); err != nil {
			return nil, nil, err
		}
		if this.Merge {
			items = append(items, doc.Body...)
		} else {
			items = append(items, &config.Directive{
				Line:    conf.Line,
				Virtual: config.Include,
				Name:    "file",
				Args:    []string{file.Name},
				Body:    doc.Body,
			})
		}
	}
	if this.Merge {
		return items, nil, nil
	}
	return nil, items, nil
}
