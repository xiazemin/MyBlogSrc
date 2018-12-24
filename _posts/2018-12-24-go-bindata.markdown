---
title: go-bindata Go 语言打包静态文件
layout: post
category: golang
author: 夏泽民
---
https://github.com/jteeuwen/go-bindata
go get -u github.com/jteeuwen/go-bindata/...
执行后仍然没有找到go-bindata命令

echo $GOPATH 
/root/devhome/golang/go_demo
cd /root/devhome/golang/go_demo/bin

对于 Go 语言开发者来说，在享受语言便利性的同时，最终编译的单一可执行文件也是我们所热衷的。但是，一旦遇到我们需要分发的东西不只有可执行文件的时候，事情就变得稍微有点复杂了，例如，需要分发个默认的配置文件；或者说是一个 Web 服务需要附带一些简单的 js/css 文件之类的。

当然，对于经验丰富的老司机们来说这都不是问题，例如 RH 系列的 RPM 是很多老司机们的选择，像我这样的新手也是觉得老司机们的这车开得好，可以很方便得管理一个分发包。但是，对于我们说的如果只有一点点文件，我就来打个 rpm 包，似乎有点脱裤放屁的感觉；同时，有时我们可能是在 OSX/Win 平台跑，这又缺乏一定的通用性。

基于很多考虑，有同学就觉得为啥不把这些简单的静态文件打包进可执行文件中，这样我还是只分发一个文件就可以了，其实我也挺喜欢这种方式的，好处就是很简单，带来的缺点当然也很明显，自然分发可执行文件就大了。Anyway，有得有失，不妨一起和我看一下这种方式是如何实现的。
<!-- more -->
在 Go 语言的 Awesome 中你可以看到很多静态打包库，但是，你却看不到 go-bindata，哈哈，这很奇怪哈。而且相比于 Awesome 中列举的这些库，go-bindata 明显更受欢迎，更流行。

go-bindata 很简单，设计理念也不难理解。它的任务就是讲静态文件封装在一个 Go 语言的 Source Code 里面，然后提供一个统一的接口，你通过这个接口传入文件路径，它将给你返回对应路径的文件数据。这也就是说它不在乎你的文件是字符型的文件还是字节型的，你自己处理，它只管包装。
To access asset data, we use the Asset(string) ([]byte, error) function which is included in the generated output.

data, err := Asset("pub/style/foo.css")
if err != nil {
	// Asset was not found.
}

// use asset data
就是一个使用 go-bindata 的例子，这里讲 etc/config.json 通过 go-bindata 封装起来，然后下面就直接使用它，这样在运行的时候我就不用关注静态文件的具体位置了

然后再 etc 目录下是放置我的配置文件：config.json，这也是我希望打包的配置文件，然后再 file/service.go 里面我是实现了需要进行文件读取的代码。
1. 使用 go-bindata 打包静态文件
前面说了，go-bindata 是将静态文件打包成 go 文件，所以第一步就是使用 go-bindata 读取配置文件，然后再生成 go 文件，具体的使用命令为：
$ go-bindata  -pkg etc -o etc/bindata.go etc/
关于 go-bindata 的参数，可以使用 go-bindata --help 查看。使用上面这条命令之后，会发现在 etc 目录下多了一个 bindata.go 文件，下一步我们就会使用它来代替 etc/config.json 文件。
使用静态文件
使用静态文件的代码也是很简单，我们只需要记住我们刚才生成 go 文件的相对路径，然后使用：
Asset("etc/config.json")
这样的语句就可以读取到静态文件了

https://tech.townsourced.com/post/embedding-static-files-in-go/
http://www.wiremoons.com/posts/2014-11-30-Serving-Up-Static-Content-from-Golang-Apps/
