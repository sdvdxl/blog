---
title: 安装Hbase分布式集群
tags:
  - hadoop
  - hbase
category: 大数据
abbrlink: 10952
date: 2017-02-15 17:44:43
---

**以下操作都是在hadoop这个用户下**

1. 下载最新版hbase，放到/home/hadoop目录下，解压，生成目录 hbase-1.2.4
2. 下载zookeeper，放到 /home/hadoop/zookeeper-3.4.9， 解压生成 zookeeper-3.4.9 目录
3. 编辑conf下hbase-site.xml
```xml
  <configuration>
  <property>
    <name>hbase.rootdir</name>
    <value>hdfs://master:9000/hbase</value>
  </property>
  <property>
    <name>hbase.cluster.distributed</name>
    <value>true</value>
  </property>
  <property>
    <name>hbase.zookeeper.quorum</name>
    <value>master,slave1,slave2</value>
  </property>
  <property>
    <name>hbase.zookeeper.property.dataDir</name>
    <value>/home/hadoop/zookeeper-3.4.9</value>
  </property>
</configuration>
```
其中 `hbase.zookeeper.quorum` 配置是hbase集群机器的名字，上面的值代表在master，slave1，和slave2上启动hbase和zookeeper
`hbase.zookeeper.property.dataDir` 是用于配置zookeeper安装目录的，这里我把zookeeper安装到了 `/home/hadoop/zookeeper-3.4.9`
4. 修改 conf下 hbase-env.sh
修改JAVA_HOME，将其指向具体的java安装目录，（最小版本要求是1.7）
5. 启动hbase
```bash
bin/start-hbase.sh
```
6. 测试
输入 `bin/hbase shell` 进入交互环境，
- 创建表 `create 't1',{NAME => 'f1', VERSIONS => 2},{NAME => 'f2', VERSIONS => 2}`
- 列出表 `list`
如果没有出错那么安装成功。
