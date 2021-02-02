package rbac

import (
	"github.com/shipengqi/example.v1/blog/dao"
)

type Interface interface {
	AddUser(name, pass, phone, email string) error
}

type rbac struct {
	dao dao.Interface
}

func New(d dao.Interface) Interface {
	return &rbac{dao: d}
}

func (r *rbac) AddUser(name, pass, phone, email string) error {
	return r.dao.AddUser(name, pass, phone, email)
}