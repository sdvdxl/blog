---
title: 配置Idea的Go开发环境
tags:
  - go
  - golang
category: golang
abbrlink: 47403
date: 2016-03-13 11:43:19
---
# 获取IDEA

1. [最新版戳这里下载](https://www.jetbrains.com/idea/download/)，下载对应平台的版本，一般来说，社区版(Community Edition)就已经足够了。
2. 安装IDEA。Windows和Mac是安装版，一步一步安装完成即可；Linux是免安装版，解压，给IDEA的执行文件添加可执行权限即可。假设安装目录是`H:\software\dev\JetBrains\IntelliJ IDEA Community Edition 14.1.1`，其他平台自行设置目录。

# 获取Go

1. Go的官方网站是[http://golang.org/](http://golang.org/)，[下载地址](https://golang.org/dl/)，但是鉴于中国网络问题，不科学上网则没法下载，各位同学可以从[这里下载](https://studygolang.com/dl)，最好下载最新的，当然这个网站(https://studygolang.com)本身就是国内较活跃的一个Go社区。
2.  下载后解压到本地目录，假设安装目录是`H:\software\dev\go`。其他平台自行设置目录。

# 获取Git和hg

因为Go get命令要使用到git或者hg，所以需要安装git和hg。
git可以[从这里下载](http://git-scm.com/), hg(执行命令是hg，实际下载的软件叫mercurial)可以[从这里下载](http://mercurial.selenic.com/)，加入git安装到`H:\software\dev\Git`，hg安装到`H:\software\dev\Mercurial`

# 配置Go和Git(hg)

配置环境变量，这里以Windows为例，其他平台请自行换成对应的路径即可。如果打开命令行分别执行以下命令都成功，那么不需要额外配置环境变量，否则配置对应的环境变量

``` bash
go version
git version
hg version
```

变量名称 | 变量值 | 说明
------ | ----- |-----
GOROOT | H:\software\dev\go | go根路径
GOPATH  | H:\software\dev\gopath | gopath可以是任何一个目录
PATH | %PATH%;%GOROOT%\bin;H:\software\dev\Git\bin;H:\software\dev\Mercuria; | 注意不要忘了加入原来的path变量

# 配置IDEA的Go环境

1. 打开IDEA，File -> Settings -> Plugins -> Browse repositiores... -> Manage repositories... ，添加自定义repository url `https://plugins.jetbrains.com/plugins/nightly/list`(nightly build)或者 `https://plugins.jetbrains.com/plugins/alpha/list`(alpha version)，添加完成之后，等待刷新完成后，输入go，选择go插件，点击安装，等待安装完成后，重启生效。网络环境不好的话，可能插件不能下载，可以直接[去idea官网下载插件](https://plugins.jetbrains.com/plugin/download?updateId=19402)，如何获取最新插件呢，这里是根据updateId来的，这个最新的id就是从上面的repository的url中获取的，用浏览器打开这个url，就会观察到以下内容

```  xml
<idea-plugin downloads="97922" size="1071401" date="1428797441000" url="">
<name>Go</name>
<id>ro.redeul.google.go</id>
<description>
<![CDATA[
Support for Go programming language. <p>Alpha pre-release of the 1.0.0 version.</p> <p>Doesn't contain all the functionality of the 0.9.x branch but has a completely reworked internals. It's faster than 0.9.x, refactoring works to some degree and has native support for gopath packages.</p> Compatibility <p>Plugin can be installed on IntelliJ platform 141.2 or greater. It corresponds to IntelliJ IDEA 14.1, WebStorm 10, PhpStorm 9</p>
]]>
</description>
<version>0.9.271</version>
<vendor email="" url="https://github.com/go-lang-plugin-org"/>
<download-url>../../plugin/download?updateId=19402</download-url>
<idea-version min="n/a" max="n/a" until-build="3999"/>
<change-notes>
<![CDATA[
<ul> <li>Initial GAE support: running dev server. <strong>Requires resetting project SDK.</strong></li> </ul>
]]>
</change-notes>
<rating>4.3</rating>
</idea-plugin>
    ```
其中 download-url 中的 updateId 就是最新的下载id。
如果还是没法下载，那么[请点击这里](http://pan.baidu.com/s/1c0o50ys)，从百度云上下载。
然后 File -> Settings -> Install plugin from disk...选择刚才下载的压缩包（不要解压），确定后重启成效。
2.  File -> Other settings -> Default Project Structure... -> Platform Settings -> SDKs -> + -> Go SDK -> 选择GOROOT路径，确定。
3. File -> New Project -> Go -> Next -> 输入Project name和Project location -> Finish -> 在项目根目录中新建main.go，添加以下内容

```go
package main
  import (
"fmt"
)

func main() {
fmt.Println("hello world")
}
```

Run -> Edit Configrations -> + -> Go Application -> File 中在原来的路径基础上添加main文件即 添加 `\main.go`，点击确定，然后运行可以看到控制台打印 hello world。

# 配置GDB debug

Run -> Edit Configrations -> Defaults -> Go GDB ->
Name：可以随便填写
GDB executeable：dbg.exe的完整路径
Application executable：填写生成的可执行文件的完整路径，路径要是windows写法，如G:\gopath\src\example\main.exe则需要两个反斜杠，就变成了G:\\gopath\\src\\example\\main.exe，或者是Unix写法 G:/gopath/src/example/main.exe,否则会提示找不到文件。这里的可执行文件必须使用`go build -gcflags "-N -l"`编译出来的，这样的文件带有debug信息并且没有被go内联优化。

# 配置保存时自动格式化代码和自动导入

这个配置需要用到IDEA的宏（所谓的宏，就是一系列操作），下面就说怎么录制这个宏。
