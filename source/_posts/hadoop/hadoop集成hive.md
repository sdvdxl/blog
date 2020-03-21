---
title: hadoop集成hive
tags:
  - hive
  - hadoop
  - hdfs
category: 大数据
abbrlink: 61872
date: 2016-03-11 16:25:27
---
# 前置条件
hadoop和yarn已经配置好并且成功运行。

# 选择版本和下载
[下载hive](http://www.apache.org/dyn/closer.cgi/hive/)，有两个版本可以选择，hive1和hive2，hive2版本MR功能已经废弃，将来版本可能会直接去掉，如果要用hive的MR功能，那么请选择hive1相应版本，否则的话选哪个都可以进行测试。

# 配置hive
1. 下载完hive后，解压，然后`sudo vi /etc/profile`编辑文件，添加环境变量

```bash
HIVE_HOME=hive目录的绝路路径
PATH=$PATH:$HIVE_HOME

export HIVE_HOME PATH
```

2. 控制台输入`source /etc/profile`使环境变量生效。

3. 进入$HIVE_HOME/conf目录，拷贝`hive-default.xml.template`并重命名为`hive-site.xml`，编辑：（如果没有该参数，请自行添加，存在则改之）

```xml
<property>
    <name>hive.exec.scratchdir</name>
    <value>/user/hive/tmp</value> <!--hdfs目录 hive目录需要手动创建,并改为777权限-->
    <description>Scratch space for Hive jobs</description>
</property>

<property>
    <name>hive.metastore.warehouse.dir</name>
    <value>/user/hive/warehouse</value> <!--hdfs目录，需要手动创建-->
    <description>location of default database for the warehouse</description>
</property>

<property>
    <name>javax.jdo.option.ConnectionURL</name>
    <value>jdbc:mysql://localhost:3306/hive?createDatabaseIfNotExist=true</value> <!--修改为自己的mysql数据库-->
    <description>JDBC connect string for a JDBC metastore</description>
</property>
<property>
  <name>javax.jdo.option.ConnectionDriverName</name>
  <value>com.mysql.jdbc.Driver</value>
  <description>Driver class name for a JDBC metastore</description>
</property>
<property>
  <name>javax.jdo.option.ConnectionUserName</name>
  <value>root</value> <!--mysql用户名-->
  <description>username to use against metastore database</description>
</property>
<property>
  <name>javax.jdo.option.ConnectionPassword</name>
  <value>root</value> <!--mysql密码-->
  <description>password to use against metastore database</description>
</property>
```
4. 然后查找value是`${system`开头的，替换成具体的本地绝对路径，可以创建一个hive用户，放到hive用户目录下。

5. 拷贝`hive-env.sh.template`并重命名`hive-env.sh`，添加如下内容，注意环境变量值更换为自己的路径。

```bash
HADOOP_HOME=/home/hadoop/hadoop-2.7.1

export JAVA_HOME=/usr/local/share/jdk1.8.0_73
```

6. 将hive-site.xml拷贝一份到hadoop的配置目录。
7. 下载mysql驱动放到hive目录下的lib目录中。
8. 启动hive，控制台输入`hive`，如果正确则输出一段信息后进入hive交互模式。

# 问题
1. > Caused by: org.datanucleus.store.rdbms.connectionpool.DatastoreDriverNotFoundException: The specified datastore driver ("com.mysql.jdbc.Driver") was not found in the CLASSPATH. Please check your CLASSPATH specification, and the name of the driver.

    **解决方法**：下载mysql驱动放到hive的lib目录

2. > org.apache.hadoop.security.AccessControlException: Permission denied: user=hive, access=WRITE, inode="/tmp/hadoop-yarn/staging/hive/.staging":hadoop:supergroup:drwxr-xr-x

    **解决方法**：在hdfs中创建tmp目录，并改为777权限。

3. > Starting Job = job_1457683200911_0001, Tracking URL = http://10.10.1.110:8088/proxy/application_1457683200911_0001/
Kill Command = /home/hadoop/hadoop-2.7.1/bin/hadoop job  -kill job_1457683200911_0001
Hadoop job information for Stage-1: number of mappers: 0; number of reducers: 0，打开管理页面发现信息：waiting for AM container to be allocated, launched and register with RM.

    **解决方法**：yarn没有配置正确，没有启动`nodemanager`，启动命令`yarn-daemon.sh start nodemanager`，启动之后即可运行job。


### 参考资料
1. [hadoop2.2完全分布式集群+hive+mysql存储元数据配置](http://blog.csdn.net/jyf211314/article/details/34110721)排版有点乱，请将就看
2. [hive导入HDFS数据](http://blog.csdn.net/z363115269/article/details/39048589) 和上面毛病类似




