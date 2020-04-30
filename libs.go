package terror

// StringSliceReverse return reverse ordered string slice
func StringSliceReverse(arr []string) []string {
	l := len(arr)
	strs := make([]string, l)

	for i, str := range arr {
		strs[l-i-1] = str
	}

	return strs
}
