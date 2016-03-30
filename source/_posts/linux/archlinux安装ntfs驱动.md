---
title: archlinux安装ntfs驱动
date: 2016-03-09 14:08:56
tags:
  - archlinux
  - arch
  - ntfs
category: linux
---
默认情况下，archlinux本身支持挂在ntfs文件系统，只不过是只读，不能写入。如果要支持ntfs系统文件，那么需要安装ntfs的驱动程序。
用命令 ` yaourt -Ss ntfs` 可以查找关于ntfs的软件包
>
extra/ntfs-3g 2014.2.15-1 [installed]
    NTFS filesystem driver and utilities
aur/disk-manager 1.0.1-3 (56)
    A tool to manage filesystems, partitions, and NTFS write mode
aur/fgetty 0.7-5 (14)
    A mingetty stripped of the printfs
aur/fgetty-pam 0.7-4 (4)
    A mingetty stripped of the printfs, patched for PAM-support.
aur/grub4dos 0.4.5c_20140822-1 (35)
    A GRUB boot loader support menu on windows(fat,ntfs)/linux(ext2,3,4)
aur/libntfs-wii 2013.1.13-1 (1)
    NTFS-3G filesystem access library (for Nintendo Gamecube/Wii homebrew
    development)
aur/ntfs-3g-ar 2014.2.15AR.3-1 (37)
    NTFS filesystem driver and utilities with experimental features
aur/ntfs-3g-fuse 2014.2.15-1 (46)
    Stable read and write NTFS driver and ntfsprogs. This package will allow
    normal users to mount NTFS Volumes.
aur/ntfs-3g_ntfsprogs-git 4695.db35a16-1 (7)
    Read and write NTFS driver and utilities - GIT version
aur/ntfs-config 1.0.1-13 (119)
    Enable/disable NTFS write support with a simple click
aur/ntfsfixboot 1.0-3 (18)
    Fix NTFS boot sector
aur/scrounge-ntfs 0.9-2 (28)
    Data recovery program for NTFS file systems
aur/ufsd-module 8.9.0-3 (9)
    Paragon NTFS & HFS for Linux driver. - ACLs removed
aur/ufsd-module-dkms 8.9.0-3 (4)
    Paragon NTFS & HFS for Linux driver. - ACLs removed. DKMS version
aur/wipefreespace 2.0-1 (3)
    Securely wipe the free space on an ext2/3/4,NTFS, XFS,ReiserFSv3,
    ReiserFSv4, FAT12/16/32,Minix,JFS and HFS+ partition or drive

可以看到有这么多相关软件包，其实只需要安装第一个也就是 `extra/ntfs-3g` 就可以了，输入命令 `yaourt -S extra/ntfs-3g`，安装完毕后重新插入ntfs分区U盘或者移动硬盘就可以进行写入操作了。
