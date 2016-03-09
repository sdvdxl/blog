---
title: Hexo+Git+Oschina+Golang+Tenxcloud打造博客
date: 2016-03-09 23:17:00
tags:
---
# 简介
*Hexo* 是一个快速、简洁且高效的博客框架。Hexo 使用 Markdown（或其他渲染引擎）解析文章，在几秒内，即可利用靓丽的主题生成静态网页。具体使用方法[参见这里](https://hexo.io/zh-cn/docs/)

*Git* 介绍和使用[参见这里](http://gitref.org/zh/index.html)

*OSCina* 开源信息平台，这里指的是[其下的Git托管平台](http://git.oschina.net/)

*Golang* 谷歌开发的一款跨平台的语言，[官方地址](golang.org)在国内无法打开，[golangtc](http://golangtc.com/)是一个Golang学习网站，可以自行查阅。

*TenxCloud* 也就是[时速云](https://www.tenxcloud.com/)，是国内最早的容器云平台之一(Container as a service)，提供丰富的容器化应用，镜像构建与发布，弹性可伸缩的容器服务，以及灵活、高性能的容器主机管理。容器化应用包括但不限于云主机，云数据库，大数据，Web应用等。

# 准备         
1. 安装hexo
2. 安装Golang，并配置Golang环境
3. 安装Git，并配置相关环境变量

# 创建Git仓库
打开[开源中国Git托管平台](http://git.oschina.net/)，(注册后)登录，点击右上角`+`号，新建项目，输入项目名，描述，如果不想公开的话，可以选择私有，其余默认即可，点击创建。然后克隆到本地。命令行切换到刚才克隆的项目根目录，输入`hexo init`，hexo博客初始化完成。输入`hexo generate`可以渲染页面，生成静态页面，默认是在public文件夹。hexo默认初始化忽略了public文件夹，我们需要修改`.gitignore`文件，删除public的记录，这样保证可以同步到git仓库中。通过git提交文件到远程仓库。

# 创建容器
这里之所以选择时速云，是因为一开始接触这类最早的就是这个平台，所以使用的还算熟悉。下面我们就在上面创建一个容器。
1. 登录之后，[选择这个镜像](https://hub.tenxcloud.com/repos/sdvdxl/golang)，此镜像集成了Git，Golang，SSH服务，点击右侧部署镜像按钮。如下图配置：
![配置](/images/other/2016-03-09_2351.png)，点击创建按钮，等待片刻即可创建成功。
2. [返回容器服务](https://console.tenxcloud.com/containers?0)，切换到北京2区，可以看到我们刚才创建的容器服务，点击如图所示的图标，进入控制台：
[进入控制台](/images/other/2016-03-09_2355.png)。
3. 进入控制台后，使用git命令`git clone 之前创建的git仓库地址`，克隆完后，进入项目目录，输入`go build server.go`，然后输入`./server &`运行服务端。
4. 现在打开容器服务视图，找到我们创建的容器，点击右侧的查看所有服务地址，点击协议为`HTTP`的那个服务地址，在打开的页面中即可看到我们的博客内容。
5. 点击绑定域名，绑定80端口域名，我们就可以通过自己的域名访问了。
注意：不要选择杭州区的服务，因为采用的是阿里云服务，所以会导致没有备案的域名没法打开。
