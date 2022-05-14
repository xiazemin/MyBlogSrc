---
title: event_scheduler
layout: post
category: storage
author: 夏泽民
---
1、  查看MySQL是否开启了事件功能

查看命令：

show variables like '%sc%';

打开event_scheduler（临时开启，MySQL服务重启后时效）

SET GLOBAL event_scheduler = ON;

永久开启方法：my.cnf中[mysqld]添加event_scheduler=on #重启服务

2、  创建事件

create event myevent on SCHEDULE every 5 second do delete from Syslog.SystemEvents where ReceivedAt<(CURRENT_TIMESTAMP() + INTERVAL -5 DAY);#删除5天前的数据

说明：

ReceivedAt：数据库Syslog.Systemevents表中的时间字段

(date,INTERVAL expr type):

date：数据库当前时间CURRENT_TIMESTAMP()

INTERVAL：关键字（间隔）

expr：具体的时间间隔（-5）

type:时间单位：

MICROSECOND

间隔单位：毫秒

SECOND

间隔单位：秒

MINUTE

间隔单位：分钟

HOUR

间隔单位：小时

DAY

间隔单位：天

WEEK

间隔单位：星期

MONTH

间隔单位：月

QUARTER

间隔单位：季度

YEAR

间隔单位：年

SECOND_MICROSECOND

复合型，间隔单位：秒、毫秒，expr可以用两个值来分别指定秒和毫秒

MINUTE_MICROSECOND

复合型，间隔单位：分、毫秒

MINUTE_SECOND

复合型，间隔单位：分、秒

HOUR_MICROSECOND

复合型，间隔单位：小时、毫秒

HOUR_SECOND

复合型，间隔单位：小时、秒

HOUR_MINUTE

复合型，间隔单位：小时分

DAY_MICROSECOND

复合型，间隔单位：天、毫秒

DAY_SECOND

复合型，间隔单位：天、秒

DAY_MINUTE

复合型，间隔单位：天、分

DAY_HOUR

复合型，间隔单位：天、小时

YEAR_MONTH

复合型，间隔单位：年、月

  

 

如果存在事件，请先删除，删除命令：drop event if exists myevent;

3、  开启事件

alter event myevent on completion preserve enable;

4、关闭事件的命令：alter event myevent on completion preserve disable;
<!-- more -->
由于用户环境有张日志表每天程序都在狂插数据，导致不到一个月时间，这张日志表就高达200多万条记录，但是日志刷新较快，里面很多日志没什么作用，就写了个定时器，定期删除这张表的数据。

首先查看mysql是否开启定时任务开关
SHOW VARIABLES LIKE ‘event_scheduler’;
Value为ON则已打开，OFF则关闭

如果是OFF，就先打开：SET GLOBAL event_scheduler = ON;
然后创建我们想要的定时器及删除存储过程

delimiter $$  
drop event if exists delplan;  
create event delplan
on schedule   
EVERY 1 DAY  
 STARTS '2019-07-02 00:00:00'  
ON COMPLETION  PRESERVE ENABLE  
do  
begin  
  delete from sys_prepose_sync_record where Sync_Time<(CURRENT_TIMESTAMP()-INTERVAL 497 HOUR);  
end $$  
delimiter;

最后把上面的语句在数据库执行一遍，大功告成，可以删除Sync_Time在497小时之前的数据了。

如果不想用该定时器了的话，可以直接在数据库事件中将状态更改为ENable即可。

使用navicat连接的数据库->事件->找到对应的事件修改
如果想要调整定时器执行时间间隔，可以直接在事件中修改

由于测试环境有张日志表没定时2分钟程序就狂插数据，导致不到1一个月时间，这张日志表就占用了6.7G的空间，但是日志刷新较快，有些日志就没什么作用，就写了个定时器，定期删除这张表的数据

    首先先查看mysql是否开启定时任务开关

    # SHOW VARIABLES LIKE 'event_scheduler';

---------------------



 

Value为ON则已打开，OFF则关闭

如果是OFF，就先打开：

# SET GLOBAL event_scheduler = ON;

然后创建我们想要的定时器

 

DELIMITER $$
DROP EVENT IF EXISTS deleteLog;
CREATE EVENT deleteLog
ON SCHEDULE EVERY 300 SECOND
ON COMPLETION PRESERVE
DO BEGIN
delete from ftp_log where TO_DAYS(now())-TO_DAYS(createOn)>2;
END$$
DELIMITER ；

 

该脚本的意思是：每300秒执行一次计划，执行的动作为删除两天前的数据

创建完成后，查看定时器

 # select * from mysql.event;


