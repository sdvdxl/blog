---
title: jds实现方式
abbrlink: 901e049d
date: 2022-01-07 22:54:02
updateDate: 2022-01-07 22:54:02
top_img:
cover:
category: java
tags:
- java
- 数据同步
keywords:
- java
- 数据同步
- sdvdxl
---

上一篇文章中我们介绍了JAVA 简单同步框架的诞生，这一次我们介绍下其是如何实现的。本文基于 `1.0.1-SNAPSHOT` 版本进行说明。

## 回顾

考虑[JDS JAVA Data Sync 数据同步框架 ](https://mp.weixin.qq.com/s?__biz=Mzg2NTU0NDQ1Ng==&mid=2247483743&idx=1&sn=2097324621ceff60950d918884ed743a&chksm=ce593138f92eb82e845c5436a8cdb8627ff521d5b7d795ad9774238ae736d80ee62eb1e70cc6&token=686176240&lang=zh_CN#rd) 中提到的几点：

1. 数据分发 ，偏移量持久化

2. 重试机制

3. 消息持久化

### 偏移量持久化

实际场景中服务会停机，服务重启的时候需要重新读取最新的数据偏移量状态（最后同步到了哪条记录），需要从最后同步的那条记录开始继续同步

### 重试

消息分发过程中可能会遇到各种问题，比如网络抖动，目标服务停止，业务逻辑出错等等。数据在某些场景下需要重发，而有些场景可能不需要重发；重发也需要考虑间隔问题，不能一直死循环重发

### 消息持久化

基于上面2条，自然还需要将消息进行持久化，否则服务重启会造成要同步的数据丢失。

## 架构

![](https://public-links.todu.top/1641607478.png?imageMogr2/thumbnail/!100p)

主要分为2大模块（逻辑上）

- Core 提供API接口和数据处理逻辑

- 存储层 提供数据持久化能力（包括待分发数据和偏移量）

### Core 模块

核心模块，提供了 API，数据处理的能力。接收用户数据，传输给底层存储；启动分发任务，调用用户逻辑进行数据处理；失败数据进行重试；分发成功调用存储层持久化偏移量；抽象存储层逻辑接口

### 存储层

该层要实现Core提供的接口，可以根据选择的组件进行具体的实现，比如实现基于Mongodb或者MySQL的存储层，甚至可以基于Kafka实现。

## 核心代码讲解

### 结构

该项目是一个标准的Maven项目。

在 `1.0.1-SNAPSHOT` 版本中，项目结构如下：

![](https://public-links.todu.top/1641624944.png?imageMogr2/thumbnail/!100p)

- api 就是Core模块

- xxx-provider 具体的存储层实现（该版本只实现了mongodb）

- 其他都是项目信息

### 代码逻辑入口

入口类为`DataSync` 主要方法有：

![](https://public-links.todu.top/1641624980.png?imageMogr2/thumbnail/!100p)

因为要配置参数，提供了一个builder `DataSyncBuilder`，负责源数据，目标地址，还有设置具体的provicer。

### 启动流程

调用 `DataSync.start()`，会调用 `init` 函数进行初始化，调用 `internalStoreDataHandler.internalStoreDataHandler` 进行provider的准备工作（比如连接数据库，创建表等），然后启动一个线程，调用 provider 获取数据 `internalStoreDataHandler.poll()`

![](https://public-links.todu.top/1641610143.png?imageMogr2/thumbnail/!100p)

其中提到的 provider 的 `internalStoreDataHandler` 对象是 `InternalStoreDataHandler`接口的一个实例，具体 provider 负责具体实现。

![](https://public-links.todu.top/1641610245.png?imageMogr2/thumbnail/!100p)

可以看到 `InternalStoreDataHandler` 接口的方法和 `DataSync`的方法基本一致，其实面向 `DataSync`的操作最终会委托给 provider 的 `InternalStoreDataHandler` 实现来处理。

### 数据分发逻辑

![](https://public-links.todu.top/1641611782.png?imageMogr2/thumbnail/!100p)

接着看 数据交给用户处理逻辑：

![](https://public-links.todu.top/1641612109.png?imageMogr2/thumbnail/!100p)

以上就是主要的数据分发逻辑，下面介绍 provider（基于mongo实现）的主要逻辑。

### Provider 设计

在 core api 模块中，给 provider 提供了一个接口 `InternalStoreDataHandler`，同时还提供了一个 抽象类 `InternalAbstractStoreDataHandler` 方便 provider 获取DataSync的一些共享数据，比如 builder 中的信息。![](https://public-links.todu.top/1641612467.png?imageMogr2/thumbnail/!100p)

### Mongodb-Provider实现

在 `mongodb-provider`模块中，`MongoStoreDataHandlerImpl`类实现了`InternalAbstractStoreDataHandler`。接下来我们来看下主要内容。



MongoStoreDataHandlerImpl 这里使用了 `spring-data-mongodb`依赖，提供了一个基于 MongoTemplate 的 构造方法，传入该参数，可以使用其操作mongodb数据库；同时还有一个collectionName参数，该参数是用于保存数据的collection名字。![](https://public-links.todu.top/1641622438.png?imageMogr2/thumbnail/!100p)

### 存储数据逻辑

存储数据需要设置来源（所有数据都可以放在同一个collection中保存），这样在取数据的时候能够区分不同数据；然后使用 mongoTemplate save 方法保存数据。

![](https://public-links.todu.top/1641622719.png?imageMogr2/thumbnail/!100p)

### 查询数据

![](https://public-links.todu.top/1641623314.png?imageMogr2/thumbnail/!100p)

### offset 处理

![](https://public-links.todu.top/1641623431.png?imageMogr2/thumbnail/!100p)

## 框架现存问题

### mongodb-provider 问题

mongodb-provider 该存储层在实现比较offset的时候，使用的是 mongodb 的ObjectID，而ObjectId的构造在分片集群下不能保证严格的插入顺序，所以这个地方只适合单击或者副本集的mongodb server。

优化的的方案是，可以使用自定义的一个collection，使用一条记录，使用inc函数，这样会严格递增，保证分片下也是严格有序，不会因为比较大小而造成数据遗漏。

### 重试机制

如果数据发送失败，现在固定写死了等待3s后重试。其实可以使用策略模式，给用户暴露一个配置重试策略的接口。封装好的重试框架有：

- spring-retry

- guava-retrying

## 总结

1. 使用的设计模式：构建器 builder 模式； 策略模式（重试策略）; xxx-provider

2. 并发控制： synchronize；LockSupport 暂停，唤醒线程

3. 面向接口编程：internalStoreDataHandler

## 参考资料

1. jds 架构 https://gitmind.cn/app/flowchart/1235286165

2. Mongodb ObjectID 说明 https://docs.mongodb.com/manual/reference/bson-types/#std-label-objectid

3. spring-retrying https://docs.spring.io/spring-batch/docs/current/reference/html/retry.html

4. guava-retrying https://github.com/rholder/guava-retrying
