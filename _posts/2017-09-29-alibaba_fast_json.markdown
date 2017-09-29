---
title: alibaba_fast_json
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
Fastjson是一个Java语言编写的高性能功能完善的JSON库。将解析json的性能提升到极致，是目前Java语言中最快的JSON库。Fastjson接口简单易用，已经被广泛使用在缓存序列化、协议交互、Web输出、Android客户端等多种应用场景。

GitHub下载地址: 
https://github.com/alibaba/fastjson

最新发布版本jar包 1.2.23 下载地址: https://search.maven.org/remote_content?g=com.alibaba&a=fastjson&v=LATEST
{% highlight scala linenos %}
import com.alibaba.fastjson.JSON
object FastJsonExp {
  def main(args:Array[String]){
  val json="{\"name\":\"chenggang\",\"age\":24}";
 //反序列化
 var userInfo=JSON.parseObject(json);
 System.out.println("name:"+userInfo.get("name")+", age:"+userInfo.get("age"));
  }
}
{% endhighlight %}