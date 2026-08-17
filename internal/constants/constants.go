package constants

import "errors"

const (
	SuperAdminRole uint8 = 0
	StaffRole      uint8 = 1
	UserRole       uint8 = 2
)

const (
	AdminGuard = "admin"
	UserGuard  = "user"
)

var (
	UserNotFound          = errors.New("User not found")
	AuthorizationHeader   = errors.New("Authorization header is required")
	BadAuthFormat         = errors.New("Authorization format must be: Bearer <token>")
	Unauthenticated       = errors.New("Unauthenticated")
	Forbidden             = errors.New("You do not have permission to access this resource.")
	AccessDenied          = errors.New("Access denied for this guard")
	LoginSuccess          = errors.New("Login successfully")
	LogoutSuccess         = errors.New("Logout successfully")
	RefreshSuccess        = errors.New("Token refreshed successfully")
	SignUpSuccess         = errors.New("User registered successfully")
	SessionAlreadyRevoked = errors.New("Your session has already been logged out. Please login again.")
	SessionExpired        = errors.New("Your session has expired. Please login again.")
	SessionNotFound       = errors.New("Session not found. Please login again.")
	ProfileFetchSuccess   = errors.New("Profile fetched successfully")
	ProfileUpdateSuccess  = errors.New("Profile updated successfully")
	NotFound              = errors.New("Record not found.")
	SomethingWentWrong    = errors.New("Something went wrong. Please try again later.")

	RoleFetchedSuccess   = errors.New("Role fetched successfully")
	RoleRetrievedSuccess = errors.New("Role retrieved successfully")
	RoleCreatedSuccess   = errors.New("Role created successfully")
	RoleUpdatedSuccess   = errors.New("Role updated successfully")
	RoleDeletedSuccess   = errors.New("Role deleted successfully")

	RoleAlreadyExists  = errors.New("Role already exists")
	InvalidCredentials = errors.New("These credentials do not match our records.")
	InactiveAccount    = errors.New("Your account is inactive. Please contact support.")
	RoleNotFound       = errors.New("Role not found")

	StaffFetchedSuccess   = errors.New("Staff fetched successfully")
	StaffRetrievedSuccess = errors.New("Staff retrieved successfully")
	StaffCreatedSuccess   = errors.New("Staff created successfully")
	StaffUpdatedSuccess   = errors.New("Staff updated successfully")
	StaffDeletedSuccess   = errors.New("Staff deleted successfully")

	StaffNotFound        = errors.New("Staff not found")
	StaffAlreadyExists   = errors.New("Staff already exists")
	StaffInactiveAccount = errors.New("Staff is inactive. Please contact support.")
)
