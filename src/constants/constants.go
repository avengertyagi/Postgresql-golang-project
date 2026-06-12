package constants

const (
	SUPER_ADMIN_ROLE uint8 = 1
	STAFF_ROLE       uint8 = 2
	USER_ROLE        uint8 = 3
)

const (
	ADMIN_GUARD = "admin"
	USER_GUARD  = "user"
)

const (
	USER_NOT_FOUND          = "User not found"
	INVALID_CREDENTIALS     = "These credentials do not match our records."
	INACTIVE_ACCOUNT        = "Your account is inactive. Please contact support."
	AUTHORIZATION_HEADER    = "Authorization header is required"
	BAD_AUTH_FORMAT         = "Authorization format must be: Bearer <token>"
	UNAUTHENTICATED         = "Unauthenticated"
	FORBIDDEN               = "You do not have permission to access this resource"
	ACCESS_DENIED           = "Access denied for this guard"
	LOGIN_SUCCESS           = "Login successful"
	LOGOUT_SUCCESS          = "Logout successful"
	REFRESH_SUCCESS         = "Token refreshed successfully"
	SIGN_UP_SUCCESS         = "User registered successfully"
	SESSION_ALREADY_REVOKED = "Your session has already been logged out. Please login again."
	SESSION_EXPIRED         = "Your session has expired. Please login again."
	SESSION_NOT_FOUND       = "Session not found. Please login again."
	PROFILE_FETCH_SUCCESS   = "Profile fetched successfully"
	PROFILE_UPDATE_SUCCESS  = "Profile updated successfully"
	NOT_FOUND               = "Record not found."

	ROLE_ALREADY_EXISTS     = "Role already exists"
	ROLE_CREATED_SUCCESS    = "Role created successfully"
	ROLE_UPDATED_SUCCESS    = "Role updated successfully"
	ROLE_DELETED_SUCCESS    = "Role deleted successfully"
	ROLE_NOT_FOUND          = "Role not found"
)
