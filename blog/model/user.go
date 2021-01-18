package model

type User struct {
	Model

	Username string
	Password string
	Phone    string
	Email    string
	Sex      uint8
}

// TableName overwrite table name `users` to `blog_user`
func (u *User) TableName() string {
	return "blog_user"
}
