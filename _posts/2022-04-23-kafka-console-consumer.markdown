---
title: kafka-console-consumer
layout: post
category: storage
author: 夏泽民
---
kafka-console-consumer.sh 脚本是一个简易的消费者控制台。该 shell 脚本的功能通过调用 kafka.tools 包下的 ConsoleConsumer 类，并将提供的命令行参数全部传给该类实现。

注意：Kafka 从 2.2 版本开始将 kafka-topic.sh 脚本中的 −−zookeeper 参数标注为 “过时”，推荐使用 −−bootstrap-server 参数。若读者依旧使用的是 2.1 及以下版本，请将下述的 --bootstrap-server 参数及其值手动替换为 --zookeeper zk1:2181,zk2:2181,zk:2181。一定要注意两者参数值所指向的集群地址是不同的。

消息消费
bin/kafka-console-consumer.sh --bootstrap-server node1:9092,node2:9092,node3:9092 --topic topicName
1
 表示从 latest 位移位置开始消费该主题的所有分区消息，即仅消费正在写入的消息。
从开始位置消费
bin/kafka-console-consumer.sh --bootstrap-server node1:9092,node2:9092,node3:9092 --from-beginning --topic topicName
1
 表示从指定主题中有效的起始位移位置开始消费所有分区的消息。
显示key消费
bin/kafka-console-consumer.sh --bootstrap-server node1:9092,node2:9092,node3:9092 --property print.key=true --topic topicName
1
 消费出的消息结果将打印出消息体的 key 和 value。
<!-- more -->
https://blog.csdn.net/qq_29116427/article/details/80206125