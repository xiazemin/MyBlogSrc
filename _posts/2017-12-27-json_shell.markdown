---
title: json_shell
layout: post
category: web
author: 夏泽民
---
<!-- more -->
解析简单json
{% highlight shell linenos %}
 #!/bin/bash
s="{\"rv\":0,\"flag\":1,\"url\":\"http://www.jinhill.com\",\"msg\":\"test\"}"
parse_json(){
 #echo "$1" | sed "s/.*\"$2\":\([^,}]*\).*/\1/"
echo "${1//\"/}" | sed "s/.*$2:\([^,}]*\).*/\1/"
}
echo $s
value=$(parse_json $s "url")
echo $value
{% endhighlight %}
解析URL Query
{% highlight shell linenos %}
 #!/bin/bash
s="http://www.zonetec.cn/WlanAuth/portal.do?appid=aaaa&apidx=0"
parse(){
 echo $1 | sed 's/.*'$2'=\([[:alnum:]]*\).*/\1/'
}
value=$(parse $s "appid")
echo $value
{% endhighlight %}
