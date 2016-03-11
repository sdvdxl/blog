---
title: 运行第一个SparkStreaming程序（及过程中问题解决）
date: 2016-03-09 12:51:43
tags:
  - zookeeper
  - spark-streaming
  - spark
category: spark
---

# 官方示例说明
按照官方文档的 [这个示例说明](http://spark.apache.org/docs/1.0.0/streaming-programming-guide.html#a-quick-example)，可以轻松的在本地的spark-shell环境中测试这个示例。示例，即为了更好的入门，那么就再说明一下。
运行这个统计单词的方式有三种，前面两种是官方文档上的指引，第三种则是用scala程序运行。

----
- ## 第一种方式, run-demo
1. 打开一个终端，打开一个终端，输入 命令 ` nc -lk 9999 `，暂时叫做 “nc终端” 吧
2. 再打开终端，切换到Spark HOME目录， 执行命令 ` bin/run-example org.apache.spark.examples.streaming.NetworkWordCount localhost 9999 `， 然后每秒会有类似一下日志循环输出
>
\-------------------------------------------
Time: 1415701382000 ms
\-------------------------------------------
\-------------------------------------------
Time: 1415701383000 ms
\-------------------------------------------

2. 在nc终端随便输入一些字符串，用空格隔开，回车，如aa aa bb c。可以在上面的Spark终端中看到有新内容输出
>
\-------------------------------------------
Time: 1415701670000 ms
\-------------------------------------------
(aa,2)
(bb,1)
(c,1)

OK，成功！

----


- ## 第二种 spark-shell 模式
下面介绍在spark-shell中输入scala代码运行的方式。
1. 同上面第一步，打开一个终端，打开一个终端，输入 命令 ` nc -lk 9999 `，暂时叫做 “nc终端” 吧
1. 再打开一个终端， 切换到Spark HOME目录下，输入 ` bin/spark-shell ` （如果你已经安装好了Spark的话，直接输入 ` spark-shell ` 即可），等待Spark启动成功，会打印信息  
 >
Spark context available as sc.
scala>

 然后输入以下语句：
 

```
import org.apache.spark.streaming._
import org.apache.spark.streaming.StreamingContext._
import org.apache.spark.api.java.function._
import org.apache.spark.streaming._
import org.apache.spark.streaming.api._

// Create a StreamingContext with a local master
val ssc = new StreamingContext(sc, Seconds(1))

// Create a DStream that will connect to serverIP:serverPort, like localhost:9999
val lines = ssc.socketTextStream("localhost", 9999)

// Split each line into words
val words = lines.flatMap(_.split(" "))
import org.apache.spark.streaming.StreamingContext._

// Count each word in each batch
val pairs = words.map(word => (word, 1))
val wordCounts = pairs.reduceByKey(_ + _)

// Print a few of the counts to the console
wordCounts.print()
ssc.start()             // Start the computation
ssc.awaitTermination()  // Wait for the computation to terminate
```
 会打印以下信息：
>
14/11/11 18:07:23 INFO MemoryStore: ensureFreeSpace(2216) called with curMem=100936, maxMem=278019440
\......
14/11/11 18:07:23 INFO DAGScheduler: Stage 91 (take at DStream.scala:608) finished in 0.004 s
14/11/11 18:07:23 INFO SparkContext: Job finished: take at DStream.scala:608, took 0.007531701 s
 \-------------------------------------------
Time: 1415700443000 ms
\-------------------------------------------

2.  同第一种方式的第3步，随便输入一些字符串，用空格隔开，回车，如aa aa bb c。可以在上面的Spark终端中看到有新内容输出
>
\-------------------------------------------
Time: 1415701670000 ms
\-------------------------------------------
(aa,2)
(bb,1)
(c,1)

 OK，成功！
-----

- ## 第三种 scala-ide编程方式
在用这种方式运行这个demo代码的时候，遇到了不少问题，记录下来，供大家参考。这个例子，请大家先根据这里记录的方式进行操作，得到一个可以运行的程序，后面我会记录遇到的问题。

1. 下载scala-ide, [下载链接](http://scala-ide.org/download/sdk.html)，下载 For Scala 2.10.4 下的对应平台的ide，解压，运行。
2. 安装sbt，[下载链接](http://www.scala-sbt.org/download.html),
3. 安装sbteclipse, [github地址](https://github.com/typesafehub/sbteclipse), 编辑 ` ~/.sbt/0.13/plugins/plugins.sbt  ` 文件， 添加以下内容 ` addSbtPlugin("com.typesafe.sbteclipse" % "sbteclipse-plugin" % "2.5.0") `，如果没有plugins目录和plugins.sbt，自行创建。
4. 用向导创建一个scala项目，并在项目根目录下创建一个build.sbt文件，添加以下内容(注意，每行正式语句之后要换行)

 ```
 name := "spark-test"


 version := "1.0"


 scalaVersion := "2.10.4"


 // set the main class for the main 'run' task
 // change Compile to Test to set it for 'test:run'
 mainClass in (Compile, run) := Some("test.SparkTest")

 libraryDependencies += "org.apache.spark" % "spark-streaming_2.10" % "1.1.0"
```

5. 创建test.SparkTest.scala文件，添加以下代码

```java
package test
import org.apache.spark.streaming._
import org.apache.spark.streaming.StreamingContext._
import org.apache.spark.SparkContext
import org.apache.spark.api.java.function._
import org.apache.spark.streaming._
import org.apache.spark.streaming.api._


object SparkTest {
  def main(args: Array[String]): Unit = {
    // Create a StreamingContext with a local master
    // Spark Streaming needs at least two working thread
    val ssc = new StreamingContext("local[2]", "NetworkWordCount", Seconds(10))
    // Create a DStream that will connect to serverIP:serverPort, like localhost:9999
    val lines = ssc.socketTextStream("localhost", 9999)
    // Split each line into words
    val words = lines.flatMap(_.split(" "))
    // Count each word in each batch
    val pairs = words.map(word => (word, 1))
    val wordCounts = pairs.reduceByKey(_ + _)
    wordCounts.print
    ssc.start
    ssc.awaitTermination
  }
}
```

6. 终端中切换目录到这个项目根目录，输入命令 ` sbt ` ， 命令运行成功后，敲入 ` eclipse ` 生成eclipse项目和项目所需依赖
7. 同第一种方式的第1,3步，
 再打开一个终端，输入 命令 ` nc -lk 9999 `。
然后运行刚才写的main程序，在nc终端中输入一些字符串，用空格隔开，回车，如aa aa bb c。可以在ide控制台中观察到
 >
\-------------------------------------------
Time: 1415701670000 ms
\-------------------------------------------
(aa,2)
(bb,1)
(c,1)

OK，成功！

----

# 下面是遇到的问题及解决方法：
### 1. 运行程序说找不到主类
解：没有在sbt文件配置主类是哪个，在` build.sbt`  文件中添加以下代码
> mainClass in (Compile, run) := Some("test.SparkTest")

 Some中就是主类的路径

### 2. java.lang.NoClassDefFoundError: scala/collection/GenTraversableOnce$class
这个问题困扰了我很长时间，一直没找到怎么解决。后来看到说是scala每次版本升级不兼容以前的版本编译的库，于是换了对应的版本的ide才正常运行。
解：scala-ide版本和现在用的spark包依赖编译的scala版本不一致， 请下载上面说过的 ` scala-ide For Scala 2.10.4` 版本。
