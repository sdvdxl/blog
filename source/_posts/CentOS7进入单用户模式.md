---
title: CentOS7 进入单用户模式
tags:
  - linux
  - centos7
  - 单用户
abbrlink: 99f7
date: 2019-06-02 12:39:46
updateDate: 2019-06-02 12:39:46
category: linux
---

## CentOS7 进入单用户模式

在出现启动项选择菜单的时候，移动到要启动的项目上，按 `e`，如图所示
![进入启动项修改信息](https://public-links.todu.top/1559449263.png?imageMogr2/thumbnail/!100p)

进入启动项配置页面，按照下图进行修改

**注意** 一定要将 `ro` 改为 `rw`，否则没有写权限；另外一定要删除 `rhgb quiet`，否则无法正常启动并初始化
![修改启动信息](https://public-links.todu.top/1559450983.png?imageMogr2/thumbnail/!100p)

等待初始化信息完成后，就可以操作了
![等待初始化信息完成后，就可以操作了](https://public-links.todu.top/1559451110.png?imageMogr2/thumbnail/!100p)

操作完成后，执行 `exec /sbin/init` 进入普通用户模式

完整示例
![完整示例](https://public-links.todu.top/2019-06-02%2012.33.01.gif?imageMogr2/thumbnail/!10p)
