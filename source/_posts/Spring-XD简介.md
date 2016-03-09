---
title: Spring-XD简介
date: 2016-03-09 15:30:24
tags:
  - spring-xd
  - spring
category: spring
---

# 简介
Spring XD is a unified, distributed, and extensible service for data ingestion, real time analytics, batch processing, and data export.
![distributed-overview.](/images/xd/xd-overview.png)
# Streams
翻译过来就是流，通过定义stream可以控制数据的流向，比如从MongoDB读取数据然后存储到HDFS中。
![stream](/images/xd/stream.png)
## 创建方式
一个简单的示例：该示例创建一个名字叫`ticktock`的stream，每秒钟产生一条时间信息然后通过管道传送到log中。
```
xd:> stream create --definition "time | log" --name ticktock
```
## 销毁Stream
```
xd:> stream destroy --name stream-name
```

# Modules
模块，当前包含`source`, `sink`,`processor`, 和`job`。

# Souces
数据源，Stream的来源，有以下几种方式：

方式|描述
----|------
File|文件方式
FTP|FTP方式
GemFire Continuous Query(gemfire-cq)|GemFire查询
GemFire source(gemfire)|GemFire文件
HTTP|http方式
JDBC Source(jdbc)|关系型数据库jdbc
JMS|JMS
Kafka|Kafka消息队列
Mail|通过接收电子邮件
MongoDB Source(mongodb)|MongoDB数据库
MQTT|MQTT
RabbitMQ|RabbitMQ消息队列
Reactor IP(reactor-ip)|Reactor IP(reactor-ip)
SFTP|SFTP
Stdout Capture| 标准输入
Syslog|系统日志
Tail | Tail程序
TCP|TCP
TCP Client(tcp-client)|TCP 客户端
Time|时间
Trigger Source(trigger)|触发器
Twitter Search(twittersearch)|Twitter搜索
Twitter Stream(twitterstream)|Twitter Stream流

# Sinks
数据源，Stream的输出，有以下几种方式：

方式|描述
----|------
Dynamic Router(router)|动态路由
File Sink(file)|文件方式
FTP Sink(ftp)|FTP方式
GemFire Server|GemFire服务器
GPFDIST|GPFDIST
Cassandra|Cassandra 数据库
Hadoop(HDFS) (hdfs)|hdfs文件系统
HDFS Dataset(Avro/Parquet) (hdfs-dataset)|hdfs文件系统中的avro或者parquet类型文件
JDBC Source(jdbc)|关系型数据库jdbc
Kafka Sink (kafka)|Kafka消息队列
Log | log 文件
Mail | Mail 发送
Mongo | Mongo数据库
MQTT Sink (mqtt) | MQTT
Null Sink(null) | null
RabbitMQ|RabbitMQ 消息队列
Redis | Redis
Shell Sink (shell) |shell
Splunk Server (splunk) | splunk
TCP Sink (tcp) |TCP

# Processors
可用的处理器包括`Aggregator``Filter``Header Enricher``HTTP Client``JSON to Tuple``Object to JSON``Script``Shell Command``Splitter``Transform`
- Aggregator -- 作用和 splitter相反，用于聚合，
- Splitter -- 用于拆解
- Filter -- 过滤器，用于中间处理数据
- Header Enricher (header-enricher) -- 用于添加头部信息
- HTTP Client -- 通过httpClient方式发送URL请求
- JSON to Tuple (json-to-tuple) -- 转换json数据到Tuple类型
- Object to JSON (object-to-json) -- 将对象转换为json格式
- Script 用于加载Groovy脚本
- Shell -- 用于加载Shell脚本
- Transform -- 用于负载类型转换

# Taps
监听器，窃听器。
不用重复定义相同的stream，然后监听此stream就可以做其他操作。并且可以用Label来分别对每个部分内容做个别名，定义Tab时候可以使用别名。如
```
stream create foo --definition "httpLabel: http | fLabel: filter --expression=payload.startsWith('A') | flibble: transform --expression=payload.toLowerCase() | log" --deploy
stream create fooTap --definition "tap:stream:foo.flibble > log" --deploy
```
上面对trasfrom部分做了一个别名，叫做`flibble`，然后下面定义一个Tap，并且最后指定是flibble这个标签，那么就是对`foo`这个stream的`flibble`做监听。

# Jobs
Job相比Stream不同点在于，Job算是静态的，Stream是动态的。Stream会持续接收数据，处理数据；Job是一次性接收数据，处理数据，如果数据改变，那么是不会进行处理的，除非有定时任务。
```
job create --name jobtest --definition 'timestampfile --directory=D:/jobs' --deploy
stream create --name time-cron --definition "trigger --cron='* * * * * *' > queue:job:jobtest" --deploy
```

# 使用counter
```
stream create foo --definition 'http --outputType=application/json | log'
stream create countName --definition 'tap:stream:foo > field-value-counter --fieldName=name' --deploy
stream deploy --name foo
http post --data {"name":"a"}
http post --data {"name":"a"}
http post --data {"name":"b"}
```
`field-value-counter list` 列出field-value-counter的名字
```
FieldValueCounter name
----------------------
countName
```

`field-value-counter display --name countName` 列出名字为`countName`的描述
```
FieldValueCounter=countName
---------------------------  -  -----
VALUE                        -  COUNT
a                            |  2
b                            |  1
```  
