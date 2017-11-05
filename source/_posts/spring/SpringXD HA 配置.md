---
title: SpringXD HA 配置
category: spring
tags:
  - xd
  - spring
abbrlink: 17742
date: 2016-06-29 16:34:48
---
[SpringXD官方文档](http://docs.spring.io/spring-xd/docs/1.3.1.RELEASE/reference/html) 上说的不是很清楚，而且有些配置（如 配置 `hadoop` namenode ha ）并没有在上面说明，只是简单的说明了怎么配置 `namenode` ，如果没有ha配置，那么在生产环境中会令人头痛。

# XD Admin HA
## 说明
在 [官方文档](http://docs.spring.io/spring-xd/docs/1.3.1.RELEASE/reference/html/#_configuring_spring_xd_for_high_availabilty_ha) 中，有说如何配置，就是通过启动多个`admin` ，然后通过 `zookeeper` 管理。
`Spring XD` 要求只有一个主节点来和 `Container` 交互，例如 `Stream` 发布等。同时，这些操作都是按顺序处理的。假如只有一个 `admin` ，那么就存在单点失败的风险，因此，在生产环境中推荐做法是启动 2 个或者多 `admin` 。注意：在有多个 `admin` 节点的时候，每个 `admin` 都可以处理 [REST](http://docs.spring.io/spring-xd/docs/1.3.1.RELEASE/reference/html/#REST-API) 请求，但是只有一个实例会作为 `Leader` 处理请求并更新运行时的状态。如果 `Leader` 宕掉，另一个可用的admin就会成为新的 `Leader` 来接管任务。当然，Spring XD 的HA不只是他自身要求ha，还需要依赖外部服务，如：`zookeeper`，`MessageBus` 等 HA 配置。

## 配置信息
如果要配置 `admin` 的 ha，那么启动多个 `admin` 即可，但是请注意，如果是在同一台机器上部署多个`admin`，需要在启动时候添加如下参数以防止和默认的端口（`9393`）冲突：
`--httpPort` 用来指定rest api端口
`--mgmtPort` 用来指定管理端口
如果在不同机器上启动，只需配置相同的配置文件，然后启动即可。

# XD Container HA
当添加 `Container` 的时候，Spring XD 可以动态水平扩展，也就是说不需要额外什么操作，只需像第一次启动 `Container` 一样输入命令 `bin/xd-container` 启动即可。

# XD Hadoop namenode HA
如果 xd 中创建 `stream` 或者其他任务是用到了 `hdfs` 的功能，那么要配置 `hadoop` 的`namenode` ，要在 `xd/config/servers.xml `中的 `spring.hadoop.fsUri` 的配置项中配置。注意，这个地方只允许配置一个 `host`，如果有备用 `namenode` ，是不允许配置在这个地方的。但是这样配置是有问题的，就是存在 `hadoop` 的 `namenode` 主从切换后 xd 的 `stream` 无法写入 `hdfs` 或者读取 `hdfs` 的故障。要解决这个问题，我们要再 `xd/config/hadoop.properties` 中配置如下配置项：
```
dfs.nameservices=MyCluster
dfs.ha.namenodes.MyCluster=nn1,nn2
dfs.namenode.rpc-address.MyCluster.nn2=hadoop-master1-host:8020
dfs.namenode.rpc-address.MyCluster.nn1=hadoop-master2-host:8020
dfs.client.failover.proxy.provider.MyCluster=org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider
```
其中，配置项中所有的 `MyCluster` 要换成自己项目中 hadoop 集群的名字，然后在 `xd/config/servers.xml` 中 `spring.hadoop.fsUri` 值配置成 `hdfs://MyCluster:8020`（注意`8020`端口换成自己配置的），`hadoop-master1-host` 和 `hadoop-master1-host` 换成自己集群的 hadoop的 master 的主机名字或者ip。这样配置后，重新启动 `admin` 和 `container` 就会自动检测 hadoop 的 `namenode` 主从，并自行切换。
如果要在 `xd-shell `中使用，需要登录shell之后（如果有安全设置，还需要先用密码登录成功），输入一下命令：
```
hadoop config props set --property dfs.nameservices=MyCluster
hadoop config props set --property dfs.ha.namenodes.MyCluster=nn1,nn2
hadoop config props set --property dfs.namenode.rpc-address.MyCluster.nn1=hadoop-master1-host:8020
hadoop config props set --property dfs.namenode.rpc-address.MyCluster.nn2=hadoop-master2-host:8020
hadoop config props set --property dfs.client.failover.proxy.provider.MyCluster=org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider
hadoop config fs --namenode hdfs://MyCluster

```
`MyCluster`，`hadoop-master1-host` 和 `hadoop-master2-host` 同上配置。对于 shell 来说，单纯配置 hadoop 的主 `namenode`也是可以的，因为这个配置只是对 shell 起作用。如果觉得每次打开shell都要输入上面几行配置太繁琐的话，可以将 `xd/config/hadoop.properties` 中配置的项目添加一份到`shell/config/hadoop.properties` 即可，这样在shell中操作hdfs，只需配置 `hadoop config fs --namenode hdfs://MyCluster:8020` 即可。

## 附言
对于如何在xd中使用hadoop的namenode ha配置，xd 官方文档中并未见说明，而是百般Goole之后得到的结果，而且由于xd资料尚少，搜索结果不佳，最后通过搜索关键词 `xd namenode fail` 才找到解决方案，[请参见这里](https://jira.spring.io/browse/XD-1745)。

