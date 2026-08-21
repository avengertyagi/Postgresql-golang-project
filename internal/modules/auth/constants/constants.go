package constants

import "errors"

var (
	LoginSuccess          = errors.New("Login successfully.")
	LogoutSuccess         = errors.New("Logout successfully.")
	RefreshSuccess        = errors.New("Token refreshed successfully.")
	SignUpSuccess         = errors.New("User registered successfully.")
	SessionAlreadyRevoked = errors.New("Your session has already been logged out. Please login again.")
	SessionExpired        = errors.New("Your session has expired. Please login again.")
	SessionNotFound       = errors.New("Session not found. Please login again.")
	ProfileFetchSuccess   = errors.New("Profile fetched successfully.")
	ProfileUpdateSuccess  = errors.New("Profile updated successfully.")
	InvalidCredentials    = errors.New("These credentials do not match our records.")
	InactiveAccount       = errors.New("Your account is inactive. Please contact support.")
	EmailAlreadyExists    = errors.New("Email already exists. Please use a different email.")
	MobileAlreadyExists   = errors.New("Mobile number already exists. Please use a different mobile number.")
)
