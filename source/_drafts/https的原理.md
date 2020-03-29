---
title: https的原理
abbrlink: d9f
category:
tags:
---

https 使用了TLS进行加密。

### SSL

Secure Sockets Layer 的缩写，中文翻译：安全套接字。

### TLS


### X.509

### 非对称加密

### 对称加密

###

tls and (ip.dst==112.48.154.240 or ip.src==112.48.154.240)


## 传输过程

### TLS握手

1. 客户端：发送 Client Hello 到服务器
1. 服务端： 发送 Server Hello 到客户端
1. 服务端： 发送 Certificate（证书）到客户端
1. 服务端： 发送 Server Key Exchange 服务端密钥交换额外数据（公钥）（跟选择的非对称加密算法有关，有的算法不需要发送额外的Server Key Exchagne 消息）
1.Server Hello Down ，换告诉客户端TLS握手结束
1. 客户端校验证书的有效性，比如签发机构，组织，有效期等，如果校验通过或者根本不进行有效性校验则继续进行下面步骤
1. 客户端发送 Client Key Exchange，Change Cipher Exchange，Encrypted Handshake Message，客户端密钥额外交换数据，交换加密算法，加密的握手信息（使用公钥加密）
1. Change Cipher Exchange，Encrypted Handshake Message，服务端发送确认交换加密算法和加密信息，确定后续使用的对称加密密钥。

上述说的是证书合法的情况，如果证书不合法，到 客户端接受到 Server Hello Down消息后，客户端就会像服务端发送一个 Alert 消息，告诉它证书不合法，接下来也就不会进行常规的密钥交换流程，整个流程就结束了。

需要注意的是，如果客户端不关心证书有效性的话，还是会继续正常的流程，还是会进行密钥交换并进行加密传输。
![证书不合法的情况](https://public-links.todu.top/1585366986.png?imageMogr2/thumbnail/!100p)


### 数据传输

握手结束后，客户端和服务端双方都确认了数据传输使用的对称加密算法和对应的密钥，然后之后所有的数据都是用这个加密算法进行加解密。

这里为什么不继续使用非对称加密（公私钥方式）对称加密，而是换成了对称加密呢？因为非对称加密运算量非常大，效率低，但是对称加密效率非常高，为了提高加解密效率，就是用了对称加密。


1. Client Hello

    包含的主要内容： 支持的最高TSL协议版本version，从低到高依次 SSLv2 SSLv3 TLSv1 TLSv1.1 TLSv1.2，当前基本不再使用低于 TLSv1 的版本;，随机数，支持的加密套件（列表）

1. Server Hello

    包含的主要内容： tls版本，随机数，确定要使用加密套件（1个）



