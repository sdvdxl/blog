---
title: 安装启动Hadoop集群
date: 2016-03-09 17:49:43
tags:
  - hadoop
  - hdfs
category: hadoop
---

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
