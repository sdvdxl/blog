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

- key  要使用的 stream 的名字，可以随便起
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

1. 过滤 ID `1576719377792-0` 以后的数据

    ```bash
    XREAD streams mystream 1576719377792-0
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

**ID 兼容** XREAD 中 ID 只写前半部分的毫秒时间戳也是兼容的

## [XRANGE](https://redis.io/commands/xrange) 命令

XRANGE 也是用于读取数据，但是可以指定区间，基本用法：

```bash
XRANGE key start end [COUNT count]
```

- start 可以指定开始的 ID，闭区间
- end 结束的 ID，闭区间
- COUNT 限制返回的数量

比如要返回 mystream `1576750346024-0` 和 `1576750348981-0` 之间的数据，可以输入：

```bash
XRANGE mystream 1576750346024-0 1576750348981-0
```

如果要返回所有的数据可以输入：

```bash
XRANGE mystream - +
```

## [XREVRANGE](https://redis.io/commands/xrevrange) 命令

XREVRANGE 和 XRANGE 基本类似，但是 XREVRANGE 是将流的顺序逆转输出。这里不需要关注太多，不过要注意的是范围参数也和 XRANGE 相反 `+` 在前， `-` 在后， 如 ：`XREVRANGE mystream + -`

其中 `-` 代表 `0-0` 的ID， `+` 代表 `18446744073709551615-18446744073709551615` 的 ID。

## [XDEL](https://redis.io/commands/xdel) 命令

XDEL 用于从 stream 中删除删除指定 ID 的记录，基本用法如下：

```bash
XDEL key ID [...]
```

ID 可以同时指定 1 到多个。返回值是删除的数量。

需要注意的是，删除并不是真正立刻就删除数据，只是给这条数据打了一个删除标签，当1个 `macro-node` 中所有数据都被删除时，才会将这个 `macro-node` 删除掉，数据也就随之真正删除。大量删除操作会引发内存碎片化情况的出现（性能不会受到影响），（在将来的Redis版本中，如果给定的宏节点达到给定数量的已删除条目，我们可能会触发节点垃圾回收。）

举例： 删除 id 为 0-0 的数据 `XDEL mystream 0-0`

## [XTRIM](https://redis.io/commands/xtrim) 命令

`XTRIM` 命令用于裁剪 stream 的长度，基本用法如下：

```bash
XTRIM key MAXLEN [~] count
```

返回结果是裁剪掉的数量。

使用的时候，需要指定期望裁剪后的最大长度，使用 `MAXLEN` 关键词来指定，比如期望裁剪后的最大长度为1000：

```bash
XTRIM key MAXLEN 1000
```

注意：这里说的是裁剪后的**最大长度**，如果 stream 本身没有达到 `MAXLEN` 的值，那么裁剪后的 stream 的长度还是其真实大小，即 `XLEN key` 所看到的结果。

另外用法中包含了 `~`，这个代表裁剪的时候近似的进行裁剪，可以多出 `MAXLEN` 几十个，但是不能少于 `MAXLEN`。

如： 执行 `XTRIM mystream MAXLEN ~ 1000` 后，查看长度 `XLEN mystream` 结果可能是大于1000的。

时间复杂度：O(N) N是要删除的数据的数量。

因为条目是在包含多个条目的宏节点中组织的，这些条目可以通过一次释放释放，所以可以在极短的时间内完成。

## []() 命令

## macro nodes

为了提高内存效率，stream 由 `macro-node`（宏节点） 组成 [`radix tree`](https://en.wikipedia.org/wiki/Radix_tree)（[基数树](https://zh.wikipedia.org/wiki/%E5%9F%BA%E6%95%B0%E6%A0%91)）。
