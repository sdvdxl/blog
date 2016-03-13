---
title: go交叉编译
date: 2016-03-13 11:41:45
tags:
  - go
  - golang
category: golang
---
1.5之前需要进入go安装目录下的src目录，然后执行
`GOOS=windows GOARCH=i386  CGO_ENABLED=0 ./make.bash --no-clean`

1.5及其之后可以直接执行` CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o 输出文件名 go文件
`
其中`GOOS`有：`windows`,`linux`,`darwin`也就是mac系统
`GOARCH`有：`adm64`和`386`分别对应64位平台和32位平台
