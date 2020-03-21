---
title: 安装启动Hadoop集群
tags:
  - hadoop
  - hdfs
category: 大数据
abbrlink: 58641
date: 2016-03-09 17:49:43
---

# 环境
UCloud云主机，2.6.32-431.11.15.el6.ucloud.x86_64
假设三台主机内网IP分别为master, 10.10.1.11和10.10.1.12，hostname分别为：10-10-1-10,10-10-1-11和10-10-1-12

# 配置JDK
本次搭建测试用的是jdk8，可以从[Oracle官网下载](http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html)对应的版本。
## 配置jdk
假设jdk解压后目录存放在`/usr/local/jdk8`，命令行输入`sudo vi /etc/profile`，添加一下内容：
```bash
JAVA_HOME=/usr/local/jdk8
CLASSPATH=.:$JAVA_HOME/jre/lib
PATH=$PATH:$JAVA_HOME/bin
```
然后输入`source /etc/profile`使环境变量生效，输入`java -version`有java版本信息输出说明配置成功，三台主机均这么配置。

# 配置hadoop用户
控制台输入`sudo useradd -m -U hadoop` 添加hadoop用户，然后输入`sudo passwd hadoop` 修改hadoop用户的密码，输入`su -l hadoop`，输入刚才设置的密码，切换到hadoop用户。

# 配置ssh
前提是已经切换到hadoop用户。
在每个主机控制台输入`ssh-keygen`回车，一直回车直到结束。最后在master主机上使用`ssh-copy-id`命令拷贝认证信息到本主机和其他两台主机，这样可以免密码登录。`ssh-copy-id hadoop@主机地址`，

# 配置网络
三台主机均在 `/etc/hosts` 添加以下映射:
```
master master
10.10.1.11 slave1
10.10.1.12 slave2
```

# 配置Hadoop
1. [下载Hadoop](http://hadoop.apache.org/releases.html)，这里下载的是`2.7.1`版本。控制输入`wget http://www.apache.org/dyn/closer.cgi/hadoop/common/hadoop-2.7.1/hadoop-2.7.1.tar.gz`，等待下载完成。
2. 控制台输入`tar -xvf hadoop-2.7.1.tar.gz`解压，会生成`hadoop-2.7.1`目录。
3. 输入`cd hadoop-2.7.1/etc/hadoop`进入配置文件目录。
4. 修改`hadoop-env.sh`中的`export JAVA_HOME=`，将等号后的内容改成上面配置的jdk绝对路径，在这里就是`/usr/local/jdk8`，修改完后应该是`export JAVA_HOME=/usr/local/jdk8`，保存退出。
5. 修改`core-site.xml`，配置config内容：

```xml
<configuration>
    <property>
        <name>hadoop.tmp.dir</name>
        <value>/home/hadoop/tmp</value>
        <description>Abase for other temporary directories.</description>
    </property>
    <property>
        <name>fs.defaultFS</name>
        <!-- 注意，这里改成自己本机的ip -->
        <value>hdfs://master:9000</value>
    </property>
    <property>
        <name>fs.default.namenode</name>
        <!-- 注意，这里改成自己本机的ip -->
        <value>hdfs://master:8082</value>
    </property>
    <property>
        <name>io.file.buffer.size</name>
        <value>4096</value>
    </property>
    <property>
        <name>hadoop.native.lib</name>
        <value>true</value>
        <description>Should native hadoop libraries, if present, be used.</description>
    </property>
</configuration>
```
6. 修改`hdfs-site.xml`，修改config内容为：

```xml
<configuration>
    <property>
        <name>dfs.nameservices</name>
        <value>cluster</value>
    </property>
    <property>
        <name>dfs.namenode.secondary.http-address</name>
        <!-- 注意修改为自己的ip -->
        <value>master:50090</value>
    </property>
    <property>
        <name>dfs.namenode.name.dir</name>
        <value>/home/hadoop/dfs/name</value>
    </property>
    <property>
        <name>dfs.datanode.data.dir</name>
        <value>/home/hadoop/dfs/data</value>
    </property>
    <property>
        <name>dfs.replication</name>
        <value>2</value>
    </property>
    <property>
        <name>dfs.webhdfs.enabled</name>
        <value>true</value>
    </property>
</configuration>
```
7. 修改`yarn-site.xml`

```xml
<configuration>

  <!-- 注意ip改为自己的 -->
    <property>
        <name>yarn.nodemanager.aux-services</name>
        <value>mapreduce_shuffle</value>
    </property>
    <property>
        <name>yarn.resourcemanager.address</name>
        <value>master:8032</value>
    </property>
    <property>
        <name>yarn.resourcemanager.scheduler.address</name>
        <value>master:8030</value>
    </property>
    <property>
        <name>yarn.resourcemanager.resource-tracker.address</name>
        <value>master:8031</value>
    </property>
    <property>
        <name>yarn.resourcemanager.admin.address</name>
        <value>master:8033</value>
    </property>
    <property>
        <name>yarn.resourcemanager.webapp.address</name>
        <value>master:8088</value>
    </property>
</configuration>
```
8. 修改`mapred-site.xml`

```xml
<configuration>
    <property>
        <name>mapreduce.framework.name</name>
        <value>yarn</value>
    </property>
    <property>
        <name>mapreduce.jobtracker.http.address</name>
        <value>master:50030</value>
    </property>
    <property>
        <name>mapreduce.jobhistory.address</name>
        <value>master:10020</value>
    </property>
    <property>
        <name>mapreduce.jobhistory.webapp.address</name>
        <value>master:19888</value>
    </property>
</configuration>
```
9. 修改`slaves`文件，添加其他两台ip

``` java
slave1
slave2
```

将hadoop目录覆盖到其余机器对应目录。
下面开始操作hadoop命令，如果遇到hadoop native错误，请查看文末`Hadoop Native 配置`部分。
10. 格式化文件系统
注意：这里的格式化文件系统并不是硬盘格式化，只是针对主服务器hdfs-site.xml的dfs.namenode.name.dir和dfs.datanode.data.dir目录做相应的清理工作。切换到Hadoop的home目录，执行`bin/hdfs namenode -format`。
11. 启动停止服务
启动`sbin/start-dfs.sh`，可以一次性启动master和slaves节点服务。`sbin/start-yarn.sh`启动yarn资源管理服务。要停止服务，用对应的`sbin/stop-dfs.sh`和`sbin/stop-dfs.sh`即可停止服务。
12. 单独启动一个datanode
增加节点或者重启节点，需要单独启动，则可使用以下命令:
`sbin/hadoop-daemon.sh start datanode`，启动nodeManager`sbin/yarn-daemon.sh start nodemanager`，当然也可以操作namenode`sbin/hadoop-daemon.sh start namenode ` `sbin/yarn-daemon.sh start resourcemanager`。
**注意**：原文中是`sbin/yarn-daemons.sh`和`sbin/hadoop-daemons.sh`，运行后发现并没有启动成功，去掉s后启动成功。



# Hadoop Native 配置
输入  `hadoop checknative` 检查Hadoop本地库版本和相关依赖信息：

```
16/03/10 12:17:56 DEBUG util.NativeCodeLoader: Trying to load the custom-built native-hadoop library...
16/03/10 12:17:56 DEBUG util.NativeCodeLoader: Failed to load native-hadoop with error: java.lang.UnsatisfiedLinkError: /home/hadoop/hadoop-2.6.3/lib/native/libhadoop.so.1.0.0: /lib64/libc.so.6: version `GLIBC_2.14' not found (required by /home/hadoop/hadoop-2.6.3/lib/native/libhadoop.so.1.0.0)
16/03/10 12:17:56 DEBUG util.NativeCodeLoader: java.library.path=/home/hadoop/hadoop-2.6.3/lib/native
16/03/10 12:17:56 WARN util.NativeCodeLoader: Unable to load native-hadoop library for your platform... using builtin-java classes where applicable
16/03/10 12:17:56 DEBUG util.Shell: setsid exited with exit code 0
Native library checking:
hadoop:  false
zlib:    false
snappy:  false
lz4:     false
bzip2:   false
openssl: false
16/03/10 12:17:56 INFO util.ExitUtil: Exiting with status 1
```

发现`/lib64/libc.so.6: version `GLIBC_2.14' not found`信息，说明该版本的Hadoop需要glibc_2.14版本。下面就安装所需的版本。
1. `mkdir glib_build && cd glib_build`
2. `wget http://ftp.gnu.org/gnu/glibc/glibc-2.14.tar.gz && wget http://ftp.gnu.org/gnu/glibc/glibc-linuxthreads-2.5.tar.bz2`
3. `tar zxf glibc-2.14.tar.gz && cd glibc-2.14 && tar jxf ../glibc-linuxthreads-2.5.tar.bz2`
4. `cd ../ && export CFLAGS="-g -O2" && ./glibc-2.14/configure --prefix=/usr --disable-profile --enable-add-ons --with-headers=/usr/include --with-binutils=/usr/bin`
5. `make`
6. `make install`
install最后会遇到错误信息：

```
CC="gcc -B/usr/bin/" /usr/bin/perl scripts/test-installation.pl /root/
/usr/bin/ld: cannot find -lnss_test1
collect2: ld returned 1 exit status
Execution of gcc -B/usr/bin/ failed!
The script has found some problems with your installation!
Please read the FAQ and the README file and check the following:
- Did you change the gcc specs file (necessary after upgrading from
  Linux libc5)?
- Are there any symbolic links of the form libXXX.so to old libraries?
  Links like libm.so -> libm.so.5 (where libm.so.5 is an old library) are wrong,
  libm.so should point to the newly installed glibc file - and there should be
  only one such link (check e.g. /lib and /usr/lib)
You should restart this script from your build directory after you've
fixed all problems!
Btw. the script doesn't work if you're installing GNU libc not as your
primary library!
make[1]: *** [install] Error 1
make[1]: Leaving directory `/root/glibc-2.14'
make: *** [install] Error 2
```

无需关注，检验是否成功
`ls -l /lib64/libc.so.6`
lrwxrwxrwx 1 root root 12 Mar 10 12:12 /lib64/libc.so.6 -> libc-2.14.so
出现了`/lib64/libc.so.6 -> libc-2.14.so`字样说明成功了。

安装openssl
`yum install openssl-static.x86_64`

# 如何修改主机名称
修改文件`/etc/sysconfig/network`
然后执行`/etc/rc.d/init.d/network restart`重启网络模块

# secondaryNameNode 配置
1. 修改masters文件（如果没有则自己创建），添加一个主机名称，用以作为secondaryNameNode。
2. 修改hdfs-site.xml的内容，删除`dfs.namenode.secondary.http-address`部分配置，添加新的配置（注意修改为自己的ip）：

```xml
<property>
	<name>dfs.http.address</name>
	<value>master:50070</value>
	<description>
		The address and the base port where the dfs namenode web ui will listen on.
		If the port is 0 then the server will start on a free port.
	</description>
</property>
<property>
	<name>dfs.namenode.secondary.http-address</name>
	<value>10.10.1.11</value>
</property>
```

###### 参考资料：
1. [Hadoop-2.5.2集群安装配置详解](http://blog.csdn.net/tang9140/article/details/42869531)
2. [基于hadoop2.2的namenode与SecondaryNameNode分开配置在不同的计算机](http://blog.csdn.net/zzu09huixu/article/details/36873669)
