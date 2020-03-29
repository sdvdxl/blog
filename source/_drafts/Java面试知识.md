---
title: Java面试知识
category: Java
tags:
  - java
  - 面试
abbrlink: b8e1
---

01 Zookeeper 部分

CAP定理

一致性（Consistency） （等同于所有节点访问同一份最新的数据副本）
可用性（Availability）（每次请求都能获取到非错的响应——但是不保证获取的数据为最新数据）
分区容错性（Partition tolerance）（以实际效果而言，分区相当于对通信的时限要求。系统如果不能在时限内达成数据一致性，就意味着发生了分区的情况，必须就当前操作在C和A之间做出选择[3]。）

ZAB协议
Zookeeper基于ZAB（Zookeeper Atomic Broadcast），实现了主备模式下的系统架构，保持集群中各个副本之间的数据一致性。
ZAB协议定义了选举（election）、发现（discovery）、同步（sync）、广播(Broadcast)四个阶段。


leader选举算法和流程



02 Redis 部分

Redis的应用场景
Redis支持的数据类型（必考）
zset跳表的数据结构（必考）
Redis的数据过期策略（必考）
Redis的LRU过期策略的具体实现
如何解决Redis缓存雪崩，缓存穿透问题
Redis的持久化机制（必考）
Redis的管道pipeline
03 Mysql 部分

事务的基本要素
事务隔离级别
如何解决事务的并发问题(脏读，幻读)？
MVCC多版本并发控制？
binlog,redolog,undolog都是什么，起什么作用？
InnoDB的行锁/表锁？
myisam和innodb的区别，什么时候选择myisam？
为什么选择B+树作为索引结构？
索引B+树的叶子节点都可以存哪些东西？
查询在什么时候不走（预期中的）索引？
sql如何优化?
explain是如何解析sql的？
order by原理
04 JVM 部分

运行时数据区域（内存模型）
垃圾回收机制
垃圾回收算法
Minor GC和Full GC触发条件
GC中Stop the world
各垃圾回收器的特点及区别
双亲委派模型
JDBC和双亲委派模型关系
05 Java 基础部分

HashMap和ConcurrentHashMap区别
ConcurrentHashMap的数据结构
高并发HashMap的环是如何产生的？
volatile作用
Atomic类如何保证原子性（CAS操作）
synchronized和Lock的区别
为什么要使用线程池？
核心线程池ThreadPoolExecutor的参数
ThreadPoolExecutor的工作流程
如何控制线程池线程的优先级
线程之间如何通信
Boolean占几个字节
jdk1.8/jdk1.7都分别新增了哪些特性？
Exception和Error
06 Spring 部分

Spring的IOC/AOP的实现
动态代理的实现方式
Spring如何解决循环依赖（三级缓存）
Spring的后置处理器
Spring的@Transactional如何实现的？
Spring的事务传播级别
BeanFactory和ApplicationContext的联系和区别
07 其他部分

高并发系统的限流如何实现？
高并发秒杀系统的设计
负载均衡如何设计？


岗位要求
计算机相关专业本科及以上学历，三年以上软件开发经验，英语流利者优先； Java基础扎实，熟悉JVM原理、Java高级特性、Java网络编程、Java多线程编程； 精通Java主流开源框架，如Spring、Dubbo、Netty等，掌握底层原理和机制； 精通MySql、MyBatis、MyCat等数据库相关技术，对SQL性能优化有经验； 了解分布式系统原理：CAP、最终一致性、幂等操作、分布式事务等； 了解大型网络应用架构：MQ、缓存、负载均衡、集群技术、数据同步、高可用、可容灾等； 良好的团队合作精神和沟通能力，能主动寻求挑战、采取行动、达成目标； 持续学习，追求卓越，能为团队引入创新的技术和方案，用创新的思路解决问题

1. Serializable

Serializable 是一个标记性接口，就是该这种类型的接口没有任何需要实现的方法。

1. 数



二叉树
平衡二叉树：AVL
红黑树




启动类加载器（Bootstrap ClassLoader）
$JAVA_HOME/jre/lib目录下的jar文件，比如 rt.jar、tools.jar，或者-Xbootclasspath系统环境变量指定目录下的路径。

扩展类加载器（Extension ClassLoader）
这个类加载器由sun.misc.Launcher$ExtClassLoader来实现，负责加载$JAVA_HOME/jre/lib/ext目录中，或者java.ext.dirs系统变量指定路径中所有的类库，允许用户将具备通用性的类库可以放到ext目录下，扩展Java SE功能。在JDK 9之后，这种扩展机制被模块化带来的天然的扩展能力所取代。

应用类加载器（App/System ClassLoader），也称作为系统类加载器，这个类加载器由sun.misc.Launcher$AppClassLoader来实现。它负责加载用户应用类路径（ClassPath）上所有的类库，开发者同样可以直接在代码中使用这个类加载器。如果应用程序中没有自定义过自己的类加载器，一般情况下这个就是程序中默认的类加载器。


## Tree

二叉树，平衡二叉树，平衡树，2-3树，红黑树，BTree，B+Tree, B*Tree
