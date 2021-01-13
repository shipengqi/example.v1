package model

type Role struct {
	Model

	Name        string
	Role        string
	Description string
}

// TableName overwrite table name `roles` to `blog_role`
func (r *Role) TableName() string {
	return "blog_role"
}
