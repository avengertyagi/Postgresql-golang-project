package middlewares

import (
	"net/http"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	"github.com/gin-gonic/gin"
)

func HasPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// role, exists := c.Get("role")
		// if exists {
		// 	roleStr, ok := role.(string)
		// 	if ok && roleStr == fmt.Sprintf("%d", constants.SuperAdminRole) {
		// 		c.Next()
		// 		return
		// 	}
		// }

		permissionsVal, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.Forbidden.Error()})
			return
		}
		permissions, ok := permissionsVal.([]string)
		if !ok || !helpers.Contains(permissions, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.Forbidden.Error()})
			return
		}
		c.Next()
	}
}
