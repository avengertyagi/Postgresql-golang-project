package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	authcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/controllers/admin"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	protected := r.Group("admin")
	protected.POST("/login", middlewares.ThrottleFailures(5, 2), authcontrollers.Login)
	protected.POST("/refresh", middlewares.ThrottleFailures(5, 2), authcontrollers.RefreshToken)
	protected.Use(middlewares.AuthMiddleware())
	protected.POST("/logout", authcontrollers.Logout)
	protected.GET("/profile", authcontrollers.GetProfile)
	protected.PUT("/profile", authcontrollers.UpdateProfile)
}
