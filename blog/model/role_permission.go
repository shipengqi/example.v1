package model

type RolePermission struct {
	Model

	RoleId       int
	PermissionId int
}

// TableName overwrite table name `role_permissions` to `blog_role_permission`
func (r *RolePermission) TableName() string {
	return "blog_role_permission"
}
