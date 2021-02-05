package model

type User struct {
	Model

	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Locked   bool   `json:"locked,omitempty"`
}

// TableName overwrite table name `users` to `blog_user`
func (u *User) TableName() string {
	return "blog_user"
}
