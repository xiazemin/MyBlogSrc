---
title: Logstash_example
layout: post
category: elasticsearch
author: 夏泽民
---
<!-- more -->
1，Elasticsearch slow-log
input {
    file {
        path => ["/var/log/elasticsearch/private_test_index_search_slowlog.log"]
        start_position => "beginning"
        ignore_older => 0
        # sincedb_path => "/dev/null"
        type => "elasticsearch_slow"
        }   
}

filter {
    grok {
        match =>  { "message" => "^\[(\d\d){1,2}-(?:0[1-9]|1[0-2])-(?:(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9])\s+(?:2[0123]|[01]?[0-9]):(?:[0-5][0-9]):(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)\]\[(TRACE|DEBUG|WARN\s|INFO\s)\]\[(?<io_type>[a-z\.]+)\]\s\[(?<node>[a-z0-9\-\.]+)\]\s\[(?<index>[A-Za-z0-9\.\_\-]+)\]\[\d+\]\s+took\[(?<took_time>[\.\d]+(ms|s|m))\]\,\s+took_millis\[(\d)+\]\,\s+types\[(?<types>([A-Za-z\_]+|[A-Za-z\_]*))\]\,\s+stats\[\]\,\s+search_type\[(?<search_type>[A-Z\_]+)\]\,\s+total_shards\[\d+\]\,\s+source\[(?<source>[\s\S]+)\]\,\s+extra_source\[[\s\S]*\]\,\s*$" }
        remove_field => ["message"]
        }   

    date {
        match => ["timestamp","dd/MMM/yyyy:HH:mm:ss Z"] 
        }   
    ruby {
        code => "event.timestamp.time.localtime"
        }   
    }

output {
     elasticsearch {
         codec => "json"
         hosts => ["127.0.0.1:9200"]
         index => "logstash-elasticsearch-slow-%{+YYYY.MM.dd}"
         user => "admin"
         password => "xxxx"
    }   

}

2，Mysql-slow log
input {
    file {
        path => "/var/lib/mysql/slow.log"
        start_position => "beginning"
        ignore_older => 0
        # sincedb_path => "/dev/null"
        type => "mysql-slow"
    }   
}
filter {
    if ([message] =~ "^(\/usr\/local|Tcp|Time)[\s\S]*") { drop {} }
    multiline {
        pattern => "^\#\s+Time\:\s+\d+\s+(0[1-9]|[12][0-9]|3[01]|[1-9])"
        negate => true
        what => "previous"
    }   
    grok {
        match => { "message" => "^\#\sTime\:\s+\d+\s+(?<datetime>%{TIME})\n+\#\s+User@Host\:\s+[A-Za-z0-9\_]+\[(?<mysql_user>[A-Za-z0-9\_]+)\]\s+@\s+(?<mysql_host>[A-Za-z0-9\_]+)\s+\[\]\n+\#\s+Query\_time\:\s+(?<query_time>[0-9\.]+)\s+Lock\_time\:\s+(?<lock_time>[0-9\.]+)\s+Rows\_sent\:\s+(?<rows_sent>\d+)\s+Rows\_examined\:\s+(?<rows_examined>\d+)(\n+|\n+use\s+(?<dbname>[A-Za-z0-9\_]+)\;\n+)SET\s+timestamp\=\d+\;\n+(?<slow_message>[\s\S]+)$"
   }   
        remove_field => ["message"]
   }   
    date {
        match => ["timestamp","dd/MMM/yyyy:HH:mm:ss Z"] 
    }   
    ruby {
        code => "event.timestamp.time.localtime"
    }   
}
output { 
    elasticsearch {
        codec => "json"
        hosts => ["127.0.0.1:9200"]
        index => "logstash-mysql-slow-%{+YYYY.MM.dd}"
        user => "admin"
        password => "xxxxx"
    }   
}

3，Nginx access.log
logstash 中内置 nginx 的正则,我们只要稍作修改就能使用 
将下面的内容写入到/opt/logstash/vendor/bundle/jruby/1.9/gems/logstash- 
patterns-core-2.0.5/patterns/grok-patterns 文件中

X_FOR (%{IPV4}|-)

NGINXACCESS %{COMBINEDAPACHELOG} \"%{X_FOR:http_x_forwarded_for}\"

ERRORDATE %{YEAR}/%{MONTHNUM}/%{MONTHDAY} %{TIME}

NGINXERROR_ERROR %{ERRORDATE:timestamp}\s{1,}\[%{DATA:err_severity}\]\s{1,}(%{NUMBER:pid:int}#%{NUMBER}:\s{1,}\*%{NUMBER}|\*%{NUMBER}) %{DATA:err_message}(?:,\s{1,}client:\s{1,}(?<client_ip>%{IP}|%{HOSTNAME}))(?:,\s{1,}server:\s{1,}%{IPORHOST:server})(?:, request: %{QS:request})?(?:, host: %{QS:server_ip})?(?:, referrer:\"%{URI:referrer})?

NGINXERROR_OTHER %{ERRORDATE:timestamp}\s{1,}\[%{DATA:err_severity}\]\s{1,}%{GREEDYDATA:err_message}

之后的 log 配置文件如下

input {
    file {
    path => [ "/var/log/nginx/www-access.log" ]
    start_position => "beginning"
    # sincedb_path => "/dev/null"
    type => "nginx_access"
    }   
}
filter {
    grok {
         match => { "message" => "%{NGINXACCESS}"}
    }
    mutate {
        convert => [ "response","integer" ]
        convert => [ "bytes","integer" ]
    }
    date {
        match => [ "timestamp","dd/MMM/yyyy:HH:mm:ss Z"]
    }   
    ruby {
        code => "event.timestamp.time.localtime"
    }   
}
output {
    elasticsearch {
        codec => "json"
        hosts => ["127.0.0.1:9200"]
        index => "logstash-nginx-access-%{+YYYY.MM.dd}"
        user => "admin"
        password => "xxxx"
    }
}

4，Nginx error.log
input {
    file {
    path => [ "/var/log/nginx/www-error.log" ]
    start_position => "beginning"
    # sincedb_path => "/dev/null"
    type => "nginx_error"
    }
}
filter {
    grok {
        match => [
                   "message","%{NGINXERROR_ERROR}",
                   "message","%{NGINXERROR_OTHER}"
                 ]
    }   
    ruby {
        code => "event.timestamp.time.localtime"
    }   
     date {
         match => [ "timestamp","dd/MMM/yyyy:HH:mm:ss"]
     } 

}

output {
    elasticsearch {
        codec => "json"
        hosts => ["127.0.0.1:9200"]
        index => "logstash-nginx-error-%{+YYYY.MM.dd}"
        user => "admin"
        password => "xxxx"
    }   
}

5，PHP error.log
input {
    file {
        path => ["/var/log/php/error.log"]
        start_position => "beginning"
        # sincedb_path => "/dev/null"
        type => "php-fpm_error"
    }   
}

filter {
    multiline {
        pattern => "^\[(0[1-9]|[12][0-9]|3[01]|[1-9])\-%{MONTH}-%{YEAR}[\s\S]+"
        negate => true
        what => "previous"
    }   
    grok {
        match => { "message" => "^\[(?<timestamp>(0[1-9]|[12][0-9]|3[01]|[1-9])\-%{MONTH}-%{YEAR}\s+%{TIME}?)\s+[A-Za-z]+\/[A-Za-z]+\]\s+(?<category>(?:[A-Z]{3}\s+[A-Z]{1}[a-z]{5,7}|[A-Z]{3}\s+[A-Z]{1}[a-z\s]{9,11}))\:\s+(?<error_message>[\s\S]+$)" }

        remove_field => ["message"]
    }   

    date {
        match => ["timestamp","dd/MMM/yyyy:HH:mm:ss Z"] 
    }   

    ruby {
        code => "event.timestamp.time.localtime"
    }   

}

output {
    elasticsearch {
        codec => "json"
        hosts => ["127.0.0.1:9200"]
        index => "logstash-php-error-%{+YYYY.MM.dd}"
        user => "admin"
        password => "xxxxx"
    }   
}

6，Php-fpm slow-log
input {
    file {
        path => ["/var/log/php-fpm/www.slow.log"]
        start_position => "beginning"
        # sincedb_path => "/dev/null"
        type => "php-fpm_slow"
    }   
}

filter {
    multiline {
        pattern => "^$"
        negate => true
        what => "previous"
    }   
    grok {
        match => { "message" => "^\[(?<timestamp>(0[1-9]|[12][0-9]|3[01]|[1-9])\-%{MONTH}-%{YEAR}\s+%{TIME})\]\s+\[[a-z]{4}\s+(?<pool>[A-Za-z0-9]{1,8})\]\s+[a-z]{3}\s+(?<pid>\d{1,7})\n(?<slow_message>[\s\S]+$)" }

        remove_field => ["message"]
    }   

    date {
        match => ["timestamp","dd-MMM-yyyy:HH:mm:ss Z"] 
    }   

    ruby {
        code => "event.timestamp.time.localtime"
    }   

}

output {

    elasticsearch {
        codec => "json"
        hosts => ["127.0.0.1:9200"]
        index => "logstash-php-fpm-slow-%{+YYYY.MM.dd}"
        user => "admin"
        password => "xxxx"
    }   
}
