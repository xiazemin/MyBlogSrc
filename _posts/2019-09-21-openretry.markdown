---
title: openretry
layout: post
category: web
author: 夏泽民
---
OpenResty® 是一个基于 Nginx 与 Lua 的高性能 Web 平台，其内部集成了大量精良的 Lua 库、第三方模块以及大多数的依赖项。用于方便地搭建能够处理超高并发、扩展性极高的动态 Web 应用、Web 服务和动态网关。
http://openresty.org/cn/

OpenResty 1.15.8.2 Released
<!-- more -->
Nginx 是俄罗斯人发明的， Lua 是巴西几个教授发明的，中国人章亦春把 LuaJIT VM 嵌入到 Nginx 中，实现了 OpenResty 这个高性能服务端解决方案。

通过 OpenResty，你可以把 nginx 的各种功能进行自由拼接， 更重要的是，开发门槛并不高，这一切都是用强大轻巧的 Lua 语言来操控。

它主要的使用场景主要是：

在 Lua 中揉和和处理各种不同的 nginx 上游输出（Proxy，Postgres，Redis，Memcached 等）

在请求真正到达上游服务之前，Lua 可以随心所欲的做复杂的访问控制和安全检测

随心所欲的操控响应头里面的信息

从外部存储服务（比如 Redis，Memcached，MySQL，Postgres）中获取后端信息，并用这些信息来实时选择哪一个后端来完成业务访问

在内容 handler 中随意编写复杂的 Web 应用，使用 同步但依然非阻塞 的方式，访问后端数据库和其他存储

在 rewrite 阶段，通过 Lua 完成非常复杂的 URL dispatch

用 Lua 可以为 nginx 子请求和任意 location，实现高级缓存机制

有用来写 WAF、有做 CDN 调度、有做广告系统、消息推送系统，还有像我们部门一样，用作 API server 的。有些还用在非常关键的业务上，比如开涛在高可用架构分享的京东商品详情页，是我知道的 ngx_lua 最大规模的应用。

ngx_openresty 目前有两大应用目标：

通用目的的 web 应用服务器。在这个目标下，现有的 web 应用技术都可以算是和 OpenResty 或多或少有些类似，比如 Nodejs, PHP 等等。ngx_openresty 的性能（包括内存使用和 CPU 效率）算是最大的卖点之一。
Nginx 的脚本扩展编程，用于构建灵活的 Web 应用网关和 Web 应用防火墙。有些类似的是 NetScaler。其优势在于 Lua 编程带来的巨大灵活性。

二．OpenResty运行原理

Nginx 采用的是 master-worker 模型，一个 master 进程管理多个 worker 进程，基本的事件处理都是放在 woker 中，master 负责一些全局初始化，以及对 worker 的管理。在OpenResty中，每个 woker 使用一个 LuaVM，当请求被分配到 woker 时，将在这个 LuaVM 里创建一个 coroutine(协程)。协程之间数据隔离，每个协程具有独立的全局变量_G。
