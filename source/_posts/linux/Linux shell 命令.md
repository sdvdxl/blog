layout: linux
title: Linux shell 命令
date: 2016-05-25 09:34:07
tags:
	- linux
	- shell
category: linux

---

## lsof
适用于ip4

```
lsof -Pnl +M -i4 | grep port
```

适用于ip6

```
lsof -Pnl +M -i6 | grep port
```

## awk
杀掉名字一样的java进程
```
jps |grep SparkSubmit | awk '{print "kill -9 " $1}' | sh
```
如果仅仅是打印命令，则后面的管道和sh不需要加，如下
```
jps |grep SparkSubmit | awk '{print "kill -9 " $1}'
```
## find
删除找到的符合条件的文件（或者目录）
```
find . -iname target -exec rm -rf {} \;
```