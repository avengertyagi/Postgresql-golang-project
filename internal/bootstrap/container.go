package bootstrap

import (
	"gorm.io/gorm"

	authcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/controllers/admin"
	authrepositories "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	authservices "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/services"

	rolecontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/role/controllers"
	rolerepositories "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
	roleservices "github.com/akshit_tyagi/postgresql_project/internal/modules/role/services"

	staffcontrollers "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/controllers"
	staffrepositories "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/repositories"
	staffservices "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/services"
)

type Container struct {
	AdminController *authcontrollers.AdminController
	RoleController  *rolecontrollers.RoleController
	StaffController *staffcontrollers.StaffController
}

func NewContainer(db *gorm.DB) *Container {
	// auth module
	authRepo := authrepositories.NewAuthRepository(db)
	authService := authservices.NewAuthService(authRepo)
	adminController := authcontrollers.NewAuthController(authService)

	// role module
	roleRepo := rolerepositories.NewRoleRepository(db)
	roleService := roleservices.NewRoleService(roleRepo)
	roleController := rolecontrollers.NewRoleController(roleService)

	// staff module (depends on both staff repo and role repo)
	staffRepo := staffrepositories.NewStaffRepository(db)
	staffService := staffservices.NewStaffService(staffRepo, roleRepo)
	staffController := staffcontrollers.NewStaffController(staffService)

	return &Container{
		AdminController: adminController,
		RoleController:  roleController,
		StaffController: staffController,
	}
}
