---
layout:  post
title:  "jekyll layout"
date:   2017-08-05
category: jekyll
tags: [octopress, jekyll]

# Author.
author: 夏泽民
---

layout: default

作用是通过 layout 告诉Jekyll生成 index.html 时，要在_layouts 目录下找 default.html 文件，然后把当前文件解析后添加到 default.html 的content 的部分，组成最终的 index.html 文件。

在我们写的markdown文档中也要设置YAML头信息，如我的这篇博文的头信息：

layout: post 
title: “Jekyll和Github搭建个人静态博客” 
date: 2016/6/26 13:03:42

categories: original

layout表示使用post布局，title 是文章标题，date是自动生成的日期，categories 是该文章生成html文件后的存放目录，也就是文章的分类属性。可以在_site/original下找到。（category 只能添加一个分类属性， categories 可以添加多个分类属性。各属性使用空格隔开）

因为文章套用的是post模板，所以title会传入 post.html 文件中的Jekyll和Github搭建个人静态博客中，成为最终 index.html 页面中的文章列表标题：

author

而 post.html 又套用了default.html 模板，而default页面中的头部又由 head.html 构成：


head页面中的title属性：

Jekyll和Github搭建个人静态博客 
就可以读取到这篇博文中的title并且设置在最终 index.html 文件中。

_posts

这个目录存放的就是我们所有的博文了。文件名字格式很重要，必须使用统一的格式：

YEAR-MONTH-DAY-title.MARKUP 
例如，2016-06-26-MakeBlog.md，写成这样文件名才会被解析。

