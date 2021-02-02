package model

type RolePermission struct {
	Model

	RoleId       int `json:"role_id"`
	PermissionId int `json:"permission_id"`
}

// TableName overwrite table name `role_permissions` to `blog_role_permission`
func (r *RolePermission) TableName() string {
	return "blog_role_permission"
}
