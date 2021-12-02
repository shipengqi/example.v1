package main

import (
	"fmt"
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// struct field must Capital letter, otherwise it cannot be mapped to the corresponding value
type State struct {
	File string
	Info *stateInfo
}

type stateInfo struct {
	Url      string
	Segments []string
}

func main()  {
	s := &State{
		File: "test.json",
		Info: &stateInfo{
			Url:      "test url",
			Segments: make([]string, 0),
		},
	}
	s.SetSegments([]string{"segment1", "segment2", "segment3"})
	err := s.Save()
	if err != nil {
		fmt.Println(err)
	}

	bytesstr := []byte(`{"file": "test", "info": {"url": "testurl"}}`)
	s2 := &State{}
	err = json.Unmarshal(bytesstr, s2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", s2)
}

func (s *State) Save() error {
	f, err := os.OpenFile(s.File, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(s.Info)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) Read() error {
	data, err := ioutil.ReadFile(s.File)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s.Info)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) Segments() []string {
	return s.Info.Segments
}

func (s *State) SetSegments(segments []string) {
	s.Info.Segments = segments
}

