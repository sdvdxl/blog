---
title: 使用autobahn-testsuite测试websocket兼容性
category: java
tags:
  - java
  - websocket
  - wstest
  - autobahn
  - test
abbrlink: 805e
date: 2019-06-16 17:53:22
updateDate: 2019-06-16 17:53:22
---

## 背景

这段时间在线上环境的日志中出现了如下错误日志：
![netty websocket extension negotiated](https://public-links.todu.top/1560676613.png?imageMogr2/thumbnail/!100p)。

关于这个异常的分析和解决参考[netty websocket extension negotiated](/posts/c920)

为了复现线上的问题，本地启动了对应的服务，然后使用 [golang](https://golang.org/) 编写了一个简单的客户端，示例如下：

```go
package main

import (
    "fmt"
    "github.com/gorilla/websocket"
    "github.com/sdvdxl/go-tools/errors"
    "io/ioutil"
    "net/url"
)

func main() {
    u := url.URL{Scheme: "ws", Host: "connectHost", Path: ""}
    conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
    errors.Panic(err)
    c, err := ioutil.ReadAll(resp.Body)
    errors.Panic(err)
    defer resp.Body.Close()
    fmt.Println(string(c))

    errors.Panic(conn.WriteMessage(websocket.TextMessage, []byte("some data\n")))
    t, c, err := conn.ReadMessage()
    errors.Panic(err)
    fmt.Println(t, string(c))
}
```

测试发现没有异常出现。然后又在 ws 建立链接之前增加了以下配置，来测试是否对服务端有影响，结果还是没有问题。

```go
// 如果设置为true，则开启内容压缩传输，这里设置为true 测试是否能影响server
websocket.DefaultDialer.EnableCompression = true
// 自定义sub protocols，测试是否能影响server
websocket.DefaultDialer.Subprotocols = []string{"ss"}
```

经过 ws 客户端配置的变更测试，均没有发现异常，所以肯定是 client 还有什么配置可能影响服务端的处理，但是对 ws 协议内容本身又没有非常了解的情况下，难以修改 client 的链接配置来做各种测试。

于是接下来对关键词 `RSV != 0 and no extension negotiated` 进行了搜索，但是结果都不尽人意，还是没有弄清楚在什么情况下会触发此异常。无意中（过程比较冗长）发现了一个叫做 [autobahn-testsuite](https://github.com/crossbario/autobahn-testsuite) 的组件，可以用来测试对ws的支持情况。

## autobahn-testsuite

### 简介

该测试套件可以对WebSocket基本功能、扩展协议、性能等进行测试。详细说明参见 https://github.com/crossbario/autobahn-testsuite#test-suite-coverage 。


### 说明

简介中已经简要描述了 testsuite 的功能，但是在本文中只是介绍 testsuite 作为 client 对 server 的测试使用，并不涉及 testsuite 的其他功能。

### 安装、运行

testsuite 提供了多种使用方式，但是最简单的就是使用配置好的 Docker 镜像的方式，除此之外，还可以使用原生的 python 执行方式，但是鉴于 python 需要环境配置，这里就以 docker 运行方式说明一下使用说明，其他运行方式可以参考本示例。

github 上有个说明，就是直接执行以下命令运行：

```bash
docker run -it --rm \
    -v "${PWD}/config:/config" \
    -v "${PWD}/reports:/reports" \
    -p 9001:9001 \
    --name fuzzingserver \
    crossbario/autobahn-testsuite
```

说明：

`-v ${PWD}/config:/config` 将当前目录下的 `config` 目录映射到容器的 `/config` 目录
`-v "${PWD}/reports:/reports"` 将当前目录下的 `reports` 目录映射到容器的 `/reports` 目录

需要注意的是，这个命令本身没有问题，但是是运行的 server 模式，就是启动了 server，而不是作为 client 测试 server ，所以需要修改为如下配置：

```bash
docker run -it --rm \
    -v "${PWD}/config:/config" \
    -v "${PWD}/reports:/reports" \
    --name fuzzingclient \
    crossbario/autobahn-testsuite wstest -m fuzzingclient -s config/fuzzingclient.json
```

以上命令中增加了 `wstest -m fuzzingclient -s config/fuzzingclient.json` 命令：

- wstest 就是测试组件的命令
- -m 运行模式，因为我们要运行client模式，所以指定为 fuzzingclient
- -s 指定配置文件

其中配置文件内容格式如下，

```json
{
    "options": {
        "failByDrop": false
    },
    "outdir": "/reports/clients",
    "servers": [
        {
            "name": "AutobahnPython测试Server",
            "url": "ws://192.168.1.6:84",
            "agent": "AutobahnServer",
            "options": {
                "version": 18
            }
        }
    ],
    "cases": [
        "*"
    ],
    "exclude-cases": [],
    "exclude-agent-cases": {}
}
```

配置说明：

主要修改的就是 servers.url 部分，这个地方修改为要测试的 ws server地址，**注意** 如果要测试是本机，则需要填写本机ip地址，不能是 `localhost` 或者 `127.0.0.1`，否则在容器内无法进行测试到本机上的server。

将其放到写入 当前目录下的 config/fuzzingclient.json，然后执行上面的命令，就可以看到输出（部分省略）：

```text
Using explicit spec file 'config/fuzzingclient.json'
Loading spec from /config/fuzzingclient.json

Using Twisted reactor class <class 'twisted.internet.epollreactor.EPollReactor'>
Using UTF8 Validator class <class 'autobahn.websocket.utf8validator.Utf8Validator'>
Using XOR Masker classes <class 'autobahn.websocket.xormasker.XorMaskerNull'>

Autobahn Fuzzing WebSocket Client (Autobahn Testsuite Version 0.8.0 / Autobahn Version 0.10.9)
Ok, will run 519 test cases against 1 servers
Cases = ['1.1.1', ...,'13.7.18']
Servers = [u'ws://192.168.1.6:84']
Running test case ID 1.1.1 for agent AutobahnServer from peer tcp4:192.168.1.6:84
Running test case ID 1.1.2 for agent AutobahnServer from peer tcp4:192.168.1.6:84
....
```

由于测试case比较多（默认是所有），所以耗时及较长，这个时候你可以喝杯茶或者干点其他事情后回来再看看。测试运行结束后程序会自动退出，没有其他输出信息。打开 `reports/clients` 目录，会发现里面生成了一些文件，将 `index.html` 用浏览器打开就可以观察测试报告了，大概长得下面这个模样：

![测试报告](https://public-links.todu.top/1560683789.png?imageMogr2/thumbnail/!100p)

说明：

1. 结果标签说明（建议自己看英文说明理解）
    - Pass 测试执行并且通过
    - Non-Strict 测试执行，但是并非完全符合规范
    - Fail 测试执行但是失败
    - Info （不知如何理解）
    - Missing 测试 case 缺失
1. 可以点击查看测试case详情
1. 可以点击查看测试case和结果详情

## 结语

本文只是对 autobahn-testsuite websocket 测试套件进行了简要说明，更多功能暂时需要自行去了解（后面如果有深入研究我会再进行补充）。

同时需要注意的是，该套件只是对 server 端进行一个测试，失败并不是代表 server 端有问题，比如 server 端 业务逻辑要求如此，或者 server 本身就是要求支持某个版本而已。
