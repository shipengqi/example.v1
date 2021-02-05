package model

type Group struct {
	Model

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Locked      bool   `json:"locked,omitempty"`
}

// TableName overwrite table name `groups` to `blog_group`
func (g *Group) TableName() string {
	return "blog_group"
}
