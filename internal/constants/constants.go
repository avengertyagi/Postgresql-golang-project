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
	UserNotFound        = errors.New("User not found.")
	AuthorizationHeader = errors.New("Authorization header is required.")
	BadAuthFormat       = errors.New("Authorization format must be: Bearer <token>.")
	Unauthenticated     = errors.New("Unauthenticated.")
	Forbidden           = errors.New("You do not have permission to access this resource.")
	AccessDenied        = errors.New("Access denied for this guard.")
	SessionNotFound     = errors.New("Session not found. Please login again.")

	NotFound           = errors.New("Record not found.")
	SomethingWentWrong = errors.New("Something went wrong. Please try again later.")
	InvalidRequestBody = errors.New("Please provide valid JSON data.")
)
