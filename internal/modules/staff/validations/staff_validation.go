package validations

import (
	"fmt"
	"regexp"
	"strconv"

	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
)

func Validate(req dto.StaffRequest) error {
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
	if req.Name != "" && (req.Name[0] == ' ' || req.Name[len(req.Name)-1] == ' ') {
		return fmt.Errorf("Name cannot start or end with a space.")
	}
	if req.Name != "" && regexp.MustCompile(`\d`).MatchString(req.Name) {
		return fmt.Errorf("Name cannot contain numbers.")
	}
	if req.Name != "" && regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(req.Name) {
		return fmt.Errorf("Name cannot contain special characters.")
	}
	if req.Name != "" && len(req.Name) > 50 {
		return fmt.Errorf("Name cannot be longer than 50 characters.")
	}
	if req.Email == "" {
		return fmt.Errorf("Email is required.")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(req.Email) {
		return fmt.Errorf("Please enter a valid email address")
	}
	if req.Password == "" {
		return fmt.Errorf("Password is required.")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("Password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.Password) {
		return fmt.Errorf("Password must contain at least one uppercase letter.")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(req.Password) {
		return fmt.Errorf("Password must contain at least one lowercase letter.")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(req.Password) {
		return fmt.Errorf("Password must contain at least one digit.")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(req.Password) {
		return fmt.Errorf("Password must contain special characters.")
	}
	if req.RoleID == 0 {
		return fmt.Errorf("Role ID is required.")
	}
	if req.RoleID < 0 {
		return fmt.Errorf("Role ID must be a positive integer.")
	}
	if req.Mobile == 0 {
		return fmt.Errorf("Mobile number is required.")
	}
	if len(strconv.Itoa(req.Mobile)) < 10 {
		return fmt.Errorf("Mobile number must be at least 10 digits long.")
	}
	if len(strconv.Itoa(req.Mobile)) > 15 {
		return fmt.Errorf("Mobile number cannot be longer than 15 digits.")
	}
	return nil
}
