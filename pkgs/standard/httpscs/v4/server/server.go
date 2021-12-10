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

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer,
			"Hi, This is an example of http service in golang!\n")
	})

	s := &http.Server{
		Addr:    ":8081",
		Handler: mux,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	err = s.ListenAndServeTLS("../../ssl/v3/server.crt", "../../ssl/v3/server.key")
	if err != nil {
		fmt.Println("ListenAndServeTLS err:", err)
	}
}
