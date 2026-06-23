package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	"github.com/gin-gonic/gin"
)

func requestID(c *gin.Context) string {
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("auth: missing authorization header",
				"request_id", requestID(c),
				"path", c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.AuthorizationHeader})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.BadAuthFormat})
			return
		}
		claims, err := helpers.ParseAccessToken(parts[1])
		if err != nil {
			slog.Warn("auth: invalid access token",
				"request_id", requestID(c),
				"error", err.Error(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.SessionNotFound})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("guard", claims.Guard)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated})
			return
		}
		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.Unauthenticated})
			return
		}
		if _, ok := roleSet[roleStr]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.Forbidden})
			return
		}
		c.Next()
	}
}

func GuardMiddleware(guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		guard, exists := c.Get("guard")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.AccessDenied})
			return
		}
		guardStr, ok := guard.(string)
		if !ok || guardStr != guardName {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.AccessDenied})
			return
		}
		c.Next()
	}
}
