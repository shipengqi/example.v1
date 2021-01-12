package model

type Model struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"` // GORM use CreatedAt, UpdatedAt to track creating/updating time by convention,
	// and GORM will set the current time when creating/updating if the fields are defined
	UpdatedAt int `json:"updated_at"`
	DeletedAt int `json:"deleted_at"` // not sure
}
