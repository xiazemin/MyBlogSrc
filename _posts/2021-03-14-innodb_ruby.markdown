---
title: innodb_ruby
layout: post
category: storage
author: 夏泽民
---
https://www.cnblogs.com/cnzeno/p/6322842.html

CREATE TABLE `t2` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `remark` varchar(50) DEFAULT NULL,
  `add_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `ix_t2_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3829 DEFAULT CHARSET=utf8

source  /Users/xiazemin/Downloads/mysql3376_backup_test_inndb_ruby.sql

https://github.com/jeremycole/innodb_ruby/wiki

 show tables;
+--------------------+
| Tables_in_zeno3376 |
+--------------------+
| t1                 |
| t2                 |
+--------------------+

 show variables like "datadir";
+---------------+-----------------------+
| Variable_name | Value                 |
+---------------+-----------------------+
| datadir       | /usr/local/var/mysql/ |
+---------------+-----------------------+

% ls /usr/local/var/mysql/zeno3376
t1.ibd	t2.ibd

https://blog.jcole.us/2013/01/10/btree-index-structures-in-innodb/

https://blog.jcole.us/2013/01/03/a-quick-introduction-to-innodb-ruby/

<!-- more -->
https://www.cnblogs.com/cnzeno/p/6322842.html

https://blog.csdn.net/Chuck_Perry/article/details/64441289?utm_source=blogxgwz39

https://www.jianshu.com/p/ad553222cfbc


git clone 

cd innodb_ruby/
 
 % sudo gem install innodb_ruby

 innodb_space -s  /usr/local/var/mysql/ibdata1 -T /usr/local/var/mysql/zeno3376/t2  space-page-type-regions
 
 2.6.0/gems/innodb_ruby-0.9.16/lib/innodb/index.rb:33:in `page': undefined method `record_describer=' for #<Innodb::Page::FspHdrXdes:0x00007fd6a610dae0> (NoMethodError)
 
 https://github.com/jeremycole/innodb_ruby/wiki
 
 show variables like "version";
+---------------+--------+
| Variable_name | Value  |
+---------------+--------+
| version       | 8.0.22 |
+---------------+--------+
1 row in set (0.00 sec)


https://web.stanford.edu/class/cs245/homeworks/b+tree/readings/MySQL_InnoDB.B+Tree.Reading3.pdf

https://blog.csdn.net/vkingnew/article/details/82775944

 % brew search mysql
==> Formulae
automysqlbackup           mysql-client              mysql-sandbox             mysql@5.7
mysql ✔                   mysql-client@5.7          mysql-search-replace      mysqltuner
mysql++                   mysql-connector-c++       mysql@5.6
==> Casks
mysql-connector-python             mysql-utilities                    navicat-for-mysql
mysql-shell                        mysqlworkbench                     sqlpro-for-mysql


% brew install  mysql@5.7

% mysql.server restart

https://www.yuque.com/alipayrrql7pkhl7/mxxd2g/omzshe


/usr/local/Cellar/mysql@5.7/5.7.32/bin/mysql.server restart
 ERROR! MySQL server PID file could not be found!
Starting MySQL
. ERROR! The server quit without updating PID file (/usr/local/var/mysql/bogon.pid).

 cat  /usr/local/var/mysql/bogon.err
 
2021-03-14T11:49:50.962562Z 0 [ERROR] [FATAL] InnoDB: Table flags are 0 in the data dictionary but the flags in file ./ibdata1 are 0x4800!
2021-03-14 19:49:50 0x105495dc0  InnoDB: Assertion failure in thread 4383661504 in file ut0ut.cc line 921


cp /usr/local/etc/my.cnf /usr/local/etc/my5.7.cnf
datadir=/usr/local/var/mysql5.7

/usr/local/Cellar/mysql@5.7/5.7.32/bin/mysqld --defaults-file=/usr/local/etc/my5.7.cnf --user=root

 [ERROR] Can't open the mysql.plugin table. Please run mysql_upgrade to create it.
2021-03-14T12:06:27.249095Z 0 [ERROR] unknown variable 'mysqlx-bind-address=127.0.0.1'
2021-03-14T12:06:27.249113Z 0 [ERROR] Aborting


#mysqlx-bind-address = 127.0.0.1


/usr/local/Cellar/mysql@5.7/5.7.32/bin/mysqld  --initalize --defaults-file=/usr/local/etc/my5.7.cnf  --user=root --console

 [ERROR] [FATAL] InnoDB: Table flags are 0 in the data dictionary but the flags in file ./ibdata1 are 0x4800!
 
 
 /usr/local/Cellar/mysql@5.7/5.7.32/support-files/mysql.server restart
 ERROR! MySQL server PID file could not be found!
Starting MySQL
. ERROR! The server quit without updating PID file (/usr/local/var/mysql/bogon.pid).


 /usr/local/Cellar/mysql@5.7/5.7.32/bin/mysqld --defaults-file=/usr/local/etc/my5.7.cnf  --user=root --datadir=/usr/local/var/mysql5.7/ --basedir=/usr/local/Cellar/mysql@5.7/5.7.32
 
 http://www.360doc.com/content/18/1003/08/835902_791541046.shtml
 
  /usr/local/Cellar/mysql@5.7/5.7.32/bin/mysqld --defaults-file=/usr/local/etc/my5.7.cnf  --user=root --datadir=/usr/local/var/mysql5.7/ --basedir=/usr/local/Cellar/mysql@5.7/5.7.32 --initialize
 
 0 [ERROR] --initialize specified but the data directory has files in it. Aborting.
 
 rm -rf /usr/local/var/mysql5.7/*
 
 A temporary password is generated for root@localhost: i8dTC=?kngxW
 

 /usr/local/Cellar/mysql@5.7/5.7.32/bin/mysqld --defaults-file=/usr/local/etc/my5.7.cnf  --user=root --datadir=/usr/local/var/mysql5.7/ --basedir=/usr/local/Cellar/mysql@5.7/5.7.32
 
 
 mysql -uroot -pi8dTC=\?kngxW
 
 mysql>  alter user'root'@'localhost' identified by '123456';
 
  mysql -uroot -p123456
  
  https://blog.csdn.net/Gordo_Li/article/details/103511374
  
  source  /Users/xiazemin/Downloads/mysql3376_backup_test_inndb_ruby.sql
  
   innodb_space -s  /usr/local/var/mysql5.7/ibdata1 -T /usr/local/var/mysql5.7/zeno3376/t2  space-page-type-regions
   
   /Library/Ruby/Gems/2.6.0/gems/innodb_ruby-0.9.16/lib/innodb/system.rb:101:in `space_by_table_name': Table /usr/local/var/mysql5.7/zeno3376/t2 not found (RuntimeError)
 
cd  /usr/local/var/mysql5.7/
 % innodb_space -s ibdata1 -T zeno3376/t2 space-indexes
id          name                            root        fseg        fseg_id     used        allocated   fill_factor
42          PRIMARY                         3           internal    1           1           1           100.00%
42          PRIMARY                         3           leaf        2           9           9           100.00%
43          ix_t2_name                      4           internal    3           1           1           100.00%
43          ix_t2_name                      4           leaf        4           4           4           100.00%


 %  innodb_space -s ibdata1 -T zeno3376/t2  space-page-type-regions
start       end         count       type
0           0           1           FSP_HDR
1           1           1           IBUF_BITMAP
2           2           1           INODE
3           17          15          INDEX
18          18          1           FREE (ALLOCATED)

　　　　start：从第几个page开始

　　　　end：从第几个page结束

　　　　count：占用了多少个page；

　　　　type: page的类型

　　　　从上面的结果可以看出：“FSP_HDR”、“IBUF_BITMAP”、“INODE”是分别占用了0,1,2号的page，从3号page开始才是存放数据和索引的页（Index），占用了3~17号的page，共15个page。

　　　　接下来，根据得到的聚集索引和辅助索引的根节点来获取索引上的其他page的信息
　　　　
innodb_space -s ibdata1 -T zeno3376/t2 -I primary -p 3 page-records      

system.rb:213:in `index_by_name': undefined method `[]' for nil:NilClass (NoMethodError)

%  innodb_space -s ibdata1 -T zeno3376/t2 -p 3 page-records
Record 126: (id=1782) → #5
Record 140: (id=1890) → #6
Record 154: (id=2101) → #7
Record 168: (id=2317) → #10
Record 182: (id=2531) → #11
Record 196: (id=2747) → #12
Record 210: (id=2964) → #15
Record 224: (id=3179) → #16
Record 238: (id=3394) → #17

　　　　上面的结果是解析聚集索引根节点页的信息，1行就代表使用了1个page，所以，叶子节点共使用了9个page，根节点使用了1个page，跟space_indexes的解析结果一致。

　　　　Record 126: (id=1782) → #5 

　　　   id = 1782 代表的就是表中id为1782的记录，因为id是主键

　　　　-> #5 代表的是指向5号page

　　　　Record 126: (id=1782) → #5： 整行的意思就是5号page的id最小值是1782，包含了1782~1889的行记录。

　　　　注意：page number并不是连续的

　　　　根据解析root得到的信息，继续解析第一个叶子节点的信息

# innodb_space -s ibdata1 -T zeno3376/t2 -p 5 page-records
Record 128: (id=1782) → (name="zeno", remark="mysql", add_time=:NULL)

Record 162: (id=1783) → (name="KIK91QJET1FCZ46EJKML", remark="H4HJO5F7W5GSSDORT8AAT", add_time="184524556-49-63 92:14:08")

.....


从上面可以看出，聚集索引的叶子节点是包含了行记录的所有数据。

　　　　同理，解析辅助索引ix_t2_name，但是需要注意的是，在解析辅助索引是，需要加上“-I ix_t2_name”

 % innodb_space -s ibdata1 -T zeno3376/t2 -I ix_t2_name -p 4 page-records
Record 127: (name="01EE2CCYUW35K0LVT5DAG2044NW") → #8
Record 196: (name="8WCS36CV56KGA8NE6OG23QFS") → #13
Record 169: (name="HQVX6ZX7H2XI") → #9
Record 235: (name="QXS8RUJF6FY") → #14


从上面可以出，辅助索引ix_t2_name的key是name列，叶子节点共使用了4个page，加上根节点，那么辅助索引ix_t2_name共使用了5个page，跟使用space_indexes解析出来的结果一致。

　　　　Record 127: (name="01EE2CCYUW35K0LVT5DAG2044NW") → #8 这条记录代表的意思是辅助索引的第1个叶子节点的page number是8，8号page的第一个key值是"01EE2CCYUW35K0LVT5DAG2044NW"

　　　　Record 196: (name="8WCS36CV56KGA8NE6OG23QFS") → #13 这条记录代表的意思是辅助索引的第2个叶子节点的page number是13，13号page的第一个key值是"8WCS36CV56KGA8NE6OG23QFS"

　　　　其它的记录如此类推……

　　　　接下来看看辅助索引的叶子节点的结构

innodb_space -s ibdata1 -T zeno3376/t2 -I ix_t2_name -p 8 page-records
Record 127: (name="01EE2CCYUW35K0LVT5DAG2044NW") → (id=1855)

Record 165: (name="02RFY8SJLQ879F2CYHI") → (id=2132)
.....

从上面可以看到叶子节点中包含可辅助索引和主键列

　　　　Record 127: (name="01EE2CCYUW35K0LVT5DAG2044NW") → (id=1855) 代表的意思就是name值为"01EE2CCYUW35K0LVT5DAG2044NW"的记录指向主键id=1855的行记录。


https://www.cnblogs.com/cnzeno/p/6322842.html

 show create table t2;
 
 CREATE TABLE `t2` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `remark` varchar(50) DEFAULT NULL,
  `add_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `ix_t2_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3829 DEFAULT CHARSET=utf8 

% innodb_space -s ibdata1 -T zeno3376/t2  data-dictionary-tables
name                            id          n_cols      type        mix_id      mix_len     cluster_name   space
SYS_DATAFILES                   14          2           1           0           64                         0
SYS_FOREIGN                     11          4           1           0           64                         0
SYS_FOREIGN_COLS                12          4           1           0           64                         0
SYS_TABLESPACES                 13          3           1           0           64                         0
SYS_VIRTUAL                     15          3           1           0           64                         0
mysql/engine_cost               34          2147483654  33          0           80                         20
mysql/gtid_executed             32          2147483651  33          0           80                         18

 % innodb_space -s ibdata1 -T zeno3376/t2    space-index-pages-free-plot
 
 kernel_require.rb:54:in `require': cannot load such file -- gnuplot (LoadError)
 https://stackoverflow.com/questions/20999055/plotting-pdf-in-gnuplot-error-cannot-open-load-file-stat-inc
 
 % innodb_space -s ibdata1 -T zeno3376/t2 space-page-type-summary
type                count       percent     description
INDEX               15          78.95       B+Tree index
FSP_HDR             1           5.26        File space header
IBUF_BITMAP         1           5.26        Insert buffer bitmap
INODE               1           5.26        File segment inode
ALLOCATED           1           5.26        Freshly allocated


 % innodb_space -s ibdata1 -T zeno3376/t2   page-illustrate -p6
 
 Legend (█ = 1 byte):
  Region Type                         Bytes    Ratio
  █ FIL Header                           38    0.23%
  █ Index Header                         36    0.22%
  █ File Segment Header                  20    0.12%
  █ Infimum                              13    0.08%
  █ Supremum                             13    0.08%
  █ Record Header                      1688   10.30%
  █ Record Data                       14038   85.68%
  █ Page Directory                      106    0.65%
  █ FIL Trailer                           8    0.05%
  ░ Garbage                               0    0.00%
    Free                                424    2.59%
    
  % innodb_space -s ibdata1 -T zeno3376/t2  -p 6 -R 15143 record-dump
Record at offset 15143

Header:
  Next record offset  : 112
  Heap number         : 212
  Type                : conventional
  Deleted             : false
  Length              : 8

System fields:
  Transaction ID: 1296
  Roll Pointer:
    Undo Log: page 290, offset 4278
    Rollback Segment ID: 46
    Insert: true

Key fields:
  id: 2100

Non-key fields:
  name: "GD4NUSRK8ZAX63O1UMW4ZBYGOM8"
  remark: "QPYW2MCBIAT69Y9OHEC7N"
  add_time: "184524556-56-84 67:86:56"
  
  % innodb_space -s ibdata1 -T zeno3376/t2  -p 5 -R 7444 record-dump
Record at offset 7444

Header:
  Next record offset  : 112
  Heap number         : 109
  Type                : conventional
  Deleted             : false
  Length              : 8

System fields:
  Transaction ID: 1296
  Roll Pointer:
    Undo Log: page 290, offset 1556
    Rollback Segment ID: 46
    Insert: true

Key fields:
  id: 1889

Non-key fields:
  name: "0O2S6OCUC99MQKM1"
  remark: "1K5GJEQ5QU83T3F"
  add_time: "184524556-52-49 32:71:04"
  
 https://www.jianshu.com/p/ad553222cfbc 
 
 
 
