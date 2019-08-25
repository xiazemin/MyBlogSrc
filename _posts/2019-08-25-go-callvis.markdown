---
title: go-callvis 生成golang调用图
layout: post
category: golang
author: 夏泽民
---
安装graphviz
$ brew install graphviz
安装go-callvis
go get -u github.com/TrueFurby/go-callvis
cd $GOPATH/src/github.com/TrueFurby/go-callvis && make
用法
$ go-callvis [flags] package
<!-- more -->
https://github.com/TrueFurby/go-callvis


$     which go-callvis
/Users/didi/goLang/bin/go-callvis

以go-callvis项目为例
$go-callvis github.com/TrueFurby/go-callvis
2019/08/25 00:29:05 http serving at http://localhost:7878
2019/08/25 00:29:06 converting dot to svg..
2019/08/25 00:29:06 serving file: /var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn

浏览器打开http://localhost:7878/
就可以看到调用关系图

如果没有focus标识，默认是main
查看 github.com/uber/go-torch/pprof  的调用
$   go-callvis -focus github.com/uber/go-torch/pprof github.com/uber/go-torch

