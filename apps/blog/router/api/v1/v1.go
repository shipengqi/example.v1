package v1

import (
	"github.com/shipengqi/example.v1/apps/blog/service"
)

var svc *service.Service

func Init(s *service.Service)  {
	svc = s
}
