package main

import "fmt"

type user struct {
	name string
	age  int
}

func main() {
	u := user{"pooky", 18}
	//Printf 格式化输出
	fmt.Printf("%+v\n", u)       // 格式化输出结构
	fmt.Printf("%#v\n", u)       // 输出值的 Go 语言表示方法
	fmt.Printf("%T\n", u)        // 输出值的类型的 Go 语言表示
	fmt.Printf("%t\n", true)     // 输出值的 true 或 false
	fmt.Printf("%b\n", 1024)     // 二进制表示
	fmt.Printf("%c\n", 11111111) // 数值对应的 Unicode 编码字符
	fmt.Printf("%d\n", 10)       // 十进制表示
	fmt.Printf("%o\n", 8)        // 八进制表示
	fmt.Printf("%q\n", 22)       // 转化为十六进制并附上单引号
	fmt.Printf("%x\n", 1223)     // 十六进制表示，用 a-f 表示
	fmt.Printf("%X\n", 1223)     // 十六进制表示，用 A-F 表示
	fmt.Printf("%U\n", 1233)     // Unicode 表示
	fmt.Printf("%b\n", 12.34)    // 无小数部分，两位指数的科学计数法6946802425218990p-49
	fmt.Printf("%e\n", 12.345)   // 科学计数法，e表示
	fmt.Printf("%E\n", 12.34455) // 科学计数法，E表示
	fmt.Printf("%f\n", 12.3456)  // 有小数部分，无指数部分
	fmt.Printf("%g\n", 12.3456)  // 根据实际情况采用%e或%f输出
	fmt.Printf("%G\n", 12.3456)  // 根据实际情况采用%E或%f输出
	fmt.Printf("%s\n", "wqdew")  // 直接输出字符串或者[]byte
	fmt.Printf("%q\n", "dedede") // 双引号括起来的字符串
	fmt.Printf("%x\n", "abczxc") // 每个字节用两字节十六进制表示，a-f 表示
	fmt.Printf("%X\n", "asdzxc") // 每个字节用两字节十六进制表示，A-F 表示
	fmt.Printf("%p\n", &u)           // 0x 开头的十六进制数表示，指针
	fmt.Println(&u)
	fmt.Println(u)
	widthIdentifierDemo()
	otherIdentifierDemo()
}

// Output:
//{name:pooky age:18}
//main.user{name:"pooky", age:18}
//main.user
//true
//10000000000
//�
//10
//10
//'\x16'
//4c7
//4C7
//U+04D1
//6946802425218990p-49
//1.234500e+01
//1.234455E+01
//12.345600
//12.3456
//12.3456
//wqdew
//"dedede"
//6162637a7863
//6173647A7863
//0xc000096440
//&{pooky 18}
//{pooky 18}


// 宽度通过一个紧跟在百分号后面的十进制数指定，如果未指定宽度，则表示值时除必需之外不作填充。精度通过（可选的）宽度后跟点号后跟的十进制数指定。
// 如果未指定精度，会使用默认精度；如果点号后没有跟数字，表示精度为0。
// %e 和 %f 的默认精度为 6。
// 示例：
// %f，   默认宽度和精度
// %9f，  宽度为 9，默认精度
// %.2f， 默认宽度，精度为 2
// %9.2f，宽度为 9，精度为 2
// %9.f， 宽度为 9，精度为 0

func widthIdentifierDemo()  {
	n := 12.34
	fmt.Println("****************************")
	fmt.Printf("%f\n", n)
	fmt.Printf("%9f\n", n)
	fmt.Printf("%.2f\n", n)
	fmt.Printf("%9.2f\n", n)
	fmt.Printf("%9.f\n", n)
	fmt.Printf("%10f\n", n)
	fmt.Printf("%.10f\n", n)
	fmt.Printf("%10.10f\n", n)
	fmt.Printf("%10.f\n", n)
	fmt.Printf("%10s\n", "a")                    // 字符串最小宽度为 10
	fmt.Printf("%-10s\n", "a")                   // 字符串最小宽度为 10（左对齐）
	fmt.Printf("%.5s\n", "0123456789")           // 字符串最大宽度为 5
	fmt.Printf("%5.7s\n","01234567890123")       // 最小宽度为5，最大宽度为7
	fmt.Printf("%-5.7s\n","01234567890123")      // 最小宽度为5，最大宽度为7（左对齐）
	fmt.Printf("%.3s\n","01234567890123")       // 如果宽度大于3，则截断
	fmt.Printf("%05s\n","aa")                    // 如果宽度小于5，就会在字符串前面补零
}

// Output:
//12.340000
//12.340000
//12.34
//    12.34
//       12
// 12.340000
//12.3400000000
//12.3400000000
//        12
//          a
// a
// 01234
// 0123456
// 0123456
//  012
// 000aa

// + 总打印数值的正负号； 对于 %q（%+q）保证只输出 ASCII 编码的字符	Printf("%+q", "中文")	"\u4e2d\u6587"
//
// - 在右侧而非左侧填充空格（左对齐该区域）
//
// # 备用格式：为八进制添加前导 0（%#o）；为十六进制添加前导 0x（%#x）或 0X（%#X），为 %p（%#p）去掉前导 0x；
// 如果可能的话，%q（%#q）会打印原始 （即反引号围绕的）字符串；
// 如果是可打印字符，%U（%#U）会写出该字符的；
// Unicode 编码形式（如字符 x 会被打印成 U+0078 'x'）；	Printf("%#U", '中')	U+4E2D
//
// ' ' (空格)为数值中省略的正负号留出空白（% d）；以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开
//
// 0 使用 0 填充为而不是空格，对于数值类型，会把填充的 0 放到正负号后面

func otherIdentifierDemo()  {
	s := "hi"
	fmt.Println("****************************")
	fmt.Printf("%s\n", s)
	fmt.Printf("%5s\n", s)
	fmt.Printf("%-5s\n", s)
	fmt.Printf("%5.7s\n", s)
	fmt.Printf("%-5.7s\n", s)
	fmt.Printf("%5.2s\n", s)
	fmt.Printf("%05s\n", s)
}

// Output:
//hi
//   hi
//hi
//   hi
//hi
//   hi
//000hi
