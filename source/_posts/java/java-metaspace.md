---
title: java-metaspace
category: Java
tags:
  - java
  - jvm
  - metaspace
abbrlink: ac66
date: 2021-01-07 11:20:27
updateDate: 2021-01-07 11:20:27
---

## 介绍

Java 类在 Java Hotspot VM 中有一个内部表示，称为类元数据，里面存放了class信息。在 Java Hotspot VM 之前的版本中，类元数据是在所谓的`永久代`生成中分配的。在 JDK 8中，`永久代` 被删除，元数据在本地内存（直接内存）中分配。默认情况下，可用于元数据的本机内存量是无限的（实际最大是本地最大可用内存）。使用 MaxMetaspaceSize 选项对用于类元数据的本机内存量设置上限。

为了避免为元数据空间太小导致的垃圾收集（FullGC），我们应该给 MetaspaceSize 选项设置一个比较高的值。为应用程序分配的类元数据量与应用程序相关，不存在用于选择 MetaspaceSize 的通用准则。MetaspaceSize 的默认大小依赖于平台，范围从12 MB 到大约20 MB。

当卸载相应的 Java 类时，类元数据被释放。由于垃圾收集，Java 类被卸载，为了卸载类和释放类元数据，就会进行FullGC。当为类元数据提交的空间达到某个级别(高水平标记)时，将导致垃圾收集。在垃圾收集之后，可以根据从类元数据中释放的空间量提高或降低高水平标记。将提高高水位线，以免过早引发另一次垃圾收集。最高水平标记最初设置为命令行选项 MetaspaceSize 的值。它根据选项 MaxMetaspaceFreeRatio 和 MinMetaspaceFreeRatio 来升高或降低。如果类元数据可用的提交空间占类元数据总提交空间的百分比大于 MaxMetaspaceFreeRatio，则高水位标记将降低。如果它小于 MinMetaspaceFreeRatio，那么高水位线将被提高。


一般情况下这个值（12-20M是足够的）。如果应用中有用到实时加载类的情况（热加载class或者调用动态脚本比如JavaScript）会增加metaspace的使用量，容量可能会不够用，会造成oom，这时候应该增大metaspace和maxmetaspace（比如200M，如果还比较频繁出现oom，则需要再次提高这个值）

## 空间不足产生的问题

-XX:CompressedClassSpaceSize 不足的情况下会：

`java.lang.OutOfMemoryError: Compressed class space`

-XX:MaxMetaspaceSize 不足会：

`java.lang.OutOfMemoryError: Metaspace`


CompressedClassSpaceSize 默认1G，和 metaspace 内存是分配在2块本地内存上的（但是都称之为元数据metaspace）

所以：元数据大小=MaxMetaspace + CompressedClassSpaceSize

metaspace 不够用，会在gc日志（如果开启了的话）打印

`[Full GC (Metadata GC Threshold) [PSYoungGen: 96K->0K(693760K)] [ParOldGen: 187487K->187486K(1398272K)] 187583K->187486K(2092032K), [Metaspace: 144184K->144184K(1226752K)], 0.2667937 secs] [Times: user=0.59 sys=0.00, real=0.27 secs]`

## 其他说明

Metaspace used 2425K, capacity 4498K, committed 4864K, reserved 1056768K
class space used 262K, capacity 386K, committed 512K, reserved 1048576K

![](https://public-links.todu.top/1609988894.png?imageMogr2/thumbnail/!100p)


- reserved reserved 是指，操作系统已经为该进程“保留”的。所谓的保留，更加接近一种记账的概念，就是操作系统承诺说一大块连续的内存已经是你这个进程的了。
- committed 进程committed的时候。当进程真的要用这个连续地址空间的时候，操作系统才会分配真正的内存。所以，这也就是意味着，这个过程会失败。
- capacity jvm 分配的大小
- used jvm 真正使用的
- class space 存放class的内存
- metaspace 包含了class和静态变量的的内存

**重点** metaspace 是用来控制第一次因元空间触发的FullGC的，如果有动态脚本执行，要设置的大一点，比如200M；MaxMetaspaceSize 是用来控制最大类元空间的，要设置的比 metaspace 大。CompressedClassSpaceSize 默认是1G，这个值只是一个上限值，不是真正一次性占用的，只有真正使用了才占用实际内存。

## 应用

一般设置metaspace和MaxMetaspaceSize就够了



**参考资料**

[Oracle 官方 Metaspace 介绍](https://docs.oracle.com/javase/8/docs/technotes/guides/vm/gctuning/considerations.html#sthref66)

https://www.jianshu.com/p/cd34d6f3b5b4

https://www.cnblogs.com/williamjie/p/9558094.html
