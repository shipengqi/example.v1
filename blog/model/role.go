package model

type Role struct {
	Model

	Name        string `json:"name"`
	Description string `json:"description"`
	Locked      bool   `json:"locked"`
}

// TableName overwrite table name `roles` to `blog_role`
func (r *Role) TableName() string {
	return "blog_role"
}
