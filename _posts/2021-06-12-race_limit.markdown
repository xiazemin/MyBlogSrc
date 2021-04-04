---
title: 内存预分配与race limit on 8128 
layout: post
category: golang
author: 夏泽民
---
go test -count=1 -race ./... -v

package utils

var buf = make([]byte, 5010241024) // 50M static buffer
 
起10000个协程

这个没有用的变量干掉了竟然是race: limit on 8128 simultaneously alive goroutines is exceeded, dying的原因

https://ifun.dev/post/golang-concurrency/
<!-- more -->
{% raw %}

{% endraw %}
<div class="container">
	<div class="row">
	<img src="{{site.url}}{{site.baseurl}}/img/jupyterSlider.png"/>
	</div>
	<div class="row">
	</div>
</div>
