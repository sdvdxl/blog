---
title: redis知识点
abbrlink: ae04
date: 2021-11-21 14:19:56
updateDate: 2021-11-21 14:19:56
tag_img:
category:
tags:
keywords:
---

入门级 Redis 试题
============

试题一：为什么使用 redis？
----------------

**分析**: 博主觉得在项目中使用 redis，主要是从两个角度去考虑: **性能**和**并发**。当然，redis 还具备可以做分布式锁等其他功能，但是如果只是为了分布式锁这些其他功能，完全还有其他中间件 (如 zookpeer 等) 代替，并不是非要使用 redis。因此，这个问题主要从性能和并发两个角度去答。
**回答**: 如下所示，分为两点
**（一）性能**
我们在碰到需要执行耗时特别久，且结果不频繁变动的 SQL，就特别适合将运行结果放入缓存。这样，后面的请求就去缓存中读取，使得请求能够**迅速响应**。

**（二）并发**
在大并发的情况下，所有的请求直接访问数据库，数据库会出现连接异常。这个时候，就需要使用 redis 做一个缓冲操作，让请求先访问到 redis，而不是直接访问数据库。

参考资料
----

**疯狂创客圈 经典图书 ： 《[Netty Zookeeper Redis 高并发实战](https://www.cnblogs.com/crazymakercircle/p/13878143.html)》 面试必备 + 面试必备 + 面试必备**

Redis 基础试题
==========

1. Redis 有哪些数据结构？
-----------------

字符串 String、字典 Hash、列表 List、集合 Set、有序集合 SortedSet。
如果你是 Redis 中高级用户，还需要加上下面几种数据结构 HyperLogLog、Geo、Pub/Sub。
如果你说还玩过 Redis Module，像 BloomFilter，RedisSearch，Redis-ML，面试官得眼睛就开始发亮了。

2. 使用过 Redis 分布式锁么，它是什么回事？
--------------------------

先拿 setnx 来争抢锁，抢到之后，再用 expire 给锁加一个过期时间防止锁忘记了释放。
但是如果在 setnx 之后执行 expire 之前进程意外 crash 或者要重启维护了，那会怎么样？
可以使用set的参数控制，比如 `SET lock 1 PX 1000 NX` 使用一条指令合并setnx和expire。

3. Redis 里面有 1 亿个 key，其中有 10w 个 key 是以某个固定的已知的前缀开头的，如何将它们全部找出来？
---------------------------------------------------------------

使用 keys 指令可以扫出指定模式的 key 列表。
对方接着追问：如果这个 redis 正在给线上的业务提供服务，那使用 keys 指令会有什么问题？
这个时候你要回答 redis 关键的一个特性：redis 的单线程的。keys 指令会导致线程阻塞一段时间，线上服务会停顿，直到指令执行完毕，服务才能恢复。这个时候可以使用 scan 指令，scan 指令可以无阻塞的提取出指定模式的 key 列表，但是会有一定的重复概率，在客户端做一次去重就可以了，但是整体所花费的时间会比直接用 keys 指令长。

4. 使用过 Redis 做异步队列么，你是怎么用的？
---------------------------

一般使用 list 结构作为队列，rpush 生产消息，lpop 消费消息。当 lpop 没有消息的时候，要适当 sleep 一会再重试。
如果对方追问可不可以不用 sleep 呢？list 还有个指令叫 blpop，在没有消息的时候，它会阻塞住直到消息到来。
如果对方追问能不能生产一次消费多次呢？使用 pub/sub 主题订阅者模式，可以实现 1:N 的消息队列。
pub/sub 有什么缺点？在消费者下线的情况下，生产的消息会丢失，得使用专业的消息队列如 rabbitmq 等。
redis 如何实现延时队列？使用 sortedset，拿时间戳作为 score，消息内容作为 key 调用 zadd 来生产消息，消费者用 zrangebyscore 指令获取 N 秒之前的数据轮询进行处理。

5. 如果有大量的 key 需要设置同一时间过期，一般需要注意什么？
----------------------------------

如果大量的 key 过期时间设置的过于集中，到过期的那个时间点，redis 可能会出现短暂的卡顿现象。一般需要在时间上加一个随机值，使得过期时间分散一些。

6. Redis 如何做持久化的？
-----------------

bgsave 做镜像全量持久化，aof 做增量持久化。因为 bgsave 会耗费较长时间，不够实时，在停机的时候会导致大量丢失数据，所以需要 aof 来配合使用。在 redis 实例重启时，优先使用 aof 来恢复内存的状态，如果没有 aof 日志，就会使用 rdb 文件来恢复。
如果再问 aof 文件过大恢复时间过长怎么办？你告诉面试官，Redis 会定期做 aof 重写，压缩 aof 文件日志大小。如果面试官不够满意，再拿出杀手锏答案，Redis4.0 之后有了混合持久化的功能，将 bgsave 的全量和 aof 的增量做了融合处理，这样既保证了恢复的效率又兼顾了数据的安全性。这个功能甚至很多面试官都不知道，他们肯定会对你刮目相看。
如果对方追问那如果突然机器掉电会怎样？取决于 aof 日志 sync 属性的配置，如果不要求性能，在每条写指令时都 sync 一下磁盘，就不会丢失数据。但是在高性能的要求下每次都 sync 是不现实的，一般都使用定时 sync，比如 1s1 次，这个时候最多就会丢失 1s 的数据。

7. Pipeline 有什么好处，为什么要用 pipeline？
---------------------------------

可以将多次 IO 往返的时间缩减为一次，前提是 pipeline 执行的指令之间没有因果相关性。使用 redis-benchmark 进行压测的时候可以发现影响 redis 的 QPS 峰值的一个重要因素是 pipeline 批次指令的数目。

8. Redis 的同步机制了解么？
------------------

从从同步。第一次同步时，主节点做一次 bgsave，并同时将后续修改操作记录到内存 buffer，待完成后将 rdb 文件全量同步到复制节点，复制节点接受完成后将 rdb 镜像加载到内存。加载完成后，再通知主节点将期间修改的操作记录同步到复制节点进行重放就完成了同步过程。

9. 是否使用过 Redis 集群，集群的原理是什么？
---------------------------

Redis Sentinal 着眼于高可用，在 master 宕机时会自动将 slave 提升为 master，继续提供服务。
Redis Cluster 着眼于扩展性，在单个 redis 内存不足时，使用 Cluster 进行分片存储。

Redis 提升试题
==========

1、在项目中缓存是如何使用的？为什么要用缓存？缓存使用不当会造成什么后果？
-------------------------------------

#### 面试官心理分析

这个问题，互联网公司必问，要是一个人连缓存都不太清楚，那确实比较尴尬。

只要问到缓存，上来第一个问题，肯定是先问问你项目哪里用了缓存？为啥要用？不用行不行？如果用了以后可能会有什么不良的后果？

这就是看看你对缓存这个东西背后有没有思考，如果你就是傻乎乎的瞎用，没法给面试官一个合理的解答，那面试官对你印象肯定不太好，觉得你平时思考太少，就知道干活儿。

#### 面试题剖析

项目中缓存是如何使用的？

这个，需要结合自己项目的业务来。

为什么要用缓存？

用缓存，主要有两个用途：高性能、高并发。

#### **高性能**

假设这么个场景，你有个操作，一个请求过来，吭哧吭哧你各种乱七八糟操作 mysql，半天查出来一个结果，耗时 600ms。但是这个结果可能接下来几个小时都不会变了，或者变了也可以不用立即反馈给用户。那么此时咋办？

缓存啊，折腾 600ms 查出来的结果，扔缓存里，一个 key 对应一个 value，下次再有人查，别走 mysql 折腾 600ms 了，直接从缓存里，通过一个 key 查出来一个 value，2ms 搞定。性能提升 300 倍。

就是说对于一些需要复杂操作耗时查出来的结果，且确定后面不怎么变化，但是有很多读请求，那么直接将查询出来的结果放在缓存中，后面直接读缓存就好。

#### **高并发**

所以要是你有个系统，高峰期一秒钟过来的请求有 1 万，那一个 mysql 单机绝对会死掉。你这个时候就只能上缓存，把很多数据放缓存，别放 mysql。缓存功能简单，说白了就是 key-value 式操作，单机支撑的并发量轻松一秒几万十几万，支撑高并发 so easy。单机承载并发量是 mysql 单机的几十倍。

缓存是走内存的，内存天然就支撑高并发。

用了缓存之后会有什么不良后果？

常见的缓存问题有以下几个：

缓存与数据库双写不一致、缓存雪崩、缓存穿透、缓存并发竞争后面再详细说明。

2、redis 和 memcached 有什么区别？redis 的线程模型是什么？为什么 redis 单线程却能支撑高并发？
--------------------------------------------------------------

#### 面试官心理分析

这个是问 redis 的时候，最基本的问题吧，redis 最基本的一个内部原理和特点，就是 redis 实际上是个单线程工作模型，你要是这个都不知道，那后面玩儿 redis 的时候，出了问题岂不是什么都不知道？

还有可能面试官会问问你 redis 和 memcached 的区别，但是 memcached 是早些年各大互联网公司常用的缓存方案，但是现在近几年基本都是 redis，没什么公司用 memcached 了。

#### 面试题剖析

redis 和 memcached 有啥区别？

#### **redis 支持复杂的数据结构**

redis 相比 memcached 来说，拥有更多的数据结构，能支持更丰富的数据操作。如果需要缓存能够支持更复杂的结构和操作， redis 会是不错的选择。

#### **redis 原生支持集群模式**

在 redis3.x 版本中，便能支持 cluster 模式，而 memcached 没有原生的集群模式，需要依靠客户端来实现往集群中分片写入数据。

#### **性能对比**

由于 redis 只使用单核，而 memcached 可以使用多核，所以平均每一个核上 redis 在存储小数据时比 memcached 性能更高。而在 100k 以上的数据中，memcached 性能要高于 redis。虽然 redis 最近也在存储大数据的性能上进行优化，但是比起 memcached，还是稍有逊色。

redis 的线程模型
-----------

redis 内部使用文件事件处理器 file event handler，这个文件事件处理器是单线程的，所以 redis 才叫做单线程的模型。它采用 IO 多路复用机制同时监听多个 socket，将产生事件的 socket 压入内存队列中，事件分派器根据 socket 上的事件类型来选择对应的事件处理器进行处理。

文件事件处理器的结构包含 4 个部分：

*   多个 socket
*   IO 多路复用程序
*   文件事件分派器
*   事件处理器（连接应答处理器、命令请求处理器、命令回复处理器）

多个 socket 可能会并发产生不同的操作，每个操作对应不同的文件事件，但是 IO 多路复用程序会监听多个 socket，会将产生事件的 socket 放入队列中排队，事件分派器每次从队列中取出一个 socket，根据 socket 的事件类型交给对应的事件处理器进行处理。

来看客户端与 redis 的一次通信过程：

[![](https://img-blog.csdnimg.cn/20210505221627697.png)](https://img-blog.csdnimg.cn/20210505221627697.png)

要明白，通信是通过 socket 来完成的，不懂的同学可以先去看一看 socket 网络编程。

> *   客户端 socket01 向 redis 的 server socket 请求建立连接，此时 server socket 会产生一个 `AE_READABLE` 事件，IO 多路复用程序监听到 server socket 产生的事件后，将该事件压入队列中。文件事件分派器从队列中获取该事件，交给`连接应答处理器`。连接应答处理器会创建一个能与客户端通信的 socket01，并将该 socket01 的 `AE_READABLE` 事件与命令请求处理器关联。
> *   假设此时客户端发送了一个 `set key value` 请求，此时 redis 中的 socket01 会产生 `AE_READABLE` 事件，IO 多路复用程序将事件压入队列，此时事件分派器从队列中获取到该事件，由于前面 socket01 的 `AE_READABLE` 事件已经与命令请求处理器关联，因此事件分派器将事件交给命令请求处理器来处理。命令请求处理器读取 socket01 的 `key value` 并在自己内存中完成 `key value` 的设置。操作完成后，它会将 socket01 的 `AE_WRITABLE` 事件与令回复处理器关联。
> *   如果此时客户端准备好接收返回结果了，那么 redis 中的 socket01 会产生一个 `AE_WRITABLE` 事件，同样压入队列中，事件分派器找到相关联的命令回复处理器，由命令回复处理器对 socket01 输入本次操作的一个结果，比如`ok`，之后解除 socket01 的 `AE_WRITABLE` 事件与命令回复处理器的关联。

这样便完成了一次通信。

为啥 redis 单线程模型也能效率这么高？

*   纯内存操作
*   核心是基于非阻塞的 IO 多路复用机制
*   单线程反而避免了多线程的频繁上下文切换问题

3、redis 都有哪些数据类型？分别在哪些场景下使用比较合适？
--------------------------------

#### 面试官心理分析

除非是面试官感觉看你简历，是工作 3 年以内的比较初级的同学，可能对技术没有很深入的研究，面试官才会问这类问题。否则，在宝贵的面试时间里，面试官实在不想多问。

其实问这个问题，主要有两个原因：

*   看看你到底有没有全面的了解 redis 有哪些功能，一般怎么来用，啥场景用什么，就怕你别就会最简单的 KV 操作；
*   看看你在实际项目里都怎么玩儿过 redis。

要是你回答的不好，没说出几种数据类型，也没说什么场景，你完了，面试官对你印象肯定不好，觉得你平时就是做个简单的 set 和 get。

#### 面试题剖析

redis 主要有以下几种数据类型：

*   string
*   hash
*   list
*   set
*   sorted set

### string

这是最简单的类型，就是普通的 set 和 get，做简单的 KV 缓存。

```
set college szu
```

### hash

这个是类似 map 的一种结构，这个一般就是可以将结构化的数据，比如一个对象（前提是**这个对象没嵌套其他的对象**）给缓存在 redis 里，然后每次读写缓存的时候，可以就操作 hash 里的**某个字段**。

```
hset person name bingo
hset person age 20
hset person id 1
hget person name
person = {
    "name": "bingo",
    "age": 20,
    "id": 1
}
```

### list

list 是有序列表，这个可以玩儿出很多花样。

比如可以通过 list 存储一些列表型的数据结构，类似粉丝列表、文章的评论列表之类的东西。

比如可以通过 lrange 命令，读取某个闭区间内的元素，可以基于 list 实现分页查询，这个是很棒的一个功能，基于 redis 实现简单的高性能分页，可以做类似微博那种下拉不断分页的东西，性能高，就一页一页走。

```
# 0 开始位置，-1 结束位置，结束位置为-1 时，表示列表的最后一个位置，即查看所有。
lrange mylist 0 -1
```

比如可以搞个简单的消息队列，从 list 头怼进去，从 list 尾巴那里弄出来。

```
lpush mylist 1
lpush mylist 2
lpush mylist 3 4 5

# 1
rpop mylist
```

### set

set 是无序集合，自动去重。

直接基于 set 将系统里需要去重的数据扔进去，自动就给去重了，如果你需要对一些数据进行快速的全局去重，你当然也可以基于 jvm 内存里的 HashSet 进行去重，但是如果你的某个系统部署在多台机器上呢？得基于 redis 进行全局的 set 去重。

可以基于 set 玩儿交集、并集、差集的操作，比如交集吧，可以把两个人的粉丝列表整一个交集，看看俩人的共同好友是谁？对吧。

把两个大 V 的粉丝都放在两个 set 中，对两个 set 做交集。

```
#-------操作一个 set-------
# 添加元素
sadd mySet 1

# 查看全部元素
smembers mySet

# 判断是否包含某个值
sismember mySet 3

# 删除某个/些元素
srem mySet 1
srem mySet 2 4

# 查看元素个数
scard mySet

# 随机删除一个元素
spop mySet

#-------操作多个 set-------
# 将一个 set 的元素移动到另外一个 set
smove yourSet mySet 2

# 求两 set 的交集
sinter yourSet mySet

# 求两 set 的并集
sunion yourSet mySet

# 求在 yourSet 中而不在 mySet 中的元素
sdiff yourSet mySet
```

### sorted set

sorted set 是排序的 set，去重但可以排序，写进去的时候给一个分数，自动根据分数排序。

```
zadd board 85 zhangsan
zadd board 72 lisi
zadd board 96 wangwu
zadd board 63 zhaoliu

# 获取排名前三的用户（默认是升序，所以需要 rev 改为降序）
zrevrange board 0 3

# 获取某用户的排名
zrank board zhaoliu
```

4、redis 的过期策略都有哪些？内存淘汰机制都有哪些？手写一下 LRU 代码实现？
-------------------------------------------

#### 面试官心理分析

如果你连这个问题都不知道，上来就懵了，回答不出来，那线上你写代码的时候，想当然的认为写进 redis 的数据就一定会存在，后面导致系统各种 bug，谁来负责？

常见的有两个问题：

**（1）往 redis 写入的数据怎么没了？**

可能有同学会遇到，在生产环境的 redis 经常会丢掉一些数据，写进去了，过一会儿可能就没了。我的天，同学，你问这个问题就说明 redis 你就没用对啊。redis 是缓存，你给当存储了是吧？

啥叫缓存？用内存当缓存。内存是无限的吗，内存是很宝贵而且是有限的，磁盘是廉价而且是大量的。可能一台机器就几十个 G 的内存，但是可以有几个 T 的硬盘空间。redis 主要是基于内存来进行高性能、高并发的读写操作的。

那既然内存是有限的，比如 redis 就只能用 10G，你要是往里面写了 20G 的数据，会咋办？当然会干掉 10G 的数据，然后就保留 10G 的数据了。那干掉哪些数据？保留哪些数据？当然是干掉不常用的数据，保留常用的数据了。

**（2）数据明明过期了，怎么还占用着内存？**

这是由 redis 的过期策略来决定。

#### 面试题剖析

**redis 过期策略**

redis 过期策略是：定期删除 + 惰性删除。

所谓定期删除，指的是 redis 默认是每隔 100ms 就随机抽取一些设置了过期时间的 key，检查其是否过期，如果过期就删除。

假设 redis 里放了 10w 个 key，都设置了过期时间，你每隔几百毫秒，就检查 10w 个 key，那 redis 基本上就死了，cpu 负载会很高的，消耗在你的检查过期 key 上了。注意，这里可不是每隔 100ms 就遍历所有的设置过期时间的 key，那样就是一场性能上的灾难。实际上 redis 是每隔 100ms 随机抽取一些 key 来检查和删除的。

但是问题是，定期删除可能会导致很多过期 key 到了时间并没有被删除掉，那咋整呢？所以就是惰性删除了。这就是说，在你获取某个 key 的时候，redis 会检查一下，这个 key 如果设置了过期时间那么是否过期了？如果过期了此时就会删除，不会给你返回任何东西。

获取 key 的时候，如果此时 key 已经过期，就删除，不会返回任何东西。

答案是：走内存淘汰机制。

**内存淘汰机制**

redis 内存淘汰机制有以下几个：

*   noeviction: 当内存不足以容纳新写入数据时，新写入操作会报错，这个一般没人用吧，实在是太恶心了。
*   allkeys-lru：当内存不足以容纳新写入数据时，在键空间中，移除最近最少使用的 key（这个是最常用的）。
*   allkeys-random：当内存不足以容纳新写入数据时，在键空间中，随机移除某个 key，这个一般没人用吧，为啥要随机，肯定是把最近最少使用的 key 给干掉啊。
*   volatile-lru：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，移除最近最少使用的 key（这个一般不太合适）。
*   volatile-random：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，随机移除某个 key。
*   volatile-ttl：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，有更早过期时间的 key 优先移除。

**手写一个 LRU 算法**

你可以现场手写最原始的 LRU 算法，那个代码量太大了，似乎不太现实。

不求自己纯手工从底层开始打造出自己的 LRU，但是起码要知道如何利用已有的 JDK 数据结构实现一个 Java 版的 LRU。

```
package top.todu.leaning.base.example;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.HashMap;
import java.util.Map;

/** 暂不考虑并发安全 */
public class LRUCache {
  private final int size;
  private final Map<String, Object> cacheMap;
  private final Deque<String> queue;

  public LRUCache(int size) {
    this.size = size;
    cacheMap = new HashMap<>(size);
    queue = new ArrayDeque<>(size);
  }

  public static void main(String[] args) {
    LRUCache cache = new LRUCache(2);
    cache.put("a", 1);
    cache.put("a", 2);
    cache.put("a", 3);
    System.out.println(cache.cacheMap);
    System.out.println(cache.queue);
    cache.put("b", 1);
    cache.put("c", 1);
    System.out.println(cache.cacheMap);
    System.out.println(cache.queue);
  }

  public void put(String key, Object value) {
    if (!cacheMap.containsKey(key)) {
      if (queue.size() >= size) {
        String last = queue.removeLast();
        cacheMap.remove(last);
      }
      queue.addFirst(key);
    }

    cacheMap.put(key, value);
  }
}
```

结果：

>{a=3}
>[a]
>{b=1, c=1}
>[c, b]

5、如何保证 redis 的高并发和高可用？redis 的主从复制原理能介绍一下么？redis 的哨兵原理能介绍一下么？
------------------------------------------------------------

#### 面试官心理分析

其实问这个问题，主要是考考你，redis 单机能承载多高并发？如果单机扛不住如何扩容扛更多的并发？redis 会不会挂？既然 redis 会挂那怎么保证 redis 是高可用的？

其实针对的都是项目中你肯定要考虑的一些问题，如果你没考虑过，那确实你对生产系统中的问题思考太少。

#### 面试题剖析

如果你用 redis 缓存技术的话，肯定要考虑如何用 redis 来加多台机器，保证 redis 是高并发的，还有就是如何让 redis 保证自己不是挂掉以后就直接死掉了，即 redis 高可用。

由于此节内容较多，因此，会分为两个小节进行讲解。- redis 主从架构 - redis 基于哨兵实现高可用 redis 实现高并发主要依靠主从架构，一主多从，一般来说，很多项目其实就足够了，单主用来写入数据，单机几万 QPS，多从用来查询数据，多个从实例可以提供每秒 10w 的 QPS。

如果想要在实现高并发的同时，容纳大量的数据，那么就需要 redis 集群，使用 redis 集群之后，可以提供每秒几十万的读写并发。

redis 高可用，如果是做主从架构部署，那么加上哨兵就可以了，就可以实现，任何一个实例宕机，可以进行主备切换。

6、redis 的持久化有哪几种方式？不同的持久化机制都有什么 优缺点？持久化机制具体底层是如何实现的？
----------------------------------------------------

#### 面试官心理分析

redis 如果仅仅只是将数据缓存在内存里面，如果 redis 宕机了再重启，内存里的数据就全部都弄丢了啊。

你必须得用 redis 的持久化机制，将数据写入内存的同时，异步的慢慢的将数据写入磁盘文件里，进行持久化。

如果 redis 宕机重启，自动从磁盘上加载之前持久化的一些数据就可以了，也许会丢失少许数据，但是至少不会将所有数据都弄丢。

这个其实一样，针对的都是 redis 的生产环境可能遇到的一些问题，就是 redis 要是挂了再重启，内存里的数据不就全丢了？能不能重启的时候把数据给恢复了？

#### 面试题剖析

持久化主要是做灾难恢复、数据恢复，也可以归类到高可用的一个环节中去，比如你 redis 整个挂了，然后 redis 就不可用了，你要做的事情就是让 redis 变得可用，尽快变得可用。

重启 redis，尽快让它对外提供服务，如果没做数据备份，这时候 redis 启动了，也不可用啊，数据都没了。

很可能说，大量的请求过来，缓存全部无法命中，在 redis 里根本找不到数据，这个时候就死定了，出现缓存雪崩问题。所有请求没有在 redis 命中，就会去 mysql 数据库这种数据源头中去找，一下子 mysql 承接高并发，然后就挂了…

如果你把 redis 持久化做好，备份和恢复方案做到企业级的程度，那么即使你的 redis 故障了，也可以通过备份数据，快速恢复，一旦恢复立即对外提供服务。

redis 持久化的两种方式

*   RDB：RDB 持久化机制，是对 redis 中的数据执行周期性的持久化。
*   AOF：AOF 机制对每条写入命令作为日志，以 append-only 的模式写入一个日志文件中，在 redis 重启的时候，可以通过回放 AOF 日志中的写入指令来重新构建整个数据集。

通过 RDB 或 AOF，都可以将 redis 内存中的数据给持久化到磁盘上面来，然后可以将这些数据备份到别的地方去，比如说阿里云等云服务。

如果 redis 挂了，服务器上的内存和磁盘上的数据都丢了，可以从云服务上拷贝回来之前的数据，放到指定的目录中，然后重新启动 redis，redis 就会自动根据持久化数据文件中的数据，去恢复内存中的数据，继续对外提供服务。

如果同时使用 RDB 和 AOF 两种持久化机制，那么在 redis 重启的时候，会使用 AOF 来重新构建数据，因为 AOF 中的数据更加完整。

**RDB 优缺点**

*   RDB 会生成多个数据文件，每个数据文件都代表了某一个时刻中 redis 的数据，这种多个数据文件的方式，非常适合做冷备，可以将这种完整的数据文件发送到一些远程的安全存储上去，比如说 Amazon 的 S3 云服务上去，在国内可以是阿里云的 ODPS 分布式存储上，以预定好的备份策略来定期备份 redis 中的数据。
*   RDB 对 redis 对外提供的读写服务，影响非常小，可以让 redis 保持高性能，因为 redis 主进程只需要 fork 一个子进程，让子进程执行磁盘 IO 操作来进行 RDB 持久化即可。·
*   相对于 AOF 持久化机制来说，直接基于 RDB 数据文件来重启和恢复 redis 进程，更加快速。
*   如果想要在 redis 故障时，尽可能少的丢失数据，那么 RDB 没有 AOF 好。一般来说，RDB 数据快照文件，都是每隔 5 分钟，或者更长时间生成一次，这个时候就得接受一旦 redis 进程宕机，那么会丢失最近 5 分钟的数据。
*   RDB 每次在 fork 子进程来执行 RDB 快照数据文件生成的时候，如果数据文件特别大，可能会导致对客户端提供的服务暂停数毫秒，或者甚至数秒。

**AOF 优缺点**

*   AOF 可以更好的保护数据不丢失，一般 AOF 会每隔 1 秒，通过一个后台线程执行一次 fsync 操作，最多丢失 1 秒钟的数据。
*   AOF 日志文件以 append-only 模式写入，所以没有任何磁盘寻址的开销，写入性能非常高，而且文件不容易破损，即使文件尾部破损，也很容易修复。
*   AOF 日志文件即使过大的时候，出现后台重写操作，也不会影响客户端的读写。因为在 rewrite log 的时候，会对其中的指令进行压缩，创建出一份需要恢复数据的最小日志出来。在创建新日志文件的时候，老的日志文件还是照常写入。当新的 merge 后日志文件 ready 的时候，在交换新老日志文件即可。
*   AOF 日志文件的命令通过非常可读的方式进行记录，这个特性非常适合做灾难性的误删除的紧急恢复。比如某人不小心用 flushall 命令清空了所有数据，只要这个时候后台 rewrite 还没有发生，那么就可以立即拷贝 AOF 文件，将最后一条 flushall 命令给删了，然后再将该 AOF 文件放回去，就可以通过恢复机制，自动恢复所有数据。
*   对于同一份数据来说，AOF 日志文件通常比 RDB 数据快照文件更大。
*   AOF 开启后，支持的写 QPS 会比 RDB 支持的写 QPS 低，因为 AOF 一般会配置成每秒 fsync 一次日志文件，当然，每秒一次 fsync，性能也还是很高的。（如果实时写入，那么 QPS 会大降，redis 性 能会大大降低）
*   以前 AOF 发生过 bug，就是通过 AOF 记录的日志，进行数据恢复的时候，没有恢复一模一样的数据出来。所以说，类似 AOF 这种较为复杂的基于命令日志 / merge / 回放的方式，比基于 RDB 每次持久化一份完整的数据快照文件的方式，更加脆弱一些，容易有 bug。不过 AOF 就是为了避免 rewrite 过程导致的 bug，因此每次 rewrite 并不是基于旧的指令日志进行 merge 的，而是基于当时内存中的数据进行指令的重新构建，这样健壮性会好很多。

**RDB 和 AOF 到底该如何选择**

*   不要仅仅使用 RDB，因为那样会导致你丢失很多数据；
*   也不要仅仅使用 AOF，因为那样有两个问题：第一，你通过 AOF 做冷备，没有 RDB 做冷备来的恢复速度更快；第二，RDB 每次简单粗暴生成数据快照，更加健壮，可以避免 AOF 这种复杂的备份和恢复机制的 bug；

### 缓存一致性

1. 先删除缓存，再更新DB；读操作比较少，没啥问题，如果读操作比较高，容易不一致： DB尚未更新完成，同时来了一个读操作，会将旧的数据写到缓存，导致不一致。
1. 先更新DB，再删除缓存，如果没有事务保证，缓存删除失败，就会出现不一致性。 如果读并发高，在更新DB之后，删除缓存之前出现读操作还是会出现短暂的不一致性； 如果没有事务保证，可以设置缓存的ttl，达到最终一致性。
1. 如果需要完全一致性，那就不要使用缓存，（全局加锁，对于使用缓存的目的也基本失去了意义，并发降低），所以使用缓存，我们只能追求缓存最终一致性

11、redis 的并发竞争问题是什么？如何解决这个问题？了解 redis 事务的 CAS 方案吗？
--------------------------------------------------

#### 面试官心理分析

这个也是线上非常常见的一个问题，就是多客户端同时并发写一个 key，可能本来应该先到的数据后到了，导致数据版本错了；或者是多客户端同时获取一个 key，修改值之后再写回去，只要顺序错了，数据就错了。

而且 redis 自己就有天然解决这个问题的 CAS 类的乐观锁方案。

### Optimistic locking using check-and-set(乐观锁)

#### 乐观锁介绍：

watch 指令在 redis 事物中提供了 CAS 的行为。为了检测被 watch 的 keys 在是否有多个 clients 同时改变引起冲突，这些 keys 将会被监控。如果至少有一个被监控的 key 在执行 exec 命令前被修改，整个事物将会回滚，不执行任何动作，从而保证原子性操作，并且执行 exec 会得到 null 的回复。

#### 乐观锁工作机制：

watch 命令会监视给定的每一个 key，当 exec 时如果监视的任一个 key 自从调用 watch 后发生过变化，则整个事务会回滚，不执行任何动作。注意 watch 的 key 是对整个连接有效的，事务也一样。如果连接断开，监视和事务都会被自动清除。当然 exec，discard，unwatch 命令，及客户端连接关闭都会清除连接中的所有监视。还有，如果 watch 一个不稳定 (有生命周期) 的 key 并且此 key 自然过期，exec 仍然会执行事务队列的指令。
[![](https://img-blog.csdnimg.cn/20210421172344584.png)](https://img-blog.csdnimg.cn/20210421172344584.png)
[![](https://img-blog.csdnimg.cn/20210421172417360.png)](https://img-blog.csdnimg.cn/20210421172417360.png)

#### redis 的 Watch 机制是什么？

> Redis Watch 命令用于监视一个 (或多个) key，如果在事务执行之前这个 (或这些) key 被其他命令所改动，那么事务将被打断。注意使用 multi 开始事务，exec 提交事务。

> 语法， redis Watch 命令基本语法如下：
> WATCH key [key …]

验证：首先开启两个 redis 客户端，客户端 1 和客户端 2.

*   1、客户端 1 中，先 set 一个值

```
redis 127.0.0.1:6379> set number 10
OK
```

*   2、客户端 1 开启 Watch 此值。

```
redis 127.0.0.1:6379> watch number
OK
```

*   3、客户端 1 开启事务，修改此值

```
redis 127.0.0.1:6379> multi
OK
redis 127.0.0.1:6379> set number 100
QUEUED
redis 127.0.0.1:6379> get number
QUEUED
redis 127.0.0.1:6379>
```

注意此时先不要 `exec` 执行

*   4、客户端 2，去修改此值

```
redis 127.0.0.1:6379> set number 500
OK
```

*   5、客户端 1，执行 `exec` 执行

```
redis 127.0.0.1:6379> exec
(nil)
redis 127.0.0.1:6379> get number
"500"
```

发现为 nil, 执行未成功，获取的值为客户端 2 修改后的值。

11、生产环境中的 redis 是怎么部署的？
-----------------------

#### 面试官心理分析析

看看你了解不了解你们公司的 redis 生产集群的部署架构，如果你不了解，那么确实你就很失职了，你的 redis 是主从架构？集群架构？用了哪种集群方案？有没有做高可用保证？有没有开启持久化机制确保可以进行数据恢复？线上 redis 给几个 G 的内存？设置了哪些参数？压测后你们 redis 集群承载多少 QPS？

兄弟，这些你必须是门儿清的，否则你确实是没好好思考过。

#### 面试题剖析

redis cluster，10 台机器，5 台机器部署了 redis 主实例，另外 5 台机器部署了 redis 的从实例， 每个主实例挂了一个从实例，5 个节点对外提供读写服务，每个节点的读写高峰 qps 可能可以达到每秒 5 万，5 台机器最多是 25 万读写请求 / s。

机器是什么配置？32G 内存 + 8 核 CPU + 1T 磁盘，但是分配给 redis 进程的是 10g 内存，一般线上生产环境，redis 的内存尽量不要超过 10g，超过 10g 可能会有问题。

5 台机器对外提供读写，一共有 50g 内存。

因为每个主实例都挂了一个从实例，所以是高可用的，任何一个主实例宕机，都会自动故障迁移，redis 从实例会自动变成主实例继续提供读写服务。

你往内存里写的是什么数据？每条数据的大小是多少？商品数据，每条数据是 10kb。100 条数据是 1mb，10 万条数据是 1g。常驻内存的是 200 万条商品数据，占用内存是 20g，仅仅不到总内存的 50%。目前高峰期每秒就是 3500 左右的请求量。

其实大型的公司，会有基础架构的 team 负责缓存集群的运维。

Redis 高级试题
==========

12. 了解什么是 Redis 的雪崩、穿透和击穿？Redis 崩溃之后会怎么样？系统该如何应对这种情况？如何处理 Redis 的穿透？
--------------------------------------------------------------------

#### 缓存雪崩

对于系统 A，假设每天高峰期每秒 5000 个请求，本来缓存在高峰期可以扛住每秒 4000 个请求，但是缓存机器意外发生了全盘宕机。缓存挂了，此时 1 秒 5000 个请求全部落数据库，数据库必然扛不住，它会报一下警，然后就挂了。此时，如果没有采用什么特别的方案来处理这个故障，DBA 很着急，重启数据库，但是数据库立马又被新的流量给打死了。

[![](https://img-blog.csdnimg.cn/img_convert/75aa68675170524b8b8252d65c2e74e5.png)](https://img-blog.csdnimg.cn/img_convert/75aa68675170524b8b8252d65c2e74e5.png)

缓存雪崩的事前事中事后的解决方案如下：

*   事前：Redis 高可用，主从 + 哨兵，Redis cluster，避免全盘崩溃。
*   事中：本地 ehcache 缓存 + hystrix 限流 & 降级，避免 MySQL 被打死。
*   事后：Redis 持久化，一旦重启，自动从磁盘上加载数据，快速恢复缓存数据。

[![](https://img-blog.csdnimg.cn/img_convert/e3d4e1cbd03b5326eb8b3f15a7fc73bc.png)](https://img-blog.csdnimg.cn/img_convert/e3d4e1cbd03b5326eb8b3f15a7fc73bc.png)

用户发送一个请求，系统 A 收到请求后，先查本地 ehcache 缓存，如果没查到再查 Redis。如果 ehcache 和 Redis 都没有，再查数据库，将数据库中的结果，写入 ehcache 和 Redis 中。

限流组件，可以设置每秒的请求，有多少能通过组件，剩余的未通过的请求，怎么办？**走降级**！可以返回一些默认的值，或者友情提示，或者空值。

好处：

*   数据库绝对不会死，限流组件确保了每秒只有多少个请求能通过。
*   只要数据库不死，就是说，对用户来说，2/5 的请求都是可以被处理的。
*   只要有 2/5 的请求可以被处理，就意味着你的系统没死，对用户来说，可能就是点击几次刷不出来页面，但是多点几次，就可以刷出来了。

#### 缓存穿透

对于系统 A，假设一秒 5000 个请求，结果其中 4000 个请求是黑客发出的恶意攻击。

黑客发出的那 4000 个攻击，缓存中查不到，每次你去数据库里查，也查不到。

举个栗子。数据库 id 是从 1 开始的，结果黑客发过来的请求 id 全部都是负数。这样的话，缓存中不会有，请求每次都 “**绕过缓存**”，直接查询数据库。这种恶意攻击场景的缓存穿透就会直接把数据库给打死。

[![](https://img-blog.csdnimg.cn/img_convert/834fcb5588972d9be273189239bf0aba.png)](https://img-blog.csdnimg.cn/img_convert/834fcb5588972d9be273189239bf0aba.png)

解决方式很简单，每次系统 A 从数据库中只要没查到，就写一个空值到缓存里去，比如`set -999 UNKNOWN`。然后设置一个过期时间，这样的话，下次有相同的 key 来访问的时候，在缓存失效之前，都可以直接从缓存中取数据。

这种方式虽然是简单，在某些场景（如数据量大的博客）下不优雅，还可能会缓存过多的空值，更加优雅的方式就是：使用 bitmap 布隆过滤

#### 缓存击穿

缓存击穿，就是说某个 key 非常热点，访问非常频繁，处于集中式高并发访问的情况，当这个 key 在失效的瞬间，大量的请求就击穿了缓存，直接请求数据库，就像是在一道屏障上凿开了一个洞。

不同场景下的解决方式可如下：

*   若缓存的数据是基本不会发生更新的，则可尝试将该热点数据设置为永不过期。
*   若缓存的数据更新不频繁，且缓存刷新的整个流程耗时较少的情况下，则可以采用基于 Redis、zookeeper 等分布式中间件的分布式互斥锁，或者本地互斥锁以保证仅少量的请求能请求数据库并重新构建缓存，其余线程则在锁释放后能访问到新缓存。
*   若缓存的数据更新频繁或者在缓存刷新的流程耗时较长的情况下，可以利用定时线程在缓存过期前主动地重新构建缓存或者延后缓存的过期时间，以保证所有的请求能一直访问到对应的缓存。

#### 缓存击穿和 缓存穿透这两者的区别：

*   缓存击穿重点在 “击” 就是某个或者是几个热点 key 穿透了缓存层
*   缓存穿透重点在 “透”：大量的请求绕过了缓存层

### Redis 布隆过滤器防恶意流量击穿缓存

#### 什么是恶意流量穿透

假设我们的 Redis 里存有一组用户的注册 email，以 email 作为 Key 存在，同时它对应着 DB 里的 User 表的部分字段。

一般来说，一个合理的请求过来我们会先在 Redis 里判断这个用户是否是会员，因为从缓存里读数据返回快。如果这个会员在缓存中不存在那么我们会去 DB 中查询一下。

现在试想，有千万个不同 IP 的请求（不要以为没有，我们就在 2018 年和 2019 年碰到了，因为攻击的成本很低）带着 Redis 里根本不存在的 key 来访问你的网站，这时我们来设想一下：

1.  请求到达 Web 服务器；
2.  请求派发到应用层 -> 微服务层；
3.  请求去 Redis 捞数据，Redis 内不存在这个 Key；
4.  于是请求到达 DB 层，在 DB 建立 connection 后进行一次查询

千万乃至上亿的 DB 连接请求，先不说 Redis 是否撑的住 DB 也会被瞬间打爆。这就是 “Redis 穿透”，它会打爆你的缓存或者是连 DB 一起打爆进而引起一系列的 “雪崩效应”。

[![](https://img-blog.csdnimg.cn/20200112145306477.png)](https://img-blog.csdnimg.cn/20200112145306477.png)

#### 怎么防

那就是使用布隆过滤器，可以把所有的 user 表里的关键查询字段放于 Redis 的 bloom 过滤器内。有人会说，这不疯了，我有 4000 万会员？so what！

你把 4000 会员放在 Redis 里是比较夸张，有些网站有 8000 万、1 亿会员呢？

> 因此我没让你直接放在 Redis 里，而是放在布隆过滤器内！

布隆过滤器内不是直接把 key,value 这样放进去的，它存放的内容是这么一个 **bitmap** 中。

#### **bitmap**

所谓的 Bit-map 就是用一个 bit 位来标记某个元素对应的 Value，通过 Bit 为单位来存储数据，可以大大节省存储空间.
所以我们可以通过一个 int 型的整数的 32 比特位来存储 32 个 10 进制的数字，那么这样所带来的好处是内存占用少、
效率很高（不需要比较和位移）比如我们要存储 5(101)、3(11) 四个数字，那么我们申请 int 型的内存空间，会有 32
个比特位。这四个数字的二进制分别对应
从右往左开始数，比如第一个数字是 5，对应的二进制数据是 101, 那么从右往左数到第 5 位，把对应的二进制数据
存储到 32 个比特位上。

第一个 5 就是 00000000000000000000000000101000
输入 3 时候 00000000000000000000000000001100

#### Bloom Filte 介绍

**1. 含义**

(1). 布隆过滤器（Bloom Filter）是由 Howard Bloom 在 1970 年提出的一种比较巧妙的概率型数据结构，它实际上是由一个很长的二进制 (0 或 1) 向量和一系列随机映射函数组成。

(2). 布隆过滤器可以用于检索一个元素是否在一个集合中。它可以告诉你某种东西**一定不存在**或者**可能存在**。当布隆过滤器说，某种东西存在时，这种东西可能不存在；当布隆过滤器说，某种东西不存在时，那么这种东西一定不存在。

(3). 布隆过滤器 优点：A. 空间效率高，占用空间少 B. 查询时间短

　　　　　　　 缺点：A. 有一定的误判率 B. 元素不能删除

**2. 原理** 　

　当一个元素被加入集合时，通过 K 个散列函数将这个元素映射成一个位数组中的 K 个点（使用多个哈希函数对**元素 key (bloom 中不存 value)** 进行哈希，算出一个整数索引值，然后对位数组长度进行取模运算得到一个位置，每个无偏哈希函数都会得到一个不同的位置），把它们置为 1。

　检索时，我们只要看看这些点是不是都是 1 就（大约）知道集合中有没有它了：① 如果这些点有任何一个为 0（如下图的 e），则被检元素一定不在；如果都是 1（如下图的 d），并不能完全说明这个元素就一定存在其中，有可能这些位置为 1 是因为其他元素的存在，这就是布隆过滤器会出现误判的原因。

如下图：

[![](https://img2020.cnblogs.com/blog/1031302/202011/1031302-20201106204833000-564795432.png)](https://img2020.cnblogs.com/blog/1031302/202011/1031302-20201106204833000-564795432.png)

补充：

　　Bloom Filter 跟 **‘’ 单哈希函数 BitMap‘’** 不同之处在于：Bloom Filter 使用了 k 个哈希函数，每个字符串跟 k 个 bit 对应，从而降低了冲突的概率。

**3. 实现**

(1). Redis 的 bitmap

　基于 redis 的 bitmap 数据结构 的相关指令来执行。

(2). RedisBloom （推荐）

　布隆过滤器可以使用 Redis 中的位图 (bitmap) 操作实现，直到 Redis4.0 版本提供了插件功能，Redis 官方提供的布隆过滤器才正式登场，布隆过滤器作为一个插件加载到 Redis Server 中，官网推荐了一个 RedisBloom 作为 Redis 布隆过滤器的 Module。

　详细安装、指令操作参考：[https://github.com/RedisBloom/RedisBloom](https://github.com/RedisBloom/RedisBloom)

　文档地址：[https://oss.redislabs.com/redisbloom/](https://oss.redislabs.com/redisbloom/)

(3). PyreBloom

　pyreBloom 是 Python 中 Redis + BloomFilter 模块，是 c 语言实现。如果觉得 Redis module 的形式部署很麻烦或者线上环境 Redis 版本不是 4.0 及以上，则可以采用这个，但是它是在 hiredis 基础上，需要安装 hiredis，且不支持重连和重试。

(4). Lua 脚本实现

　详见：[https://github.com/erikdubbelboer/redis-lua-scaling-bloom-filter](https://github.com/erikdubbelboer/redis-lua-scaling-bloom-filter)

(5). guvua 包自带的布隆过滤器

#### Bloom Filte 应用场景

**1. 解决缓存穿透**

(1). 含义

　业务请求中数据缓存中没有，DB 中也没有，导致类似请求直接跨过缓存，反复在 DB 中查询，与此同时缓存也**不会**得到更新。（详见：[https://www.cnblogs.com/yaopengfei/p/13878124.html）](https://www.cnblogs.com/yaopengfei/p/13878124.html%EF%BC%89)

(2). 解决思路

　事先把存在的 key 都放到 redis 的 **Bloom Filter** 中，他的用途就是存在性检测，如果 BloomFilter 中不存在，那么数据一定不存在；如果 BloomFilter 中存在，实际数据也有可能会不存在。

**剖析：**** 布隆过滤器可能会误判，放过部分请求，当不影响整体，所以目前该方案是处理此类问题最佳方案。**

**2. 黑名单校验**

　识别垃圾邮件，只要发送者在黑名单中的，就识别为垃圾邮件。假设黑名单的数量是数以亿计的，存放起来就是非常耗费存储空间的，布隆过滤器则是一个较好的解决方案。把所有黑名单都放在布隆过滤器中，再收到邮件时，判断邮件地址是否在布隆过滤器中即可。

ps：

　　如果用哈希表，每存储一亿个 email 地址，就需要 1.6GB 的内存（用哈希表实现的具体办法是将每一个 email 地址对应成一个八字节的信息指纹，然后将这些信息指纹存入哈希表，由于哈希表的存储效率一般只有 50%，因此一个 email 地址需要占用十六个字节。一亿个地址大约要 1.6GB，即十六亿字节的内存）。因此存贮几十亿个邮件地址可能需要上百 GB 的内存。而 Bloom Filter 只需要哈希表 1/8 到 1/4 的大小就能解决同样的问题。

**3. Web 拦截器**

(1). 含义

　如果相同请求则拦截，防止重复被攻击。

(2). 解决思路

　用户第一次请求，将请求参数放入布隆过滤器中，当第二次请求时，先判断请求参数是否被布隆过滤器命中，从而提高缓存命中率。

#### 布隆过滤器其他场景

　　比如有如下几个需求：

　　1、原本有 10 亿个号码，现在又来了 10 万个号码，要快速准确判断这 10 万个号码是否在 10 亿个号码库中？

　　解决办法一：将 10 亿个号码存入数据库中，进行数据库查询，准确性有了，但是速度会比较慢。

　　解决办法二：将 10 亿号码放入内存中，比如 Redis 缓存中，这里我们算一下占用内存大小：10 亿 * 8 字节 = 8GB，通过内存查询，准确性和速度都有了，但是大约 8gb 的内存空间，挺浪费内存空间的。

　　2、接触过爬虫的，应该有这么一个需求，需要爬虫的网站千千万万，对于一个新的网站 url，我们如何判断这个 url 我们是否已经爬过了？

　　解决办法还是上面的两种，很显然，都不太好。

　　3、同理还有垃圾邮箱的过滤。

　　那么对于类似这种，大数据量集合，如何准确快速的判断某个数据是否在大数据量集合中，并且不占用内存，**布隆过滤器**应运而生了。

​ 4、假设我用 python 爬虫爬了 4 亿条 url，需要去重？

#### 给 Redis 安装 Bloom Filter

Redis 从 4.0 才开始支持 bloom filter，因此本例中我们使用的是 Redis5.4。

Redis 的 bloom filter 下载地址在这：[https://github.com/RedisLabsModules/redisbloom.git](https://github.com/RedisLabsModules/redisbloom.git)

```
git clone https://github.com/RedisLabsModules/redisbloom.git
cd redisbloom
make # 编译
```

让 Redis 启动时可以加载 bloom filter 有两种方式：

**手工加载式：**

```
redis-server --loadmodule ./redisbloom/rebloom.so
```

**每次启动自加载：**

编辑 Redis 的 redis.conf 文件，加入：

```
loadmodule /soft/redisbloom/redisbloom.so
```

Like this:

[![](https://img-blog.csdnimg.cn/20200112235518154.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpZmV0cmFnZWR5,size_16,color_FFFFFF,t_70)](https://img-blog.csdnimg.cn/20200112235518154.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpZmV0cmFnZWR5,size_16,color_FFFFFF,t_70)

#### 在 Redis 里使用 Bloom Filter

基本指令：

bf.reserve {key} {error_rate} {size}

```
127.0.0.1:6379> bf.reserve userid 0.01 100000
OK
```

上面这条命令就是：创建一个空的布隆过滤器，并设置一个期望的错误率和初始大小。{error_rate} 过滤器的错误率在 0-1 之间，如果要设置 0.1%，则应该是 0.001。**该数值越接近 0，内存消耗越大，对 cpu 利用率越高**。

bf.add {key} {item}

```
127.0.0.1:6379> bf.add userid '181920'
(integer) 1
```

上面这条命令就是：往过滤器中添加元素。如果 key 不存在，过滤器会自动创建。

bf.exists {key} {item}

```
127.0.0.1:6379> bf.exists userid '101310299'
(integer) 1
```

上面这条命令就是：判断指定 key 的 value 是否在 bloomfilter 里存在。存在：返回 1，不存在：返回 0。

#### 引入 Redis 布隆过滤器防止缓存穿透

缓存穿透（大量查询一个不存在的 key）定义：

> 缓存穿透，是指查询一个数据库中不一定存在的数据；

正常使用缓存查询数据的流程是，依据 key 去查询 value，数据查询先进行缓存查询，如果 key 不存在或者 key 已经过期，再对数据库进行查询，并把查询到的对象，放进缓存。如果数据库查询对象为空，则不放进缓存。

如果每次都查询一个不存在 value 的 key，由于缓存中没有数据，所以每次都会去查询数据库；当对 key 查询的并发请求量很大时，每次都访问 DB，很可能对 DB 造成影响；并且由于缓存不命中，每次都查询持久层，那么也失去了缓存的意义。

**缓存穿透**　解决方法

**第一种是缓存层缓存空值**

将数据库中的空值也缓存到缓存层中，这样查询该空值就不会再访问 DB，而是直接在缓存层访问就行。

但是这样有个弊端就是缓存太多空值占用了更多的空间，可以通过给缓存层空值设立一个较短的过期时间来解决，例如 60s。

**第二种是布隆过滤器**

将数据库中所有的查询条件，放入布隆过滤器中，

当一个查询请求过来时，先经过布隆过滤器进行查，如果判断请求查询值存在，则继续查；如果判断请求查询不存在，直接丢弃。

这里看 **Bloom Filter。**

我们先看看一般业务缓存流程：

[![](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAACbklEQVRoQ+2aMU4dMRCGZw6RC1CSSyQdLZJtKQ2REgoiRIpQkCYClCYpkgIESQFIpIlkW+IIcIC0gUNwiEFGz+hlmbG9b1nesvGW++zxfP7H4/H6IYzkwZFwQAUZmpJVkSeniFJKA8ASIi7MyfkrRPxjrT1JjZ8MLaXUDiJuzwngn2GJaNd7vyP5IoIYY94Q0fEQIKIPRGS8947zSQTRWh8CwLuBgZx479+2BTkHgBdDAgGAC+fcywoyIFWqInWN9BSONbTmFVp/AeA5o+rjKRJ2XwBYRsRXM4ZXgAg2LAPzOCDTJYQx5pSIVlrC3EI45y611osMTHuQUPUiYpiVooerg7TWRwDAlhSM0TuI+BsD0x4kGCuFSRVzSqkfiLiWmY17EALMbCAlMCmI6IwxZo+INgQYEYKBuW5da00PKikjhNNiiPGm01rrbwDwofGehQjjNcv1SZgddALhlJEgwgJFxDNr7acmjFLqCyJuTd6LEGFttpmkYC91Hrk3s1GZFERMmUT01Xv/sQljjPlMRMsxO6WULwnb2D8FEs4j680wScjO5f3vzrlNJszESWq2LYXJgTzjZm56MCHf3zVBxH1r7ftU1splxxKYHEgoUUpTo+grEf303rPH5hxENJqDKQEJtko2q9zGeeycWy3JhpKhWT8+NM/sufIhBwKI+Mta+7pkfxKMtd8Qtdbcx4dUQZcFCQ2I6DcAnLUpf6YMPxhIDDOuxC4C6djoQUE6+tKpewWZ1wlRkq0qUhXptKTlzv93aI3jWmE0Fz2TeujpX73F9TaKy9CeMk8vZusfBnqZ1g5GqyIdJq+XrqNR5AahKr9CCcxGSwAAAABJRU5ErkJggg==)](http://p1.pstatp.com/large/pgc-image/ad836008a83d44ea807ec3699f235311)

先查询缓存，缓存不命中再查询数据库。 然后将查询结果放在缓存中即使数据不存在，也需要创建一个缓存，用来防止穿库。这里需要区分一下数据是否存在。 如果数据不存在，缓存时间可以设置相对较短，防止因为主从同步等问题，导致问题被放大。

这个流程中存在薄弱的问题是，当用户量太大时，我们会缓存大量数据空数据，并且一旦来一波冷用户，会造成雪崩效应。 对于这种情况，我们产生第二个版本流程: redis 过滤冷用户缓存流程

[![](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAACbklEQVRoQ+2aMU4dMRCGZw6RC1CSSyQdLZJtKQ2REgoiRIpQkCYClCYpkgIESQFIpIlkW+IIcIC0gUNwiEFGz+hlmbG9b1nesvGW++zxfP7H4/H6IYzkwZFwQAUZmpJVkSeniFJKA8ASIi7MyfkrRPxjrT1JjZ8MLaXUDiJuzwngn2GJaNd7vyP5IoIYY94Q0fEQIKIPRGS8947zSQTRWh8CwLuBgZx479+2BTkHgBdDAgGAC+fcywoyIFWqInWN9BSONbTmFVp/AeA5o+rjKRJ2XwBYRsRXM4ZXgAg2LAPzOCDTJYQx5pSIVlrC3EI45y611osMTHuQUPUiYpiVooerg7TWRwDAlhSM0TuI+BsD0x4kGCuFSRVzSqkfiLiWmY17EALMbCAlMCmI6IwxZo+INgQYEYKBuW5da00PKikjhNNiiPGm01rrbwDwofGehQjjNcv1SZgddALhlJEgwgJFxDNr7acmjFLqCyJuTd6LEGFttpmkYC91Hrk3s1GZFERMmUT01Xv/sQljjPlMRMsxO6WULwnb2D8FEs4j680wScjO5f3vzrlNJszESWq2LYXJgTzjZm56MCHf3zVBxH1r7ftU1splxxKYHEgoUUpTo+grEf303rPH5hxENJqDKQEJtko2q9zGeeycWy3JhpKhWT8+NM/sufIhBwKI+Mta+7pkfxKMtd8Qtdbcx4dUQZcFCQ2I6DcAnLUpf6YMPxhIDDOuxC4C6djoQUE6+tKpewWZ1wlRkq0qUhXptKTlzv93aI3jWmE0Fz2TeujpX73F9TaKy9CeMk8vZusfBnqZ1g5GqyIdJq+XrqNR5AahKr9CCcxGSwAAAABJRU5ErkJggg==)](http://p1.pstatp.com/large/pgc-image/136c03c25a734629a3fce3520fa637f5)

我们将数据库里面，命中的用户放在 redis 的 set 类型中，设置不过期。 这样相当把 redis 当作数据库的索引，只要查询 redis，就可以知道是否数据存在。 redis 中不存在就可以直接返回结果。 如果存在就按照上面提到一般业务缓存流程处理。

聪明的你肯定会想到更多的问题：

1.  redis 本身可以做缓存，为什么不直接返回数据呢？
2.  如果数据量比较大，单个 set，会有性能问题？
3.  业务不重要，将全量数据放在 redis 中，占用服务器大量内存。投入产出不成比例？

问题 1 需要区分业务场景，结果数据少，我们是可以直接使用 redis 作为缓存，直接返回数据。 结果比较大就不太适合用 redis 存放了。比如 ugc 内容，一个评论里面可能存在上万字，业务字段多。

redis 使用有很多技巧。bigkey 危害比较大，无论是扩容或缩容带来的内存申请释放， 还是查询命令使用不当导致大量数据返回，都会影响 redis 的稳定。这里就不细谈原因及危害了。 解决 bigkey 方法很简单。我们可以使用 hash 函数来分桶，将数据分散到多个 key 中。 减少单个 key 的大小，同时不影响查询效率。

问题 3 是 redis 存储占用内存太大。因此我们需要减少内存使用。 重新思考一下引入 redis 的目的。 redis 像一个集合，整个业务就是验证请求的参数是否在集合中。

[![](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAACbklEQVRoQ+2aMU4dMRCGZw6RC1CSSyQdLZJtKQ2REgoiRIpQkCYClCYpkgIESQFIpIlkW+IIcIC0gUNwiEFGz+hlmbG9b1nesvGW++zxfP7H4/H6IYzkwZFwQAUZmpJVkSeniFJKA8ASIi7MyfkrRPxjrT1JjZ8MLaXUDiJuzwngn2GJaNd7vyP5IoIYY94Q0fEQIKIPRGS8947zSQTRWh8CwLuBgZx479+2BTkHgBdDAgGAC+fcywoyIFWqInWN9BSONbTmFVp/AeA5o+rjKRJ2XwBYRsRXM4ZXgAg2LAPzOCDTJYQx5pSIVlrC3EI45y611osMTHuQUPUiYpiVooerg7TWRwDAlhSM0TuI+BsD0x4kGCuFSRVzSqkfiLiWmY17EALMbCAlMCmI6IwxZo+INgQYEYKBuW5da00PKikjhNNiiPGm01rrbwDwofGehQjjNcv1SZgddALhlJEgwgJFxDNr7acmjFLqCyJuTd6LEGFttpmkYC91Hrk3s1GZFERMmUT01Xv/sQljjPlMRMsxO6WULwnb2D8FEs4j680wScjO5f3vzrlNJszESWq2LYXJgTzjZm56MCHf3zVBxH1r7ftU1splxxKYHEgoUUpTo+grEf303rPH5hxENJqDKQEJtko2q9zGeeycWy3JhpKhWT8+NM/sufIhBwKI+Mta+7pkfxKMtd8Qtdbcx4dUQZcFCQ2I6DcAnLUpf6YMPxhIDDOuxC4C6djoQUE6+tKpewWZ1wlRkq0qUhXptKTlzv93aI3jWmE0Fz2TeujpX73F9TaKy9CeMk8vZusfBnqZ1g5GqyIdJq+XrqNR5AahKr9CCcxGSwAAAABJRU5ErkJggg==)](http://p3.pstatp.com/large/pgc-image/790bf4b04c9a4a0581607613ffe1886c)

这个结构就像洗澡的时候用的双向阀门：左边热水，右边冷水。大部分的编程语言都内置了 filter。 拿 python 举例，filter 函数用于过滤序列， 过滤掉不符合条件的元素，返回由符合条件元素组成的列表。

12. Redis 的并发竞争问题是什么？如何解决这个问题？了解 Redis 事务的 CAS 方案吗？
---------------------------------------------------

简单的讲：就是**多客户端同时并发写**一个 key，可能本来应该先到的数据后到了，导致数据版本错了；或者是多客户端同时获取一个 key，修改值之后再写回去，只要顺序错了，数据就错了。

而且 Redis 自己就有天然解决这个问题的 CAS 类的乐观锁方案，使用版本号进行控制，cas 的思想这里就不详细说了。

13. Redis 集群模式的工作原理能说一下么？在集群模式下，Redis 的 key 是如何寻址的？分布式寻址都有哪些算法？了解一致性 hash 算法吗？
------------------------------------------------------------------------------

#### Redis cluster 介绍

*   自动将数据进行分片，每个 master 上放一部分数据
*   提供内置的高可用支持，部分 master 不可用时，还是可以继续工作的

在 Redis cluster 架构下，每个 Redis 要放开两个端口号，比如一个是 6379，另外一个就是 加 1w 的端口号，比如 16379。

16379 端口号是用来进行节点间通信的，也就是 cluster bus 的东西，cluster bus 的通信，用来进行故障检测、配置更新、故障转移授权。cluster bus 用了另外一种二进制的协议， `gossip` 协议，用于节点间进行高效的数据交换，占用更少的网络带宽和处理时间。

#### 集群节点间的内部通信机制

#### 基本通信原理

集群元数据的维护有两种方式：集中式、Gossip 协议。Redis cluster 节点间采用 gossip 协议进行通信。

**集中式**是将集群元数据（节点信息、故障等等）几种存储在某个节点上。集中式元数据集中存储的一个典型代表，就是大数据领域的`storm`。它是分布式的大数据实时计算引擎，是集中式的元数据存储的结构，底层基于 zookeeper（分布式协调的中间件）对所有元数据进行存储维护。

[![](https://img-blog.csdnimg.cn/img_convert/b6eb2ae3caaf46151a2a82dec17e2a20.png)](https://img-blog.csdnimg.cn/img_convert/b6eb2ae3caaf46151a2a82dec17e2a20.png)

Redis 维护集群元数据采用另一个方式， `gossip` 协议，所有节点都持有一份元数据，不同的节点如果出现了元数据的变更，就不断将元数据发送给其它的节点，让其它节点也进行元数据的变更。

[![](https://img-blog.csdnimg.cn/img_convert/809bad17847d818c833b8376934e9890.png)](https://img-blog.csdnimg.cn/img_convert/809bad17847d818c833b8376934e9890.png)

**集中式**的**好处**在于，元数据的读取和更新，时效性非常好，一旦元数据出现了变更，就立即更新到集中式的存储中，其它节点读取的时候就可以感知到；**不好**在于，所有的元数据的更新压力全部集中在一个地方，可能会导致元数据的存储有压力。

gossip 好处在于，元数据的更新比较分散，不是集中在一个地方，更新请求会陆陆续续打到所有节点上去更新，降低了压力；不好在于，元数据的更新有延时，可能导致集群中的一些操作会有一些滞后。

*   10000 端口：每个节点都有一个专门用于节点间通信的端口，就是自己提供服务的端口号 + 10000，比如 7001，那么用于节点间通信的就是 17001 端口。每个节点每隔一段时间都会往另外几个节点发送 `ping` 消息，同时其它几个节点接收到 `ping` 之后返回`pong`。
*   交换的信息：信息包括故障信息，节点的增加和删除，hash slot 信息等等。

#### gossip 协议

gossip 协议包含多种消息，包含 `ping` , `pong` , `meet` , `fail` 等等。

*   meet：某个节点发送 meet 给新加入的节点，让新节点加入集群中，然后新节点就会开始与其它节点进行通信。

```
Redis-trib.rb add-node
```

其实内部就是发送了一个 gossip meet 消息给新加入的节点，通知那个节点去加入我们的集群。

*   ping：每个节点都会频繁给其它节点发送 ping，其中包含自己的状态还有自己维护的集群元数据，互相通过 ping 交换元数据。
*   pong：返回 ping 和 meeet，包含自己的状态和其它信息，也用于信息广播和更新。
*   fail：某个节点判断另一个节点 fail 之后，就发送 fail 给其它节点，通知其它节点说，某个节点宕机啦。

#### ping 消息深入

ping 时要携带一些元数据，如果很频繁，可能会加重网络负担。

每个节点每秒会执行 10 次 ping，每次会选择 5 个最久没有通信的其它节点。当然如果发现某个节点通信延时达到了`cluster_node_timeout / 2`，那么立即发送 ping，避免数据交换延时过长，落后的时间太长了。比如说，两个节点之间都 10 分钟没有交换数据了，那么整个集群处于严重的元数据不一致的情况，就会有问题。所以 `cluster_node_timeout` 可以调节，如果调得比较大，那么会降低 ping 的频率。

每次 ping，会带上自己节点的信息，还有就是带上 1/10 其它节点的信息，发送出去，进行交换。至少包含 `3` 个其它节点的信息，最多包含 `总节点数减 2` 个其它节点的信息。

#### 分布式寻址算法

*   hash 算法（大量缓存重建）
*   一致性 hash 算法（自动缓存迁移）+ 虚拟节点（自动负载均衡）
*   Redis cluster 的 hash slot 算法

#### hash 算法

来了一个 key，首先计算 hash 值，然后对节点数取模。然后打在不同的 master 节点上。一旦某一个 master 节点宕机，所有请求过来，都会基于最新的剩余 master 节点数去取模，尝试去取数据。这会导致**大部分的请求过来，全部无法拿到有效的缓存**，导致大量的流量涌入数据库。

[![](https://img-blog.csdnimg.cn/img_convert/47a6d4348acdcdf6c88096b368fc75c9.png)](https://img-blog.csdnimg.cn/img_convert/47a6d4348acdcdf6c88096b368fc75c9.png)

#### 一致性 hash 算法

一致性 hash 算法将整个 hash 值空间组织成一个虚拟的圆环，整个空间按顺时针方向组织，下一步将各个 master 节点（使用服务器的 ip 或主机名）进行 hash。这样就能确定每个节点在其哈希环上的位置。

来了一个 key，首先计算 hash 值，并确定此数据在环上的位置，从此位置沿环**顺时针 “行走”**，遇到的第一个 master 节点就是 key 所在位置。

在一致性哈希算法中，如果一个节点挂了，受影响的数据仅仅是此节点到环空间前一个节点（沿着逆时针方向行走遇到的第一个节点）之间的数据，其它不受影响。增加一个节点也同理。

燃鹅，一致性哈希算法在节点太少时，容易因为节点分布不均匀而造成**缓存热点**的问题。为了解决这种热点问题，一致性 hash 算法引入了虚拟节点机制，即对每一个节点计算多个 hash，每个计算结果位置都放置一个虚拟节点。这样就实现了数据的均匀分布，负载均衡。

[![](https://img-blog.csdnimg.cn/img_convert/af760135a43f89e98d75304ff7e768a9.png)](https://img-blog.csdnimg.cn/img_convert/af760135a43f89e98d75304ff7e768a9.png)

#### Redis cluster 的 hash slot 算法

Redis cluster 有固定的 `16384` 个 hash slot，对每个 `key` 计算 `CRC16` 值，然后对 `16384` 取模，可以获取 key 对应的 hash slot。

Redis cluster 中每个 master 都会持有部分 slot，比如有 3 个 master，那么可能每个 master 持有 5000 多个 hash slot。hash slot 让 node 的增加和移除很简单，增加一个 master，就将其他 master 的 hash slot 移动部分过去，减少一个 master，就将它的 hash slot 移动到其他 master 上去。移动 hash slot 的成本是非常低的。客户端的 api，可以对指定的数据，让他们走同一个 hash slot，通过 `hash tag` 来实现。

任何一台机器宕机，另外两个节点，不影响的。因为 key 找的是 hash slot，不是机器。

[![](https://img-blog.csdnimg.cn/img_convert/8894a88dac85b7e0abfcfd069b2bcbd8.png)](https://img-blog.csdnimg.cn/img_convert/8894a88dac85b7e0abfcfd069b2bcbd8.png)

14 Redis cluster 的高可用与主备切换原理
----------------------------

Redis cluster 的高可用的原理，几乎跟哨兵是类似的。

#### 判断节点宕机

如果一个节点认为另外一个节点宕机，那么就是`pfail`，**主观宕机**。如果多个节点都认为另外一个节点宕机了，那么就是`fail`，**客观宕机**，跟哨兵的原理几乎一样，sdown，odown。

在 `cluster-node-timeout` 内，某个节点一直没有返回`pong`，那么就被认为`pfail`。

如果一个节点认为某个节点 `pfail` 了，那么会在 `gossip ping` 消息中， `ping` 给其他节点，如果**超过半数**的节点都认为 `pfail` 了，那么就会变成`fail`。

#### 从节点过滤

对宕机的 master node，从其所有的 slave node 中，选择一个切换成 master node。

检查每个 slave node 与 master node 断开连接的时间，如果超过了`cluster-node-timeout * cluster-slave-validity-factor`，那么就**没有资格**切换成`master`。

#### 从节点选举

每个从节点，都根据自己对 master 复制数据的 offset，来设置一个选举时间，offset 越大（复制数据越多）的从节点，选举时间越靠前，优先进行选举。

所有的 master node 开始 slave 选举投票，给要进行选举的 slave 进行投票，如果大部分 master node `（N/2 + 1）` 都投票给了某个从节点，那么选举通过，那个从节点可以切换成 master。

从节点执行主备切换，从节点切换为主节点。

#### 与哨兵比较

整个流程跟哨兵相比，非常类似，所以说，Redis cluster 功能强大，直接集成了 replication 和 sentinel 的功能。
