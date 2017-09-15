---
title: markdown-table
layout: post
category: jekyll
author: 夏泽民
---
1. 方案一：
markdown原生语法可以生成表格,在字段左右加｜分隔，第二行 | -   |-:| :-----: |
例如
<code>

| 水果|价格|数量|

    | -   |-:| :--: |
    
    | 香蕉        | $1      |   5    |
    
    | 苹果        | $1      |   6    |
    
    | 草莓        | $1      |   7    |
    
<code>

显示效果

| 水果|价格|数量|
    | -   |-:| :--: |
    | 香蕉        | $1      |   5    |
    | 苹果        | $1      |   6    |
    | 草莓        | $1      |   7    |


|Name|Academy     |  score| 
| - |:-| -- | 
|Harry Potter     | Gryffindor    | 90 |
|Hermione Granger | Gryffindor | 100 |
|Draco Malfoy     | Slytherin  | 90|

<!-- more -->
1. 方案二
markdown和html语法兼容，可以使用html的table
例如：
{% highlight html linenos %}
<table>
        <tr>
            <th>设备</th>
            <th>设备文件名</th>
            <th>文件描述符</th>
            <th>类型</th>
        </tr>
        <tr>
            <th>键盘</th>
            <th>/dev/stdin</th>
            <th>0</th>
            <th>标准输入</th>
        </tr>
</table>
{% endhighlight %}
<table>
        <tr>
            <th>设备</th>
            <th>设备文件名</th>
            <th>文件描述符</th>
            <th>类型</th>
        </tr>
        <tr>
            <th>键盘</th>
            <th>/dev/stdin</th>
            <th>0</th>
            <th>标准输入</th>
        </tr>
</table>
1. 方案三
excel转markdown工具
https://link.zhihu.com/?target=http://fanfeilong.github.io/exceltk0.0.4.7z
