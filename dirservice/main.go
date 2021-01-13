package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"strings"
)

var rootPath string

func init()  {
	flag.StringVar(&rootPath, "root", "", "Specifies the root path.")
}

func main()  {

	if err := preCheck(); err != nil {
		log.Fatalf("precheck err: %s", err)
	}

	r := InitRouter()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Listen err: %s", err)
	}
}

func preCheck() error {
	if strings.TrimSpace(rootPath) == "" {
		return errors.New("root path is invalid")
	}

	return nil
}

