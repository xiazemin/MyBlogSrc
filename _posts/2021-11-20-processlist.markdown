---
title:  批量 Kill  mysql processlist 进程
layout: post
category: mysql
author: 夏泽民
---
https://www.cnblogs.com/bianxj/articles/9605067.html
<!-- more -->
 如果大批量的操作能够通过一系列的select 语句产生，那么理论上就能对这些结果批量处理。

          但是mysql并没有提供eval这样的对结果集进行分析操作的功能。索引只能将select结果保存到临时文件中，然后再执行临时文件中的指令。

具体过程如下

1、通过information_schema.processlist表中的连接信息生成需要处理掉的MySQL连接的语句临时文件，然后执行临时文件中生成的指令

复制代码
mysql> select concat('KILL ',id,';') from information_schema.processlist where user='root';
+------------------------+
| concat('KILL ',id,';') 
+------------------------+
| KILL 3101;             
| KILL 2946;             
+------------------------+
2 rows in set (0.00 sec)
 
mysql>select concat('KILL ',id,';') from information_schema.processlist where user='root' into outfile '/tmp/a.txt';
Query OK, 2 rows affected (0.00 sec)
 
mysql>source /tmp/a.txt;
Query OK, 0 rows affected (0.00 sec)
复制代码
 

2、杀掉当前所有的MySQL连接

 

mysqladmin -uroot -p processlist|awk -F "|" '{print $2}'|xargs -n 1 mysqladmin -uroot -p kill   
          

      杀掉指定用户运行的连接，这里为sa

   

mysqladmin -uroot -p processlist|awk -F "|" '{if($3 == "sa")print $2}'|xargs -n 1 mysqladmin -uroot -p kill
 

    3、通过shell脚本实现

#杀掉锁定的MySQL连接
for id in `mysqladmin processlist|grep -i locked|awk '{print $1}'`
do
   mysqladmin kill ${id}
done
