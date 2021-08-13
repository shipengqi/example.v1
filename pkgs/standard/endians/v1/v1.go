package main

import (
	"encoding/binary"
	"fmt"
)

// 网络字节序有两种形式表示：
//
// 小端序：低位字节在前，高位字节在后，计算机方便处理。
// 大端序：高位字节在前，低位字节在后，与人的使用习惯相同。
func main() {
	a := uint32(0x01020304)
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, a)
	fmt.Println(arr)

	binary.LittleEndian.PutUint32(arr, a)
	fmt.Println(arr)
}

// Output:
// [1 2 3 4]
// [4 3 2 1]
