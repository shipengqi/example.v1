package model

type Resource struct {
	Model

	Name        string
	Description string
}

// TableName overwrite table name to `blog_resource`
func (r *Resource) TableName() string {
	return "blog_resource"
}
