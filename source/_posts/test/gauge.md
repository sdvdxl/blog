---
title: 自动测试工具Gauge
tags:
  - gauge
  - 测试
category: test
abbrlink: 33607
date: 2017-07-31 14:49:23
---

Gauge 是一款轻量级、跨平台自动化测试工具集。规则文件语法可以使用markdown语法编写。另外还可以使用你喜欢的语言来编写业务规则代码比如 `go`、 `java`、`ruby`等语言。Gauge还提供了输出插件，可以将执行结果导出为html或者xml，甚至flash，方便查看。

# Quick Start
有小伙伴不喜欢翻译官方文档的内容，让我开门见山，直接实战。那我们直接按照步骤来创建一个 JAVA 版的测试项目。前提是已经安装好了 gauge ，html-report 和 JAVA 插件，如果需要安装帮助，请参考下面的安装步骤或者直接参考官网的[安装文档](https://docs.getgauge.io/plugins.html#installation)。

## 创建项目
命令行执行 `gauge init java`，执行成功后悔看到下面的日志输出:

```
Downloading java.zip
...
Copying Gauge template java to current directory ...
Successfully initialized the project. Run specifications with "gauge specs/".
```

这样就创建了一个基于 JAVA 语言的测试项目，并且提供了demo 示例，切换到项目目录下，直接执行 `gauge run specs` 进行测试。会有以下输出:

```
# Specification Heading
  ## Vowel counts in single word	 ✔ ✔
  ## Vowel counts in multiple word	 ✔ ✔

Successfully generated html-report to => /Users/du/Desktop/test/reports/html-report
Specifications:	1 executed	1 passed	0 failed	0 skipped
Scenarios:	2 executed	2 passed	0 failed	0 skipped

Total time taken: 93ms
```

测试已经完成，并且会在项目目录下 `reports/html-report` 生成测试报告，直接用浏览器打开 index 即可查看测试报告:
![](http://note.youdao.com/yws/public/resource/9b09c5bfd845d727a234b59d5c599fa8/xmlnote/FE1AD5BE7CC14F31BCB44BFAC0E59B9D/1352) 

这是报告的首页，上面标注了总共多少个 spec 文件，失败个数，成功个数，跳过多少个。并且可以通过左下角的搜索或者直接点标题来查看对应的测试文件或者用例。
![](http://note.youdao.com/yws/public/resource/9b09c5bfd845d727a234b59d5c599fa8/xmlnote/FB2E7F3C889840A6B6A47A82A68DAFDC/1355)

每个spce 的测试报告都会将测试的用例结果打印出来并且用不同颜色表示出来，绿色则是成功，红色则是失败。其中一个大的圆点表示的是一个测试场景，圆点下每个绿色线条代表一个具体的测试用例。

下面我们来尝试修改一下  `The word "gauge" has "3" vowels` 这个用例。打开 `specs/example.spec` 这个文件，在第18行我们就会找到上面提到的这个测试用例，然后打开 `src/test/java/StepImplementation.java` 这个文件，第21行就是这个测试用例的实现代码。 可以看到，这个测试用例有两个参数，`word` 和 `expectedCount`，现在我们将 `specs/example.spec` 中的这个测试用例的参数  "expectedCount" 改为 1，然后重新运行 `gauge run specs` ，执行过程中控制台抛出了异常信息，接收后，打开 生成的 html 文件查看，有没有通过的测试用例: 
![](http://note.youdao.com/yws/public/resource/9b09c5bfd845d727a234b59d5c599fa8/xmlnote/2E49A843026C460993E0456FDF7EAD43/1357)

尝试增加一个场景和一个用例:
在 `specs/example.spec` 中最后追加以下代码:

```markdown
测试场景
---------------

* 新增的一个测试"成功"
```

然后在 `StepImplementation.java` 中增加下面的方法:

```java
    @Step("新增的一个测试<param>")
    public void testParam(String param) {
        assertEquals("成功",param);
    }
```

运行 `gauge run specs`，执行结束后打开报告，我们新定义的场景和用例豆子行成功了；现在再修改一下 `sample.spce`文件，在最后再追加以下代码:

```markdown
 新增的一个测试"失败"
```
运行 `gauge run specs`，执行结束后打开报告，我们新定义的场景和用例豆子行失败了，因为参数值和预期不一致，我们的代码抛出了异常，则代表该测试用例是失败的。

细心的同学可能发现，每个场景开头都有 `Vowels in English language are "aeiou".` 这个测试用例，这是因为这个测试用例是在 spec 层级下，会在执行过程中追加到每个场景的开始。另外会发现有切换窗口的动作，并且在失败的用例中有图片生成，这个则是默认配置，失败则截取屏幕图片的配置造成的，如果不想截图，则可以修改 `env/default/default.properties` 文件第13行的 screenshot_on_failure 属性为 false 即可。

至此，我们学会了如何运行测试用例，如何编写测试用例，如果需要其他知识点，比如数据到不同范围，如何在用例失败的情况下继续执行该场景的其他用例，请参考下面的介绍和官网知识。

# 主要特点
- A rich markup based on markdown
- Simple, Flexible and Rich Syntax
- Business Language Tests : Supports the concept of executable documentation.
- Consistent Cross Platform/Language Support for writing test code. Currently supported languages.
- Open Source, so it could be shared freely and improved by others as well.
- A modular architecture with Plugins support.
- Extensible through Plugins and Hackable.
- Supports External Data Sources.
- Helps you create Maintainable and Understandable test suites.
IDE Support.

# 概念(术语)
## Specifications (spec)
一个spec就是一个spec文件，用来定义规则。
使用markdown的一级标题来定义该文件的说明 ，比如:
```
Specification name
==================
```
或者
```
# Specification name
```

## Scenarios
一个 scenarios 可以认为是一个组，定义在spec文件中。使用markdown的二级标题定义
```
Scenarios
----------------
```
或者
```
## Scenarios
```

## Tags 
可以给 Spec 文件和 Scenarios 打标签。

```
Specification name
================

tags: s1,test1


Scenarios
----------------

tags: s1, s2
```
## [Steps](https://docs.getgauge.io/syntax.html#steps)
这个就是具体的测试用例。使用markdown的 `*` 声明一个用例。

```
* Step Name
```
step 还支持参数，包括静态参数，动态参数，表格参数  ，特殊参数
## 静态参数
静态参数用双引号包裹:
```
* Check "product" exists
```
## 动态参数
动态参数使用尖括号包裹，跟下面的特殊参数结合使用

```
* Check <product> exists
```

## 表格参数
表格参数，就是一个二维表格所构成的数据。

```
* Step that takes a table
   | id  | name    |
   |-----|---------|
   | 123 | John    |
   | 456 | Mcclain |
```
** 注意，定义和表格之间没有空行，表格是直接跟在定义的语句下面 **

## 特殊参数
```
<prefix:value>
```
prefix 支持 `file` 和 `table`

```
// file 
* Check if <file:/work/content.txt> is visible

// table
* Check if the users exist <table:/Users/john/work/users.csv>
```


## 注释
没有特殊的语法，任何一般的文本都可以作为注释

## 其他
支持图片，连接


# 安装
首先要安装Gauge程序，可以从这里[https://getgauge.io/get-started.html](https://getgauge.io/get-started.html)下载对应平台的安装程序。Windows平台没有测试。
## Mac
mac用户可以使用brew安装，`brew update && brew install gauge`，也可以从上面的地址下载安装包，然后执行安装即可。

## Linux
linux用户需要下载linux程序，下载下来是个安装包，假设我们放到 `/opt/gauge/gauge-0.9.0-linux.x86_64.zip` (最好创建一个目录，因为解压后不会生成单独的目录)。切换到 `/opt/gauge` 执行 `unzip gauge-0.9.0-linux.x86_64.zip`，得到以下文件：
```

├── bin
│   ├── gauge
│   └── gauge_screenshot
└── install.sh
```

执行 `./install.sh` ，一直按回车，使用默认配置即可安装完成。

## 其他环境
如果要使用java，则需要安装jdk，如果要使用C#，则需要安装.net环境，总之，需要什么语言，就需要安装什么环境。

安装完成后，命令行输入`gauge`，如果出现 gauge 的帮助信息，则说明安装成功，否则请按照[官方文档](https://getgauge.io/get-started.html)进行安装。

## 插件
Gauge 支持一些插件，比如生成java和执行代码的插件，生成html报告的插件等等，插件支持:
- java
- ruby
- flash
- go
- js
- html-report
- xml-report
- python
- spectacle
(可能还有其他的，暂时没有去搜索。)

插件安装方式: `gauge install <plugin-name>`，如果下载太慢，可以手动下载，然后使用本地安装方式安装 `gauge install --file <path-plugin>`

# 实战
这里我们使用的是JAVA方式。所以这里需要安装JAVA插件 `gauge install java`，并且我们要以html方式展示报告，需要安装 html report 插件 `gauge install html-report`

执行 `gauge init java` 会生成一个java的gauge项目:
```
├── env
│   └── default
│       ├── default.properties
│       └── java.properties
├── libs
├── manifest.json
├── specs
│   └── example.spec
└── src
    └── test
        └── java
            └── StepImplementation.java
```
想要执行测试，执行命令 `gauge run specs` 即可，   最后会输出执行结果，同时会在项目下生成`reports/html-report` html报告，用浏览器打开即可查看结果。

## 结构说明
目录 | 说明
---------|---------
env | 环境变量，可以配置不通参数，适用于不通环境，比如测试环境还是线上环境
libs |  其他要依赖的java类库要放到这里
manifest.json | 项目配置文件
spces | 测试用例描述文件存放的地方
src | java 代码放置地方

## 命令说明
上面我们使用了命令 `gauge run spces` 来运行测试用例，`run` 后面跟的 `spces`就是项目中的 `specs` 目录，代表执行这个目录下的所有文件，也可以单独执行一个文件执行，比如 `gauge run specs/example.spec`。

## 文件编写
用例定义在spec文件中，具体实现的代码则是在java代码中实现的。

# 进阶
