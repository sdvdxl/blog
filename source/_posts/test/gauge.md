---
title: 自动测试工具Gauge
date: 2017-07-31 14:49:23
tags:
  - gauge
  - 测试
category: test
----------------

Gauge 是一款轻量级、跨平台自动化测试工具集。规则文件语法可以使用markdown语法编写。另外还可以使用你喜欢的语言来编写业务规则代码比如 `go`、 `java`、`ruby`等语言。Gauge还提供了输出插件，可以将执行结果导出为html或者xml，甚至flash，方便查看。

# 主要特点
- A rich markup based on markdown
- Simple, Flexible and Rich Syntax
- Business Language Tests : Supports the concept of executable documentation.
- Consistent Cross Platform/Language Support for writing test code. Currently supported languages.
- Open Source, so it could be shared freely and improved by others as well.
- A modular architecture with Plugins support.
- Extensible through Plugins and Hackable.
- Supports External Data Sources.
- Helps you create Maintainable and Understandable test suites.
IDE Support.

# 概念(术语)
## Specifications (spec)
一个spec就是一个spec文件，用来定义规则。
使用markdown的一级标题来定义该文件的说明 ，比如:
```
Specification name
==================
```
或者
```
# Specification name
```

## Scenarios
一个 scenarios 可以认为是一个组，定义在spec文件中。使用markdown的二级标题定义
```
Scenarios
----------------
```
或者
```
## Scenarios
```

## Tags 
可以给 Spec 文件和 Scenarios 打标签。

```
Specification name
================

tags: s1,test1


Scenarios
----------------

tags: s1, s2
```
## [Steps](https://docs.getgauge.io/syntax.html#steps)
这个就是具体的测试用例。使用markdown的 `*` 声明一个用例。

```
* Step Name
```
step 还支持参数，包括静态参数，动态参数，表格参数  ，特殊参数
## 静态参数
静态参数用双引号包裹:
```
* Check "product" exists
```
## 动态参数
动态参数使用尖括号包裹，跟下面的特殊参数结合使用

```
* Check <product> exists
```

## 表格参数
表格参数，就是一个二维表格所构成的数据。

```
* Step that takes a table
   | id  | name    |
   |-----|---------|
   | 123 | John    |
   | 456 | Mcclain |
```
** 注意，定义和表格之间没有空行，表格是直接跟在定义的语句下面 **

## 特殊参数
```
<prefix:value>
```
prefix 支持 `file` 和 `table`

```
// file 
* Check if <file:/work/content.txt> is visible

// table
* Check if the users exist <table:/Users/john/work/users.csv>
```


## 注释
没有特殊的语法，任何一般的文本都可以作为注释

## 其他
支持图片，连接


# 安装
首先要安装Gauge程序，可以从这里[https://getgauge.io/get-started.html](https://getgauge.io/get-started.html)下载对应平台的安装程序。Windows平台没有测试。
## Mac
mac用户可以使用brew安装，`brew update && brew install gauge`，也可以从上面的地址下载安装包，然后执行安装即可。

## Linux
linux用户需要下载linux程序，下载下来是个安装包，假设我们放到 `/opt/gauge/gauge-0.9.0-linux.x86_64.zip` (最好创建一个目录，因为解压后不会生成单独的目录)。切换到 `/opt/gauge` 执行 `unzip gauge-0.9.0-linux.x86_64.zip`，得到以下文件：
```

├── bin
│   ├── gauge
│   └── gauge_screenshot
└── install.sh
```

执行 `./install.sh` ，一直按回车，使用默认配置即可安装完成。

## 其他环境
如果要使用java，则需要安装jdk，如果要使用C#，则需要安装.net环境，总之，需要什么语言，就需要安装什么环境。

安装完成后，命令行输入`gauge`，如果出现 gauge 的帮助信息，则说明安装成功，否则请按照[官方文档](https://getgauge.io/get-started.html)进行安装。

# 进阶
