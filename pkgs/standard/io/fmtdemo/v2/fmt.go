package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	_, _ = fmt.Fprintf(buf, "\t%s: %s\n", "test1", "test1")
	_, _ = fmt.Fprintf(buf, "\t%s: %s\n", "test2", "test2")
	fmt.Println(buf.String())
}

// 	test1: test1
//	test2: test2
