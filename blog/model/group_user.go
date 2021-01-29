package model

type GroupUser struct {
	Model

	GroupId int
	UserId  int
}

// TableName overwrite table name to `blog_group_user`
func (g *GroupUser) TableName() string {
	return "blog_group_user"
}
