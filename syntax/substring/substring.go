package substring

import (
	"unicode/utf8"

	"github.com/thinkeridea/go-extend/exunicode/exutf8"
)

// 字符串截取的几种方式

// 通过 go 内置的 slice 语法截取字符串
// 这种方式可以完美处理 ASCII 单字节字符串的截取，
// 但是中文一般会占用多个字节，utf8 是 3 个字节，如果使用 slice 语法截取字符串，获得的是乱码
func SliceSubString(s string, start, length int) string {
	return s[start:length]
}

// 使用 []rune 类型转换后，再按切片语法截取
// 因为类型转换带来了内存分配，产生新的字符串，性能较差
func RuneSubString(s string, start, length int) string {
	rs := []rune(s)
	return string(rs[start:length])
}

// utf8.DecodeRuneInString 可以转换单个字符，并给出字符占用的字节数
// 这种方式的性能比 RuneSubString 要快好几倍
func DecodeRuneInStringSubString(s string, start, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}

	return s[:n]
}

// range 是按字符迭代的，并不是字节
func RangeSubString(s string, start, length int) string {
	var i, n int
	for i = range s { // 这里 i 是字节下标
		if n == length {
			break
		}
		n++
	}

	return s[:i]
}

// 使用 go-extend 扩展库 exutf8.RuneIndexInString
func ExtRuneIndexInStringSubString(s string, start, length int) string  {
	n, _ := exutf8.RuneIndexInString(s, length)
	return s[:n]
}


// 使用 go-extend 扩展库 exutf8.RuneSubString
// 这是最易用的方式
func ExtRuneSubString(s string, start, length int) string  {
	return exutf8.RuneSubString(s, start, length)
}
