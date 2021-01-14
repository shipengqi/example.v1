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

	flag.Parse()

	if err := preCheck(); err != nil {
		log.Fatalf("precheck err: %s", err)
	}

	log.Println("root path: ", rootPath)

	http.HandleFunc("/files", getSummary)
	http.HandleFunc("/statistics", getSummary)
	addr := ":8080"
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Listen err: %s", err)
	}

	log.Println("Listen addr: ", addr)
}

func preCheck() error {
	if strings.TrimSpace(rootPath) == "" {
		return errors.New("root path is invalid")
	}

	return nil
}
