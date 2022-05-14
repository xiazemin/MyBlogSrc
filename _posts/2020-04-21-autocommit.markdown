---
title: autocommit
layout: post
category: storage
author: 夏泽民
---
ysql中set autocommit=0与start transaction区别
set autocommit=0指事务非自动提交，自此句执行以后，每个SQL语句或者语句块所在的事务都需要显示"commit"才能提交事务。

1、不管autocommit 是1还是0 
     START TRANSACTION 后，只有当commit数据才会生效，ROLLBACK后就会回滚。
2、当autocommit 为 0 时
    不管有没有START TRANSACTION。
    只有当commit数据才会生效，ROLLBACK后就会回滚。
3、如果autocommit 为1 ，并且没有START TRANSACTION 。
    调用ROLLBACK是没有用的。即便设置了SAVEPOINT。
<!-- more -->
set autocommit =1;

DELIMITER $$
DROP PROCEDURE IF EXISTS  mmmm $$  
CREATE DEFINER=`tms_admin`@`%` PROCEDURE `mmmm`()
BEGIN


      
    DECLARE t_error INTEGER DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION,SQLWARNING,NOT FOUND  SET t_error=1;
    


    START TRANSACTION;
    SAVEPOINT p1;  

            update t_user set instance_seq_id = 'tt00006';
            update t_user set instance_seq_id000 = 'tt00007';


            IF t_error = 1 THEN 
                ROLLBACK to p1;
            ELSE 
                COMMIT; 
            END IF;

END$$
DELIMITER ;

CALL `tms_inst_tt00003`.`mmmm`();

set autocommit =0;

DELIMITER $$
DROP PROCEDURE IF EXISTS  mmmm $$  
CREATE DEFINER=`tms_admin`@`%` PROCEDURE `mmmm`()
BEGIN


      
    DECLARE t_error INTEGER DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION,SQLWARNING,NOT FOUND  SET t_error=1;
    


    START TRANSACTION;
    SAVEPOINT p1;  

            update t_user set instance_seq_id = 'tt00006';
            update t_user set instance_seq_id000 = 'tt00007';


            IF t_error = 1 THEN 
                ROLLBACK to p1;
            ELSE 
                COMMIT; 
            END IF;

END$$
DELIMITER ;

CALL `tms_inst_tt00003`.`mmmm`();

set autocommit =0;

DELIMITER $$
DROP PROCEDURE IF EXISTS  mmmm $$  
CREATE DEFINER=`tms_admin`@`%` PROCEDURE `mmmm`()
BEGIN


      
    DECLARE t_error INTEGER DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION,SQLWARNING,NOT FOUND  SET t_error=1;

    SAVEPOINT p1;  

            update t_user set instance_seq_id = 'tt00006';
            update t_user set instance_seq_id000 = 'tt00007';


            IF t_error = 1 THEN 
                ROLLBACK to p1;
            ELSE 
                COMMIT; 
            END IF;

END$$
DELIMITER ;

autocoomit是事务，根据mysql的文档如果等于1是立即提交。但在transction中只有遇到commit或rollback才提交

set autocommit=0,
当前session禁用自动提交事物，自此句执行以后，每个SQL语句或者语句块所在的事务都需要显示"commit"才能提交事务。

start transaction

指的是启动一个新事务。

 

     在默认的情况下，MySQL从自动提交（autocommit）模式运行，这种模式会在每条语句执行完毕后把它作出的修改立刻提交给数据库并使之永久化。事实上，这相当于把每一条语句都隐含地当做一个事务来执行。如果你想明确地执行事务，需要禁用自动提交模式并告诉MySQL你想让它在何时提交或回滚有关的修改。
执行事务的常用办法是发出一条START TRANSACTION（或BEGIN）语句挂起自动提交模式，然后执行构成本次事务的各条语句，最后用一条 COMMIT语句结束事务并把它们作出的修改永久性地记入数据库。万一在事务过程中发生错误，用一条ROLLBACK语句撤销事务并把数据库恢复到事务开 始之前的状态。

       START TRANSACTION语句"挂起"自动提交模式的含义是：在事务被提交或回滚之后，该模式将恢复到开始本次事务的 START TRANSACTION语句被执行之前的状态。（如果自动提交模式原来是激活的，结束事务将让你回到自动提交模式；如果它原来是禁用的，结束 当前事务将开始下一个事务。）
如果是autocommit模式  ，autocommit的值应该为 1 ，不autocommit 的值是 0 ；请在试验前 确定autocommit 的模式是否开启

set autocommit
指事务非自动提交，此句执行以后，每个SQL语句或者语句块所在的事务都需要显示"commit"才能提交事务。

小实验
做了一个小实验，打开一个窗口，开了一个连接，set autocommit = 0开冲。

发现怎么不起作用？？？
经过同事指点，又开了另外的一个窗口，建立了数据库连接。

发现这个新建立的窗口是看不到ppp这一行的，也就是看不到之前那个窗口的那次insert结果。


直到之前的那个窗口提交后（commit后），才能在这个窗口看得到。

思考
其实一个窗口对应着一个数据库连接，那么在同一个窗口里 = 在同一个连接里 = 在同一个事务里。
那么问题来了，事务的隔离性隔离的到底是什么。从上面看的话，隔离的是这个数据库的多个连接。也就是说事务和连接是一对一的关系。

tips
set autocommit 和 START TRANSACTION
不管autocommit 是1还是0 ，START TRANSACTION 后，只有当commit数据才会生效，ROLLBACK后就会回滚。
当autocommit 为 0 时，不管有没有START TRANSACTION。
只有当commit数据才会生效，ROLLBACK后就会回滚。
为什么不推荐set autocommit=0
如果使用set autocommit=0，如果数据库是长连接，这就导致接下来的查询都在事务中，出现了长事务，长事务意味着系统里面会存在很老的事务视图。由于这些事务随时可能访问数据库里面的任何数据，所以这个事务提交之前，数据库里面它可能用到的回滚记录都必须保留，这就会导致大量占用存储空间。除了对回滚段的影响，长事务还占用锁资源，也可能拖垮整个库。

如何避免长事务对业务的影响？
出自林晓斌老师的mysql45讲。
首先，从应用开发端来看：

确认是否使用了set autocommit=0。这个确认工作可以在测试环境中开展，把MySQL的general_log开起来，然后随便跑一个业务逻辑，通过general_log的日志来确认。一般框架如果会设置这个值，也就会提供参数来控制行为，你的目标就是把它改成1。

确认是否有不必要的只读事务。有些框架会习惯不管什么语句先用begin/commit框起来。我见过有些是业务并没有这个需要，但是也把好几个select语句放到了事务中。这种只读事务可以去掉。

业务连接数据库的时候，根据业务本身的预估，通过SET MAX_EXECUTION_TIME命令，来控制每个语句执行的最长时间，避免单个语句意外执行太长时间。（为什么会意外？在后续的文章中会提到这类案例）

其次，从数据库端来看：

监控 information_schema.Innodb_trx表，设置长事务阈值，超过就报警/或者kill；
Percona的pt-kill这个工具不错，推荐使用；
在业务功能测试阶段要求输出所有的general_log，分析日志行为提前发现问题；
如果使用的是MySQL 5.6或者更新版本，把innodb_undo_tablespaces设置成2（或更大的值）。如果真的出现大事务导致回滚段过大，这样设置后清理起来更方便。

set autocommit=0,

当前session禁用自动提交事物，自此句执行以后，每个SQL语句或者语句块所在的事务都需要显示"commit"才能提交事务。

start transaction

指的是启动一个新事务。

在默认的情况下，MySQL从自动提交（autocommit）模式运行，这种模式会在每条语句执行完毕后把它作出的修改立刻提交给数据库并使之永久化。事实上，这相当于把每一条语句都隐含地当做一个事务来执行。如果你想明确地执行事务，需要禁用自动提交模式并告诉MySQL你想让它在何时提交或回滚有关的修改。

执行事务的常用办法是发出一条START TRANSACTION（或BEGIN）语句挂起自动提交模式，然后执行构成本次事务的各条语句，最后用一条 COMMIT语句结束事务并把它们作出的修改永久性地记入数据库。万一在事务过程中发生错误，用一条ROLLBACK语句撤销事务并把数据库恢复到事务开 始之前的状态。

START TRANSACTION语句"挂起"自动提交模式的含义是：在事务被提交或回滚之后，该模式将恢复到开始本次事务的 START TRANSACTION语句被执行之前的状态。（如果自动提交模式原来是激活的，结束事务将让你回到自动提交模式；如果它原来是禁用的，结束 当前事务将开始下一个事务。）

如果是autocommit模式  ，autocommit的值应该为 1 ，不autocommit 的值是 0 ；请在试验前 确定autocommit 的模式是否开启