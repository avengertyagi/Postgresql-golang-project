package constants

import "errors"

var (
	ThemeFetchedSuccess       = errors.New("Theme fetched successfully.")
	ThemeRetrievedSuccess     = errors.New("Theme retrieved successfully.")
	ThemeCreatedSuccess       = errors.New("Theme created successfully.")
	ThemeUpdatedSuccess       = errors.New("Theme updated successfully.")
	ThemeStatusUpdatedSuccess = errors.New("Theme status updated successfully.")
	ThemeDeletedSuccess       = errors.New("Theme deleted successfully.")

	ThemeAlreadyExists = errors.New("Theme already exists.")
	InvalidCredentials = errors.New("These credentials do not match our records.")
	InactiveAccount    = errors.New("Your account is inactive. Please contact support.")
	ThemeNotFound      = errors.New("Theme not found.")
	ThemeIdNotFound    = errors.New("The selected theme id does not exist.")
)
