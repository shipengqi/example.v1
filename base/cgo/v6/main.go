package main

// __attribute__ constructor/destructor 若函数被设定为 constructor 属性,则该函数会在 main 函数执行之前被自动的执行
// 拥有此类属性的函数经常隐式的用在程序的初始化数据方面


/*
#include<stdio.h>
__attribute__((constructor)) void before_main() {
   printf("before main\n");
}
*/
import "C"

import "log"
func main()  {
	log.Printf("hello world!")
}

// Output:
// before main
// 2021/08/16 17:22:16 hello world!