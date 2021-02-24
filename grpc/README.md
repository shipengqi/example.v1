# gRPC

RPC 代指远程过程调用（Remote Procedure Call），它的调用包含了传输协议和编码（对象序列号）协议等等。

## 常见 RPC 框架

- [gRPC](https://grpc.io/)
- [Thrift](https://github.com/apache/thrift)
- [Rpcx](https://github.com/smallnest/rpcx)
- [Dubbo](https://github.com/apache/incubator-dubbo)

| \      | 跨语言 | 多 IDL | 服务治理 | 注册中心 | 服务管理 |
| ------ | ------ | ------ | -------- | -------- | -------- |
| gRpc   | √      | ×      | ×        | ×        | ×        |
| Thrift | √      | ×      | ×        | ×        | ×        |
| Rpcx   | ×      | √      | √        | √        | √        |
| Dubbo  | ×      | √      | √        | √        | √        |


## Protobuf

Protocol Buffers 是一种与语言、平台无关，可扩展的序列化结构化数据的方法，常用于通信协议，数据存储等等。相较于 JSON、XML，它更小、更快、更简单。

```protobuf
syntax = "proto3";

// 定义 SearchService，包含 RPC 方法 Search，入参为 `SearchRequest` 消息，出参为 `SearchResponse` 消息
service SearchService {
    rpc Search (SearchRequest) returns (SearchResponse);
}

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}

message SearchResponse {
    ...
}
```

第一行（非空的非注释行）声明使用 `proto3` 语法。如果不声明，将默认使用 `proto2` 语法。

[Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)

### 数据类型

| .proto Type | Java Type  | Go Type |
| ----------- | ---------- | ------- |
| double      | double     | float64 |
| float       | float      | float32 |
| int32       | int        | int32   |
| int64       | long       | int64   |
| uint32      | int        | uint32  |
| uint64      | long       | uint64  |
| sint32      | int        | int32   |
| sint64      | long       | int64   |
| fixed32     | int        | uint32  |
| fixed64     | long       | uint64  |
| sfixed32    | int        | int32   |
| sfixed64    | long       | int64   |
| bool        | boolean    | bool    |
| string      | String     | string  |
| bytes       | ByteString | []byte  |

## gRPC

gRPC 是一个 基于 HTTP/2 协议设计的 RPC 框架，它采用了 Protobuf 作为 IDL

### 特点

1. HTTP/2
2. Protobuf
3. 客户端、服务端基于同一份 IDL
4. 移动网络的良好支持
5. 支持多语言

![](gRPC-flow.png)

1. 客户端（gRPC Sub）调用 A 方法，发起 RPC 调用
2. 对请求信息使用 Protobuf 进行对象序列化压缩（IDL）
3. 服务端（gRPC Server）接收到请求后，解码请求体，进行业务逻辑处理并返回
4. 对响应结果使用 Protobuf 进行对象序列化压缩（IDL）
5. 客户端接受到服务端响应，解码请求体。回调被调用的 A 方法，唤醒正在等待响应（阻塞）的客户端调用并返回响应结果


## 安装

### gRPC

```bash
go get -u google.golang.org/grpc
```

### Protocol Buffers v3

```bash
wget https://github.com/google/protobuf/releases/download/v3.5.1/protobuf-all-3.5.1.zip
unzip protobuf-all-3.5.1.zip
cd protobuf-3.5.1/
./configure
make
make install
```

检查是否安装成功

```bash
protoc --version
```

若出现以下错误，执行 `ldconfig` 命名就能解决这问题

```bash
protoc: error while loading shared libraries: libprotobuf.so.15: cannot open shared object file: No such file or directory
```

Protocol Buffers Libraries 的默认安装路径在 `/usr/local/lib`。而我们安装了一个新的动态链接库，`ldconfig` 一般在系统启动时运行，所以现在
会找不到这个 lib，因此我们要手动执行 `ldconfig`，让动态链接库为系统所共享，它是一个动态链接库管理命令。

### Protoc Plugin

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

### 生成
```bash
protoc --go_out=plugins=grpc:. *.proto
```

- `plugins=plugin1+plugin2`：指定要加载的子插件列表

我们定义的 proto 文件是涉及了 RPC 服务的，而默认是不会生成 RPC 代码的，因此需要给出 `plugins` 参数传递给 `protoc-gen-go`，告诉它，
支持 RPC（这里指定了 gRPC）

- `--go_out=.`：设置 Go 代码输出的目录

该指令会加载 `protoc-gen-go` 插件达到生成 Go 代码的目的，生成的文件以 `.pb.go` 为文件后缀

- `:` （冒号）：

冒号充当分隔符的作用，后跟所需要的参数集。如果这处不涉及 RPC，命令可简化为：

```
$ protoc --go_out=. *.proto
```

### gRPC-gateway

```bash
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```


## 流式 RPC

gRPC 的流式，分为三种类型：

- Server-side streaming RPC：服务器端流式 RPC
- Client-side streaming RPC：客户端流式 RPC
- Bidirectional streaming RPC：双向流式 RPC


### 为什么用 Streaming RPC

在使用 Simple RPC 时，有如下问题：

- 数据包过大造成的瞬时压力
- 接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）

Streaming RPC 使用场景：

- 大规模数据包
- 实时场景


```protobuf
syntax = "proto3";

package proto;

service StreamService {
    rpc List(StreamRequest) returns (stream StreamResponse) {};

    rpc Record(stream StreamRequest) returns (StreamResponse) {};

    rpc Route(stream StreamRequest) returns (stream StreamResponse) {};
}

message StreamPoint {
  string name = 1;
  int32 value = 2;
}

message StreamRequest {
  StreamPoint pt = 1;
}

message StreamResponse {
  StreamPoint pt = 1;
}
```

关键字 `stream`，声明其为一个流方法。

- List：服务器端流式 RPC
- Record：客户端流式 RPC
- Route：双向流式 RPC


### Server-side streaming RPC

服务器端流式 RPC，是单向流，指 Server 为 Stream 相应，而 Client 为普通 RPC 请求：

server：

```go
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
```

`stream.Send` 方法是 protoc 生成的：

```go
type StreamService_ListServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

func (x *streamServiceListServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}
```

最终调度内部的 `SendMsg` 方法，`SendMsg` 的执行过程：

- 消息体（对象）序列化
- 压缩序列化后的消息体
- 对正在传输的消息体增加 5 个字节的 header
- 判断压缩+序列化后的消息体总字节长度是否大于预设的 `maxSendMessageSize`（预设值为 `math.MaxInt32`），若超出则提示错误
- 写入给流的数据集

client：

```go
func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}
```

`stream.Recv()` 方法：

```go
type StreamService_ListClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

func (x *streamServiceListClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
```

`RecvMsg` 会从流中读取完整的 gRPC 消息体：

1. `RecvMsg` 是阻塞等待的
2. `RecvMsg` 当流成功/结束（调用了 `Close`）时，会返回 `io.EOF`
3. `RecvMsg` 当流出现任何错误时，流会被中止，错误信息会包含 RPC 错误码。而在 `RecvMsg` 中可能出现如下错误：
  - `io.EOF`
  - `io.ErrUnexpectedEOF`
  - `transport.ConnectionError`
  - `google.golang.org/grpc/codes`

默认的 `MaxReceiveMessageSize` 值为 `1024 _ 1024 _ 4`，建议不要超出

切换到 stream 目录下，验证：

```bash
$ go run server/server.go

$ go run client/client.go
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2018
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2019
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2020
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2021
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2022
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2023
2021/01/22 13:07:44 resp: pj.name: gRPC Stream Client: List, pt.value: 2024
```

### Client-side streaming RPC

客户端通过流式发起多次 RPC 请求给服务端，服务端发起一次响应给客户端：

server：
```go
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}
```

`stream.SendAndClose` 方法有什么用？

上面的代码中对每一个 `Recv` 都进行了处理，当发现 `io.EOF` (流关闭) 后，需要将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 `Recv`。

client：

```go
func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}
```

`stream.CloseAndRecv` 和 `stream.SendAndClose` 是配套使用的流方法。

切换到 stream 目录下，验证：

```bash
$ go run client/client.go
2021/01/22 13:14:28 resp: pj.name: gRPC Stream Server: Record, pt.value: 1

$ go run server/server.go
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
2021/01/22 13:14:28 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 2018
```

### Bidirectional streaming RPC

由客户端以流式的方式发起请求，服务端同样以流式的方式响应请求。

server：

```go
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}
```

client：

```go
func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	stream.CloseSend()

	return nil
}
```

切换到 stream 目录下，验证：

```bash
$ go run server/server.go
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018
2021/01/22 13:20:42 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 2018

$ go run client/client.go
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 0
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 1
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 2
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 3
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 4
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 5
2021/01/22 13:20:42 resp: pj.name: gPRC Stream Client: Route, pt.value: 6
```


## gRPC TLS

不带证书的 client：

```go
conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
```

`grpc.WithInsecure()` 方法

```go
func WithInsecure() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.insecure = true
	})
}
```

`DialOption` 禁用了安全传输。

### 证书生成

#### 私钥

```
openssl ecparam -genkey -name secp384r1 -out server.key
```

#### 自签公钥

```
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

填写信息：

```
Country Name (2 letter code) []:
State or Province Name (full name) []:
Locality Name (eg, city) []:
Organization Name (eg, company) []:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:grpc-example
Email Address []:
```

注意：**Common Name 和 `credentials.NewClientTLSFromFile("../../conf/server.pem", "grpc-example")` 中
的 server name 一致**，否则校验会失败。

### client

```go
	c, err := credentials.NewClientTLSFromFile("../../conf/server.pem", "grpc-example")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
```

### server

```go
	c, err := credentials.NewServerTLSFromFile("../../conf/server.pem", "../../conf/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}

	server := grpc.NewServer(grpc.Creds(c))
```

运行 `go run client.go` 输出：
```bash
$ go run client.go
2021/02/24 17:52:16 resp: Hello, gRPC Server
```

Client 是基于 Server 端的证书和服务名称来建立请求的。这样的话，就需要将 Server 的证书通过各种手段给到 Client 端，
否则是无法完成这项任务的。

为了保证证书的可靠性和有效性，需要引入 CA 颁发的根证书。
