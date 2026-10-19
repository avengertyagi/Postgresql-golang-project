package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	rolecontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/role/controllers"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.RouterGroup) {
	protected := r.Group("role")
	protected.GET("/", middlewares.HasPermission("role-list"), rolecontrollers.List)
	protected.POST("/create", middlewares.HasPermission("role-create"), rolecontrollers.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("role-edit"), rolecontrollers.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("role-update"), rolecontrollers.Update)
	protected.PATCH("/status/:id", middlewares.HasPermission("role-status"), rolecontrollers.UpdateStatus)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("role-delete"), rolecontrollers.Delete)
}
