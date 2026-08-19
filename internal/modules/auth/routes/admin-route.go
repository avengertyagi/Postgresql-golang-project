package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	authcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/controllers/admin"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup, authController *authcontrollers.AdminController) {
	protected := r.Group("admin")
	protected.POST("/login", middlewares.ThrottleFailures(5, 2), authController.Login)
	protected.POST("/refresh", middlewares.ThrottleFailures(5, 2), authController.RefreshToken)
	protected.Use(middlewares.AuthMiddleware())
	protected.Use(middlewares.GuardMiddleware(constants.AdminGuard))
	protected.POST("/logout", authController.Logout)
	protected.GET("/profile", authController.GetProfile)
	protected.PUT("/profile", authController.UpdateProfile)
}
