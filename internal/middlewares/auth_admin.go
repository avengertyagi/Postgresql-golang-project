package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	authrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(repo authrepo.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("auth: missing authorization header",
				"request_id", requestid.Get(c),
				"path", c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.AuthorizationHeader.Error()})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.BadAuthFormat.Error()})
			return
		}

		claims, err := helpers.ParseAccessToken(parts[1])
		if err != nil {
			slog.Warn("auth: invalid access token",
				"request_id", requestid.Get(c),
				"error", err.Error(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "statusCode": http.StatusUnauthorized, "message": constants.SessionNotFound.Error()})
			return
		}
		user, err := repo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			slog.Warn("auth: user not found",
				"request_id", requestid.Get(c),
				"user_id", claims.UserID,
			)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "statusCode": http.StatusNotFound, "message": constants.UserNotFound.Error()})
			return
		}
		if !user.Status {
			slog.Warn("auth: inactive user",
				"request_id", requestid.Get(c),
				"user_id", claims.UserID,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.InactiveAccount.Error()})
			return
		}
		if user.UserType != constants.SuperAdminRole && user.UserType != constants.StaffRole {
			slog.Warn("auth: unauthorized access",
				"request_id", requestid.Get(c),
				"user_id", claims.UserID,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "statusCode": http.StatusForbidden, "message": constants.Unauthorized.Error()})
			return
		}

		slog.Debug("auth: token verified", "request_id", requestid.Get(c), "user_id", claims.UserID, "guard", claims.Guard)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("guard", claims.Guard)
		c.Set("permissions", claims.Permissions)
		c.Next()
	}
}
