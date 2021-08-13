package app

import (
	"github.com/astaxie/beego/validation"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Info().Msgf("%s | %s", err.Key, err.Message)
	}

	return
}
