package conv

type User struct {
	Name string
}

func (u *User) String() string {
	return u.Name
}


