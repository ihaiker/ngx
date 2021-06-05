package hooks

import "github.com/ihaiker/ngx/v2/config"

type (
	Hooker interface {
		Execute(conf *config.Configuration) (err error)
	}
	Hookers []Hooker
)

func New(hookers ...Hooker) Hookers {
	hs := Hookers{}
	for _, hooker := range hookers {
		hs = append(hs, hooker)
	}
	return hs
}

func (this Hookers) Hooker(hooker Hooker) Hookers {
	this = append(this, hooker)
	return this
}

func (h Hookers) Execute(conf *config.Configuration) (err error) {
	for _, hooker := range h {
		if err = hooker.Execute(conf); err != nil {
			return
		}
	}
	return
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
