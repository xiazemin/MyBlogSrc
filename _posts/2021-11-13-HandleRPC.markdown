---
title: grpc HandleRPC
layout: post
category: golang
author: 夏泽民
---
服务器需要为每个连接保存各自的数据。连接创建时初始化数据，连接断开时清理数据。
这里利用了连接统计的接口，不知道是否是最适当的实现方式?

服务器创建时添加 StatsHandler 选项，输入一个 stats.Handler 的实现。

-   s := grpc.NewServer()
+   s := grpc.NewServer(grpc.StatsHandler(&statshandler{}))

statshandler 需实现4个方法，只用到2个连接相关的方法，TagConn() 和 HandleConn(),
另外2个 TagRPC() 和 HandleRPC() 用于RPC统计, 实现为空。
<!-- more -->

https://blog.csdn.net/jq0123/article/details/78895737
