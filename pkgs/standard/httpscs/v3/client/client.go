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

	caCrt, err := ioutil.ReadFile("../../ssl/v3/ca.crt")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("../../ssl/v3/client.crt", "../../ssl/v3/client.key")
	if err != nil {
		fmt.Println("tls.LoadX509KeyPair err:", err)
		return
	}

	ts := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
			Certificates: []tls.Certificate{cliCrt},
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
