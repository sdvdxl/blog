---
title: Java中的锁
abbrlink: c8e8
date: 2020-03-20 15:23:48
updateDate: 2020-03-20 15:23:48
category: Java
tags:
---

## 公平锁/非公平锁

### 公平锁

严格按照锁的申请顺序获得锁，但是因为严格按照顺序获得锁，唤醒申请锁的线程会带来一定性能损失，所以性能不如非公平锁高。

ReentrantLock 默认是非公平锁，但是可以通过参数配置成公平锁 `new ReentrantLock(true);`

### 非公平锁

对应于公平锁，就有非公平锁，是指获得锁的顺序并非严格按照申请所锁的顺序获得，后面申请的可能先获得锁。

非公平锁可以一定程度减少唤醒线程的等待时间，从而提高性能。

缺陷是，部分线程可能一直抢不到锁，那么就会一直等待（也就是线程饥饿）。

synchronized 是非公平锁。

ReentrantLock 默认是非公平锁。

### 可重入锁/非可重入锁

如果同一个线程可以重复获得已获得的锁，那么就说这个锁是可重入锁；反之，就是非可重入锁。

可重入锁可以一定程度避免死锁的发生。

乐观锁/悲观锁

自旋锁

分段锁

独享锁/共享锁

互斥锁/读写锁

偏向锁/轻量级锁/重量级锁


