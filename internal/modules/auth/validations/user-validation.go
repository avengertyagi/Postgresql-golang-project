package validations

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
)

func UserRegisterValidation(req dto.UserRegisterRequest) error {
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
	if req.Mobile == 0 {
		return fmt.Errorf("Mobile number is required.")
	}
	if len(strconv.Itoa(req.Mobile)) < 10 {
		return fmt.Errorf("Mobile number must be at least 10 digits long.")
	}
	if len(strconv.Itoa(req.Mobile)) > 15 {
		return fmt.Errorf("Mobile number cannot be longer than 15 digits.")
	}
	if req.CountryCode == "" {
		return fmt.Errorf("Country code is required.")
	}
	if req.Dob == "" {
		return fmt.Errorf("Date of birth is required.")
	}
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return fmt.Errorf("Date of birth must be in the format YYYY-MM-DD.")
	}
	if dob.Year() < 1900 {
		return fmt.Errorf("Year in date of birth must be 1900 or later.")
	}
	if dob.After(time.Now()) {
		return fmt.Errorf("Date of birth cannot be a future date.")
	}
	return nil
}
