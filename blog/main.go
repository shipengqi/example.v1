package main

import (
	"fmt"
	"net/http"

	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/router"
)

func main()  {
	r := router.Init()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        r,
		ReadTimeout:    setting.ServerSettings().ReadTimeout,
		WriteTimeout:   setting.ServerSettings().WriteTimeout,
		MaxHeaderBytes: 1 << 20, // max request header bytes
	}
	fmt.Println("server addr: ", s.Addr)

	_ = s.ListenAndServe()
}
