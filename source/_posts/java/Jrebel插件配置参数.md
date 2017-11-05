---
title: Jrebel插件配置参数
tags:
  - jrebel
  - java
category: java
abbrlink: 15199
date: 2016-03-09 14:10:18
---
# JVM 参数
`-javaagent:/path/jrebel.jar `
# Spring
## Spring Bean/Core/MVC/Security/Webflow/WS
`-Drebel.spring_plugin=true`
## Spring Data
`-Drebel.spring_data_plugin=true`
## Struts
`-Drebel.struts2_plugin=true`
# Hibernate
`-Drebel.hibernate_plugin=true`
# Hibernate Validator
`-Drebel.hibernate_validator_plugin=true`
# MyBatis
`-Drebel.mybatis_plugin=true`
# Logback
`-Drebel.logback_plugin=true`
# Log4J 2
`-Drebel.log4j2_plugin=true`
#Groovy
`-Drebel.groovy_plugin=true`
# Jruby
`-Drebel.jruby_plugin=true`
# GWT
`-Drebel.gwt_plugin=true`
# 参考
[plugins](http://manuals.zeroturnaround.com/jrebel/misc/frameworks.html)
