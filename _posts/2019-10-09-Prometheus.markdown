---
title: Prometheus
layout: post
category: golang
author: 夏泽民
---
https://github.com/prometheus/prometheus
Prometheus是由SoundCloud开发的开源监控报警系统和时序列数据库(TSDB)。Prometheus使用Go语言开发，是Google BorgMon监控系统的开源版本。
2016年由Google发起Linux基金会旗下的原生云基金会(Cloud Native Computing Foundation), 将Prometheus纳入其下第二大开源项目。
Prometheus目前在开源社区相当活跃。
Prometheus和Heapster(Heapster是K8S的一个子项目，用于获取集群的性能数据。)相比功能更完善、更全面。Prometheus性能也足够支撑上万台规模的集群。
<!-- more -->
Prometheus的特点
多维度数据模型。
灵活的查询语言。
不依赖分布式存储，单个服务器节点是自主的。
通过基于HTTP的pull方式采集时序数据。
可以通过中间网关进行时序列数据推送。
通过服务发现或者静态配置来发现目标服务对象。
支持多种多样的图表和界面展示，比如Grafana等。
官网地址：https://prometheus.io/
<img src="{{site.url}}{{site.baseurl}}/img/prometheus.png"/>

基本原理
Prometheus的基本原理是通过HTTP协议周期性抓取被监控组件的状态，任意组件只要提供对应的HTTP接口就可以接入监控。不需要任何SDK或者其他的集成过程。这样做非常适合做虚拟化环境监控系统，比如VM、Docker、Kubernetes等。输出被监控组件信息的HTTP接口被叫做exporter 。目前互联网公司常用的组件大部分都有exporter可以直接使用，比如Varnish、Haproxy、Nginx、MySQL、Linux系统信息(包括磁盘、内存、CPU、网络等等)。

服务过程
Prometheus Daemon负责定时去目标上抓取metrics(指标)数据，每个抓取目标需要暴露一个http服务的接口给它定时抓取。Prometheus支持通过配置文件、文本文件、Zookeeper、Consul、DNS SRV Lookup等方式指定抓取目标。Prometheus采用PULL的方式进行监控，即服务器可以直接通过目标PULL数据或者间接地通过中间网关来Push数据。
Prometheus在本地存储抓取的所有数据，并通过一定规则进行清理和整理数据，并把得到的结果存储到新的时间序列中。
Prometheus通过PromQL和其他API可视化地展示收集的数据。Prometheus支持很多方式的图表可视化，例如Grafana、自带的Promdash以及自身提供的模版引擎等等。Prometheus还提供HTTP API的查询方式，自定义所需要的输出。
PushGateway支持Client主动推送metrics到PushGateway，而Prometheus只是定时去Gateway上抓取数据。
Alertmanager是独立于Prometheus的一个组件，可以支持Prometheus的查询语句，提供十分灵活的报警方式。
三大套件
Server 主要负责数据采集和存储，提供PromQL查询语言的支持。
Alertmanager 警告管理器，用来进行报警。
Push Gateway 支持临时性Job主动推送指标的中间网关。

三.安装pushgateway
pushgateway是为了允许临时作业和批处理作业向普罗米修斯公开他们的指标。
由于这类作业的存在时间可能不够长, 无法抓取到, 因此它们可以将指标推送到推网关中。
Prometheus采集数据是用的pull也就是拉模型，这从我们刚才设置的5秒参数就能看出来。但是有些数据并不适合采用这样的方式，对这样的数据可以使用Push Gateway服务。
它就相当于一个缓存，当数据采集完成之后，就上传到这里，由Prometheus稍后再pull过来。
我们来试一下，首先启动Push Gateway

mkdir -p /home/chenqionghe/promethues/pushgateway
cd !$
docker run -d -p 9091:9091 --name pushgateway prom/pushgateway
访问http://10.211.55.25:9091 可以看到pushgateway已经运行起来了

四.安装Grafana展示
Grafana是用于可视化大型测量数据的开源程序，它提供了强大和优雅的方式去创建、共享、浏览数据。
Dashboard中显示了你不同metric数据源中的数据。
Grafana最常用于因特网基础设施和应用分析，但在其他领域也有用到，比如：工业传感器、家庭自动化、过程控制等等。
Grafana支持热插拔控制面板和可扩展的数据源，目前已经支持Graphite、InfluxDB、OpenTSDB、Elasticsearch、Prometheus等。

我们使用docker安装

docker run -d -p 3000:3000 --name grafana grafana/grafana

五.安装AlterManager
Pormetheus的警告由独立的两部分组成。
Prometheus服务中的警告规则发送警告到Alertmanager。
然后这个Alertmanager管理这些警告。包括silencing, inhibition, aggregation，以及通过一些方法发送通知，例如：email，PagerDuty和HipChat。
建立警告和通知的主要步骤：

创建和配置Alertmanager
启动Prometheus服务时，通过-alertmanager.url标志配置Alermanager地址，以便Prometheus服务能和Alertmanager建立连接。
创建和配置Alertmanager

mkdir -p /home/chenqionghe/promethues/alertmanager
cd !$

Prometheus 作为生态圈 Cloud Native Computing Foundation（简称：CNCF）中的重要一员,其活跃度仅次于 Kubernetes, 现已广泛用于 Kubernetes 集群的监控系统中。本文将简要介绍 Prometheus 的组成和相关概念，并实例演示 Prometheus 的安装，配置及使用，以便开发人员和云平台运维人员可以快速的掌握 Prometheus。


作为新一代的监控框架，Prometheus 具有以下特点：

强大的多维度数据模型：
时间序列数据通过 metric 名和键值对来区分。
所有的 metrics 都可以设置任意的多维标签。
数据模型更随意，不需要刻意设置为以点分隔的字符串。
可以对数据模型进行聚合，切割和切片操作。
支持双精度浮点类型，标签可以设为全 unicode。
灵活而强大的查询语句（PromQL）：在同一个查询语句，可以对多个 metrics 进行乘法、加法、连接、取分数位等操作。
易于管理： Prometheus server 是一个单独的二进制文件，可直接在本地工作，不依赖于分布式存储。
高效：平均每个采样点仅占 3.5 bytes，且一个 Prometheus server 可以处理数百万的 metrics。
使用 pull 模式采集时间序列数据，这样不仅有利于本机测试而且可以避免有问题的服务器推送坏的 metrics。
可以采用 push gateway 的方式把时间序列数据推送至 Prometheus server 端。
可以通过服务发现或者静态配置去获取监控的 targets。
有多种可视化图形界面。
易于伸缩。
需要指出的是，由于数据采集可能会有丢失，所以 Prometheus 不适用对采集数据要 100% 准确的情形。但如果用于记录时间序列数据，Prometheus 具有很大的查询优势，此外，Prometheus 适用于微服务的体系架构。

Prometheus 组成及架构
Prometheus 生态圈中包含了多个组件，其中许多组件是可选的：

Prometheus Server: 用于收集和存储时间序列数据。
Client Library: 客户端库，为需要监控的服务生成相应的 metrics 并暴露给 Prometheus server。当 Prometheus server 来 pull 时，直接返回实时状态的 metrics。
Push Gateway: 主要用于短期的 jobs。由于这类 jobs 存在时间较短，可能在 Prometheus 来 pull 之前就消失了。为此，这次 jobs 可以直接向 Prometheus server 端推送它们的 metrics。这种方式主要用于服务层面的 metrics，对于机器层面的 metrices，需要使用 node exporter。
Exporters: 用于暴露已有的第三方服务的 metrics 给 Prometheus。
Alertmanager: 从 Prometheus server 端接收到 alerts 后，会进行去除重复数据，分组，并路由到对收的接受方式，发出报警。常见的接收方式有：电子邮件，pagerduty，OpsGenie, webhook 等。
一些其他的工具。

其大概的工作流程是：

Prometheus server 定期从配置好的 jobs 或者 exporters 中拉 metrics，或者接收来自 Pushgateway 发过来的 metrics，或者从其他的 Prometheus server 中拉 metrics。
Prometheus server 在本地存储收集到的 metrics，并运行已定义好的 alert.rules，记录新的时间序列或者向 Alertmanager 推送警报。
Alertmanager 根据配置文件，对接收到的警报进行处理，发出告警。
在图形界面中，可视化采集数据。

数据模型

Prometheus 中存储的数据为时间序列，是由 metric 的名字和一系列的标签（键值对）唯一标识的，不同的标签则代表不同的时间序列。

metric 名字：该名字应该具有语义，一般用于表示 metric 的功能，例如：http_requests_total, 表示 http 请求的总数。其中，metric 名字由 ASCII 字符，数字，下划线，以及冒号组成，且必须满足正则表达式 [a-zA-Z_:][a-zA-Z0-9_:]*。
标签：使同一个时间序列有了不同维度的识别。例如 http_requests_total{method="Get"} 表示所有 http 请求中的 Get 请求。当 method="post" 时，则为新的一个 metric。标签中的键由 ASCII 字符，数字，以及下划线组成，且必须满足正则表达式 [a-zA-Z_:][a-zA-Z0-9_:]*。
样本：实际的时间序列，每个序列包括一个 float64 的值和一个毫秒级的时间戳。
格式：<metric name>{<label name>=<label value>, …}，例如：http_requests_total{method="POST",endpoint="/api/tracks"}。
四种 Metric 类型

Prometheus 客户端库主要提供四种主要的 metric 类型：

Counter

一种累加的 metric，典型的应用如：请求的个数，结束的任务数， 出现的错误数等等。
例如，查询 http_requests_total{method="get", job="Prometheus", handler="query"} 返回 8，10 秒后，再次查询，则返回 14。

Gauge

一种常规的 metric，典型的应用如：温度，运行的 goroutines 的个数。
可以任意加减。
例如：go_goroutines{instance="172.17.0.2", job="Prometheus"} 返回值 147，10 秒后返回 124。

Histogram

可以理解为柱状图，典型的应用如：请求持续时间，响应大小。
可以对观察结果采样，分组及统计。

Summary

类似于 Histogram, 典型的应用如：请求持续时间，响应大小。
提供观测值的 count 和 sum 功能。
提供百分位的功能，即可以按百分比划分跟踪结果。
instance 和 jobs

instance: 一个单独 scrape 的目标， 一般对应于一个进程。

jobs: 一组同种类型的 instances（主要用于保证可扩展性和可靠性）


什么是TSDB？
TSDB(Time Series Database)时序列数据库，我们可以简单的理解为一个优化后用来处理时间序列数据的软件，并且数据中的数组是由时间进行索引的。

时间序列数据库的特点

大部分时间都是写入操作。
写入操作几乎是顺序添加，大多数时候数据到达后都以时间排序。
写操作很少写入很久之前的数据，也很少更新数据。大多数情况在数据被采集到数秒或者数分钟后就会被写入数据库。
删除操作一般为区块删除，选定开始的历史时间并指定后续的区块。很少单独删除某个时间或者分开的随机时间的数据。
基本数据大，一般超过内存大小。一般选取的只是其一小部分且没有规律，缓存几乎不起任何作用。
读操作是十分典型的升序或者降序的顺序读。
高并发的读操作十分常见。
常见的时间序列数据库

TSDB项目	官网
influxDB	https://influxdata.com/
RRDtool	http://oss.oetiker.ch/rrdtool/
Graphite	http://graphiteapp.org/
OpenTSDB	http://opentsdb.net/
Kdb+	http://kx.com/
Druid	http://druid.io/
KairosDB	http://kairosdb.github.io/
Prometheus	https://prometheus.io/
什么是Prometheus?
Prometheus是由SoundCloud开发的开源监控报警系统和时序列数据库(TSDB)。Prometheus使用Go语言开发，是Google BorgMon监控系统的开源版本。

2016年由Google发起Linux基金会旗下的原生云基金会(Cloud Native Computing Foundation), 将Prometheus纳入其下第二大开源项目。Prometheus目前在开源社区相当活跃。

Prometheus和Heapster(Heapster是K8S的一个子项目，用于获取集群的性能数据。)相比功能更完善、更全面。Prometheus性能也足够支撑上万台规模的集群。

Prometheus的特点
多维度数据模型。
灵活的查询语言。
不依赖分布式存储，单个服务器节点是自主的。
通过基于HTTP的pull方式采集时序数据。
可以通过中间网关进行时序列数据推送。
通过服务发现或者静态配置来发现目标服务对象。
支持多种多样的图表和界面展示，比如Grafana等。
Prometheus相关组件
Prometheus生态系统由多个组件组成，它们中的一些是可选的。多数Prometheus组件是Go语言写的，这使得这些组件很容易编译和部署。

Prometheus Server
主要负责数据采集和存储，提供PromQL查询语言的支持。

客户端SDK
官方提供的客户端类库有go、java、scala、python、ruby，其他还有很多第三方开发的类库，支持nodejs、php、erlang等。

Push Gateway
支持临时性Job主动推送指标的中间网关。

PromDash
使用Rails开发可视化的Dashboard，用于可视化指标数据。

Exporter
Exporter是Prometheus的一类数据采集组件的总称。它负责从目标处搜集数据，并将其转化为Prometheus支持的格式。与传统的数据采集组件不同的是，它并不向中央服务器发送数据，而是等待中央服务器主动前来抓取。

Prometheus提供多种类型的Exporter用于采集各种不同服务的运行状态。目前支持的有数据库、硬件、消息中间件、存储系统、HTTP服务器、JMX等。

alertmanager
警告管理器，用来进行报警。

prometheus_cli
命令行工具。

其他辅助性工具
多种导出工具，可以支持Prometheus存储数据转化为HAProxy、StatsD、Graphite等工具所需要的数据存储格式。

Prometheus的基本原理是通过HTTP协议周期性抓取被监控组件的状态，任意组件只要提供对应的HTTP接口就可以接入监控。不需要任何SDK或者其他的集成过程。这样做非常适合做虚拟化环境监控系统，比如VM、Docker、Kubernetes等。输出被监控组件信息的HTTP接口被叫做exporter 。目前互联网公司常用的组件大部分都有exporter可以直接使用，比如Varnish、Haproxy、Nginx、MySQL、Linux系统信息(包括磁盘、内存、CPU、网络等等)。

Prometheus服务过程大概是这样：

Prometheus Daemon负责定时去目标上抓取metrics(指标)数据，每个抓取目标需要暴露一个http服务的接口给它定时抓取。Prometheus支持通过配置文件、文本文件、Zookeeper、Consul、DNS SRV Lookup等方式指定抓取目标。Prometheus采用PULL的方式进行监控，即服务器可以直接通过目标PULL数据或者间接地通过中间网关来Push数据。

Prometheus在本地存储抓取的所有数据，并通过一定规则进行清理和整理数据，并把得到的结果存储到新的时间序列中。

Prometheus通过PromQL和其他API可视化地展示收集的数据。Prometheus支持很多方式的图表可视化，例如Grafana、自带的Promdash以及自身提供的模版引擎等等。Prometheus还提供HTTP API的查询方式，自定义所需要的输出。

PushGateway支持Client主动推送metrics到PushGateway，而Prometheus只是定时去Gateway上抓取数据。

Alertmanager是独立于Prometheus的一个组件，可以支持Prometheus的查询语句，提供十分灵活的报警方式。

Prometheus适用的场景
Prometheus在记录纯数字时间序列方面表现非常好。它既适用于面向服务器等硬件指标的监控，也适用于高动态的面向服务架构的监控。对于现在流行的微服务，Prometheus的多维度数据收集和数据筛选查询语言也是非常的强大。Prometheus是为服务的可靠性而设计的，当服务出现故障时，它可以使你快速定位和诊断问题。它的搭建过程对硬件和服务没有很强的依赖关系。

Prometheus不适用的场景
Prometheus它的价值在于可靠性，甚至在很恶劣的环境下，你都可以随时访问它和查看系统服务各种指标的统计信息。 如果你对统计数据需要100%的精确，它并不适用，例如：它不适用于实时计费系统。

Prometheus官网：https://prometheus.io/

安装Prometheus
Prometheus官方给出了多重部署方案，比如：Docker容器、Ansible、Chef、Puppet、Saltstack等。

Prometheus用Golang实现，因此具有天然可移植性(支持Linux、Windows、macOS和Freebsd)。这里直接使用预编译的二进制文件部署，开箱即用。


