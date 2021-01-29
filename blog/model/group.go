package model

type Group struct {
	Model

	Name        string
	Description string
}

// TableName overwrite table name `groups` to `blog_group`
func (g *Group) TableName() string {
	return "blog_group"
}