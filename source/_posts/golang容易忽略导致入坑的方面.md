---
title: golang容易忽略导致入坑的方面
date: 2016-03-15 18:19:19
tags:
  - golang
category: golang
---
# main 函数

main函数必须定义在main包里，不可导出。当main函数调用完毕，程序就会立即退出，不会等待运行中的goroutine运行完毕。



# 初始化函数

` func init() { … }`

此函数无参，不需显示调用。

此函数根据依赖顺序执行初始化顺序，并且是顺序执行，如：

A 依赖 B， B 依赖 C，并且A，B，C 都有init函数，那么初始化顺序是C ，B ， A



# iota

- iota 只能用在const定义的常量上

- iota在一个const范围，不管中间的变量有没有用到iota，iota都是持续递增加1的，如



``` go

const (

       aa = iota

       bb

       cc = 3

       dd = iota

)

```

dd 是3



- 如果重新定义const，那么值会是重新从0开始，如



``` go

const (

       ee = iota

)

```

# 类型断言

x.(T) x为变量， T为要判断的类型，其中有两种用法

- `a := x.(T)` ，如果x是类型T的实现类型变量，那么可以得到转换后的值a，否则panic: interface conversion: (x的类型) is not package.T: missing method （T的方法）

- `a, ok := x.(T)`，如果x是类型T的实现的类型变量，那么ok为true， 转换成功，否则，ok为false，说明转换失败。

# 从Panic中恢复

需要用到build-in recover()方法，函数原型：

`func recover() interface{}`

要捕获panic，必须在有可能发生panic的地方的上面进行处理，


``` go

defer func(){

 if err := recover(); err !=nil {

   //处理错误

 }

}()

```

一般来说，会将上面的代码放在函数（或者方法）的开始的地方。

** 注意 ， recover函数只能用在defer中，否则不会进行错误捕获。**
