package model

type Model struct {
	ID        uint `gorm:"primary_key;autoIncrement" json:"id"`
	Deleted   bool `json:"deleted"`
	CreatedAt uint `json:"created_at"` // GORM use CreatedAt, UpdatedAt to track creating/updating time by convention,
	// and GORM will set the current time when creating/updating if the fields are defined
	UpdatedAt uint `json:"updated_at"`
	DeletedAt uint `json:"deleted_at"` // not sure
}

type UserRBAC struct {
	Groups []Group
	Roles  []Role
}

type UserPermission struct {
	URL    string
	Method string
}
