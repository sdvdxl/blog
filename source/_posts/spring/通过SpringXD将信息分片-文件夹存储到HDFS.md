---
title: '通过SpringXD将信息分片|文件夹存储到HDFS'
date: 2016-03-16 17:29:03
tags:
  - spring xd
  - hdfs
  - kafka
category: spring
---
	kafka --topic=kafka_test --zkconnect=10.10.1.20:2181 --queueSize=64 |hdfs --inputType=application/json --idleTimeout=10000 --partitionPath=dateFormat('yyyy/MM/dd/HH/mm')
