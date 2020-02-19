---
title: 使用Browserify解决浏览器加载require没有被定义
layout: post
category: web
author: 夏泽民
---
https://github.com/browserify/browserify

Nodejs的模块是基于CommonJS规范实现的，可不可以应用在浏览器环境中呢？

var math = require('math');
math.add(2, 3);
　　第二行math.add(2, 3)，在第一行require('math')之后运行，因此必须等math.js加载完成。也就是说，如果加载时间很长，整个应用就会停在那里等。这对服务器端不是一个问题，因为所有的模块都存放在本地硬盘，可以同步加载完成，等待时间就是硬盘的读取时间。但是，对于浏览器，这却是一个大问题，因为模块都放在服务器端，等待时间取决于网速的快慢，可能要等很长时间，浏览器处于"假死"状态

　　而browserify这样的一个工具，可以把nodejs的模块编译成浏览器可用的模块，解决上面提到的问题。本文将详细介绍Browserify

 

实现
　　Browserify是目前最常用的CommonJS格式转换的工具

　　请看一个例子，b.js模块加载a.js模块

复制代码
// a.js
var a = 100;
module.exports.a = a;

// b.js
var result = require('./a');
console.log(result.a);
复制代码
　　index.html直接引用b.js会报错，提示require没有被定义

复制代码
//index.html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
</head>
<body>
<script src="b.js"></script>    
</body>
</html>
<!-- more -->
这时，就要使用Browserify了

【安装】

　　使用下列命令安装browserify

npm install -g browserify
【转换】

　　使用下面的命令，就能将b.js转为浏览器可用的格式bb.js

$ browserify b.js > bb.js
　　查看bb.js，browserify将a.js和b.js这两个文件打包为bb.js，使其在浏览器端可以运行
　　
原理
　　Browserify到底做了什么？安装一下browser-unpack，就能清楚原理了

$ npm install browser-unpack -g
　　然后，使用下列命令，将前面生成的bb.js解包

$ browser-unpack < bb.js

　　可以看到，browerify将所有模块放入一个数组，id属性是模块的编号，source属性是模块的源码，deps属性是模块的依赖

　　因为b.js里面加载了a.js，所以deps属性就指定./a对应1号模块。执行的时候，浏览器遇到require('./a')语句，就自动执行1号模块的source属性，并将执行后的module.exports属性值输出

　　browerify将a.js和b.js打包，并生成bb.js，browser-unpack将bb.js解包，是一个逆向的过程。但实际上，bb.js依然存在
　　
Browserify
browserify的官网是http://browserify.org/，他的用途是将前端用到的众多资源（css,img,js,...) 打包成一个js文件的技术。

比如在html中引用外部资源的时候，原来我们可能这样写

  <script src="/static/libs/landing/js/bootstrap.min.js"></script>
  <script src="/static/libs/landing/js/jquery.flexslider-min.js"></script>
  <script src="/static/libs/landing/js/jquery.nav.js"></script>
  <script src="/static/libs/landing/js/jquery.appear.js"></script>
  <script src="/static/libs/landing/js/headhesive.min.js"></script>
  <script src="/static/libs/jquery/jquery-qrcode/jquery.qrcode.js"></script>
  <script src="/static/libs/jquery/jquery-qrcode/qrcode.js"></script>
  <script src="/static/libs/landing/js/scripts.js"></script>
但是有了 browserify 的帮助，就可以把这些通通压缩成一句

<script src="/bundle.js"></script>
而且不用担心，jQuery或者underscore等等库的冲突问题。

虽然这项技术也是最近几年才流行起来的，但是它迅速的在前端领域流行了起来。另一个跟browserify比较类似的是webpack，但这篇文章不打算介绍它，因为主页感觉不如browserify做的专业。

安装
安装起来很简单，不过首先你还需要需要先把nodejs装上。

npm install -g browserify
借助browserify你可以使用nodejs中常用到的require, module.exports功能。

简单入门
来个很简单的例子。

先创建一个hello.js文件，内容如下

module.exports = 'Hello world';
然后在创建一个entry.js文件，内容

var msg = require('./hello.js')console.log("MSG:", msg);
最后使用browserify进行进行打包

browserify entry.js -o bundle.js
然后entry.js和hello.js就被打包成一个bundle.js文件了。

写一个简单的index.html验证下效果

<!DOCTYPE html><html>
    <head>
        <meta charset="utf-8" />
        <title>index</title>
    </head>
    <body>
        <script src="bundle.js"></script>
    </body></html>
然后用浏览器打开该文件，F12开启调试选项。应该会看到有一句MSG: Hello world被打印出来了。

这就是最简单的打包应用了。

打包npm中的库
先创建一个package.json文件，内容最简单的写下。

 {
     "name": "study_browserify"
 }
接着安装jquery库

 npm i --save jquery
其中--save的意思是将jquery信息保存到package.json文件中。

修改下我们之前创建的hello.js文件成

 module.exports = function(){     var $ = require('jquery')
     $(function(){
         $("body").html("Hello world, jquery version: " + $.fn.jquery);
     })
 };
entry.js文件也要稍微修改下

 var hello = require('./hello.js')
 hello()
查看效果

这时打开index.html,你会看到页面上有字了，出现了Hello world, jquery version ....

这样子做的好处有很多，即使这个页面你又引用了别的jquery也不会和hello.js里面使用到的冲突。因为打包出来的bundle.js把引用的jquery放到的局部变量里面。

利用gulp工具自动打包
gulp也是前端开发人员常用的一个工具，用起来效果就像Makefile似的。gulp的主页是http://gulpjs.com/ 主页那叫一个简洁。

gulp的配置文件是gulpfile.js,按照我提供的内容先创建一个，具体怎么使用可以之后再去看官网。

var gulp = require('gulp');var rename = require('gulp-rename');var browserify = require('gulp-browserify');

gulp.task('build', function(){    return gulp.src('./entry.js')
        .pipe(browserify({
        }))
        .pipe(rename('bundle.js'))
        .pipe(gulp.dest('./'))
});

gulp.task('default', ['build'], function(){    return gulp.watch(['./*.js', '!./bundle.js'], ['build'])
})
之后安装下依赖库

npm i -g gulpnpm i --save-dev gulp gulp-rename gulp-browserify
当前目录下启动gulp，效果就是每次你修改了js文件，都会调用browserify打包一次。

打包HTML资源
这个时候用到了另外一个库 stringify,有了这个库的帮忙，就可以这么着用require("./hello.html") 是不是很酷炫。

首先还是安装 npm i --save-dev stringify

之后需要稍微修改下gulpfile.js

原来这个样子

gulp.task('build', ['lint'], function(){    return gulp.src('./entry.js')
        .pipe(browserify({ 
        })) 
        .pipe(rename('bundle.js'))
        .pipe(gulp.dest('./'))
});
增加几行代码，需要改造成这样. 第一行的require可以放到最上面。

var stringify = require('stringify');

gulp.task('build', ['lint'], function(){    return gulp.src('./entry.js')
        .pipe(browserify({
            transform: [
                stringify(['.html']),
            ],  
        })) 
        .pipe(rename('bundle.js'))
        .pipe(gulp.dest('./'))
});
为了验证效果。我们添加一个文件 hello.html

内容简单的写下

<strong>Hello</strong><span style="color:blue">World</span>
接着修改下hello.js,改成

module.exports = function(){    var $ = require('jquery')
    $(function(){
        $("body").html(require('./hello.html'));
    })  
};
重新打包，并再次刷新index.html 那个网页，就可以看到加粗的Hello，以及变蓝的World了。

添加静态代码检查
默认情况下，出现的一些低级错误，browserify是检查不到的。此时可以用js比较流行的代码检查工具jshint,官网是 http://jshint.com/

jshint相比较jslint配置少了不少，不过依然很多，闲麻烦的话，可以直接用我的。 下面的内容直接保存为文件 .jshintrc. 注意前面有个.

{
  "camelcase": true,
  "curly": true,
  "freeze": true,
  "immed": true,
  "latedef": "nofunc",
  "newcap": false,
  "noarg": true,
  "noempty": true,
  "nonbsp": true,
  "nonew": true,
  "undef": true,
  "unused": true,
  "trailing": true,
  "maxlen": 120,
  "asi": true,
  "esnext": true,
  "laxcomma": true,
  "laxbreak": true,
  "node": true,
  "globals": {
    "describe": false,
    "it": false,
    "before": false,
    "beforeEach": false,
    "after": false,
    "afterEach": false,
    "Promise": true
  }}
之后修改gulpfile.js文件为

var gulp = require('gulp');var rename = require('gulp-rename');var browserify = require('gulp-browserify');var jshint = require('gulp-jshint');

gulp.task('build', ['lint'], function(){    return gulp.src('./entry.js')
        .pipe(browserify({
        }))
        .pipe(rename('bundle.js'))
        .pipe(gulp.dest('./'))
});

gulp.task('lint', ['jshint'])
gulp.task('jshint', function(){    return gulp.src(['./*.js', '!./bundle.js'])
        .pipe(jshint())
        .pipe(jshint.reporter('jshint-stylish'))
})

gulp.task('default', ['build'], function(){    return gulp.watch(['./*.js', '!./bundle.js'], ['build'])
})
然后安装几个新增的依赖

npm i --save-dev gulp-jshint jshint jshint-stylish
重新运行gulp, 然后故意把entry.js文件改的错一点。你就会发现编辑器开始提示你错误了。


