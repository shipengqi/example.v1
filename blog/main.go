package main

import (
	"fmt"
	"net/http"

	"github.com/shipengqi/example.v1/blog/model"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/router"
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/shipengqi/example.v1/tree/main/blog
// @license.name MIT
// @license.url https://github.com/shipengqi/example.v1/blob/main/LICENSE
func main()  {
	settings, err := setting.Init("C:\\code\\example.v1\\blog\\conf\\app.debug.ini")
	if err != nil {
		fmt.Println(err)
	}
	_, err = log.Init(settings.Log)
	if err != nil {
		fmt.Println(err)
	}
	model.Init()
	r := router.Init()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        r,
		ReadTimeout:    setting.ServerSettings().ReadTimeout,
		WriteTimeout:   setting.ServerSettings().WriteTimeout,
		MaxHeaderBytes: 1 << 20, // max request header bytes
	}

	log.Info().Msgf("server addr %s", s.Addr)

	_ = s.ListenAndServe()
}
