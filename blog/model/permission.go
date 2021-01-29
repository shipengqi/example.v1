package model

type Permission struct {
	Model

	ResourceId  int
	OperationId int
	Name        string
	Description string
}

// TableName overwrite table name `permissions` to `blog_role_permission`
func (p *Permission) TableName() string {
	return "blog_permission"
}
