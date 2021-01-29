package model

type Operation struct {
	Model

	Name        string
	Description string
}

// TableName overwrite table name to `blog_operation`
func (o *Operation) TableName() string {
	return "blog_operation"
}
