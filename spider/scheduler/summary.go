package scheduler

import "github.com/shipengqi/example.v1/spider/component"

// SchedSummary 代表调度器摘要的接口类型
type SchedSummary interface {
	// Struct 用于获得摘要信息的结构化形式
	Struct() SummaryStruct
	// String 用于获得摘要信息的字符串形式
	String() string
}

// SummaryStruct 代表调度器摘要的结构
type SummaryStruct struct {
	RequestArgs     RequestArgs               `json:"request_args"`
	DataArgs        DataArgs                  `json:"data_args"`
	ComponentArgs   ComponentArgsSummary      `json:"component_args"`
	Status          string                    `json:"status"`
	Downloaders     []component.SummaryStruct `json:"downloaders"`
	Analyzers       []component.SummaryStruct `json:"analyzers"`
	Pipelines       []component.SummaryStruct `json:"pipelines"`
	ReqBufferPool   BufferPoolSummaryStruct   `json:"request_buffer_pool"`
	RespBufferPool  BufferPoolSummaryStruct   `json:"response_buffer_pool"`
	ItemBufferPool  BufferPoolSummaryStruct   `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct   `json:"error_buffer_pool"`
	NumURL          uint64                    `json:"url_number"`
}

// BufferPoolSummaryStruct 代表缓冲池的摘要类型
type BufferPoolSummaryStruct struct {
	BufferCap       uint32 `json:"buffer_cap"`
	MaxBufferNumber uint32 `json:"max_buffer_number"`
	BufferNumber    uint32 `json:"buffer_number"`
	Total           uint64 `json:"total"`
}
