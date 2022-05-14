---
title: Nightingale
layout: post
category: algorithm
author: 夏泽民
---
与Open-Falcon的异同
由于Nightingale与Open-Falcon有很深的血缘关系，想必大家对二者的对比会较为感兴趣。总体来看，理念上有很强的延续性，架构和产品体验上的改进巨大。

与Open-Falcon的不同点
告警引擎重构为推拉结合模式，通过推模式保证大部分策略判断的效率，通过拉模式支持了nodata告警，去除原来的nodata组件，简化系统部署难度
引入了服务树，对endpoint进行层级管理，去除原来扁平的HostGroup，同时干掉告警模板，告警策略直接与服务树节点绑定，大幅提升灵活度和易用性
干掉原来的基于数据库的索引库，改成内存模式，单独抽出一个index模块处理索引，避免了原来MySQL单表达到亿级的尴尬局面，索引基于内存之后效率也大幅提升
存储模块Graph，引入facebook的Gorilla，即内存TSDB，近期几个小时的数据默认存内存，大幅提升数据查询效率，硬盘存储方式仍然使用rrdtool
告警引擎judge模块通过心跳机制做到了故障自动摘除，再也不用担心单个judge挂掉导致部分策略失效的问题，index模块也是采用类似方式保证可用性
客户端中内置了日志匹配指标抽取能力，web页面上支持了日志匹配规则的配置，同时也支持读取目标机器特定目录下的配置文件的方式，让业务指标监控更为易用
将portal(falcon-plus中的api)、uic、dashboard、hbs、alarm合并为一个模块：monapi，简化了系统整体部署难度，原来的部分模块间调用变成进程内方法调用，稳定性更好
与Open-Falcon的相同点
数据模型没有变化，仍然是metric、endpoint、tags的组织方式，所以agent基本是可以复用的，Nightingale中的agent叫collector，融合了原来Open-Falcon的agent和falcon-log-agent的逻辑，各种监控插件也都是可以复用的
数据流向和整体处理逻辑是类似的，仍然使用灵活的推模型，分为数据存储和告警判断两条链路，
<!-- more -->
https://www.bookstack.cn/read/Nightingale/a194607102c18979.md

https://www.bookstack.cn/read/Nightingale/d7a2753d1138e6af.md