---
layout: post
title:  "jekyll macdown使用"
date:   2017-08-05
category: jekyll
tags: [octopress, jekyll]

# Author.
author: 夏泽民
---
Jekyll中一篇文章就是一个文件，所有需要发布的文章都要放在_posts文件夹内。Jekyll对于文章的文件名也是有要求的，系统会根据文件名来生成每篇文章的链接地址。具体的格式为：YYYY-MM-DD-文章标题.markdown 其中YYYY为4位年份，MM是两位的月份，DD是两位的日期。

在使用Markdown撰写文章之前我们需要先设置头信息。头信息需要根据YAML的格式写在两行三虚线之间。在头信息可以设置预定义的全局变量的值，Jekyll会根据变量的值来生成文章页面。

layout使用指定的模版文件，不加扩展名。模版文件放在_layouts目录下。

title文章的标题。

date发布文章的时间。

categories将文章设置成不同的属性。系统在生成页面时会根据多个属性值来生成文章地址。以上设置会生http://.../jekyll/update/...格式的文章链接。

tags标签，一篇文章可以设置多个标签，使用空格分割。


Jekyll是支持图片和其它资源文件的

$ls
404.html	Gemfile.lock	_posts		about.md
Gemfile		_config.yml	_site		index.md

Jekyll 使用 Liquid 模板语言，{{ page.title }} 表示文章标题，{{ content }} 表示文章内容

# 1.文件结构
_config.yml：用于保存配置。（jekyll会自动加载这些配置）

_includes文件夹：存放可以重复利用的文件，可以被其他的文件包含（方法：{\% include 文件名 \%}）

_layouts文件夹：存放模板文件（标签{{ content }}将content插入页面中）。

_posts文件夹：存放实际的博客文章内容（文件名格式：年-月-日-标题.md）。

_site文件夹：存放最终生成的文件（其他的目录都会被拷贝到最终文件的目录下。所以css,images等目录都可以放在根目录下）。

YAML头信息（可选的）：（文章只要包含YAML头，yekyll就会将其转换成html文件）设置一些预定义的变量，或你自己定义的变量。

# 2.常用命令（命令行输入）

$ jekyll build     ：当前文件夹中的内容将会生成到 ./site 文件夹中。

$ jekyll build --destination <destination>   ：当前文件夹中的内容会生成到指定文件夹中。

$ jekyll build --source <source>--destination <destination>  ：将指定源文件夹中的内容生成到指定文件夹中。

$ jekyll build --watch  ：查看更改，再生成。

$ jekyll serve      ：启动服务器，使用本地预览，运行在http://localhost:4000/。（jekyll集成了一个服务器）

$ jekyll serve --watch     ：先查看变更在启动服务器。

可以在_config.yml文件中添加配置，jekyll会自动获取其中的配置，例如：

source:_source

destination:_deploy

等同于命令：jekyll build --source _source --destination _deploy

# 3.jekyll原理

jekyll使用Liquid语言

Liquid语言使用2种标记（Output和Tag）：Output：{{content}}，Tag：{\% content \%}

Liquid过滤器：将左边字符串通过过滤器得到想要的结果并输出。

过滤器示例
Liquid的标准过滤器：

date - 格式化日期

capitalize - 将输入语句的首字母大写

downcase - 将输入字符串转为小写

upcase - 将输入字符串转为大写

first - 得到传递数组的第一个元素

last - 得到传递数组的最后一个元素

join - 将数组中的元素连成一串，中间通过某些字符分隔

sort - 对数组元素进行排序

map - 从一个给定属性中映射/收集一个数组

size - 返回一个数组或字符串的大小

escape - 对一串字符串进行编码

escape_once - 返回一个转义的html版本，而不影响现有的转义文本

strip_html - 去除一串字符串中的所有html标签

strip_newlines - 从字符串中去除所有换行符(\n)

newline_to_br - 将所有的换行符(\n)换成 html 的换行标记

replace - 匹配每一处指定字符串并替换，如 {{ 'foofoo' | replace:'foo','bar' }} #=> 'barbar'

replace_first - 匹配第一处指定的字符串并替换，如 {{ 'barbar' | replace_first:'bar','foo' }} #=> 'foobar'

remove - 删除每一处匹配字符串，如 {{ 'foobarfoobar' | remove:'foo' }} #=> 'barbar'

remove_first - 删除第一处匹配的字符串，如 {{ 'barbar' | remove_first:'bar' }} #=> 'bar'

truncate - 将一串字符串截断为x个字符

truncatewords - 将一串字符串截断为x个单词

prepend - 在一串字符串前面加上指定字符串，如 {{ 'bar' | prepend:'foo' }} #=> 'foobar'

append - 在一串字符串后面加上指定字符串，如 {{ 'foo' | append:'bar' }} #=> 'foobar'

minus - 减，如 {{ 4 | minus:2 }} #=> 2

plus - 加，如 {{ '1' | plus:'1' }} #=> '11', {{ 1 | plus:1 }} #=> 2

times - 乘，如 {{ 5 | times:4 }} #=> 20

divided_by - 除，如 {{ 10 | divided_by:2 }} #=> 5

split - 将一串字符串根据匹配模式分割成数组，如 {{ "a~b" | split\:~ }} #=> \['a','b'\]

modulo - 余数，如 {{ 3 | modulo:2 }} #=> 1
tag标签：

assign- 创建一个变量

capture- 块标记，把一些文本捕捉到一个变量中（如：把一系列字符串连接为一个字符串，并将其存储到变量中）

case- 块标记，标准的 case 语句

comment- 块标记，将一块文本作为注释

if- 标准的 if/else 块

unless- if 语句的简版

include- 包含其他的模板

raw- 暂时性的禁用的标签的解析（展示一些可能产生冲突的内容）

cycle- 用于循环轮换值，如颜色或 DOM 类

for- 用于循环 For loop（for 。。。 in 。。。  limit:int使你可以限制接受的循环项个数；offset:int可以可以让你从循环集合的第 n 项开始；reversed让你可以翻转循环）
jekyll新增的过滤器：

date_to_string - 日期转化为短格式

date_to_long_string - 日期转化为长格式

number_of_words - 统计字数（{{ page.content | number_of_words }}）

array_to_sentence_string - 数组转换为句子（列举标签时：{{ page.tags | array_to_sentence_string }}）

markdownify - 将makedown格式字符串转换成HTML

jsonify - data to JSON
jekyll新增标签：

highlight 语言 linenos（行号，可选）- 块标签，代码高亮 

post_url - 使用某篇博文的超链接（不需要写文件后缀）

gist - github gist显示代码（gist的介绍和使用 ）（{\% gist 5555251 \%}）

# 4.书写博客
引用图片或其他资源：新建一个文件夹存放，在博文中的引用方式：{{site.url}}表示站点的根目录
`![实例图片]({{site.url}}／assets/image.jpeg)`

其他的资源引用也是一样的。

# 5.创建博文目录

一个简单的例子，使用的是Liquid模板语言。
<ul>
{\% for post in site.posts \%}
<li>
<a href="{{post.url}}"> {{post.title}}</a>
</li>
{\% endfor \%}
</ul>

创建目录

# 6、分页

在_config.yml里边加一行，并填写每页需要几行：

paginate:5

对需要带有分页页面的配置： paginate_path:"blog/page:num"
blog/index.html将会读取这个设置，把他传给每个分页页面，然后从第2页开始输出到blog/page:num，:num是页码。如果有 12 篇文章并且做如下配置paginate: 5， Jekyll会将前 5 篇文章写入blog/index.html，把接下来的 5 篇文章写入blog/page2/index.html，最后 2 篇写入blog/page3/index.html。

# 7、草稿

草稿是你还在创作中不想发表的文章。

创建一个名为_drafts的文件夹