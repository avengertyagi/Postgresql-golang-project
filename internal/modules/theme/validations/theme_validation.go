package validations

import (
	"fmt"
	"regexp"

	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
)

func Validate(req dto.ThemeRequest) error {
	if req.Name == "" {
		return fmt.Errorf("Name is required.")
	}
	if req.Name != "" && len(req.Name) < 3 {
		return fmt.Errorf("Name must be at least 3 characters long.")
	}
	if req.Name != "" && len(req.Name) > 50 {
		return fmt.Errorf("Name cannot exceed 50 characters.")
	}
	if req.Name != "" && !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(req.Name) {
		return fmt.Errorf("Name can only contain letters and spaces.")
	}
	if req.Name != "" && regexp.MustCompile(`\s{2,}`).MatchString(req.Name) {
		return fmt.Errorf("Name cannot contain multiple consecutive spaces.")
	}
	if req.Description == "" {
		return fmt.Errorf("Description is required.")
	}
	if req.Description != "" && len(req.Description) < 3 {
		return fmt.Errorf("Description must be at least 3 characters long.")
	}
	if req.Description != "" && len(req.Description) > 200 {
		return fmt.Errorf("Description cannot exceed 200 characters.")
	}
	if req.Image != "" && regexp.MustCompile(`.*\.(jpg|jpeg|png|gif)$`).MatchString(req.Image) == false {
		return fmt.Errorf("Image must be a valid image URL (jpg, jpeg, png, gif).")
	}
	return nil
}
