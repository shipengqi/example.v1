package main

import (
	"fmt"
	"reflect"
)

type StructField2 struct {
	Age int
}

type StructField3 struct {
	Address string
}

type StructField4 struct {
	Sex string
}

type Coordinate struct {
	IntFiled1    int32
	IntFiled2    int32
	StringFiled1 string
	StringFiled2 string
	StructField1 struct {
		Name string
	}
	StructField2 StructField2
	StructField3 *StructField3
	StructField4 *StructField4
}

func main() {
	coordinate := Coordinate{
		IntFiled1:    0,
		IntFiled2:    1,
		StringFiled1: "field1",
		StringFiled2: "field2",
		StructField4: &StructField4{
			Sex: "man",
		},
	}
	v := reflect.ValueOf(coordinate)
	t := reflect.TypeOf(coordinate)
	for num := 0; num < v.NumField(); num++ {
		fmt.Printf("type: %s, key: %s, value: %v\n", t.Field(num).Type, t.Field(num).Name, v.Field(num))
		f, _ := t.FieldByName(t.Field(num).Name)
		fmt.Printf("type: %s\n", f.Type)
	}
}
