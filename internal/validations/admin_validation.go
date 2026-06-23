package validations

import (
	"fmt"
	"regexp"

	usermodel "github.com/akshit_tyagi/postgresql_project/internal/models/user"
)

func AdminLoginValidation(request usermodel.AdminLoginRequest) error {
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
