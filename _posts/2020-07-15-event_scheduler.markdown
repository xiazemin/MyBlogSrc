---
title: event_scheduler mysql定时器定时清理表数据
layout: post
category: storage
author: 夏泽民
---
1.首先先查看mysql是否开启定时任务开关

    SHOW VARIABLES LIKE 'event_scheduler';
2.Value为ON则已打开，OFF则关闭

 如果是OFF，就先打开：

#  SET GLOBAL event_scheduler = ON;
<!-- more -->
3.然后创建我们想要的定时器

DELIMITER
DROPEVENTIFEXISTSdeleteLog;CREATEEVENTdeleteLogONSCHEDULEEVERY300SECONDONCOMPLETIONPRESERVEDOBEGINdeletefromtsgcloginlogwhereTODAYS(now())−TODAYS(logintm)>2;END
DROPEVENTIFEXISTSdeleteLog;CREATEEVENTdeleteLogONSCHEDULEEVERY300SECONDONCOMPLETIONPRESERVEDOBEGINdeletefromtsgcloginlogwhereTODAYS(now())−TODAYS(logintm)>2;END
 
DELIMITER ; 

该脚本的意思是：每300秒执行一次计划，执行的动作为删除两天前的数据

CREATE PROCEDURE `prc_del_um_time_stat`(IN date_inter int) COMMENT '自动删除日志'
BEGIN
delete from um_time_stat where (TO_DAYS(NOW()) - TO_DAYS(create_time))>=date_inter;
END;

CREATE EVENT `auto_delete_um_time_stat` ON SCHEDULE EVERY 1 DAY STARTS '2019-05-08 00:00:00'
ON COMPLETION PRESERVE ENABLE COMMENT '自动删除30天以前的日志统计数据' DO call prc_del_um_time_stat(30);

该脚本的意思是：每天凌晨执行一次计划，执行的动作为删除30天前的数据


4.创建完成后，查看定时器(需要root用户)

select * from  mysql.event;

https://www.cnblogs.com/c840136/articles/2388512.html