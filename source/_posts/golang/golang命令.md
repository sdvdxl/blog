---
title: golang命令
tags:
  - go
  - golang
category: golang
abbrlink: 55216
date: 2016-03-13 12:00:37
---
Go 提供了很多好用的命令，比如可以用go get 直接从网络上导入包，下面介绍一些常用的Go命令
# go get
## 基本用法
go get 命令用来直接下载并安装网络上的包到GOPATH中。事实上它依赖于版本控制工具，比如常用的Git和Mercurial。比如我们要使用`github.com/golang/text`这个包，那么可以执行命令 `go get github.com/golang/text`，稍后片刻（依照网络情况时间长短不一），那么这个包源码就会放到 `GOPATH/src`对应的目录下，编译好的文件会放到`GOPATH/pkg`对应的目录下，如果有可执行性代码，那么会将编译好的可执行性文件放到`GOPATH/bin`目录下。
## 高级用法
像刚才例子中提到的`github.com/golang/text`这个包下又有子包，也想get下来，该如何做呢？也许你想没办法，一个一个get吧，的确这是一个不错的方法，但是有一个更为高效和优雅的方式，那就是 `...`，这个符号代表旗下子目录，那我们可以这样操作`go get github.com/golang/text ...`，就可以将text和其子包同时get下来了。
## 参数
- `-v` 可以显示正在get的包
- `-x` 可以显示正在执行的命令
# go build
编译go文件。
可以切换到要build的包中执行`go build`，也可以直接 "go build" + "包名"，如要build包 `github.com/golang/text`，可以切换到text包目录下执行 `go build`，也可以直接执行 `go build github.com/golang/text`，这样就可以build text包了，如果要build子包，那么输入`go build ./...`。`build`命令会将main函数编译成可执行性文件，如果没有main，那么没有额外的文件产生。

# go clean
`clean`和`build`作用相反，是将build出来的可执行性文件清除掉。

# go install
`install` 命令是将包编译成二进制文件并放到`GOPATH/pkg/目标平台/`对应的目录下。比如我们自己从github上克隆了`github.com/golang/text`，要使用它的话，需要进入到text包中，然后执行 `go install`；加入有子包也要安装的话，输入`go install ./...`，就会在`GOPATH/pkg/目标平台/`产生后缀为`.a`的文件。当然也可以和build一样输入完整的包路径。

# go list
这个命令用来查看所引用的包。
## 基本用法
1. 直接跟包名，如：`github.com/golang/text`，会打印这个包本身；也可以切换到该包目录下执行`go list`，效果相同。
2. 可以直接跟"包名"+"/..."，比如要查看`github.com/golang/text`这个包的导入情况，执行`go list github.com/golang/text/...`会打印本包的go文件。
## 参数
`-f`可以指定打印格式，默认是`go list -f '{% raw %} {{.ImportPath}} {% endraw %}'`， `-f`后面的参数可以用下面结构体的属性：

```go
 type Package struct {
    Dir           string // directory containing package sources
    ImportPath    string // import path of package in dir
    ImportComment string // path in import comment on package statement
    Name          string // package name
    Doc           string // package documentation string
    Target        string // install path
    Goroot        bool   // is this package in the Go root?
    Standard      bool   // is this package part of the standard Go library?
    Stale         bool   // would 'go install' do anything for this package?
    Root          string // Go root or Go path dir containing this package

    // Source files
    GoFiles        []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
    CgoFiles       []string // .go sources files that import "C"
    IgnoredGoFiles []string // .go sources ignored due to build constraints
    CFiles         []string // .c source files
    CXXFiles       []string // .cc, .cxx and .cpp source files
    MFiles         []string // .m source files
    HFiles         []string // .h, .hh, .hpp and .hxx source files
    SFiles         []string // .s source files
    SwigFiles      []string // .swig files
    SwigCXXFiles   []string // .swigcxx files
    SysoFiles      []string // .syso object files to add to archive

    // Cgo directives
    CgoCFLAGS    []string // cgo: flags for C compiler
    CgoCPPFLAGS  []string // cgo: flags for C preprocessor
    CgoCXXFLAGS  []string // cgo: flags for C++ compiler
    CgoLDFLAGS   []string // cgo: flags for linker
    CgoPkgConfig []string // cgo: pkg-config names

    // Dependency information
    Imports []string // import paths used by this package
    Deps    []string // all (recursively) imported dependencies

    // Error information
    Incomplete bool            // this package or a dependency has an error
    Error      *PackageError   // error loading package
    DepsErrors []*PackageError // errors loading dependencies

    TestGoFiles  []string // _test.go files in package
    TestImports  []string // imports from TestGoFiles
    XTestGoFiles []string // _test.go files outside package
    XTestImports []string // imports from XTestGoFiles
}
```
如果想打印比较全的信息，也有一个参数可以使用：`-json`，这会将上面部分信息以json格式打印出来，读者可以自行实验。

# go run
这个命令可以直接运行go文件（带main函数），不需要编译成二进制可执行性文件。


