---
title: prom-client
layout: post
category: node
author: 夏泽民
---
https://github.com/siimon/prom-client
这是一个支持histogram, summaries, gauges and counters四种数值格式的prometheus nodejs客户端。
停止轮询默认metrics

要停止采集默认metrics，你需要调用调用函数并传给clearInterval。

const client = require('prom-client');

clearInterval(client.collectDefaultMetrics());

// Clear the register
client.register.clear();

<!-- more -->
https://www.jianshu.com/p/46dfef9ff582?utm_source=oschina-app


