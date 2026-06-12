package validation

import (
	"fmt"
	"regexp"

	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
)

func UserSignUpValidation(request authmodel.UserSignUpRequest) error {
	if request.Name == "" {
		return fmt.Errorf("Name is required.")
	}
	if request.Name != "" && len(request.Name) < 3 {
		return fmt.Errorf("Name must be at least 3 characters long.")
	}
	if request.Name != "" && !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(request.Name) {
		return fmt.Errorf("Name can only contain letters and spaces.")
	}
	if request.Name != "" && regexp.MustCompile(`\s{2,}`).MatchString(request.Name) {
		return fmt.Errorf("Name cannot contain multiple consecutive spaces.")
	}
	if request.Name != "" && (request.Name[0] == ' ' || request.Name[len(request.Name)-1] == ' ') {
		return fmt.Errorf("Name cannot start or end with a space.")
	}
	if request.Name != "" && regexp.MustCompile(`\d`).MatchString(request.Name) {
		return fmt.Errorf("Name cannot contain numbers.")
	}
	if request.Name != "" && regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(request.Name) {
		return fmt.Errorf("Name cannot contain special characters.")
	}
	if request.Name != "" && len(request.Name) > 50 {
		return fmt.Errorf("Name cannot be longer than 50 characters.")
	}
	if request.Email == "" {
		return fmt.Errorf("Email is required.")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(request.Email) {
		return fmt.Errorf("Please enter a valid email address")
	}
	if request.Password == "" {
		return fmt.Errorf("Password is required.")
	}
	if len(request.Password) < 8 {
		return fmt.Errorf("Password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(request.Password) {
		return fmt.Errorf("Password must contain at least one uppercase letter.")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(request.Password) {
		return fmt.Errorf("Password must contain at least one lowercase letter.")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(request.Password) {
		return fmt.Errorf("Password must contain at least one digit.")
	}
	return nil
}
