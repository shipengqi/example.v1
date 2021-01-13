package model

type User struct {
	Model

	Username string
	Password string
	Sex      string
	Phone    string
	Email    string
}

// TableName overwrite table name `users` to `blog_user`
func (u *User) TableName() string {
	return "blog_user"
}
