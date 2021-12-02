package main

import (
	"encoding/json"
	"fmt"
)

type IPAM struct {
	Subnets             *map[string]*BitMap
}

func main() {
	bitmap := NewBitMap(256)

	bytes, err := json.Marshal(bitmap)
	if err != nil {
		fmt.Println("Marshal: ", err)
		return
	}
	fmt.Println(bytes)

	bitmap2 := NewBitMap(256)
	err = json.Unmarshal(bytes, bitmap2)
	if err != nil {
		fmt.Println("Unmarshal: ", err)
		return
	}
	bitmap2.Count = 10
	fmt.Println(bitmap2)

	s := &map[string]*BitMap{}
	(*s)["test"] = bitmap2
	fmt.Println((*s)["test"])
	fmt.Println(s)
}

// Output:
// [123 34 98 105 116 115 34 58 34 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 65 34 44 34 99 111 117 110 116 34 58 48 44 34 99 97 112 34 58 50 53 54 125]
// &{[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 10 256}
// &{[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 10 256}
