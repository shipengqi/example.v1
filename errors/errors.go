package errors

import (
	"bytes"
	"strings"
)

const (
	ERROR_PACKAGE_DEFAULT ErrorPackage = "example.v1"
)

const (
	ERROR_TYPE_DEFAULT ErrorType = "example.v1 error"
)

// ErrorPackage 代表错误所属 package
type ErrorPackage string

// ErrorType 代表错误类型
type ErrorType string

type V1Error interface {
	// Package 用于获得错误所属 package
	Package() ErrorPackage
	// Type 用于获得错误的类型
	Type() ErrorType
	// Error 用于获得错误提示信息
	Error() string
}

type errorImpl struct {
	pkg     ErrorPackage
	errType ErrorType
	msg     string
	fullMsg string
}

func NewV1Error(errType ErrorType, errMsg string) V1Error {
	return &errorImpl{
		pkg:     ERROR_PACKAGE_DEFAULT,
		errType: errType,
		msg:     strings.TrimSpace(errMsg),
	}
}

func NewV1ErrorWithPkg(pkg ErrorPackage, errType ErrorType, errMsg string) V1Error {
	return &errorImpl{
		pkg:     pkg,
		errType: errType,
		msg:     strings.TrimSpace(errMsg),
	}
}

func (e *errorImpl) Package() ErrorPackage {
	return e.pkg
}

func (e *errorImpl) Type() ErrorType {
	return e.errType
}

func (e *errorImpl) Error() string {
	if e.fullMsg == "" {
		e.genFullErrMsg()
	}
	return e.fullMsg
}

func (e *errorImpl) genFullErrMsg() {
	var buffer bytes.Buffer
	if e.pkg != "" {
		buffer.WriteString(string(e.pkg))
		buffer.WriteString(" error: ")
	}

	if e.errType != "" {
		buffer.WriteString(string(e.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(e.msg)
	e.fullMsg = buffer.String()
	return
}

func WrapWithType(errType ErrorType, err error) V1Error {
	return NewV1Error(errType, err.Error())
}

func Wrap(err error) V1Error {
	return NewV1Error(ERROR_TYPE_DEFAULT, err.Error())
}
