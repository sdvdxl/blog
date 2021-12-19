---
title: java-web程序模拟time_wait
category: java
tags:
  - java
  - linux
  - tcp
  - close_wait
keywords:
  - java
  - linux
  - tcp
  - close_wait
abbrlink: b3f1
date: 2021-12-01 22:25:01
updateDate: 2021-12-01 22:25:01
cover: 'https://img2.baidu.com/it/u=3068497985,619330922&fm=253&fmt=auto&app=138&f=PNG?w=402&h=416'
---

# tcp close_wait 模拟

https://start.spring.io/ 创建一个基本的web工程

![image-20211201213924617](https://tva1.sinaimg.cn/large/008i3skNly1gwynwip0c9j31fq0u0juq.jpg)

程序只添加一个java代码：



```java
package top.todu.top.demo;

import java.util.concurrent.TimeUnit;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class Controller {
  @GetMapping("/")
  public void c() {
    // 无限等待，故意造成请求未结束
    try {
      TimeUnit.DAYS.sleep(24);
    } catch (InterruptedException e) {
      e.printStackTrace();
    }
  }
}

```

执行 ` ./mvnw clean package -Dmaven.test.skip=true` 打包，运行（扔到Linux上）

本地 curl一下，服务器上 执行 netstat -anpl|grep 8080

![image-20211201220540734](https://tva1.sinaimg.cn/large/008i3skNly1gwyonrhgd1j311u0u0789.jpg)

断开curl请求：

这个请求从ESTABLISHED 变成了TIME_WAIT

![image-20211201220914190](https://tva1.sinaimg.cn/large/008i3skNly1gwyorghyuoj316o04eaak.jpg)

再重复执行几次curl，断开请求，观察服务器连接状况，CLOSE_WAIT 的个数就是我们curl的次数

![image-20211201221342277](https://tva1.sinaimg.cn/large/008i3skNly1gwyow50g4vj315k07ygnf.jpg)

## 结论

web server 如果没有处理完当前请求，客户端也没有断开连接的话，该条连接（server 和 client端）状态都是 ESTABLISHED；

如果 server 没有处理完请求，client 主动断开，那么 server 端状态会变成 CLOSE_WAIT；因为 按照 TCP 断开连接四次挥手协议（如下图），client 发送 FIN，server 收到 回复 ACK，然后server进入 CLOSE_WAIT 状态，但是 server 没有处理完数据，没有进行 socket close，导致server端未发送FIN。

![img](https://tva1.sinaimg.cn/large/008i3skNly1gwyp48dq0uj31ek0u0jul.jpg)
