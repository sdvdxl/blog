---
title: Spark基础知识
tags:
  - spark
category: spark
abbrlink: 20157
date: 2016-03-09 13:59:50
---


# Spark基本概念
-  RDD——Resillient Distributed Dataset A Fault-Tolerant Abstraction for In-Memory Cluster Computing弹性分布式数据集。
- Operation——作用于RDD的各种操作分为transformation和action。
- Job——作业，一个JOB包含多个RDD及作用于相应RDD上的各种operation。
- Stage——一个作业分为多个阶段。
- Partition——数据分区， 一个RDD中的数据可以分成多个不同的区。
- DAG——Directed Acycle graph，有向无环图，反应RDD之间的依赖关系。
- Narrow dependency——窄依赖，子RDD依赖于父RDD中固定的data partition。
- Wide Dependency——宽依赖，子RDD对父RDD中的所有data partition都有依赖。
- Caching Managenment——缓存管理，对RDD的中间计算结果进行缓存管理以加快整 体的处理速度。
