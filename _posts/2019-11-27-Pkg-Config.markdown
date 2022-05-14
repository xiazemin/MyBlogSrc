---
title: Pkg-Config cgo
layout: post
category: golang
author: 夏泽民
---
https://github.com/xlab/c-for-go

https://dave.cheney.net/2016/01/18/cgo-is-not-go


https://www.ardanlabs.com/blog/2013/08/using-c-dynamic-libraries-in-go-programs.html

https://www.ardanlabs.com/blog/2013/08/using-cgo-with-pkg-config-and-custom.html

构建了一个 C 语言动态库，并编写了一个使用该动态库的 Go 程序。但其中的 Go 代码只有动态库和程序在同一个文件夹下才能正确工作。

这个限制导致无法使用 go get 命令直接下载，编译，并安装这个程序的工作版本。我不想代码需要预先安装依赖或者在调用 go-get 之后执行任何脚本和命令才能正确运行。Go 的工具套件不会把 C 语言的动态库拷贝到 bin 目录下，因为我无法在 go-get 命令完成后，就运行程序。这简直是不可接受的，必须有办法让我能够在运行完 Go-get 之后，就获得一个正确运行程序。

解决这个问题，需要两个步骤。第一步，我需要使用包配置文件 (package configuration file) 来指定 CGO 的编译器和链接器。第二步，我需要为操作系统设置一个环境变量，让它能在不需要将二进制文件拷贝到 bin 目录下，找到二进制文件。

如果你找找看，你会发现有些标准库同样也有一个包配置 (.pc) 文件。一个名为 pkg-config 的特殊程序被构建工具（如 gcc）用于从这些文件中检索信息。
<!-- more -->
https://www.freedesktop.org/wiki/Software/pkg-config/
文件头部的 prefix 的变量是最重要的。这个变量指定库和头文件被安装的基础目录 (base folder)。

另外一个需要注意的事情是，你不能使用环境变量来帮助指定一条路径位置。如果你这么做，构建工具在定位它所需要的任何文件都会有类似的问题 (you will have problems with the build tools locating any of the files it needs.)。 这个环境变量最终会一个字符串的形式提供给编译工具。请记住这一点，因为它很重要。

以下参数在终端运行这个 pkg-config 命令：

pkg-config – cflags – libs libcrypto
这些参数要求 pkg-config 程序显示 libcrypto 这个 .pc 类型文件所设定的编译器和链接器参数。

这是应该返回的：

-lcrypto -lz
pkg-config – cflags – libs MagickWand

-fopenmp -DMAGICKCORE_HDRI_ENABLE=0 -DMAGICKCORE_QUANTUM_DEPTH=16
-I/usr/local/include/ImageMagick-6  -L/usr/local/lib -lMagickWand-6.Q16
-lMagickCore-6.Q16

我直接指定了编译器和链接器的参数。头文件和动态链接库的位置是通过相对路径找到的。

package main

/*
#cgo CFLAGS: -I../DyLib
#cgo LDFLAGS: -L. -lkeyboard
#include <keyboard.h>
*/
import "C"
这是修改后的代码。在这个代码中，我告诉 CGO 使用 pkg-config 程序来寻找编译和链接的参数。包配置文件的名字在结尾处被指定。

package main

/*
#cgo pkg-config: – define-variable=prefix=. GoingGoKeyboard
#include <keyboard.h>
*/
import "C"
注意一下，pkg-config 程序使用 -define-variable 参数。这个设置是让一切运转的诀窍。让我们马上回过头来看看。

对我们的包配置文件，运行 pkg-config 程序：

pkg-config – cflags – libs GoingGoKeyboard

-I$GOPATH/src/github.com/goinggo/keyboard/DyLib
-L$GOPATH/src/github.com/goinggo/keyboard/DyLib -lkeyboard
如果仔细观察调用的输出，你会看到些我告诉你的错误的用法。$GOPATH 环境变量是运行时提供的。

打开在 pkgconfig 目录下的包配置文件，你会看到 pkg-config 程序没有撒谎。在文件的头部，我正在使用 $GOPATH 设置一条路径的前缀路径 (prefix variable)。 那为什么一切都有效？

我直接指定了编译器和链接器的参数。头文件和动态链接库的位置是通过相对路径找到的。

package main

/*
#cgo CFLAGS: -I../DyLib
#cgo LDFLAGS: -L. -lkeyboard
#include <keyboard.h>
*/
import "C"
这是修改后的代码。在这个代码中，我告诉 CGO 使用 pkg-config 程序来寻找编译和链接的参数。包配置文件的名字在结尾处被指定。

package main

/*
#cgo pkg-config: – define-variable=prefix=. GoingGoKeyboard
#include <keyboard.h>
*/
import "C"
注意一下，pkg-config 程序使用 -define-variable 参数。这个设置是让一切运转的诀窍。让我们马上回过头来看看。

对我们的包配置文件，运行 pkg-config 程序：

pkg-config – cflags – libs GoingGoKeyboard

-I$GOPATH/src/github.com/goinggo/keyboard/DyLib
-L$GOPATH/src/github.com/goinggo/keyboard/DyLib -lkeyboard
如果仔细观察调用的输出，你会看到些我告诉你的错误的用法。$GOPATH 环境变量是运行时提供的。

打开在 pkgconfig 目录下的包配置文件，你会看到 pkg-config 程序没有撒谎。在文件的头部，我正在使用 $GOPATH 设置一条路径的前缀路径 (prefix variable)。 那为什么一切都有效？

6

让我们使用在 main.go 代码中相同的选项运行这个程序：

pkg-config – cflags – libs GoingGoKeyboard – define-variable=prefix=.

-I./DyLib
-L./DyLib -lkeyboard
你看到有什么不同吗？在第一次运行 pkg-config 程序时，我们获得的路径中使用 $GOPAHT 这样一个字符串的，因为这就是前缀变量的设置方式。第二次运行时，我们将前缀变量的值覆盖到当前目录， 得到我们想要的返回。

还记得我们在使用 Go 工具之前设置的环境变量吗？

PKG_CONFIG_PATH=$GOPATH/src/github.com/goinggo/keyboard/pkgconfig
PKG_CONFIG_PATH 环境变量告诉 pkg-config 程序，它可以在哪里找到不在任何默认位置的软件包配置文件。我们的 GoingGoKeyboard.pc 文件就是这样被 pkg-config 程序找到的。

最后一个要解释的谜团是，操作系统如何找到运行我们程序所需要的动态库。还记得我们在使用 Go 工具之前设置的这个环境变量吗？

export DYLD_LIBRARY_PATH=$GOPATH/src/github.com/goinggo/keyboard/DyLib
DYLD_LIBRARY_PATH 环境变量告诉操作系统在哪里还可以查找动态库。

在 /usr/local 文件夹中安装动态库可以使事情保持简单。默认情况下，所有构建工具都配置为在这个文件夹中查找。但是，如果对自己的或第三方库文件使用默认位置，需要在运行 Go 工具之前执行额外的安装步骤。通过使用包配置文件，向 pkg-config 程序传递所需的选项，使用 CGO 的 Go 程序可以部署安装即可运行的构建。

还有一个我们提到的好处，你可以使用这种技术来将第三方库安装到一个临时的路径下进行测试使用。这让你在不想使用这个第三库时，可以很方便地进行移除。

cgo is not Go
To steal a quote from JWZ,

Some people, when confronted with a problem, think “I know, I’ll use cgo.”
Now they have two problems.

Recently the use of cgo came up on the Gophers’ slack channel and I voiced my concerns that using cgo, especially on a project that is intended to showcase Go inside an organisation was a bad idea. I’ve said this a number of times, and people are probably sick of hearing my spiel, so I figured that I’d write it down and be done with it.

cgo is an amazing technology which allows Go programs to interoperate with C libraries. It’s a tremendously useful feature without which Go would not be in the position it is today. cgo is key to ability to run Go programs on Android and iOS.

However, and to be clear these are my opinions, I am not speaking for anyone else, I think cgo is overused in Go projects. I believe that when faced with reimplementing a large piece of C code in Go, programmers choose instead to use cgo to wrap the library, believing that it is a more tractable problem. I believe this is a false economy.

Obviously, there are some cases where cgo is unavoidable, most notably where you have to interoperate with a graphics driver or windowing system that is only available as a binary blob. But those cases where cgo’s use justifies its trade-offs are fewer and further between than many are prepared to admit.

Here is an incomplete list of trade-offs you make, possibly without realising them, when you base your Go project on a cgo library.

Slower build times
When you import "C" in your Go package, go build has to do a lot more work to build your code. Building your package is no longer simply passing a list of all the .go files in scope to a single invocation of go tool compile, instead:

The cgo tool needs to be invoked to generate the C to Go and Go to C thunks and stubs.
Your system C compiler has to be invoked for every C file in the package.
The individual compilation units are combined together into a single .o file.
The resulting .o file take a trip through the system linker for fix-ups against shared objects they reference.
All this work happens every time you compile or test your package, which is constantly, if you’re actively working in that package. The Go tool parallelises some of this work where possible, but your packages’ compile time just grew to include a full rebuild of all that C code.

It’s possible to work around this by pushing the cgo shims out into their own package, avoiding the compile time hit, but now you’ve had to restructure your application to work around a problem that you didn’t have before you started to use cgo.

Oh, and you have to debug C compilation failures on the various platforms your package supports.

Complicated builds
One of the goals of Go was to produce a language who’s build process was self describing; the source of your program contains enough information for a tool to build the project. This is not to say that using a Makefile to automate your build workflow is bad, but before cgo was introduced into a project, you may not have needed anything but the go tool to build and test. Afterwards, to set all the environment variables, keep track of shared objects and header files that may be installed in weird places, now you do.

Keep in mind that Go supports platforms that don’t ship with make out of the box, so you’ll have to dedicate some time to coming up with a solution for your Windows users.

Oh, and now your users have to have a C compiler installed, not just a Go compiler. They also have to install the C libraries your project depends on, so you’ll be taking on that support cost as well.

Cross compilation goes out the window
Go’s support for cross compilation is best in class. As of Go 1.5 you can cross compile from any supported platform to any other platform with the official installer available on the Go project website.

By default cgo is disabled when cross compiling. Normally this isn’t a problem if your project is pure Go. When you mix in dependencies on C libraries, you either have to give up the option to cross compile your product, or you have to invest time in finding and maintaining cross compilation C toolchains for all your targets.

Maybe if you work on a product that only communicates with clients over TCP sockets and you intend to run it in a SaaS model it’s reasonable to say that you don’t care about cross compilation. However, if you’re making a product which others will use, possibly integrated into their products, maybe it’s a monitoring solution, maybe it’s a client for your SaaS service, then you’ve locked them out of being able to easily cross compile.

The number of platforms that Go supports continues to grow. Go 1.5 added support for 64 bit ARM and PowerPC. Go 1.6 adds support for 64 bit MIPS, and IBM’s s390 architecture is touted for Go 1.7. RISC-V is in the pipeline. If your product relies on a C library, not only do you have the all problems of cross compilation described above, you also have to make sure the C code you depend on works reliably on the new platforms Go is supporting — and you have to do that with the limited debuggability a C/Go hybrid affords you. Which brings me to my next point.

You lose access to all your tools
Go has great tools; we have the race detector, pprof for profiling code, coverage, fuzz testing, and source code analysis tools. None of those work across the cgo blood/brain barrier.

Conversely excellent tools like valgrind don’t understand Go’s calling conventions or stack layout.  On that point, Ian Lance Taylor’s work to integrate clang’s memory sanitiser to debug dangling pointers on the C side will be of benefit for cgo users in Go 1.6.

Combing Go code and C code results in the intersection of both worlds, not the union; the memory safety of C, and the debuggability of a Go program.

Performance will always be an issue
C code and Go code live in two different universes, cgo traverses the boundary between them. This transition is not free and depending on where it exists in your code, the cost could be inconsequential, or substantial.

C doesn’t know anything about Go’s calling convention or growable stacks, so a call down to C code must record all the details of the goroutine stack, switch to the C stack, and run C code which has no knowledge of how it was invoked, or the larger Go runtime in charge of the program.

To be fair, Go doesn’t know anything about C’s world either. This is why the rules for passing data between the two have become more onerous over time as the compiler becomes better at spotting stack data that is no longer considered live, and the garbage collector becomes better at doing the same for the heap.

If there is a fault while in the C universe, the Go code has to recover enough state to at least print a stack trace and exit the program cleanly, rather than barfing up a core file.

Managing this transition across call stacks, especially where signals, threads and callbacks are involved is non trivial, and again Ian Lance Taylor has done a huge amount of work in Go 1.6 to improve the interoperability of signal handling with C.

The take away is that the transition between the C and Go world is non trivial, and it will never be free from overhead.

C calls the shots, not your code
It doesn’t matter which language you’re writing bindings or wrapping C code with; Python, Java with JNI, some language using libFFI, or Go via cgo; it’s C’s world, you’re just living in it.

Go code and C code have to agree on how resources like address space, signal handlers, and thread TLS slots are to be shared — and when I say agree, I actually mean Go has to work around the C code’s assumption. C code that can assume it always runs on one thread, or blithely be unprepared to work in a multi threaded environment at all.

You’re not writing a Go program that uses some logic from a C library, instead you’re writing a Go program that has to coexist with a belligerent piece of C code that is hard to replace, has the upper hand negotiations, and doesn’t care about your problems.

Deployment gets more complicated
Any presentation on Go to a general audience will contain at least one slide with these words:

Single, static binary

This is Go’s ace in the hole that has lead it to become a poster child of the movement away from virtual machines and managed runtimes. Using cgo, you give that up.

Depending on your environment, it’s probably possible to build your Go project into a deb or rpm, and assuming your other dependencies are also packaged, add them as an install dependency and push the problem off the operating system’s package manager. But that’s several significant changes to a build and deploy process that was previously as straight forward as go build && scp.

It is possible to compile a Go program entirely statically, but it is by no means simple and shows that the ramifications of including cgo in your project will ripple through your entire build and deploy life cycle.

Choose wisely
To be clear, I am not saying that you should not use cgo. But before you make that Faustian bargain, please consider carefully the qualities of Go that you’ll be giving up in return.

Related Posts:
Cross compilation with Go 1.5
An introduction to cross compilation with Go
Cross compilation just got a whole lot better in Go 1.5
An introduction to cross compilation with Go 1.1

为了能够重用已有的C语言库，我们在使用Golang开发项目或系统的时候难免会遇到Go和C语言混合编程，这时很多人都会选择使用cgo。话说cgo这个东西可算得上是让人又爱又恨，好处在于它可以让你快速重用已有的C语言库，无需再用Golang重造一遍轮子，而坏处就在于它会在一定程度上削弱你的系统性能。关于cgo的种种劣迹，Dave Cheney大神在他的博客上有一篇专门的文章《cgo is not Go》，感兴趣的同学可以看一看。但话说回来，有时候为了快速开发满足项目需求，使用cgo也实在是不得已而为之。

     在Golang中使用cgo调用C库的时候，如果需要引用很多不同的第三方库，那么使用#cgo CFLAGS:和#cgo LDFLAGS:的方式会引入很多行代码。首先这会导致代码很丑陋，最重要的是如果引用的不是标准库，头文件路径和库文件路径写死的话就会很麻烦。一旦第三方库的安装路径变化了，Golang的代码也要跟着变化，所以使用pkg-config无疑是一种更为优雅的方法，不管库的安装路径有何变化，我们都不需要修改Go代码，接下来本博主就用一个简单的例子来说明如何在cgo命令中使用pkg-config。

     首先假定我们在路径/home/ubuntu/third-parties/hello下安装了一个名称为hello的第三方C语言库，其目录结构如下所示，在hello_world.h中只定义了一个接口函数hello，该函数接收一个char *字符串作为变量并调用printf将其打印到标准输出。

# tree /home/ubuntu/third-parties/hello/
/home/ubuntu/third-parties/hello/
├── include
│   └── hello_world.h
└── lib
    ├── libhello.so
    └── pkgconfig
       └── hello.pc

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

     完成pkg-config描述文件的创建后，还需要将该描述文件的路径信息添加到PKG_CONFIG_PATH环境变量中，只有这样pkg-config才能正确获取这个C语言库的相关信息。此外，我们还需要将该C语言库的库文件路径添加到LD_LIBRARY_PATH环境变量中，具体命令如下：

# export PKG_CONFIG_PATH=/home/ubuntu/third-parties/hello/lib/pkgconfig
# pkg-config --list-all | grep libhello
libhello    libhello - The hello library just for testing pkgconfig
# export LD_LIBRARY_PATH=/home/ubuntu/third-parties/hello/lib

     在完成以上一系列准备工作之后，我们就可以开始编写Golang代码了，以下是Golang调用C语言接口的代码示例，我们只需要#cgo pkg-config: libhello和#include < hello_world.h >两行语句即可实现对hello函数的调用。如果C语言库的安装路径发生了变化，只需修改hello.pc这个描述文件即可，Golang代码无需重新修改和编译。

package main
// #cgo pkg-config: libhello
// #include < stdlib.h >
// #include < hello_world.h >
import "C"
import (
"unsafe"
)
func main() {
msg := "Hello, world!"
cmsg := C.CString(msg)
C.hello(cmsg)
C.free(unsafe.Pointer(cmsg))
}

     最后，编译该程序代码，查看可执行程序是否正确链接了C语言库，执行程序验证能否正确调用库函数功能。

# go build hello_world.go 
# ldd hello_world
linux-vdso.so.1 =>  (0x00007ffff63d3000)
libhello.so => /home/ubuntu/third-parties/hello/lib/libhello.so (0x00007fc31c0e1000)
libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fc31bec3000)
libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fc31bafe000)
       /lib64/ld-linux-x86-64.so.2 (0x00007fc31c2e3000)
# ./hello_world 
Hello, world!

     在以上步骤中需要关注的有两个地方：1）创建C语言库的pkg-config配置文件并将配置文件的路径添加到环境变量PKG_CONFIG_PATH中；2）C语言库文件的路径添加到环境变量LD_LIBRARY_PATH中，如果没有这一步，Go语言程序可以编译成功，但是可执行文件无法正确连接到C语言库，会出现如下情况：

# ldd hello_world
linux-vdso.so.1 =>  (0x00007fffa49e2000)
libhello.so => not found
libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007feb0fe93000)
libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007feb0face000)
       /lib64/ld-linux-x86-64.so.2 (0x00007feb100b1000)
       

Cgo可以创建出能够调用C代码的Go包。

在go命令行中使用cgo
使用cgo书写普通的Go代码，导入一个伪包“C”。Go代码就可以关联到如C.size_t的类型，如C.stdout的变量，或者如C.putchar的方法。

如果“C”包的导入紧接在注释之后，那个这个注释，被称为前导码【preamble】，就会作为编译Go包中的C语言部分的头文件。例如：

// #include <stdio.h>
// #include <errno.h>
import "C"
前导码【preamble】可以包含任何C语言代码，包括方法和变量声明及定义。这些代码后续可能会由Go代码引用，就像它们定义在了名为“C”的Go包中。所有在前导码中声明的名字都可能被使用，即使它们以小写字母开头。有一个例外：前导码中的static变量不能被Go代码引用；而static方法则可以被Go引用。

可以查看$GOROOT/misc/cgo/stdio和$GOROOT/misc/cgo/gmp中的例子。也可以在https://golang.org/doc/articles/c_go_cgo.html 中了解使用cgo的介绍。

CFLAGS， CPPFLAGS， CXXFLAGS， FFLAGS，和LDFLAGS可以在注释区域中，使用#cgo伪指令来定义，以改变C，C++，或Fortran编译器的行为。被多个指令定义的值会被联系在一起。指令还可以包含一个构建约束的列表，将其对系统的影响限制为满足其中一个约束条件。可以查看https://golang.org/pkg/go/build/#hdr-Build_Constraints 了解约束语法的详情。例如：

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"
或者，CPPFLAGS和LDFLAGS也可通过pkg-config工具获取，使用#cgo pkg-config:指令，后面跟上包名即可。比如：

// #cgo pkg-config: png cairo
// #include <png.h>
import "C"
默认的pkg-config工具可以通过设置PKG_CONFIG环境变量来改变。

当编译时，CGO_CFLAGS，CGO_CPPFLAGS，CGO_CXXFLAGS，CGO_FFLAGS和CGO_LDFLAGS这些环境变量都会从指令中提取出来，并加入到flags中。包特定的flags需要使用指令来设置，而不是通过环境变量，所以这些构建可以在未更改的环境中也能正常运行。

一个包中的所有cgo CPPFLAGS和CFLAGS指令会被连接起来，并用来编译包中的C文件。一个包中的所有CPPFLAGS和CXXFLAGS指令会被连接起来，并用来编译包中的C++文件。一个包中的所有CPPFLAGS和FFLAGS指令会被连接起来，并用来编译包中的Fortran文件。在这个程序中任何包内的所有LDFLAGS指令会被连接起来，并在链接时使用。所有的pkg-config指令会被连接起来，并同时发送给pkg-config，以添加到每个适当的命令行标志集中。

当cgo指令被转化【parse】时，任何出现${SRCDIR}字符串的地方，都会被替换为包含源文件的目录的绝对路径。这就允许预编译的静态库包含在包目录中，并能够正确的链接。例如，如果foo包在/go/src/foo目录下：

// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
就会被展开为：

// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
心得： // #cgo LDFLAGS: 可用来链接静态库。-L指定静态库所在目录，-l指定静态库文件名，注意静态库文件名必须有lib前缀，但是这里不需要写，比如上面的-lfoo实际找的是libfoo.a文件

当Go tool发现一个或多个Go文件使用了特殊导入“C”包，它就会在目录中寻找其它非Go文件，并将其编译为Go包的一部分。任何.c, .s, 或者.S文件会被C编译器编译。任何.cc, .cpp, 或者.cxx文件会被C++编译器编译。任何.f, .F, .for或者.f90文件会被fortran编译器编译。任何.h, .hh, .hpp或者.hxx文件不会被分开编译，但是，如果这些头文件修改了，那么C和C++文件就会被重新编译。默认的C和C++编译器有可能会分别被CC和CXX环境变量改变，那些环境变量可能包含命令行选项。

cgo tool对于本地构建默认是启用的。而在交叉编译时默认是禁用的。你可以在运行go tool时，通过设置CGO_ENABLED环境变量来控制它：设为1表示启用cgo， 设为0关闭它。如果cgo启用，go tool将会设置构建约束“cgo”。

在交叉编译时，你必须为cgo指定一个C语言交叉编译器。你可以在使用make.bash构建工具链时，设置CC_FOR_TARGET环境变量来指定，或者是在你运行go tool时设置CC环境变量来指定。与此相似的还有作用于C++代码的CXX_FOR_TARGET和CXX环境变量。

Go引用C代码
在Go文件中，C的结构体属性名是Go的关键字，可以加上下划线前缀来访问它们：如果x指向一个拥有属性名"type"的C结构体，那么就可以用x._type来访问这个属性。对于那些不能在Go中表达的C结构体字段（例如位字段或者未对齐的数据），会在Go结构体中被省略，而替换为适当的填充，以到达下一个属性，或者结构体的结尾。

标准C数字类型，可以通过如下名字访问：

C.char, 
C.schar(signed char), 
C.uchar(unsigned char), 
C.short, 
C.ushort(unsigned short), 
C.int, 
C.uint(unsigned int), 
C.long, 
C.ulong(unsigned long), 
C.longlong(long long), 
C.ulonglong(unsigned long long), 
C.float, 
C.double,
C.complexfloat(complex float), 
C.complexdouble(complex double)
C的void*类型由Go的unsafe.Pointer表示。C的__int128_t和__uint128_t由[16]byte表示。

如果想直接访问一个结构体，联合体，或者枚举类型，请加上如下前缀：struct_, union_, 或者enum_。比如C.struct_stat。

心得：但是如果C结构体用了typedef struct设置了别名，则就不需要加上前缀，可以直接C.alias访问该类型。

任何一个C类型T的尺寸大小，都可以通过C.sizeof_T获取，比如C.sizeof_struct_stat。

因为在通常情况下Go并不支持C的联合体类型【union type】，所以C的联合体类型，由一个等长的Go byte数组来表示。

Go结构体不能嵌入具有C类型的字段。

Go代码不能引用发生在非空C结构体末尾的零尺寸字段。如果要获取这个字段的地址（这也是对于零大小字段唯一能做的操作），你必须传入结构体的地址，并加上结构体的大小，才能算出这个零大小字段的地址。

cgo会将C类型转换为对应的，非导出的的Go类型。因为转换是非导出的，一个Go包就不应该在它的导出API中暴露C的类型：在一个Go包中使用的C类型，不同于在其它包中使用的同样C类型。

可以在多个赋值语境中，调用任何C函数（甚至是void函数），来获取返回值（如果有的话），以及C errno变量作为Go error（如果方法返回void，则使用 _ 来跳过返回值）。例如：

n, err = C.sqrt(-1)
_, err := C.voidFunc()
var n, err = C.sqrt(1)
调用C的方法指针目前还不支持，然而你可以声明Go变量来引用C的方法指针，然后在Go和C之间来回传递它。C代码可以调用来自Go的方法指针。例如：

package main

// typedef int (*intFunc) ();
//
// int
// bridge_int_func(intFunc f)
// {
//		return f();
// }
//
// int fortytwo()
// {
//	    return 42;
// }
import "C"
import "fmt"

func main() {
	f := C.intFunc(C.fortytwo)
	fmt.Println(int(C.bridge_int_func(f)))
	// Output: 42
}
在C中，一个方法参数被写为一个固定大小的数组，实际上需要的是一个指向数组第一个元素的指针。C编译器很清楚这个调用习惯，并相应的调整这个调用，但是Go不能这样做。在Go中，你必须显式的传入指向第一个元素的指针：C.f(&C.x[0])。

在Go和C类型之间，通过拷贝数据，还有一些特殊的方法转换。用Go伪代码定义如下：

// Go string to C string
// The C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CString(string) *C.char

// Go []byte slice to C array
// The C array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CBytes([]byte) unsafe.Pointer

// C string to Go string
func C.GoString(*C.char) string

// C data with explicit length to Go string
func C.GoStringN(*C.char, C.int) string

// C data with explicit length to Go []byte
func C.GoBytes(unsafe.Pointer, C.int) []byte
作为一个特殊例子，C.malloc并不是直接调用C的库函数malloc，而是调用一个Go的帮助函数【helper function】，该函数包装了C的库函数malloc，并且保证不会返回nil。如果C的malloc指示内存溢出，这个帮助函数会崩溃掉程序，就像Go自己运行时内存溢出一样。因为C.malloc从不失败，所以它不会返回包含errno的2值格式。
