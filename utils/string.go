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
func StringArrtoString(dat []string) string {
	all := ""
	for _, v := range dat {
		all += v
	}
	return all
}
