package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	admincontroller "github.com/akshit_tyagi/postgresql_project/internal/controllers/admin"
	rolecontroller "github.com/akshit_tyagi/postgresql_project/internal/controllers/role"
	tenantcontroller "github.com/akshit_tyagi/postgresql_project/internal/controllers/tenant"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.POST("/login", middlewares.ThrottleFailures(5, 2), admincontroller.Login)
	r.POST("/refresh", middlewares.ThrottleFailures(5, 2), admincontroller.RefreshToken)
	protected := r.Group("")
	protected.POST("/tenant", tenantcontroller.Create)
	protected.Use(middlewares.AuthMiddleware())
	protected.Use(middlewares.GuardMiddleware(constants.AdminGuard))
	{
		protected.POST("/logout", admincontroller.Logout)
		protected.GET("/profile", admincontroller.GetProfile)

		//Roles Routes
		protected.GET("/roles", rolecontroller.GetAll)
		protected.POST("/role/create", rolecontroller.Create)
		protected.GET("/role/edit/:id", rolecontroller.GetByID)
		protected.PUT("/role/update/:id", rolecontroller.Update)
		protected.DELETE("/role/destroy/:id", rolecontroller.Delete)
	}
}
