---
title: Prometheus grafana
layout: post
category: golang
author: 夏泽民
---
https://github.com/prometheus/client_golang
https://github.com/grpc-ecosystem/go-grpc-prometheus

Server-side
    // After all your registrations, make sure all of the Prometheus metrics are initialized.
    grpc_prometheus.Register(myServer)
    // Register Prometheus metrics handler.    
    http.Handle("/metrics", promhttp.Handler())
    
Client-side
 grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
 
 https://github.com/zsais/go-gin-prometheus
 p := ginprometheus.NewPrometheus("gin")
 
 为了方便client library的使用提供了四种数据类型： Counter, Gauge, Histogram, Summary, 简单理解就是Counter对数据只增不减，Gauage可增可减，Histogram,Summary提供跟多的统计信息。
<!-- more -->
https://www.cnblogs.com/gaorong/p/7881203.html

https://skyingzz.github.io/2020/01/19/prometheus-client-go/

https://blog.csdn.net/huosenbulusi/article/details/107529961


下载解压
在官网 下在最新版本压缩包：prometheus-2.3.0.darwin-amd64.tar.gz

解压并进入该目录：

tar xvfz prometheus-*.tar.gz
cd prometheus-*
修改配置文件
本文将配置文件修改为检测Prometheus自身：
修改配置文件为：


global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']

.访问web
http://190.168.3.250:9090访问自己的状态页面
可以通过访问localhost:9090验证Prometheus自身的指标：190.168.3.250:9090/metrics

1.配置Prometheus监控本身
全局配置文件简介
有关配置选项的完整，请参阅：https://prometheus.io/docs/prometheus/latest/configuration/configuration/
Prometheus以scrape_interval规则周期性从监控目标上收集数据，然后将数据存储到本地存储上。scrape_interval可以设定全局也可以设定单个metrics。
Prometheus以evaluation_interval规则周期性对告警规则做计算，然后更新告警状态。evaluation_interval只有设定在全局。

global：全局配置
alerting：告警配置
rule_files：告警规则
scrape_configs：配置数据源，称为target，每个target用job_name命名。又分为静态配置和服务发现

global:
默认抓取周期，可用单位ms、smhdwy #设置每15s采集数据一次，默认1分钟
[ scrape_interval: <duration> | default = 1m ]
默认抓取超时
[ scrape_timeout: <duration> | default = 10s ]
估算规则的默认周期 # 每15秒计算一次规则。默认1分钟
[ evaluation_interval: <duration> | default = 1m ]
和外部系统（例如AlertManager）通信时为时间序列或者警情（Alert）强制添加的标签列表
external_labels:
[ <labelname>: <labelvalue> ... ]

规则文件列表
rule_files:
[ - <filepath_glob> ... ]

抓取配置列表
scrape_configs:
[ - <scrape_config> ... ]

Alertmanager相关配置
alerting:
alert_relabel_configs:
[ - <relabel_config> ... ]
alertmanagers:
[ - <alertmanager_config> ... ]

远程读写特性相关的配置
remote_write:
[ - <remote_write> ... ]
remote_read:
[ - <remote_read> ... ]

vi prometheus.yml
下面就是拉取自身服务采样点数据配置
scrape_configs:
别监控指标，job名称会增加到拉取到的所有采样点上，同时还有一个instance目标服务的host：port标签也会增加到采样点上

job_name: 'prometheus'
覆盖global的采样点，拉取时间间隔5s
scrape_interval: 5s
static_configs:
targets: ['localhost:9090']
最下面，静态配置监控本机，采集本机9090端口数据
https://www.infoq.cn/article/uj12knworcwg0kke8zfv

Grafana安装
安装步骤参考grafana官方文档：http://docs.grafana.org/installation/mac/

使用homebrew安装：

brew update
brew install grafana
安装成功后你可以使用默认配置启动：

brew services start grafana
登录
使用默认账号密码 ： admin/admin 登录

https://www.jianshu.com/p/7954823d6e65
https://blog.csdn.net/Holly_walker/article/details/103820509

https://blog.csdn.net/Kammingo/article/details/105225838
https://segmentfault.com/a/1190000018372409
https://cloud.tencent.com/developer/article/1744617

https://www.yisu.com/zixun/7612.html

https://dbaplus.cn/news-134-3247-1.html
https://zhuanlan.zhihu.com/p/117719823
https://www.bladewan.com/2020/01/03/prometheus_1/


https://www.cnblogs.com/netonline/p/8289411.html