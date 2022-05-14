---
title: gometalinter
layout: post
category: golang
author: 夏泽民
---
当然代码质量分析是devops中持续集成部分非常重要的一个环节。涉及到团队协作的时候，很多公司会有自己的一套规则，最熟悉的是阿里巴巴的java代码参考手册，专家总结，大家按照规则去写代码。但是对于devops，有了规则远远不够，还需要提高代码检查的自动化程度，以及与其他的环节合作衔接上。
那么如果没有代码质量把控这个环节，那么可能带来以下问题：

抽象不够，大量重复的代码。
变量命名各种风格。驼峰派，下划线派等。
if else 嵌套太深。
不写注释，或是注释写的不合格。
函数里代码太长
潜在的bug，这就是各种语言的一些语法上的检查了。
最终导致维护困难。这就是程序员经常所说的前任留下的坑，往往结果就是新的接任者选择重构。其实对单位也是一直财力和人力的浪费。另外，更无法满足devops的快速上线，快速诊断，快速迭代的目标。
<!-- more -->
SonarQube 是一个开源的代码分析平台, 用来持续分析和评测项目源代码的质量。 通过SonarQube我们可以检测出项目中重复代码， 潜在bug， 代码风格问题，缺乏单元测试等问题， 并通过一个web ui展示出来。

SonarQube在devops领域很容易与各种ci各种集成。不过这是我们以后会谈到的问题。
支持不少主流语言，但是没有golang。所以引入了下面的话题。

golang代码质量检查分析工具
gometalinter 简介
该工具基本上集成了目前市场上所有的检测工具，然后可以并发的帮你静态分析你的代码

  deadcode
  dupl
  errcheck
  gas
  goconst
  gocyclo
  goimports
  golint
  gosimple
  gotype
  gotypex
  ineffassign
  interfacer
  lll
  maligned
  megacheck
  misspell
  nakedret
  safesql
  staticcheck
  structcheck
  unconvert
  unparam
  unused
  varcheck
  vet
gometalinter安装
go get github.com/alecthomas/gometalinter
gometalinter --install --update
执行上面的两个命令即可。安装非常简单。

gometalinter 的使用
cd 到go项目下，执行 gometalinter ./...
即检查所有目录的go文件，此时vendor目录下的也会检测
如果是想指定指定目录，执行gometalinter + 文件夹名。

bogon:telegraf gaohj$ gometalinter web
web/status.go:165::warning: Errors unhandled.,LOW,HIGH (gas)
web/status.go:165:10:warning: error return value not checked (w.Write([]byte("welcome telegraf for rc"))) (errcheck)
web/status.go:212:19:warning: w can be io.Writer (interfacer)
web/status.go:205:25:warning: do not pass a nil Context, even if a function permits it; pass context.TODO if you are unsure about which Context to use (SA1012) (megacheck)
web/status.go:205:2:warning: 'if err != nil { return err }; return nil' can be simplified to 'return err' (S1013) (megacheck)
vscode集成gometalinter
vscode 默认使用的是golint，如果想用gometalinter替换golint，直接打开
设置项，
在用户设置里添加"go.lintTool": "gometalinter"即可。


总结
其实接下来会做：

gometalinter与SonarQube的配合
SonarQube与CI工具的配合，从而构建devops的整个流程。

用几种：

1，go vet

2，golint

3，staticcheck

4，go/analysis

     addlint filename.go

第 4 种命令如下：

$ addlint: reports integer additions

Usage: addlint [-flag] [package]


Flags:  -V      print version and exit
  -all
        no effect (deprecated)
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -source
        no effect (deprecated)
  -tags string
        no effect (deprecated)
  -trace string
        write trace log to this file
  -v    no effect (deprecated)

更多

5，go/parser

6，go/ast

如果您想使用它，可以在 github.com/fatih/addlint repo 中找到此处编写的所有代码。如果您有更多疑问，请 go/analysis 务必加入Gophers Slack #tools 频道，其中许多 Go 开发人员会讨论问题和问题 go/analysis

GolangCI-Lint是一个lint聚合器，它的速度很快，平均速度是gometalinter的5倍。它易于集成和使用，具有良好的输出并且具有最小数量的误报。而且它还支持go modules。最重要的是免费开源。

下面公司或者产品都使用了golangci-lint，例如：Google、Facebook、Red Hat OpenShift、Yahoo、IBM、Xiaomi、Samsung、Arduino、Eclipse Foundation、WooCart、Percona、Serverless、ScyllaDB、NixOS、The New York Times和Istio。

安装
CI安装
大多数安装都是为CI（continuous integration）准备的，强烈推荐安装固定版本的golangci-lint。

// 二进制文件将会被安装在$(go env GOPATH)/bin/golangci-lint目录
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin vX.Y.Z
// 或者安装它到./bin/目录
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z
// 在alpine Linux中，curl不是自带的，你需要使用下面命令
wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z
上述命令执行完成后，你可以使用golangci-lint --version来查看它的版本。

本地安装
建议不要对CI管道进行本地安装，仅在本地开发环境中以这种方式安装linter。
Windows，MacOS和Linux上，用命令：go get -u github.com/golangci/golangci-lint/cmd/golangci-lint。当你的Go版本不低于1.11时，你可以获取golangci-lint的指定版本。例如：

GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0
在MacOS上面，你还可以使用brew进行安装。

brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
本地安装完成后，你使用golangci-lint --version是看不到它的版本信息的。命令行会提示错误：Error：unknown flag： --version。要想看到可以执行下面的命令。

cd $(go env GOPATH)/src/github.com/golangci/golangci-lint/cmd/golangci-lint
go install -ldflags "-X 'main.version=$(git describe --tags)' -X 'main.commit=$(git rev-parse --short HEAD)' -X 'main.date=$(date)'"
如果你的机器是Windows，你可以在Git Bash中运行上述命令。这时再使用golangci-lint --version就会看到类似的信息：

golangci-lint has version v1.16.0-11-g692dacb built from 692dacb on 2019年04月28日 14:30:25
使用
golangci-lint run [目录]/[文件名]，例如检测下面的go文件。

package main

import (
	"fmt"
	"unsafe"
)

type Book struct {
	Age   int
	Title string
}

func main() {
	b := Book{}
	fmt.Println("book size:", unsafe.Sizeof(b))
}
执行命令golangci-lint run main.go，可以在命令行看到下面提示：

main.go:14:2: SA4006: this value of `b` is never used (staticcheck)
        b := Book{}
        ^
支持的linter
可以通过命令golangci-lint help linters查看它支持的linters。你可以传入参数-E/--enable来使某个linter可用，也可以使用-D/--disable参数来使某个linter不可用。例如：

golangci-lint run --disable-all -E errcheck

goreporter
A Golang tool that does static analysis, unit testing, code review and generate code quality report. This is a tool that concurrently runs a whole bunch of those linters and normalises their output to a report:

Supported linters
Supported template
Installing
Credits
Quickstart
Example
Summary
UnitTest
SimpleCode
DeadCode & CopyCode
Credits
Supported linters
unittest - Golang unit test status.
deadcode - Finds unused code.
gocyclo - Computes the cyclomatic complexity of functions.
varcheck - Find unused global variables and constants.
structcheck - Find unused struct fields.
aligncheck - Warn about un-optimally aligned structures.
errcheck - Check that error return values are used.
copycode(dupl) - Reports potentially duplicated code.
gosimple - Report simplifications in code.
staticcheck - Statically detect bugs, both obvious and subtle ones.
godepgraph - Godepgraph is a program for generating a dependency graph of Go packages.
misspell - Correct commonly misspelled English words... quickly.
Supported template
html template file which can be loaded via -t <file> .
Installing
There are two options for installing goreporter.

Install a stable version, eg. go get -u github.com/wgliang/goreporter/tree/version-1.0.0 . I will generally only tag a new stable version when it has passed the Travis regression tests. The downside is that the binary will be called goreporter.version-1.0.0 .
Install from HEAD with: go get -u github.com/wgliang/goreporter . This has the downside that changes to goreporter may break.
Quickstart
Install goreporter (see above).