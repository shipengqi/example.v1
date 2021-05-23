package model



type UserRBAC struct {
	U      User
	Groups []Group
	Roles  []Role
}

type UserPermission struct {
	URL    string
	Method string
}
