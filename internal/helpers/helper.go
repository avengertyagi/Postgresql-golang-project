package helpers

import "unicode"

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Ucfirst(s string) string {
	if s == "" {
		return ""
	}
	rune := []rune(s)
	rune[0] = unicode.ToUpper(rune[0])
	return string(rune)
}
