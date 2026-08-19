package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	staffcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/controllers"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(r *gin.RouterGroup, staffController *staffcontrollers.StaffController) {

	protected := r.Group("staff")
	protected.GET("/", middlewares.HasPermission("staff-list"), staffController.List)
	protected.POST("/create", middlewares.HasPermission("staff-create"), staffController.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("staff-edit"), staffController.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("staff-update"), staffController.Update)
	protected.PATCH("/status/:id", middlewares.HasPermission("staff-status"), staffController.UpdateStatus)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("staff-delete"), staffController.Delete)
}
