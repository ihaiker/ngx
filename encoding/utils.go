package encoding

import (
	"regexp"
	"strings"
)

func split2(value, sep string) (string, string) {
	ary := strings.SplitN(value, sep, 2)
	if len(ary) < 2 {
		return value, ""
	}
	return ary[0], ary[1]
}

func compileSplit2(str, sep string) (a, b string) {
	outs := regexp.MustCompile(sep).Split(str, 2)
	a = outs[0]
	if len(outs) > 1 {
		b = outs[1]
	}
	return
}

func index(args []string, index int) string {
	if index > len(args)-1 {
		return ""
	}
	return args[index]
}
