---
title: logstash输出到elasticsearch
layout: post
category: elasticsearch
author: 夏泽民
---
<!-- more -->
output配置中elasticsearch配置
action
index 给一个文档建立索引
delete 通过id值删除一个文档（这个action需要指定一个id值）
create 插入一条文档信息，如果这条文档信息在索引中已经存在，那么本次插入工作失败
update 通过id值更新一个文档。更新有个特殊的案例upsert，如果被更新的文档还不存在，那么就会用到upsert
示例：

action => "index"

index
写入事件所用的索引。可以动态的使用%{foo}语法，它的默认值是： 
“logstash-%{+YYYY.MM.dd}”，以天为单位分割的索引，使你可以很容易的删除老的数据或者搜索指定时间范围内的数据。

索引不能包含大写字母。推荐使用以周为索引的ISO 8601格式，例如logstash-%{+xxxx.ww}

示例：

 index => "tomcat_logs_index_%{+YYYY.MM.dd}"

hosts
是一个数组类型的值

意http协议使用的是http地址，端口是9200，示例：

hosts => ["192.168.102.209:9200", "192.168.102.216:9200"]

document_type
定义es索引的type，一般你应该让同一种类型的日志存到同一种type中，比如debug日志和error日志存到不同的type中

如果不设置默认type为logs

template
如果你愿意，你可以设置指向你自己模板的路径。如果没有设置，那么默认的模板会被使用

template_name
这个配置项用来定义在Elasticsearch中模板的命名

注意删除旧的模板

curl -XDELETE <http://localhost:9200/_template/OldTemplateName?pretty>

template_overwrite
布尔类型 默认为false 
设置为true表示如果你有一个自定义的模板叫logstash，那么将会用你自定义模板覆盖默认模板logstash

manage_template
布尔类型 默认为true 
设置为false将关闭logstash自动管理模板功能 
比如你定义了一个自定义模板，更加字段名动态生成字段，那么应该设置为false

order参数
ELK Stack 在入门学习过程中，必然会碰到自己修改定制索引映射(mapping)乃至模板(template)的问题。 
这时候，不少比较认真看 Logstash 文档的新用户会通过下面这段配置来制定自己的模板策略：

output {
    elasticsearch {
        host => "127.0.0.1"
        manage_template => true
        template => "/path/to/mytemplate"
        template_name => "myname"
    }
}

然而随后就发现，自己辛辛苦苦修改出来的模板，通过 curl -XGET ‘http://127.0.0.1:9200/_template/myname’ 看也确实上传成功了，但实际新数据索引创建出来，就是没生效！

这个原因是：Logstash 默认会上传一个名叫 logstash 的模板到 ES 里。如果你在使用上面这个配置之前，曾经运行过 Logstash（一般来说都会），那么 ES 里就已经存在这么一个模板了。你可以curl -XGET ‘http://127.0.0.1:9200/_template/logstash’ 验证。

这个时候，ES 里就变成有两个模板，logstash 和 myname，都匹配 logstash-* 索引名，要求设置一定的映射规则了。

ES 会按照一定的规则来尝试自动 merge 多个都匹配上了的模板规则，最终运用到索引上

其中要点就是：template 是可以设置 order 参数的！而不写这个参数，默认的 order 值就是 0。order 值越大，在 merge 规则的时候优先级越高。

所以，解决这个问题的办法很简单：在你自定义的 template 里，加一行，变成这样：

{
    "template" : "logstash-*",
    "order" : 1,
    "settings" : { ... },
    "mappings" : { ... }
}

当然，其实如果只从 Logstash 配置角度出发，其实更简单的办法是：直接修改原来默认的 logstash 模板，然后模板名称也不要改，就好了：

output {
    elasticsearch {
        host => "127.0.0.1"
        manage_template => true
        template_overwrite => true
    }
}

为elasticsearch配置模板
在使用logstash收集日志的时候，我们一般会使用logstash自带的动态索引模板，虽然无须我们做任何定制操作，就能把我们的日志数据推送到elasticsearch索引集群中

但是在我们查询的时候，就会发现，默认的索引模板常常把我们不需要分词的字段，给分词了，这样以来，我们的比较重要的聚合统计就不准确了：

所以这时候，就需要我们自定义一些索引模板了

在logstash与elasticsearch集成的时候，总共有如下几种使用模板的方式：

使用默认自带的索引模板 ，大部分的字段都会分词，适合开发和时候快速验证使用

在logstash收集端自定义配置模板，因为分散在收集机器上，维护比较麻烦

在elasticsearc服务端自定义配置模板，由elasticsearch负责加载模板，可动态更改，全局生效，维护比较容易

使用默认自带的索引模板
ElasticSearch默认自带了一个名字为”logstash”的模板，默认应用于Logstash写入数据到ElasticSearch使用

优点：最简单，无须任何配置

缺点：无法自定义一些配置，例如：分词方式

在logstash收集端自定义配置模板
使用第二种，适合小规模集群的日志收集

需要在logstash的output插件中使用template指定本机器上的一个模板json路径， 例如 template => “/tmp/logstash.json”

优点：配置简单

缺点：因为分散在Logstash Indexer机器上，维护起来比较麻烦

在elasticsearc服务端自定义配置模板
manage_template => false//关闭logstash自动管理模板功能  
template_name => "xxx"//映射模板的名字  
第三种需要在elasticsearch的集群中的config/templates路径下配置模板json，在elasticsearch中索引模板可分为两种

静态模板
适合索引字段数据固定的场景，一旦配置完成，不能向里面加入多余的字段，否则会报错

优点：scheam已知，业务场景明确，不容易出现因字段随便映射从而造成元数据撑爆es内存，从而导致es集群全部宕机，维护比较容易，可动态更改，全局生效。

缺点：字段数多的情况下配置稍繁琐

一个静态索引模板配置例子如下：

{  
  "xxx" : {  
      "template": "xxx-*",  
        "settings": {  
            "index.number_of_shards": 3,  
            "number_of_replicas": 0   
        },  
    "mappings" : {  
      "logs" : {  
        "properties" : {  
          "@timestamp" : { //这是专门给kibana用的一个字段，时间索引
            "type" : "date",  
            "format" : "dateOptionalTime",  
            "doc_values" : true  
          },  
          "@version" : {  
            "type" : "string",  
            "index" : "not_analyzed",  
            "doc_values" : true      
          },  
          "id" : {  
            "type" : "string",  
            "index" : "not_analyzed"  
          },  
          "name" : {  
            "type" : "string",  
            "index" : "not_analyzed"  
          }
        }  
      }  
    }  
  }  
}  

动态模板
适合字段数不明确，大量字段的配置类型相同的场景，多加字段不会报错

优点：可动态添加任意字段，无须改动scheaml，

缺点：如果添加的字段非常多，有可能造成es集群宕机

一个动态索引模板配置例子如下：

{  
  "template" : "xxx-*",  
  "settings" : {  
   "index.number_of_shards": 5,  
   "number_of_replicas": 0    

},  
  "mappings" : {  
    "_default_" : {  
      "_all" : {"enabled" : true, "omit_norms" : true},  
      "dynamic_templates" : [ {  
        "message_field" : {  
          "match" : "message",  
          "match_mapping_type" : "string",  
          "mapping" : {  
            "type" : "string", "index" : "analyzed", "omit_norms" : true,  
            "fielddata" : { "format" : "disabled" }  
          }  
        }  
      }, {  
        "string_fields" : {  
          "match" : "*",  
          "match_mapping_type" : "string",  
          "mapping" : {  
            "type" : "string", "index" : "not_analyzed", "doc_values" : true  
          }  
        }  
      } ],  
      "properties" : {  
        "@timestamp": { "type": "date" },  
        "@version": { "type": "string", "index": "not_analyzed" }, 
        "geoip"  : {  
          "dynamic": true,  
          "properties" : {  
            "ip": { "type": "ip" },  
            "location" : { "type" : "geo_point" },  
            "latitude" : { "type" : "float" },  
            "longitude" : { "type" : "float" }  
          }  
        }  
      }  
    }  
  }  
}  

只设置message字段分词，其他的字段默认都不分词

模板结构
通用设置 主要是模板匹配索引的过滤规则，影响该模板对哪些索引生效
settings：配置索引的公共参数，比如索引的replicas，以及分片数shards等参数
mappings：最重要的一部分，在这部分中配置每个type下的每个field的相关属性，比如field类型（string,long,date等等），是否分词，是否在内存中缓存等等属性都在这部分配置
aliases：索引别名，索引别名可用在索引数据迁移等用途上。
例子：

{
  "logstash" : {
    "order" : 0,
    "template" : "logstash-*",
    "settings" : {
      "index" : {
        "refresh_interval" : "5s"
      }
    },
    "mappings" : {
      "_default_" : {
        "dynamic_templates" : [ {
          "message_field" : {
            "mapping" : {
              "fielddata" : {
                "format" : "disabled"
              },
              "index" : "analyzed",
              "omit_norms" : true,
              "type" : "string"
            },
            "match_mapping_type" : "string",
            "match" : "message"
          }
        }, {
          "string_fields" : {
            "mapping" : {
              "fielddata" : {
                "format" : "disabled"
              },
              "index" : "analyzed",
              "omit_norms" : true,
              "type" : "string",
              "fields" : {
                "raw" : {
                  "ignore_above" : 256,
                  "index" : "not_analyzed",
                  "type" : "string"
                }
              }
            },
            "match_mapping_type" : "string",
            "match" : "*"
          }
        } ],
        "_all" : {
          "omit_norms" : true,
          "enabled" : true
        },
        "properties" : {
          "@timestamp" : {
            "type" : "date"
          },
          "geoip" : {
            "dynamic" : true,
            "properties" : {
              "ip" : {
                "type" : "ip"
              },
              "latitude" : {
                "type" : "float"
              },
              "location" : {
                "type" : "geo_point"
              },
              "longitude" : {
                "type" : "float"
              }
            }
          },
          "@version" : {
            "index" : "not_analyzed",
            "type" : "string"
          }
        }
      }
    },
    "aliases" : { }
  }
}

我们创建一个自定义Template动态模板，这个模板指定匹配所有以”go_logsindex“开始的索引，并且指定允许添加新字段，匹配所有string类型的新字段会创建一个raw的嵌套字段，这个raw嵌套字段类型也是string，但是是not_analyzed不分词的（主要用于解决一些analyzed的string字段无法做统计，但可以使用这个raw嵌套字段做统计）

{
  "template": "go_logs_index_*",
  "order":0,
  "settings": {
      "index.number_of_replicas": "1",
      "index.number_of_shards": "5",
      "index.refresh_interval" : "10s"
  },
  "mappings": {
    "_default_": {
      "_all": {
        "enabled": false
      },
      "dynamic_templates": [
        {
          "my_template": {
            "match_mapping_type": "string",
            "mapping": {
              "type": "string",
              "fields": {
                "raw": {
                  "type": "string",
                  "index": "not_analyzed"
                }
              }
            }
          }
        }
      ]
    },
    "go": {
      "properties": {
        "timestamp": {
          "type": "string",
          "index": "not_analyzed"
        },
        "msg": {
          "type": "string",
          "analyzer": "ik",
          "search_analyzer": "ik_smart"
        },
        "file": {
          "type": "string",
          "index": "not_analyzed"
        },
        "line": {
          "type": "string",
          "index": "not_analyzed"
        },
        "threadid": {
          "type": "string",
          "index": "not_analyzed"
        },
        "info": {
          "type": "string",
          "index": "not_analyzed"
        },
        "type": {
          "type": "string",
          "index": "not_analyzed"
        },
        "@timestamp": {
          "format": "strict_date_optional_time||epoch_millis",
          "type": "date"
        },
        "@version": {
          "type": "string",
          "index": "not_analyzed"
        }
      }
    }
  }
}
总结
第三种方式统一管理Template最好，推荐使用第三种方式，但是具体问题具体分析。例如场景是Logstash 和ElasticSearch都在一台服务器，第二种就比较好

定制索引模板，是搜索业务中一项比较重要的步骤，需要注意的地方有很多，比如： 
1.字段数固定吗 
2.字段类型是什么 
3.分不分词 
4.索引不索引 
5.存储不存储 
6.排不排序 
7.是否加权
