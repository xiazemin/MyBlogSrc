---
title: gometalinter
layout: post
category: golang
author: 夏泽民
---
GitHub地址：https://github.com/alecthomas/gometalinter

gometalinter安装和使用
1、安装

go get github.com/alecthomas/gometalinter
gometalinter --install --update
2、使用

cd 到go项目下，执行 gometalinter ./...

即检查所有目录的go文件，此时vendor目录下的也会检测。

如果是想指定指定目录，执行gometalinter + 文件夹名。

 

goland集成
1、File Watchers开启



2、引用gometalinter



gofmt 保存的时候自动 格式化go代码

goimports  保存的时候自动导入处理包

gometalinter 保存的时候自动检查go语法
<!-- more -->
https://github.com/alecthomas/gometalinter

https://github.com/rshipp/awesome-malware-analysis

SonarQube 是一个开源的代码分析平台, 用来持续分析和评测项目源代码的质量。 通过SonarQube我们可以检测出项目中重复代码， 潜在bug， 代码风格问题，缺乏单元测试等问题， 并通过一个web ui展示出来。

vscode集成gometalinter
vscode 默认使用的是golint，如果想用gometalinter替换golint，直接打开
设置项，
在用户设置里添加"go.lintTool": "gometalinter"即可。

https://segmentfault.com/a/1190000013553309?utm_source=tag-newest
