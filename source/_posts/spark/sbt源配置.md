---
title: sbt源配置
p: spark/sbt源配置
tags:
  - sbt
  - scala
category: spark
abbrlink: 51898
date: 2016-06-20 11:21:13
---

# sbt 默认源下载有点慢，我们可以调教它，让它从我们自己配置的源下载。

## 配置源
在 `.sbt` （默认是在用户名下）文件夹中创建 `repositories` 文件，然后添加如下内容:     
```
[repositories]
  local
  my: http://o8r69qphn.bkt.clouddn.com/
  Nexus osc: http://maven.oschina.net/content/groups/public/
  Nexus osc thirdparty: http://maven.oschina.net/content/repositories/thirdparty/
  central: http://central.maven.org/maven2/



[ivy]
  local
  my: http://o8r69qphn.bkt.clouddn.com/
  Nexus osc: http://maven.oschina.net/content/groups/public/
  Nexus osc thirdparty: http://maven.oschina.net/content/repositories/thirdparty/
  central: http://central.maven.org/maven2/
  TypeSafe: https://oss.sonatype.org/content/repositories/releases/
  #proxy库
  typesafe-ivy-releases: http://dl.bintray.com/typesafe/ivy-releases/
  typesafe-maven-releases: http://dl.bintray.com/typesafe/maven-releases/
  typesafe-sbt-plugin-releases: http://dl.bintray.com/sbt/sbt-plugin-releases/

  #group库
  ivy-releases : typesafe-ivy-releases,typesafe-sbt-plugin-releases
```

上面 `repositories` 是说加载maven镜像中的库文件从这个标签下的路径中找， `local` 代表从本地中找，默认是 `.M2` 中，下面的都是自定义源，名字随便取。
下面 `ivy` 是加载ivy库的。

## 修改加载配置项
单纯修改上面的源还不足以让sbt加载我们的源。打开 sbt 软件安装位置下的 `conf/sbtopts` 文件，在其中添加：    
```
-Dsbt.override.build.repos=true
```
然后就可以生效了。

** 没有深入研究，如果有错还请指出 **
