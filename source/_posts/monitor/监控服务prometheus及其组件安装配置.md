---
title: prometheus机器组件安装配置
category: 监控
tags:
  - monitor
  - 监控
  - Prometheus
  - alert
abbrlink: d995
date: 2019-12-30 11:54:55
updateDate: 2019-12-30 11:54:55
---

# Prometheus 安装

准备：

```bash
mkdir -p /data/logs /data/soft/monitor
```

## 安装核心服务 Prometheus

下载 [Prometheus](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/prometheus-2.14.0.linux-amd64.tar.gz) 到 /data/soft 目录,并解压到 /data/soft/monitor/prometheus

执行命令为:

```bash
mkdir -p /data/soft/monitor && mkdir -p /data/logs && cd /data/soft

wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/prometheus-2.14.0.linux-amd64.tar.gz && tar -xvf prometheus-2.14.0.linux-amd64.tar.gz && mv prometheus-2.14.0.linux-amd64 monitor/prometheus
```

启动

```bash
cd /data/soft/monitor/prometheus
setsid ./prometheus --config.file="prometheus.yml" --web.enable-lifecycle >> /data/logs/prometheus.log 2>&1
```

上面使用默认配置启动了Prometheus，至于配置文件下面会再进行修改。

检查是否启动成功，输入命令 `curl -XGET -sI localhost:9090/metrics|head -n1` 如果输出 `HTTP/1.1 200 OK` 则代表成功，否则查看日志 `/data/logs/prometheus.log` 定位问题。

## 安装Linux监控组件

该组件需要安装在被监控的主机上，如果主机有多台，那么应该分别进行安装并启动。

下载 [node_export](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/node_exporter-0.18.1.linux-amd64.tar.gz) 到 目录 /data/soft/monitor 并解压，命令：

```bash
cd /data/soft
wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/node_exporter-0.18.1.linux-amd64.tar.gz

tar -xvf node_exporter-0.18.1.linux-amd64.tar.gz && mv node_exporter-0.18.1.linux-amd64 monitor/node_exporter
```

启动 node_exporter

```bash
cd /data/soft/monitor/node_exporter

setsid ./node_exporter >> /data/logs/node_exporter.log 2>&1
```

检查是否启动成功, `curl -XGET -sI localhost:9100/metrics|head -n1` 如果输出 `HTTP/1.1 200 OK` 则启动成功，否则查看日志 `/data/logs/node_exporter.log` 定位问题。

## 安装 redis 监控组件

下载 [redis-exporter](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/redis_exporter-v1.3.4.linux-amd64.tar.gz) 到目录 `/data/soft` 并解压，命令：

```bash
cd /data/soft
wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/redis_exporter-v1.3.4.linux-amd64.tar.gz

tar -xvf redis_exporter-v1.3.4.linux-amd64.tar.gz && mv redis_exporter-v1.3.4.linux-amd64 monitor/redis_exporter
```

运行 redis_exporter

```bash
cd /data/soft/monitor/redis_exporter
setsid ./redis_exporter -redis.addr='redis://192.168.25.150:6379' --check-keys=db1=stream*,db0=stream* >> /data/logs/redis_exporter.log 2>&1
```

**注意** `--check-keys` 参数是要监控的 redis 对应 db 的 key名字，比如要监控 `db0` 中的 `stream` 开头的key，则加入参数 `--check-keys='db1=stream*'`，多个用英文逗号隔开，这里需要根据 iot 服务部署时候设置的redis db 设置，默认是1。

`-redis.addr` 是redis 地址。

检查是否启动成功, `curl -XGET -sI localhost:9121/metrics|head -n1` 如果输出 `HTTP/1.1 200 OK` 则启动成功，否则查看日志 `/data/logs/redis_exporter.log` 定位问题。

## 安装大盘显示组件 Grafana

下载 [grafana](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/grafana-6.4.4.linux-amd64.tar.gz) 到目录 `/data/soft` 并解压，命令：

```bash
cd /data/soft
wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/grafana-6.4.4.linux-amd64.tar.gz

tar -xvf grafana-6.4.4.linux-amd64.tar.gz && mv grafana-6.4.4 monitor/grafana
```

启动：

```bash
cd /data/soft/monitor/grafana
setsid bin/grafana-server -config conf/defaults.ini >> /data/logs/grafana.log 2>&1
```

检查是否启动成功, `curl -XGET -sI localhost:3000/login|head -n1` 如果输出 `HTTP/1.1 200 OK` 则启动成功，否则查看日志 `/data/logs/grafana.log` 定位问题。

## 安装告警组件 AlertManager

下载 [alertmanager](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/alertmanager-0.19.0.linux-amd64.tar.gz) 到目录 `/data/soft` 并解压，命令：

```bash
cd /data/soft
wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/alertmanager-0.19.0.linux-amd64.tar.gz

tar -xvf alertmanager-0.19.0.linux-amd64.tar.gz && mv alertmanager-0.19.0.linux-amd64 monitor/alertmanager
```

启动：

```bash
cd /data/soft/monitor/alertmanager
setsid ./alertmanager --config.file="alertmanager.yml" >> /data/logs/alertmanager.log 2>&1
```

检查是否启动成功, `curl -XGET -sI http://localhost:9093|head -n1` 如果输出 `HTTP/1.1 200 OK` 则启动成功，否则查看日志 `/data/logs/grafana.log` 定位问题。

## 安装 健康检查附加组件

该组件可以针对 http，tcp 等端口进行检查，从而获得服务运行状态并进行告警。

下载 [black_exporter](https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/blackbox_exporter-0.16.0.linux-amd64.tar.gz) 到目录 `/data/soft` 并解压，命令：

```bash
cd /data/soft
wget https://hekr-files.oss-cn-shanghai.aliyuncs.com/soft/prometheus/linux/blackbox_exporter-0.16.0.linux-amd64.tar.gz

tar -xvf blackbox_exporter-0.16.0.linux-amd64.tar.gz && mv blackbox_exporter-0.16.0.linux-amd64 monitor/blackbox_exporter
```

启动：

```bash
cd /data/soft/monitor/blackbox_exporter
setsid ./blackbox_exporter --config.file="blackbox.yml" >> /data/logs/blackbox_exporter.log 2>&1
```

检查是否启动成功, `curl -XGET -sI http://localhost:9115 |head -n1` 如果输出 `HTTP/1.1 200 OK` 则启动成功，否则查看日志 `/data/logs/blackbox_exporter.log` 定位问题。

### 配置检查项

将如下配置文件写入（覆盖） `balckbox.yml` :

```yml
modules:
  http_spring:
    prober: http
    timeout: 3s
    http:
      fail_if_body_not_matches_regexp:
        - '"status":"UP"'
  http_2xx:
    prober: http
    timeout: 3s
  tcp_connect:
    prober: tcp
    timeout: 3s
```

使用命令： `curl -XPOST  http://localhost:9115/-/reload` 使其重新加载生效。

## 检查服务是否都已经启动

检查已经启动服务： `ps -ef|grep -v grep|egrep "prometheus|grafana|exporter|alarmmanager"`。

## 增加监控配置

修改 `/data/soft/monitor/prometheus/prometheus.yml`，替换为如下内容，注意 **&#123;&#123;&#125;&#125;** 内的内容需要动态替换为对应内容。

```yml
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - localhost:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  - "rules/*.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:9090']

  - job_name: redis_exporter
    static_configs:
    - targets: ['{{localhost:9121}}']

  - job_name: iot_cloud_os_core
    metrics_path: /actuator/prometheus
    static_configs:
    - targets: ['{{ core ip:port }}']

  - job_name: linux_host
    static_configs:
    - targets: ['{{ localhost:9100 }}']

  - job_name: emqx
    static_configs:
    - targets: ['{{localhost:9540}}']

  # iot_cloud_os_core 健康检查
  - job_name: health_iot_cloud_os_core
    metrics_path: /probe
    params:
      module: [http_spring]
    static_configs:
      - targets:
        - '{{localhost:8001}}/actuator/health'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: "{{192.168.25.192:9115}}"

  # redis 端口检查
  - job_name: health_redis_port
    metrics_path: /probe
    params:
      module: [tcp_connect]
    static_configs:
      - targets:
        - '{{localhost:6379}}'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: "{{192.168.25.192:9115}}"

  # mongodb 端口检查
  - job_name: health_mongodb_port
    metrics_path: /probe
    params:
      module: [tcp_connect]
    static_configs:
      - targets:
        - '{{localhost:27017}}'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: "{{192.168.25.192:9115}}"

  # rocketmq 端口检查
  - job_name: health_mongodb_port
    metrics_path: /probe
    params:
      module: [tcp_connect]
    static_configs:
      - targets:
        - '{{localhost:9876}}'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: "{{192.168.25.192:9115}}"
```

保存后输入 `curl -XPOST localhost:9090/-/reload` 使Prometheus重新加载配置，如果没有输出则正确。

## 添加告警配置

创建目录 `mkdir -p /data/soft/monitor/prometheus/rules` 并将以下文件内容保存 `/data/soft/monitor/prometheus/rules` 目录下：

health.yml

```yml
groups:
- name: 监控检查
  rules:
  - alert: "健康检查失败"
    expr: probe_success!=1
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "健康检查失败 {{$labels.job}}  {{$labels.instance}}"
      description: "健康检查失败 {{$labels.job}}  {{$labels.instance}}"
```
