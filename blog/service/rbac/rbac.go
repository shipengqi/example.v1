package rbac

import (
	"github.com/shipengqi/example.v1/blog/dao"
)

type Interface interface {

}

type rbac struct {
	dao dao.Interface
}

func New(d dao.Interface) Interface {
	return &rbac{dao: d}
}