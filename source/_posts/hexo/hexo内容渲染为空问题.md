---
title: hexo 内容渲染为空问题
category: hexo
tags:
  - hexo
abbrlink: 45df
date: 2020-07-16 09:51:01
updateDate: 2020-07-16 09:51:01
---

如果使用的node版本大于12，并且hexo版本是4.2.0及其以下，在使用 `hexo g` 生成静态页面的时候，会发现生成的内容全部为空，要想解决这个问题需要升级一下hexo版本到 `4.2.1` 及其以上，因为`4.2.0`不兼容node最新版本。

issue见[渲染问题](https://github.com/hexojs/hexo/issues/4289)。

另一种解决方法是降低node版本，可以使用 `12` 这个版本同样可以解决（网友测试）。
