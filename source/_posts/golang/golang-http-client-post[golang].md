---
title: golang-http-client-post
tags:
  - go
  - golang
  - http
category: golang
abbrlink: 19514
date: 2016-03-13 11:11:55
---
---
title: Golang环境搭建
date: 2016-03-13 11:10:48
tags:
  - go
  - golang
category: golang
---
# 下载Go
[go官方网站](http://golang.org)在大陆已经被和谐，要访问，如果没有梯子，这里有个[传送门](http://www.tvdaili.com/),可以在线代理访问。首先就是要下载go开发程序了，建议[在此处下载](http://golangtc.com/download)对应版本。

# 配置Go环境
    ## Window 环境
下载对应的windows版本（注意64位和32位系统），然后解压得到go目录，假如名字就叫go，绝对路径是E:\go。右键计算机(xp 是我的电脑，windows8是这台电脑)，选择属性，选择左侧的高级系统设置，接下来选择环境变量，出现环境变量的对话框。上面是当前用户的环境变量，也就是说配置的变量只是针对当前用户生效；下面是系统变量，对于整个系统的所有用户生效。我习惯于配置成系统变量，在此也是用系统变量举例。
点击系统变量下的新建按钮，变量名填写 GOROOT,变量值填写上面解压后的go路径，在这也就是 E:\go，然后点击确定。这个变量是用来配置go的home目录。再点击新建按钮，变量名填写GOPATH，变量值填写E:\gopath(gopath要存在，当然也可以选择其他文件夹)，点击确定。这个变量是用来指定go查找包的路径，也是用go get 命令所需安装的位置目录。然后找到PATH变量，点击编辑，变量值最后添加一个半角分号，然后再输入 %GOROOT%\bin，点击确定。可以关闭所有窗口了。
在开始菜单中打开cmd命令行（也可以用快捷键，windows键+R，然后输入cmd回车）。在窗口中输入go回车，如果有go相关的帮助打印，则说明配置成功，否则没有成功，重新校验上面的环境变量配置是否正确。
    ## Linux 环境
    假设对应版本的go解压后绝对路径是/home/user/go。
    vi（或者你喜欢的其他编辑器）打开 ~/.bashrc文件，在后面加入
    >
    GOROOT=/home/user/go
    GOPATH=/home/user/gopath
    PATH=$PATH:$GOROOT/bin
    export GOROOT GOPATH PATH

    然后 运行命令 `soure ~/.bashrc`
    输入 go ，如果打印go帮助文档，说明配置成功，否则检查环境变量配置是否正确。
# 编写第一个Go程序
用你喜欢的编辑器编辑一个文件，加入文件名叫 hello.go 。
敲入以下代码
``` go
package main
import (
"fmt"
)
func main(){
    fmt.println("Hello World!")
}
```
然后再该文件当前目录下，输入 go build hello.go，在目录下会生成一个 hello(windows下是hello.exe)文件，运行该文件，可以看到控制台输入出 HelloWorld,运行成功。
