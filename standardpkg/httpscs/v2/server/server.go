package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w,
		"Hi, This is an example of https service in golang!")
}

func main()  {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServeTLS(":8081", "../../ssl/v2/server.crt",
		"../../ssl/v2/server.key", nil)
}
