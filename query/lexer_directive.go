package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ihaiker/ngx/v2/config"
	"regexp"
	"strings"
)

type directive struct {
	Pos   lexer.Position
	Regex *string `"."[ ("*") |  @Regex`
	Name  *string `| (@Ident|("[" @String "]")) ]`
	Args  *Args   `[("("@@")")]`
	Index *Index  `["["@@"]"]`
}
type directives []directive

//指令检索
func (expr directives) call(items config.Directives) (matched config.Directives, err error) {
	dir := expr[0]
	if matched, err = dir.call(items); err != nil {
		return
	}
	if len(expr[1:]) == 0 {
		return
	}
	return expr[1:].call(matched)
}

func (this *directive) call(items config.Directives) (matched config.Directives, err error) {

	//匹配 "."
	if this.Name == nil && this.Regex == nil {
		matched = items
	} else if this.Name != nil { //名称匹配子指令名称
		for _, item := range items {
			for _, body := range item.Body {
				if body.Name == *this.Name {
					//不检查参数匹配或者参数匹配正确
					if this.Args == nil || this.Args.match(body.Args) {
						matched = append(matched, body)
					}
				}
			}
		}
	} else if this.Regex != nil { //正则匹配紫属性
		regex := (*this.Regex)[1 : len(*this.Regex)-1]
		for _, item := range items {
			for _, body := range item.Body {
				if match, _ := regexp.MatchString(regex, body.Name); match {
					//不检查参数匹配或者参数匹配正确
					if this.Args == nil || this.Args.match(body.Args) {
						matched = append(matched, body)
					}
				}
			}
		}
	}

	//没有匹配到相应内容
	if len(matched) == 0 {
		return
	}

	if this.Index != nil {
		if this.Index.Split == "" {
			index := *this.Index.Start
			if index < 0 {
				index = len(matched) + index
			}
			if index > len(matched)-1 || index < 0 {
				return nil, fmt.Errorf("index out of: %s\n%s^",
					this.Pos.Filename, strings.Repeat(" ", 14+this.Pos.Offset+this.Index.Pos.Offset))
			}
			matched = matched[index : index+1]
		} else {
			start := 0
			end := len(matched)
			if this.Index.Start != nil {
				start = *this.Index.Start
			}
			if this.Index.End != nil {
				end = *this.Index.End
			}
			if end < 0 {
				end = len(matched) + end
			}
			if !(0 <= start && start < len(matched)) {
				return nil, fmt.Errorf("index out of: %s\n%s^",
					this.Pos.Filename, strings.Repeat(" ", 14+this.Pos.Offset+this.Index.Pos.Offset))
			}
			if !(start <= end && end <= len(matched)) {
				return nil, fmt.Errorf("index out of: %s\n%s^",
					this.Pos.Filename, strings.Repeat(" ", 14+this.Pos.Offset+this.Index.Pos.Offset))
			}
			matched = matched[start:end]
		}
	}
	return
}

type Index struct {
	Pos   lexer.Position
	Start *int   `[@Number]`
	Split string `[@(":")]`
	End   *int   `[@Number]`
}

type Args struct {
	Pos      lexer.Position
	Left     Arg    `@@`
	Operator string `[Space (@("&""&") | @("|""|")) Space`
	Right    *Arg   `@@]`
}

func (this *Args) match(args []string) bool {
	left := this.Left.match(args)
	if this.Right != nil {
		if left && this.Operator == "||" {
			return left
		}
		right := this.Right.match(args)
		return left && right
	}
	return left
}

type Arg struct {
	Pos        lexer.Position
	Comparison string  `(([@("!" | "@" | "^" | "$")]`
	Value      *string `@String)`
	Regex      *string `|@Regex`
	Group      *Args   `|("(" @@ ")"))`
}

func (this *Arg) matchComparison(comparison string, items []string) bool {

	if this.Regex != nil {
		for _, item := range items {
			if match, _ := regexp.MatchString(*this.Regex, item); match {
				return true
			}
		}
		return false
	}

	switch this.Comparison {
	case "!":
		return !this.matchComparison("", items)
	case "@":
		for _, item := range items {
			if strings.Contains(*this.Value, item) {
				return true
			}
		}
	case "^":
		for _, item := range items {
			if strings.HasSuffix(*this.Value, item) {
				return true
			}
		}
	case "$":
		for _, item := range items {
			if strings.HasSuffix(*this.Value, item) {
				return true
			}
		}
	case "":
		for _, item := range items {
			if *this.Value == item {
				return true
			}
		}
	}
	return false
}

func (this *Arg) match(items []string) bool {
	if this.Group != nil {
		return this.Group.match(items)
	}
	return this.matchComparison(this.Comparison, items)
}
