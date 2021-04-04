---
title: kafka 指定 partition
layout: post
category: golang
author: 夏泽民
---
Partition：Topic的分区，每个topic可以有多个分区，分区的作用是做负载，提高kafka的吞 吐量。同一个topic在不同的分区的数据是不重复的，partition的表现形式就是一个一个的⽂件夹！

 kafka中有几个原则：

1.partition在写入的时候可以指定需要写入的partition，如果有指定，则写入对应的partition。
2.如果没有指定partition，但是设置了数据的key，则会根据key的值hash出一个partition。
3.如果既没指定partition，又没有设置key，则会采用轮询⽅式，即每次取一小段时间的数据写入某
个partition，下一小段的时间写入下一个partition

在同⼀个消费者组中，每个消费者实例可以消费多个分区，但是每个分区最多只 能被消费者组中的⼀个实例消费。

https://www.bookstack.cn/read/topgoer/68f58bcd82d11624.md
<!-- more -->
https://juejin.cn/post/6844903903113248776

https://qixiang-liu.github.io/post/golang%E7%AC%94%E8%AE%B0/2019-10-12-go%E6%93%8D%E4%BD%9Ckafka/

https://studygolang.com/articles/17912?fr=sidebar

https://blog.csdn.net/feelwing1314/article/details/81097167

kafka的原生API可以使用consumer.assign(partitions)来订阅指定分区，spring kafka的API有没有相应的方法？我只找到用@KafkaListener(topicPartitions ={@TopicPartition(topic = "topic1", partitions = { "0", "1" }))}注解实现的，但这种方式实现时topic、partitions的值都必须为常量


https://segmentfault.com/q/1010000016811772

