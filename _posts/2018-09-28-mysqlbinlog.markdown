---
title: mysqlbinlog 格式解析
layout: post
category: storage
author: 夏泽民
---
1.什么时候写binlog
在说明什么时候写binlog前，先简单介绍下binlog的用途。binlog是二进制日志文件，用于记录mysql的数据更新或者潜在更新(比如DELETE语句执行删除而实际并没有符合条件的数据)，在mysql主从复制中就是依靠的binlog。在mysql中开启binlog需要设置my.cnf中的log_bin参数，另外也可以通过binlog_do_db
指定要记录binlog的数据库和binlog_ignore_db指定不记录binlog的数据库。对运行中的mysql要启用binlog可以通过命令SET SQL_LOG_BIN=1来设置。设置完成，我们就可以来测试binlog了。

需要注意下innodb引擎中的redo/undo log与mysql binlog是完全不同的日志，它们主要有以下几个区别：

a）层次不同。redo/undo log是innodb层维护的，而binlog是mysql server层维护的，跟采用何种引擎没有关系，记录的是所有引擎的更新操作的日志记录。innodb的redo/undo log更详细的说明可以参见姜承尧的《mysql技术内幕-innodb存储引擎》一书中相关章节。
b）记录内容不同。redo/undo日志记录的是每个页的修改情况，属于物理日志+逻辑日志结合的方式（redo log物理到页，页内采用逻辑日志，undo log采用的是逻辑日志），目的是保证数据的一致性。binlog记录的都是事务操作内容，比如一条语句DELETE FROM TABLE WHERE i > 1之类的，不管采用的是什么引擎，当然格式是二进制的，要解析日志内容可以用这个命令mysqlbinlog -vv BINLOG。
c）记录时机不同。redo/undo日志在事务执行过程中会不断的写入;而binlog仅仅在事务提交后才写入到日志，之前描述有误，binlog是在事务最终commit前写入的，多谢anti-semicolon 指出。当然，binlog什么时候刷新到磁盘跟参数sync_binlog相关。
显然，我们执行SELECT等不涉及数据更新的语句是不会记binlog的，而涉及到数据更新则会记录。要注意的是，对支持事务的引擎如innodb而言，必须要提交了事务才会记录binlog。

binlog刷新到磁盘的时机跟sync_binlog参数相关，如果设置为0，则表示MySQL不控制binlog的刷新，由文件系统去控制它缓存的刷新，而如果设置为不为0的值则表示每sync_binlog次事务，MySQL调用文件系统的刷新操作刷新binlog到磁盘中。设为1是最安全的，在系统故障时最多丢失一个事务的更新，但是会对性能有所影响，一般情况下会设置为100或者0，牺牲一定的一致性来获取更好的性能。

通过命令SHOW MASTER LOGS可以看到当前的binlog数目。如下面就是我机器上的mysql的binlog情况，第一列是binlog文件名，第二列是binlog文件大小。可以通过设置expire_logs_days来指定binlog保留时间，要手动清理binlog可以通过指定binlog名字或者指定保留的日期，命令分别是:purge master logs to BINLOGNAME;和purge master logs before DATE;。

......
| mysql-bin.000018 |       515 |
| mysql-bin.000019 |       504 |
| mysql-bin.000020 |       107 |
+------------------+-----------+
2 binlog格式解析
2.1 binlog文件格式简介
binlog格式分为statement，row以及mixed三种，mysql5.5默认的还是statement模式，当然我们在主从同步中一般是不建议用statement模式的，因为会有些语句不支持，比如语句中包含UUID函数，以及LOAD DATA IN FILE语句等，一般推荐的是mixed格式。暂且不管这三种格式的区别，看看binlog的存储格式是什么样的。binlog是一个二进制文件集合，当然除了我们看到的mysql-bin.xxxxxx这些binlog文件外，还有个binlog索引文件mysql-bin.index。如官方文档中所写，binlog格式如下：

binlog文件以一个值为0Xfe62696e的魔数开头，这个魔数对应0xfe 'b''i''n'。
binlog由一系列的binlog event构成。每个binlog event包含header和data两部分。
header部分提供的是event的公共的类型信息，包括event的创建时间，服务器等等。
data部分提供的是针对该event的具体信息，如具体数据的修改。
从mysql5.0版本开始，binlog采用的是v4版本，第一个event都是format_desc event 用于描述binlog文件的格式版本，这个格式就是event写入binlog文件的格式。关于之前版本的binlog格式，可以参见http://dev.mysql.com/doc/internals/en/binary-log-versions.html
接下来的event就是按照上面的格式版本写入的event。
最后一个rotate event用于说明下一个binlog文件。
binlog索引文件是一个文本文件，其中内容为当前的binlog文件列表。比如下面就是一个mysql-bin.index文件的内容。
/var/log/mysql/mysql-bin.000019
/var/log/mysql/mysql-bin.000020
/var/log/mysql/mysql-bin.000021
接下来分析下几种常见的event，其他的event类型可以参见官方文档。event数据结构如下：

+=====================================+
| event  | timestamp         0 : 4    |
| header +----------------------------+
|        | type_code         4 : 1    |
|        +----------------------------+
|        | server_id         5 : 4    |
|        +----------------------------+
|        | event_length      9 : 4    |
|        +----------------------------+
|        | next_position    13 : 4    |
|        +----------------------------+
|        | flags            17 : 2    |
|        +----------------------------+
|        | extra_headers    19 : x-19 |
+=====================================+
| event  | fixed part        x : y    |
| data   +----------------------------+
|        | variable part              |
+=====================================+
2.2 format_desc event
下面是我在FLUSH LOGS之后新建的一个全新的binlog文件mysql-bin.000053，从binlog第一个event也就是format_desc event开始分析（mysql日志是小端字节序）：

root@ubuntu:/var/log/mysql# hexdump -C mysql-bin.000053

00000000  fe 62 69 6e b8 b2 7f 56  0f 04 00 00 00 67 00 00  |.bin...V.....g..|
00000010  00 6b 00 00 00 01 00 04  00 35 2e 35 2e 34 36 2d  |.k.......5.5.46-|
00000020  30 75 62 75 6e 74 75 30  2e 31 34 2e 30 34 2e 32  |0ubuntu0.14.04.2|
00000030  2d 6c 6f 67 00 00 00 00  00 00 00 00 00 00 00 00  |-log............|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 13  |................|
00000050  38 0d 00 08 00 12 00 04  04 04 04 12 00 00 54 00  |8.............T.|
00000060  04 1a 08 00 00 00 08 08  08 02 00                 |...........|

对照官方文档中的说明来看下format_desc event格式：

+=====================================+
| event  | timestamp         0 : 4    |
| header +----------------------------+
|        | type_code         4 : 1    | = FORMAT_DESCRIPTION_EVENT = 15
|        +----------------------------+
|        | server_id         5 : 4    |
|        +----------------------------+
|        | event_length      9 : 4    | >= 91
|        +----------------------------+
|        | next_position    13 : 4    |
|        +----------------------------+
|        | flags            17 : 2    |
+=====================================+
| event  | binlog_version   19 : 2    | = 4
| data   +----------------------------+
|        | server_version   21 : 50   |
|        +----------------------------+
|        | create_timestamp 71 : 4    |
|        +----------------------------+
|        | header_length    75 : 1    |
|        +----------------------------+
|        | post-header      76 : n    | = array of n bytes, one byte per event
|        | lengths for all            |   type that the server knows about
|        | event types                |
+=====================================+
前面4个字节是固定的magic number,值为0x6e6962fe。接着是一个format_desc event，先看下19个字节的header。这19个字节中前4个字节0x567fb2b8是时间戳，第5个字节0x0f是event type，接着4个字节0x00000004是server_id，再接着4个字节0x00000067是长度103，然后的4个字节0x0000006b是下一个event的起始位置107，接着的2个字节的0x0001是flag（1为LOG_EVENT_BINLOG_IN_USE_F，标识binlog还没有关闭，binlog关闭后，flag会被设置为0），这样4+1+4+4+4+2=19个字节的公共头就完了(extra_headers暂时没有用到)。然后是这个event的data部分，event的data分为Fixed data和Variable data两部分，其中Fixed data是event的固定长度和格式的数据，Variable data则是长度变化的数据，比如format_desc event的Fixed data长度是0x54=84个字节。下面看下这84=2+50+4+1+27个字节的分配：开始的2个字节0x0004为binlog的版本号4，接着的50个字节为mysql-server版本，如我的版本是5.5.46-0ubuntu0.14.04.2-log，与SELECT version();查看的结果一致。接下来4个字节是binlog创建时间，这里是0；然后的1个字节0x13是指之后所有event的公共头长度，这里都是19；接着的27个字节中每个字节为mysql已知的event（共27个）的Fixed data的长度；可以发现format_desc event自身的Variable data部分为空。

2.3 rotate event
接着我们不做额外操作，直接FLUSH LOGS，可以看到一个rotate event，此时的binlog内容如下:

......
00000060  ................................. c2 b3 7f 56 04  |..............V.|
00000070  04 00 00 00 2b 00 00 00  96 00 00 00 00 00 04 00  |....+...........|
00000080  00 00 00 00 00 00 6d 79  73 71 6c 2d 62 69 6e 2e  |......mysql-bin.|
00000090  30 30 30 30 35 34                                 |000054|
00000096
前面的内容跟之前的几乎一致，除了format_desc event的flag从0x0001变成了0x0000。然后从0x567fb3c2开始是一个rotate event。依照前面的分析，前面19个字节为event的header，其event type是0x04，长度为0x2b=43，下一个event起始位置为0x96=150，然后是flag为0x0000，接着是event data部分，首先的8个字节为Fixed data部分，记录的是下一个binlog的位置偏移4，而余下来的43-19-8=16个字节为Variable data部分，记录的是下一个binlog的文件名mysql-bin.000054。对照mysqlbinlog -vv mysql-bin.000053可以验证。

ssj@ubuntu:/var/log/mysql$ mysqlbinlog -vv mysql-bin.000053 
...
# at 4
#151227 17:43:20 server id 4  end_log_pos 107   Start: binlog v 4, server v 5.5.46-0ubuntu0.14.04.2-log created 151227 17:43:20
BINLOG '
uLJ/Vg8EAAAAZwAAAGsAAAAAAAQANS41LjQ2LTB1YnVudHUwLjE0LjA0LjItbG9nAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAEzgNAAgAEgAEBAQEEgAAVAAEGggAAAAICAgCAA==
'/*!*/;
# at 107
#151227 17:47:46 server id 4  end_log_pos 150   Rotate to mysql-bin.000054  pos: 4
...
2.4 query event
刷新binlog，设置binlog_format=statement，创建一个表CREATE TABLEtt(ivarchar(100) DEFAULT NULL) ENGINE=InnoDB, 然后在测试表tt中插入一条数据insert into tt values('abc')，会产生3个event，包括2个query event和1个xid event。其中2个query event分别是BEGIN以及INSERT 语句，而xid event则是事务提交语句（xid event是支持XA的存储引擎才有的，因为测试表tt是innodb引擎的，所以会有。如果是myisam引擎的表，也会有BEGIN和COMMIT,只不过COMMIT会是一个query event而不是xid event）。

mysql> show binlog events in 'mysql-bin.000060';
+------------------+-----+-------------+-----------+-------------+--------------------------------------------------------+
| Log_name         | Pos | Event_type  | Server_id | End_log_pos | Info                                                   |
+------------------+-----+-------------+-----------+-------------+--------------------------------------------------------+
| mysql-bin.000060 |   4 | Format_desc |         4 |         107 | Server ver: 5.5.46-0ubuntu0.14.04.2-log, Binlog ver: 4 |
| mysql-bin.000060 | 107 | Query       |         4 |         175 | BEGIN                                                  |
| mysql-bin.000060 | 175 | Query       |         4 |         266 | use `test`; insert into tt values('abc')               |
| mysql-bin.000060 | 266 | Xid         |         4 |         293 | COMMIT /* xid=138 */                                   |
+------------------+-----+-------------+-----------+-------------+--------------------------------------------------
binlog如下：

.......
0000006b  01 9d 82 56 02 04 00 00  00 44 00 00 00 af 00 00  |...V.....D......|
0000007b  00 08 00 26 00 00 00 00  00 00 00 04 00 00 1a 00  |...&............|
0000008b  00 00 00 00 00 01 00 00  00 00 00 00 00 00 06 03  |................|
0000009b  73 74 64 04 21 00 21 00  08 00 74 65 73 74 00 42  |std.!.!...test.B|
000000ab  45 47 49 4e                                       |EGIN|

000000af  01 9d 82 56 02 04 00 00  00 5b 00 00 00 0a 01 00  |...V.....[......|
000000bf  00 00 00 26 00 00 00 00  00 00 00 04 00 00 1a 00  |...&............|
000000cf  00 00 00 00 00 01 00 00  00 00 00 00 00 00 06 03  |................|
000000df  73 74 64 04 21 00 21 00  08 00 74 65 73 74 00 69  |std.!.!...test.i|
000000ef  6e 73 65 72 74 20 69 6e  74 6f 20 74 74 20 76 61  |nsert into tt va|
000000ff  6c 75 65 73 28 27 61 62  63 27 29                 |lues('abc')|

0000010a  01 9d 82 56 10 04 00 00  00 1b 00 00 00 25 01 00  |...V.........%..|
0000011a  00 00 00 8a 00 00 00 00  00 00 00                 |...........|
抛开format_desc event，从0000006b开始分析第一个query event。头部跟之前的event一样，只是query event的type为0x02，长度为0x44=64，下一个event位置为0xaf=175。flag为8，接着是data部分，从format_desc event我们可以知道query event的Fixed data部分为13个字节，因此也可以算出Variable data部分为64-19-13=32字节。

Fixed data：首先的4个字节0x00000026为执行该语句的thread id，接下来的4个字节是执行的时间0(以秒为单位)，接下来的1个字节0x04是语句执行时的默认数据库名字的长度，我这里数据库是test，所以长度为4.接着的2个字节0x0000是错误码（注：通常情况下错误码是0表示没有错误，但是在一些非事务性表如myisam表执行INSERT...SELECT语句时可能插入部分数据后遇到duplicate-key错误会产生错误码1062，或者是事务性表在INSERT...SELECT出错不会插入部分数据，但是在执行过程中CTRL+C终止语句也可能记录错误码。slave db在复制时会执行后检查错误码是否一致，如果不一致，则复制过程会中止）,接着2个字节0x001a为状态变量块的长度26。
Variable data：从0x001a之后的26个字节为状态变量块（这个暂时先不管），然后是默认数据库名test，以0x00结尾，然后是sql语句BEGIN，接下来就是第2个query event的内容了。
第二个query event与第一个格式一样，只是执行语句变成了insert into tt values('abc')。

第三个xid event为COMMIT语句。前19个字节是通用头部，type是16。data部分中Fixed data为空，而variable data为8个字节，这8个字节0x000000008a是事务编号（注意事务编号不一定是小端字节序，因为是从内存中拷贝到磁盘的，所以这个字节序跟机器相关）。

2.5 table_map event & write_rows event
这两个event是在binlog_format=row的时候使用，设置binlog_format=row，然后创建一个测试表

CREATE TABLE `trow` (
  `i` int(11) NOT NULL,
  `c` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`i`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1`

执行语句INSERT INTO trow VALUES(1, NULL), (2, 'a')，这个语句会产生一个query event，一个table_map event、一个write_rows event以及一个xid event。


mysql> show binlog events in 'mysql-bin.000074';
| Log_name         | Pos | Event_type  | Server_id | End_log_pos | Info                                                   |
+------------------+-----+-------------+-----------+-------------+--------------------------------------------------------+
| mysql-bin.000074 |   4 | Format_desc |         4 |         107 | Server ver: 5.5.46-0ubuntu0.14.04.2-log, Binlog ver: 4 |
| mysql-bin.000074 | 107 | Query       |         4 |         175 | BEGIN                                                  |
| mysql-bin.000074 | 175 | Table_map   |         4 |         221 | table_id: 50 (test.trow)                               |
| mysql-bin.000074 | 221 | Write_rows  |         4 |         262 | table_id: 50 flags: STMT_END_F                         |
| mysql-bin.000074 | 262 | Xid         |         4 |         289 | COMMIT /* xid=245 */                                   
对应的mysql-bin.000074数据如下：

...
#query event (BEGIN)
0000006b  95 2a 85 56 02 04 00 00  00 44 00 00 00 af 00 00  |.*.V.....D......|
0000007b  00 08 00 26 00 00 00 00  00 00 00 04 00 00 1a 00  |...&............|
0000008b  00 00 00 00 00 01 00 00  00 00 00 00 00 00 06 03  |................|
0000009b  73 74 64 04 21 00 21 00  08 00 74 65 73 74 00 42  |std.!.!...test.B|
000000ab  45 47 49 4e                                       |EGIN|

#table_map event
000000af  95 2a 85 56 13 04 00 00  00 2e 00 00 00 dd 00 00  |.*.V............|
000000bf  00 00 00 32 00 00 00 00  00 01 00 04 74 65 73 74  |...2........test|
000000cf  00 04 74 72 6f 77 00 02  03 0f 02 0a 00 02        |..trow........|

#write_rows event
000000dd  95 2a 85 56 17 04 00 00  00 29 00 00 00 06 01 00  |.*.V.....)......|
000000ed  00 00 00 32 00 00 00 00  00 01 00 02 ff fe 01 00  |...2............|
000000fd  00 00 fc 02 00 00 00 01  61                       |........a|

#xid event
00000106  95 2a 85 56 10 04 00 00  00 1b 00 00 00 21 01 00  |.*.V.........!..|
00000116  00 00 00 f5 00 00 00 00  00 00 00                 |...........|
0x0000006b-0x000000ae为query event，语句是BEGIN,前面已经分析过。

table_map event
0x0000000af开始为table_map event。除去头部19个字节，Fixed data为8个字节，前面6个字节0x32=50为table id，接着2个字节0x0001为flags。

Variable data部分，首先1个字节0x04为数据库名test的长度，然后5个字节是数据库名test+结束符。接着1个字节0x04为表名长度，接着5个字节为表名trow+结束符。接着1个字节0x02为列的数目。而后是2个列的类型定义，分别是0x03和0x0f（列的类型MYSQL_TYPE_LONG为0x03，MYSQL_TYPE_VARCHAR为0x0f)。接着是列的元数据定义，首先0x02表示元数据长度为2，因为MYSQL_TYPE_LONG没有元数据，而MYSQL_TYPE_VARCHAR元数据长度为2。接着的0x000a就是MYSQL_TYPE_VARCHAR的元数据，表示我们在定义表时的varchar字段c长度为10，最后一个字节0x02为掩码，表示第一个字段i不能为NULL。关于列的类型以及元数据等更详细的信息可以参见http://dev.mysql.com/doc/internals/en/table-map-event.html。

write_rows event
从0x000000dd开始为write_rows event，除去头部19个字节，前6个字节0x32也是table id，然后两个字节0x0001为flags。接着的1个字节0x02为表中列的数目。然后1个字节0xff各个bit标识各列是否存在值，这里表示都存在。

接着的就是插入的各行数据了。第1个字节0xfe的各个bit标识该行变化之后各列是否为NULL，为NULL记为1.这里表示第1列不为NULL，因为第一行数据插入的是(1,NULL)。接下来是各列的数据，第一列是MYSQL_TYPE_LONG,长度为4个字节，所以0x00000001就是这个值。第二列是NULL不占字节。接下来是第二行，先是0xfc标识两列都不为NULL，先读取第一列的4个字节0x00000002也就是我们插入的数字2，然后读取第二列，先是一个字节的长度0x01，然后是内容0x61也就是字符'a'。到此，write_rows event也就分析完了。rows相关的event还有update_rows event和delete_rows event等，欲了解更多可以参见官方文档。

最后是xid event，之前已经分析过，不再赘述。

2.6 intvar event
intvar event在binlog_format=statement时使用到，用于自增键类型auto_increment，十分重要。intval event的Fixed data部分为空，而Variable data部分为9个字节，第1个字节用于标识自增事件类型 LAST_INSERT_ID_EVENT = 1 or INSERT_ID_EVENT = 2，余下的8个字节为自增ID。
创建一个测试表 create table tinc (i int auto_increment primary key, c varchar(10)) engine=innodb;，然后执行一个插入语句INSERT INTO tinc(c) values('abc');就可以看到intvar event了，这里的自增事件类型为INSERT_ID_EVENT。而如果用语句INSERT INTO tinc(i, c) VALUES(LAST_INSERT_ID()+1, 'abc')，则可以看到自增事件类型为LAST_INSERT_ID_EVENT的intvar event。

| Log_name         | Pos | Event_type  | Server_id | End_log_pos | Info                                                   |
+------------------+-----+-------------+-----------+-------------
| mysql-bin.000079 |   4 | Format_desc |         4 |         107 | Server ver: 5.5.46-0ubuntu0.14.04.2-log, Binlog ver: 4 |
| mysql-bin.000079 | 107 | Query       |         4 |         175 | BEGIN                                                  |
| mysql-bin.000079 | 175 | Intvar      |         4 |         203 | INSERT_ID=1                                            |
| mysql-bin.000079 | 203 | Query       |         4 |         299 | use `test`; insert into tinc(c) values('abc')          |
| mysql-bin.000079 | 299 | Xid         |         4 |         326 | COMMIT /* xid=263 */                                
3 简单总结
上面提到，binlog有三种格式，各有优缺点：

statement：基于SQL语句的模式，某些语句和函数如UUID, LOAD DATA INFILE等在复制过程可能导致数据不一致甚至出错。
row：基于行的模式，记录的是行的变化，很安全。但是binlog会比其他两种模式大很多，在一些大表中清除大量数据时在binlog中会生成很多条语句，可能导致从库延迟变大。
mixed：混合模式，根据语句来选用是statement还是row模式。
不同版本的mysql在主从复制要慎重，虽然mysql5.0之后都用的V4版本的binlog了，估计还是会有些坑在里面，特别是高版本为主库，低版本为从库时容易出问题。在主从复制时最好还是主库从库版本一致，至少是大版本一致
<!-- more -->
<img src="{{site.url}}{{site.baseurl}}/img/mysqlbinlogformat.png"/>
二进制有两个最重要的使用场景: 
    其一：MySQL Replication在Master端开启binlog，Mster把它的二进制日志传递给slaves来达到master-slave数据一致的目的。 
    其二：自然就是数据恢复了，通过使用mysqlbinlog工具来使恢复数据。 
    二进制日志包括两类文件：二进制日志索引文件（文件名后缀为.index）用于记录所有的二进制文件，二进制日志文件（文件名后缀为.00000*）记录数据库所有的DDL和DML(除了数据查询语句)语句事件
常用binlog日志操作命令
    1.查看所有binlog日志列表
      mysql> show master logs;

    2.查看master状态，即最后(最新)一个binlog日志的编号名称，及其最后一个操作事件pos结束点(Position)值
      mysql> show master status;

    3.刷新log日志，自此刻开始产生一个新编号的binlog日志文件
      mysql> flush logs;
      注：每当mysqld服务重启时，会自动执行此命令，刷新binlog日志；在mysqldump备份数据时加 -F 选项也会刷新binlog日志；

    4.重置(清空)所有binlog日志
      mysql> reset master;

binlog日志与数据库文件在同目录中(我的环境配置安装是选择在/usr/local/mysql/data中)
      在MySQL5.5以下版本使用mysqlbinlog命令时如果报错，就加上 “--no-defaults”选项
    
      # /usr/local/mysql/bin/mysqlbinlog /usr/local/mysql/data/mysql-bin.000013
        下面截取一个片段分析：

         ...............................................................................
         # at 552
         #131128 17:50:46 server id 1  end_log_pos 665   Query   thread_id=11    exec_time=0     error_code=0 ---->执行时间:17:50:46；pos点:665
         SET TIMESTAMP=1385632246/*!*/;
         update zyyshop.stu set name='李四' where id=4              ---->执行的SQL
         /*!*/;
         # at 665
         #131128 17:50:46 server id 1  end_log_pos 692   Xid = 1454 ---->执行时间:17:50:46；pos点:692 
         ...............................................................................

         注: server id 1     数据库主机的服务号；
             end_log_pos 665 pos点
             thread_id=11    线程号


    2.上面这种办法读取出binlog日志的全文内容较多，不容易分辨查看pos点信息，这里介绍一种更为方便的查询命令：

      mysql> show binlog events [IN 'log_name'] [FROM pos] [LIMIT [offset,] row_count];

             选项解析：
               IN 'log_name'   指定要查询的binlog文件名(不指定就是第一个binlog文件)
               FROM pos        指定从哪个pos起始点开始查起(不指定就是从整个文件首个pos点开始算)
               LIMIT [offset,] 偏移量(不指定就是0)
               row_count       查询总条数(不指定就是所有行)

             截取部分查询结果：
             *************************** 20. row ***************************
                Log_name: mysql-bin.000021  ----------------------------------------------> 查询的binlog日志文件名
                     Pos: 11197 ----------------------------------------------------------> pos起始点:
              Event_type: Query ----------------------------------------------------------> 事件类型：Query
               Server_id: 1 --------------------------------------------------------------> 标识是由哪台服务器执行的
             End_log_pos: 11308 ----------------------------------------------------------> pos结束点:11308(即：下行的pos起始点)
                    Info: use `zyyshop`; INSERT INTO `team2` VALUES (0,345,'asdf8er5') ---> 执行的sql语句
             *************************** 21. row ***************************
                Log_name: mysql-bin.000021
                     Pos: 11308 ----------------------------------------------------------> pos起始点:11308(即：上行的pos结束点)
              Event_type: Query
               Server_id: 1
             End_log_pos: 11417
                    Info: use `zyyshop`; /*!40000 ALTER TABLE `team2` ENABLE KEYS */
             *************************** 22. row ***************************
                Log_name: mysql-bin.000021
                     Pos: 11417
              Event_type: Query
               Server_id: 1
             End_log_pos: 11510
                    Info: use `zyyshop`; DROP TABLE IF EXISTS `type`

      这条语句可以将指定的binlog日志文件，分成有效事件行的方式返回，并可使用limit指定pos点的起始偏移，查询条数；
      
      A.查询第一个(最早)的binlog日志：
        mysql> show binlog events\G; 
    
      B.指定查询 mysql-bin.000021 这个文件：
        mysql> show binlog events in 'mysql-bin.000021'\G;

      C.指定查询 mysql-bin.000021 这个文件，从pos点:8224开始查起：
        mysql> show binlog events in 'mysql-bin.000021' from 8224\G;

      D.指定查询 mysql-bin.000021 这个文件，从pos点:8224开始查起，查询10条
        mysql> show binlog events in 'mysql-bin.000021' from 8224 limit 10\G;

      E.指定查询 mysql-bin.000021 这个文件，从pos点:8224开始查起，偏移2行，查询10条
        mysql> show binlog events in 'mysql-bin.000021' from 8224 limit 2,10\G;


五、恢复binlog日志实验(zyyshop是数据库)
    1.假设现在是凌晨4:00，我的计划任务开始执行一次完整的数据库备份：

      将zyyshop数据库备份到 /root/BAK.zyyshop.sql 文件中：
      # /usr/local/mysql/bin/mysqldump -uroot -p123456 -lF --log-error=/root/myDump.err -B zyyshop > /root/BAK.zyyshop.sql
从binlog日志恢复数据
      
      恢复语法格式：
      # mysqlbinlog mysql-bin.0000xx | mysql -u用户名 -p密码 数据库名

        常用选项：
          --start-position=953                   起始pos点
          --stop-position=1437                   结束pos点
          --start-datetime="2013-11-29 13:18:54" 起始时间点
          --stop-datetime="2013-11-29 13:21:53"  结束时间点
          --database=zyyshop                     指定只恢复zyyshop数据库(一台主机上往往有多个数据库，只限本地log日志)
            
        不常用选项：    
          -u --user=name              Connect to the remote server as username.连接到远程主机的用户名
          -p --password[=name]        Password to connect to remote server.连接到远程主机的密码
          -h --host=name              Get the binlog from server.从远程主机上获取binlog日志
          --read-from-remote-server   Read binary logs from a MySQL server.从某个MySQL服务器上读取binlog日志

      小结：实际是将读出的binlog日志内容，通过管道符传递给mysql命令。这些命令、文件尽量写成绝对路径；

      A.完全恢复(本例不靠谱，因为最后那条 drop database zyyshop 也在日志里，必须想办法把这条破坏语句排除掉，做部分恢复)
        # /usr/local/mysql/bin/mysqlbinlog  /usr/local/mysql/data/mysql-bin.000021 | /usr/local/mysql/bin/mysql -uroot -p123456 -v zyyshop 

      B.指定pos结束点恢复(部分恢复)：
        @ --stop-position=953 pos结束点
        注：此pos结束点介于“导入实验数据”与更新“name='李四'”之间，这样可以恢复到更改“name='李四'”之前的“导入测试数据”
        # /usr/local/mysql/bin/mysqlbinlog --stop-position=953 --database=zyyshop /usr/local/mysql/data/mysql-bin.000023 | /usr/local/mysql/bin/mysql -uroot -p123456 -v zyyshop