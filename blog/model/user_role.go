package model

type UserRole struct {
	Model

	UserId int
	RoleId int
}

// TableName overwrite table name `user_roles` to `blog_user_role`
func (u *UserRole) TableName() string {
	return "blog_user_role"
}
