package tabby

func AddPrefix(strs []string, prefix string) []string {
	var result []string
	for _, str := range strs {
		result = append(result, prefix+str)
	}
	return result
}
