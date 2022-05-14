---
title: redis slow log
layout: post
category: storage
author: 夏泽民
---
<!-- more -->
redis是目前最流行的缓存系统，因其丰富的数据结构和良好的性能表现，被各大公司广泛使用。尽管redis性能极佳，但若不注意使用方法，极容易出现慢查询，慢查询多了或者一个20s的慢查询会导致操作队列（redis是单进程）堵塞，最终引起雪崩甚至整个服务不可用。对于慢查询语句，redis提供了相关的配置和命令。
配置有两个：slowlog-log-slower-than 和 slowlog-max-len。slowlog-log-slower-than是指当命令执行时间（不包括排队时间）超过该时间时会被记录下来，单位为微秒，比如通过下面的命令，就可以记录执行时长超过20ms的命令了。

config set slowlog-log-slower-than 20000
slowlog-max-len是指redis可以记录的慢查询命令的总数，比如通过下面的命令，就可以记录最近100条慢查询命令了。

config set slowlog-max-len 100
操作慢查询的命令有两个：slowlog get [len] 和 slowlog reset。slowlog get [len]命令获取指定长度的慢查询列表。

redis 127.0.0.1:6379> slowlog get 2
1) 1) (integer) 14
   2) (integer) 1309448221
   3) (integer) 15
   4) 1) "ping"
2) 1) (integer) 13
   2) (integer) 1309448128
   3) (integer) 30
   4) 1) "slowlog"
      2) "get"
      3) "100"
上面返回了两个慢查询命令，其中每行的含义如下：

第一行是一个慢查询id。该id是自增的，只有在 redis server 重启时该id才会重置。
第二行是慢查询命令执行的时间戳
第三行是慢查询命令执行耗时，单位为微秒
第四行是慢查询命令的具体内容。
slowlog reset命令是清空慢日志队列。

elasticsearch用来存储解析后的redis slowlog，kibana用于图形化分析，beats用于收集redis slowlog。
这里着重讲一下beats，它是一系列轻量级的数据收集产品统称，目前官方提供了filebeat、packetbeat、heartbeat、metricbeat等，可以用来收集日志文件、网络包、心跳包、各类指标数据等。像我们这次要收集的redis slowlog，官方还没有提供相关工具，需要我们自己实现，但借助beats的一系列脚手架工具，我们可以方便快速的创建自己的rsbeat---redis slowlog beat。

rsbeat原理简介
接下来我们先讲解一下rsbeat的实现原理，一图胜千言，我们先来看下它的工作流。

rsbeat工作流
我们由下往上分析：

最下面是我们要分析的redis server列表
再往上便是rsbeat，它会与这些redis server建立连接并定期去拉取 slowlog。
在启动时，rsbeat会发送下面的命令到每一台redis server，来完成slowlog的配置，这里设置记录最近执行时长超过20ms的500条命令。
config set slowlog-log-slower-than 20000
config set slowlog-max-len 500
slowlog reset
然后rsbeat会定时去拉取每台redis server的慢查询命令
slowlog get 500
slowlog reset
注意之类slowlog reset是因为此次已经将所有的慢日志都取出了，下次获取时取最新生成的，防止重复计算。

rsbeat将解析的慢日志发布到elasticsearch中进行存储
通过kibana进行slowlog的图形化分析
rsbeat的整个工作流到这里已经介绍完毕了，是不是很简单呢？下面我们来简单看一下rsbeat的核心代码实现。

rsbeat核心代码讲解
rsbeat已经在github上开源了，感兴趣的同学可以自己去下下来使用。下面我们分析的代码位于beater/rsbeat.go，这也是rsbeat的核心文件。

func poolInit(server string, slowerThan int) *redis.Pool {
    return &redis.Pool{
        MaxIdle:     3,
        MaxActive:   3,
        IdleTimeout: 240 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", server, redis.DialConnectTimeout(3*time.Second), redis.DialReadTimeout(3*time.Second))
            if err != nil {
                logp.Err("redis: error occurs when connect %v", err.Error())
                return nil, err
            }
            c.Send("MULTI")
            c.Send("CONFIG", "SET", "slowlog-log-slower-than", slowerThan)
            c.Send("CONFIG", "SET", "slowlog-max-len", 500)
            c.Send("SLOWLOG", "RESET")
            r, err := c.Do("EXEC")

            if err != nil {
                logp.Err("redis: error occurs when send config set %v", err.Error())
                return nil, err
            }

            logp.Info("redis: config set %v", r)
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            logp.Info("redis: PING")
            return err
        },
    }
}
poolInit方法是rsbeat初始化时进行的操作，这里也就是发送slowlog配置的地方，代码很简单，就不展开解释了。

func (bt *Rsbeat) redisc(beatname string, init bool, c redis.Conn, ipPort string) {
    defer c.Close()
    logp.Info("conn:%v", c)

    c.Send("SLOWLOG", "GET")
    c.Send("SLOWLOG", "RESET")
    logp.Info("redis: slowlog get. slowlog reset")

    c.Flush()
    reply, err := redis.Values(c.Receive()) // reply from GET
    c.Receive()                             // reply from RESET

    logp.Info("reply len: %d", len(reply))

    for _, i := range reply {
        rp, _ := redis.Values(i, err)
        var itemLog itemLog
        var args []string
        redis.Scan(rp, &itemLog.slowId, &itemLog.timestamp, &itemLog.duration, &args)
        argsLen := len(args)
        if argsLen >= 1 {
            itemLog.cmd = args[0]
        }
        if argsLen >= 2 {
            itemLog.key = args[1]
        }
        if argsLen >= 3 {
            itemLog.args = args[2:]
        }
        logp.Info("timestamp is: %d", itemLog.timestamp)
        t := time.Unix(itemLog.timestamp, 0).UTC()

        event := common.MapStr{
            "type":           beatname,
            "@timestamp":     common.Time(time.Now()),
            "@log_timestamp": common.Time(t),
            "slow_id":        itemLog.slowId,
            "cmd":            itemLog.cmd,
            "key":            itemLog.key,
            "args":           itemLog.args,
            "duration":       itemLog.duration,
            "ip_port":        ipPort,
        }

        bt.client.PublishEvent(event)
    }
}
redisc方法实现了定时从redis server拉取最新的slowlog列表，并将它们转化为elasticsearch中可以存储的数据后，发布到elasticsearch中。这里重点说下每一个字段的含义：

@timestamp是指当前时间戳。
@log_timestamp是指慢日志命令执行的时间戳。
slow_id是该慢日志的id。
cmd是指执行的 redis 命令，比如sadd、scard等等。
key是指redis key的名称
args是指 redis 命令的其他参数，通过 cmd、key、args我们可以完整还原执行的redis命令。
duration是指redis命令执行的具体时长，单位是微秒。
ip_port是指发生命令的 redis server 地址。
有了这些字段，我们就可以用kibana来愉快的进行可视化数据分析了。
