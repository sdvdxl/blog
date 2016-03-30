---
title: sink-hdfs
date: 2016-03-14 16:10:23
tags:
  - hdfs
  - hadoop
  - spring
  - spring xd
category: spring
---

hdfs根据时间自动划分文件夹
`stream create --name dataset1 --definition "kafka --topic=kafka_test --zkconnect=10.10.1.20:2181 --queueSize=64  |hdfs --inputType=application/json --idleTimeout=10000 --partitionPath=dateFormat('yyyy/MM/dd/HH/mm')" --deploy`
其中，--partitionPath=dateFormat('yyyy/MM/dd/HH/mm')用来指定划分 策略，这个是说用年(四位)/月(两位)/天(2位)/时(2位)/分(2位)这种格式来划分
