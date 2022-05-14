---
title: WebAssembly Go
layout: post
category: golang
author: 夏泽民
---
已经安装Go 1.11及以上版本。

Getting Started
编辑main.go

package main

import "fmt"

func main() {
    fmt.Println("Hello, Go WebAssembly!")
}
把main.go build成WebAssembly(简写为wasm)二进制文件

GOOS=js GOARCH=wasm go build -o lib.wasm main.go
把JavaScript依赖拷贝到当前路径

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

ls  /usr/local/go/misc/wasm/
go_js_wasm_exec wasm_exec.html  wasm_exec.js

创建一个index.html文件，并引入wasm_exec.js文件，调用刚才build的lib.wasm

<html>
    <head>
        <meta charset="utf-8">
        <script src="wasm_exec.js"></script>
        <script>
            const go = new Go();
            WebAssembly.instantiateStreaming(fetch("lib.wasm"), go.importObject).then((result) => {
                go.run(result.instance);
            });
        </script>
    </head>
    <body></body>
</html>
创建server.go监听8080端口，serve当前路径

package main

import (
  "flag"
  "log"
  "net/http"
)

var (
  listen = flag.String("listen", ":8080", "listen address")
  dir    = flag.String("dir", ".", "directory to serve")
)

func main() {
  flag.Parse()
  log.Printf("listening on %q...", *listen)
  err := http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir)))
  log.Fatalln(err)
}
启动服务

go run server.go
在浏览器访问localhost:8080,打开浏览器console，就可以看到输出Hello, Go WebAssembly!。

reference
https://github.com/golang/go/wiki/WebAssembly
<!-- more -->
编写main.go

package main

import (
  "strconv"
  "syscall/js"
)
// 传入value1, value2, result三个元素的id，将value1+value2结果赋给result元素
func add(ids []js.Value) {
  // 根据id获取输入值
  value1 := js.Global().Get("document").Call("getElementById", ids[0].String()).Get("value").String()
  value2 := js.Global().Get("document").Call("getElementById", ids[1].String()).Get("value").String()
 
  int1, _ := strconv.Atoi(value1)
  int2, _ := strconv.Atoi(value2)
  // 将相加结果set给result元素
  js.Global().Get("document").Call("getElementById", ids[2].String()).Set("value", int1+int2)
}

// 添加监听事件
func registerCallbacks() {
  js.Global().Set("add", js.NewCallback(add))
}

func main() {
  c := make(chan struct{}, 0)
  println("Go WebAssembly Initialized!")
  registerCallbacks()

  <-c
}

将main.go编译成lib.wasm

GOOS=js GOARCH=wasm go build -o lib.wasm main.go
在index.html中调用lib.wasm

<html>
  <head>
    <meta charset="utf-8">
    <script src="wasm_exec.js"></script>
    <script>
      if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer();
          return await WebAssembly.instantiate(source, importObject);
        };
      }

      const go = new Go();
      let mod, inst;
      WebAssembly.instantiateStreaming(fetch("lib.wasm"), go.importObject).then(async (result) => {
        mod = result.module;
        inst = result.instance;
        await go.run(inst)
      });
    </script>
  </head>
  <body>
    <input type="text" id="value1"/>
    <input type="text" id="value2"/>
    <button type="button" id="add" onClick="add('value1', 'value2', 'result');">add</button>
    <input type="text" id="result"/>
  </body>
</html>
打开server，在浏览器打开即可调用WebAssembly二进制文件执行。

go run server.go
示例代码GitHub
https://github.com/wlchn/go-webassembly
reference
https://tutorialedge.net/golang/go-webassembly-tutorial/

https://github.com/golang/go/wiki/WebAssembly
https://godoc.org/syscall/js
https://github.com/siongui/godom/
https://github.com/gopherjs/gopherjs

WebAssembly（也称为wasm）将很快改变这种情况。使用WebAssembly可以用任何语言编写Web应用程序。在本文中，我们将了解如何编写Go程序并使用wasm在浏览器中运行它们。

但首先，什么是WebAssembly
webassembly.org 将其定义为“基于堆栈的虚拟机的二进制指令格式”。这是一个很好的定义，但让我们将其分解为我们可以轻松理解的内容。

从本质上讲，wasm是一种二进制格式; 就像ELF，Mach和PE一样。唯一的区别是它适用于虚拟编译目标，而不是实际的物理机器。为何虚拟？因为不同于 C/C++ 二进制文件，wasm二进制文件不针对特定平台。因此，您可以在Linux，Windows和Mac中使用相同的二进制文件而无需进行任何更改。 因此，我们需要另一个“代理”，它将二进制文件中的wasm指令转换为特定于平台的指令并运行它们。通常，这个“代理”是一个浏览器，但从理论上讲，它也可以是其他任何东西。

这为我们提供了一个通用的编译目标，可以使用我们选择的任何编程语言构建Web应用程序！只要我们编译为wasm格式，我们就不必担心目标平台。就像我们编写一个Web应用程序一样，但是现在我们有了用我们选择的任何语言编写它的优势。

你好 WASM
让我们从一个简单的“hello world”程序开始，但是要确保您的Go版本至少为1.11。我们可以这样写：

package main

import (
    "fmt"
)

func main() {
    fmt.Println("hello wasm")
}
保存为test.go。看起来像是一个普通的Go程序。现在让我们将它编译为wasm平台程序。我们需要设置GOOS和GOARCH。

$GOOS=js GOARCH=wasm go build -o test.wasm test.go
现在我们生成了 wasm 二进制文件。但与原生系统不同，我们需要在浏览器中运行它。为此，还需要再做一点工作来实现这一目标：

Web服务器来运行应用
一个index.html文件，其中包含加载wasm二进制文件所需的一些js代码。
还有一个js文件，它作为浏览器和我们的wasm二进制文件之间的通信接口。

现在Go目录中已经包含了html和js文件，因此我们将其复制过来。

$cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
$cp "$(go env GOROOT)/misc/wasm/wasm_exec.html" .

DOM API
但首先，要使Go代码与浏览器进行交互，我们需要一个DOM API。我们有syscall/js库来帮助我们解决这个问题。它是一个非常简单却功能强大的DOM API形式，我们可以在其上构建我们的应用程序。在我们制作应用程序之前，让我们快速了解它的一些功能。

回调
为了响应DOM事件，我们声明了回调并用这样的事件将它们连接起来：

import "syscall/js"

// Declare callback
cb := js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
    // handle event
})


// Hook it up with a DOM event
js.Global().Get("document").
    Call("getElementById", "myBtn").
    Call("addEventListener", "click", cb)


// Call cb.Release() on your way out.
更新DOM
要从Go中更新DOM，我们可以

import "syscall/js"

js.Global().Get("document").
        Call("getElementById", "myTextBox").
        Set("value", "hello wasm")
您甚至可以调用JS函数并操作本机JS对象，如 FileReader或Canvas。查看syscall/js文档以获取更多详细信息。

正确的 Web 应用程序
接下来我们将构建一个小应用程序，它将获取输入的图像，然后对图像执行一些操作，如亮度，对比度，色调，饱和度，最后将输出图像发送回浏览器。 每个效果都会有滑块，用户可以更改这些效果并实时查看目标图像的变化。

首先，我们需要从浏览器获取输入的图像给到我们的Go代码，以便可以处理它。为了有效地做到这一点，我们需要采取一些不安全的技巧，这里跳过具体细节。拥有图像后，它完全在我们的控制之下，我们可以自由地做任何事情。下面是图像加载器回调的简短片段，为简洁起见略有简化：

onImgLoadCb = js.NewCallback(func(args []js.Value) {
    reader := bytes.NewReader(inBuf) // inBuf is a []uint8 slice where our image is loaded
    sourceImg, _, err := image.Decode(reader)
    if err != nil {
        // handle error
    }
    // Now the sourceImg is an image.Image with which we are free to do anything!
})

js.Global().Set("loadImage", onImgLoadCb)
然后我们从效果滑块中获取用户值，并操纵图像。我们使用了很棒的bild库。下面是回调的一小部分：

import "github.com/anthonynsimon/bild/adjust"

contrastCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
    delta := ev.Get("target").Get("valueAsNumber").Float()
    res := adjust.Contrast(sourceImg, delta)
})

js.Global().Get("document").
        Call("getElementById", "contrast").
        Call("addEventListener", "change", contrastCb)
在此之后，我们将目标图像编码为jpeg并将其发送回浏览器。这是完整的应用程序

由于Go是一种垃圾收集语言，因此整个运行时都在wasm二进制文件中。因此，二进制文件通常有几MB的大小。与C/Rust等其他语言相比，这仍然是一个痛点; 因为向浏览器发送MB级数据并不理想。但是，如果wasm规范本身支持GC，那么这可能会改变。
Go中的Wasm支持正式进行试验。syscall/js API本身也在不断变化，未来可能会发生变化。如果您发现错误，请随时在我们issues报告问题。
与所有技术一样，WebAssembly也不是一颗银弹。有时，简单的JS更快更容易编写。然而，wasm规范本身正在开发中，并且即将推出更多功能。线程支持就是这样一个特性。


环境
需要golang版本高于go1.11, 本文golang版本:

升级go 版本后idea golang 标准库也标红 解决办法，将 /usr/local/go/ 也加入gopath 问题解决

$ go version
> go version go1.11.1 darwin/amd64
js中调用golang函数案例
本案例基于goland IDE编写, 为了获取syscall/js库的自动提示, 需要对IDE进行如下设置:
<img src="{{site.url}}{{site.baseurl}}/img/system_js.jpg"/>

点击File -> Project Structure 打开项目结构配置窗口。 加载go1.13.5

解决方法是 File-Invalidate Caches 然后重启IDEA。


./webassembly.go:54:35: undefined: js.NewCallback


js.Global().Set("jsFunctionName", js.NewCallback(goCallback))
should be as:

js.Global().Set("jsFunctionName", js.FuncOf(goCallback))
Note the signature of goCallback has changed and now since Go 1.12, there is a support for return values. For example, here is how to expose a simple add function:

// function definition
func add(this js.Value, i []js.Value) interface{} {
  return js.ValueOf(i[0].Int()+i[1].Int())
}

// exposing to JS
js.Global().Set("add", js.FuncOf(add))

https://stackoverflow.com/questions/55800163/golangs-syscall-js-js-newcallback-is-undefined

./main.go:22:34: cannot use add (type func([]js.Value)) as type func(js.Value, []js.Value) interface {} in argument to js.FuncOf

//https://godoc.org/syscall/js
//https://github.com/golang/go/wiki/WebAssembly#getting-started

https://github.com/vugu/vugu

在每个浏览器里面，无论Chrome，Firefox，Safari，Edge，能够运行的语言就是Javascript。为了能够让其他语言的代码在浏览器中运行，WebAssembly被创造出来。它拥有更好性能，更小的size，能够更快的加载和执行。我们无需编写WebAssembly的代码，只需要将其他高级语言编译成WebAssembly，这样就能在浏览器中复用大量的其他语言现有的代码。

WebAssembly仍在持续的发展，还有大量的特性即将到来。其最早发明出来是为了将C++的转译成JS，然后在浏览器中运行起来，这样就能把大量现有的C++代码在浏览器中复用。被转译后的JS代码比原生的JS代码要慢，Mozilla的工程师发现一种类型系统，可以让被转译后的JS运行得更快，这就是asm.js. 同时，其他浏览器厂商发现asm.js的运行速度非常快，也把这种优化加入到他们的浏览器引擎中。这仅仅是开始，工程师们仍在持续努力，但是，不是将其他语言编译成JS，而是一种新的语言，那就是WebAssembly。

最小可用产品

WebAssembly不仅仅支持C/C++，同时也希望支持更多的高级语言，因此，需要一个语言无关的编译目标，就像汇编语言一样，支持任何语言编译成汇编语言。这个编译目标有如下的特点：

跟具体的平台无关，因此不同平台的不同浏览器都能运行WebAssembly。
拥有足够快的运行速度，能够带来足够流畅的交互体验。
加载速度要足够快，因此，需要编译目标能够被压缩，减小加载内容的大小
能够手动的管理，分配内存。我们知道C/C++一类的语言支持指针的特性，通过指针可以读写特定地址的内存；为了安全考虑，还要对限制特定地址的内存进行操作。出于以上的亮点，WebAssembly使用了线性内存模型。
通过以上的特点，保证了WebAssembly能够在生产环境中使用起来。

如何应对繁重的桌面应用

我们知道，大量的桌面应用，像PS，AutoCAD，这些应用非常的庞大，对性能要求非常苛刻。要先让他们在浏览器里面运行起来非常的难，因此需要更多的特性来确保更佳的性能：

首当其冲的，是需要支持多线程。现代的计算机都是多核的，通过多线程能够更好的利用计算机的计算能力。
SIMD(单指令多数据)。通过SIMD，能够将一组内存划分成不同的执行单元，就像多核一样。
64位寻址。借助64位寻址，能够使用更多的内存，这对一些内存敏感性的应用是非常有利的。
流式编译。前面提到了，提升加载的速度，其实我们有更好的办法，就是刚下载的时候就开始编译，这将是巨大的提升。
HTTP缓存。如何两个浏览器加载相同的WebAssembly代码，将会编译成相同的机器码，因此可以将编译后的机器码保存在HTTP缓存中，这样就可以跳过编译的过程，复用机器码。
现状
多线程：一个草案已经接近完成，其中的关键SharedArrayBuffers，已经被否决了。
SIMD：正在开发中...
64位寻址：wasm-64即将登场
流式编译：Firefox已经在2017年支持，其他浏览器也即将支持
虽然这些特性仍在开发中，但是我们能够看到已经有大量的桌面应用在浏览器中运行起来，其中最大的幕后功臣就是WebAssembly。

WebAssembly与JavaScript

对于很多的web应用场景，我们可能只需要在一些性能敏感的部分，使用WebAssembly。因此，某些模块需要用WebAssembly来编写，然后替换掉那些JS写的部分。一个例子就是Firefox中的source map library的parser，它用WebAssembly编写，比原来用JS编写的快11倍。为了能让这种场景下，WebAssembly更好的发挥作用，有更多的要求：

JS和WASM能够更快的相互调用。因为要将WASM代码作为模块继承到现存的JS应用中，需要他们能够更快的相互调用，Firefox中已经有了巨大的提升
快速而容易的数据转换。在JS和WASM相互调用时，需要传递数据，要想实现上面的两个目标，非常的难：WASM只理解数字，那就需要将各种数据格式转换成数字
ES module。集成WASM模块，通常在JS中使用import，export关键词，因此，浏览器需要内置ES module。
工具链。在JS中，可以使用npm，brower等工具，但是在WASM中，好像没有这个工具...
兼容性。前端开发，都逃不了兼容性的问题。
现状
Firefox中，JS和WASM能够很快的调用
引用类型草案登场，其增加了一种新的，WASM函数能够接收和返回的类型，这个类型引用一个外部的object，可以是JS的Object。
一个ES module的草案被提及，浏览器厂商正在支持。
Rust生态的wasm-pack能够像npm一样支持包管理
借助wasm2js工具，能够让WASM在旧版的浏览器中得到支持
通过以上的特性以及正在开发中的功能，WASM的能力得到释放，接下来就是如何再现有的Web生态中使用WASM。

应用

在前端开发中，大量涉及的框架及编译成JS的语言都将是WASM发挥作用的场景。所以就有两种选择了：1，使用WASM来重写现有的Web框架；2，将Reasonml，Elm等语言编译成WASM。为了实现这些功能，需要WASM提供更多高级语言的特性，包括：

GC。首先，提供GC功能对重写web框架是非常有优势的。例如：使用WASM重写React中的diff功能，借助多线程，手动的内存分配，能够提供以前无法现象的高性能，但是当你跟JS 对象交互时，例如组件，仍然需要GC来减轻开发的负担。
异常处理。很多的高级语言，如C/C++提供异常处理，在某些特定场景下非常有用，同时JS也有异常处理，当WASM和JS互操作时，也需要有异常处理的支持。
debug。这个就不多说
现状
JS拥有Typed Objects 草案，WASM拥有GC草案。通过这两个草案，JS和WASM都能够清晰的知道一个对象的结构以及如何去存储，使用，回收。
异常处理。目前还在开发阶段。
debug。目前，大多数浏览器已经支持。

https://hacks.mozilla.org/2018/10/webassemblys-post-mvp-future/

https://cloud.tencent.com/developer/article/1369556

able Of Contents
Introduction
Starting Point
Function Registration
Components
Building a Router
A Full Example
Challenges Going Forward
Conclusion
JavaScript Frontend frameworks have undoubtedly helped to push the boundaries of what was previously possible in the context of a browser. Ever more complex applications have come out built on top of the likes of React, Angular and VueJS to name but a few and there’s the well known joke about how a new frontend framework seems to come out every day.

However, this pace of development is exceptionally good news for developers around the world. With each new framework, we discover better ways of handling state, or rendering efficiently with things like the shadow DOM.

The latest trend however, seems to be moving towards writing these frameworks in languages other than JavaScript and compiling them into WebAssembly. We’re starting to see major improvements in the way that JavaScript and WebAssembly communicates thanks to the likes of Lin Clark and we’ll undoubtedly see more major improvements as WebAssembly starts to become more prominent in our lives.

Introduction
So, in this tutorial, I thought it would be a good idea to build the base of an incredibly simple frontend framework written in Go that compiles into WebAssembly. At a minimum, this will include the following features:

Function Registration
Components
Super Simplistic-Routing
I’m warning you now though that these are going to be incredibly simple and nowhere near production ready. If this is article is somewhat popular, I’ll hopefully be taking it forward however, and trying to build something that meets the requirements of a semi-decent frontend framework.

Github: The full source code of this project can be found here: elliotforbes/oak. If you fancy contributing to the project, feel free, I’d be happy to get any pull requests!

Starting Point
Right, let’s dive into our editor of choice and start coding! The first thing we’ll want to do is create a really simple index.html that will act as our entry point for our frontend framework:

 1<!doctype html>
 2<!--
 3Copyright 2018 The Go Authors. All rights reserved.
 4Use of this source code is governed by a BSD-style
 5license that can be found in the LICENSE file.
 6-->
 7<html>
 8
 9<head>
10    <meta charset="utf-8">
11    <title>Go wasm</title>
12    <script src="./static/wasm_exec.js"></script>
13    <script src="./static/entrypoint.js"></script>
14</head>
15<body>    
16
17  <div class="container">
18    <h2>Oak WebAssembly Framework</h2>
19  </div>
20</body>
21
22</html>
You’ll notice these have 2 js files being imported at the top, these allow us to execute our finished WebAssembly binary. The first of which is about 414 lines long so, in the interest of keeping this tutorial readable, I recommend you download it from here: https://github.com/elliotforbes/oak/blob/master/examples/blog/static/wasm_exec.js

The second is our entrypoint.js file. This will fetch and run the lib.wasm that we’ll be building very shortly.

1// static/entrypoint.js
2const go = new Go();
3WebAssembly.instantiateStreaming(fetch("lib.wasm"), go.importObject).then((result) => {
4    go.run(result.instance);
5});
Finally, now that we have that out of the way, we can start diving into some Go code! Create a new file called main.go which will contain the entry point to our Oak Web Framework!

1// main.go
2package main
3
4func main() {
5    println("Oak Framework Initialized")
6}
This is as simple as it gets. We’ve created a really simple Go program that should just print out Oak Framework Initialized when we open up our web app. To verify that everything works, we need to compile this using the following command:

1$ GOOS=js GOARCH=wasm go build -o lib.wasm main.go
This should then build our Go code and output our lib.wasm file which we referenced in our entrypoint.js file.

Awesome, if everything worked, then we are ready to try it out in the browser! We can use a really simple file server like this:

 1// server.go
 2package main
 3
 4import (
 5    "flag"
 6    "log"
 7    "net/http"
 8)
 9
10var (
11    listen = flag.String("listen", ":8080", "listen address")
12    dir    = flag.String("dir", ".", "directory to serve")
13)
14
15func main() {
16    flag.Parse()
17    log.Printf("listening on %q...", *listen)
18    log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
19}
You can then serve your application by typing go run server.go and you should be able to access your app from http://localhost:8080.

Function Registration
Ok, so we’ve got a fairly basic print statement working, but in the grand scheme of things, I don’t quite think that qualifies it as a Web Framework just yet.

Let’s take a look at how we can build functions in Go and register these so we can call them in our index.html. We’ll create a new utility function which will take in both a string which will be the name of our function as well as the Go function it will map to.

Add the following to your existing main.go file:

1// main.go
2import "syscall/js"
3
4// RegisterFunction
5func RegisterFunction(funcName string, myfunc func(i []js.Value)) {
6    js.Global().Set(funcName, js.NewCallback(myfunc))
7}
So, this is where things start to become a bit more useful. Our framework now allows us to register functions so users of the framework can start creating their own functionality.

Other projects using our framework can start to register their own functions that can subsequently be used within their own frontend applications.

Components
So, I guess the next thing we need to consider adding to our framework is the concept of components. Basically, I want to be able to define a components/ directory within a project that uses this, and within that directory I want to be able to build like a home.gocomponent that features all the code needed for my homepage.

So, how do we go about doing this?

Well, React tends to feature classes that feature render() functions which return the HTML/JSX/whatever code you wish to render for said component. Let’s steal this and use it within our own components.

I essentially want to be able to do something like this within a project that uses this framework:

1package components
2
3type HomeComponent struct{}
4
5var Home HomeComponent
6
7func (h HomeComponent) Render() string {
8    return "<h2>Home Component</h2>"
9}
So, within my components package, I define a HomeComponent which features a Render() method which returns our HTML.

In order to add components to our framework, we’ll keep it simple and just define an interface to which any components we subsequently define will have to adhere to. Create a new file called components/comopnent.go within our Oak framework:

1// components/component.go
2package component
3
4type Component interface {
5    Render() string
6}
What happens if we want to add new functions to our various components? Well, this allows us to do just that. We can use the oak.RegisterFunction call within the initfunction of our component to register any functions we want to use within our component!

 1package components
 2
 3import (
 4    "syscall/js"
 5
 6    "github.com/elliotforbes/oak"
 7)
 8
 9type AboutComponent struct{}
10
11var About AboutComponent
12
13func init() {
14    oak.RegisterFunction("coolFunc", CoolFunc)
15}
16
17func CoolFunc(i []js.Value) {
18    println("does stuff")
19}
20
21func (a AboutComponent) Render() string {
22    return `<div>
23                        <h2>About Component Actually Works</h2>
24                        <button onClick="coolFunc();">Cool Func</button>
25                    </div>`
26}
When we combine this with a router, we should be able to see our HTML being rendered to our page and we should be able to click that button which calls coolFunc() and it will print out does stuff within our browser console!

Awesome, let’s see how we can go about building a simple router now.

Building a Router
Ok, so we’ve got the concept of components within our web framework down. We’ve almost finished right?

Not quite, the next thing we’ll likely need is a means to navigate between different components. Most frameworks seem to have a <div> with a particular id that they bind to and render all their components within, so we’ll steal that same tactic within Oak.

Let’s create a router/router.go file within our oak framework so that we can start hacking away.

Within this, we’ll want to map string paths to components, we wont do any URL checking, we’ll just keep everything in memory for now to keep things simple:

 1// router/router.go
 2package router
 3
 4import (
 5    "syscall/js"
 6
 7    "github.com/elliotforbes/oak/component"
 8)
 9
10type Router struct {
11    Routes map[string]component.Component
12}
13
14var router Router
15
16func init() {
17    router.Routes = make(map[string]component.Component)
18}
So within this, we’ve created a new Router struct which contains Routes which are a map of strings to the components we defined in the previous section.

Routing won’t be a mandatory concept within our framework, we’ll want users to choose when they wish to initialize a new router. So let’s create a new function that will register a Link function and also bind the first route in our map to our <div id="view"/> html tag:

 1// router/router.go
 2// ...
 3func NewRouter() {
 4    js.Global().Set("Link", js.NewCallback(Link))
 5    js.Global().Get("document").Call("getElementById", "view").Set("innerHTML", "")
 6}
 7
 8func RegisterRoute(path string, component component.Component) {
 9    router.Routes[path] = component
10}
11
12func Link(i []js.Value) {
13    println("Link Hit")
14
15    comp := router.Routes[i[0].String()]
16    html := comp.Render()
17
18    js.Global().Get("document").Call("getElementById", "view").Set("innerHTML", html)
19}
You should notice, we’ve created a RegisterRoute function which allows us to register a path to a given component.

Our Link function is also pretty cool in the sense that it will allow us to navigate between various components within a project. We can specify really simple <button>elements to allow us to navigate to registered paths like so:

1<button onClick="Link('link')">Clicking this will render our mapped Link component</button>
Awesome, so we’ve got a really simple router up and running now, if we wanted to use this in a simple application we could do so like this:

 1// my-project/main.go
 2package main
 3
 4import (
 5    "github.com/elliotforbes/oak"
 6    "github.com/elliotforbes/oak/examples/blog/components"
 7    "github.com/elliotforbes/oak/router"
 8)
 9
10func main() {
11    // Starts the Oak framework
12    oak.Start()
13
14    // Starts our Router
15    router.NewRouter()
16    router.RegisterRoute("home", components.Home)
17    router.RegisterRoute("about", components.About)
18
19    // keeps our app running
20    done := make(chan struct{}, 0)
21    <-done
22}
A Full Example
With all of this put together, we can start building really simple web applications that feature components and routing. If you want to see a couple of examples as to how this works, then check out the examples within the official repo: elliotforbes/oak/examples

Challenges Going Forward
The code in this framework is in no way production ready, but I’m hoping this post kicks off good discussion as to how we can start building more production ready frameworks in Go.

If nothing else, it starts the journey of identifying what still has to be done to make this a viable alternative to the likes of React/Angular/VueJS, all of which are phenomenal frameworks that massively speed up developer productivity.

I’m hoping this article motivates some of you to go off and start looking at how you can improve on this incredibly simple starting point.

Conclusion
If you enjoyed this tutorial, then please feel free to share it to your friends, on your twitter, or wherever you feel like, it really helps the site and directly supports me writing more!

I’m also on YouTube, so feel free to subscribe to my channel for more Go content! - TutorialEdge.

The full source code for the Oak framework can be found here: github.com/elliotforbes/oak. Feel free to submit PRs!


https://www.vugu.org/doc/start
https://github.com/vugu/vugu