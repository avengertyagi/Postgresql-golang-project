package validation

import (
	"fmt"

	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
)

func ValidateRole(req authmodel.RoleRequest) error {
	if req.Name == "" {
		return fmt.Errorf("Role name is required.")
	}
	if len(req.Permission) == 0 {
		return fmt.Errorf("Role permission is required.")
	}
	return nil
}
