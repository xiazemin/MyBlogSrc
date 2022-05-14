---
title: MySQL中的Text类型
layout: post
category: mysql
author: 夏泽民
---
TEXT类型是一种特殊的字符串类型，包括TINYTEXT、TEXT、MEDIUMTEXT和LONGTEXT

以上各类型无须指定长度！
允许的长度是指实际存储的字节数，而不是实际的字符个数，比如假设一个中文字符占两个字节，那么TEXT 类型可存储 65535/2 = 32767 个中文字符，而varchar(100)可存储100个中文字符，实际占200个字节，但varchar(65535) 并不能存储65535个中文字符，因为已超出表达范围。

char长度固定， 即每条数据占用等长字节空间；适合用在身份证号码、手机号码等定。超过255字节的只能用varchar或者text。
varchar可变长度，可以设置最大长度；适合用在长度可变的属性。
text不设置长度， 当不知道属性的最大长度时，适合用text， 能用varchar的地方不用text。
如果都可以选择，按照查询速度： char最快， varchar次之，text最慢。
<!-- more -->
https://blog.csdn.net/SlowIsFastLemon/article/details/106383776