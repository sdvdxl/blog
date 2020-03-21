---
title: MySql避免重复插入记录
tags:
  - mysql
category: 数据库
abbrlink: 23635
date: 2016-03-09 14:02:47
---

# 方案一：使用ignore关键字
如果是用主键primary或者唯一索引unique区分了记录的唯一性,避免重复插入记录可以使用：
``` sql
insert ignore into table_name(email,phone,user_id) values('test9@163.com','99999','9999')
```
这样当有重复记录就会忽略,执行后返回数字0,还有个应用就是复制表,避免重复记录：
``` sql
insert ignore into table(name)  select  name from table2
```

# 方案二：使用Replace

replace的语法格式为：
1. `replace into table_name(col_name, ...) values(...) `
2. `replace into table_name(col_name, ...) select ... `
3. `replace into table_name set col_name=value, ...`

## 算法说明：
REPLACE的运行与INSERT很相像,但是如果旧记录与新记录有相同的值，则在新记录被插入之前，旧记录被删除，即：
1. 尝试把新行插入到表中
2. 当因为对于主键或唯一关键字出现重复关键字错误而造成插入失败时：
  - 从表中删除含有重复关键字值的冲突行
  - 再次尝试把新行插入到表中
旧记录与新记录有相同的值的判断标准就是：表有一个PRIMARY KEY或UNIQUE索引，否则，使用一个REPLACE语句没有意义。该语句会与INSERT相同，因为没有索引被用于确定是否新行复制了其它的行。
## 返回值：
REPLACE语句会返回一个数，来指示受影响的行的数目。该数是被删除和被插入的行数的和。受影响的行数可以容易地确定是否REPLACE只添加了一行，或者是否REPLACE也替换了其它行：检查该数是否为1（添加）或更大（替换）。
## 示例:
eg:(phone字段为唯一索引)
`replace  into table_name(email,phone,user_id) values('test569','99999','123')`
另外：在 SQL Server 中可以这样处理：
``` sql
if not exists (select phone from t where phone= '1')
    insert into t(phone, update_time) values('1', getdate())
else
    update t set update_time = getdate() where phone= '1'
```
更多信息[请看](http://dev.mysql.com/doc/refman/5.1/zh/sql-syntax.html#replace)

# 方案三：ON DUPLICATE KEY UPDATE
如‍上所写，你也可以在`INSERT INTO.....`后面加上 `ON DUPLICATE KEY UPDATE`方法来实现。如果您指定了`ON DUPLICATE KEY UPDATE`，并且插入行后会导致在一个UNIQUE索引或PRIMARY KEY中出现重复值，则执行旧行UPDATE。例如，如果列a被定义为UNIQUE，并且包含值1，则以下两个语句具有相同的效果：
``` sql
INSERT INTO table (a,b,c) VALUES (1,2,3) ON DUPLICATE KEY UPDATE c=c+1; `
UPDATE table SET c=c+1 WHERE a=1;
```
如果行作为新记录被插入，则受影响行的值为1；如果原有的记录被更新，则受影响行的值为2。注释：如果列b也是唯一列，则INSERT与此UPDATE语句相当：
``` sql
UPDATE table SET c=c+1
WHERE a=1 OR b=2 LIMIT 1;
```
如果a=1 OR b=2与多个行向匹配，则只有一个行被更新。通常，您应该尽量避免对带有多个唯一关键字的表使用ON DUPLICATE KEY子句。
您可以在UPDATE子句中使用VALUES(col_name)函数从INSERT...UPDATE语句的INSERT部分引用列值。换句话说，如果没有发生重复关键字冲突，则UPDATE子句中的VALUES(col_name)可以引用被插入的col_name的值。本函数特别适用于多行插入。VALUES()函数只在INSERT...UPDATE语句中有意义，其它时候会返回NULL。
``` sql
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)
ON DUPLICATE KEY UPDATE c=VALUES(a)+VALUES(b);
```
本语句与以下两个语句作用相同：
``` sql
INSERT INTO table (a,b,c) VALUES (1,2,3)  ON DUPLICATE KEY UPDATE c=3;
INSERT INTO table (a,b,c) VALUES (4,5,6)  ON DUPLICATE KEY UPDATE c=9;
```
当您使用`ON DUPLICATE KEY UPDATE`时，`DELAYED`选项被忽略。
注：[来源](http://www.cnblogs.com/zeroone/archive/2012/04/18/2454728.html)
