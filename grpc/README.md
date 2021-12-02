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

### 根证书

根证书（root certificate）是属于根证书颁发机构（CA）的公钥证书。我们可以通过验证 CA 的签名从而信任 CA ，任何人都可以得到 CA 的证书（含公钥），
用以验证它所签发的证书（客户端、服务端）

它包含的文件如下：

- 公钥
- 密钥

### 生成 Key

```
# 生成 CA Key
openssl genrsa -out ca.key 2048

# 生成 CA 证书
openssl req -new -x509 -days 7200 -key ca.key -out ca.crt
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:
State or Province Name (full name) []:
Locality Name (eg, city) []:
Organization Name (eg, company) []:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:grpc-example
Email Address []:

# 生成 CSR
# CSR 是 Cerificate Signing Request 的英文缩写，为证书请求文件。
# 主要作用是 CA 会利用 CSR 文件进行签名使得攻击者无法伪装或篡改原有证书
openssl req -new -key server.key -out server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) []:
State or Province Name (full name) []:
Locality Name (eg, city) []:
Organization Name (eg, company) []:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:grpc-example
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:

# 基于 CA 签发 server 证书
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.crt

# 生成 client Key
openssl ecparam -genkey -name secp384r1 -out client.key

# 生成 client CSR
openssl req -new -key client.key -out client.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:CN
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:Shanghai
Organization Name (eg, company) [Default Company Ltd]:MF
Organizational Unit Name (eg, section) []:RA
Common Name (eg, your name or your server's hostname) []:grpc-example
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:

# 基于 CA 签发 client 证书
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.crt
```

#### 整理目录

将生成的一堆文件，请按照以下目录结构存放：

```
$ tree ssl 
ssl
├── ca
│   ├── ca.crt
│   ├── ca.key
│   └── ca.srl
├── client
│   ├── client.crt
│   ├── client.csr
│   └── client.key
└── server
    ├── server.crt
    ├── server.csr
    └── server.key
```


### server

```go
    // 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair("../../ssl/server/server.crt", "../../ssl/server/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}
    // 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../ssl/ca/ca.crt")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
    // 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},        // 设置证书链，允许包含一个或多个
		ClientAuth:   tls.RequireAndVerifyClientCert, // 必须校验客户端的证书
		ClientCAs:    certPool,                       // 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})
	// 创建 grpc Server 对象
	server := grpc.NewServer(grpc.Creds(c))
```

`ClientAuth` 可以填的参数:

```go
const (
	NoClientCert ClientAuthType = iota
	RequestClientCert
	RequireAnyClientCert
	VerifyClientCertIfGiven
	RequireAndVerifyClientCert
)
```

### client

```go
	cert, err := tls.LoadX509KeyPair("../../ssl/client/client.crt", "../../ssl/client/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../ssl/ca/ca.crt")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "grpc-example",
		RootCAs:      certPool,
	})
	
	// 创建与 server 的连接
	conn, err := grpc.Dial(fmt.Sprintf(":%s", PORT), grpc.WithTransportCredentials(c))
```

简单流程：

1. Client 通过请求得到 Server 端的证书
2. 使用 CA 认证的根证书对 Server 端的证书进行可靠性、有效性等校验
3. 校验 ServerName 是否可用、有效

在设置了 `tls.RequireAndVerifyClientCert` 模式的情况下，Server 也会使用 CA 认证的根证书对 Client 端的证书进行可靠性、有效性等校验。
也就是两边都会进行校验。

## 拦截器
gRPC 中，可分为两种 RPC 方法，与拦截器的对应关系是：

- 普通方法：一元拦截器（`grpc.UnaryInterceptor`）
- 流方法：流拦截器（`grpc.StreamInterceptor`）

### grpc.UnaryInterceptor

```go
func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
	return func(o *options) {
		if o.unaryInt != nil {
			panic("The unary server interceptor was already set and may not be reset.")
		}
		o.unaryInt = i
	}
}
```
函数原型：
```go
type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
```

通过查看源码可得知，要完成一个拦截器需要实现 `UnaryServerInterceptor` 方法。形参如下：

- `ctx context.Context`：请求上下文
- `req interface{}`：RPC 方法的请求参数
- `info *UnaryServerInfo`：RPC 方法的所有信息
- `handler UnaryHandler`：RPC 方法本身

### grpc.StreamInterceptor

```go
func StreamInterceptor(i StreamServerInterceptor) ServerOption
```

函数原型：
```go
type StreamServerInterceptor func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
```

StreamServerInterceptor 与 UnaryServerInterceptor 形参的意义是一样的。

### 如何实现多个拦截器

另外，可以发现 gRPC 本身居然只能设置一个拦截器，要实现多个拦截器，可以采用开源
项目 [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware) 。

```go
import "github.com/grpc-ecosystem/go-grpc-middleware"

myServer := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        ...
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
       ...
    )),
)
```

## http

gRPC 的协议是基于 HTTP/2 的，因此应用程序能够在单个 TCP 端口上提供 HTTP/1.1 和 gRPC 接口服务（两种不同的流量）

```go
	http.ListenAndServeTLS(fmt.Sprintf(":%s", PORT), certFile, keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检测请求协议是否为 HTTP/2
			// 判断 Content-Type 是否为 application/grpc（gRPC 的默认标识位）
			// 根据协议的不同转发到不同的服务处理
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}

			return
		}),
	)
```

## 认证

```
type PerRPCCredentials interface {
    GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    RequireTransportSecurity() bool
}
```

在 gRPC 中默认定义了 `PerRPCCredentials`，是 gRPC 默认提供用于自定义认证的接口，它的作用是将所需的安全认证信息添加到每个 RPC 方法的上下
文中。其包含 2 个方法：

- `GetRequestMetadata`：获取当前请求认证所需的元数据（metadata）
- `RequireTransportSecurity`：是否需要基于 TLS 认证进行安全传输


## Deadline

gRPC Deadlines 的用法。


当未设置 Deadlines 时，将采用默认的 DEADLINE_EXCEEDED（这个时间非常大）

如果产生了阻塞等待，就会造成大量正在进行的请求都会被保留，并且所有请求都有可能达到最大超时

这会使服务面临资源耗尽的风险，例如内存，这会增加服务的延迟，或者在最坏的情况下可能导致整个进程崩溃


## 链路追踪

OpenTracing 通过提供平台无关、厂商无关的API，使得开发人员能够方便的添加（或更换）追踪系统的实现

不过 OpenTracing 并不是标准。因为 CNCF 不是官方标准机构，但是它的目标是致力为分布式追踪创建更标准的 API 和工具

### Topic

#### Trace

一个 trace 代表了一个事务或者流程在（分布式）系统中的执行过程

#### Span

一个 span 代表在分布式系统中完成的单个工作单元。也包含其他 span 的 “引用”，这允许将多个 spans 组合成一个完整的 Trace

每个 span 根据 OpenTracing 规范封装以下内容：

- 操作名称
- 开始时间和结束时间
- key:value span Tags
- key:value span Logs
- SpanContext

#### Tags

Span tags（跨度标签）可以理解为用户自定义的 Span 注释。便于查询、过滤和理解跟踪数据

#### Logs

Span logs（跨度日志）可以记录 Span 内特定时间或事件的日志信息。主要用于捕获特定 Span 的日志信息以及应用程序本身的其他调试或信息输出

#### SpanContext

SpanContext 代表跨越进程边界，传递到子级 Span 的状态。常在追踪示意图中创建上下文时使用

#### Baggage Items

Baggage Items 可以理解为 trace 全局运行中额外传输的数据集合

### Zipkin

[Zipkin](https://github.com/openzipkin/zipkin) 是分布式追踪系统。它的作用是收集解决微服务架构中的延迟问题所需的时序数据。它管理这些数据的收集和查找

Zipkin 的设计基于 [Google Dapper](https://research.google/pubs/pub36356/) 论文。

跟踪系统中通常有四个组件，Zipkin 包括：

- Recorder(记录器)：记录跟踪数据
- Reporter (or collecting agent)(报告器或收集代理)：从记录器收集数据并将数据发送到 UI 程序
- Tracer：生成跟踪数据
- UI：负责在图形 UI 中显示跟踪数据

#### 运行

```bash
docker run -d -p 9411:9411 openzipkin/zipkin
```

### gRPC + Opentracing + Zipkin
实现 gRPC 通过 Opentracing 标准 API 对接 Zipkin，再通过 Zipkin 去查看数据


## gRPC-gateway

如果希望用 Rpc 作为内部 API 的通讯，同时也想对外提供 Restful Api，写两套又太繁琐不符合，就可以使用 [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)。

官方文档：<https://grpc-ecosystem.github.io/grpc-gateway/>

```bash
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```
