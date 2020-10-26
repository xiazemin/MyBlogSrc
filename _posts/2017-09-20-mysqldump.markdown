---
title: mysqldump
layout: post
category: storage
author: 夏泽民
---
<!-- more -->
mysqldump  -P端口  -hIP -u用户名 -p密码 表名 库名 > 目标文件.sql

mysqldump: [Warning] Using a password on the command line interface can be insecure.

mysqldump  -P端口  -hIP -u用户名 -p 表名 库名 > 目标文件.sql

然后输入密码