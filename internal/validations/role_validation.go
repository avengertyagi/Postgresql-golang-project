package validations

import (
	"fmt"

	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/models/role"
)

func ValidateRole(req rolemodel.RoleRequest) error {
	if req.Name == "" {
		return fmt.Errorf("Role name is required.")
	}
	if len(req.PermissionIDs) == 0 {
		return fmt.Errorf("Role permission is required.")
	}
	return nil
}
