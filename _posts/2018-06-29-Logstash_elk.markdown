---
title: Logstash_elk
layout: post
category: elasticsearch
author: 夏泽民
---
 Logstash是一个开源的服务器端数据处理管道，可以同时从多个源获取数据。面对海量的日志量，rsyslog和sed，awk等日志收集，处理工具已经显的力不从心。logstash是一个整合型的框架，可以用以日志的收集，存储，索引构建（一般这个功能被ES取代）。
<!-- more -->
<img src="{{site.url}}{{site.baseurl}}/img/Logstash_elk.jpeg"/>
 logstash 的服务器端从redis/kafka/rabbitmq等（broker）消息队列获取数据。一条数据一条数据的清洗。清洗完成后发送给elasticsearch集群。再在kibana上显示。

  对于logstash而言，他的所有功能都是基于插件来完成的。input,filter,output等等都是这样。

输入插件：
    提取所有的数据。 数据通常以多种格式散布或分布在各个系统上。logstash支持各种输入。可以从常见的各类源中获取事件。轻松从日志，指标，web应用程序和各种AWS服务中进行采集，所有的这些都是以连续流的方式进行。
 过滤器：
    随着数据传输而来。logstash过滤器解析每一个事件，识别命名字段以构建结构。将它们转换为通用格式，从而更轻松，更快捷的分析。常用的过滤数据手段。（1）、用grok从非结构化数据导出结构。（2）、从IP地址解读地理坐标。（3）、匿名PII数据，完全排除敏感字段。 （4）、简化数据源，格式或模式。
 产出：
  当然，elasticsearch是首选的输出对象，但是还有很多的选项。可根据需要路由。

 插件介绍
      常见的插件。在诸多的input插件当中有file，udp，http，kafka，rabbitmq，beats等等
   file：
    文件流事件。类似与 tail -n 1 -f 开始阅读。当然也可以从头开始阅读。利用了sincedb记录文件的状态，包括文件的inode号，主设备号，从设备号，文件内字节偏移量。所以文件不能改名字。一改名字就会使用新的sincedb。
 input {
        file {
                path => ["/var/log/messages"]
                type => "system"
                start_position => "beginning"
             }
}
 udp：
    通过udp将消息作为事件通过网络读取。唯一需要的配置项是port，它指定udp端口logstash将监听事件流。
演示一下。
    先安装 collectd（epel源中）    yum -y install collectd      collectd是一个性能监控程序。
编辑配置文件 /etc/collectd.conf
将想要监控的内容取消注释。 （一定要取消 network）
 Hostname    "node-1"
LoadPlugin syslog
LoadPlugin battery
LoadPlugin cpu
LoadPlugin df
LoadPlugin disk
LoadPlugin interface
LoadPlugin load
LoadPlugin memory
LoadPlugin network
<Plugin network>
	<Server "192.168.40.133" "25826"> #192.168.40.133 是logstash监听的地址，25826 是logstash监听端口
 	</Server>
</Plugin>
Include "/etc/collectd.d"

启动即可。
systemctl  start   collectd
接下来写一个udp.conf 
 input {
        udp {
                port => 25826
                codec => collectd {}
                type  => "collectd"
        }
}
 
output {
        stdout {
                codec   => rubydebug
        }
}
执行命令。开始收集。
# logstash -f udp.conf

Filter插件：
    用于在将event发往output之前，对其实现一定的处理功能。

grok：
    用于分析并结构化文本数据。目前是logstash中非结构化数据转化为结构化数据的不二之选。可处理 syslog，Apache，nginx格式的日志。
 Filebeat：
     Filebeat是ELK的一部分。可以使Logstash，Elasticsearch和Kibana无缝协作。filebeat的功能是转发和汇总日志与文件。  filebeat可以读取并转发日志行，如果出现中断，还会在一切恢复正常后，从中断的位置继续开始。
  
  beat插件（这是一个input插件）
input {
  beats {
    port => 5044
  }
}
 
output {
  elasticsearch {
    hosts => "localhost:9200"
    manage_template => false
    index => "%{[@metadata][beat]}-%{+YYYY.MM.dd}"
    document_type => "%{[@metadata][type]}"
  }
}
output插件：
     这一般是把数据存储下来的，有email，csv。当然还有大名鼎鼎的 Elasticsearch