package util

// StrSliceContains returns true if slice of strings contains a given string
func StrSliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
