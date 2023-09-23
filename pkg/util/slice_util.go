package util

func IsStrInSlice(str string, strSlice []string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}
	return false
}
