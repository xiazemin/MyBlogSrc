---
title: mysql instant add colum
layout: post
category: storage
author: 夏泽民
---
“Instant ADD COLUMN”，即“瞬加字段功能”
鹅厂工程师通过扩展MySQL InnoDB的存储格式，可以把原来几个小时才能完成的给表加字段命令，在1秒之内执行完成，更新TB级的表都是毛毛雨，有效地提高了数据库的管理效率，降低运维成本。
随着MySQL新版本的发布，陈福荣和梁飞龙将该特性提交到MySQL 8.0.12

DDL的改进-快速增加新列
     实现原理
ONline DDL 除了inplace和copy之外新加了 instant：过程不加表锁，能极大的提高执行能力和并发能力。

 instant：并没有对表做大的操作，只是将变更记录到数据字典里。只是分为了两种格式：new，old。

                 1.add column 操作之前的数据依旧保存在老格式里；之后的数据保存的新格式里。

                 2. 第一次做ddl操作，会 instant flag 里记录下 “字段数量”，和新增的column的字段值。

    8.0.11 add column不支持的特性
不支持将列加到中间，只支持将列加到最后。

不支持compress格式的表

如果在add column之前有fulltext，那么也不支持instant addcolumn

在数据字典表空间里的表也不支持instant add column

临时表也不支持instant add column
<!-- more -->
官方文档列出了一些可以快速ddl的操作，大体包括:

修改索引类型
Add column (limited) 当一条alter语句中同时存在不支持instant的ddl时，则无法使用 只能顺序加列 不支持压缩表 不支持包含全文索引的表 不支持临时表，临时表只能使用copy的方式执行DDL 不支持那些在数据词典表空间中创建的表
修改/删除列的默认值
修改索引类型
修改ENUM/SET类型的定义 存储的大小不变时 向后追加成员
增加或删除类型为virtual的generated column
RENAME TABLE操

MySQL8.0.12引入的快速加列特性能解决很大部分加列带来的问题：

对超级大表的加列操作通常可能耗时几个小时甚至数天的时间
在ddl的过程中产生的临时表会占用磁盘空间
ddl带来的复制延迟问题
具体的worklog为: WL#11250 - Support Instant Add Column,
使用
ALTER语句增加了新的语法INSTANT，你可以显式地指定，但MySQL自身也会自动选择合适的算法。所以这个特性通常对用户是透明的。
增加了一些新的information_schema表来展示相关信息:
I_S.innodb_tables.instant_cols
I_S.innodb_columns.has_default/default_value


root@test 03:54:47>show create table t1\G
*************************** 1. row ***************************
       Table: t1
Create Table: CREATE TABLE `t1` (
  `a` int(11) NOT NULL,
  `b` int(11) DEFAULT NULL,
  `c` int(11) DEFAULT NULL,
  `d` int(11) DEFAULT '11',
  PRIMARY KEY (`a`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
1 row in set (0.00 sec)
 
 
c,d 是instant added...
 
root@test 03:55:40>select table_id, name, pos, len, has_default,default_value from information_schema.innodb_columns where has_default = 1\G
*************************** 1. row ***************************
     table_id: 1143
         name: c
          pos: 2
          len: 4
  has_default: 1
default_value: NULL
*************************** 2. row ***************************
     table_id: 1143
         name: d
          pos: 3
          len: 4
  has_default: 1
default_value: 8000000b
2 rows in set (0.00 sec)
 
 
root@test 03:56:35>select instant_cols from information_schema.innodb_tables where name = 'test/t1';
+--------------+
| instant_cols |
+--------------+
|            2 |
+--------------+

实现
记录格式修改
快速加列特性，在增加列时，实际上只是修改了元数据，原来存储在文件中的行记录并没有被修改。当行格式为redundent类型时，记录解析是不依赖元数据的，可以自解析，
但如果行格式是dynamic或者compact类型，由于行内不存储元数据，尤其是列的个数信息，其记录的解析需要依赖元数据的辅助。因此为了支持动态加列功能，需要对行格式做一定的修改

其大体思路为:

如果表上从未发生过instant add column, 则行格式维持不变。
如果发生过instant ddl, 那么所有新的记录上都被特殊标记了一个flag, 同时在行内存储了列的个数
由于只支持往后顺序加列，通过列的个数就可以知道这个行记录中包含了哪些列的信息
我们先来看看典型compact行类型的记录组织结构:

+--------------------------------+---------------------+---------------+
| Non-null variable-length array | SQL-null flags/bitmap | Extra 5 bytes |


其中extra 5 bytes包含如下信息:

+-----------+---------------+----------+-------------+-----------------+
| Info bits | Records owned | Heap No. | Record type | Next record ptr |
+-----------+---------------+----------+-------------+-----------------+

extra info中包含的信息如下:
a) Info bits: 4 bits

0x10: REC_INFO_MIN_REC_FLAG
0x20: REC_INFO_DELETED_FLAG
其中还有两个bit是未使用的
b) Record owned: 4 bits
REC_NEW_N_OWNED

c) Heap No. : 13 bits

d) Record type: 3 bits

REC_STATUS_ORDINARY:叶子节点记录
REC_STATUS_NODE_PTR:非叶子节点记录
REC_STATUS_INFIMUM/REC_STATUS_SUPREMUM 系统记录
e) Next record ptr: 2 bytes

为了支持instant add column, 使用了info bits中的一个bit位，如果被设置，表示这条记录是第一次instant add column后插入的, flag为:
Ox80: REC_INFO_INSTANT_FLAG

当flag被设置时，在记录中就会使用１或２个字节来存储列的个数

+--------------------------+----------------+---------------+---------------+
| Non-null variable-length |                |               |               |
| array                    | SQL-null flags | fields number | Extra 5 bytes |
+--------------------------+----------------+---------------+---------------+
对于redundent类型，由于已经有了列个数信息，无需进行修改

数据词典信息
对数据词典进行了扩展并记录:

在第一次instant add column之前的列个数
每次加的列的默认值
通过这些信息加上记录上的额外信息，可以正确解析出记录上的数据

数据词典:
a) dd::Table::se_private_data::instant_col:在第一次instant ADD COLUMN之前表上面的列的个数
b) dd::Partition::se_private_data::instant_col, 和a类似，存储分区表上instant col的个数，但有所不同的是，分区表上的分区之间
可能存在不同列的个数。因为我们单独truncate一个分区，而truncate操作会清空instant标记，因此b)中存储的instant_col不应该比a)中每个分区上的instant_col要小 
c) dd::Column::se_private_data::default_null, 表示默认值为NULL
d) dd::Column::se_private_data::default, 当默认值不为null时，这里存储默认值
DD_instant_col_val_coder
--- column default value需要从innodb类型byte转换成se_private_data中的text类型(char), 使用一个类型DD_instant_col_val_coder来辅助转换
example: 0XFF => 0x0F, 0x0F

在将表load到内存建立表对象dict_table_t和索引对象dict_index_t时，有几个关键成员要载入进来，因为会用于辅助解析记录

dict_table_t::n_instant_cols 第一次instant add column之前的非虚拟列个数，（包含系统列
dict_index_t::instant_cols flag用于标示是否存在Instant column
dict_index_t::n_instant_nullable： 第一次instant add column之前的可为null的列个数
dict_col_t::instant_default: 存储默认值及其长度, 当解析数据时看到Instant column, 会直接引用到这里的数据指针

载入逻辑:

ha_innobase::open
|-->dd_open_table
    |--> dd_open_table_one
上述提到的几个变量会被设置.
DDL
检查表是否支持instant ddl

ha_innodb::check_if_supported_inplace_alter()
    innobase_support_instant
    innopart_support_instant
            dict_table_t::support_instant_add()
condition:

不是压缩表
不是data dictionary tablespace
不是全文索引表
不是临时表
除此之外, 新增列还要确保不改变列的顺序

当判定可以立刻加列时，仅仅需要修改数据词典信息即可

ha_innobase::commit_inplace_alter_table
|--> dd_commit_inplace_instant
    |--> dd_commit_instant_part
    |--> dd_commit_instant_table
     1. dd::TABLE中记录instant column的个数
     2. 存储新的列的默认值
Note:

truncate操作会重置instant标记
ha_innobase::truncate_impl
    dd_clear_instant_table
2.重建表的话，新不的表将不包含instant列

select
查询的关键在于如何正确的解析出记录中的每一行(对于不在其中的instant column，填默认值即可), 关键的函数是:

rec_init_offsets
|-->rec_init_offsets_comp_ordinary
    |-->rec_init_null_and_len_comp
何时填写instant add的default值:
default值存储在dict_col_t::dict_col_default_t

将默认值填到返回的记录中:

row_sel_store_mysql_field_func
    rec_get_nth_field_instant
    rec_get_nth_field_instant: 封装了列值： 如果是记录中的，则从记录中读取，否则返回其默认值
insert
记录在插入之前从tuple转换成physical record:

rec_convert_dtuple_to_rec
    rec_convert_dtuple_to_rec_new
           rec_convert_dtuple_to_rec_comp
当表上有instant column时
   1. 会占用1（如果列个数小于REC_N_FIELDS_ONE_BYTE_MAX）或者2个字节来存储列个数
   2. 在记录的Info bits字段设置REC_INFO_INSTANT_FLAG，表示这个记录是instant add column之后创建的
update
对于update，不会把default的值转换成inline的，除非去更新包含default值的列（row_upd_changes_field_size_or_external）

对于update的回滚做了特殊处理：

如果回滚的值从non-default到default值，那么这个是不会存储到列里面去的。（dtuple_t::ignore_trailing_default()）
