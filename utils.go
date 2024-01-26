package tabby

func AddPrefix(strs []string, prefix string) []string {
	var result []string
	for _, str := range strs {
		result = append(result, prefix+str)
	}
	return result
}

func MapKeys[keyT string, valueT string](m map[keyT]valueT) []keyT {
	var keys []keyT
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
