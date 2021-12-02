package main

// C 代码注释和 import "C" 之间不能有空行

/*
int add(int a, int b) {
    return a + b;
}
*/
import "C"
import "fmt"
func main()  {
	a := C.int(1)
	b := C.int(2)
	value := C.add(a, b)
	fmt.Printf("value: %d\n", int(value))
}

// Output:
// value: 3