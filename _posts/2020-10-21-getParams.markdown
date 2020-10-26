---
title: getParams yaf
layout: post
category: php
author: 夏泽民
---
get('name') //获取参数(不仅仅是get方法，也可以是post方法)，没有返回NULL,需要传入一个参数名，字符串形式，也只能获取到单个的参数
getPost() //获取post参数，
getQuery() //获取url地址及参数,不需要传入参数 /User/User/index/name/huyouheng/age/23
getParam('name') //得到指定的参数
getParams() //得到传入的所有参数
getRequestUri() //得到请求的url,其实得到的和 getQuery()一致的
getMethod() //得到请求的方法
getFiles() //上传的文件
<!-- more -->
https://www.jianshu.com/p/3157563e87be
http://www.shixinke.com/php/yaf-request-and-response


