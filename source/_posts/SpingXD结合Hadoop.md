---
title: SpingXD结合Hadoop
date: 2016-03-11 12:34:40
tags:
	- spring
	- spring xd
	- spring-xd
	- hadoop
	- hdfs
	- xd
category:
	- spring
	- hadoop
---
# 前言
之前的文章介绍了[Spring XD](/Spring-XD简介)，[以分布式方式运行Spring-XD](/以分布式方式运行Spring-XD)和[安装启动Hadoop集群](/安装启动Hadoop集群)的文章。本文将简单介绍（刚学，很LOW(⊙﹏⊙)b）SpringXD和Hadoop，HDFS结合配置和使用方法。

# 准备
已经按照上述提到的两篇文章或者其他资料搭建并运行了SpringXD和Hadoop。

# 配置SpringXD
1. 运行命令`jps`找到`AdminServerApplication`和`ContainerServerApplication`两项，结束进程`kill 进程pid`。
2. 编辑SpingXD的配置文件`server.yaml`，在`spring`节点下增加以下hadoop的配置信息：

```yaml
hadoop:
	<!-- 注意换成自己的hdfs地址 -->
    fsUri: hdfs://10.10.1.110:8020
    resourceManagerHost: 10.10.1.110
    resourceManagerPort: 8032
    yarnApplicationClasspath:
```
然后启动xd-admin`bin/xd-admin`和xd-container`bin/xd-container`。
3. 切换到hadoop用户下，使用命令`hadoop fs -mkdir /xd`创建目录，然后更改权限`hadoop fs -chmod -R 777 /xd`（如果没有配置hadoop的环境变量，则请进入hadoop的目录使用`bin/hadoop`命令代替`hadoop`）
4. 打开新的控制台，进入xd-shell交互环境，假如根据上面的Spring-XD配置文章配置了安全措施，那么还需要执行下面的命令`admin config server --uri http://xd-adminIP:9393 --username 用户名 --password 密码`进行授权后登录。
5. 创建stream，向hdfs中写入数据`stream create --name myhdfsstream1 --definition "time | hdfs" --deploy`，用命令`hadoop fs ls /xd/myhdfsstream1`即可看到有临时文件生成。
