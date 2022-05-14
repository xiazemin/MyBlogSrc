---
title: browserify
layout: post
category: node
author: 夏泽民
---
http://browserify.org/#install
npm install -g browserify

Now recursively bundle up all the required modules starting at main.js into a single file called bundle.js with the browserify command:

browserify main.js -o bundle.js
<script src="bundle.js"></script>
<!-- more -->
Browserify 可以让你使用类似于 node 的 require() 的方式来组织浏览器端的 Javascript 代码，通过预编译让前端 Javascript 可以直接使用 Node NPM 安装的一些库。

browserify 是如何工作的
Browserify从你给你的入口文件开始,寻找所有调用require()方法的地方, 然后沿着抽象语法树,通过 detective 模块来找到所有请求的模块. (其实这个意思就是说,它require里还有require,还有require,所有的require像一棵树一样,然后沿着这棵树,通过detective来找到所有的模块)

每一个require()调用里都传入一个字符串作为参数,browserify把这个字符串解析成文件的路径然后递归的查找文件直到整个依赖树都被找到.

每个被require()的文件,它的名字都会被映射到内部的id,最后被整合到一个javascript文件中.

这就意味着最后打包生成了文件已经包含了所有能让你的应用跑起来所需要的东西.

查看更多browserify的用法,可以看文档的编译器管道部分.

https://www.cnblogs.com/liulangmao/p/4920534.html
