package validations

import (
	"fmt"
	"regexp"

	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func ValidateRole(req rolemodel.RoleRequest) error {
	if req.Name == "" {
		return fmt.Errorf("Name is required.")
	}
	if req.Name != "" && len(req.Name) < 3 {
		return fmt.Errorf("Name must be at least 3 characters long.")
	}
	if req.Name != "" && !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(req.Name) {
		return fmt.Errorf("Name can only contain letters and spaces.")
	}
	if req.Name != "" && regexp.MustCompile(`\s{2,}`).MatchString(req.Name) {
		return fmt.Errorf("Name cannot contain multiple consecutive spaces.")
	}
	if len(req.PermissionIDs) == 0 {
		return fmt.Errorf("Role permission is required.")
	}
	return nil
}
