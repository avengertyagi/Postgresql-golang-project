package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	rolecontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/role/controllers"
	rolerepositories "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
	roleservices "github.com/akshit_tyagi/postgresql_project/internal/modules/role/services"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.RouterGroup) {

	roleRepo := rolerepositories.NewRoleRepository(config.DB)
	roleService := roleservices.NewRoleService(roleRepo)
	roleController := rolecontrollers.NewRoleController(roleService)

	protected := r.Group("role")
	protected.GET("/", middlewares.HasPermission("role-list"), roleController.List)
	protected.POST("/create", middlewares.HasPermission("role-create"), roleController.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("role-edit"), roleController.GetByID)
	protected.PATCH("/update/:id", middlewares.HasPermission("role-update"), roleController.Update)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("role-delete"), roleController.Delete)
}
