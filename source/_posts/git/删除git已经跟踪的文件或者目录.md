---
title: 删除git已经跟踪的文件或者目录
tags:
  - git
category: git
abbrlink: 18216
date: 2016-03-13 11:09:29
---

如果第一次提交的时候，没有在gitignore文件中添加忽略文件，那么这些文件（目录也是文件）就会被git跟踪，push的时候也会被推送到远程。所以最好就是一开始在commit之前先添加到gitignore中。

如果文件已经被跟踪且被推送到远程，可以按照下面方法解决：
1.` rm -rf 文件`
2. `git rm -r --cached 要忽略的文件`
3.` git add -A (添加所有)`
4.` git push origin 分支`

如果同名的文件过多，如：.class 文件被提交了，那么如果这样一个一个显然效率太低，可以按照下面方法操作
1. `find . -iname 文件名 -exec rm -rf {}\;`
2. 重复上面的步骤，文件名替换为下一个要删除的文件名
3. 修改gitignore，添加忽略文件
4.   `git rm -r --cached 要忽略的文件`
5. `git add -A`
6. `git push origin 分支`
