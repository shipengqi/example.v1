# https server and client
## 最简单的 http server

### 证书生成

```bash
# 生成服务器端的私钥
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus
..........+++
...................................................................................+++
e is 65537 (0x10001)

# 生成服务器端证书
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:https-example
Email Address []:

```

### server and client

代码在 `standardpkg/httpscs/v1` 目录下。

## 验证服务端证书

多数时候，需要对服务端的证书进行校验，而不是设置 `tls.Config` 的 `InsecureSkipVerify` 为 `true` 来跳过验证。

```bash
# 生成 CA Key
$ openssl genrsa -out ca.key 2048

# 生成 CA 证书
$ openssl req -new -x509 -days 7200 -key ca.key -out ca.crt
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:https-example  
Email Address []:

# 生成 CSR
# CSR 是 Cerificate Signing Request 的英文缩写，为证书请求文件。
# 主要作用是 CA 会利用 CSR 文件进行签名使得攻击者无法伪装或篡改原有证书
$ openssl req -new -key server.key -out server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:https-example
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:

# 基于 CA 签发 server 证书
$ openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.crt
```

### server and client

代码在 `standardpkg/httpscs/v2` 目录下。

运行 `go run client.go`，会报错：

```bash
error: Get "https://localhost:8081": x509: certificate is valid for https-example, not localhost
```

可以修改 `Transport`：

```go
	ts := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
			ServerName: "https-example",
		},
	}
```


## 双向验证证书

服务端可以要求对客户端的证书进行校验，以更严格识别客户端的身份，限制客户端的访问。

```
# 生成客户端的私钥
$ openssl genrsa -out client.key 2048

# 生成 client CSR
$ openssl req -new -key client.key -out client.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:https-example
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:

# 基于 CA 签发 client 证书
$ openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.crt
```

### server and client

代码在 `standardpkg/httpscs/v3` 目录下。
