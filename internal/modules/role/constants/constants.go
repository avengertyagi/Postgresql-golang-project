package constants

import "errors"

var (
	RoleFetchedSuccess   = errors.New("Role fetched successfully.")
	RoleRetrievedSuccess = errors.New("Role retrieved successfully.")
	RoleCreatedSuccess   = errors.New("Role created successfully.")
	RoleUpdatedSuccess   = errors.New("Role updated successfully.")
	RoleDeletedSuccess   = errors.New("Role deleted successfully.")

	RoleAlreadyExists  = errors.New("Role already exists.")
	InvalidCredentials = errors.New("These credentials do not match our records.")
	InactiveAccount    = errors.New("Your account is inactive. Please contact support.")
	RoleNotFound       = errors.New("Role not found.")
	RoleIdNotFound     = errors.New("The selected role id does not exist.")
)
