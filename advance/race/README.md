# 竞争检测

数据竞争是并发系统中最常见，同时也最难处理的 Bug 类型之一。数据竞争会在两个 Go 程并发访问同一个变量， 
且至少有一个访问为写入时产生。

## 数据竞争检测器

Go 内建了数据竞争检测器。要使用它，将 `-race` 标记添加到 go 命令之后：

```bash
go test -race mypkg    // 测试并检测 mypkg 包
go run -race mysrc.go  // 运行并检测 mysrc.go
go build -race mycmd   // 构建并检测 mycmd
go install -race mypkg // 安装并检测 mypkg
```

### 选项

GORACE 环境变量可以设置竞争检测的选项：

```bash
GORACE="option1=val1 option2=val2"
```

选项：

- `log_path`（默认为 `stderr`）：竞争检测器会将其报告写入名为 `log_path.pid` 的文件中。
  特殊的名字 stdout 和 stderr 会将报告分别写入到标准输出和标准错误中。
- `exitcode`（默认为 66）：当检测到竞争后使用的退出状态。
- `strip_path_prefix`（默认为 ""）：从所有报告文件的路径中去除此前缀， 让报告更加简洁。
- `history_size`（默认为 1）：每个 Go 程的内存访问历史为 `32K * 2**history_size` 个元素。
  增加该值可避免在报告中避免 "failed to restore the stack"（栈恢复失败）的提示，但代价是会增加内存的使用。
- `halt_on_error`（默认为 0）：控制程序在报告第一次数据竞争后是否退出。 例如：
  ```bash
  GORACE="log_path=/tmp/race/report strip_path_prefix=/my/go/sources/" go test -race
  ```

## 运行时开销

竞争检测的代价因程序而异，但对于典型的程序，内存的使用会增加 5 到 10 倍， 而执行时间会增加 2 到 20 倍。