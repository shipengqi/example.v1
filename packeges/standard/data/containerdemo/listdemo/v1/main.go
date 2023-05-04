package main

import (
	"container/list"
	"fmt"
)

func main() {
	l := list.New()
	for i := 0; i < 5; i++ {
		l.PushBack(i)
	}
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println("[1] list e: ", e.Value)
	}
	fmt.Println("")
	fmt.Println("The first element of list: ", l.Front().Value)
	fmt.Println("The last element of list: ", l.Back().Value)
	// insert a element after the first element
	l.InsertAfter(6, l.Front())

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println("[2] list e: ", e.Value)
	}
	fmt.Println("")
	// Swap the header two elements
	l.MoveBefore(l.Front().Next(), l.Front())
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println("[3] list e: ", e.Value)
	}

	fmt.Println("")

	l2 := list.New()
	l2.PushBackList(l) // insert l list to l2 list end
	for e := l2.Front(); e != nil; e = e.Next() {
		fmt.Println("[4] list2 e: ", e.Value)
	}

	fmt.Println("")
	// clean l list
	l.Init()

	fmt.Println("list length: ", l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println("[5] list e: ", e.Value)
	}
}

// Output:
// [1] list e:  0
// [1] list e:  1
// [1] list e:  2
// [1] list e:  3
// [1] list e:  4
//
// The first element of list:  0
// The last element of list:  4
// [2] list e:  0
// [2] list e:  6
// [2] list e:  1
// [2] list e:  2
// [2] list e:  3
// [2] list e:  4
//
// [3] list e:  6
// [3] list e:  0
// [3] list e:  1
// [3] list e:  2
// [3] list e:  3
// [3] list e:  4
//
// [4] list2 e:  6
// [4] list2 e:  0
// [4] list2 e:  1
// [4] list2 e:  2
// [4] list2 e:  3
// [4] list2 e:  4
//
// list length:  0