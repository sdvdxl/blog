---
title: xattr 命令清除 app 的隔离属性
abbrlink: bf39800a
date: 2022-12-30 10:47:32
updateDate: 2022-12-30 10:47:32
top_img:
cover:
category: mac
tags:
  - mac
  - app
keywords: 
  - macos
  - app
  - 损坏
---

macOS 中如果安装一些修改版或破解版软件，通过拖拽方式放到 application 目录，在打开 app 的时候可能提示程序无法运行，通常我们在系统 preference - security 里可以手动允许运行，但是有时候使用这种方法也无法打开，这是由于系统识别到这个 app 可能有问题所以给它加上了 `com.apple.quarantine` 隔离属性阻止了他的运行。

如果我们需要运行它，就需要删除 app 的 `com.apple.quarantine` 属性，可以使用 `xattr` 来处理。

macOS 上的文件不只有 "normal" 属性，也有 "读", "写", "执行" 等其他属性。常规属性可以用 `ls -l myfile` 命令来查看。除此之外还可以定义**扩展属性**，扩展属性的修改可以用 `xattr` 来处理。

命令语法如下：

```
xattr [options] attributes [files]
```

可用 option：

```
   -c  CLear all Atrributes.
   -d  Delete the given attribute.
   -h  Help.
   -l  By default, the first two command forms either display just the attribute names or
       values, respectively. The -l option causes both the attribute names and corresponding
       values to be displayed. For hex display of values, the output is preceeded with the hex
       offset values and followed by ASCII display, enclosed by '|'.
   -p  Print the value associated with the given attribute.
   -r  If a file argument is a directory, act as if the entire contents of the directory
       recursively were also specified (so that every file in the directory tree is acted upon).
   -s  If a file argument is a symbolic link, act on the symbolic link itself, rather than
       the file that the symbolic link points at.
   -v  Force the the file name to be displayed, even for a single file.
   -w  Write a given attribute name with a value.
   -x  Force the attribute value to be displayed in the hexadecimal representation.   
```

显示一个文件的存在的扩展属性：

```
xattr myfile
```

显示一个文件的存在的扩展属性及这个属性的值：

```
xattr -l myfile
```

给一个文件添加 `com.example.color` 属性：

```
xattr -w com.example.color myfile
```

清除一个文件的 `com.example.color` 属性：

```
xattr -d com.example.color myfile
```

清除一个文件夹内所有文件的 `com.example.color` 属性：

```
xattr -d -r com.example.color mydir
```

清除一个文件的所有扩展属性：

```
xattr -c myfile
```

---

对于被屏蔽的 app 需要清除`com.apple.quarantine` 这一属性，处理过程如下：

首先查看 app 的现有属性：

```
xattr /path/to/MyApp.app
```

如果返回的结果有 `com.apple.quarantine` 属性则执行下面命令删除：

```
sudo xattr -r -d com.apple.quarantine /path/to/MyApp.app
```

执行后就可以正常打开 app 了。
