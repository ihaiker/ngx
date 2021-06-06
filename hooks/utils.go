package hooks

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
