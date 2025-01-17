---
title: go中main
tags:
  - go
  - golang
category: golang
abbrlink: 14904
date: 2016-03-13 11:16:46
---

想必很多朋友在入门的时候都是拿main开始，而不是test，我也喜欢这样，我想可能是main比较为人熟知的用法吧，test在go中也是非常友好的，不需要依赖其他库就可方便使用。既然都偏向于main方法的开始和入门，那么这个博文就说一下go语言main相关的事情。  
原本只打算写一下main包的拆分和运行方式，突然想到还有其他一些注意地方，那么就一并记录一下，其他的如果使用过程中遇到了，再进行记录。
**以下示例都是在GOPATH下进行**
# main函数定义
想要作为程序的运行入口，那么这个函数必须明明为main，同时要放到main包。main函数声明极其简单，如下：
```go
func main(){
    //...
}
```
这样就声明并定义好了程序的运行的入口函数，不需要其他额外的参数和返回值。加上包的声明，完整的main文件就如下格式：
```go
package main
func main() {
}
```
另外main文件名字可以随意命名，只需要后缀是`.go`就可以了，如：`a.go`, `main.go`, `server.go`等都是可以的。

# main 的运行
要运行main文件，go提供了2中方式：
1. `go run main.go` 其中main.go 就是要运行的main函数所在的文件
2. `go build `命令，文件名可以省略，也可以加上， 还可以用用参数 `-o` 指定编译后的文件，如 `go build -o main.exe main.go` 就是把main.go文件编译成main.exe可执行文件，然后直接执行main.exe就可以运行了。

# main 包文件（内容）的拆分
假如觉得一个main文件中放太多东西有点杂乱，那么可以把main函数和其他内容拆开，放到不同的文件中（这里指的都是main这个包中），文件名字随意。比如我们有个sum函数，那么可以把它放到math.go这个main包的文件中，然后main函数独立在main.go文件中，main方法可以直接调用main包（相同名字的包）的函数、方法或者变量。那么如何执行呢，当然就可以用上述的方法运行。但是如果用`go run main.go`这种方式，那么可能遇到一个问题：找不到sum这个函数。因为run的只是main.go 这一个文件，没有加载其他文件的内容，自然就找不到sum这个函数了，那么加载sum这个函数内容，自然就可以了，对应执行方式就变成了`go run main.go math.go`，也就是说，要把相关以来的内容相关的文件也要加到run 后面。用`go build` 如果不指定文件名字，那么go会自行加载依赖项目，可以顺利执行。


