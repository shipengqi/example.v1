# CGO

如果在 Go 代码中出现了 `import "C"` 语句则表示使用了 CGO 特性，紧跟在这行语句前面的注释是一种特殊语法，里面包含的是正常的C语言代码。
当确保 CGO 启用的情况下，还可以在当前目录中包含 C/C++ 对应的源文件。

Go 是强类型语言，所以 **cgo 中传递的参数类型必须与声明的类型完全一致，而且传递前必须用 ”C” 中的转化函数转换成对应的 C 类型，不能直接传入 Go 中类
型的变量**。同时通过虚拟的 C 包导入的 C 语言符号并不需要是大写字母开头，它们不受 Go 语言的导出规则约束。

https://chai2010.cn/advanced-go-programming-book/ch2-cgo/ch2-03-cgo-types.html
https://www.jianshu.com/p/c48e6cd84ff5
https://github.com/xianlubird/mydocker
https://zhuanlan.zhihu.com/p/23456448
```go
cmd := &exec.Cmd{
Path:   "/proc/self/exe",
Args:   []string{"setns"},
}
cmd.Start()
```