---
title: redis-stream
category: redis
tags:
  - redis
  - mq
  - stream
abbrlink: 6d35
date: 2019-12-16 18:47:32
updateDate: 2019-12-16 18:47:32
---

# Redis Stream


## 何为 redis stream

[Redis Stream](https://redis.io/topics/streams-intro) 是 redis 5.0 版本引入的一种新数据类型，可以认为是一个消息队列，但是相比List实现的消息队列功能又更为强大。下面结合官方教程进行简单介绍。详细介绍可以[参见这里](https://redis.io/topics/streams-intro)。

## Redis 安装

上面介绍说过，stream是5.0开始才有的特性，如果要使用stream，那么首先要先安装（或者升级）redis到5.0以后，这里我们可以从官网下载截止目前[最新的版本5.0.7](http://download.redis.io/releases/redis-5.0.7.tar.gz)，如果官网下载速度慢，可以[从这里下载](https://public-links.todu.top/soft/redis-5.0.7.tar.gz)。

我这里的环境是 `OSX 10.15.2`，其他linux系统理论上都可以使用，windows 未安装测试。

下载完成后解压 命令：`tar -xvf redis-5.0.7.tar.gz` 并得到 目录`redis-5.0.7` 结构如下：

<img src="https://public-links.todu.top/1576644448.png?imageMogr2/thumbnail/!100p" alt="redis-5源码目录结构" style="zoom:50%;" />

执行命令`cd redis-5.0.7` 进入到目录内，执行命令 `make` 进行编译，当最后输出：

> Hint: It's a good idea to run 'make test' ;)

的时候，说明编译通过了，并且在 `src` 目录中会生成 `redis-server` 文件。

执行 `src/redis-server` 启动redis，在另一个终端中输入 `src/redic-cli` 进入交互界面。

## Stream 的使用

### [XADD](https://redis.io/commands/xadd) 命令

`XADD` 命令可以在指定的一个stream中追加新的数据，主要用法如下：

`XADD key ID field string [field string ...]`

- key  要使用的 stream 的名字
- ID 要生成的id，可以自己指定，格式 `1526919030474-55`，使用 `-` 分隔，前后都是64bit的正整数，注意，后面插入的数据的id必须大于上次的id，比如 先插入`1001-1`，在插入 `1-1` 是不合法的。也可以使用 `*` 自动生成，当自动生成的时候，格式是：前部分是unix毫秒时间戳，后面是在这1ms内的序号（从0开始），一般来说都使用自动生成。
- field 和 string  是一对，必须成对出现，类似于map，field是key的名字，string 是值，只能使用string类型。只要成对，后面可以1到多个。

示例：

输入： `XADD mystream * sensor-id 1234 temperature 19.8`， 输出 `"1576646561236-0"`（类似结构，这个值是我执行时候的当前时间戳）

解释： 将 `sensor-id 1234 temperature 19.8` 追加到名为 `mystream` 的stream中，并且id使用自动生成。这个地方是模拟了一个传感器上报数据，传感器id是 1234， 当前温度 19.8。

类似的，我们可以再增加几条数据： `XADD mystream * sensor-id 0001 temperature 10.1`，`XADD mystream * sensor-id 0002 temperature 20.1`，同样输出一个自动生成的id。

## [XLEN](https://redis.io/commands/xlen) 命令

`XLEN` 命令可以查看对应的 stream 的长度，也就是该stream内的元素个数；注意这里的元素是指通过 `XADD` 这个命令插入的消息，比如上面 `XADD mystream * sensor-id 1234 temperature 19.8` `sensor-id 1234 temperature 19.8` 是一个元素。

基本用法：

XLEN key

key 就是要查看的 stream 的名字。

要查看上面我们插入的名为 `mystream` stream 的长度，则可以输入 `XLEN mystream`，输出的数字就是该stream当前长度。

## [XREAD](https://redis.io/commands/xread) 命令

XREAD 命令可以读取stream中的数据，基本用法如下：

XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] ID [ID ...]

解释：

- XREAD 命令
- COUNT 一次性读取的数量，后面跟正整数
- BLOCK 读取的时候，如果没有符合条件消息，则进行阻塞等待的时间毫秒数，如果直到超时都没有收到新的消息就返回空，否则返回新的消息
- STREAMS 固定格式
- key stream 名字，可以指定1到多个stream读取
- ID 要开始过滤的消息ID，此 ID 和 XADD 命令产生的 ID 是同一个东西。注意，这里的过滤条件是 大于 ID。如果指定了多个stream，这里也要分别指定每个stream要开始的ID。

### XREAD 举例

1. 要读取上面我们添加到 stream `mystream` 的所有数据：

    ```bash
    XREAD streams mystream 0-0
    ```

1. 自定义返回的数据条数

    上面的语句执行后会将整个 stream 内的数据都读取出来，为避免一次性读取出所有，我们可以 `COUNT` 关键词来约束返回的条数，这个约束的是 *最多返回的条数*

    ```bash
    XREAD COUNT 1 streams mystream $
    ```

    `COUNT` 后面的数字是约束要返回的条数，这里我们指定为 `1` ，代表最多返回1条。 如果 stream 内有2条，但是指定 `COUNT 10`，那么也是返回2条。

1. 要读取stream最新的消息

    ```bash
    XREAD streams mystream $
    ```

    注意，这里用 `$` 这个特殊的 ID 来说明要读取的是最新的消息。

    这个示例会直接返回 `(nil)`，因为没有加 `BLOCK` 是直接返回的，即使是没有数据。

1. 读取最新的消息，如果没有新消息则阻塞10s

    ```bash
    XREAD BLOCK 10000 streams mystream $
    ```

    `BLOCK` 后面的数值 10000 是毫秒数。 如果在这个10s内没有新的数据流入，那么10s后输出

    > (nil)
    > (10.15s)

    重新执行 `XREAD BLOCK 10000 streams mystream $`，并且在另一个终端中进入redis-cli交互，输入 `XADD mystream * sensor-id 1234 temperature 19.8` 则就会打印出刚才输入的内容。

    ```bash
    1) 1) "mystream"
    2) 1) 1) "1576719377792-0"
         2) 1) "sensor-id"
            2) "1234"
            3) "temperature"
            4) "19.8"
    (1.75s)
    ```
1.
