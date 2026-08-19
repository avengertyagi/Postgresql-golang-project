package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/bootstrap"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	authroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/routes"
	roleroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/role/routes"
	staffroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/routes"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, c *bootstrap.Container) {
	adminProtected := r.Group("admin")
	adminProtected.Use(middlewares.AuthMiddleware())
	adminProtected.Use(middlewares.GuardMiddleware(constants.AdminGuard))

	authroutes.AdminRoutes(r, c.AdminController)
	roleroutes.RoleRoutes(adminProtected, c.RoleController)
	staffroutes.StaffRoutes(adminProtected, c.StaffController)
}
