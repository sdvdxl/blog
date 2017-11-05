---
title: Apache-Sqoop 安装
tags:
  - sqoop
category: sqoop
abbrlink: 27704
date: 2016-03-09 12:51:43
---
# 准备
首先当然是[下载sqoop](http://archive.apache.org/dist/sqoop/1.4.4/sqoop-1.4.4.bin__hadoop-2.0.4-alpha.tar.gz)
sqoop 依赖以下软件,点击链接可以直接下载
>
[hadoop](http://ftp.yz.yamagata-u.ac.jp/pub/network/apache/hadoop/common/hadoop-2.6.0/hadoop-2.6.0.tar.gz)
[accumulo](http://ftp.yz.yamagata-u.ac.jp/pub/network/apache/accumulo/1.6.2/accumulo-1.6.2-bin.tar.gz)
[apache-hive](http://ftp.tsukuba.wide.ad.jp/software/apache/hive/hive-1.0.0/apache-hive-1.0.0-bin.tar.gz)
[hbase](http://ftp.kddilabs.jp/infosystems/apache/hbase/hbase-1.0.0/hbase-1.0.0-bin.tar.gz)
[zookeeper](http://mirror.bit.edu.cn/apache/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz)

# 配置
## 配置JAVA环境变量
JAVA_HOME=/home/du/software/dev/jdk1.7.0_45
``` bash
export JAVA_HOME=/usr/install/java  #此处换成自己的jdk目录
export CLASSPATH=.:$JAVA_HOME/jre/lib
export PATH=$PATH:$JAVA_HOME/bin
```

## 配置sqoop运行依赖
``` bash
export HADOOP_COMMON_HOME=/home/du/software/dev/hadoop-2.6.0
export HADOOP_MAPRED_HOME=$HADOOP_COMMON_HOME/share/hadoop/mapreduce
export ZOOKEEPER_HOME=/home/du/software/dev/zookeeper-3.4.6
export ACCUMULO_HOME=/usr/install/accumulo-1.6.2
export HIVE_HOME=/usr/install/apache-hive-1.0.0-bin
export HCAT_HOME=/usr/install/apache-hive-1.0.0-bin/hcatalog
export HBASE_HOME=/usr/install/hbase-1.0.0
export SQOOP_HOME=/usr/install/sqoop-1.4.4.bin__hadoop-2.0.4-alpha
```

# 测试
切换到sqoop目录，运行 `bin/sqoop help`， 如果打印帮助文档则说明成功。
