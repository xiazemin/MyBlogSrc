---
title: mysql optimizer trace
layout: post
category: storage
author: 夏泽民
---
explain为mysql提供语句的执行计划信息。可以应用在select、delete、insert、update和place语句上。explain的执行计划，只是作为语句执行过程的一个参考，实际执行的过程不一定和计划完全一致，但是执行计划中透露出的讯息却可以帮助选择更好的索引和写出更优化的查询语句。

EXPLAIN输出项（可参考mysql5.7文档）
备注:当使用FORMAT=JSON， 返回的数据为json结构时，JSON Name为null的不显示

https://blog.csdn.net/aecuhty88306453/article/details/102196581
<!-- more -->
https://blog.csdn.net/weixin_38004638/article/details/106427205

step 1

 使用explain 查看执行计划， 5.6后可以加参数 explain format=json xxx 输出json格式的信息

 

step 2 

使用profiling详细的列出在每一个步骤消耗的时间，前提是先执行一遍语句。

 

#打开profiling 的设置
SET profiling = 1;
SHOW VARIABLES LIKE '%profiling%';

#查看队列的内容
show profiles;  
#来查看统计信息
show profile block io,cpu for query 3;

step 3  

 

Optimizer trace是MySQL5.6添加的新功能，可以看到大量的内部查询计划产生的信息, 先打开设置，然后执行一次sql,最后查看`information_schema`.`OPTIMIZER_TRACE`的内容

 

#打开设置
SET optimizer_trace='enabled=on';  
#最大内存根据实际情况而定， 可以不设置
SET OPTIMIZER_TRACE_MAX_MEM_SIZE=1000000;
SET END_MARKERS_IN_JSON=ON;
SET optimizer_trace_limit = 1;
SHOW VARIABLES LIKE '%optimizer_trace%';

#执行所需sql后，查看该表信息即可看到详细的执行过程
SELECT * FROM `information_schema`.`OPTIMIZER_TRACE`;

https://www.cnblogs.com/mydriverc/p/7086542.html

optimizer_trace。前两个经常被人使用，由于第三个难度较大，大家使用的较少，下面简单说下如何使用。
 
 http://www.bubuko.com/infodetail-2662893.html