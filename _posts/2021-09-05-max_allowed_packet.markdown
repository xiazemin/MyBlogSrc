---
title: max_allowed_packet
layout: post
category: storage
author: 夏泽民
---
服务器的日志一直报Packet for query is too large (7632997 > 4194304). You can change this value on the server by setting the max_allowed_packet’ variable.的解决方法

max_allowed_packet 值设置过小将导致单个记录超过限制后写入数据库失败，且后续记录写入也将失败，为了数据完整性，需要考虑到事务因素。

MySQL的一个系统参数问题:max_allowed_packet，其默认值为1048576(1M)，查询：show VARIABLES like ‘%max_allowed_packet%’;
<!-- more -->
https://www.cnblogs.com/jimloveq/p/10609487.html
