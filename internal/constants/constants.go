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

const (
	UserNotFound          = "User not found"
	AuthorizationHeader   = "Authorization header is required"
	BadAuthFormat         = "Authorization format must be: Bearer <token>"
	Unauthenticated       = "Unauthenticated"
	Forbidden             = "You do not have permission to access this resource"
	AccessDenied          = "Access denied for this guard"
	LoginSuccess          = "Login successful"
	LogoutSuccess         = "Logout successful"
	RefreshSuccess        = "Token refreshed successfully"
	SignUpSuccess         = "User registered successfully"
	SessionAlreadyRevoked = "Your session has already been logged out. Please login again."
	SessionExpired        = "Your session has expired. Please login again."
	SessionNotFound       = "Session not found. Please login again."
	ProfileFetchSuccess   = "Profile fetched successfully"
	ProfileUpdateSuccess  = "Profile updated successfully"
	NotFound              = "Record not found."
	SomethingWentWrong    = "Something went wrong. Please try again later."

	RoleFetchedSuccess   = "Role fetched successfully"
	RoleRetrievedSuccess = "Role retrieved successfully"
	RoleCreatedSuccess   = "Role created successfully"
	RoleUpdatedSuccess   = "Role updated successfully"
	RoleDeletedSuccess   = "Role deleted successfully"
)

var RoleAlreadyExists = errors.New("Role already exists")
var InvalidCredentials = errors.New("These credentials do not match our records.")
var InactiveAccount = errors.New("Your account is inactive. Please contact support.")
var RoleNotFound = errors.New("Role not found")
