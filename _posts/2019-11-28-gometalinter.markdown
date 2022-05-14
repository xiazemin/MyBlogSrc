---
title: gometalinter
layout: post
category: golang
author: 夏泽民
---
SonarQube 是一个开源的代码分析平台, 用来持续分析和评测项目源代码的质量。 通过SonarQube我们可以检测出项目中重复代码， 潜在bug， 代码风格问题，缺乏单元测试等问题， 并通过一个web ui展示出来。
<!-- more -->
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

https://github.com/uartois/sonar-golang

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

Run it:

$ goreporter -p [projtectRelativelyPath] -d [reportPath] -e [exceptPackagesName] -r [json/html]  {-t templatePathIfHtml}
Example
$ goreporter -p ../goreporter -d ../goreporter -t ./templates/template.html

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

https://github.com/golangci/golangci-lint/


1.测试驱动开发（TDD）
如果说要找一个最能提高代码质量同时还要减少bug的实践练习恐怕就非TDD莫属了。它的优点是适用于任何类型的项目和敏捷开发。其历史可以追溯到很早以前，但是直到XP的普及它才渐渐为人所知。当作为能自动化构建和测试实践的持续集成周期的一部分运作的时候，它被称为单元测试。

很多开发人员并不知道该怎么提高这方面的能力，这需要培训和教育。而且这是一个学习和积累的过程，不要想着能一夜吃成个胖子。

2.验收测试驱动开发（ATDD）
这是基于TDD单元测试之后的一个新的水平。这不但表明了验收标准，而且还能在开发工作开始之前自动执行开发需求。在很多情况下，需要专业测试人员和客户携手共同参与到测试中去。

3.持续集成（CI）
这能确保新代码不会干扰到已经存在的代码。如果再加上TDD和ATDD一起创建一个自动化、可重复的的测试套件，将会大幅度提高其使用价值。

4.结对编程
有关于结对编程的争论似乎已经偃旗息鼓了，同样的人们实际应用的例子也越来越少。这不可谓不是一个遗憾。因为在即时的代码审查上，两个脑袋总比一个管用。它也允许开发人员将注意力全部灌注到手头的工作上——不必分心于电话、邮件、短信等等，因为我们的partner会搞定。

5.代码审查
如果没办法结对编程，那么退而求其次，至少得进行一次代码审查。最好代码一写好就能落实到位一个轻量级流程的代码审查。我们在学校里学的那种又大又正规的流程其实并不实际——只有NASA（ 美国宇航局）这种不差钱的土豪才买得起。所以换个轻量级的流程，意味着只需20%的成本就能享受80%的相同效果。

6.静态分析工具
以前人人都不看好所谓的静态分析工具。现在则好了很多，虽然它们仍然并不能真正替代代码审查，但是其使用成本比较低。当然可能需要购买许可证，但是一旦将它们设置进系统中之后，以后每一次我们输入代码，它们都会一丝不苟兢兢业业地检查并且快速提示发现的所有错误。

7.编码标准
老实说我并不怎么喜欢编码标准。从我的经验来看，很多团队在讨论编码标准上面浪费了太多的时间，而且一旦确定了某种标准，这往往会损害一部分开发人员的利益。不过如果我们能克服这些问题，那么绝对会有意想不到的效果。

首先建立一个讨论小组——应该以一种面对面的形式，不要通过电子邮件和电话——讨论出编码标准里应该包含哪些内容。找到需要讨论的地方，规分为不同的类别：少许定位为必选项目，推荐项目的数量可以较前者多点，候选项目则可以更多。在候选组里的需要经过深思熟虑之后才能放到推荐组和必选组中。剩下的第四组则是明确不能成为编程标准的内容。


 
每隔三至四个月检查一下这些标准，看看有没有需要从候选组提升到推荐组，或者从推荐组放到必选组的，要是发现什么已经不适应当前工作的项目，那就尽快删除或者降级。

此外，我们不应该将编码标准当做代码审查的一部分，而是两手都要抓，两手都要硬，万一不得不遗漏其中之一，可以借助自动化工具，例如运行静态分析工具，自动执行代码标准来检查代码。

8.自动化
其实就目前而言，我们提出的大多数意见和建议，是能够自动化执行而且也应该被自动化执行的，但是可惜的是这个概念还没有深入人心。从长远看，非自动化就意味着需要耗费大量的时间，而且成本更高。虽然自动化看似在短期内需要投入大量的成本，但是从整体上而言，其实是节约了成本的。

9.重构（以及重构工具）
重构的目的就在于提高代码质量，当然更重要的是，改善整体的设计。如果重构之后不能达到上述目的，那么说明你的思路错了。我们可以在重构的时候摒弃自动化单元测试，而且很多人也是这么做的，但是这等同于高空走钢丝的时候下面没有安全防护网——一旦失足便万劫不复。如果是装备了“安全防护网”的重构不但毋须占用大量时间，而且还能频繁运行。

以上这些对于能提高代码质量显然是显而易见的。还有一些虽然也在图片名单上，但是并不那么为大家所认可，不过我认为它们也值得包括进去。

10.展示和说明（早期）
也许你会奇怪这怎么能提高代码质量呢，请不要怀疑，It does。因为定期展示相关潜在客户对于软件的要求，能促使开发人员不断地将他们的代码保持在最接近发布的状态，这也使得开发进程更快、更细致。

第二个原因则是能收集更多周期性的反馈，指引我们正确的方向。

最后，如果一个开发人员害怕将他的工作展现给用户和客户看，这是一个非常危险的信号，最好停下来好好自我检查一番。

11.用户测试
用户测试能让我们从另一角度进行测试，以便尽早发现问题。

与第10点相同，碎片化的处理模式能提供更为细致的步骤。无论是在工作规划上还是在改进代码上面，这都给了我们一个机会，能在做每一个决定之前都可以重新调整和矫正航向。

12.团队凝聚力
关于团队的凝聚力其重要性不言而喻，因为一个团队一旦失去了凝聚力，那么大家就会各执己见，各施其力。要想不如此，我们就必须要在开发目标和如何设计代码以及如何改进代码上面的观点达成一致。


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
html template file which can be loaded via-t <file>.
Installing
There are two options for installing goreporter.

Install a stable version, eg.go get -u github.com/wgliang/goreporter/tree/version-1.0.0. I will generally only tag a new stable version when it has passed the Travis regression tests. The downside is that the binary will be calledgoreporter.version-1.0.0.
Install from HEAD with:go get -u github.com/wgliang/goreporter. This has the downside that changes to goreporter may break.
Quickstart
Install goreporter (see above).

Run it:

$ goreporter -p [projtectRelativelyPath] -d [reportPath] -e [exceptPackagesName] -r [json/html]  {-t templatePathIfHtml}
Example
$ goreporter -p ../goreporter -d ../goreporter -t ./templates/template.html
Summary
summary

UnitTest
unittest

SimpleCode
simplecode

DeadCode & CopyCode
deadcodeandcopycode

Credits
Templates is designed by liufeifei

Logo is desigend by xuri


Golang 代码的静态分析尤其简单，为什么？Golang 语法方面，我们在调用一个 package 的时候，需要 import。同时，如果我们代码中对一个 package 没有引用而又把包 import 进来，这个时候 Golang 的编译器是会报错的。这两点意味着我们可以通过分析 import 语句块来得到我们包内部的调用关系。

当然尽管上面原理说的比较简单，但是实现起来还是有一些坑。没关系，我帮你们实现了。链接：github.com/legendtkl/g… 。取名 DAG 也是因为包之间的依赖其实就是一个 DAG。

1. 安装
go get -u github.com/legendtkl/godag
我在实现的时候包里面没有引用第三方包，这样也是为了不能翻墙的同学 go get 无障碍。

2. 使用
以 beego 项目为例，分析其内部的包之间的调用关系。

godag --pkg_name=github.com/astaxie/beego --pkg_path=/Users/kltao/code/go/src/github.com/astaxie/beego --depth=1 --dot_file_path=a.dot
godag 支持四个参数：

pkg_name: 要分析的 package 名称，必填
pkg_path: package 存放在本地的目录，必填
depth: 分析的代码深度。举个例子，如果 depth 为 1，我们则会分析包：beego/cache，beego/context；而如果 depth 为 2，则分析这些包：beego/cache/redis，beego/cache/ssdb 等。
dotfilepath: 是我们输出的 dot 文件，下面细说
3. 输出
3.1 dot
godag 会输出一个 .dot 文件。dot 是一种绘图语言，它可以方便你采用图形的方式快速、直观地表达一些想法，比如描述某个问题的解决方案，构思一个程序的流程，澄清一堆貌似散乱无章的事物之间的联系。举个列子，下面列出 dot 文件以及对应的流程图。

digraph graphname {
     a -> b -> c;
     b -> d;
 }
对应流程图。


关于 dot 文件更详细的信息可以参考：DOT(graph description language) 。我们生产的 dot 文件内容如下。

digraph G {
    "beego/utils" -> "beego/session"
    "beego/logs" -> "beego/logs"
    "beego/utils" -> "beego"
}
3.2 可视化
dot 文件的可视化可以使用 graphviz。graphviz 的安装比较简单，比如 Mac 上安装命令如下

brew install graphviz
dot 文件可视化，使用 dot 命令，比如生成 png 图。

dot -Tpng godag.dot > godag.png


https://github.com/legendtkl/godag

https://github.com/antonmedv/watch

