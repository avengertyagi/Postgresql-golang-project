package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	admincontroller "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/controllers/admin"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	protected := r.Group("admin")
	protected.POST("/login", middlewares.ThrottleFailures(5, 2), admincontroller.Login)
	protected.POST("/refresh", middlewares.ThrottleFailures(5, 2), admincontroller.RefreshToken)
	protected.Use(middlewares.AuthMiddleware())
	protected.Use(middlewares.GuardMiddleware(constants.AdminGuard))
	protected.POST("/logout", admincontroller.Logout)
	protected.GET("/profile", admincontroller.GetProfile)
}
