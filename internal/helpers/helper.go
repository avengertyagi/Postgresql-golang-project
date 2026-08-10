package helpers

import (
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
)

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

func ParseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
