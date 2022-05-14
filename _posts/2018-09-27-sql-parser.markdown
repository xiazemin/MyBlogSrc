---
title: sql-parser
layout: post
category: storage
author: 夏泽民
---
c++
https://github.com/hyrise/sql-parser
http://torpedro.github.io/tech/c++/sql/parser/2016/02/27/c++-sql-parser.html
https://github.com/hyrise/sql-parser/tree/master/example
http://www.sqlparser.com/sql-parser-c.php
jave
https://github.com/JSQLParser/JSqlParser
go
https://github.com/xwb1989/sqlparser
GSP（全称General SQL Parser）。他是一款专业的SQL引擎，适用于各种数据库。 
http://www.sqlparser.com/
一、检查语法

我们先讲讲下面的代码做了哪些事： 
1. 定义一个简单的create语句（我们故意把name1的类型错误的设置成varchar2） 
2. 创建一个MySQL解析器实例 
3. 将sql语句传递给解析器 
4. 解析器开始检查语法 
5. 判断检查结果，0表示语法正确，1表示语法有错误，并获取返回的错误信息 
二、格式化SQL

通常我们编写的SQL语句有一些杂乱，虽然自己不觉得。在Navicat中我们可以很方便的用快捷键format，得到美观、便于阅读的SQL。用General SQL Parser，我们同样可以很容易做到。 
我们先讲讲下面的代码做了哪些事： 
1. 前三步和“检查语法”是一样的 
2. 然后调用解析方法，注意是parse()方法 
3. 实例化格式化工具类 
4. 用格式化工具工厂类格式化SQL，获取结果 
三、提取多条SQL

通常我们一个sql文件中，不单单只有一条sql语句，并且一般会有很多注释，那我们怎么提取出每一条SQL语句呢？直接split(“;”)，注释怎么处理呢？ 
用General SQL Parser，我们很简单便捷的做到。 
我们看看我们需要做哪些： 
1. 前三步还是和前面一样 
2. 从解析器中获取statements 
3. 遍历statements 
4. 获取的TCustomSqlStatement就是每条sql实例 
<!-- more -->
