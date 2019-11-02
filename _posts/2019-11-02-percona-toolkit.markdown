---
title: percona-toolkit
layout: post
category: storage
author: 夏泽民
---
MySQL数据库是轻量级、开源数据库的佼佼者，其功能和管理，健壮性与Oracle相比还是有相当的差距。因此有很多功能强大第三方的衍生产品，如percona-toolkit，XtraBackup等等。percona-toolkit是一组高级命令行工具的集合，可以查看当前服务的摘要信息，磁盘检测，分析慢查询日志，查找重复索引，实现表同步等等。这个工具套件对DBA及运维人员着实不可多得。
https://www.percona.com/doc/percona-toolkit/2.1/pt-heartbeat.html
https://github.com/percona/percona-toolkit
<!-- more -->
1、percona-toolkit的主要功能
  Verify MySQL replication integrity by checking master and replica data consistency
  Efficiently archive rows
  Find duplicate indexes
  Summarize MySQL servers
  Analyze queries from logs and tcpdump
  Collect vital system information when problems occur

 

2、安装需求及步骤
  下载链接：http://www.percona.com/software/percona-toolkit
  需求
    * Perl v5.8 or newer
    * Bash v3 or newer
    * Core Perl modules like Time::HiRes
    
    # perl --version |head -2 #检查perl版本
    # bash --version          #检查bash版本

  快速安装步骤(缺省/usr/local/bin路径下，过程略)
    # tar zxvf percona-toolkit-<version>.tar.gz
    # cd percona-toolkit-<version>
    # perl Makefile.PL   (安装到非缺省目录 perl Makefile.PL PREFIX=${HOME})
    # make
    # make test
    # make install

 

3、主要工具介绍
  如果是非源码安装或源码安装是未指定路径，缺省情况下所有的pt相关的工具位于/usr/bin目录下，以pt-开头。
  获取有关命令行的帮助信息，直接在shell提示符下输入命令行与--hlep即可。如： /usr/bin/pt-upgrade --help
  # ls -hltr /usr/bin/pt-*
  pt-upgrade
    #该命令主要用于对比不同mysql版本下SQL执行的差异，通常用于升级前进行对比。
    #会生成SQL文件或单独的SQL语句在每个服务器上执行的结果、错误和警告信息等。
  
  pt-online-schema-change
    #功能为支持在线变更表构，且不锁定原表，不阻塞原表的DML操作。
    #该特性与Oracle的dbms_redefinition在线重定义表原理基本类似。
    
  pt-mysql-summary
    #对连接的mysql服务器生成一份详细的配置情况以及sataus信息
    #在尾部也提供当前实例的的配置文件的信息
    
  pt-mext
    #并行查看SHOW GLOBAL STATUS的多个样本的信息。
    #pt-mext会执行你指定的COMMAND，并每次读取一行结果，把空行分割的内容保存到一个一个的临时文件中，最后结合这些临时文件并行查看结果。
  
  pt-kill
    #Kill掉符合指定条件mysql语句
  
  pt-ioprofile
    #pt-ioprofile的原理是对某个pid附加一个strace进程进行IO分析
    
  pt-fingerprint
    #用于生成查询指纹。主要将将sql查询生成queryID，pt-query-digest中的ID即是通过此工具来完成的。
    #类似于Oracle中的SQL_ID，涉及绑定变量，字面量等
  
  pt-find
    #用与查找mysql表并执行指定的命令，类似于find命令
  
  pt-fifo-split
    #模拟切割文件并通过管道传递给先入先出队列而不用真正的切割文件
    
  pt-deadlock-logger
    #用于监控mysql服务器上死锁并输出到日志文件，日志包含发生死锁的时间、死锁线程id、死锁的事务id、发生死锁时事务执行时间等详细信息。
  
  pt-archiver
    #将mysql数据库中表的记录归档到另外一个表或者文件
    #该工具具只是归档旧的数据，对线上数据的OLTP查询几乎没有影响。
    #可以将数据插入另外一台服务器的其他表中，也可以写入到一个文件中，方便使用load data infile命令导入数据。
  
  pt-agent
    #基于Percona Cloud的一个客户端代理工具
  
  pt-visual-explain
    #用于格式化explain的输出
  
  pt-variable-advisor
    #用于分析mysql系统变量可能存在的一些问题，可以据此评估有关参数的设置正确与否。
  
  pt-stalk
    #用于收集mysql数据库故障时的相关信息便于后续诊断处理。
    
  pt-slave-delay
    #用于设定从服务器落后于主服务器的时间间隔。
    #该命令行通过启动和停止复制sql线程来设置从落后于主指定时间。

  pt-sift
    #用于浏览pt-stalk生成的文件。
    
  pt-show-grants
    #将当前实例的用户权限全部输出，可以用于迁移数据库过程中重建用户。
    
  pt-query-digest
    #用于分析mysql服务器的慢查询日志，并格式化输出以便于查看和分析。
  
  pt-pmp
    #为查询程序执行聚合的GDB堆栈跟踪，先进性堆栈跟踪，然后将跟踪信息汇总。
    
  pt-index-usage
    #从log文件中读取查询语句，并用分析当前索引如何被使用。
    #完成分析之后会生成一份关于索引没有被查询使用过的报告，可以用于分析报告考虑剔除无用的索引。
    
  pt-heartbeat
    #用于监控mysql复制架构的延迟。
    #主要是通过在主库上的--update线程持续更新指定表上的一个时间戳，从库上--monitor线程或者--check线程检查主库更新的时间戳并与当前系统时间对比，得到延迟值。
  
  pt-fk-error-logger
    #将外键相关的错误信息记录到日志或表。
  
  pt-duplicate-key-checker
     #功能为从mysql表中找出重复的索引和外键，这个工具会将重复的索引和外键都列出来
     #同时也可以生成相应的drop index的语句
    
  pt-diskstats
    #类似于iostat，打印磁盘io统计信息，但是这个工具是交互式并且比iostat更详细。可以分析从远程机器收集的数据。
  
  pt-config-diff
    #用于比较mysql配置文件和服务器变量
    #至少2个配置源需要指定，可以用于迁移或升级前后配置文件进行对比
    
  pt-align
    #格式化输出
    
  pt-slave-find
    #连接mysql主服务器并查找其所有的从，然后打印出所有从服务器的层级关系。
    
  pt-table-checksum
    #用于校验mysql复制的一致性。
    #该工具主要是高效的查找数据差异，如果存在差异性，可以通过pt-table-sync来解决。
MySQL的同步模式有三种，从官方文档可知Replication：

With asynchronous replication, the master writes events to its binary log and slaves request them when they are ready. There is no guarantee that any event will ever reach any slave.
With fully synchronous replication, when a master commits a transaction, all slaves also will have committed the transaction before the master returns to the session that performed the transaction. The drawback of this is that there might be a lot of delay to complete a transaction.
Semisynchronous replication falls between asynchronous and fully synchronous replication. The master waits only until at least one slave has received and logged the events. It does not wait for all slaves to acknowledge receipt, and it requires only receipt, not that the events have been fully executed and committed on the slave side.
在asynchronous和Semisynchronous模式下，主库与从库之间必然存在一定延迟，当延迟大的话，从库的查询就可能查询到旧的数据，或者查询不到数据（比如主库插入的数据尚未同步到从库）。

那么如何判断主库与从库之间是否存在延迟呢？延迟又是多少？当我们在从库执行show slave status\G查看复制状态时，其中有一个字段是Seconds_Behind_Master，从字面上理解其意思是当前从库落后主库的秒数。

上Google搜索下关键词MySQL Seconds_Behind_Master，可以看到很多博文不推荐通过Seconds_Behind_Master去监控数据库是否存在延迟。先看下MySQL官方文档关于Seconds_Behind_Master的解释：

When the slave is actively processing updates, this field shows the difference between the current timestamp on the slave and the original timestamp logged on the master for the event currently being processed on the slave.
When no event is currently being processed on the slave, this value is 0.
Seconds_Behind_Master的计算规则是：从库当前的时间戳 - 从库SQL线程当前正在执行的事件记录的主库的时间戳
当从库sql线程上没有事件被执行时，Seconds_Behind_Master的值为0.
In essence, this field measures the time difference in seconds between the slave SQL thread and the slave I/O thread. If the network connection between master and slave is fast, the slave I/O thread is very close to the master, so this field is a good approximation of how late the slave SQL thread is compared to the master. If the network is slow, this is not a good approximation; the slave SQL thread may quite often be caught up with the slow-reading slave I/O thread, so Seconds_Behind_Master often shows a value of 0, even if the I/O thread is late compared to the master. In other words, this column is useful only for fast networks.

如果从库与主库之间的网络传输速度很快，即IO线程能一直从主库获取到最新的二进制日志，那么Seconds_Behind_Master其实也就表示了SQL线程和IO线程之间的时间差。 如果网络很慢，且SQL线程总是能跟上IO线程，那么Seconds_Behind_Master就会经常出现0值（即此时IO线程仍在接收主库传来的二进制日志，而SQL线程已经把已收到的主库二进制日志处理完了），即使IO线程正在接收的二进制日志并非是最新的。

This time difference computation works even if the master and slave do not have identical clock times, provided that the difference, computed when the slave I/O thread starts, remains constant from then on. Any changes—including NTP updates—can lead to clock skews that can make calculation of Seconds_Behind_Master less reliable.

Seconds_Behind_Master的计算不受主库和从库之间时间不一致的影响，主库和从库之间的时间差在IO线程启动时就计算好了。任何对主库、从库时间调整的操作都可能导致Seconds_Behind_Master的计算结果不可靠。所以Seconds_Behind_Master的真正计算规则是：从库当前的时间戳 - 从库SQL线程当前正在执行的事件记录的主库的时间戳 - 从库与主库的时间差

This field is NULL (undefined or unknown) if the slave SQL thread is not running, or if the slave I/O thread is not running or is not connected to the master. For example, if the slave I/O thread is running but is not connected to the master and is sleeping for the number of seconds given by the CHANGE MASTER TO statement or –master-connect-retry option (default 60) before reconnecting, the value is NULL. This is because the slave cannot know what the master is doing, and so cannot say reliably how late it is.

Seconds_Behind_Master的有可能为NULL值，当SQL线程或IO线程没有在执行，或者IO线程在执行但连不上主库时。

The value of Seconds_Behind_Master is based on the timestamps stored in events, which are preserved through replication. This means that if a master M1 is itself a slave of M0, any event from M1’s binary log that originates from M0’s binary log has M0’s timestamp for that event. This enables MySQL to replicate TIMESTAMP successfully. However, the problem for Seconds_Behind_Master is that if M1 also receives direct updates from clients, the Seconds_Behind_Master value randomly fluctuates because sometimes the last event from M1 originates from M0 and sometimes is the result of a direct update on M1.

Seconds_Behind_Master的计算是基于二进制日志中记录的事件的时间戳，所以在M0->M1->S0这样的主从结构中，当M0的二进制日志传递到S0时，其二进制日志中的事件的时间戳是保持不变的。

我试着搜索了下MySQL的源码，找到了Seconds_Behind_Master的计算相关的代码，在sql/slave.cc中：

    /*
      Seconds_Behind_Master: if SQL thread is running and I/O thread is
      connected, we can compute it otherwise show NULL (i.e. unknown).
    */
    if ((mi->slave_running == MYSQL_SLAVE_RUN_CONNECT) &&
        mi->rli.slave_running)
    {
      long time_diff= ((long)(time(0) - mi->rli.last_master_timestamp)
                       - mi->clock_diff_with_master);
      /*
        Apparently on some systems time_diff can be <0. Here are possible
        reasons related to MySQL:
        - the master is itself a slave of another master whose time is ahead.
        - somebody used an explicit SET TIMESTAMP on the master.
        Possible reason related to granularity-to-second of time functions
        (nothing to do with MySQL), which can explain a value of -1:
        assume the master's and slave's time are perfectly synchronized, and
        that at slave's connection time, when the master's timestamp is read,
        it is at the very end of second 1, and (a very short time later) when
        the slave's timestamp is read it is at the very beginning of second
        2. Then the recorded value for master is 1 and the recorded value for
        slave is 2. At SHOW SLAVE STATUS time, assume that the difference
        between timestamp of slave and rli->last_master_timestamp is 0
        (i.e. they are in the same second), then we get 0-(2-1)=-1 as a result.
        This confuses users, so we don't go below 0: hence the max().

        last_master_timestamp == 0 (an "impossible" timestamp 1970) is a
        special marker to say "consider we have caught up".
      */
      protocol->store((longlong)(mi->rli.last_master_timestamp ?
                                 max(0, time_diff) : 0));
    }
    else
    {
      protocol->store_null();
    }
可见Seconds_Behind_Master的计算规则正如文档中所描述，其值可能为NULL、0、大于0。

既然使用Seconds_Behind_Master去判断主从延迟不可靠，那么用什么办法好呢？Google上找到的文章大多提到了percona工具中的pt-heartbeat。

其工作原理为：

在主库上创建一张表，然后由一个进程间隔一定时间不断地去更新这张表记录的当前主库的时间戳
由另一个进程间隔一定时间不断地去从库查询这张表记录的主库时间戳，再跟当前时间戳相减得到主从复制的延迟
现在我准备使用pt-heartbeat来做主从延迟监控，但得先在开发环境部署下，看看是否有遇到什么坑。这里我使用了MySQL Docker镜像来搭建测试环境很方便。

运行如下两个命令，我就启动了一个主库和从库了：

docker run -v /root/test/my.cnf:/etc/my.cnf --name master -e MYSQL_ROOT_PASSWORD=123456 -d mysql/mysql-server:5.5
docker run -v /root/test/slave.cnf:/etc/my.cnf --name slave -e MYSQL_ROOT_PASSWORD=123456 -d mysql/mysql-server:5.5
进入docker实例：

docker exec -it master  bash
mysql -uroot -p123456
这样就可以开始进行复制的配置了，主从配置好后，就需要配置pt-heart。

以守护进程的方式运行，定时更新记录的时间戳：

pt-heartbeat -D delay_test --update -h 172.17.0.2 -u test_dev -p 123 --daemonize --create-table
以守护进程的方式运行，计算主从之间的延迟，并将结果输出到delay.log文件中：

pt-heartbeat -D delay_test --monitor --frames 5s,15s,30s -h 172.17.0.3 -u test_dev --daemonize -p 123 > delay.log
关于使用到的参数，请自行查看pt-heartbeat文档。这里有一点需要提醒的是，如果执行pt-heartbeat -D delay_test --monitor --frames 5s,15s,30s -h 172.17.0.3 -u test_dev -p 123 > delay.log，即不以守护进程运行，且把标准输出重定向到文件delay.log中，当我们执行tail -f delay.log时，并不能看到一行行的输出记录，而会看到一段段的输出，每一段之间的输出会相隔挺久的。这是因为输出到文件delay.log的内容被缓冲起来了，如果执行``pt-heartbeat -D delay_test –monitor –frames 5s,15s,30s -h 172.17.0.3 -u test_dev -p 123`，我们会看到记录一行行地在终端上输出：

bash# ./pt-heartbeat -D delay_test --monitor --frames 5s,15s,30s -h 172.17.0.3 -u test_dev -p 123
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
0.00s [  0.00s,  0.00s,  0.00s ]
这是因为标准输出默认是行缓冲，关于IO缓冲可自行查看《UNIX环境高级编程》P116 或Google一下。但我发现执行pt-heartbeat -D delay_test --monitor --frames 5s,15s,30s -h 172.17.0.3 -u test_dev --daemonize -p 123 > delay.log，即作为守护进程运行时，tail -f delay.log看到的却是一行行的输出记录，这是为何？按理说应该是一段段的输出才对，于是粗略看了下pt-heartbeat的源码，在其daemonize函数中，有这样一行代码：$OUTPUT_AUTOFLUSH = 1;，这行代码就是用于开启自动冲洗流，相关文档：Special Variables

这里还有个问题，pt-heartbeat的输出是不带时间的，可以简单修改下pt-heartbeat的源码：

简单的更改，第6243行改为如下：

my $output = time . ' ' . sprintf $format, $delay, @vals, $pk_val;
严格来说，应该这样改，即使用计算延迟的时间戳：

//get_delay函数的返回值，第6108行
$get_delay = sub {
    my ($sth) = @_;
    $sth->execute();
    PTDEBUG && _d($sth->{Statement});
    my ($ts, $hostname, $server_id) = $sth->fetchrow_array();
    my $now = time;
    PTDEBUG && _d("Heartbeat from server", $server_id, "\n",
    " now:", ts($now, $utc), "\n",
    "  ts:", $ts, "\n",
    "skew:", $skew);
    my $delay = $now - unix_timestamp($ts, $utc) - $skew;
    PTDEBUG && _d('Delay', sprintf('%.6f', $delay), 'on', $hostname);

    # Because we adjust for skew, if the ts are less than skew seconds
    # apart (i.e. replication is very fast) then delay will be negative. 
    # So it's effectively 0 seconds of lag.
    $delay = 0.00 if $delay < 0; 

    $sth->finish();
    return ($delay, $hostname, $pk_val, $now);// 多返回一个值$now，用于计算延迟的时间戳
};

//
if ( $o->get('monitor') ) {
    $heartbeat_sth ||= $dbh->prepare($heartbeat_sql);
    my ($delay, $hostname, $master_server_id, $curr_time) = $get_delay->($heartbeat_sth);//增加$curr_time变量接收

    unshift @samples, $delay;
    pop @samples if @samples > $limit;

    # Calculate and print results
    my @vals = map {
       my $bound = min($_, scalar(@samples));
       sum(@samples[0 .. $bound-1]) / $_;
    } @$frames;

    my $output = $curr_time . ' ' . sprintf $format, $delay, @vals, $pk_val;//将时间戳加入到输出结果中

    ......
}
现在的输出结果如下：

# ./pt-heartbeat -D delay_test --monitor --frames 5s,15s,30s -h 172.17.0.3 -u test_dev -p 123
1480511974.5048 0.00s [  0.00s,  0.00s,  0.00s ]
1480511975.5018 0.00s [  0.00s,  0.00s,  0.00s ]
1480511976.50172 0.00s [  0.00s,  0.00s,  0.00s ]
1480511977.50172 0.00s [  0.00s,  0.00s,  0.00s ]
1480511978.50182 0.00s [  0.00s,  0.00s,  0.00s ]
监控日志有了，但最好能够可视化的展示出来，这里我选择使用ELK的Kibana来做展示。logstash读取监控日志，入到ELK后，在kibana上创建面板观察。

logstash的配置脚本如下：

input {
        file {
                path => "/path_to_log/delay-*.log"
        }
}
filter {
        # add host info
        mutate {
                remove_field => ["host"]
        }
        # 我是在一台机器上启动多个pt-heartbeat进程监测多台从库的，log文件名为delay-{从库ip}.log
        grok {
                match => { "path" => "%{GREEDYDATA}/delay-%{GREEDYDATA:host}\.log" }
        }
        # 解析监测的日志文件
        grok {
                match => [ "message", "%{NUMBER:timestamp:int} %{NUMBER:currentDelay:float}s \[\s+%{NUMBER:SecondDelay5:float}s,\s+%{NUMBER:SecondDelay15:float}s,\s+%{NUMBER:SecondDelay30:float}s\s+\]" ]
        }
        # 以监测日志中的时间戳为准
        date {
                match => [ "timestamp", "UNIX" ]
                remove_field => [ "timestamp" ]
        }
        mutate {
                remove_field => ["@version", "message", "tags", "path"]
        }
}
output {
        elasticsearch {
                hosts => ["es.com:80"]
                index => "log_mysql_repl_delay-%{+yyyy.MM.dd}"
                document_type => "%{+YYYY.MM.dd}"
                workers => 2
                idle_flush_time => 5
        }
        # 测试使用
        #stdout {
        #       codec => rubydebug
        #}
}
制作的面板如下：



从图中可以看出，有一台从库确实存在较多的延迟，线上运营也确实有人反馈过一些关于从库延迟的问题（比如发表了回答后刷新页面看不到自己的回答）。这台从库其实一直有磁盘负载高的告警，但一直没在意，很有可能是这个问题导致了这台从库复制存在较多延迟，接下来得解决下这台从库的磁盘告警问题了
