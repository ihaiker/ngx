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
)

// include后置处理
func Include(merge bool, search Walk) Hook {
	return func(conf *config.Configuration) (err error) {
		err = walk(conf, merge, search)
		return
	}
}

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
			err = fmt.Errorf("The include file not found: %s", strings.Join(args, ","))
		}
		return
	}
}

func walk(conf *config.Configuration, merge bool, search Walk) error {
	if conf.Body != nil {
		for i := 0; i < len(conf.Body); i++ {
			if len(conf.Body) == i {
				break
			}

			item := conf.Body[i]

			if item.Name == "include" {
				ds, err := includes(search, merge, item)
				if err != nil {
					return err
				}
				if merge {
					conf.Body = append(conf.Body[:i], append(ds, conf.Body[i+1:]...)...)
					i += len(ds) - 1
				} else {
					item.Body = append(item.Body, ds...)
				}
			} else {
				subConf := &config.Configuration{Source: "", Body: item.Body}
				if err := walk(subConf, merge, search); err != nil {
					return err
				} else {
					item.Body = subConf.Body
				}
			}
		}
	}
	return nil
}

//处理include文件
func includes(search Walk, merge bool, node *config.Directive) (config.Directives, error) {
	files, err := search(node.Args...)
	if err != nil {
		return nil, err
	}
	var doc *config.Configuration
	includeFiles := config.Directives{}
	for _, file := range files {
		if doc, err = config.ParseIO(file.Stream); err != nil {
			return nil, err
		}

		if err = Include(merge, search)(doc); err != nil {
			return nil, err
		}

		if merge {
			includeFiles = append(includeFiles, doc.Body...)
		} else {
			includeFiles = append(includeFiles, &config.Directive{
				Line:    0,
				Virtual: config.Include,
				Name:    "file",
				Args:    []string{file.Name},
				Body:    doc.Body,
			})
		}
	}
	return includeFiles, nil
}
