package rbac

import (
	"github.com/shipengqi/example.v1/apps/blog/dao"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	"github.com/shipengqi/example.v1/apps/blog/pkg/utils"
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
	encoded := utils.EncodeMD5WithSalt(pass, setting.AppSettings().Salt)
	return r.dao.AddUser(name, encoded, phone, email)
}
