package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	rolecontroller "github.com/akshit_tyagi/postgresql_project/internal/modules/role/controllers"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.RouterGroup) {
	{
		protected := r.Group("role")
		protected.GET("/", middlewares.HasPermission("role-list"), rolecontroller.GetAll)
		protected.POST("/create", middlewares.HasPermission("role-create"), rolecontroller.Create)
		protected.GET("/edit/:id", middlewares.HasPermission("role-edit"), rolecontroller.GetByID)
		protected.PUT("/update/:id", middlewares.HasPermission("role-update"), rolecontroller.Update)
		protected.DELETE("/destroy/:id", middlewares.HasPermission("role-delete"), rolecontroller.Delete)
	}
}
