package model

type Permission struct {
	Model

	ResourceId  int    `json:"resource_id"`
	OperationId int    `json:"operation_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName overwrite table name `permissions` to `blog_role_permission`
func (p *Permission) TableName() string {
	return "blog_permission"
}
