---
title: dockertest
layout: post
category: k8s
author: 夏泽民
---
如果你使用的是 Go，则可以使用 dockertest，一个可以管理和编排 Go 测试文件中的容器的库。

从 Go 文件管理测试基础设施容器，允许我们控制在每个测试中需要的服务（例如，某些包正在使用数据库而不是 Redis，为这个测试运行 Redis 没有意义）

go get -u github.com/ory/dockertest/v3
<!-- more -->
https://mp.weixin.qq.com/s/08lnsEbxqD1ziX7BmBPQOQ
