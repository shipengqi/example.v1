package e

// Common errors
var (
	OK                = add(0, "OK")

	ErrInternalServer = add(-500, "Internal server error")
	ErrBadRequest     = add(-400, "Bad request")
	ErrUnauthorized   = add(-401, "Unauthorized")
	ErrForbidden      = add(-403, "Forbidden")
	ErrNothingFound   = add(-404, "Not found") // not found
	ErrTokenExpired   = add(-4011, "token is expired")
	ErrTokenMalformed = add(-4012, "token malformed")
	ErrNotValidYet    = add(-4013, "token not valid yet")
	ErrTokenInvalid   = add(-4014, "invalid token")
	ErrGenTokenFailed = add(-4015, "generate token failed")
)
