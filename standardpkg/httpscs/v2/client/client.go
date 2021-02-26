package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main()  {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile("../../ssl/v2/ca.crt")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	// ts := &http.Transport{
	// 	TLSClientConfig: &tls.Config{RootCAs: pool},
	// }

	ts := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
			ServerName: "https-example",
		},
	}

	client := &http.Client{Transport: ts}

	resp, err := client.Get("https://localhost:8081")

	// resp, err := http.Get("https://localhost:8081")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}


// go run client.go 运行 client
// error: Get "https://localhost:8081": x509: certificate signed by unknown authority
// 这是因为 Client 端默认也是要对服务端传过来的数字证书进行校验

// 设置 tls.Config 的 InsecureSkipVerify 为 true
// go run client.go 运行 client
// Hi, This is an example of https service in golang!
