package validations

import (
	"fmt"

	authmodel "github.com/akshit_tyagi/postgresql_project/internal/models"
)

func ValidateRole(req authmodel.RoleRequest) error {
	if req.Name == "" {
		return fmt.Errorf("Role name is required.")
	}
	if len(req.PermissionIDs) == 0 {
		return fmt.Errorf("Role permission is required.")
	}
	return nil
}
