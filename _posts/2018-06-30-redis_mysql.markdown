---
title: redis_mysql
layout: post
category: storage
author: 夏泽民
---
redis的介绍：
Redis的主从复制功能非常强大，一个master可以拥有多个slave，而一个slave又可以拥有多个slave，如此下去，形成了强大的多级服务器集群架构。下面是关于redis主从复制的一些特点：
1.master可以有多个slave
2.除了多个slave连到相同的master外，slave也可以连接其他slave形成图状结构
3.主从复制不会阻塞master。也就是说当一个或多个slave与master进行初次同步数据时，master可以继续处理client发来的请求。相反slave在初次同步数据时则会阻塞不能处理client的请求。
4.主从复制可以用来提高系统的可伸缩性,我们可以用多个slave 专门用于client的读请求，比如sort操作可以使用slave来处理。也可以用来做简单的数据冗余
5.可以在master禁用数据持久化，只需要注释掉master 配置文件中的所有save配置，然后只在slave上配置数据持久化。

下面介绍下主从复制的过程
       当设置好slave服务器后，slave会建立和master的连接，然后发送sync命令。无论是第一次同步建立的连接还是连接断开后的重新连 接，master都会启动一个后台进程，将数据库快照保存到文件中，同时master主进程会开始收集新的写命令并缓存起来。后台进程完成写文件 后，master就发送文件给slave，slave将文件保存到磁盘上，然后加载到内存恢复数据库快照到slave上。接着master就会把缓存的命 令转发给slave。而且后续master收到的写命令都会通过开始建立的连接发送给slave。从master到slave的同步数据的命令和从 client发送的命令使用相同的协议格式。当master和slave的连接断开时slave可以自动重新建立连接。如果master同时收到多个 slave发来的同步连接命令，只会使用启动一个进程来写数据库镜像，然后发送给所有slave。

二、     配置
下面我演示下怎样在多台服务器上进行Redis数据主从复制。我假设有两台服务器，一台是Linux操作系统（局域网IP：192.168.1.4，master服务器），一台是Linux操作系统（局域网IP：192.168.1.5，slave服务器）

配置slave服务器很简单，只需要在配置文件(redis.conf)中加入如下配置
bind  192.168.1.5(从服务器,此处默认是127.0.0.1，请修改成本机的IP地址，要不然，客户端无法进行访问)

slaveof 192.168.1.4 6379  (映射到主服务器上)

如果是在一台机器上面配置主从关系，那么还需要修改从服务器的默认端口号，同样也在redis.conf中进行修改。

MYsql 介绍
Replication原理  
Mysql 的 Replication 是一个异步的复制过程，从一个MySQL节点（称之为Master）复制到另一个MySQL节点（称之Slave）。在 Master 与 Slave 之间的实现整个复制过程主要由三个线程来完成，其中两个线程（SQL 线程和 I/O 线程）在 Slave 端，另外一个线程（I/O 线程）在 Master 端。
  
要实现 MySQL 的 Replication ，首先必须打开 Master 端的 Binary Log，因为整个复制过程实际上就是 Slave 从 Master 端获取该日志然后再在自己身上完全顺序的执行日志中所记录的各种操作。
  
看上去MySQL的Replication原理非常简单，总结一下： 
     * 每个从仅可以设置一个主。 
    * 主在执行sql之后，记录二进制log文件（bin-log）。 
    * 从连接主，并从主获取binlog，存于本地relay-log，并从上次记住的位置起执行sql，一旦遇到错误则停止同步。 
   
从这几条Replication原理来看，可以有这些推论： 
     * 主从间的数据库不是实时同步，就算网络连接正常，也存在瞬间，主从数据不一致。 
    * 如果主从的网络断开，从会在网络正常后，批量同步。 
    * 如果对从进行修改数据，那么很可能从在执行主的bin-log时出现错误而停止同步，这个是很危险的操作。所以一般情况下，非常小心的修改从上的数据。 
    * 一个衍生的配置是双主，互为主从配置，只要双方的修改不冲突，可以工作良好。 
    * 如果需要多主的话，可以用环形配置，这样任意一个节点的修改都可以同步到所有节点。

配置过程如下
主从复制配置

步骤如下：

主服务器：从服务器ip地址分别为
192.168.145.222、192.168.145.226  
1、修改主服务器master:
vi /etc/my.cnf  
[mysqld]  
    log-bin=mysql-bin   #[必须]启用二进制日志  
    server-id=222      #[必须]服务器唯一ID，默认是1，一般取IP最后一段  
2、修改从服务器slave:
vi /etc/my.cnf  
[mysqld]  
     log-bin=mysql-bin   #[不是必须]启用二进制日志  
     server-id=226      #[必须]服务器唯一ID，默认是1，一般取IP最后一段  
3、重启两台服务器的mysql
systemctl restart mariadb  
4、在主服务器上建立帐户并授权slave:
mysql  
mysql>GRANT REPLICATION SLAVE ON *.* to 'mysync'@'%' identified by 'q123456'; //一般不用root帐号，“%”表示所有客户端都可能连，只要帐号，密码正确，此处可用具体客户端IP代替，如192.168.145.226，加强安全。  
5、登录主服务器的mysql，查询master的状态
mysql>show master status;  
 +------------------+----------+--------------+------------------+  
 | File             | Position | Binlog_Do_DB | Binlog_Ignore_DB |  
 +------------------+----------+--------------+------------------+  
 | mysql-bin.000004 |      308 |              |                  |  
 +------------------+----------+--------------+------------------+  
 1 row in set (0.00 sec)  

注：执行完此步骤后不要再操作主服务器MYSQL，防止主服务器状态值变化

6、配置从服务器Slave：
注意mysql-bin.000004和308是第五步中的File和Position
mysql>change master to master_host='192.168.145.222',master_user='mysync',master_password='q123456',master_log_file='mysql-bin.000004',master_log_pos=308; //注意mysql-bin.000004和308是第五步中的File和  
mysql>start slave; //启动从服务器复制功能  
7、检查从服务器复制功能状态：
mysql> show slave status\G  
************************** 1. row ***************************  
  
            Slave_IO_State: Waiting for master to send event  
            Master_Host: 192.168.2.222  //主服务器地址  
            Master_User: mysync   //授权帐户名，尽量避免使用root  
            Master_Port: 3306    //数据库端口，部分版本没有此行  
            Connect_Retry: 60  
            Master_Log_File: mysql-bin.000004  
            Read_Master_Log_Pos: 600     //#同步读取二进制日志的位置，大于等于Exec_Master_Log_Pos  
            Relay_Log_File: ddte-relay-bin.000003  
            Relay_Log_Pos: 251  
            Relay_Master_Log_File: mysql-bin.000004  
            Slave_IO_Running: Yes    //此状态必须YES  
            Slave_SQL_Running: Yes     //此状态必须YES  
                  ......  
注：Slave_IO及Slave_SQL进程必须正常运行，即YES状态，否则都是错误的状态(如：其中一个NO均属错误)。

以上操作过程，主从服务器配置完成。

 9、主从服务器测试：

主服务器Mysql，建立数据库，并在这个库中建表插入一条数据：
mysql> create database hi_db;  
  Query OK, 1 row affected (0.00 sec)  
  
  mysql> use hi_db;  
  Database changed  
  
  mysql>  create table hi_tb(id int(3),name char(10));  
  Query OK, 0 rows affected (0.00 sec)  
   
  mysql> insert into hi_tb values(001,'bobu');  
  Query OK, 1 row affected (0.00 sec)  
  
  mysql> show databases;  
   +--------------------+  
   | Database           |  
   +--------------------+  
   | information_schema |  
   | hi_db                |  
   | mysql                |  
   | test                 |  
   +--------------------+  
   4 rows in set (0.00 sec)  
  
从服务器Mysql查询：  
  
   mysql> show databases;  
  
   +--------------------+  
   | Database               |  
   +--------------------+  
   | information_schema |  
   | hi_db                 |       //I'M here，大家看到了吧  
   | mysql                 |  
   | test          |  
  
   +--------------------+  
   4 rows in set (0.00 sec)  
  
   mysql> use hi_db  
   Database changed  
   mysql> select * from hi_tb;           //查看主服务器上新增的具体数据  
   +------+------+  
   | id   | name |  
   +------+------+  
   |    1 | bobu |  
   +------+------+  
   1 row in set (0.00 sec)  

<!-- more -->
实现数据同步的方式有以下几种：
1. 代码的方式控制，但是这种方式缺点是需要写入数据库的同时，也要向redis中写数，造成代码冗余和业务领域编程人员无法专注与业务逻辑。
2. 配置mybatis拦截器，逻辑上是可行的，但是要考虑到只是实时数据可以插入或者删除，但是历史数据不能及时同步到redis中，尤其是代码异常，数据不能及时同步到redis中。
3. 使用mysql的udf，通俗来说来说就是利用表的trigger，自动触发函数库，实现对redis进行操作，但是一般项目都是java居多，调c++代码显得很麻烦。
4. 使用redis作为mysql的二级缓存，实现org.apache.ibatis.cache.Cache接口写个MybatisRedisCache这样的类。修改mysql数据后，可以直接刷新redis缓存数据。比第二种稍好些。
5. 将redis作为关系型数据库的从库，从可行性和逻辑作用来看，在这种使用场景之下，确实是一种主从的关系，redis负责取数，传统数据库负责对数据的增删改查。实现原理是模拟数据库的主从同步机制，读取数据库二进制日志文件，将数据同步到redis中。canal-阿里开源项目，实现了mysql和redis的逻辑上主从同步。
mysql被oracle收购了，所以以后可能会收费，他的另一个版本mariadb开发者的另一个开源版本，他是从mysql衍生出来的，所以他在某种意义上说，就是mysql。配置好my.ini：
[client]
port = 7020
[mysqld]
character_set_server = utf8
server_id = 90
port = 7020
log_bin = mysql_bin
binlog-format = ROW
sql_mode = NO_ENGIN_SUBSTITUTION_TRANS_TABLES
    配置startup.bat文件：

cd bin
mysql --defaults-file=../my.ini --user=root
    点击startup.bat启动，创建canal用户，赋予服务器和sql权限：

        1. 创建数据库表:tb_person_file

        2. 服务器权限select,replication client,replication slave付给canal，并且赋予一定的增删改查权限。

    重启动服务器，startup.bat 

2. 开启redis,本机使用docker运行redis镜像
使用canal.deployer这个版本就可以，配置好canal，将数据信息配置到canal/conf/example/instance.properties

#################################################
## mysql serverId
canal.instance.mysql.slaveId=2020
# position info
canal.instance.master.address=192.168.1.109:7020
canal.instance.master.journal.name=
canal.instance.master.position=
canal.instance.master.timestamp=
 这里配置的是我的数据库所在的服务器，启动canal

 5. 启动canal-client，与canal服务连接，将数据同步到redis中：

向数据库中插入一体数据。通过加断点的方式，定位到：

调用redis工具类，将数据保存到redis中

数据已经保存到redis里边。 这样，我们实现的mysql数据实时同步到redis中，解决了数据需要重复使用并且访问关系型数据库耗时的问题。
使用canal的优势在于模拟了数据库的主从同步模式，通过mysql中数据日志对数据进行同步：

    1. 数据的实时同步，包括增删改查等操作。

    2. 服务重启以后，根据之前的日志文件，将历史数据同步到redis中。
