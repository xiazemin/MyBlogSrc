---
title: cgo Pkg-Config
layout: post
category: golang
author: 夏泽民
---
通过import "C"语句启动了CGO特性，go build命令会在编译和链接阶段启动gcc编译器

//void SayHello(_GoString_ s); //Go1.10中CGO新增的预定义C语言类型，用来表示Go语言字符串
import "C"

如果在一个go文件中出现了import "C" 指令则表示将调用cgo命令生成的对应的中间文件，
在保证go build 没问题的情况下执行如下命令就可以生成中间文件

go tool cgo main.go
生成的中间文件在_obj目录下

为了在C语言中使用Go语言定义的函数，我们需要将Go代码编译为一个C静态库

go build -buildmode=c-archive -o SayHello.a cgoTest.go

Go生成C动态库

go build -buildmode=c-shared -o SayHello.so cgoTest.go

CGO提供了CFLAGS/CPPFLAGS/CXXFLAGS三种参数，其中CFLAGS对应C语言编译参数(以.c后缀名)、 CPPFLAGS对应C/C++ 代码编译参数(.c,.cc,.cpp,.cxx)、CXXFLAGS对应纯C++编译参数(.cc,.cpp,*.cxx)

链接参数：LDFLAGS
链接参数主要包含要链接库的检索目录和要链接库的名字。因为历史遗留问题，链接库不支持相对路径，我们必须为链接库指定绝对路径。 cgo 中的 ${SRCDIR} 为当前目录的绝对路径。经过编译后的C和C++目标文件格式是一样的，因此LDFLAGS对应C/C++共同的链接参数
<!-- more -->
CGO在使用C/C++资源的时候一般有三种形式：直接使用源码；链接静态库；链接动态库。

直接使用源码就是在import "C"之前的注释部分包含C代码，或者在当前包中包含C/C++源文件。链接静态库和动态库的方式比较类似，都是通过在LDFLAGS选项指定要链接的库方式链接


package main
 
//#cgo CFLAGS: -I./dirname
//#cgo LDFLAGS: -L${SRCDIR}/dirname -lfilename
//
//#include "filename.h"
import "C"
import "fmt"
 
func main() {
     fmt.Println(C.filename_func())
}

http://www.cppcns.com/jiaoben/golang/470638.html

Go 代码只有动态库和程序在同一个文件夹下才能正确工作。这个限制导致无法使用 go get 命令直接下载，编译，并安装这个程序的工作版本。

Go 的工具套件不会把 C 语言的动态库拷贝到 bin 目录下，因为我无法在 go-get 命令完成后，就运行程序。

解决这个问题，需要两个步骤。第一步，我需要使用包配置文件 (package configuration file) 来指定 CGO 的编译器和链接器。第二步，我需要为操作系统设置一个环境变量，让它能在不需要将二进制文件拷贝到 bin 目录下，找到二进制文件。

有些标准库同样也有一个包配置 (.pc) 文件。一个名为 pkg-config 的特殊程序被构建工具（如 gcc）用于从这些文件中检索信息。

 pkg-config 命令：

pkg-config – cflags – libs libcrypto

cd $HOME
export GOPATH=$HOME/keyboard
export PKG_CONFIG_PATH=$GOPATH/src/github.com/goinggo/keyboard/pkgconfig
export DYLD_LIBRARY_PATH=$GOPATH/src/github.com/goinggo/keyboard/DyLib
go get Github.com/goinggo/keyboard

以前
头文件和动态链接库的位置是通过相对路径找到的。

package main

/*
#cgo CFLAGS: -I../DyLib
#cgo LDFLAGS: -L. -lkeyboard
#include <keyboard.h>
*/
import "C"

我告诉 CGO 使用 pkg-config 程序来寻找编译和链接的参数。包配置文件的名字在结尾处被指定。

package main

/*
#cgo pkg-config: – define-variable=prefix=. GoingGoKeyboard
#include <keyboard.h>
*/
import "C"
注意一下，pkg-config 程序使用 -define-variable 参数。这

对我们的包配置文件，运行 pkg-config 程序：

pkg-config – cflags – libs GoingGoKeyboard

-I$GOPATH/src/github.com/goinggo/keyboard/DyLib
-L$GOPATH/src/github.com/goinggo/keyboard/DyLib -lkeyboard

https://zhuanlan.zhihu.com/p/264984763

在Golang中使用cgo调用C库的时候，如果需要引用很多不同的第三方库，那么使用#cgo CFLAGS:和#cgo LDFLAGS:的方式会引入很多行代码。首先这会导致代码很丑陋，最重要的是如果引用的不是标准库，头文件路径和库文件路径写死的话就会很麻烦。一旦第 三方库的安装路径变化了，Golang的代码也要跟着变化，所以使用pkg-config无疑是一种更为优雅的方法，不管库的安装路径有何变化，我们都不 需要修改Go代码，接下来本博主就用一个简单的例子来说明如何在cgo命令中使用pkg-config。

为了保证pkg-config能够找到这个C语言库，我们要为这个库生成一个描述文件，也就是lib/pkgconfig目录下的hello.pc，其内容如下，有不了解该配置文件内容的看客们可以去搜索一下pkg-config的相关文档。
 
# cat hello.pc 
prefix=/home/ubuntu/third-parties/hello
exec_prefix=${prefix}
libdir=${exec_prefix}/lib
includedir=${exec_prefix}/include
Name: hello
Description: The hello library just for testing pkgconfig
Version: 0.1
Libs: -lhello -L${libdir}
Cflags: -I${includedir}
 
       完成pkg-config描述文件的创建后，还需要将该描述文件的路径信息添加到PKG_CONFIG_PATH环境变量中，只有这样 pkg-config才能正确获取这个C语言库的相关信息。此外，我们还需要将该C语言库的库文件路径添加到LD_LIBRARY_PATH环境变量中， 具体命令如下：
 
# export PKG_CONFIG_PATH=/home/ubuntu/third-parties/hello/lib/pkgconfig
# pkg-config --list-all | grep libhello
libhello    libhello - The hello library just for testing pkgconfig
# export LD_LIBRARY_PATH=/home/ubuntu/third-parties/hello/lib
 
       在完成以上一系列准备工作之后，我们就可以开始编写Golang代码了，以下是Golang调用C语言接口的代码示例，我们只需要#cgo pkg-config: libhello和#include < hello_world.h >两行语句即可实现对hello函数的调用。如果C语言库的安装路径发生了变化，只需修改hello.pc这个描述文件即可，Golang代码无需重新修改和编译
       
       http://t.zoukankan.com/mokliu-p-5538926.html