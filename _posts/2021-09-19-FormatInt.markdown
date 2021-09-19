---
title: Int转string几种方式性能
layout: post
category: golang
author: 夏泽民
---
Go语言内置int转string至少有3种方式：

fmt.Sprintf("%d",n)

strconv.Itoa(n)

strconv.FormatInt(n,10)
<!-- more -->
strconv.FormatInt()效率最高，fmt.Sprintf()效率最低

https://blog.csdn.net/flyfreelyit/article/details/79701577

