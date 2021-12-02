package errors

import (
	"bytes"
	"strings"
)

type ErrorType string

type CrawlerError interface {
	Type() ErrorType
	Error() string
}
type CrawlerErrorImpl struct {
	errType ErrorType
	msg     string
	fullMsg string
}

const (
	ERROR_TYPE_DOWNLOAD ErrorType = "download err"
	ERROR_TYPE_ANALYSIS ErrorType = "analysis err"
	ERROR_TYPE_PIPELINE ErrorType = "pipeline err"
	ERROR_TYPE_SCHEDULE ErrorType = "schedule err"
)

func NewCrawlerError(t ErrorType, msg string) *CrawlerErrorImpl {
	return &CrawlerErrorImpl{
		errType: t,
		msg:     strings.TrimSpace(msg),
	}
}

func (c *CrawlerErrorImpl) Type() ErrorType {
	return c.errType
}

func (c *CrawlerErrorImpl) Error() string {
	if c.fullMsg == "" {
		c.genFullMsg()
	}

	return c.fullMsg
}


func (c *CrawlerErrorImpl) genFullMsg() string {
	var buf bytes.Buffer
	buf.WriteString("crawler err: ")

	if c.errType != "" {
		buf.WriteString(string(c.errType))
		buf.WriteString(": ")
	}

	buf.WriteString(c.msg)
	c.fullMsg = buf.String()
	return c.fullMsg
}