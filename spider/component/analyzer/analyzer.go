package analyzer

import (
	"github.com/shipengqi/example.v1/spider/component"
	"github.com/shipengqi/example.v1/spider/component/stub"
)

type AnalyzerImpl struct {
	stub.ComponentInternal

	resParsers []component.ParseResponse
}

func New() {}
