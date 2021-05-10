package config

import "regexp"

type Filter func(current, previous string) bool

func (self Filter) And(cf ...Filter) Filter {
	return func(current, previous string) bool {
		if !self(current, previous) {
			return false
		}
		for _, filter := range cf {
			if !filter(current, previous) {
				return false
			}
		}
		return true
	}
}

func (self Filter) Or(cf ...Filter) Filter {
	return func(current, previous string) bool {
		out := self(current, previous)
		for _, filter := range cf {
			out = out || filter(current, previous)
		}
		return out
	}
}

var (
	vailCharRegexp        = regexp.MustCompile("\\S")
	ValidChars     Filter = func(current, previous string) bool {
		return vailCharRegexp.MatchString(current)
	}

	In = func(chars ...string) Filter {
		return func(current, previous string) bool {
			for _, char := range chars {
				if char == current {
					return true
				}
			}
			return false
		}
	}

	Not = func(cf Filter) Filter {
		return func(current, previous string) bool {
			return !cf(current, previous)
		}
	}
)
