# IO

`Reader` 接口的定义：
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

`Read` 将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数 `n`（`0 <= n <= len(p)`） 以及任何遇到的错误。

即使 `Read` 返回的 `n < len(p)`，它也会在调用过程中占用 `len(p)` 个字节作为暂存空间。若可读取的数据不到 `len(p)` 个字节，`Read` 会
返回可用数据，而不是等待更多数据。

注意：`Read` 在成功读取 `n > 0` 个字节，如果遇到一个错误或 `EOF` (`end-of-file`)，就会返回读取的字节数。所以在处理 IO 错误前，应该先处理读取
到的数据，如：

```go
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}
```

`Writer` 接口的定义：
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

`Write` 将 `len(p)` 个字节从 `p` 中写入到基本数据流中。返回从 `p` 中被写入的字节数 `n`（`0 <= n <= len(p)`）以及任何遇到的引起写入提前
停止的错误。若 `Write` 返回的 `n < len(p)`，它就必须返回一个 `非 nil` 的错误。

`fmt` 标准库中，有一组函数：`Fprint/Fprintf/Fprintln`，第一个参数 `io.Wrtier` 类型，也就是说将数据格式化输出到 `io.Writer` 中。

先看 `fmt.Println` 函数的源码。
```go
func Println(a ...interface{}) (n int, err error) {
    return Fprintln(os.Stdout, a...)
}
```
`fmt.Println` 会将内容输出到标准输出中。

运行示例：
```bash
go run main.go reader.go writer.go
```