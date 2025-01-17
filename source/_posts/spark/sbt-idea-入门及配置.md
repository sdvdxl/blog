---
title: sbt-idea-入门及配置
tags:
  - spark
  - scala
category: spark
abbrlink: 34842
date: 2016-05-10 09:35:48
---

# Java 环境配置
这个就不多说了，这是前提条件，请自行安装后配置正确，如果不清楚请自行搜索 **java 环境变量配置** 相关问题。
# Scala 配置
首先要要配置 scala 环境。从[官方](http://scala-lang.org/)地址下载，这里我们使用scala2.10.6版本，所以[从这里下载](http://scala-lang.org/download/2.10.6.html)对应的平台版本。Windows请下载 [scala-2.10.6.zip](http://downloads.lightbend.com/scala/2.10.6/scala-2.10.6.zip) ，MacOS和Linux请下载 [scala-2.10.6.tgz](http://downloads.lightbend.com/scala/2.10.6/scala-2.10.6.tgz) 。

下载完成后，解压到一个目录，然后配置环境变量 `SCALA_HOME` ，把scala的解压后的绝对路径配置到 `SCALA_HOME` ，然后增加 `PATH` 的配置。以下以 WIndows 和 Linux 举例来说，假如这里解压后得到的文件夹为 `scala-2.10.6`。

## Windows 环境
假如 `scala-2.10.6` 文件夹放在了 `C:\`目录下，那么新增 `SCALA_HOME` 的值为 `C:\scala-2.10.6` ，然后找到 `PATH` 这个环境变量， 在已有的值后面添加 `;%SCALA_HOME%\bin` （注意前面的分号），重新打开一个新的命令行窗口即可操作。关于windows环境变量的其他说明情自行补脑，这里就不多啰嗦了。

## Linux 和 Mac 环境
这里假定我们 `scala-2.10.6` 的 scala 目录存放在 `/usr/local/` 下，绝对路径就是 `/usr/local/scala-2.10.6` 。

### Bash 环境
bash 环境下可以修改 `~/.bashrc` 或者 `/etc/profile` 文件，添加一下内容：             

```bash   
SCALA_HOME=/usr/local/scala-2.10.6
PATH=$PATH:$SCALA_HOME/bin

export SCALA_HOME PATH
```
然后执行 `source ~/.bashrc` 或者 `source /etc/profile` 即可。

### Zsh 环境
如果你的终端bash用的是 zsh ，那么需要在 `~/.zshrc` 文件中增加上述内容，然后执行 `source ~/.zshrc` 即可。

# Sbt 配置
[官方地址](http://www.scala-sbt.org/)，[从这里下载](http://www.scala-sbt.org/download.html)，如果是MacOS的话，可以依照这里提示的方法进行快捷安装，如果是其他平台或者想手动配置，[点击这里直接下载](https://dl.bintray.com/sbt/native-packages/sbt/0.13.11/sbt-0.13.11.zip)即可。

下载完成并解压，得到文件夹 `sbt`。

## Windows 环境
假如我们将 `sbt` 文件夹放到了 `C:\`目录下。新增环境变量 `SBT_HOME` 值为 `C:\sbt`， 在 `PATH` 变量值后面添加 `;%SBT_HOME%\bin`，重新打开一个新的命令行窗口即可。
## Linux 或者 Mac 环境
这里假定我们将 `sbt` 目录放在了 `/usr/local/` 目录下。同上面配置 scala 环境变量一样。
### Bash 环境
编辑 `~/.bashrc` 或者 `/etc/profile` ，新增以下内容:
     
```bash
SBT_HOME=/usr/local/sbt
PATH=$PATH:$SBT_HOEM/bin

export SBT_HOME PATH
```
然后执行 `source ~/.bashrc` 或者 `source /etc/profile` 即可生效。
### Zsh 环境
编辑 `~/.zshrc`， 添加上面的内容并保存后，执行 `source ~/.zshrc` 即可生效。

## 示例程序
### 项目结构
在随意一个地方创建一个文件夹，名字为 `spark-sbt-demo` ，下面是目录结构：

```
spark-sbt-demo                                                                                                                                                                         
├── build.sbt
├── project
│   ├── build.properties
│   └── plugins.sbt
└── src
    ├── main
        ├── scala
           └── WordCount.scala
```
build.sbt 文件中添加如下内容（注意每行要用空行隔开）：

```scala
name := "spark-sbt-demo"

version := "1.0"

scalaVersion := "2.10.6"

organization := "spark.demo"

version := "1.0.0-SNAPSHOT"

libraryDependencies += "org.apache.spark" % "spark-core_2.10" % "1.6.1" % "provided"
```

WordCount.scala 文件内容：

```scala
import org.apache.spark.{SparkConf, SparkContext}

/**
  * Created by sdvdxl on 16/5/11.
  */
object WordCount {
  def main(args: Array[String]) {
    val conf = new SparkConf().setAppName("spark-sbt-demo").setMaster("local[*]")
    val sc  = new SparkContext(conf)
    sc.textFile("src/main/scala/WordCount.scala").flatMap(_.split(" ")).map(word=>(word,1)).reduceByKey(_+_).foreach(println)
    sc.stop()
  }
}

```

build.properties 文件内容 ：

```
sbt.version = 0.13.11
```

plugins.sbt 文件先保持为空。

至此，我们已经创建了一个 sbt 机构的项目。

接下来会说明使用sbt下载依赖，使用 idea 创建 sbt 项目，在idea中如何运行sbt管理的 spark app。

## Sbt 的基本使用
上面我们创建了一个用 sbt 管理的 spark app 项目，如果想要提交到spark中运行，那么需要打包成jar包，好在 sbt 本身或者插件提供了这样的功能。

### 应用打包
打开命令行，切换到该项目目录下，然后输入 `sbt` 之后，进入 sbt 的交互中，然后输入 `package` ，开始打包，最后如果看到类似
> [info] Done packaging.
  [success] Total time: 11 s, completed 2016-5-11 12:32:09
  
字样，那么说明打包成功，打成的 jar 包在上面的日志中可以找到。

### 第三方 jar 统一打包
在写应用的时候，我们不只是用到 `spark` 自身的 jar 包，还会用到好多其他第三方类库，那么，在提交应用到 spark 运行的时候，这些第三方依赖也需要一并提交上去，否则会出现找不到类的问题。如果依赖少的话，直接将这些 jar 包直接一个一个提交上去也没问题，但是一旦依赖了大量的类库，这种方法显然是低效费力的，那么怎么才能将这些所有的第三方依赖打成一个 jar 包呢？

sbt 本身没有提供这样的功能，但是我们可以依靠相应的插件完成此操作。记得上面有个文件内容留空的 `plugins.sbt` 文件吗？这个文件中可以配置我们想要完成特定功能的插件，现在我们在其中添加如下内容：

```scala
addSbtPlugin("com.eed3si9n" % "sbt-assembly" % "0.14.2")
```

然后重新 进入 sbt 交互式环境，输入 `assemblyPackageDependency` 回车，稍后将看到类似如下输出：    
> [info] Done packaging.
  [success] Total time: 41 s, completed 2016-5-11 13:36:37

这样就成功的将所有依赖的第三方类库打包到一个 jar 包中了，具体打包的文件可以在上面的日志中看到。

## 使用 idea 创建 sbt 项目
### 安装插件
使用 idea 创建 sbt 项目需要安装 `scala` 和 `sbt` 插件。
打开idea的首选项，然后找到 `Plugins` ，点击 `Browser repositores...` 按钮，输入 `scala` 搜索，然后找到 `scala` 和 `sbt` 的插件进行安装，如下图所示：
![scala-sbt-plugins](https://public-links.todu.top/images/scala/scala-sbt-plugin.png)
安装完成后重启idea。
### 创建 sbt 项目
File -> New -> Project... 打开项目创建向导：
![创建sbt项目](https://public-links.todu.top/images/scala/sbt-project-1.png)
创建完成后，等待idea刷新项目，目录结构大体如下（project/project 和 target相关没有列出）：

```
spark-sbt-demo                                                                                                                                                                
├── build.sbt
├── project
│   ├── build.properties
│   └── plugins.sbt
└── src
    ├── main
    │   ├── java
    │   ├── resources
    │   ├── scala
    │   └── scala-2.11
    └── test
        ├── java
        ├── resources
        ├── scala
        └── scala-2.11

```

- `plugins.sbt` 文件放置插件配置
- `build.sbt` 是整体的项目配置信息
- `build.properties` 可以设置 sbt 版本
- `java` 目录存放 java 文件
- `scala` 目录存放 scala 文件
- `resources` 目录用来存放配置文件
- `test` 相关目录用来存放测试相关文件
## 在 idea 中 运行 spark app
上面我们介绍了如何使用 idea 项目向导创建一个 sbt 项目，现在我们来说一下如何在 idea 中直接运行 sbt 构建的 spark app。

这里我们使用一开始我们创建的那个项目，使用 idea 导入功能，File -> Open 找到项目目录打开即可。
在 `WordCount.scala` 文件中右键，选择 `Run WordCount` ，开始运行，但是结果可能不是我们所期望的：

```java
Exception in thread "main" java.lang.NoClassDefFoundError: org/apache/spark/SparkConf
	at WorldCount$.main(WorldCount.scala:8)
	at WorldCount.main(WorldCount.scala)
	at sun.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
	at sun.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:62)
	at sun.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
	at java.lang.reflect.Method.invoke(Method.java:497)
	at com.intellij.rt.execution.application.AppMain.main(AppMain.java:144)
Caused by: java.lang.ClassNotFoundException: org.apache.spark.SparkConf
	at java.net.URLClassLoader.findClass(URLClassLoader.java:381)
	at java.lang.ClassLoader.loadClass(ClassLoader.java:424)
	at sun.misc.Launcher$AppClassLoader.loadClass(Launcher.java:331)
	at java.lang.ClassLoader.loadClass(ClassLoader.java:357)
	... 7 more

Process finished with exit code 1
```

这是为什么呢？原因是我们在 `build.sbt` 中配置的 spark 依赖是这样的：

```scala
libraryDependencies += "org.apache.spark" % "spark-core_2.10" % "1.6.1" % "provided"
```
注意到后面的 `provided` 了吗？这个代表打包或者运行的时候不会将这个 jar 包的文件包含进去（注意：spark app 要求这样，注意不要把spark相关的jar包包含进去）。这样导致我们无法再 idea 中调试或者运行 spark app。

解决方案还是有的，sbt 和 maven（也是一个项目管理的软件）一样，提供了模块开发功能，我们定义两个模块，一个模块就是我们上面我们做好的，另一个是用来运行的，这个里面包含了运行时类库，配置如下：
1. 创建一个名为 `main` 的文件夹，把项目中的 `src` 文件夹移动到这个目录下
2. 在项目根目录下创建名为 `run` 的文件夹
3. 修改项目根目录下的 `build.sbt` 文件，内容为：
```
name := "spark-sbt-demo"

version := "1.0"

scalaVersion := "2.10.4"

organization := "spark.demo"

version := "1.0.0-SNAPSHOT"

libraryDependencies += "org.apache.spark" % "spark-core_2.10" % "1.6.1" % "provided"

lazy val root = (project in file(".")).aggregate(main, run)

lazy val main = (project in file("main"))

lazy val run = (project in file("run")).dependsOn(main)
```
4. 在子项目 `main` 创建 `build.sbt` 内容为：
```scala
libraryDependencies += "org.apache.spark" % "spark-core_2.10" % "1.6.1" % "provided"
```
5. 在子项目 `run` 创建 `build.sbt` 内容为 ：
```
libraryDependencies += "org.apache.spark" % "spark-core_2.10" % "1.6.1"
```
6. 配置运行参数，如下图：
![idea-spark-run-config](https://public-links.todu.top/images/scala/idea-spark-run-config.png)
然后选择上面的运行配置，运行即可。这里可能会碰到一个异常：

```java
Exception in thread "main" org.apache.hadoop.mapred.InvalidInputException: Input path does not exist: file:/Users/du/workspace/hekr/spark-sbt-demo/src/main/scala/WorldCount.scala
	at org.apache.hadoop.mapred.FileInputFormat.listStatus(FileInputFormat.java:251)
	at org.apache.hadoop.mapred.FileInputFormat.getSplits(FileInputFormat.java:270)
	......
```

这是由于上面我们修改了改程序main文件的位置，导致找不到该文件所致，请自行设置为一个存在的文件路径或者修改为 `main/src/main/scala/WorldCount.scala` 重新运行即可成功。

## Sbt 本地依赖库存储位置配置

抽空再补上，其实就是建立一个连接，先自行思考方案。

# 项目下载

没有源码下载的都是耍流氓，[点这里下载](https://public-links.todu.top/files/sparkspark-sbt-demo.tar.gz)

