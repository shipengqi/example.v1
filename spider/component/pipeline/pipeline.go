package pipeline

import (
	"github.com/shipengqi/example.v1/spider/component"
	"github.com/shipengqi/example.v1/spider/component/stub"
)

type PipelineImpl struct {
	stub.ComponentInternal

	itemProcessors []component.ProcessItem

	failFast bool
}

func New() {}
