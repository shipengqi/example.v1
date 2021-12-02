package model

type Resource struct {
	Model

	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// TableName overwrite table name to `blog_resource`
func (r *Resource) TableName() string {
	return "blog_resource"
}
