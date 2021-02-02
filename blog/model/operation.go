package model

type Operation struct {
	Model

	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName overwrite table name to `blog_operation`
func (o *Operation) TableName() string {
	return "blog_operation"
}
