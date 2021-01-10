package errno

var (
	// Common errors
	OK                  = New(0, "OK")
	InternalServerError = New(10001, "Internal server error")
	ErrBind             = New(10002, "Error occurred while binding the request body to the struct.")

	ErrValidation = New(20001, "Validation failed.")
	ErrDatabase   = New(20002, "Database error.")
	ErrToken      = New(20003, "Error occurred while signing the JSON web token.")

	// user errors
	ErrEncrypt           = New(20101, "Error occurred while encrypting the user password.")
	ErrUserNotFound      = New(20102, "The user was not found.")
	ErrTokenInvalid      = New(20103, "The token was invalid.")
	ErrPasswordIncorrect = New(20104, "The password was incorrect.")
)
