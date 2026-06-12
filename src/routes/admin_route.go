package route

import (
	"github.com/akshit_tyagi/postgresql_project/src/constants"
	admincontrollers "github.com/akshit_tyagi/postgresql_project/src/controllers/admin"
	middleware "github.com/akshit_tyagi/postgresql_project/src/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", middleware.ThrottleFailures(5, 2), admincontrollers.Login)
	r.POST("/refresh", admincontrollers.RefreshToken)
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.GuardMiddleware(constants.ADMIN_GUARD))
	{
		protected.POST("/logout", admincontrollers.Logout)
		protected.GET("/profile", admincontrollers.GetProfile)
		//protected.PUT("/profile", admincontrollers.UpdateProfile)
	}
}
