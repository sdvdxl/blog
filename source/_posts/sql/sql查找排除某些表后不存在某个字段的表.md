---
title: sql查找排除某些表后不存在某个字段的表
tags:
  - mysql
  - sql
category: 数据库
abbrlink: 6879
date: 2016-03-09 14:03:56
---

``` sql
SELECT table_name FROM information_schema.tables
WHERE table_schema='database_name'
	AND table_name NOT LIKE 'table_name'
	AND table_name NOT IN(SELECT col.table_name  FROM information_schema.`COLUMNS` col
	       WHERE col.table_name IN (SELECT table_name FROM information_schema.tables
					WHERE table_schema='database_name'
					AND table_name NOT LIKE 'table_name'
					AND col.column_name='gmt_modified');
```
