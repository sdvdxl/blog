---
title: sparkApp提交到SpringXD出现错误可能情况及解决方法
date: 2016-03-16 11:52:19
tags:
  - spark
  - spring xd
category: spring
---
编写spark代码：

```java
package com.demo;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import org.apache.spark.SparkConf;
import org.apache.spark.api.java.JavaPairRDD;
import org.apache.spark.api.java.JavaRDD;
import org.apache.spark.api.java.JavaSparkContext;
import org.apache.spark.api.java.function.FlatMapFunction;
import org.apache.spark.api.java.function.Function2;
import org.apache.spark.api.java.function.PairFunction;
import scala.Tuple2;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by sdvdxl on 2016/3/14.
 */
public class SparkCalcDemo {
    private static final String HADOOP_URL = "hdfs://10.10.1.110:8020/";

    public static void main(String[] args) throws Exception {
        SparkConf conf = new SparkConf().setAppName("test").setMaster("local[1]");
        JavaSparkContext sc = new JavaSparkContext(conf);
        JavaRDD<String> textFile = sc.textFile(HADOOP_URL + "/xd/dataset1/2016/03/14/15/01", 1);
        JavaRDD<String> words = textFile.flatMap(new FlatMapFunction<String, String>() {
            public Iterable<String> call(String s) {
                List<String> list = new ArrayList<String>();

                JSONObject jobj = JSON.parseObject(new String(org.apache.commons.codec.binary.Base64.decodeBase64(s.substring(1, s.length() - 1))));
                list.add(jobj.getString("name"));
                list.add(jobj.getString("random"));
                return list;
            }
        });
        JavaPairRDD<String, Integer> pairs = words.mapToPair(new PairFunction<String, String, Integer>() {
            public Tuple2<String, Integer> call(String s) {
                return new Tuple2<String, Integer>(s, 1);
            }
        });
        JavaPairRDD<String, Integer> counts = pairs.reduceByKey(new Function2<Integer, Integer, Integer>() {
            public Integer call(Integer a, Integer b) {
                return a + b;
            }
        });

        counts.foreach(tuple2 ->
                System.out.println(tuple2._1 + " : " + tuple2._2));
    }
}
```

pom依赖：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	<modelVersion>4.0.0</modelVersion>

	<groupId>kafka-demo</groupId>
	<artifactId>kafka-demo</artifactId>
	<version>1.0-SNAPSHOT</version>

	<dependencies>
		<dependency>
			<groupId>junit</groupId>
			<artifactId>junit</artifactId>
			<version>4.7</version>
			<scope>test</scope>
		</dependency>

		<dependency>
			<groupId>com.alibaba</groupId>
			<artifactId>fastjson</artifactId>
			<version>1.2.8</version>
		</dependency>
		<dependency>
			<groupId>org.apache.spark</groupId>
			<artifactId>spark-core_2.11</artifactId>
			<version>1.6.1</version>
			<!--<exclusions>
				<exclusion>
					<groupId>com.fasterxml.jackson.module</groupId>
					<artifactId>jackson-module-scala_2.10</artifactId>
				</exclusion>

			</exclusions>-->
		</dependency>
		<!--dependency>
			<groupId>com.fasterxml.jackson.module</groupId>
			<artifactId>jackson-module-scala_2.10</artifactId>
			<version>2.7.2</version>
		</dependency>-->
		<dependency>
			<groupId>commons-codec</groupId>
			<artifactId>commons-codec</artifactId>
			<version>1.10</version>
		</dependency>

        <dependency>
            <groupId>org.apache.camel</groupId>
            <artifactId>camel-base64</artifactId>
            <version>2.16.2</version>
        </dependency>


    </dependencies>

	<build>
		<plugins>
			<plugin>
				<groupId>org.apache.maven.plugins</groupId>
				<artifactId>maven-compiler-plugin</artifactId>
				<version>3.5.1</version>
				<configuration>
					<source>1.8</source>
					<target>1.8</target>
				</configuration>
			</plugin>

            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-dependency-plugin</artifactId>
                <configuration>
                    <outputDirectory>${project.build.directory}/lib</outputDirectory>
                    <excludeTransitive>false</excludeTransitive> <!-- 表示是否不包含间接依赖的包 -->
                    <stripVersion>false</stripVersion> <!-- 去除版本信息 -->
                </configuration>

                <executions>
                    <execution>
                        <id>copy-dependencies</id>
                        <phase>package</phase>
                        <goals>
                            <goal>copy-dependencies</goal>
                        </goals>
                        <configuration>
                            <!-- 拷贝项目依赖包到lib/目录下 -->
                            <outputDirectory>${project.build.directory}/lib</outputDirectory>
                            <excludeTransitive>false</excludeTransitive>
                            <stripVersion>false</stripVersion>
                        </configuration>
                    </execution>
                </executions>
            </plugin>

            <!-- 项目资源插件 -->
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-resources-plugin</artifactId>
                <version>2.6</version>
                <executions>
                    <execution>
                        <id>copy-resources</id>
                        <phase>package</phase>
                        <goals>
                            <goal>copy-resources</goal>
                        </goals>
                        <configuration>
                            <encoding>UTF-8</encoding>
                            <!-- 拷贝项目src/main/resources/下，除.bat以外的所有文件到conf/目录下 -->
                            <outputDirectory>${project.build.directory}/conf</outputDirectory>
                            <resources>
                                <resource>
                                    <directory>src/main/resources/</directory>
                                    <filtering>true</filtering>
                                    <excludes>
                                        <exclude>*.bat</exclude>
                                    </excludes>
                                </resource>
                            </resources>
                        </configuration>
                    </execution>
                    <execution>
                        <id>copy-command</id>
                        <phase>package</phase>
                        <goals>
                            <goal>copy-resources</goal>
                        </goals>
                        <configuration>
                            <encoding>UTF-8</encoding>
                            <!-- 只拷贝项目src/main/resources/目录下的.bat文件到输出目录下 -->
                            <outputDirectory>${project.build.directory}</outputDirectory>
                            <resources>
                                <resource>
                                    <directory>src/main/resources/</directory>
                                    <filtering>true</filtering>
                                    <includes>
                                        <include>*.bat</include>
                                    </includes>
                                </resource>
                            </resources>
                        </configuration>
                    </execution>
                </executions>
            </plugin>

            <!-- 打包插件 -->
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-jar-plugin</artifactId>
                <version>2.4</version>
                <configuration>
                    <archive>
                        <!-- 生成MANIFEST.MF的设置 -->
                        <manifest>
                            <!-- 为依赖包添加路径, 这些路径会写在MANIFEST文件的Class-Path下 -->
                            <addClasspath>true</addClasspath>
                            <classpathPrefix>lib/</classpathPrefix>
                            <!-- jar启动入口类-->
                            <mainClass>com.some.package.some.class.Main</mainClass>
                        </manifest>
                        <manifestEntries>
                            <!-- 在Class-Path下添加配置文件的路径 -->
                            <Class-Path>conf/</Class-Path>
                        </manifestEntries>
                    </archive>
                    <includes>
                        <!-- 打jar包时，只打包class文件 -->
                        <include>**/*.class</include>
                    </includes>
                </configuration>
            </plugin>
		</plugins>
	</build>
	<repositories>
		<repository>
			<id>oschina</id>
			<url>http://maven.oschina.net/content/groups/public</url>
		</repository>
		<repository>
			<id>mavenspring</id>
			<url>http://maven.springframework.org/release</url>
		</repository>
		<repository>
			<id>jcenter</id>
			<url>http://jcenter.bintray.com</url>
		</repository>
		<repository>
			<id>spring-milestones</id>
			<name>Spring Milestones</name>
			<url>http://repo.spring.io/milestone</url>
			<snapshots>
				<enabled>false</enabled>
			</snapshots>
		</repository>
		<repository>
			<id>spring-release</id>
			<name>Spring Releases</name>
			<url>http://repo.spring.io/libs-release</url>
			<snapshots>
				<enabled>false</enabled>
			</snapshots>
		</repository>
		<repository>
			<id>spring-snapshots</id>
			<name>Spring Snapshots</name>
			<url>http://repo.spring.io/snapshot</url>
			<snapshots>
				<enabled>true</enabled>
			</snapshots>
			<releases>
				<enabled>false</enabled>
			</releases>
		</repository>
	</repositories>

</project>
```

定义xd job:

```
job create --name sparkAppDemo --definition "sparkapp --mainClass=com.demo.SparkCalcDemo --appJar=/home/spark/spark-app.jar --master=local[1]" --deploy
```

加载job

```
job launch sparkAppDemo
```

然后出现以下类似的错误，主要是：`...redis:queue-inbound-channel-adapter...`错误

```
2016-03-16T10:16:10+0800 1.3.1.RELEASE INFO DeploymentSupervisor-0 zk.ZKJobDeploymentHandler - Deployment status for job 'sparkAppDemo': DeploymentStatus{state=deployed}
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark application 'com.demo.SparkCalcDemo' finished with exit code: 1
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: Exception in thread "main" java.lang.SecurityException: Invalid signature file digest for Manifest main attributes
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.security.util.SignatureFileVerifier.processImpl(SignatureFileVerifier.java:284)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.security.util.SignatureFileVerifier.process(SignatureFileVerifier.java:238)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.util.jar.JarVerifier.processEntry(JarVerifier.java:316)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.util.jar.JarVerifier.update(JarVerifier.java:228)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.util.jar.JarFile.initializeVerifier(JarFile.java:383)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.util.jar.JarFile.getInputStream(JarFile.java:450)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.JarIndex.getJarIndex(JarIndex.java:137)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$JarLoader$1.run(URLClassPath.java:839)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$JarLoader$1.run(URLClassPath.java:831)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.security.AccessController.doPrivileged(Native Method)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$JarLoader.ensureOpen(URLClassPath.java:830)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$JarLoader.<init>(URLClassPath.java:803)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$3.run(URLClassPath.java:530)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath$3.run(URLClassPath.java:520)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.security.AccessController.doPrivileged(Native Method)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath.getLoader(URLClassPath.java:519)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath.getLoader(URLClassPath.java:492)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath.getNextLoader(URLClassPath.java:457)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at sun.misc.URLClassPath.getResource(URLClassPath.java:211)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.net.URLClassLoader$1.run(URLClassLoader.java:365)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.net.URLClassLoader$1.run(URLClassLoader.java:362)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.security.AccessController.doPrivileged(Native Method)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.net.URLClassLoader.findClass(URLClassLoader.java:361)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.lang.ClassLoader.loadClass(ClassLoader.java:424)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.lang.ClassLoader.loadClass(ClassLoader.java:357)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.lang.Class.forName0(Native Method)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at java.lang.Class.forName(Class.java:348)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at org.apache.spark.deploy.SparkSubmit$.org$apache$spark$deploy$SparkSubmit$$runMain(SparkSubmit.scala:538)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at org.apache.spark.deploy.SparkSubmit$.doRunMain$1(SparkSubmit.scala:166)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at org.apache.spark.deploy.SparkSubmit$.submit(SparkSubmit.scala:189)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at org.apache.spark.deploy.SparkSubmit$.main(SparkSubmit.scala:110)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: 	at org.apache.spark.deploy.SparkSubmit.main(SparkSubmit.scala)
2016-03-16T10:16:14+0800 1.3.1.RELEASE ERROR inbound.job:sparkAppDemo-redis:queue-inbound-channel-adapter1 tasklet.SparkTasklet - Spark Logger: Using Spark's default log4j profile: org/apache/spark/log4j-defaults.properties
```

仔细看的话，在上面有一句`Exception in thread "main" java.lang.SecurityException: Invalid signature file digest for Manifest main attributes`错误，这个错误是由于导出的jar包结构信息不正确导致的。用eclipse的导出runnable jar 功能导出的jar包就没问题了。

另外如果有依赖的jar包没哟被加载进去，则会在最上方出现`java.lang.NoClassDefFoundError: `类似信息。

相关资料：[以分布式方式运行Spring-XD](http://www.todu.top/2016/03/09/以分布式方式运行Spring-XD)
