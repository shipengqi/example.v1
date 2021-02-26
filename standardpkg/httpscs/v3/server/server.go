package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	_, _ = fmt.Fprintf(w,
		"Hi, This is an example of http service in golang!\n")
}

func main()  {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile("../../ssl/v2/ca.crt")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    ":8081",
		Handler: &myhandler{},
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
