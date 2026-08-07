package middlewares

import (
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	"github.com/gin-gonic/gin"
)

func HasPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions := c.MustGet("permissions").([]string)
		if !helpers.Contains(permissions, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.Forbidden})
			return
		}
		c.Next()
	}
}
