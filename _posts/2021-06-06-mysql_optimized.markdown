---
title: Select tables optimized away
layout: post
category: storage
author: 夏泽民
---
EXPLAIN   SELECT MAX(`year`) FROM st_sch_recruit_info info

输出的结果里,Extra列输出了"Select tables optimized away"语句。

这个在MySQL的手册里面没有任何提及，不过看其他各列的数据大概能猜到意思：SELECT操作已经优化到不能再优化了（MySQL根本没有遍历表或索引就返回数据了）。

在MySQL官方站点翻到两段相关的描述，印证了上述观点，原文如下：
For explains on simple count queries (i.e. explain select count(*) from people) the extra section will read "Select tables optimized away." This is due to the fact that MySQL can read the result directly from the table internals and therefore does not need to perform the select.

https://blog.csdn.net/persistencegoing/article/details/91441084
<!-- more -->
1.optimize table

如果MySQL没有选中正确的索引，有可能是因为表经常被更改。这会影响统计数据。如果时间允许(表在此期间是锁定的)，我们可以通过重新构建表来帮助解决这个问题。

2.analyze table

analyze table需要的时间更少，尤其是在InnoDB中。分析将更新索引统计信息并帮助生成更好的查询计划。

3.使用hint

比如使用关键字use index(index-name)、force index(index-name)

select c1 from abce use index(idx_c1) where ... ;

select c1 from abce force index(idx_c1) where ... ;

4.忽略索引

如果使用了错误的索引，可以尝试使用关键字来忽略使用被选中的索引。比如，让sql忽略掉主键：

select id from abce ignore index(primary) where ... ;

5.修改业务的逻辑结构，从而修改sql语句

6.否定索引的使用

select id, type, age from abce

where type=12345 and age > 3

order by id+0;#这样id+0是函数运算，就不使用id上的索引了

7.修改sql的结构，使得优化器认为选择该索引的代价会更高

select id, type, age from abce

where type=12345 and age > 3

order by id;

#修改成

select id, type, age from abce

where type=12345 and age > 3

order by id,type, age;

这样，可能机会导致优化器认为选择使用id上的索引的代价会更高。从而选择其他索引。

(以上id只是用作示例，并不一定是主键的列)

8.有时候可能是bug，报bug，等待官方修复。

https://blog.csdn.net/weixin_39562340/article/details/113232364
