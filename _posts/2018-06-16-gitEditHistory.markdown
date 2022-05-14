---
title: 查看某行代码的git修改历史
layout: post
category: web
author: 夏泽民
---
使用:  ~/gitEditHis.sh file code
<!-- more -->
{% highlight bash linenos %}
 #!/bin/bash
version=`git log $1  |grep commit |awk '{print $2}' `
for i in $version;
do
 echo $i ;
 git checkout $i;
git blame $1  |grep $2;
 done
git checkout master
{% endhighlight %}
