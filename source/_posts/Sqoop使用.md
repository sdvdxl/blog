---
title: Sqoop使用
date: 2016-03-09 12:51:43
tags:
  - sqoop
	- hdfs
	- mysql
category: sqoop
---

`sqoop help` 查看帮助信息
`sqoop help COMMAND` 查看 COMMAND具体的帮助，如要查看 list-databases 命令的用法，则使用 `sqoop help list-databases` 查看。

主要可用的命令如下：

|命令 | 功能|
|----|----|
| help | List available commands |
| import|  Import a table from a database to HDFS|
| list-databases | List available databases on a server |
| list-tables |  List available tables in a database |

主要参数说明

参数 |说明
----|-----
--connect| 用来指定jdbc链接url，如mysql的: jdbc:mysql://ip:port/database
--password | 指定密码， 安全起见，建议使用 -P 参数，交互式填写密码或者使用 --password-file参数
--password-file  | 指定密码的文件，从该文件中读取密码
--username  | 指定用户名

用help查看帮助，使用示例：
list-databases 是列出所有的数据库，sqoop help list-databases· 查看使用方法


使用示例，查看 本机上的mysql中的数据库
./sqoop  list-databases --connect jdbc:mysql://127.0.0.1:3306/test --username username -P
这样直接操作会提示找不到驱动，我们需要把对应的mysql驱动jar包放到$SQOOP/lib目录下，然后再次执行就可以了，或者用参数 -libjars 指定驱动jar包路径。

# 配置项说明
按照此处的配置项进行可避免文末的错误，如果遇到错误请参考文末错误说明和解决方法。
1. sqoop 要使用对应的hadoop版本，如使用的hadoo版本是2.0.4，那么对应的sqoop版本就要使用文件名包含hadoop2.0.4的信息的版本。
2. SQOOP_HOME   环境变量关系到sqoop运行时选择的版本问题，所以该变量请配置成正确的版本路径。如果配置成了别的，虽然执行命令是在正确的路径下执行，而真实运行的版本却是其他的版本，该问题可以通过运行sqoop version 查看，此问题比较隐晦，要注意。
3. 执行sqoop所对应的SQOOP_HOME 文件要和hdfs文件系统上的一致，否则会产生找不到对应库文件的错误。
4. 在/etc/hosts 文件中增加 archeagle 到 hdfs节点ip的映射，否则sqoop会用默认的ip映射，会连接不上。
5. 用户权限问题，可以在 文件 hadoop/etc/hadoop/hdfs-site.xml中增加或者修改 配置
    ``` xml
    <property>
     <name>dfs.namenode.acls.enabled</name>
     <value>false</value>
   </property>
   <property>
     <name>dfs.permissions</name>
     <value>false</value>
   </property>
   ```
6. hdfs 集群要启动yarn服务。


# import 的使用
常用参数说明

|    参数| 说明    |
|-----|-----|
| -fs  | 指定hdfs节点  |
| --target-dir | 要到处到hdfs文件系统上的文件路径 |
| --table | 要导出的表名 |
| --connect | jdbc url |
| --username | 数据库用户名 |
| -P | 从控制台输入密码 |

使用示例 ：
``` bash
bin/sqoop  import -fs hdfs://192.168.6.63:9000 --target-dir /user/admin/export_test_admin_user11  --table admin_user --connect jdbc:mysql://192.168.6.201:3306/test --username username -P
```
## 增量导入 [原始链接](http://sqoop.apache.org/docs/1.4.5/SqoopUserGuide.html#_incremental_imports)
主要参数如下：

参数 | 说明
-----|-----
--incrementa | 增量方式， 有两种方式，lastmodified和append
--last-value | 以lastmodified方式的增量追加，要指定时间；append则要指定偏移id
--check-column | 要检查的字段， 即以哪个字段为标准计算增量范围
--append | 指定以增量方式追加

使用增量导入（以时间为标识作参考）：
``` bash
bin/sqoop  import -fs hdfs://192.168.6.63:9000 --target-dir /user/admin/export_test_admin_user11  --table admin_user --connect jdbc:mysql://192.168.6.201:3306/forseti_core --username forseti -P--incremental lastmodified --check-column gmt_create --last-value '2012-02-01 11:0:00' --verbose --append
```
使用增量导入（以id为标识作为参考）：
``` bash
bin/sqoop  import -fs hdfs://192.168.6.63:9000 --target-dir /user/admin/export_test_admin_user11  --table admin_user --connect jdbc:mysql://192.168.6.201:3306/forseti_core --username forseti -P--incremental append --check-column id --verbose --append
```

## 使用select语句(-e或者--query参数)
如果使用这个参数，那么可以执行自定义语句，比如可以执行join操作等其他复杂sql语句，但是语句中where是必须的，而且where后面要加 $CONDITIONS 参数。sql语句本身可以用单引号包裹，但是如果sql语句中已经包含了单引号，那么可以用双引号包裹。另外，使用了这个参数，那么参数 --split-by 在import命令中是必须的，而且该参数后面指定的字段必须出现在sql查询结果中。因为通过观察sqoop执行过程中输出的执行sql可以发现，它是在原有的sql上包裹一层，如下示例中，结果就变成了 SELECT MIN(gmt_modified), MAX(gmt_modified) FROM (select id from admin_user where  (1 = 1) ) AS t1。
使用示例：
``` bash
bin/sqoop  import  --connect jdbc:mysql://192.168.6.201:3306/test --username username -P -e "select id from test where $CONDITIONS" --split-by id
```


# job 使用

## 主要参数

| 参数 | 说明 |
|------ |------|
| --create <job-id> | Create a new saved job |
| --delete <job-id> | Delete a saved job |
| --exec <job-id> | Run a saved job |
| --help | Print usage instructions |
| --list | List saved jobs |
| --show <job-id> | Show the parameters for a saved job |
| -fs <local|namenode:port> | specify a namenode |
| -libjars <comma separated list of jars> | specify comma separated jar files to include in the classpath. |
|-conf <configuration file> | specify an application configuration file |

## 创建Job示例：
``` bash
bin/sqoop job --create export_mysql_table -- import --table admin_user --connect jdbc:mysql://192.168.6.201:3306/forseti_core
```

## 执行Job示例：
``` bash
bin/sqoop job -fs hdfs://192.168.6.63:9000 --exec  export_mysql_table --  --username forseti -P --target-dir /user/admin/export_test_admin_user11112
```
## 执行带密码的任务
有密码要求的任务，如果不存储密码的话，每次执行任务都要求手动输入密码，如果是定时任务，那么这个肯定是不合理的。默认metastore是不保存密码的，如果需要保存，则在conf/sqoop-site.xml增加或者取消注释如下内容

``` xml
 <property>
    <name>sqoop.metastore.client.record.password</name>
    <value>true</value>
    <description>If true, allow saved passwords in the metastore.
    </description>
  </property>
```

# 错误解决

- ERROR tool.ImportTool: Encountered IOException running import job: java.io.FileNotFoundException: File does not exist: hdfs://192.168.6.63:9000/home/du/software/dev/sqoop-1.4.5.bin__hadoop-0.20/lib/ant-contrib-1.0b3.jar

        在不同机器或者用户下执行sqoop，会查找hadoop集群指定的节点上的hdfs目录中的这个文件，比如我是用在/home/du/software/dev/sqoop-1.4.5.bin__hadoop-0.20 下执行的sqoop，并且SQOOP_HOME配置的也是这个路径，那么到hdfs://192.168.6.63:9000上就会查找/home/du/software/dev/sqoop-1.4.5.bin__hadoop-0.20/lib这个路径下的ant-contrib-1.0b3.jar这个文件，解决方法就是在hdfs上创建对应目录，并把sqoop拷贝到对应目录，目录结构和执行sqoop的目录结构一样即可。

- Exception in thread "main" java.lang.IncompatibleClassChangeError: Found interface org.apache.hadoop.mapreduce.JobContext, but class was expected

        使用的hadoop版本问题，从2.6.0切换到2.4.0 解决

- ERROR tool.ImportTool: Encountered IOException running import job: org.apache.hadoop.security.AccessControlException: Permission denied: user=du, access=WRITE, inode="/user":admin:supergroup:drwxr-xr-x

- ERROR manager.SqlManager: Error executing statement: java.sql.SQLException: Access denied for user 'forseti'@'192.168.6.165' (using password: YES)

        很明显是mysql的用户登陆失败，填写正确的用户名和密码即可解决该问题。

- 15/03/05 17:40:10 INFO mapreduce.Job: Running job: job_1425543105230_0006
15/03/05 17:40:44 INFO ipc.Client: Retrying connect to server: archeagle/220.250.64.20:43175. Already tried 0 time(s); maxRetries=3
15/03/05 17:41:04 INFO ipc.Client: Retrying connect to server: archeagle/220.250.64.20:43175. Already tried 1 time(s); maxRetries=3
15/03/05 17:41:24 INFO ipc.Client: Retrying connect to server: archeagle/220.250.64.20:43175. Already tried 2 time(s); maxRetries=3
15/03/05 17:41:44 INFO mapred.ClientServiceDelegate: Application state is completed. FinalApplicationStatus=FAILED. Redirecting to job history server
15/03/05 17:41:44 ERROR tool.ImportTool: Encountered IOException running import job: java.io.IOException: Job status not available
	at org.apache.hadoop.mapreduce.Job.updateStatus(Job.java:322)
	at org.apache.hadoop.mapreduce.Job.isComplete(Job.java:599)
	at org.apache.hadoop.mapreduce.Job.monitorAndPrintJob(Job.java:1344)
	at org.apache.hadoop.mapreduce.Job.waitForCompletion(Job.java:1306)
	at org.apache.sqoop.mapreduce.ImportJobBase.doSubmitJob(ImportJobBase.java:186)
	at org.apache.sqoop.mapreduce.ImportJobBase.runJob(ImportJobBase.java:159)
	at org.apache.sqoop.mapreduce.ImportJobBase.runImport(ImportJobBase.java:247)
	at org.apache.sqoop.manager.DirectMySQLManager.importTable(DirectMySQLManager.java:92)
	at org.apache.sqoop.tool.ImportTool.importTable(ImportTool.java:497)
	at org.apache.sqoop.tool.ImportTool.run(ImportTool.java:601)
	at org.apache.sqoop.Sqoop.run(Sqoop.java:143)
	at org.apache.hadoop.util.ToolRunner.run(ToolRunner.java:70)
	at org.apache.sqoop.Sqoop.runSqoop(Sqoop.java:179)
	at org.apache.sqoop.Sqoop.runTool(Sqoop.java:218)
	at org.apache.sqoop.Sqoop.runTool(Sqoop.java:227)
	at org.apache.sqoop.Sqoop.main(Sqoop.java:236)

	       在运行sqoop的主机hosts文件增减加hadoop节点ip映射 192.168.6.63 archeagle

- 使用--direct参数
 Error: java.io.IOException: Cannot run program "mysqldump": error=2, No such file or directory
	at java.lang.ProcessBuilder.start(ProcessBuilder.java:1047)
	at java.lang.Runtime.exec(Runtime.java:617)
	at java.lang.Runtime.exec(Runtime.java:485)
	at org.apache.sqoop.mapreduce.MySQLDumpMapper.map(MySQLDumpMapper.java:405)
	at org.apache.sqoop.mapreduce.MySQLDumpMapper.map(MySQLDumpMapper.java:49)
	at org.apache.hadoop.mapreduce.Mapper.run(Mapper.java:145)
	at org.apache.hadoop.mapred.MapTask.runNewMapper(MapTask.java:764)
	at org.apache.hadoop.mapred.MapTask.run(MapTask.java:340)
	at org.apache.hadoop.mapred.YarnChild$2.run(YarnChild.java:167)
	at java.security.AccessController.doPrivileged(Native Method)
	at javax.security.auth.Subject.doAs(Subject.java:415)
	at org.apache.hadoop.security.UserGroupInformation.doAs(UserGroupInformation.java:1548)
	at org.apache.hadoop.mapred.YarnChild.main(YarnChild.java:162)
Caused by: java.io.IOException: error=2, No such file or directory
	at java.lang.UNIXProcess.forkAndExec(Native Method)
	at java.lang.UNIXProcess.<init>(UNIXProcess.java:186)
	at java.lang.ProcessImpl.start(ProcessImpl.java:130)
	at java.lang.ProcessBuilder.start(ProcessBuilder.java:1028)
	... 12 more
	- ERROR tool.ImportTool: Encountered IOException running import job: org.apache.hadoop.ipc.RemoteException(org.apache.hadoop.ipc.RetriableException): org.apache.hadoop.hdfs.server.namenode.SafeModeException: Cannot delete /user/du/.staging/job_1425543105230_0010. Name node is in safe mode.
The reported blocks 0 needs additional 963 blocks to reach the threshold 0.9990 of total blocks 963.
The number of live datanodes 0 has reached the minimum number 0. Safe mode will be turned off automatically once the thresholds have been reached.
	at org.apache.hadoop.hdfs.server.namenode.FSNamesystem.checkNameNodeSafeMode(FSNamesystem.java:1199)
	at org.apache.hadoop.hdfs.server.namenode.FSNamesystem.deleteInternal(FSNamesystem.java:3336)
	at org.apache.hadoop.hdfs.server.namenode.FSNamesystem.deleteInt(FSNamesystem.java:3296)
	at org.apache.hadoop.hdfs.server.namenode.FSNamesystem.delete(FSNamesystem.java:3280)
	at org.apache.hadoop.hdfs.server.namenode.NameNodeRpcServer.delete(NameNodeRpcServer.java:727)

            hdfs上(用户)目录不存在。

- INFO ipc.Client: Retrying connect to server: arch57/220.250.64.20:56564. Already tried 2 time(s); maxRetries=3
15/03/10 15:47:55 INFO mapred.ClientServiceDelegate: Application state is completed. FinalApplicationStatus=SUCCEEDED. Redirecting to job history server
15/03/10 15:47:55 ERROR tool.ImportTool: Encountered IOException running import job: java.io.IOException: Job status not available
	at org.apache.hadoop.mapreduce.Job.updateStatus(Job.java:322)
	at org.apache.hadoop.mapreduce.Job.isComplete(Job.java:599)
	at org.apache.hadoop.mapreduce.Job.monitorAndPrintJob(Job.java:1344)
	at org.apache.hadoop.mapreduce.Job.waitForCompletion(Job.java:1306)
	at org.apache.sqoop.mapreduce.ImportJobBase.doSubmitJob(ImportJobBase.java:186)
	at org.apache.sqoop.mapreduce.ImportJobBase.runJob(ImportJobBase.java:159)
	at org.apache.sqoop.mapreduce.ImportJobBase.runImport(ImportJobBase.java:247)
	at org.apache.sqoop.manager.SqlManager.importTable(SqlManager.java:665)
	at org.apache.sqoop.manager.MySQLManager.importTable(MySQLManager.java:118)
	at org.apache.sqoop.tool.ImportTool.importTable(ImportTool.java:497)
	at org.apache.sqoop.tool.ImportTool.run(ImportTool.java:601)
	at org.apache.sqoop.tool.JobTool.execJob(JobTool.java:228)
	at org.apache.sqoop.tool.JobTool.run(JobTool.java:283)
	at org.apache.sqoop.Sqoop.run(Sqoop.java:143)
	at org.apache.hadoop.util.ToolRunner.run(ToolRunner.java:70)
	at org.apache.sqoop.Sqoop.runSqoop(Sqoop.java:179)
	at org.apache.sqoop.Sqoop.runTool(Sqoop.java:218)
	at org.apache.sqoop.Sqoop.runTool(Sqoop.java:227)
	at org.apache.sqoop.Sqoop.main(Sqoop.java:236)

        在执行sqoop的机器的hosts增加 arch57 这个主机ip映射（PS:arch57 是一台hadoop机器的名字）
