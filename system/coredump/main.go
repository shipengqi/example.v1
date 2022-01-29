package main

import "math/rand"

// https://xie.infoq.cn/article/470960b2a41d7b33ff9fe0d6e
// GOTRACEBACK 控制程序崩溃时输出的详细程度。 它可以采用不同的值：
//
//
//
// none 不显示任何 goroutine 栈 trace。
// single, 默认选项，显示当前 goroutine 栈 trace。
// all 显示所有用户创建的 goroutine 栈 trace。
// system 显示所有 goroutine 栈 trace,甚至运行时的 trace。
// crash 类似 system, 而且还会生成 core dump。
//
// core dump 是通过 SIGABRT 信号触发。
// core dump 可以通过 delve 或者 gdb。

func main() {
	var sum int
	for {
		n := rand.Intn(1e6)
		sum += n
		if sum % 42 == 0 {
			panic("panic for GOTRACEBACK")
		}
	}
}

// 该程序将很快崩溃
// panic: panic for GOTRACEBACK
//
// goroutine 1 [running]:
// main.main()
//	C:/Code/example.v1/system/coredump/main.go:21 +0x78
//
// 我们无法从栈 trace 中分辨出崩溃所涉及的值。增加日志或许是一种解决方案，但是我们并不总是知道在何处添加日志。

// 添加环境变量 GOTRACEBACK=crash 再运行它。由于我们现在已打印出所有 goroutine，包括运行时，因此输出更加详细。 并输出 core dump：
// GOROOT=C:\Program Files\Go #gosetup
// GOPATH=C:\Code\gowork #gosetup
// "C:\Program Files\Go\bin\go.exe" build -o C:\Users\shipeng\AppData\Local\Temp\GoLand\___1go_build_github_com_shipengqi_example_v1_system_coredump.exe github.com/shipengqi/example.v1/system/coredump #gosetup
// C:\Users\shipeng\AppData\Local\Temp\GoLand\___1go_build_github_com_shipengqi_example_v1_system_coredump.exe #gosetup
// panic: panic for GOTRACEBACK
//
// goroutine 1 [running]:
// panic({0x4408c0, 0x45e5f8})
//	C:/Program Files/Go/src/runtime/panic.go:1147 +0x3a8 fp=0xc000047f58 sp=0xc000047e98 pc=0x40ea08
// main.main()
//	C:/Code/example.v1/system/coredump/main.go:21 +0x78 fp=0xc000047f80 sp=0xc000047f58 pc=0x43be58
// runtime.main()
//	C:/Program Files/Go/src/runtime/proc.go:255 +0x217 fp=0xc000047fe0 sp=0xc000047f80 pc=0x411437
// runtime.goexit()
//	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000047fe8 sp=0xc000047fe0 pc=0x435921
//
// goroutine 2 [force gc (idle)]:
// runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
//	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000043fb0 sp=0xc000043f90 pc=0x4117d6
// runtime.goparkunlock(...)
//	C:/Program Files/Go/src/runtime/proc.go:372
// runtime.forcegchelper()
//	C:/Program Files/Go/src/runtime/proc.go:306 +0xb1 fp=0xc000043fe0 sp=0xc000043fb0 pc=0x411671
// runtime.goexit()
//	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000043fe8 sp=0xc000043fe0 pc=0x435921
// created by runtime.init.7
//	C:/Program Files/Go/src/runtime/proc.go:294 +0x25
//
// goroutine 3 [GC sweep wait]:
// runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
//	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000045fb0 sp=0xc000045f90 pc=0x4117d6
// runtime.goparkunlock(...)
//	C:/Program Files/Go/src/runtime/proc.go:372
// runtime.bgsweep()
//	C:/Program Files/Go/src/runtime/mgcsweep.go:163 +0x88 fp=0xc000045fe0 sp=0xc000045fb0 pc=0x3fc7e8
// runtime.goexit()
//	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000045fe8 sp=0xc000045fe0 pc=0x435921
// created by runtime.gcenable
//	C:/Program Files/Go/src/runtime/mgc.go:181 +0x55
//
// goroutine 4 [GC scavenge wait]:
// runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
//	C:/Program Files/Go/src/runtime/proc.go:366 +0xd6 fp=0xc000055f80 sp=0xc000055f60 pc=0x4117d6
// runtime.goparkunlock(...)
//	C:/Program Files/Go/src/runtime/proc.go:372
// runtime.bgscavenge()
//	C:/Program Files/Go/src/runtime/mgcscavenge.go:265 +0xcd fp=0xc000055fe0 sp=0xc000055f80 pc=0x3fa8ed
// runtime.goexit()
//	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000055fe8 sp=0xc000055fe0 pc=0x435921
// created by runtime.gcenable
//	C:/Program Files/Go/src/runtime/mgc.go:182 +0x65
//
//
// Delve
// Delve 是使用 Go 编写的 Go 程序的调试器。 它可以通过在用户代码以及运行时中的任意位置添加断点来逐步调试，甚至可以使用以二进制文件
// 和 core dump 为参数的命令 dlv core 调试 core dump。
