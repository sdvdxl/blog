---
title: 【草稿】centos7安装haproxy
p: haproxy/centos7安装haproxy.md
category: Haproxy
tags:
  - haproxy
  - 负载均衡
  - lb
  - loadblance
abbrlink: 7dba
date: 2019-12-21 22:38:50
updateDate: 2019-12-21 22:38:50
---

# centos7 安装 Haproxy

[下载 haproyx2.1.0](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/haproxy/haproxy-2.1.0.tar.gz)

```bash
wgete https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/haproxy/haproxy-2.1.0.tar.gz
```

解压

```bash
tar -xvf haproxy-2.1.0.tar.gz
```

安装

为了支持 SSL ，必须安装 OpenSSL： `yum -y install openssl-devel.x86_64`

编译：

```bash
cd haproxy-2.1.0
make -j $(nproc) TARGET=linux-glibc USE_RT=1 ARCH=x86_64 USE_OPENSSL=1 ADDLIB=-lz USE_LINUX_TPROXY=1 USE_ZLIB=1
```

如果报错：(只有在内核低于3.0才会发生)：

```text
A raw syscall is useless, setns() is only supported in linux >= 3.0.
```

则需要指定参数： `USE_NS=`

保存以下配置到配置文件 `haproxy.cfg`：

```ini
global
  ulimit-n 10400009
  maxconn 99999
  maxpipes 99999
  tune.maxaccept 500
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice


defaults
  log global
  mode http
  option dontlognull
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms
# Listen to all MQTT requests (port 1883)
listen mqtt
  mode tcp
  bind *:1883 #要绑定的网卡和端口
  option tcplog
  # balance mode (to choose which MQTT server to use)
  balance leastconn
  server node_1 10.10.164.229:1883 check # node_1 为服务器节点别名， 绑定 后端的ip和端口，并开启健康检查
  option clitcpka # TCP keep-alive
  timeout client 3h #By default TCP keep-alive interval is 2hours in OS kernal, 'cat /proc/sys/net/ipv4/tcp_keepalive_time'
  timeout server 3h #By default TCP keep-alive interval is 2hours in OS kernal
  option tcplog
```

启动

./haproxy -f haproxy.cfg
