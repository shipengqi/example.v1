package errno

// Common errors
var (
	OK                  = add(0, "OK")

	InternalServerError = add(-500, "Internal server error")
	NothingFound        = add(-404, "not found") // not found
)
