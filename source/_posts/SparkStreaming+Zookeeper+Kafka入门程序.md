---
title: SparkStreaming+Zookeeper+Kafka入门程序
date: 2016-03-09 12:51:43
tags:
  - zookeeper
	- spark-streaming
	- spark
category: spark
---

## 准备工作：
- 安装 [spark](http://spark.apache.org/)
- 安装 [zookeeper](http://zookeeper.apache.org/)
- 安装 [kafka](http://kafka.apache.org/)

## 开始工作
#### 1. 启动zookeeper
 打开终端，切换到 `zookeeper HOME` 目录， 进入conf文件夹，拷贝一份 `zoo_sample.cfg` 副本并重命名为 ` zoo.cfg `
 切换到上级的bin目录中，执行 `./zkServer.sh start` 启动zookeeper，会有日志打印
 > Starting zookeeper ... STARTED

 然后用 `./zkServer.sh status` 查看状态，如果有下列信息输出，则说明启动成功
 > Mode: standalone

 如果要停止zookeeper，则运行 `./zkServer stop` 即可

#### 2. 启动kafka
打开终端，切换到 `kafka HOME` 目录,运行 `bin/kafka-server-start.sh config/server.properties` 会有以下类似日志输出
  >
[2014-11-12 17:38:13,395] INFO [ReplicaFetcherManager on broker 0] Removed fetcher for partitions [test,0] (kafka.server.ReplicaFetcherManager)
[2014-11-12 17:38:13,420] INFO [ReplicaFetcherManager on broker 0] Removed fetcher for partitions [test,0] (kafka.server.ReplicaFetcherManager)

#### 3. 启动kafka生产者
重新打开一个终端，暂叫做 生产者终端，方便后面引用说明。切换到 `kafka HOME` 目录,运行 `bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test` 创建一个叫 `test` 的主题。

#### 4. 编写scala应用程序

``` scala

    package test
    import java.util.Properties
    import kafka.producer._
    import org.apache.spark.streaming._
    import org.apache.spark.streaming.StreamingContext._
    import org.apache.spark.streaming.kafka._
    import org.apache.spark.SparkConf


    object KafkaWordCount {
      def main(args: Array[String]) {
    //    if (args.length < 4) {
    //      System.err.println("Usage: KafkaWordCount <zkQuorum>     <group> <topics> <numThreads>")
    //      System.exit(1)
     //    }

    //    StreamingExamples.setStreamingLogLevels()

    //val Array(zkQuorum, group, topics, numThreads) = args
    val zkQuorum = "localhost:2181"
    val group = "1"
    val topics = "test"
    val numThreads = 2

    val sparkConf = new SparkConf().setAppName("KafkaWordCount").setMaster("local[2]")
    val ssc =  new StreamingContext(sparkConf, Seconds(2))
    ssc.checkpoint("checkpoint")

    val topicpMap = topics.split(",").map((_,numThreads)).toMap
    val lines = KafkaUtils.createStream(ssc, zkQuorum, group, topicpMap).map(_._2)
    val words = lines.flatMap(_.split(" "))

    val pairs = words.map(word => (word, 1))

    val wordCounts = pairs.reduceByKey(_ + _)

    //val wordCounts = words.map(x => (x, 1L))
    //  .reduceByKeyAndWindow(_ + _, _ - _, Minutes(10), Seconds(2), 2)
    wordCounts.print()

    ssc.start()
    ssc.awaitTermination()
  }
}

```

`build.sbt` 文件中添加依赖
 >
 libraryDependencies += "org.apache.spark" % "spark-streaming_2.10" % "1.1.0"
>
libraryDependencies += "org.apache.spark" % "spark-streaming-kafka_2.10" % "1.1.0"

启动scala程序，然后在 上面第2步的 生产者终端中输入一些字符串，如  `sdfsadf a aa a a a a a a a a` ，在ide的控制台上可以看到有信息输出
 >
4/11/12 16:38:22 INFO scheduler.DAGScheduler: Stage 195 (take at DStream.scala:608) finished in 0.004 s
\-------------------------------------------
Time: 1415781502000 ms
\-------------------------------------------
(aa,1)
(a,9)
(sdfsadf,1)

说明程序成功运行。
