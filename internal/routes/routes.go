package routes

import (
	"github.com/akshit_tyagi/postgresql_project/internal/middlewares"
	authroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/routes"
	roleroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/role/routes"
	socketroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/socket/routes"
	staffroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/routes"
	themeroutes "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/routes"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	authroutes.AdminRoutes(r)

	adminProtected := r.Group("admin")
	adminProtected.Use(middlewares.AuthMiddleware())

	roleroutes.RoleRoutes(adminProtected)
	staffroutes.StaffRoutes(adminProtected)
	themeroutes.ThemeRoutes(adminProtected)

	userProtected := r.Group("user")
	authroutes.UserRoutes(userProtected)
	socketroutes.SocketRoutes(userProtected)
}
