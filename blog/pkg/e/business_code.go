package e

// Business errors
var (
	ErrExistTag    = New(10001, "tag already exists")
	ErrNotExistTag = New(10002, "tag not exists")
)
