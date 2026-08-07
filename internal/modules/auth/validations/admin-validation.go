package validations

import (
	"fmt"
	"regexp"

	request "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/request/admin"
)

func AdminLoginValidation(request request.LoginRequest) error {
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
