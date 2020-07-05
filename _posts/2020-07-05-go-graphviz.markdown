---
title: go-graphviz
layout: post
category: golang
author: 夏泽民
---
产生Graphviz 文件
用 Golang 产生Graphviz 文件的封装方法如下
https://github.com/awalterschulze/gographviz
我们使用它的范例代码如下：

package main

import (
    "fmt"
    "github.com/awalterschulze/gographviz"
)

func main() {
    graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
    graph := gographviz.NewGraph()
    gographviz.Analyse(graphAst, graph)
    graph.AddNode("G", "a", nil)
    graph.AddNode("G", "b", nil)
    graph.AddEdge("a", "b", true, nil)
    fmt.Println(graph.String())
}
使用 dot 生成 png 的命令行如下：

dot 11.gv -T png -o 11.png
<!-- more -->
{% raw %}
完整的 Go 代码如下：

 

package main

import (
    "bytes"
    "fmt"
    "github.com/awalterschulze/gographviz"
    "io/ioutil"
    "os/exec"
)

func main() {
    graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
    graph := gographviz.NewGraph()
    gographviz.Analyse(graphAst, graph)
    graph.AddNode("G", "a", nil)
    graph.AddNode("G", "b", nil)
    graph.AddEdge("a", "b", true, nil)
    fmt.Println(graph.String())

    // 输出文件
    ioutil.WriteFile("11.gv", []byte(graph.String()), 0666)

    // 产生图片
    system("dot 11.gv -T png -o 12.png")
}

//调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) {
    cmd := exec.Command(`/bin/sh`, `-c`, s) //调用Command函数
    var out bytes.Buffer                    //缓冲字节

    cmd.Stdout = &out //标准输出
    err := cmd.Run()  //运行指令 ，做判断
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%s", out.String()) //输出执行结果
}
{% endraw %}

goimportdot : 一个帮你迅速了解 golang 项目结构的工具

https://github.com/yqylovy/goimportdot

go get -u github.com/yqylovy/goimportdot
goimportdot -pkg=yourpackagename > pkg.dot 
dot -Tsvg pkg.dot >pkg.svg


画包的依赖关系图
goimportdot -pkg=github.com/yqylovy/goimportdot >goimportdot.dot
dot -Tsvg goimportdot.dot >goimportdot.svg
goimportdot -pkg=github.com/beego >beego.dot
dot -Tsvg beego.dot >beego.svg

https://github.com/ilikeorangutans/grails-service-visualizer
https://ilikeorangutans.github.io/2014/05/03/using-golang-and-graphviz-to-visualize-complex-grails-applications/


https://github.com/skratchdot/open-golang


go-callvis 是github上一个开源项目，可以用来查看golang代码调用关系。
安装go-callvis
go get -u github.com/TrueFurby/go-callvis
cd $GOPATH/src/github.com/TrueFurby/go-callvis && make

查看github.com/github/orchestrator/go/http 这个package下面的调用关系：

$ go-callvis -focus github.com/github/orchestrator/go/http  github.com/github/orchestrator/go/cmd/orchestrator


https://www.cnblogs.com/lanyangsh/p/10011093.html

https://blog.csdn.net/qq_34857250/article/details/100643339

https://blog.csdn.net/lanyang123456/article/details/84425565
http://www.mamicode.com/info-detail-2529639.html

https://www.bbsmax.com/A/8Bz8Gmmydx/

http://www.360doc.com/content/18/0328/16/51898798_741004626.shtml



