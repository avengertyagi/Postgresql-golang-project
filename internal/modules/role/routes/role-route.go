package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	rolecontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/role/controllers"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.RouterGroup, roleController *rolecontrollers.RoleController) {

	protected := r.Group("role")
	protected.GET("/", middlewares.HasPermission("role-list"), roleController.List)
	protected.POST("/create", middlewares.HasPermission("role-create"), roleController.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("role-edit"), roleController.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("role-update"), roleController.Update)
	protected.PATCH("/status/:id", middlewares.HasPermission("role-status"), roleController.UpdateStatus)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("role-delete"), roleController.Delete)
}
