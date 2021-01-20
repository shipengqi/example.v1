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

Protocol Buffers Libraries的默认安装路径在 `/usr/local/lib`。而我们安装了一个新的动态链接库，`ldconfig` 一般在系统启动时运行，所以现在
会找不到这个 lib，因此我们要手动执行 `ldconfig`，让动态链接库为系统所共享，它是一个动态链接库管理命令。

### Protoc Plugin

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

### 
```bash
protoc --go_out=plugins=grpc:. *.proto
```
## gRPC-gateway

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
