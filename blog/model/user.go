package model

type User struct {
	Model

	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Locked   bool   `json:"locked"`
}

// TableName overwrite table name `users` to `blog_user`
func (u *User) TableName() string {
	return "blog_user"
}
