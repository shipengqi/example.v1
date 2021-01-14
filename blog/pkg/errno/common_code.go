package errno

// Common errors
var (
	OK = add(0, "OK")

	ErrInternalServer = add(-500, "Internal server error")
	ErrBadRequest     = add(-400, "Bad request")
	ErrUnauthorized   = add(-401, "Unauthorized")
	ErrForbidden      = add(-403, "Forbidden")
	ErrNothingFound   = add(-404, "Not found") // not found
)
