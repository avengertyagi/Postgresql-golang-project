package migrations

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	user "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	permission "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	personalaccesstoken "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	role "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func Migrate() error {
	return config.Migrate(
		&permission.Permission{},
		&role.Role{},
		&user.User{},
		&user.TempUser{},
		&user.OTP{},
		&personalaccesstoken.PersonalAccessToken{},
	)
}
