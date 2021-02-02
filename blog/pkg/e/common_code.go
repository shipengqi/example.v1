package e

// Common errors
var (
	OK                = add(0, "OK")

	ErrInternalServer = add(-500, "Internal server error")
	ErrBadRequest     = add(-400, "Bad request")
	ErrUnauthorized   = add(-401, "Unauthorized")
	ErrForbidden      = add(-403, "Forbidden")
	ErrNothingFound   = add(-404, "Not found") // not found
	ErrTokenExpired   = add(-4011, "Token is expired")
	ErrTokenMalformed = add(-4012, "Token malformed")
	ErrNotValidYet    = add(-4013, "Token not valid yet")
	ErrTokenInvalid   = add(-4014, "Invalid token")
	ErrGenTokenFailed = add(-4015, "Generate token failed")
	ErrClaimsType     = add(-4016, "Claims type error")
	ErrUserLocked     = add(-5001, "User is locked")
	ErrUserDeleted    = add(-5002, "User is deleted")
	ErrPassWrong      = add(-5003, "Password is wrong")
)
