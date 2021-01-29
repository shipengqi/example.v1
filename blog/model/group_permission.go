package model

type GroupPermission struct {
	Model

	GroupId      int
	PermissionId int
}

// TableName overwrite table name to `blog_group_permission`
func (g *GroupPermission) TableName() string {
	return "blog_group_permission"
}
