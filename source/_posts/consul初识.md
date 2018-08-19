---
title: consul初识
category: consul
tags:
  - consul
  - 服务发现
abbrlink: '1597'
date: 2018-08-14 09:48:35
updateDate: 2018-08-14 09:48:35
---

## consul 简介

[consul](https://www.consul.io/intro/index.html) 是一个支持多数据中心，加密数据交换，KV存储，健康检查的服务发现组件。[官方对比了其他类似的组件](https://www.consul.io/intro/vs/index.html)，比如知名的[ZooKeeper](http://zookeeper.apache.org/) 、[etcd](https://coreos.com/etcd/)、[Eureka](https://github.com/Netflix/eureka)甚至是刚刚推出不久的[Istio](https://istio.io/)。总的来讲consul有如下优势：

- 支持多数据中心
- 强一致性
- 支持kv存储
- 支持健康检查（可以通过自定义脚本）
- 支持多节点（raft选举）
- 支持http(s)方式注册发现服务
- 支持DNS方式发现服务并提供负载均衡（高级用法）
- 支持外部服务注册
- 测试部署容易（仅一个2进制文件 [go](https://golang.org) 编写，支持跨平台），同时支持server和client模式，还支持开发（dev）模式

当然其缺点也是有的，当前我认为的最大的缺点是负载均衡和故障自动转移上。简单来说，一个agent节点如果宕机，那么其上面的 所有服务都会消失，在整个集群中不可见；并且如果要将集群中的服务注销，必须通过注册的时候的agentAPI进行反注册才可以成功注销，否则虽然提示是200（成功），实际上却没有任何效果。

## 下载和安装

官方提供了各平台的[二进制文件的下载链接](https://www.consul.io/downloads.html)，可以下载对应平台的文件进行安装，这个是最简单的方式。当然也可以通过源码进行编译，不过这个需要Go语言的支持。下面的操作均在CentOS7（64位）上进行。

这里我们下载Linux 64位平台，版本为**1.2.2**的[2进制文件](https://releases.hashicorp.com/consul/1.2.2/consul_1.2.2_linux_amd64.zip)并放到`/opt/consul`目录下，文件名为 `consul_1.2.2_linux_amd64.zip`，切换目录并执行解压命令 `cd /opt/consul && unzip consul_1.2.2_linux_amd64.zip` 会得到一个`consul`文件。

## 运行

输入 `./consul`即可得到如下提示则说明已经准备好组件了。

```shell

Usage: consul [--version] [--help] <command> [<args>]

Available commands are:
    agent          Runs a Consul agent
    catalog        Interact with the catalog
    connect        Interact with Consul Connect
    event          Fire a new event
    exec           Executes a command on Consul nodes
    force-leave    Forces a member of the cluster to enter the "left" state
    info           Provides debugging information for operators.
    intention      Interact with Connect service intentions
    join           Tell Consul agent to join cluster
    keygen         Generates a new encryption key
    keyring        Manages gossip layer encryption keys
    kv             Interact with the key-value store
    leave          Gracefully leaves the Consul cluster and shuts down
    lock           Execute a command holding a lock
    maint          Controls node or service maintenance mode
    members        Lists the members of a Consul cluster
    monitor        Stream logs from a Consul agent
    operator       Provides cluster-level tools for Consul operators
    reload         Triggers the agent to reload configuration files
    rtt            Estimates network round trip time between nodes
    snapshot       Saves, restores and inspects snapshots of Consul server state
    validate       Validate config files/directories
    version        Prints the Consul version
    watch          Watch for changes in Consul
```

## 启动服务

consul 提供了3中部署方式

- dev 开发模式
- 单机模式
- 集群模式

作为入门呢，我们先使用开发（dev）模式。注意的是，改种方式不建议用于生产环境，因为该模式下所有的数据都不会进行持久化，都是在内存中的，该节点停止，那么配置的所有数据都会消失不见。但这却不失为一种快速体验功能的方式。

执行启动命令

```shell
./consul agent -dev
```

会有类似以下信息输出：

```shell
==> Starting Consul agent...
==> Consul agent running!
           Version: 'v1.2.2'
           Node ID: '016013da-7e3c-3757-8412-4199d46e025a'
         Node name: 'todu'
        Datacenter: 'dc1' (Segment: '<all>')
            Server: true (Bootstrap: false)
       Client Addr: [127.0.0.1] (HTTP: 8500, HTTPS: -1, DNS: 8600)
      Cluster Addr: 127.0.0.1 (LAN: 8301, WAN: 8302)
           Encrypt: Gossip: false, TLS-Outgoing: false, TLS-Incoming: false

==> Log data will now stream in as it occurs:
```

这样就是已经成功启动了一个dev模式的agent。

### KV 存储

执行 `./consul kv put a 'test value'` 可以在consul的存储kv系统中设置一个key为`a`值为`test value`的配置。
执行 `./consul kv get a` 可以读取我们刚才设置的这个值。

### 服务注册

通过agent的api可以直接像该agent注册服务，比如：

```shell
curl --request PUT \
  --url http://localhost:8500/v1/agent/service/register \
  --header 'accept: application/json' \
  --header 'content-type: application/json' \
  --data '{
  "ID": "myweb1",
  "Name": "myweb",
  "Tags": [
    "primary",
    "v1"
  ],
  "Address": "127.0.0.1",
  "Port": 8080,
  "Meta": {
    "redis_version": "4.0"
  },
  "EnableTagOverride": false
}'
```

这样，就像这个节点注册了一个名为 `myweb`，id为`myweb1`的服务。

通过 `./consul catalog services`可以看到列表中出现了我们的服务名字：

```shell
consul
myweb
```

### web页面

上面是通过命令，简单介绍了consul的两个基本特性 `kv存储` 和 `服务管理`。有个比较简单的方式就是可以通过web页面操作kv值和查看services。浏览器打开 `http://localhost:8500`即可，默认会跳转到services的页面，如图所示：
![consul-ui-start](http://public-links.qiniudn.com/image/consul/consul-ui-start.png)
点击Nodes链接，即可查看node节点简易信息：
![consul-ui-nodes](http://public-links.qiniudn.com/image/consul/consul-ui-nodes.png)
点击Key/Value即可查看kv存储：
![consul-ui-kv](http://public-links.qiniudn.com/image/consul/consul-ui-kv.png)
点击key还可以进行编辑：
![consul-ui-kv-edit](http://public-links.qiniudn.com/image/consul/consul-ui-kv-edit.png)

通过以上简单介绍，相比我们已经对consul有了初步认识。进阶操作和高级特性会通过后面的文章进行介绍。