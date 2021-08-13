package model

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	Disabled   bool   `json:"disabled"`
}

// TableName overwrite table name `blogs` to `blog_tag`
func (tag *Tag) TableName() string {
	return "blog_tag"
}
