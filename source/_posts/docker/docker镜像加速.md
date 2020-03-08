---
title: docker镜像加速
abbrlink: d2f5
date: 2020-03-08 20:38:11
updateDate: 2020-03-08 20:38:11
category: docker
tags:
  - docker
  - 加速
  - 镜像
---

在国内网络环境下使用docker，pull 镜像的时候基本就是很慢，甚至直接网络连接超时，如果有国内镜像加速，那么pull的速度会非常快。

现在docker镜像加速效果比较好的有[阿里云](https://aliyun.com)和[华为云](https://huaweicloud.com)。我个人现在主要使用[华为云镜像](https://mirrors.huaweicloud.com/)加速，下面介绍华为云加速配置方式。

## 安装docker-ce

### CentOS 配置

1. 若已经安装过docker，需要先删掉，之后再安装依赖

    ```bash
    sudo yum remove docker docker-common docker-selinux docker-engine
    sudo yum install -y yum-utils device-mapper-persistent-data lvm2
    ```

1. 根据版本不同，下载repo文件。您使用的发行版

    ```bash
    wget -O /etc/yum.repos.d/docker-ce.repo https://mirrors.huaweicloud.com/docker-ce/linux/centos/docker-ce.repo
    ```

    替换软件仓库地址

    ```bash
    sudo sed -i 's+download.docker.com+mirrors.huaweicloud.com/docker-ce+' /etc/yum.repos.d/docker-ce.repo
    ```

1. 更新索引文件并安装

    ```bash
    sudo yum makecache fast
    sudo yum install docker-ce
    ```

### Ubuntu 配置

1. 若已经安装过docker，需要先删掉，之后再安装依赖

    ```bash
    sudo apt-get remove docker docker-engine docker.io
    sudo apt-get install apt-transport-https ca-certificates curl gnupg2 software-properties-common
    ```

1. 信任Docker的GPG公钥

    ```bash
    curl -fsSL https://mirrors.huaweicloud.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
    ```

    对于amd64架构的计算机，添加软件仓库:

    ```bash
    sudo add-apt-repository "deb [arch=amd64] https://mirrors.huaweicloud.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
    ```

    对于树莓派或其它Arm架构计算机，请运行

    ```bash
    echo "deb [arch=armhf] https://mirrors.huaweicloud.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

1. 更新索引文件并安装

```bash
sudo apt-get update
sudo apt-get install docker-ce
```

## 其他Linux发行版

参考华为[镜像官方说明](https://mirrors.huaweicloud.com/)，在列表中找到 docker-ce 然后点开查看说明，如下图
![docker-ce使用说明](https://public-links.todu.top/1583673330.png?imageMogr2/thumbnail/!100p)

## mac安装

mac 需要安装[docker for mac](https://hub.docker.com/editions/community/docker-ce-desktop-mac/)，[点击下载最新稳定版本](https://download.docker.com/mac/stable/Docker.dmg)

这里安装的是 `19.03.5` 版本。

## 配置加速

如果要使用加速则需要登录华为云平台后才能使用，如果没有账号需要先[注册](https://activity.huaweicloud.com/2020feb_promotion/invite.html?fromuser=be03913b5e2242da844857a208c62262&fromacct=16523daa-93f8-4268-965a-3a07d7f50100&needGalaxy=true)再登录。

登录成功后，刷新刚才的镜像列表页面，选择 `DockerHub官方镜像`获取专属自己的加速地址 如下图：

![docker加速镜像](https://public-links.todu.top/1583673856.png?imageMogr2/thumbnail/!100p)
![专属加速地址](https://public-links.todu.top/1583673948.png?imageMogr2/thumbnail/!100p)

### Linux 配置

执行页面上类似如下命令执行，注意需要拷贝网页上的命令并执行，下面只是示例。

```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<- 'EOF'
{
    "registry-mirrors": ["https://专属地址.mirror.swr.myhuaweicloud.com"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### mac配置

运行docker后，点击菜单栏docker图标，选择 `Preferences`，在 Docker Engine 配置项中添加加速配置，其中加速地址换成自己的，注意配置格式是json，配置错误会导致启动失败，所以最好先将原配置文件考出来备份。

```json
"registry-mirrors": [
    "https://专属加速地址.mirror.swr.myhuaweicloud.com"
  ]
```

配置如下图：

![mac docker加速配置](https://public-links.todu.top/1583671505.png?imageMogr2/thumbnail/!100p)

配置完成后点击 `Apply & Restart` 重启生效。

如果配置文件错误导致docker启动失败，可以手动编辑 `~/.docker/daemon.json` 进行恢复。

## 最后

配置完成后，就可以体验到蹭蹭蹭的加速效果了。
