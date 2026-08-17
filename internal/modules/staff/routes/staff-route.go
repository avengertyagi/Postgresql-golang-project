package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	staffcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/controllers"
	staffrepositories "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/repositories"
	staffservices "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/services"

	"github.com/gin-gonic/gin"
)

func StaffRoutes(r *gin.RouterGroup) {

	staffRepo := staffrepositories.NewStaffRepository(config.DB)
	staffService := staffservices.NewStaffService(staffRepo)
	staffController := staffcontrollers.NewStaffController(staffService)

	protected := r.Group("staff")
	protected.GET("/", middlewares.HasPermission("staff-list"), staffController.List)
	protected.POST("/create", middlewares.HasPermission("staff-create"), staffController.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("staff-edit"), staffController.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("staff-update"), staffController.Update)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("staff-delete"), staffController.Delete)
}
