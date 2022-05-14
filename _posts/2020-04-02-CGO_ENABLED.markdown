---
title: CGO_ENABLED
layout: post
category: golang
author: 夏泽民
---
在macOS下启用CGO_ENABLED的交叉编译
在启用CGO_ENABLED的情况下，尝试使用下面命令进行Windows平台的交叉编译：

$ CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -x -v -ldflags "-s -w"
出现错误如下：

# runtime/cgo
gcc_libinit_windows.c:7:10: fatal error: 'windows.h' file not found
安装mingw-w64
# piao @ PiaodeMacBook-Pro in ~ [11:10:19]
$ brew install mingw-w64
==> Downloading https://homebrew.bintray.com/bottles/mingw-w64-5.0.4_1.mojave.bottle.tar.gz
Already downloaded: /Users/piao/Library/Caches/Homebrew/downloads/954c462f9298678f85a2ca518229e941d1daed366c84c339900c756e7ca8ad25--mingw-w64-5.0.4_1.mojave.bottle.tar.gz
==> Pouring mingw-w64-5.0.4_1.mojave.bottle.tar.gz
🍺  /usr/local/Cellar/mingw-w64/5.0.4_1: 7,915 files, 747.7MB
# piao @ PiaodeMacBook-Pro in ~ [11:10:56]
$ which x86_64-w64-mingw32-gcc
/usr/local/bin/x86_64-w64-mingw32-gcc
编译x64
可执行文件
$ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o test_x64.exe
静态库
$ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -buildmode=c-archive -x -v -ldflags "-s -w" -o bin/x64/x64.a main.go
动态库
将-buildmode=c-archive改为-buildmode=c-shared即可

编译x86
可执行文件
$ CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -x -v -ldflags "-s -w" -o test_x86.exe
静态库
$ CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -buildmode=c-archive -x -v -ldflags "-s -w" -o bin/x86/x86.a main.go
动态库
将-buildmode=c-archive改为-buildmode=c-shared即可

熟悉golang的人都知道，golang交叉编译很简单的，只要设置几个环境变量就可以了
# mac上编译linux和windows二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 
 
# linux上编译mac和windows二进制
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
 
# windows上编译mac和linux二进制
SET CGO_ENABLED=0 SET GOOS=darwin SET GOARCH=amd64 go build main.go
SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build main.go
GOOS和GOARCH的值有哪些，可以网上搜，很多的

但是交叉编译是不支持CGO的，也就是说如果你的代码中存在C代码，是编译不了的，比如说你的程序中使用了sqlite数据库，在编译go-sqlite驱动时按照上面的做法是编译不通过的

需要CGO支持的，要将CGO_ENABLED的0改为1，也就是CGO_ENABLED=1，此外还需要设置编译器，例如我想在linux上编译arm版的二进制，需要这样做：
# Build for arm
CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc go build
这个arm-linux-gnueabi-gcc是个啥东西，怎么安装，如果你系统是ubuntu的话，可以按照下面命令安装：
sudo apt-get install g++-arm-linux-gnueabi
sudo apt-get install gcc-arm-linux-gnueabi
安装成功后就可以编译了，但是如果你想编译mac版呢，或者想在mac上编译linux版，window版咋办，一个一个安装效率太慢，系统命令可以安装还好，系统命令不支持，那就得自己去搜，然后找到地址，下载，安装，费时又费力

github上有这个工具：https://github.com/karalabe/xgo

它是一个docker镜像，里面集成了各种平台的编译器，按照它的教程，很轻松的可以编译出各个平台的二进制文件，安装的时候比较耗时，需要下载大概1个G的数据，但是效果可是杠杠的

默认是编译所有平台的二进制的，会有些耗时，如果只需要某个特定平台的二进制，可以使用-targets参数

注意：是-targets而不是--targets，我自己测试的时候--targets是会失败的

附：golang如何让编译生产的二进制文件变小

把Go程序变小的办法是：
go build -ldflags "-s -w" (go install类似)
-s去掉符号表（然后panic时候的stack trace就没有任何文件名/行号信息了，
这个等价于普通C/C++程序被strip的效果），
-w去掉DWARF调试信息，得到的程序就不能用gdb调试了。

比如，server.go是一个简单的http server，用了net/http包。
$ go build server.go
$ ls -l server
-rwxr-xr-x 1 minux staff 4507004 2012-10-25 14:16 server
$ go build -ldflags "-s -w" server.go
$ ls -l server
-rwxr-xr-x 1 minux staff 2839932 2012-10-25 14:16 server
-s和-w也可以分开使用，一般来说如果不打算用gdb调试，-w基本没啥损失。
-s的损失就有点大了

<!-- more -->
在macOS和Linux下gcc，在window下需要安装MinGW。同时需要保证环境变量CGO_ENABLED被设置为1，这是表示cgo是否被启用状态。在本地构建时CGO_ENABLED默认启用，在交叉构建cgo是默认禁用的。比如交叉构建ARM环境运行GO程序，需要手动设置CGO_ENABLED环境变量。

主角登场

首先下载go，搭建本地的go开发环境。下载地址https://golang.google.cn/dl/，尽量选择比较新的版本，因为版本不一样，对于cgo的支持程度有一些区别，为了避免出现一些不必要的问题，建议使用最新稳定版本。在下使用的是1.12.7。

提炼关键

安装go的时候，主要配置以下几个环境变量：

GOROOT：go安装的根目录(例如:D:\go\)，安装程序会自动写入系统环境变量，如果没有，请自行加入系统环境变量中。

GOBIN：go的可执行文件的存放目录(%GOROOT%\bin)

PATH：系统的环境变量，需要将go环境变量拼接到此变量后面。

GOPATH：go的工作空间，包含go的开发目录和依赖包目录，如果没有配置请手动配置它，GOPATH工作空间主要有三个子目录：

src：包含go的源码文件

pkg：包对象，编译好的库文件

bin：可执行文件目录

配置好后可运行go env命令来查看go环境是否配置正确




MinGW是什么

MinGW是“Minimalist GNU for Windows”的缩写，是原生Microsoft Windows应用程序的极简主义开发环境。MinGW提供了一个完整的开源编程工具集，适用于本机MS-Windows应用程序的开发，并且不依赖于任何第三方C-Runtime DLL。MinGW编译器提供对Microsoft C运行时功能和某些特定于语言的运行时的访问。MinGW是Minimalist，它不会，也绝不会尝试为MS-Windows上的POSIX应用程序部署提供POSIX运行时环境。如果您希望在此平台上部署POSIX应用程序，请考虑使用Cygwin。

也就是说很多开源库需要的库是由MinGW提供，在window环境下如果想用这些库，就需要MinGW的帮助。后续的编译opencv也需要用到。





Window环境配置

安装go，下载地址为：https://dl.google.com/go/go1.12.7.windows-amd64.msi，安装是比较简单的，点击安装程序按提示安装，需要注意的是记住你go的安装目录，后面需要用到。我的安装目录是：D:\go\。检查下go是否能在控制台下运行，调出运行框，在键盘上按下”win图标键”+”R” 输入cmd回车


输入命令go version，如果输入go的版本信息那说明go安装成功了，如果没有输出有可能是go没有加入到PATH中，需要将D:\go\bin加入到系统的环境变量PATH中。


配置Go的工作空间，需要配置环境变量GOPATH，就是你开发go项目的项目目录，最新版本在安装的时候就自动配置好了，早期版本是需要手动配置，但是建议还是手动配置下，因为系统默认是C:\Users\当前用户名\go

配置如下：

 



开发go的IDE

本人使用的是LiteIDE，简单好用，轻量级，而且还能编译出不同操作系统下的可执行文件，在window下开发挺好用的。当然每个人都可根据个人喜好下载自己熟悉的IDE进行开发。

使用LiteIDE如果已经安装了go并且要升级go到新版本，这时候LiteIDE的代码提示功能可能就不能用，需要重新再去编译下gocode，然后覆盖到LiteIDE安装目录下的bin目录。下载地址：

http://sourceforge.net/projects/liteide/files/x36/liteidex36.windows-qt5.9.5.zip/download

默认使用系统环境变量配置




这里也介绍下IntelliJ IDEA的配置吧，这个ide挺耗内存，个人感觉挺重，不是很喜欢用，但是它的功能比较强大，提示功能也比较好，因此介绍下这个IDE的配置吧。

IntelliJ IDEA需要安装go的插件才能进行go语言的开发，可以在Marketplace上下载，也可以下载插件到本地之后，在本地安装插件。建议下载插件到本地

https://nchc.dl.sourceforge.net/project/liteide/x36/liteidex36.windows-qt5.9.5.zip

从插件商店下载不一定能下载得下来。

首先完成go的环境配置之后，并且也配置了GOPATH的工作空间配置后，本人配置为D:\go-project，在IntelliJ IDEA新建项目，新建Go项目



选择我们的GOPATH作为工作空间


在GOPATH下新建目录src、pkg、bin这三个目录，这三个目录上述已经有介绍了。


新建一个go的编译配置，名称自己取一个。Run kind选择Directory，Directory选择D:\go-project\src，这里注意是要选择GOPATH下的src目录。Output directory选择D:\go-project\bin，这个你可以根据自己喜欢配置输出可执行文件的目录，我是习惯放GOPATH工作空间下的bin目录。在src下新建一个main.go的源码文件，输出hello world测试我们配置是否能正常。





这里还需要注意的是IntelliJ IDEA也是有GOPATH的配置的，只是它默认使用的是系统GOPATH的配置，当然你可以不指定使用系统的GOPATH，自己去修改自身GOPATH的路径。在file–settings–Languages & Frameworks下，使用自己指定的GOPATH时，Use GOPATH that’s defined in system environment前面的勾要去掉。


这样配置的意思我大概做个解释，GOPATH相当于我们一个大项目src就是放项目的目录，所以你编译类型选目录的时候，就要把路径指向src目录。Src目录下又可以包含多个独立的子项目，用目录区分开来，我个人实际开发中就是这样去划分的，比如我GOPATH的src目录下有2个项目，一个叫app核心代码库，一个叫web站点代码，web会引用app库，app为web项目提供各种类、函数，那么这时你要新建一个go 编译配置，因为这时候你启动的项目变成时web了，而不是默认src目录下的main.go了。于是我们新建了一个goweb的编译配置如下：


编写下app/lib.go文件代码很简单就是输出hello world


编写web/main.go代码如下，并运行成功输入hello world。至此所有配置完成。


跨平台编译方案

工作中经常使用的go服务器基本都是Linux操作系统，所以本人的编译方案就是在window开发完成，并调试没问题后，将go的源码同步到Linux服务器，当然Linux服务器上也需要用go的环境配置好后，然后再Linux上执行go build -i编译成Linux下的可执行文件。为此本人利用go开发了一套自动部署的系统，完成开发环境自动同步到服务器上，并完成自动部署的功能，后续如果有时间可以将此项目写成教程公开出来。


IntelliJ IDEA在Window环境下也是可以编译Linux可执行文件的，但是不支持cgo就是了，当你的go项目没有用到cgo的情况下，是可以直接在Window下编译可以在Linux系统下运行的程序的。需要加入环境变量，配置关键如下：

1.GOARCH：操作系统架构，值有amd64(64位操作系统)、386(32位系统)、arm(arm系统)。

2.GOOS：操作系统名，值有linux、windows、freebsd、darwin。

3.CGO_ENABLED：cgo是否启用，如果是window编译不包含cgo的项目，必须设置为0。


安装MinGW,看自己的系统是32位还是64位，如果是64位建议下载MinGW-w64,本人系统就是64位的所以安装的是MinGW-w64。下载地址：http://www.mingw-w64.org/doku.php/download


Msys2也是可以，这个包括了MinGw和Cygwin,相当于这两者的结合体。需要配置，后续有机会讲些下Msys2的配置。

安装时设置选择如下：


安装目录依据自己喜好而定，之后一直next直到安装完成，安装可能会比较慢，因为需要下载所需要的文件，耐心等待安装完成。

需要将bin目录加入到PATH中，环境配置:



Linux环境配置

安装go，下载地址https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz，解压到目录。

解压命令：

tar zxvf go1.12.7.linux-amd64.tar.gz

解压完成后，目录为/home/wjp/go



为go配置环境变量：

执行命令：

sudo vim /etc/profile

在最底部加入新配置

export GOROOT=/home/wjp/go

export GOPATH=/home/wjp/go-src #go的工作空间

export GOBIN=$GOROOT/bin

export PATH=$PATH:$GOBIN:$GOPATH/bin


加载配置执行命令

source /etc/profile

测试环境是否配置成功执行命令

go version，如果输出版本则安装成功


由于cgo是调用gcc进行编译的，而Linux系统本身就自带gcc，如果没有自行安装下就可以了。

Linux如果是服务器版的，那开发go效率不高，建议是在window下开发完成后，在拷贝到Linux上进行编译生成。除非你用的是Linux桌面版，可以使用一些IDE进行开发，这里就不做介绍了。