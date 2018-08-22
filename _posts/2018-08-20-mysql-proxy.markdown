---
title: mysql-proxy
layout: post
category: storage
author: 夏泽民
---
https://downloads.mysql.com/archives/proxy/
<!-- more -->
mysql-proxy是官方提供的mysql中间件产品可以实现负载平衡，读写分离，failover等，但其不支持大数据量的分库分表且性能较差。下面介绍几款能代替其的mysql开源中间件产品，Atlas，cobar，tddl，让我们看看它们各自有些什么优点和新特性吧。

360 Atlas

Atlas是由 Qihoo 360, Web平台部基础架构团队开发维护的一个基于MySQL协议的数据中间层项目。它是在mysql-proxy 0.8.2版本的基础上，对其进行了优化，增加了一些新的功能特性。360内部使用Atlas运行的mysql业务，每天承载的读写请求数达几十亿条。
Altas架构：
Atlas是一个位于应用程序与MySQL之间，它实现了MySQL的客户端与服务端协议，作为服务端与应用程序通讯，同时作为客户端与MySQL通讯。它对应用程序屏蔽了DB的细节，同时为了降低MySQL负担，它还维护了连接池。



以下是一个可以参考的整体架构，LVS前端做负载均衡，两个Altas做HA,防止单点故障。


Altas的一些新特性：
1.主库宕机不影响读
主库宕机，Atlas自动将宕机的主库摘除，写操作会失败，读操作不受影响。从库宕机，Atlas自动将宕机的从库摘除，对应用没有影响。在mysql官方的proxy中主库宕机，从库亦不可用。
2.通过管理接口，简化管理工作，DB的上下线对应用完全透明，同时可以手动上下线。
3.自己实现读写分离
（1）为了解决读写分离存在写完马上就想读而这时可能存在主从同步延迟的情况，Altas中可以在SQL语句前增加 /*master*/ 就可以将读请求强制发往主库
7.平滑重启
通过配置文件中设置lvs-ips参数实现平滑重启功能，否则重启Altas的瞬间那些SQL请求都会失败。该参数前面挂接的lvs的物理网卡的ip，注意不是虚ip。平滑重启的条件是至少有两台配置相同的Atlas且挂在lvs之后。
source： https://github.com/Qihoo360/Atlas 

ALIBABA COBAR

Cobar是阿里巴巴（B2B）部门开发的一种关系型数据的分布式处理系统，它可以在分布式的环境下看上去像传统数据库一样为您提供海量数据服务。那么具体说说我们为什么要用它，或说cobar--能干什么？以下是我们业务运行中会存在的一些问题：
1.随着业务的进行数据库的数据量和访问量的剧增，需要对数据进行水平拆分来降低单库的压力，而且需要高效且相对透明的来屏蔽掉水平拆分的细节。
2.为提高访问的可用性，数据源需要备份。
3.数据源可用性的检测和failover。
4.前台的高并发造成后台数据库连接数过多，降低了性能，怎么解决。 
针对以上问题就有了cobar施展自己的空间了，cobar中间件以proxy的形式位于前台应用和实际数据库之间，对前台的开放的接口是mysql通信协议。将前台SQL语句变更并按照数据分布规则转发到合适的后台数据分库，再合并返回结果，模拟单库下的数据库行为。 

Cobar应用举例
应用架构：

应用介绍：
1.通过Cobar提供一个名为test的数据库，其中包含t1,t2两张表。后台有3个MySQL实例(ip:port)为其提供服务，分别为：A,B,C。
2.期望t1表的数据放置在实例A中，t2表的数据水平拆成四份并在实例B和C中各自放两份。t2表的数据要具备HA功能，即B或者C实例其中一个出现故障，不影响使用且可提供完整的数据服务。
cabar优点总结：
1.数据和访问从集中式改变为分布：
（1）Cobar支持将一张表水平拆分成多份分别放入不同的库来实现表的水平拆分
（2）Cobar也支持将不同的表放入不同的库
（3） 多数情况下，用户会将以上两种方式混合使用
注意！：Cobar不支持将一张表，例如test表拆分成test_1,test_2, test_3.....放在同一个库中，必须将拆分后的表分别放入不同的库来实现分布式。
2.解决连接数过大的问题。
3.对业务代码侵入性少。
4.提供数据节点的failover,HA：
(1)Cobar的主备切换有两种触发方式，一种是用户手动触发，一种是Cobar的心跳语句检测到异常后自动触发。那么，当心跳检测到主机异常，切换到备机，如果主机恢复了，需要用户手动切回主机工作，Cobar不会在主机恢复时自动切换回主机，除非备机的心跳也返回异常。
(2)Cobar只检查MySQL主备异常，不关心主备之间的数据同步，因此用户需要在使用Cobar之前在MySQL主备上配置双向同步。
cobar缺点：
开源版本中数据库只支持mysql，并且不支持读写分离。
source： http://code.alibabatech.com/wiki/display/cobar/Home 

TAOBAO TDDL

淘宝根据自己的业务特点开发了TDDL（Taobao Distributed Data Layer 外号:头都大了 ©_Ob）框架，主要解决了分库分表对应用的透明化以及异构数据库之间的数据复制，它是一个基于集中式配置的 jdbc datasource实现，具有主备，读写分离，动态数据库配置等功能。
TDDL所处的位置（tddl通用数据访问层，部署在客户端的jar包，用于将用户的SQL路由到指定的数据库中）：


淘宝很早就对数据进行过分库的处理， 上层系统连接多个数据库，中间有一个叫做DBRoute的路由来对数据进行统一访问。DBRoute对数据进行多库的操作、数据的整合，让上层系统像操作一个数据库一样操作多个库。但是随着数据量的增长，对于库表的分法有了更高的要求，例如，你的商品数据到了百亿级别的时候，任何一个库都无法存放了，于是分成2个、4个、8个、16个、32个……直到1024个、2048个。好，分成这么多，数据能够存放了，那怎么查询它？这时候，数据查询的中间件就要能够承担这个重任了，它对上层来说，必须像查询一个数据库一样来查询数据，还要像查询一个数据库一样快（每条查询在几毫秒内完成），TDDL就承担了这样一个工作。在外面有些系统也用DAL（数据访问层） 这个概念来命名这个中间件。
下图展示了一个简单的分库分表数据查询策略：

主要优点：
1.数据库主备和动态切换
2.带权重的读写分离
3.单线程读重试
4.集中式数据源信息管理和动态变更
5.剥离的稳定jboss数据源
6.支持mysql和oracle数据库
7.基于jdbc规范，很容易扩展支持实现jdbc规范的数据源
8.无server,client-jar形式存在，应用直连数据库
9.读写次数,并发度流程控制，动态变更
10.可分析的日志打印,日志流控，动态变更
TDDL必须要依赖diamond配置中心（diamond是淘宝内部使用的一个管理持久配置的系统，目前淘宝内部绝大多数系统的配置，由diamond来进行统一管理，同时diamond也已开源）。
TDDL动态数据源使用示例说明： http://rdc.taobao.com/team/jm/archives/1645 
diamond简介和快速使用： http://jm.taobao.org/tag/diamond%E4%B8%93%E9%A2%98/ 
TDDL源码： https://github.com/alibaba/tb_tddl 
TDDL复杂度相对较高。当前公布的文档较少，只开源动态数据源，分表分库部分还未开源，还需要依赖diamond，不推荐使用。
终其所有，我们研究中间件的目的是使数据库实现性能的提高。具体使用哪种还要经过深入的研究，严谨的测试才可决定。
二
、mysql-proxy实现读写分离

1、安装mysql-proxy

实现读写分离是有lua脚本实现的，现在mysql-proxy里面已经集成，无需再安装

下载：http://dev.mysql.com/downloads/mysql-proxy/

Shell

tar zxvf mysql-proxy-0.8.3-linux-glibc2.3-x86-64bit.tar.gz

mv mysql-proxy-0.8.3-linux-glibc2.3-x86-64bit /usr/local/mysql-proxy

2、配置mysql-proxy，创建主配置文件

Shell

cd /usr/local/mysql-proxy

mkdir lua #创建脚本存放目录

mkdir logs #创建日志目录

cp share/doc/mysql-proxy/rw-splitting.lua ./lua #复制读写分离配置文件

cp share/doc/mysql-proxy/admin-sql.lua ./lua #复制管理脚本

vi /etc/mysql-proxy.cnf #创建配置文件

[mysql-proxy]

user=root #运行mysql-proxy用户

admin-username=proxy #主从mysql共有的用户

admin-password=123.com #用户的密码

proxy-address=192.168.0.204:4000 #mysql-proxy运行ip和端口，不加端口，默认4040

proxy-read-only-backend-addresses=192.168.0.203 #指定后端从slave读取数据

proxy-backend-addresses=192.168.0.202 #指定后端主master写入数据

proxy-lua-script=/usr/local/mysql-proxy/lua/rw-splitting.lua #指定读写分离配置文件位置

admin-lua-script=/usr/local/mysql-proxy/lua/admin-sql.lua #指定管理脚本

log-file=/usr/local/mysql-proxy/logs/mysql-proxy.log #日志位置

log-level=info #定义log日志级别，由高到低分别有(error|warning|info|message|debug)

daemon=true #以守护进程方式运行

keepalive=true #mysql-proxy崩溃时，尝试重启

保存退出！

chmod 660 /etc/mysql-porxy.cnf

3、修改读写分离配置文件

Shell

vi /usr/local/mysql-proxy/lua/rw-splitting.lua

if not proxy.global.config.rwsplit then

proxy.global.config.rwsplit = {

min_idle_connections = 1, #默认超过4个连接数时，才开始读写分离，改为1

max_idle_connections = 1, #默认8，改为1

is_debug = false

}

end

4、启动mysql-proxy

Shell

/usr/local/mysql-proxy/bin/mysql-proxy --defaults-file=/etc/mysql-proxy.cnf

netstat -tupln | grep 4000 #已经启动

tcp 0 0 192.168.0.204:4000 0.0.0.0:* LISTEN 1264/mysql-proxy

关闭mysql-proxy使用：killall -9 mysql-proxy

5、测试读写分离

1>.在主服务器创建proxy用户用于mysql-proxy使用，从服务器也会同步这个操作

Shell

mysql> grant all on *.* to 'proxy'@'192.168.0.204' identified by '123.com';

2>.使用客户端连接mysql-proxy

Shell

mysql -u proxy -h 192.168.0.204 -P 4000 -p123.com

创建数据库和表，这时的数据只写入主mysql，然后再同步从slave，可以先把slave的关了，看能不能写入，这里我就不测试了，下面测试下读的数据！

Shell

mysql> create table user (number INT(10),name VARCHAR(255));

mysql> insert into test values(01,'zhangsan');

mysql> insert into user values(02,'lisi');

3>.登陆主从mysq查看新写入的数据如下，

Shell

mysql> use test;

Database changed

mysql> select * from user;

+--------+----------+

| number | name |

+--------+----------+

| 1 | zhangsan |

| 2 | lisi |

+--------+----------+

4>.再登陆到mysql-proxy，查询数据，看出能正常查询

Shell

mysql -u proxy -h 192.168.0.204 -P 4000 -p123.com

mysql> use test;

mysql> select * from user;

+--------+----------+

| number | name |

+--------+----------+

| 1 | zhangsan |

| 2 | lisi |

+--------+----------+

5>.登陆从服务器关闭mysql同步进程，这时再登陆mysql-proxy肯定会查询不出数据

Shell

slave stop；

6>.登陆mysql-proxy查询数据，下面看来，能看到表，查询不出数据

Shell

mysql> use test;

Database changed

mysql> show tables;

+----------------+

| Tables_in_test |

+----------------+

| user |

+----------------+

mysql> select * from user;

ERROR 1146 (42S02): Table 'test.user' doesn't exist

配置成功！真正实现了读写分离的效果！
所
使用服务器列表

tpl01 NAT 29. 159 mysql(源码） /etc/my.cnf 提供mysql服务（主） yw009 数据库主从+代理搭建 mysql 5.6.39

tpl02 NAT 29. 152 mysql(源码） /etc/my.cnf mysql服务（从） yw009 数据库主从+代理搭建 mysql 5.6.39

tpl03 NAT 29. 160 mysql-proxy /usr/local/mysql-proxy 代理服务 yw009 数据库主从+代理搭建 0.8.5

work NAT 29. 158 mysql客户端 yw009 数据库主从+代理搭建 \

步骤一：安装mysql-proxy

1)下载mysql-proxy 在github.com 上下载0.8.5版

16 tar -xvf mysql-proxy-rel-0.8.5.tar

18 yum -y install lua

54 yum install lua-devel

33 yum install cmake

34 yum install make

41 yum -y install gcc openssl-devel pcre-devel zlib-devel ncurses-devel

44 yum -y install libtool

63 yum -y install gcc gcc-c++

79 yum search flex

80 yum install flex.x86_64

51 yum install mysql-devel

68 yum -y install glib2-devel

70 yum -y install libevent-devel

94 ./autogen.sh

96 ./configure --prefix=/usr/local/mysql-proxy

make

make install

79 yum search flex

80 yum install flex.x86_64

126 mkdir /usr/local/mysql-proxy/lua

127 cp lib/rw-splitting.lua /usr/local/mysql-proxy/lua/

128 cp lib/admin-sql.lua /usr/local/mysql-proxy/lua/

131 cp -r lib/proxy /usr/local/mysql-proxy/lua/

132 cd /usr/local/mysql-proxy/

135 vim lua/rw-splitting.lua

2) 搭建数据库主从 Master (tpl01) ,Slave (tpl02) ；可参考其他篇博客

3）启动mysql-proxy服务

139 bin/mysql-proxy -P 192.168.29.160:3306 -b 192.168.29.159:3306 -r 192.168.29.152:3306 -s lua/rw-splitting.lua &

启动后可确认监听状态：

netstat -anptu | grep mysql;

为了每次开机后能够自动运行mysql-proxy,可以将相关操作写到/etc/rc.local配置文件内：

步骤二： 测试读写分离

1) 在MySQL Master服务器上设置用户授权

以root 用户为例，允许其从192.168.29.0/24 网段的客户机远程访问。首先登入Master服务器添加下列授权：

mysql > grant all on*.* to root@'192.168.29.%' identified by '123qwe';

因为此前已配置mysql库的主从同步，Slave上的root授权会自动更新：

2）从客户机work 访问MySQL数据库

注意连接的是mysql-proxy服务器，而并不是Master或 Slave:

测试数据库写入操作：

mysql > create database proxydb;

mysql > use proxydb;

mysql > create table proxytb( id int(4), host varchar(48));

mysql > insert into proxytb values(1, 'aa'), (2, 'bb');

mysql > select * from proxytb;

3) 在 Master 和 Slave 确认新建的库，表

4） 观察MySQL 代理访问的网络连接

在 Master上可看到来自Slave 和proxy代理的网路连接：

netstat -anptu | grep mysql

在Proxy代理上 可以看到与MySQL读，写服务器的网络连接：

netstat -anptu | grep mysql

5) 怎么才能确定读的数据是确实无疑 来自从数据库呢？

可以 使用 root 用户登入 slave 修改一条记录（这样主从数据库数据就不一致了，我们方便观察）

mysql> update proxytb set host='ee' where id=2;

在 客户机work 上 连接mysql 查看数据：

mysql > select * from proxytb;
