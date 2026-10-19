package routes

import (
	middlewares "github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	themecontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/controllers"
	"github.com/gin-gonic/gin"
)

func ThemeRoutes(r *gin.RouterGroup) {
	protected := r.Group("theme")
	protected.GET("/", middlewares.HasPermission("theme-list"), themecontrollers.List)
	protected.POST("/create", middlewares.HasPermission("theme-create"), themecontrollers.Create)
	protected.GET("/edit/:id", middlewares.HasPermission("theme-edit"), themecontrollers.GetByID)
	protected.PUT("/update/:id", middlewares.HasPermission("theme-update"), themecontrollers.Update)
	protected.PATCH("/status/:id", middlewares.HasPermission("theme-status"), themecontrollers.UpdateStatus)
	protected.DELETE("/destroy/:id", middlewares.HasPermission("theme-delete"), themecontrollers.Delete)
}
