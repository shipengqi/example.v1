package model

type GroupUser struct {
	Model

	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
}

// TableName overwrite table name to `blog_group_user`
func (g *GroupUser) TableName() string {
	return "blog_group_user"
}
