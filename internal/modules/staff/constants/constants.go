package constants

import "errors"

var (
	StaffFetchedSuccess   = errors.New("Staff fetched successfully.")
	StaffRetrievedSuccess = errors.New("Staff retrieved successfully.")
	StaffCreatedSuccess   = errors.New("Staff created successfully.")
	StaffUpdatedSuccess   = errors.New("Staff updated successfully.")
	StaffDeletedSuccess   = errors.New("Staff deleted successfully.")

	StaffNotFound        = errors.New("Staff not found.")
	StaffAlreadyExists   = errors.New("Staff already exists.")
	StaffInactiveAccount = errors.New("Staff is inactive. Please contact support.")
	RoleIdNotFound       = errors.New("The selected role id does not exist.")
)
