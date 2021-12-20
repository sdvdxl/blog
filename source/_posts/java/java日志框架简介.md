---
title: java日志框架简介
draft: true
date: 2021-12-19 13:10:21
updateDate: 2021-12-19 13:10:21
top_img:
cover:
tags:
- java
- slf4j
- logback
- log4j
- 日志
keywords:
- java
- slf4j
- logback
- log4j
- 日志
- sdvdxl
---

# java日志框架简介

> 这个图片和本文没有任何关系，唯一有用的可能是让你的眼睛放松一下。

<img title="" src="https://public-links.todu.top/1639903881.png?imageMogr2/thumbnail/!100p" alt="" data-align="inline">

日志框架主要分为3大类

1. 门面日志框架（Logging Facade Framework）

2. 实现日志框架（简称日志框架 Logging Framework）

3. 日志桥接框架（或者叫做适配）

---

## 门面日志框架 Logging Facade

简而言之门面日志框架只提供接口，不负责日志的最终解析和输出。用户使用的时候只需要关心这个门面日志的接口，不需要调用其他日志框架接口。

主要阵营：

- slf4j（Simple Logging Facade for Java http://www.slf4j.org ）

- commons-logging （apache）

## 实现日志框架

最终负责日志解析和内容输出的框架，比如输入 log.inf("一行日志，用户名：{}", username)，它负责解析和输出为(实际跟日志格式定义有关)：

> 2021-12-19 21:28:30 INFO main 一行日志，用户名：sdvdxl 

主要阵营：

- Log4j
- JUL(jdk-logging)
- Log4j2
- Logback

## 桥接日志框架

多种日志实现框架混用情况下，需要借助桥接类进行日志的转换，最后统一成一种进行输出（下面主要针对的是slf4j的桥接）

- *slf4j-log4j12-${latest.stable.version}.jar* 绑定log4j 1.x版本

- *slf4j-jdk14-${latest.stable.version}.jar* 绑定 JUL包括 JDK 1.4 版本的 log

- *slf4j-nop-${latest.stable.version}.jar* 绑定无操作，忽略所有日志信息

- *slf4j-simple-${latest.stable.version}.jar* 绑定简单日志输出，仅仅将 INFO 及其以上的信息打印到 System.err

- *slf4j-jcl-${latest.stable.version}.jar* 绑定到 JCL （[Apache Commons Logging](http://commons.apache.org/logging/) ），将所有slf4j的日志委托给 jcl 打印

- *logback-classic-\${logback.version}.jar*  需要依赖 logback-core，原生实现了slf4j的绑定

下图摘录自slf4j官网

http://www.slf4j.org/manual.html#swapping

该图说明了使用 slf4j 、适配器和 日志框架的对应关系

![](https://public-links.todu.top/1639915616.png?imageMogr2/thumbnail/!100p)



## 总结

本篇只是对日志框架的组成进行了简单介绍， slf4j 和 适配器相关的内容放到下一篇。至于JCL（Apache Commons Logging），我并不打算细作研究，简单的对比也放到下篇中。
