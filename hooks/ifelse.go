package hooks

import (
	"github.com/ihaiker/ngx/v2/config"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type IfElseHooker struct {
	vars *Variables
}

func (this *IfElseHooker) SetVariables(variables *Variables) {
	this.vars = variables
}

func (this IfElseHooker) is() func(*config.Directive) bool {
	return func(item *config.Directive) bool {
		return "@else" == item.Name || "@elseif" == item.Name
	}
}

func (this IfElseHooker) Execute(item *config.Directive, next Next) (current config.Directives, children config.Directives, err error) {
	var elseItem *config.Directive
	elseIfs := config.Directives{}

	for i := 0; ; i++ {
		if item := next(this.is()); item == nil {
			break
		} else if item.Name == "@else" {
			elseItem = item
		} else {
			elseIfs = append(elseIfs, item)
		}
	}

	if this.check(item) {
		current = item.Body
	} else {
		for _, elseIf := range elseIfs {
			if this.check(elseIf) {
				current = elseIf.Body
			}
		}
	}
	if current == nil && elseItem != nil {
		current = elseItem.Body
	}
	return
}

func (this IfElseHooker) check(item *config.Directive) bool {
	if len(item.Args) == 0 {
		return false
	}
	val, err := this.vars.Get(item.Args[0])
	if err != nil {
		return false
	}
	vv := reflect.ValueOf(val)
	if vv.IsValid() && vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch len(item.Args) {
	case 1:
		{
			if !vv.IsValid() {
				return false
			}
			if vv.Kind() == reflect.Bool {
				return vv.Bool()
			}
			return vv.String() != "false" && vv.String() != ""
		}
	case 2:
		{
			if vv.IsValid() && vv.Kind() == reflect.Bool &&
				!vv.Bool() && item.Args[1] == "notTrue" {
				return true
			}

			return (!vv.IsValid() && item.Args[1] == "notNil") ||
				(vv.IsValid() && item.Args[1] == "isNil")
		}
	default:
		{
			if !vv.IsValid() {
				return false
			}

			if vv.Kind() == reflect.String {
				return this.checkString(vv.String(), item.Args[1], item.Args[2])

			} else if _is(vv.Type(), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
				return this.checkInt(float64(vv.Int()), item.Args[1], item.Args[2:]...)

			} else if _is(vv.Type(), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64) {
				return this.checkInt(float64(vv.Uint()), item.Args[1], item.Args[2:]...)

			} else if _is(vv.Type(), reflect.Float32, reflect.Float64) {
				return this.checkInt(vv.Float(), item.Args[1], item.Args[2:]...)
			}
		}
	}
	return false
}

func (this *IfElseHooker) checkInt(value float64, compare string, numbers ...string) bool {
	switch compare {
	case ">":
		i, err := strconv.ParseFloat(getArys(numbers, 0), 64)
		if err != nil {
			return false
		}
		return i > value
	case ">=":
		i, err := strconv.ParseFloat(getArys(numbers, 0), 64)
		if err != nil {
			return false
		}
		return i >= value
	case "==":
		i, err := strconv.ParseFloat(getArys(numbers, 0), 64)
		if err != nil {
			return false
		}
		return i == value
	case "<":
		i, err := strconv.ParseFloat(getArys(numbers, 0), 64)
		if err != nil {
			return false
		}
		return i < value
	case "<=":
		i, err := strconv.ParseFloat(getArys(numbers, 0), 64)
		if err != nil {
			return false
		}
		return i <= value
	case "in": // .s1 in [1..2)
		r := regexp.MustCompile(`(\[|\()([-]?(\d*\.)?\d+)\.\.([-]?(\d*\.)?\d+)(\]|\))`)
		gs := r.FindStringSubmatch(getArys(numbers, 0))
		startInclude := gs[1] == "["
		start, err := strconv.ParseFloat(gs[2], 64)
		if err != nil {
			return false
		}
		endInclude := gs[1] == "]"
		end, err := strconv.ParseFloat(gs[4], 64)
		if err != nil {
			return false
		}
		if startInclude && endInclude {
			return start <= value && value <= end
		} else if startInclude {
			return start <= value && value < end
		} else if endInclude {
			return start < value && value <= end
		} else {
			return start < value && value < end
		}
	case "notIn":
		return this.checkInt(value, "in", numbers...)
	}
	return false
}

func (this *IfElseHooker) checkString(value, compare, pattern string) bool {
	switch compare {
	case "match":
		ok, _ := regexp.MatchString(pattern, value)
		return ok
	case "startWith", "start_with":
		return strings.HasPrefix(value, pattern)
	case "endWith", "end_with":
		return strings.HasSuffix(value, pattern)
	case "contains":
		return strings.Contains(value, pattern)
	case "containsAny", "contains_any":
		return strings.ContainsAny(value, pattern)
	case "equal":
		return value == pattern
	case "not_equal":
		return value != pattern
	}
	return false
}
