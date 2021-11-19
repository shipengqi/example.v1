package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type State struct {
	File string
	Info *stateInfo
}

type stateInfo struct {
	Url      string
	Segments []string
}

func main()  {

	bytesstr := []byte(`{"file": "test", "info": {"url": "testurl"}}`)
	s2 := &State{}
	err := json.Unmarshal(bytesstr, s2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", s2)
	fmt.Printf("%+v\n", s2.Info)
}

