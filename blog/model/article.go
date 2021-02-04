package model

type Article struct {
	Model

	TagID         uint   `json:"tag_id" gorm:"index"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
}

// TableName overwrite table name to `blog_article`
func (a *Article) TableName() string {
	return "blog_article"
}
