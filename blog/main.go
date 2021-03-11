package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/router"
	"github.com/shipengqi/example.v1/blog/service"
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/shipengqi/example.v1/tree/main/blog
// @license.name MIT
// @license.url https://github.com/shipengqi/example.v1/blob/main/LICENSE
func main() {

	settings, err := setting.Init("C:\\code\\example.v1\\blog\\conf\\app.debug.ini")
	if err != nil {
		panic(err)
	}
	filename, err := log.Init(settings.Log)
	if err != nil {
		panic(err)
	}

	svc := service.New(settings)
	r := router.Init(svc)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", settings.Server.HttpPort),
		Handler:        r,
		ReadTimeout:    setting.ServerSettings().ReadTimeout,
		WriteTimeout:   setting.ServerSettings().WriteTimeout,
		MaxHeaderBytes: 1 << 20, // max request header bytes
	}

	log.Info().Msgf("server addr %s", s.Addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal().Msgf("Listen err: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	var ctx context.Context
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	for {
		sig := <-quit
		log.Info().Msgf("get a signal %s", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel = context.WithTimeout(context.Background(), 35*time.Second)
			if err := s.Shutdown(ctx); err != nil {
				log.Fatal().Msgf("Shutdown: %s", err)
			}
			log.Info().Msg("Blog server exit")
			log.Info().Msgf("Additional logging details can be found in:\n    %s", filename)
			svc.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
