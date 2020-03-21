---
title: linux下根据端口号查询对应进程
tags:
  - linux
  - 端口
  - 进程
category: Linux
abbrlink: 10727
date: 2016-03-15 18:22:37
---

适用于ip4

```
lsof -Pnl +M -i4 | grep port
```

适用于ip6

```
lsof -Pnl +M -i6 | grep port
```
