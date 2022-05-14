---
title: go-callvis
layout: post
category: golang
author: 夏泽民
---
一、 https://github.com/TrueFurby/go-callvis

二、quick start，详见Readme文件

三、注意事项

1. 安装完成go-callvis在$GOPATH/bin下，需要加入到PATH，否则找不到

2. go get 失败的可能要手动下，或者配置ssh

四、run

1. go-callvis [flags] <main package>

2. 访问 http://localhost:7878/
<!-- more -->
https://www.ctolib.com/mip/go-callvis.html

使用
go-callvis github.com/项目具体路径 | dot -Tpng -o syncthing.png

解析的是main包
go-callvis -group pkg,type -focus [想要分析的包（确定在后面的路径中）] github.com/项目具体路径 | dot -Tpng -o syncthing.png

https://studygolang.com/articles/17941


~/goLang$go-callvis github.com/xiazemin/ast/tag |dot -Tpng -o tag.png

/usr/local/go/src/syscall/zsyscall_darwin_amd64.go:137:134: too few arguments in call to syscall6
/usr/local/go/src/syscall/zsyscall_darwin_amd64.go:197:142: too few arguments in call to rawSyscall6
/usr/local/go/src/syscall/zsyscall_darwin_amd64.go:218:181: too few arguments in call to syscall6
/usr/local/go/src/syscall/zsyscall_darwin_amd64.go:240:145: too few arguments in call to syscall6
/usr/local/go/src/syscall/zsyscall_darwin_amd64.go:287:169: too few arguments in call to syscall6

https://www.kutu66.com//GitHub/article_149310

其中src目录是一个go package，运行go-callvis 时就需要先cd src/，然后再执行命令：

go-callvis  -group pkg,type md52id
复制代码
md52id 是package name，已在go.mod中声明，pakage name是一个必须要带的参数。

运行命令，默认会打开浏览器加载地址http://localhost:7878

https://www.lagou.com/lgeduarticle/96545.html

go-callvis 是一个开发工具，其目的是通过使用来自函数调用关系图的数据及其与包和类型的关系来对程序进行可视概览。 这在你只是试图理解别人的代码结构，或在代码复杂性增加的大型项目中特别有用。

[TOC]

缺点
github项目上的文档写的不是很清晰，我尝试了一下，没用
图画的很乱，有时候完全摸不到头绪
目前版本不支持go module
官方示例
用法
github上的图例


三个例子
docker


使用
go-callvis github.com/项目具体路径 | dot -Tpng -o syncthing.png

解析的是main包
go-callvis -group pkg,type -focus [想要分析的包（确定在后面的路径中）] github.com/项目具体路径 | dot -Tpng -o syncthing.png

1、直接按照官网的命令安装的话：

go get -u github.com/TrueFurby/go-callvis
cd $GOPATH/src/github.com/TrueFurby/go-callvis && make
在第二个命令运行后会出现dep命令不存在的错误，也就是需要先安装dep；

2、现在安装dep，按照github官网安装：

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


注意安装过程中会出现的提示行：Will install into /root/go/bin

测试是否安装成功：输入命令dep，仍然找不到该命令；
再输入whereis dep命令，发现显示为空；找到原因：即要将dep的环境变量/root/go/bin添加到$PATH；
添加PATH命令：export PATH=$PATH:/root/go/bin；source $HOME/.profile
查看是否添加环境变量成功：echo $PATH，发现已经成功添加进去；


输入命令dep，发现安装成功：
https://blog.csdn.net/u011947630/article/details/89193706
https://studygolang.com/articles/26211

https://www.jianshu.com/p/bd43569d524b

$brew install graphviz

dyld: Library not loaded: /usr/lib/libcrypto.35.dylib
  Referenced from: /Library/Developer/CommandLineTools/usr/lib/libserf-1.0.dylib
  Reason: image not found
  

https://gitlab.com/graphviz/graphviz/
https://www.jianshu.com/p/669c6e61b1e7


/Applications/CMake.app/Contents/bin/cmake .


CMake Error: The following variables are used in this project, but they are set to NOTFOUND.
Please set them or make sure they are set and tested correctly in the CMake files:
LTDL_INCLUDE_DIR (ADVANCED)
   used as include directory in directory /Users/didi/goLang/src/github.com/graphviz/lib/gvc
   
   
https://github.com/xflr6/graphviz
http://graphviz.org/
http://macappstore.org/graphviz/

Error: Failed to download resource "netpbm"
Failure while executing; `svn checkout https://svn.code.sf.net/p/netpbm/code/stable /Users/didi/Library/Caches/Homebrew/netpbm--svn --quiet -r 3806` exited with . Here's the output:
dyld: Library not loaded: /usr/lib/libcrypto.35.dylib
  Referenced from: /Library/Developer/CommandLineTools/usr/lib/libserf-1.0.dylib
  Reason: image not found
  
  
  $otool -L /usr/lib/libcrypto.dylib
/usr/lib/libcrypto.dylib:
	/usr/lib/libcrypto.0.9.8.dylib (compatibility version 0.9.8, current version 0.9.8)
	/System/Library/PrivateFrameworks/TrustEvaluationAgent.framework/Versions/A/TrustEvaluationAgent (compatibility version 1.0.0, current version 25.0.0)
	/usr/lib/libz.1.dylib (compatibility version 1.0.0, current version 1.2.5)
	/usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1225.1.1)
	
	$  ls /usr/lib/libcrypto.
libcrypto.0.9.7.dylib  libcrypto.0.9.8.dylib  libcrypto.dylib

https://mithun.co/hacks/library-not-loaded-libcrypto-1-0-0-dylib-issue-in-mac/

解决方法：
brew switch openssl 1.0.2q
 
如果你不知道要切换为的openssl版本是什么也没关系，直接输入后会提示你已经安装的可用版本是多少

$  brew info graphviz
graphviz: stable 2.44.0, HEAD
Graph visualization software from AT&T and Bell Labs
https://www.graphviz.org/
Not installed
From: https://mirrors.aliyun.com/homebrew/homebrew-core.git/Formula/graphviz.rb
==> Dependencies
Build: pkg-config ✔
Required: gd ✘, gts ✘, libpng ✘, libtool ✘, pango ✘
==> Options
--HEAD
	Install HEAD version
==> Analytics
install: 34,884 (30 days), 133,517 (90 days), 460,710 (365 days)
install-on-request: 26,975 (30 days), 100,475 (90 days), 339,821 (365 days)
build-error: 0 (30 days)


https://github.com/Homebrew/homebrew-core/blob/master/Formula/graphviz.rb

    url "https://www2.graphviz.org/Packages/stable/portable_source/graphviz-2.44.0.tar.gz"
    
    https://www2.graphviz.org/Packages/stable/portable_source/
    
    
vi graphviz.rb

把内容copy过来
$brew unlink graphviz
$brew  install graphviz.rb
Updating Homebrew...

https://segmentfault.com/a/1190000018875410
https://www.jianshu.com/p/aadb54eac0a8

https://www.mobibrw.com/2014/1268

$brew tap homebrew/versions
Updating Homebrew...
^C
Error: homebrew/versions was deprecated. This tap is now empty as all its formulae were migrated.

https://www.cnblogs.com/weiki-nttdata/p/5080746.html

https://graphviz.org/download/source/

https://www2.graphviz.org/Packages/stable/portable_source/

https://github.com/angelj-a/axel
$./configure
$make
$make install


$axel -n 20 https://www2.graphviz.org/Packages/stable/portable_source/graphviz-2.42.1.tar.gz

$git tag
2.38.0
2.42.2
2.42.3
2.42.4
2.44.0
2.44.1
LAST_LIBGRAPH
Nightly

$git checkout 2.44.0

$automake
-bash: automake: command not found

$autoconf
/Library/Developer/CommandLineTools/usr/bin/m4:configure.ac:18: cannot open `./version.m4': No such file or directory
autom4te: /usr/bin/m4 failed with exit status: 1

Makefile.am: 是一些编译的选项及要进行编译的文件项等，例如：
bin_PROGRAMS=test
lib_LIBRARIES = libhand.a
libhand_a_SOURCES = hand.c

Makefile.in: 在automake手册中是这样说：while automake is in charge of creating Makefile.ins from Makefile.ams and configure.ac. 意思是Makefile.in是由Makefile.am和configure.ac的基础之上而生成的。

Makefile: 使用生成的configure脚本根据Makefile.in中的内容进行生成的最终进行程序或库编译的规则文件；


cd graphviz-2.42.1/
./configure
make
make install

$dot -h
Error: dot: option -h unrecognized

Usage: dot [-Vv?] [-(GNE)name=val] [-(KTlso)<val>] <dot files>
(additional options for neato)    [-x] [-n<v>]
(additional options for fdp)      [-L(gO)] [-L(nUCT)<val>]
(additional options for memtest)  [-m<v>]
(additional options for config)  [-cv]

 -V          - Print version and exit
 -v          - Enable verbose mode
 -Gname=val  - Set graph attribute 'name' to 'val'
 -Nname=val  - Set node attribute 'name' to 'val'
 
 安装成功
 
 
 /usr/local/go/src/fmt/format.go:9:2: could not import unicode/utf8 (invalid package name: "")
 
 ~/goLang$go build -o go-callvis github.com/ofabry/go-callvis/
 
 https://ofabry.github.io/go-callvis/
 https://github.com/ofabry/go-callvis
 
 
 $go-callvis github.com/xiazemin/ast/tag
2020/07/04 23:18:08 http serving at http://localhost:7878
2020/07/04 23:18:08 converting dot to svg..
2020/07/04 23:18:08 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/go-callvis_export.svg

$go-callvis github.com/ofabry/go-callvis/
2020/07/04 23:20:07 http serving at http://localhost:7878
2020/07/04 23:20:09 converting dot to svg..
2020/07/04 23:20:09 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/go-callvis_export.svg
   
   
$go-callvis github.com/xiazemin/ast/tag |dot -Tpng -o tag.png
Format: "png" not recognized. Use one of: canon cmap cmapx cmapx_np dot dot_json eps fig gv imap imap_np ismap json json0 mp pic plain plain-ext pov ps ps2 svg svgz tk vdx vml vmlz xdot xdot1.2 xdot1.4 xdot_json

go-callvis github.com/xiazemin/ast/tag |dot -Tsvg -o tag.svg
2020/07/04 23:22:44 http serving at http://localhost:7878
2020/07/04 23:22:45 converting dot to svg..
2020/07/04 23:22:45 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/go-callvis_export.svg

$go-callvis -format "svg" -file "./tag.svg" github.com/xiazemin/ast/tag

$go-callvis -format "svg" -file "./tag.svg" github.com/ofabry/go-callvis/
2020/07/04 23:31:11 writing dot output..
2020/07/04 23:31:11 converting dot to svg..


$go-callvis -include "github.com/beego/bee/cmd,github.com/beego/bee/cmd/commands,github.com/beego/bee/config,github.com/beego/bee/generate/swaggergen,github.com/beego/bee/utils" github.com/beego/bee
2020/07/04 23:55:28 http serving at http://localhost:7878
2020/07/04 23:55:30 converting dot to svg..
2020/07/04 23:55:30 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/go-callvis_export.svg

只有main

$go-callvis -focus "github.com/beego/bee/logger" github.com/beego/bee
2020/07/05 00:00:08 http serving at http://localhost:7878
2020/07/05 00:00:09 converting dot to svg..
2020/07/05 00:00:09 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/go-callvis_export.svg


重点是logger包

https://www.ctolib.com/mip/go-callvis.html
