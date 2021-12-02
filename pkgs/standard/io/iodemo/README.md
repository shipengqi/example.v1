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

实现 Reader 和 Writer 接口的类型：
- `os.File` 同时实现了 `io.Reader` 和 `io.Writer`
- `strings.Reader` 实现了 `io.Reader`
- `bufio.Reader/Writer` 分别实现了 `io.Reader` 和 `io.Writer`
- `bytes.Buffer` 同时实现了 `io.Reader` 和 `io.Writer`
- `bytes.Reader` 实现了 `io.Reader`
- `compress/gzip.Reader/Writer` 分别实现了 `io.Reader` 和 `io.Writer`
- `crypto/cipher.StreamReader/StreamWriter` 分别实现了 `io.Reader` 和 `io.Writer`
- `crypto/tls.Conn` 同时实现了 `io.Reader` 和 `io.Writer`
- `encoding/csv.Reader/Writer `分别实现了 `io.Reader` 和 `io.Writer`
- `mime/multipart.Part` 实现了 `io.Reader`
- `net/conn` 分别实现了 `io.Reader` 和 `io.Writer`

ReaderAt 接口的定义：
```go
type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
```

`ReadAt` 从基本输入源的偏移量 `off` 处开始，将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数 `n`（`0 <= n <= len(p)`）以及任何遇到的错误。

当 `ReadAt` 返回的 `n < len(p)` 时，它就会返回一个 `非 nil` 的错误来解释为什么没有返回更多的字节。`ReadAt` 比 `Read` 更严格。

即使 `ReadAt` 返回的 `n < len(p)`，它也会在调用过程中使用 `p` 的全部作为暂存空间。**若可读取的数据不到 `len(p)` 字节，`ReadAt` 就会阻塞，直
到所有数据都可用或一个错误发生**。

若 `n = len(p)` 个字节从输入源的结尾处由 `ReadAt` 返回，`Read` 可能返回 `err == EOF` 或者 `err == nil`。

```go
reader := strings.NewReader("example.v1")
p := make([]byte, 6)
n, err := reader.ReadAt(p, 2)
if err != nil {
    panic(err)
}
fmt.Printf("%s, %d\n", p, n) // ample., 6
```

WriterAt 接口的定义：
```go
type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
}
```

`WriteAt` 从 `p` 中将 `len(p)` 个字节写入到偏移量 `off` 处的基本数据流中。它返回从 `p` 中被写入的字节数 `n`（`0 <= n <= len(p)`）以及任何遇到的引起写入提
前停止的错误。若 **`WriteAt` 返回的 `n < len(p)`，它就必须返回一个 `非 nil` 的错误**。

若被写区域没有重叠，可对相同的目标并行执行 `WriteAt` 调用。

```go
file, err := os.Create("writeAt.txt")
if err != nil {
    panic(err)
}
defer file.Close()
file.WriteString("hello, overwrite")
n, err := file.WriteAt([]byte("example.v1"), 7)
if err != nil {
    panic(err)
}
fmt.Println(n) // 10
```
`writeAt.txt` 中的内容： `hello, example.v1`。

运行示例：
```bash
go run main.go reader.go writer.go
```