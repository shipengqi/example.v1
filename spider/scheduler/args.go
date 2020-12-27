package scheduler

import "github.com/shipengqi/example.v1/spider/component"

// Args 代表参数容器的接口类型
type Args interface {
	// Check 用于自检参数的有效性
	// 若结果值为 nil，则说明未发现问题，否则就意味着自检未通过
	Check() error
}

// RequestArgs 代表请求相关的参数容器的类型。
type RequestArgs struct {
	// AcceptedDomains 代表可以接受的 URL 的主域名的列表
	// URL 主域名不在列表中的请求都会被忽略
	AcceptedDomains []string `json:"accepted_primary_domains"`
	// maxDepth 代表了需要被爬取的最大深度
	// 实际深度大于此值的请求都会被忽略
	MaxDepth uint32 `json:"max_depth"`
}

func (r *RequestArgs) Check() error {
	if r.AcceptedDomains == nil {
		return nil
	}
	return nil
}

// DataArgs 代表数据相关的参数容器的类型
type DataArgs struct {
	// ReqBufferCap 代表请求缓冲器的容量
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// ReqMaxBufferNumber 代表请求缓冲器的最大数量
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// RespBufferCap 代表响应缓冲器的容量
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// RespMaxBufferNumber 代表响应缓冲器的最大数量
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// ItemBufferCap 代表条目缓冲器的容量
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// ItemMaxBufferNumber 代表条目缓冲器的最大数量
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// ErrorBufferCap 代表错误缓冲器的容量
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	// ErrorMaxBufferNumber 代表错误缓冲器的最大数量
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}

func (d *DataArgs) Check() error {
	return nil
}

// ComponentArgsSummary 代表组件相关的参数容器的摘要类型
type ComponentArgsSummary struct {
	DownloaderListSize int `json:"downloader_list_size"`
	AnalyzerListSize   int `json:"analyzer_list_size"`
	PipelineListSize   int `json:"pipeline_list_size"`
}

// ComponentArgs 代表组件相关的参数容器的类型
type ComponentArgs struct {
	// Downloaders 代表下载器列表。
	Downloaders []component.Downloader
	// Analyzers 代表分析器列表。
	Analyzers []component.Analyzer
	// Pipelines 代表条目处理管道管道列表。
	Pipelines []component.Pipeline
}

func (c *ComponentArgs) Check() error {
	return nil
}