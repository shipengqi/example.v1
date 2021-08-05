package main

import (
	"encoding/binary"
	"fmt"
)

// 序列化的核心原理
// 下面的实现的几个特点：
//
// 顺序写，效率高。
// 不会把属性key序列化进去。
// 对[0 0 0 0 0 0 0 5]这样占用字节数多，但是值（5）很小的情况下有点浪费。
// 实现复杂。

func main() {
	var pkt = struct {
		Source   uint32
		Sequence uint64
		Data     []byte
	}{
		Source:   257,
		Sequence: 5,
		Data:     []byte("hello world"),
	}

	// 为了方便观看，使用大端序
	endian := binary.BigEndian

	buf := make([]byte, 1024) // buffer
	i := 0
	endian.PutUint32(buf[i:i+4], pkt.Source)
	i += 4
	endian.PutUint64(buf[i:i+8], pkt.Sequence)
	i += 8
	// 由于data长度不确定，必须先把长度写入buf, 这样在反序列化时就可以正确的解析出data
	dataLen := len(pkt.Data)
	endian.PutUint32(buf[i:i+4], uint32(dataLen))
	i += 4
	// 写入数据data
	copy(buf[i:i+dataLen], pkt.Data)
	i += dataLen
	fmt.Println(buf[0:i])
	fmt.Println("length", i)


	// 反序列化
	var output = struct {
		Source   uint32
		Sequence uint64
		Data     []byte
	}{}
	j := 0
	output.Source = endian.Uint32(buf[j : j+4])
	j += 4
	output.Sequence = endian.Uint64(buf[j : j+8])
	j += 8
	dataLen2 := endian.Uint32(buf[j : j+4])
	j += 4
	output.Data = make([]byte, dataLen2)
	copy(output.Data, buf[j:j+int(dataLen2)])
	fmt.Printf("Src:%d Seq:%d Data:%s\n", output.Source, output.Sequence, output.Data) // Src:257 Seq:5 Data:hello world
}

// Output:
// [0 0 1 1 0 0 0 0 0 0 0 5 0 0 0 11 104 101 108 108 111 32 119 111 114 108 100]
// length 27
//

// [0 0 1 1 0 0 0 0 0 0 0 5 0 0 0 11 104 101 108 108 111 32 119 111 114 108 100]
// 0 0 1 1 -> 257                 4 个字节
// 0 0 0 0 0 0 0 5 -> 5           8 个字节
// 0 0 0 11 -> 11                 4 个字节
// 104 101 108 108 111 32 119 111 114 108 100 -> hello world  11 个字节
// 一共 27 个字节
