---
title: archlinux安装
date: 2016-03-09 14:05:03
tags:
  - arch
  - archlinux
  - linux
category: linux
---

# 前言
本人也是第一次安装archlinux，严格来说是第一次安装成功，记录一下，既为自己也为新手。此方式是非UEFI模式，并且分区表使DOS的MBR方式，GPT分区表没有测试。以后也许会在虚拟机中测试过进行补充。

# 温馨提示
建议现在虚拟机中安装几次，直到安装成功，并且可以正常开机，上网，打开桌面环境，有十足把握之后再在物理机上安装，以免中间出现问题又没办法解决。并在虚拟机安装的过程中记录遇到的问题，以便日后参考。同时安装的时候记得备份重要文件，以免安装错误导致文件丢失。

# 准备安装介质
 1. 首先准备archlinux镜像，如果没有可以[点击这里下载](https://www.archlinux.org/download/)，最好选择中国的镜像服务，比如[网易的](http://mirrors.163.com/archlinux/iso/2015.01.01/)。下载完成后校验一下MD5值（官方文件的MD5值在md5sums.txt 这个文件中），如果相同那么可以进行下一步了；如果不相同需要重新下载并校验，不推荐在MD5值不同的情况下继续进行，因为不知道会发生什么问题。
 2. 刻录至U盘。如果用的是linux系统或者Mac系统（话说这么优雅的系统为啥要换呢，也可能是双系统吧），可以使用 `dd` 命令。把U盘插入计算机， 输入命令 `ls -al /dev/sd*`, 一般sdb是你的U盘，也请先做好文件备份。 现在假定U盘是 /dev/sdb, archlinux的文件路径是 /home/user/archlinux.iso,那么输入命令(需root权限) `dd -if=/home/user/archlinux.iso -of=/dev/sdb`，然后等待命令执行完毕，如果没有任何提示，则代表成功了。

## 环境准备
   ** 再次提醒，做好文件备份 **
 1. 如果没有分区的话，先进行分区，并进行格式化，如果已经操作过了或者想重用上次系统(Linux)的分区，可以直接进入第2步。
   - 分区进行时：
敲入命令 `fdisk /dev/sda` (假定操作的磁盘时sda，请自行确认好，此操作要格外小心)。输入 `m` 可以查看帮助， `n` 是新建一个分区， `d` 是删除一个已有分区。如果想新建一个DOS分区表，则输入 `o`，已经有分区表，想重新分区的话，按 `d` ，直到删除所有分区。分区方案可以按照以下来： `/boot`  大概需要  `200M`， `/` 可以分配 `15G` ～ `40G`， `/var` `8G`～ `20G` (可选)  ， `/tmp` `4G` ～ `8G` (可选)，其余分给 `/home`分区（强烈建议单独分区，以后重装系统可以不用拷贝主目录下的资料了）， `swap` `4G` ～ `8G`（可选）。新建分区输入 `n`，默认（p主分区）即可，然后默认（分区号1），接下来也是默认扇区既可以，然后选择大小，可以输入G,M,K单位的大小，我们输入 `+200M`，确定；然后创建根分区，按`n`,一路下来，大小选择输入 `+15G`,确定，根分区创建完成。如果分区少于4个，可以按照上面步骤，直到分区创建完成；但是如果分区多于4个，就要创建扩展分区，然后再创建逻辑分区了。扩展分区的创建和上面一样，只不过在选择分区格式的时候不是输入 `p` 了，而是 `e`，其余一样的。创建逻辑分区的时候输入 `l` （英文L的小写字母），剩下的步骤也是和创建主分区一样的啦。所有分区创建完成后，输入 `w` ，上面的一系列操作才会真正写入磁盘，再次之前都是在内存中，所以，在按 `w` 之前，还是有后悔药吃的，但是按下之后，那就定格了。切记！
   - 格式化分区：
格式化分区的命令是 `mkfs.xxx`，输入 `mkfs.`，按 `Tab` 键可以看到有如下格式：  `mkfs.bfs`       `mkfs.ext2`      `mkfs.ext4`      `mkfs.jfs`       `mkfs.reiserfs` `mkfs.cramfs`    `mkfs.ext3`      `mkfs.ext4dev`  ` mkfs.minix`     `mkfs.xfs`。咦，好像没有swap分区格式，swap分区格式化的命令是 `mkswap` 啦。输入命令 `mkfs.ext4 /dev/sda1` 将 `/boot` 分区格式化成ext4格式的分区，根分区和其他非swap分区用此方法依次格式化，用 `mkswap /dev/sdax` 格式化上面分的swap分区，x是分swap分区所得的号码。


 ## 安装
  1. mount 相关分区。 `mkdir /mnt/home /mnt/tmp /mnt/var /mnt/boot`  创建home，tmp，var，boot挂载点目录，然后 `mount /dev/sdax /mnt/xx` x代表分区号，xx代表目录，把的分区挂载到相应挂载点上。
  2. 修改 `/etc/pacman.d/mirrorlist` 的镜像列表，可以删除所有的，然后输入
>
Server = http://mirrors.163.com/archlinux/$repo/os/x86_64

        保存退出。
3. 执行 `pacstrap /mnt base` 命令进行基础安装。
4. 生成fstab。 `genfstab -p /mnt >> /mnt/etc/fstab`, 查看一下/mnt/etc/fstab 内容格式是否正确，有无重复内容，如有请先订正。格式大体如下：
 >
\#
\# /etc/fstab: static file system information
\#
\# <file system>	<dir>	<type>	<options>	<dump>	<pass>
\# UUID=ee8bae58-9428-4917-b63e-0258d19a4567
/dev/sda5           	/         	ext4      	rw,relatime,data=ordered0 1
 \# UUID=cbac48fe-3345-4cba-96ec-acdbdc56d0ad
/dev/sda9           	/home     	ext4      	rw,relatime,data=ordered0 2
\# UUID=59e210c2-fced-4cdd-b631-d9a50ba82312
/dev/sda7           	/tmp      	ext4      	rw,relatime,data=ordered0 2
5.  切换到新系统的root目录下，命令  `arch-choot /mnt`
6. 设置主机名 `echo your_hostname > /etc/hostname` ， your_hostname换成你想要的，最好是纯英文。
7. 设置时区。 `ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime`。
8.  修改 `/etc/locale.gen` ， 添加以下内容
>
en_US.UTF-8 UTF-8
zh_CN.UTF-8 UTF-8

    执行 `locale-gen`，
9. 执行 `echo LANG=zh_CN.UTF-8 > /etc/locale.conf`
10. 设置键盘映射和字体，文件在 `/etc/vconsole.conf`，在这就保持默认配置了。
11. 设置root密码 `passwd` 然后输入密码，再输入一次确认。
12. 安装引导程序，这里用grub。 `pacman -Sy grub`，安装完成后，执行  `pacman-db-upgrade`， 然后再执行 `grub-install --target=i386-pc --recheck --debug /dev/sda` 安装grub引导到sda上。最后执行 `grub-mkconfig -o /boot/grub/grub.cfg` ，生成引导配置。
13. 重启， 执行 `reboot`。如果成功安装的话，会出现grub引导选择系统菜单，选择默认的进入，输入root用户名，输入密码，登录成功。至此，安装已经完成，接下来是配置。


## 配置
  1. 网络配置
     - 查看网络设备名称 `ls /sys/class/net`, 记住所看到的网卡接口名称，假定叫做eth0
     - 启用网络接口 ip link set eth0 up
     - 检查结果状态 `ip link show dev eth0 ` 如果打印
     >
  enp3s0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP mode DEFAULT    group default qlen 1000
    link/ether 00:e0:66:cb:e2:1e brd ff:ff:ff:ff:ff:ff

     类似内容，说明启用成功。
  2.  创建或编辑 `/etc/systemd/network/dhcp.network` ,添加以下内容：
    >
[Match]
Name=en*
 [Network]
 DHCP=v4


  3.  启用网络服务 `systemctl enable systemd-resolved`
  4.  编辑 `/etc/resolv.conf` 配置dns ， 添加以下内容：
>
nameserver 8.8.8.8
nameserver 4.4.4.4

  如果你的IP段在192.168.xxx.yyy,则再添加 nameserver 192.168.xxx.1

  5. 执行 `dhcpd` 启用dhcp，要开机自动启动dhcp服务，则执行 `systemctl enable dhcpd`
  基本环境配置已经完成。

# 桌面环境配置
安装 fxce4


pacman -S xorg xorg-server
pacman -S slim #登录管理器
pacman -S xfce4
pacman -S xfce4-goodies
pacman -S fortune-mode
pacman -S gamin

  1. 创建用户 `useradd -Um du`
  2. 设置密码 `passwd du`
  3. 切换用户  `su -l du`
  4. 输入 `startxfce4` 可以进入xfce桌面了

# 美化显示：
  - 字体
    1. 首先可以从windowns上或者其他地方准备字体文件，然后 `cp *.ttf ~/.fonts/`
    2.  建立字体缓存
 `mkfontscale`
 `mkfontdir`
 `fc-cache -fv`

  - 输入法 [传送门](https://wiki.archlinux.org/index.php/Fcitx_%28%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87%29#.E8.BE.93.E5.85.A5.E6.B3.95.E6.A8.A1.E5.9D.97)
输入法就安装fcitx小企鹅输入法了
    1. 安装输入法

        `pacman -S fcitx`  
    2. 配置输入法
        安装输入法其他模块
        `fcitx-ui-light` Fcitx 的轻量 UI.
`fcitx-fbterm` Fbterm 对 Fcitx 的支持。
`fcitx-table-extra` Fcitx 的一些额外码表支持，包括仓颉 3, 仓颉 5, 粤拼, 速成, 五笔, 郑码等等
`fcitx-table-other` Fcitx 的一些更奇怪的码表支持，包括 Latex, Emoji, 以及一大堆不明字符等等。
`kcm-fcitx` KDE 的 Fcitx 输入法模块
    3. 启动桌面环境时候启用输入法
        在 .bashrc 文件中加入如下代码
        >
    export GTK_IM_MODULE=fcitx
 export QT_IM_MODULE=fcitx
 export XMODIFIERS="@im=fcitx"

    退出用户，重新登陆，可以欢快的使用输入法了。

ps:
浙大源：`Server = http://mirrors.zju.edu.cn/archlinux/$repo/os/$arch`
网易源：`Server = http://mirrors.163.com/archlinux/$repo/os/x86_64`
北京交大：`Server = http://mirror.bjtu.edu.cn/ArchLinux/$repo/os/x86_64`
