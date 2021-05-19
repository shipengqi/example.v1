package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

const reset = "\033[0m"

func main() {
	defer fmt.Println("in main")
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[Recover]: \n%s\n%s%s\n", err, stack(3), reset)
		}
	}()

	var t *int
	*t = 10
}

func stack(skip int) string {
	var buffer bytes.Buffer
	buffer.WriteString("stack:\n")
	st := make([]uintptr, 32)
	// skip the first {skip} invocations
	count := runtime.Callers(skip, st)
	callers := st[:count]
	frames := runtime.CallersFrames(callers)
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
	return buffer.String()
}

// Output:
// recover:  runtime error: invalid memory address or nil pointer dereference
// in main
