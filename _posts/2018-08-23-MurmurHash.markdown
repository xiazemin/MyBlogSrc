---
title: MurmurHash
layout: post
category: algorithm
author: 夏泽民
---
https://github.com/aappleby/smhasher
看Jedis的主键分区哈希时，看到了名字很萌很陌陌的MurmurHash，谷歌一看才发现Redis，Memcached，Cassandra，HBase，Lucene都用它。
其实是 multiply and rotate的意思，因为算法的核心就是不断的"x *= m; x = rotate_left(x,r);"
<!-- more -->
MurmurHash算法：高运算性能，低碰撞率，由Austin Appleby创建于2008年，现已应用到Hadoop、libstdc++、nginx、libmemcached等开源系统。2011年Appleby被Google雇佣，随后Google推出其变种的CityHash算法。 
官方网站：https://sites.google.com/site/murmurhash/ 
MurmurHash算法，自称超级快的hash算法，是FNV的4-5倍。

