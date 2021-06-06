package main

import (
	"fmt"
	"strings"
)

// SplitN slices s into substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings

func main()  {
	test1 := "str1.str2.str3.str4"

	result1 := strings.SplitN(test1, ".", 1)
	fmt.Println(result1) // [str1.str2.str3.str4]
	fmt.Println("")

	result2 := strings.SplitN(test1, ".", 0)
	fmt.Println(result2) // []
	fmt.Println("")

	result3 := strings.SplitN(test1, ".", 2)
	fmt.Println(result3) // [str1 str2.str3.str4]
	fmt.Println("")

	result4 := strings.SplitN(test1, ".", -1)
	fmt.Println(result4) // [str1 str2 str3 str4]
	fmt.Println("")

	result5 := strings.SplitN(test1, ".", 3)
	fmt.Println(result5) // [str1 str2 str3.str4]
	fmt.Println("")

	result6 := strings.SplitN(test1, ".", 4)
	fmt.Println(result6) // [str1 str2 str3 str4]
	fmt.Println("")

	result7 := strings.SplitN(test1, ".", 5)
	fmt.Println(result7) // [str1 str2 str3 str4]
	fmt.Println("")
}
