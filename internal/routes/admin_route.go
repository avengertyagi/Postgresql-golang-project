package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/controllers/admin"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", middlewares.ThrottleFailures(5, 2), admin.Login)
	r.POST("/refresh", middlewares.ThrottleFailures(5, 2), admin.RefreshToken)
	protected := r.Group("")
	protected.Use(middlewares.AuthMiddleware())
	protected.Use(middlewares.GuardMiddleware(constants.AdminGuard))
	{
		protected.POST("/logout", admin.Logout)
		protected.GET("/profile", admin.GetProfile)

		protected.POST("/role/store", admin.CreateRole)
	}
}
