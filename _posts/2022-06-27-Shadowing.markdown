---
title: Shadowing
layout: post
category: golang
author: 夏泽民
---
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest

什么是变量隐藏呢？
 
就是当年在后面重新声明了前面已经声明的同名变量时，后面的变量值会遮蔽前面的变量值，虽然这两个变量同名但值却不一样。这样是很容易产生问题的。
<!-- more -->
1，可以用vet,Go 1.12 以上的版本需要

go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
 
go vet -vettool=$(which shadow)

2，也可以用goland的tool工具
3，还可以用golangci-lint工具,也可以检测出来。
https://blog.csdn.net/dianxin113/article/details/115009686
