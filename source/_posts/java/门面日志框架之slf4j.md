---
title: 门面日志框架之slf4j
abbrlink: c95e5ae8
date: 2021-12-20 23:45:40
updateDate: 2021-12-20 23:45:40
top_img:
cover:
category:
tags:
keywords:
---

[在上一篇中，我们介绍了各种日志框架和其作用](http://mp.weixin.qq.com/s?__biz=Mzg2NTU0NDQ1Ng==&mid=2247483717&idx=1&sn=213bc8bfe80292950d1d1641691a9412&chksm=ce593122f92eb8347e2d821d0ecf66fa164d2f042b50822c16d782212883624313e40f58b1a6#rd)，本篇将介绍门面日志框架中的slf4j。

## 简介

slf4j 是 Simple Logging Facade for Java 的缩写（java简单门面日志），其官方网站为 http://www.slf4j.org/ 。它作为各种实现日志框架的一个接口存在，允许用户在部署的时候方便的替换具体的实现日志框架，而不需要重新修改代码中的log部分。

## 其他门面框架对比

除了slf4j，还有一个非常常见的门面框架JCL（Apache Jakarta Commons Logging）。相比JCL， slf4j使用的是静态绑定，JCL使用的是ClassLoader方式。slf4j的方式避免了 JCL 中存在的因为Classloader导致的内存泄露问题，并且更为高效。

## 版本

最新稳定 [1.7.32](https://mvnrepository.com/artifact/org.slf4j/slf4j-api/1.7.32)， 开发版最新为2.0.0（使用了java8的 serviceloader 方式加载绑定）。这里我们使用的是最新稳定版 1.7.32 。

## 使用

使用slf4j方式非常简单，只需要引入 `*slf4j-api-${latest.stable.version}.jar*` 即可，maven方式：

```xml
<dependency>
    <groupId>org.slf4j</groupId>
    <artifactId>slf4j-api</artifactId>
    <version>1.7.32</version>
</dependency>
```

gradle方式: `implementation 'org.slf4j:slf4j-api:1.7.32'`

hello world:

```java
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class HelloWorld {
  public static void main(String[] args) {
    Logger logger = LoggerFactory.getLogger(HelloWorld.class);
    logger.info("Hello World");
  }
}
```

编译并运行Hello World将在控制台上打印以下输出。

> SLF4J: Failed to load class "org.slf4j.impl.StaticLoggerBinder".
> SLF4J: Defaulting to no-operation (NOP) logger implementation
> SLF4J: See http://www.slf4j.org/codes.html#StaticLoggerBinder for further details.

因为只引入了接口，但是没有具体的日志实现，所以log4j提示没有绑定具体的日志实现，默认就是不输出。

## 引入实现日志框架

实际开发中，肯定会引入一个日志框架来进行日志的打印。常用的日志框架是 logback 和 log4j2。

相比log4j2，logback 原生支持slf4j，只需要引入 logback-classic 即可。这里也使用 logback 作为日志实现框架。

引入maven依赖：

```xml
<dependency>
    <groupId>ch.qos.logback</groupId>
    <artifactId>logback-classic</artifactId>
    <version>1.2.7</version>
</dependency>
```

同时在 `src/main/resources`文件夹下创建文件 `logback.xml`，内容为：

```xml
<?xml version="1.0" encoding="UTF-8"?>


<configuration>
  <property name="FILE_NAME" value="mqtt"/>

  <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
    <encoder>
      <!--      @formatter:off-->
      <Pattern>%date{"yyyy-MM-dd'T'HH:mm:ss.SSSZ"} %highlight(%-5level) [%thread] %logger{35}.%M:%line - %msg %n
      </Pattern>
<!--      @formatter:on-->
      <charset>UTF-8</charset>
    </encoder>
  </appender>


  <root>
    <level value="DEBUG"/>
    <appender-ref ref="CONSOLE"/>
  </root>
</configuration>
```

重新运行HelloWorld，就会看到控制台输出了日志内容：

> 2021-12-23T00:06:09.773+0800 INFO  [main] top.todu.log.slf4j.demo.HelloWorld.main:9 - Hello World 

### slf4j 对其他框架的支持

slf4j 对其他框架提供了比较广泛的支持

- *slf4j-log4j12-${latest.stable.version}.jar* 绑定log4j 1.x版本，使用 log4j1 输出日志

- *slf4j-jdk14-${latest.stable.version}.jar* 绑定 JUL包括 JDK 1.4 版本的 log，使用jdk输出日志

- *slf4j-nop-${latest.stable.version}.jar* 绑定无操作，忽略所有日志信息

- *slf4j-simple-${latest.stable.version}.jar* 绑定简单日志输出，仅仅将 INFO 及其以上的信息打印到 System.err

- *slf4j-jcl-${latest.stable.version}.jar* 绑定到 JCL （[Apache Commons Logging](http://commons.apache.org/logging/) ），将所有slf4j的日志委托给 jcl 打印

- *logback-classic-${logback.version}.jar* 需要依赖 logback-core，原生实现了slf4j的绑定

- log4j-slf4j-impl 绑定到 log4j2，使用log4j打印日志

下图摘录自slf4j官网

http://www.slf4j.org/manual.html#swapping

该图说明了使用 slf4j 、适配器和 日志框架的对应关系

![](https://public-links.todu.top/1639915616.png)

### slf4j 绑定原理

我们从上面这个HelloWorld入手。

1. Logger logger = LoggerFactory.getLogger(HelloWorld.class); 查看 `getLogger` 方法
   
   ![](https://public-links.todu.top/1640234551.png?imageMogr2/thumbnail/!100p)
   
   进入之后看到内部调用了 `getLogger` 方法
   
   ![](https://public-links.todu.top/1640234604.png?imageMogr2/thumbnail/!100p)
   
   继续跟进（进入 `getLogger` 这个方法）
   
   ![](/Users/du/Library/Application%20Support/marktext/images/2021-12-23-12-45-15-image.png)

分为2步

1. 获取logFactory

2. 使用这个logFactory 获取真正的logger

具体话， iLoggerFactory 实际是一个接口，只有1个方法，那就是 `getLogger`

![](https://public-links.todu.top/1640234925.png?imageMogr2/thumbnail/!100p)

这个 ILoggerFactory 也是绑定的关键，其他日志框架通过实现这个接口，可以适配到slf4j上，从而跟实现日志框架绑定。

下面接着看，slf4j是怎么选择日志框架的？

回到刚才 `getILoggerFactory`的位置，进入 `getILoggerFactory`方法

![](https://public-links.todu.top/1640235627.png?imageMogr2/thumbnail/!100p)

这个地方分为2部分

1. 第一部分是为了初始化。还记得单例模式实现方式中的 **双重校验**（不要忘记加关键字 `volatile`）吗？没错，这里使用了这个技巧，进行了延迟初始化，只有第一次调用的时候才会进行 `performInitialization()`方法调用。

2. 第二部分是对不同的状态进行判断和对应逻辑处理。这个地方 `INITIALIZATION_STATE` 的状态是由第一步初始化工作决定的。

我们接着看初始化过程是怎么进行的，继续跟进到 `perforInitialization` 方法，有个bind，这个是真正进行绑定的方法；下面那段是校验版本兼容性，这里就不细说了。

![](https://public-links.todu.top/1640269277.png?imageMogr2/thumbnail/!100p)

跟进到 `bind`方法：

![](https://public-links.todu.top/1640270326.png?imageMogr2/thumbnail/!100p)

可以分为6部分：

1. 查找运行时所有静态绑定类的路径，如果找到多个，输出信息

2. 调用`StaticLoggerBinder.getSingleton();`进行绑定；注意 `StaticLoggerBinder` 这个类并不在slf4j api 包中，是由日志框架或者适配框架进行提供，并且包名和类名都要完全一致，这个就是静态绑定的精髓。所以，如果环境中有多个`StaticLoggerBinder`，那么实际运行哪个，是由jvm决定的（可以认为是随机的）。比如我们只引入logback-classic，那么这个地方就是logback提供的这个类，getInstance() 的时候会实际调用logback去进行初始化。成功后设置状态为初始化成功，并且打印实际绑定的是哪个框架。如果没有任何适配器，那么这个地方运行时就缺少`StaticLoggerBinder`类，抛出 `NoClassDefFoundError` 错误后就会进入第3步

3. 进入这一步，说明没有框架可以绑定，就会设置为 NOP_FALLBACK_INITIALIZATION 状态（没有实际绑定的一个状态）

4. 这一步说明版本信息不正确，绑定失败

5. 其他异常绑定失败

6. 进行一些清理工作

最重要的就是2和3，完成之后我们跳回之前的代码

![](https://public-links.todu.top/1640271385.png?imageMogr2/thumbnail/!100p)

1. 这一步初始化完成，进行2

2. 如果状态是成功初始化，则直接返回绑定的框架的LoggerFactory，后续处理都交给实现框架（或者适配器）处理（这个跟具体的框架有关）

3. 如果是没有绑定，则状态是 `NOP_FALLBACK_INITIALIZATION`，就会返回一个 `NOPLoggerFactory` 实例，这个实例的方法都是空实现，也就不会输出任何日志内容；

## 总结

1. slf4j 使用 编译和运行时的类绑定到具体的日志框架上

2. slf4j中使用了双重检查机制延迟初始化

## 思考

### 为什么slf4j可以使用静态绑定，或者什么场景下可以使用类似实现方式？

考虑到日志框架，配置是全局的，也就是只需要初始化一次配置即可，而不需要用户手动进行初始化。

场景：如果有遇到全局使用同一份配置，或者由同一份配置的工厂产生对象，则可以是使用此方式

### 在判断状态的时候，slf4j 判断了 `ONGOING_INITIALIZATION` 状态，为什么双重检查还会出现这个状态？

按理说不同线程会被锁住，相同线程会按顺序执行，应该不会出现状态不一致问题。

根据作者说的问题描述 http://jira.qos.ch/browse/SLF4J-97 ，这个地方可能跟logback有关，尚未深入研究。
