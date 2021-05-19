package query

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ihaiker/ngx/v2/config"
	"regexp"
	"strings"
)

type Directive struct {
	Pos   lexer.Position
	Regex *string `"."[ ("*") |  @Regex`
	Name  *string `| (@Ident|("[" @String "]")) ]`
	Args  *Args   `[("("@@")")]`
	Index *Index  `["["@@"]"]`
}

func (d *Directive) match(item *config.Directive) bool {
	if d.Name != nil && *d.Name != item.Name {
		return false
	}
	if d.Regex != nil {
		regex := (*d.Regex)[1 : len(*d.Regex)-1]
		if match, _ := regexp.MatchString(regex, item.Name); !match {
			return false
		}
	}
	if d.Args != nil {
		return d.Args.match(item.Args)
	}
	return true
}

func (this *Directive) Select(items config.Directives) (matched config.Directives, err error) {

	for _, item := range items {
		if this.match(item) {
			matched = append(matched, item)
		}
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
	Operator string `[(([Whitespace] @("&""&") [Whitespace]) | ([Whitespace] @("|""|") [Whitespace]))`
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
	Comparison string  `(([@("!" | "@" | "^" | "$" | "&")]`
	Value      *string `@String)`
	Group      *Args   `|("(" @@ ")"))`
}

func (this *Arg) matchComparison(comparison string, items []string) bool {
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
	case "&":
		for _, item := range items {
			if match, err := regexp.MatchString(*this.Value, item); err == nil && match {
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
