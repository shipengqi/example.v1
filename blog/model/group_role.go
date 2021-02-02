package model

type GroupRoles struct {
	Model

	GroupId int `json:"group_id"`
	RoleId  int `json:"role_id"`
}

// TableName overwrite table name to `blog_group_role`
func (g *GroupRoles) TableName() string {
	return "blog_group_role"
}
