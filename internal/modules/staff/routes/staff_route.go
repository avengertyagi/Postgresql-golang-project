package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	staffcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/controllers"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(r *gin.RouterGroup) {
	protected := r.Group("staff")
	protected.GET("/", middlewares.HasPermission("staff-list"), staffcontrollers.List)
	protected.POST("/create", middlewares.HasPermission("staff-create"), staffcontrollers.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("staff-edit"), staffcontrollers.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("staff-update"), staffcontrollers.Update)
	protected.PATCH("/status/:id", middlewares.HasPermission("staff-status"), staffcontrollers.UpdateStatus)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("staff-delete"), staffcontrollers.Delete)
}
