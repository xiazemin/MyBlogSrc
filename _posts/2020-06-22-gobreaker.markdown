---
title: gobreaker
layout: post
category: golang
author: 夏泽民
---
https://github.com/sony/gobreaker

当我们依赖的服务出现问题时，可以及时容错。一方面可以减少依赖服务对自身访问的依赖，防止出现雪崩效应；另一方面降低请求频率以方便上游尽快恢复服务。

熔断器的应用也非常广泛。除了在我们应用中，为了请求服务时使用熔断器外，在 web 网关、微服务中，也有非常广泛的应用。本文将从源码角度学习sony 开源的一个熔断器实现 github/sony/gobreaker。（代码注释可以从github/lpflpf/gobreaker 查看)
<!-- more -->
https://blog.lpflpf.cn/passages/circuit-breaker/
