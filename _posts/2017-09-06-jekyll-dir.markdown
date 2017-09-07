---
layout: post
title:  "jekyll mac 安装"
date:   2017-08-05
category: jekyll
tags: [octopress, jekyll]

# Author.
author: 夏泽民
---
Jekyll 的核心是一个文本转换引擎。它的核心是把你零散的文件、文本组合起来，形成一个个网页，最终呈现在浏览器上展现出来。

一个最基础的 Jekyll 博客，会拥有下面的目录结构：

.
├── _config.yml
├── _drafts
|   ├── begin-with-the-crazy-ideas.textile
|   └── on-simplicity-in-technology.markdown
├── _includes
|   ├── footer.html
|   └── header.html
├── _layouts
|   ├── default.html
|   └── post.html
├── _posts
│   └── 2013-08-07-welcome-to-jekyll.markdown
├── _site
└── index.html
这些目录的介绍如下：

文件/目录	描述
_config.yml	存储配置数据。很多全局的配置或者指令写在这里。
_drafts	存放为发表的文章。这些是没有日期的文件。
_includes	存放一些组件。可以通过{\% include file.ext \%} 来引用。
_layouts	布局。
_posts	存放写文章，格式化为：YEAR-MONTH-DAY-title.md。
_site	最终生成的博客文件就在这里。
index.html	博客的主页。
other	例如静态文件 CSS，Images 和其他。
只要我们把自己需要的文件放到博客目录下，通过jekyll build，该目录就会被复制到_site里面。