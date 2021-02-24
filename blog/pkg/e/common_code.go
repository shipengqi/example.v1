package e

// Common errors
var (
	OK                = add(0, "OK")

	ErrInternalServer = add(-500, "Internal server error")
	ErrBadRequest     = add(-400, "Bad request")
	ErrUnauthorized   = add(-401, "Unauthorized")
	ErrForbidden      = add(-403, "Forbidden")
	ErrNothingFound   = add(-404, "Not found") // not found
	ErrMultiFormErr   = add(-4001, "Multipart Form error")
	ErrCheckImage     = add(-4002, "Image check error")
	ErrTokenExpired   = add(-4011, "Token is expired")
	ErrTokenMalformed = add(-4012, "Token malformed")
	ErrNotValidYet    = add(-4013, "Token not valid yet")
	ErrTokenInvalid   = add(-4014, "Invalid token")
	ErrUserNotFound   = add(-4041, "User not found")
	ErrUserLocked     = add(-5001, "User is locked")
	ErrUserDeleted    = add(-5002, "User is deleted")
	ErrPassWrong      = add(-5003, "Password is wrong")
	ErrUploadImage    = add(-5004, "Upload image error")
)
