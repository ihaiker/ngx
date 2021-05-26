package stdlib

func init() {
	_methods["and"] = and
	_methods["or"] = or
	_methods["if"] = ifelse
	_methods["not"] = not
}

func and(left, right bool) bool {
	return left && right
}

func or(left, right bool) bool {
	return left || right
}

func ifelse(check bool, left, right interface{}) interface{} {
	if check {
		return left
	} else {
		return right
	}
}

func not(check bool) bool {
	return !check
}
