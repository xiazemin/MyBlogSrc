---
title: mysql获取表结构
layout: post
category: storage
author: 夏泽民
---
//SELECT TABLE_SCHEMA, TABLE_NAME FROM INFORMATION_SCHEMA.TABLES;
//SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='schema' AND TABLE_NAME='table';
//SHOW FULL COLUMNS FROM 'table';
//show create table 'table';
<!-- more -->
sqlingo 就是通过SHOW FULL COLUMNS FROM 'table'; 来获取表结构生成orm的结构体，来保证代码和数据库的一致性的。

基于此思路，开发了能够保持表结构一致性的sqlc
 https://github.com/xiazemin/sqlc
通过show create table 'table';
来获取表结构



https://blog.csdn.net/liushawn520/article/details/89010610

