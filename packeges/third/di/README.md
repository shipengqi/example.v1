# 依赖注入

IoC 是一种设计思想，其核心作用是降低代码的耦合度。依赖注入是一种实现控制反转且用于解决依赖性问题的设计模式。

依赖注入处理的关键问题是解耦，解耦在代码工程学中的好处显而易见：代码扩展性强，可维护性增强以及更容易的进行单元测试。

小项目并不需要依赖注入，当项目庞大到一定程度，结构之间的关系变得非常复杂时，手动创建每个依赖，然后层层组装起来的方式就会变得异常繁琐，并且容易出错。

## 依赖注入框架

Go 的依赖注入框架有两类：

1. 一类是通过反射在运行时进行依赖注入，Uber 开源的 [dig](https://github.com/uber-go/dig) 和 Meta 的 [inject](https://github.com/facebookarchive/inject)。
2. 一类是通过 generate 进行代码生成， 官方的 [wire](https://github.com/google/wire)。

## Dig
