package utils

// StringArraytoByteArr ...
func StringArraytoByteArr(dat []string) []byte {
	all := ""
	for _, v := range dat {
		all += v
	}

	return []byte(all)
}

// StringArrtoString ...
func StringArrtoString(dat []string) string { return string(StringArraytoByteArr(dat)) }
