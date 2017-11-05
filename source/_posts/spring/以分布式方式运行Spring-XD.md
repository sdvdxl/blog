---
title: 以分布式方式运行Spring-XD
tags:
  - spring
  - xd
  - spring-xd
  - spring xd
category: spring
abbrlink: 54199
date: 2016-03-09 14:18:33
---
*主要以官方文档说明进行配置*
# 简介
Spring XD分布式运行环境（DIRT）支持以分布式方式运行多个跨节点的任务。参见[Getting Started](http://docs.spring.io/spring-xd/docs/1.3.1.RELEASE/reference/html/#getting-started)获取以单个节点运行方式的信息。

XD的分布式运行架构主要由以下组件构成：
- Admin 主要用于管理Stream，Job的发布，用户操作，和提供运行时相关的状态，系统统计和分析的REST服务
- Container 托管发布的模块（Stream处理任务）和批量任务
- ZooKeeper 提供所有XD运行时的信息。追踪Container信息，如：modules，jobs发布情况，steam定义，发布状态等。
- Spring Batch Job Repository Database --这个要求要配置一个关系型数据库。XD包含了HSQLDB，但是不推荐用在生产环境中。XD支持任何JDBC型数据库。
- A Message Broker --用于数据传输。XD的数据传输模块设计成了插拔式。当前XD版本支持`Rabbit MQ`和`Redis `，这两个都支持stream和job过程产生的数据的传输，`Kafka`仅支持steam产生的数据传输。请注意：job使用Kafka作为数据传输是不稳定的。这个项目必须要配置一个作为数据传输的插件（推荐Redis）。
- Analytics Repository -- XD目前用Redis作为counters和gauges分析的存储方式。
XD的分布式运行环境概览如下：
![distributed-runtime-overview](/images/xd/distributed-runtime-overview.png)

# Server Configuration
默认查找`$XD_HOME/config/servers.yml`文件，作为配置文件。
但是可以使用`XD_CONFIG_LOCATION`环境变量改变配置文件夹，使用`XD_CONFIG_NAME`改变配置文件位置，如：
```
export XD_CONFIG_LOCATION=file:/xd/config/
export XD_CONFIG_NAME=region1-servers
```
* 注意，`XD_CONFIG_LOCATION`最后的`/`是必须的。

# Database Configuration
MySQL，PostGresql选其中一个配置即可，当然还有Oracle也可以配置，但是在这里没有列出，可以[]参考官方文档](http://docs.spring.io/spring-xd/docs/1.3.1.RELEASE/reference/html/#_database_configuration)
`xd-singlenode`模式是使用了一个嵌入式`HSQLDB `数据库，运行分布式模式的时候，可以使用独立的`HSQLDB`，但是仅仅推荐在学习和开发的时候使用它，正式环境最好使用其他比如`MySQL`，`Postgres`等等数据库。
* 注意：如果在stream模块中使用除了`Postgres`和`HSQLDB`数据库，那么需要把对应的驱动放到` $XD_HOME/lib`目录。
`servers.yml`文件中已经注释了一部分jdbc配置信息，可以按需更改。

## `MySQL`配置
```
spring:
    datasource:
    url: jdbc:mysql://yourDBhost:3306/yourDB
    username: yourUsername
    password: yourPassword
    driverClassName: com.mysql.jdbc.Driver
```
## Postgresql配置
```
spring:
  datasource:
    url: jdbc:postgresql://yourDBhost:5432/yourDB
    username: yourUsername
    password: yourPassword
    driverClassName: org.postgresql.Driver
```
## Redis
stream和job的数据传输需要用到（Rabbit MQ也可以）（当用作数据分析时候也需要），这里推荐统一使用redis作为配置，
```
spring:
  redis:
   port: 6379
   host: localhost
   pool:
     maxIdle: 8 # max idle connections in the pool
     minIdle: 0 # min idle connections in the pool
     maxActive: -1 # no limit to the number of active connections
     maxWait: 30000 # time limit to get a connection - only applies if maxActive is finite
```
### 安装redis
从[官网](redis.io)下载最新的redis，然后解压，进入redis根目录，
```
cd deps
make hiredis jemalloc linenoise lua
cd ..
make install
```
注意：依赖于gcc和make（Ubuntu系列如果没有安装 apt-get install gcc make）

# 开启页面登录密码保护
默认ui管理界面是没有安全配置的，不需密码即可访问，为了安全起见，我们可以设置登录用户和密码。
```
spring:
  profiles: admin                                                     (1)
security:
  basic:
    enabled: true                                                     (2)
    realm: SpringXD                                                  
  user:
    name: yourAdminUsername
    password: yourAdminPassword
    role: ADMIN, VIEW, CREATE       
```

注意：`spring.batch.initializer.enabled`默认是true，会使Spring Bath初始化表结构。

# 启动admin
admin只会有一个，用来协调container和管理相关stream，job等。
`xd/bin/xd-admin`

# 启动container
container可以启动多个，也就是组成多个节点，由此构成分布式运行环境。
```
xd/bin/xd-container
```

# 创建Stream
进入`$XD_HOME/shell/`,控制台输入`bin/xd-shell`，进去xd命令行交互模式，然后输入
```
stream create --name foo --definition 'time | log' --deploy
```
即可看见admin日志（控制台没关闭的话也可以看到）有时间信息输出。

下面这个stream是从kafka读取信息，然后传输到log里，所以需要配置kafka，请查阅kafka相关资料。

```
stream create --name kafkaDevice --definition 'kafka --outputType=text/plain --zkconnect=10.10.1.20:2181 --topic=kafka_test --offsetStorage=redis | log ' --deploy
```

# 末
最后贴出一份比较完整的配置：

```
security:
  basic:
    enabled: true
    realm: SpringXD
  user:
    name: hekr
    password: hekr
    # 必须配置角色才会生效
    role: ADMIN, VIEW, CREATE
spring:
    redis:
       port: 6379
       host: 127.0.0.1
       pool:
          maxIdle: 8 # max idle connections in the pool
          minIdle: 0 # min idle connections in the pool
          maxActive: -1 # no limit to the number of active connections
          maxWait: 30000 # time limit to get a connection - only applies if maxActive is finite
       #sentinel:
       #   master: mymaster
       #   nodes: 127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381
    batch:
        isolationLevel: ISOLATION_SERIALIZABLE
        # clobType:
        dbType: MYSQL
        maxVarcharLength: 2500
        tablePrefix: BATCH_
        validateTransactionState: true
        initializer:
          enabled: true
    datasource:
        url: jdbc:mysql://localhost:3306/xd
        username: root
        password: hekr
        driverClassName: com.mysql.jdbc.Driver
        testOnBorrow: true
        validationQuery: select 1
zk:
  namespace: xd
  client:
     connect: 10.10.1.20:2181
     sessionTimeout: 60000
     connectionTimeout: 30000
     initialRetryWait: 1000
     retryMaxAttempts: 3
xd:
    transport: redis
```

##注意：
1. 配置权限后，进入xd-shell会显示`server-unknown:>`， 需要配置一下admin server才能进入交互

```
admin config server --uri http://服务器地址:端口(默认9393) --username 用户名 --password 密码
```
2. `testOnBorrow`默认是`true`，如果配置为true或者没有配置，则需要配置正确的`validationQuery`，如果配置不正确则会有类似如下异常出现
```java
Command failed org.springframework.xd.rest.client.impl.SpringXDException: Could not get JDBC Connection; nested exception is java.sql.SQLException: Failed to validate a newly established connection.
```
