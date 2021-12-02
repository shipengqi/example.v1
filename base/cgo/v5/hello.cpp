#include <iostream>

extern "C" {
    #include "hello.h"
}
// 使用 c++ 实现 SayHello
int SayHello() {
    std::cout<<"Hello World";
    return 0;
}