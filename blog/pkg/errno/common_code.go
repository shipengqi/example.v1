package errno

// Common errors
var (
	OK                = add(0, "OK")

	ErrInternalServer = add(-500, "Internal server error")
	ErrBadRequest     = add(-400, "bad request")
	ErrNothingFound   = add(-404, "not found") // not found
)
