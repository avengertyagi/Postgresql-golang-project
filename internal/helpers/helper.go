package helpers

import (
	"strconv"

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

func ParseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func StringOrNA(s *string) string {
	if s == nil || *s == "" {
		return "NA"
	}
	return *s
}
