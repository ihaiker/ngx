package config

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

var (
	colorComment = color.New(color.FgYellow, color.Concealed)
	colorKey     = color.New(color.FgBlue)
	colorValue   = color.New(color.FgHiGreen)
)

type (
	Virtual string

	Directive struct {
		Line    int        `json:"line"`
		Virtual Virtual    `json:"virtual,omitempty"`
		Name    string     `json:"name"`
		Args    []string   `json:"args,omitempty"`
		Body    Directives `json:"body,omitempty"`
	}
	Directives    []*Directive
	Configuration struct {
		Source string     `json:"source"`
		Body   Directives `json:"body"`
	}
)

var Include Virtual = "include"

func New(name string, args ...string) *Directive {
	return &Directive{Name: name, Args: args}
}
func Body(name string, body ...*Directive) *Directive {
	return &Directive{Name: name, Body: body}
}
func Config(body ...*Directive) *Configuration {
	return &Configuration{
		Source: "",
		Body:   body,
	}
}

func (d *Directive) String() string {
	return d.Pretty(0)
}

func (d *Directive) noBody() bool {
	if len(d.Body) == 0 {
		return true
	} else {
		for _, body := range d.Body {
			if body.Virtual == "" {
				return false
			}
		}
		return true
	}
}

func (d *Directive) AddBody(name string, args ...string) *Directive {
	body := New(name, args...)
	d.AddBodyDirective(body)
	return body
}

func (d *Directive) AddArgs(args ...string) *Directive {
	d.Args = append(d.Args, args...)
	return d
}

func (d *Directive) AddBodyDirective(directives ...*Directive) {
	if d.Body == nil {
		d.Body = make([]*Directive, 0)
	}
	d.Body = append(d.Body, directives...)
}

func (d Directive) Clone() *Directive {
	nd := &Directive{
		Line:    d.Line,
		Virtual: d.Virtual,
		Name:    d.Name,
	}
	for _, arg := range d.Args {
		nd.Args = append(nd.Args, arg)
	}
	if d.Body != nil && len(d.Body) > 0 {
		nd.Body = d.Body.Clone()
	}
	return nd
}

func (d *Directive) Pretty(prefix int, colorize ...bool) string {
	//fix for delete directive.
	if d.Name == "" {
		return ""
	}

	colored := len(colorize) == 1 && colorize[0]
	prefixString := strings.Repeat(" ", prefix*4)
	if d.Name == "#" {
		if colored {
			return colorComment.Sprintf("%s# %s", prefixString, d.Args[0])
		}
		return fmt.Sprintf("%s# %s", prefixString, d.Args[0])
	} else if d.Virtual != "" {
		return ""
	} else {
		out := bytes.NewBufferString(prefixString)
		if colored {
			out.WriteString(colorKey.Sprint(d.Name))
		} else {
			out.WriteString(d.Name)
		}
		splitLine := (len(d.Name) + len(strings.Join(d.Args, " "))) > 80

		for i, arg := range d.Args {
			out.WriteByte(' ')
			if i != 0 && splitLine {
				out.WriteRune('\n')
				out.WriteString(strings.Repeat(" ", len(d.Name)+prefix*4))
			}

			if strings.ContainsAny(arg, "\r\n") {
				if colored {
					arg = colorValue.Sprint(arg)
				}
				out.WriteRune('`')
				out.WriteString(arg)
				out.WriteRune('`')
			} else if strings.ContainsRune(arg, '"') {
				if colored {
					arg = colorValue.Sprint(arg)
				}
				out.WriteRune('\'')
				out.WriteString(arg)
				out.WriteRune('\'')
			} else if arg == "" || strings.ContainsRune(arg, '\'') ||
				strings.ContainsAny(arg, "{};") || strings.ContainsAny(arg, "\t ") {
				if colored {
					arg = colorValue.Sprint(strconv.Quote(arg))
					out.WriteString(arg)
				} else {
					out.WriteString(strconv.Quote(arg))
				}
			} else {
				if colored {
					arg = colorValue.Sprint(arg)
				}
				out.WriteString(arg)
			}
		}

		if d.noBody() {
			out.WriteString(";")
		} else {
			out.WriteString(" {")
			for _, body := range d.Body {
				bodyItem := body.Pretty(prefix+1, colorize...)
				if bodyItem != "" {
					out.WriteString("\n")
					out.WriteString(bodyItem)
				}
			}
			out.WriteString(fmt.Sprintf("\n%s}", prefixString))
		}
		return out.String()
	}
}

func (cfg *Configuration) Pretty(colorize ...bool) string {
	out := bytes.NewBufferString("")
	for i, item := range cfg.Body {
		if i != 0 {
			out.WriteByte('\n')
		}
		itemString := item.Pretty(0, colorize...)
		_, _ = out.WriteString(itemString)
	}
	return out.String()
}

func (ds *Directives) Get(name string) *Directive {
	var cur *Directive
	for _, d := range *ds {
		if d.Name == name {
			cur = d
		}
	}
	return cur
}

func (ds *Directives) Gets(name string) (ret []*Directive) {
	ret = make([]*Directive, 0)
	for _, d := range *ds {
		if d.Name == name {
			ret = append(ret, d)
		}
	}
	return
}

func (ds Directives) Index(idx int) *Directive {
	if idx < 0 || idx > len(ds)-1 {
		return nil
	}
	return ds[idx]
}

func (ds *Directives) Insert(d *Directive, idx int) {
	*ds = append((*ds)[:idx], append([]*Directive{d}, (*ds)[idx:]...)...)
}

func (ds *Directives) Append(d *Directive) {
	*ds = append(*ds, d)
}

func (ds Directives) Names() []string {
	names := make([]string, 0)
DLOOP:
	for _, d := range ds {
		for _, name := range names {
			if d.Name == name {
				continue DLOOP
			}
		}
		names = append(names, d.Name)
	}
	return names
}

func (ds Directives) Clone() Directives {
	outs := make([]*Directive, len(ds))
	for i, d := range ds {
		outs[i] = d.Clone()
	}
	return outs
}
