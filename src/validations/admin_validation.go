package validation

import (
	"fmt"
	"regexp"

	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
)

func AdminLoginValidation(request authmodel.AdminLoginRequest) error {
	if request.Email == "" {
		return fmt.Errorf("Email is required.")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(request.Email) {
		return fmt.Errorf("Please enter a valid email address")
	}
	if request.Password == "" {
		return fmt.Errorf("Password is required.")
	}
	return nil
}
