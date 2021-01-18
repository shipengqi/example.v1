package errno

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	ERRNO_ENV_KEY_PRINT_STACK = "IS_PRINT_STACK"
	ERRNO_ENV_KEY_STACK_SKIP  = "ERROR_STACK_SKIP"
)

// register codes.
var (
	_codes       = make(map[int]struct{})

	stackSkip    = 3
	isPrintStack bool
)

func init() {
	if os.Getenv(ERRNO_ENV_KEY_PRINT_STACK) == "true" {
		isPrintStack = true
	}
	s := os.Getenv(ERRNO_ENV_KEY_STACK_SKIP)
	if s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			stackSkip = i
		}
	}
}

// New new a errno.Codes by int value.
func New(e int, msg string) Errno {
	if e <= 0 {
		panic("code must be greater than zero")
	}
	return add(e, msg)
}

func add(e int, msg string) Errno {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("code: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return &Code{
		code:    e,
		message: msg,
		err:     nil,
	}
}

type Errno interface {
	// sometimes Error return Code and message
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	// RawError return the origin error.
	RawError() error
	// Stack return the error stack info.
	Stack() string
}

type stack struct {
	data    string
	callers []uintptr
}

type Code struct {
	stack

	code    int
	message string
	err     error
}

func (e *Code) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("code: ")
	buffer.WriteString(strconv.FormatInt(int64(e.Code()), 10))
	buffer.WriteString(", msg: ")
	buffer.WriteString(e.Message())
	return buffer.String()
}

// RawError return the origin error
func (e *Code) RawError() error {
	return e.err
}

// Code return error code
func (e *Code) Code() int { return e.code }

// Message return error message
func (e *Code) Message() string {
	return e.message
}

// Stack return the function call stack
func (e *Code) Stack() string {
	return e.data
}

func (e *Code) genStackTrace() string {
	if !isPrintStack {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("call stack:\n")
	st := make([]uintptr, 32)
	// skip the first {skip} invocations
	count := runtime.Callers(stackSkip, st)
	e.callers = st[:count]
	frames := runtime.CallersFrames(e.callers)
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}
		if !strings.Contains(frame.File, "runtime/") {
			buffer.WriteString(fmt.Sprintf("%s\n\t%s:%d\n",
				frame.Func.Name(), frame.File, frame.Line))
		}
	}
	e.data = buffer.String()
	return e.data
}

// Wrap Wrap error
func Wrap(err error, msg string) Errno {
	return wrapErr(err, msg)
}

// Wrapf Wrap error
func Wrapf(err error, args ...interface{}) Errno {
	return wrapErr(err, args...)
}

// String parse code string to Errno.
func String(e string) Errno {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return &Code{
			code:    ErrInternalServer.Code(),
			message: e,
		}
	}
	return &Code{code: i}
}

// Cause cause from error to Errno.
func Cause(e error) Errno {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Errno)
	if ok {
		return ec
	}
	return String(e.Error())
}

// Is returns true if error is a Errno error
func Is(code Errno, err error) bool {
	return Cause(err).Code() == code.Code()
}

// IsNothingFoundError returns true if error contains a IsNothingFoundError error
func IsNothingFoundError(err error) bool {
	if err, ok := err.(Errno); ok {
		if err == ErrNothingFound {
			return true
		}
	}
	return err == ErrNothingFound
}

func SetPrintStack(isPrint bool) {
	isPrintStack = isPrint
}

func SetErrorStackSkip(skip int) {
	stackSkip = skip
}

func wrapErr(err error, args ...interface{}) Errno {
	msg := normalize(args...)
	if err == nil {
		err = errors.New(msg)
	}
	e, ok := err.(*Code)
	if !ok {
		e = &Code{
			code:    ErrInternalServer.Code(),
			message: msg,
			err:     err,
		}
	}
	e.genStackTrace()
	if e.message == "" {
		e.message = err.Error()
	}
	return e
}

// fmtErrMsg used to format error message
func normalize(args ...interface{}) string {
	if len(args) > 1 {
		return fmt.Sprintf(args[0].(string), args[1:]...)
	}
	if len(args) == 1 {
		if v, ok := args[0].(string); ok {
			return v
		}
		if v, ok := args[0].(error); ok {
			return v.Error()
		}
	}
	return ""
}
