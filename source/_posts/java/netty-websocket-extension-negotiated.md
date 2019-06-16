---
category: java
title: netty websocket extension negotiated
tags:
  - java
  - websocket
  - netty
abbrlink: c920
date: 2019-06-16 10:56:48
---

# netty websocket extension negotiated

这段时间在线上环境的日志中出现了如下错误日志：
![netty websocket extension negotiated](https://public-links.todu.top/1560676613.png?imageMogr2/thumbnail/!100p)。

日志所在模块是使用 [netty](https://netty.io/) 写的websocket服务，具体
关键配置代码如下：

![netty websocket 关键配置代码](https://public-links.todu.top/1560676932.png?imageMogr2/thumbnail/!100p)

其中 **1** 是 netty 的 websocket 配置，**2** 则是自己的业务处理逻辑，我们只需要关心 **1** 即可。

## **未完，待补充**
