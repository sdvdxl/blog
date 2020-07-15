---
title: Spring Mail使用线程池异步发送卡住的问题
category: Java
tags:
  - java
  - spring
  - mail
  - 线程
abbrlink: '2083'
---

## 背景

线上有用户反馈收欧洲用户不到邮件。

## 排查过程

查看了近几日的日志，发现邮件相关的异常：

```java
Mail server connection failed; nested exception is javax.mail.MessagingException: Can't send command to SMTP host;
  nested exception is:
        java.net.SocketException: Connection closed by remote host. Failed messages: javax.mail.MessagingException: Can't send command to SMTP host;
  nested exception is:
        java.net.SocketException: Connection closed by remote host,e.class = class org.springframework.mail.MailSendException
```

```java
Mail server connection failed; nested exception is javax.mail.MessagingException: Could not connect to SMTP host: smtpdm.aliyun.com, port: 465;
  nested exception is:
        javax.net.ssl.SSLHandshakeException: Remote host closed connection during handshake. Failed messages: javax.mail.MessagingException: Could not connect to SMTP host: smtpdm.aliyun.com, port: 465;
  nested exception is:
        javax.net.ssl.SSLHandshakeException: Remote host closed connection during handshake,e.class = class org.springframework.mail.MailSendException
```

，并且日志中打印了消息队列推送过来的邮件内容（使用Rabbit作为消息队列将消息推送到发邮件的模块）。根据日志错误描述，是网络问题导致邮件发送失败。

然后我又在线上环境中通过忘记密码功能将重置密码的验证码发送到邮箱，结果等了半天也没收到，并且日志中也没有相关错误，但是还是会打印收到的邮件内容；重试了几次之后一直都是收不到。为了验证是不是阿里云邮箱服务问题（我们使用了阿里云的邮件服务），本地进行了测试，发现可以发出去并且可以顺利收到。

这就很奇怪了，为啥海外线上就发不出去了呢？查阅了业务上的源代码，有一个地方值得怀疑，就是发送的时候使用了线程池，并且为了防止发送失败，自己写了个重试策略，也就是这个线程池实际使用的是 `ScheduledThreadPoolExecutor`。

既然服务收到了邮件内容，但是没有发出去，那么是不是发送邮件的线程卡住了呢？于是使用 `jstack pid` 命令dump了线程快照，主要内容如下：

![thread-dump](https://public-links.todu.top/1589644656.png?imageMogr2/thumbnail/!100p)

其中，**标签1** 表明了当前确实发送邮件的线程池（需要自定义线程池的ThreadFactory，可以使用google-guava的ThreadFactoryBuilder 来方便的构建）；标签2说明已经走到了spring的mail方法，标签3说明确实是我们发送邮件的业务代码。
