---
title: pkgconfig
layout: post
category: linux
author: 夏泽民
---
用第三方库，就少不了要使用到第三方的头文件和库文件。我们在编译、链接的时候，必须要指定这些头文件和库文件的位置。

        对于一个比较大第三方库，其头文件和库文件的数量是比较多的。如果我们一个个手动地写，那将是相当麻烦的。所以，pkg-config就应运而生了。pkg-config能够把这些头文件和库文件的位置指出来，给编译器使用。如果你的系统装有gtk，可以尝试一下下面的命令$pkg-config --cflags gtk+-2.0。可以看到其输出是gtk的头文件的路径。

        我们平常都是这样用pkg-config的。$gcc main.c `pkg-config --cflags --libs gtk+-2.0` -o main

        上面的编译命令中，`pkg-config --cflags --libs gtk+-2.0`的作用就如前面所说的，把gtk的头文件路径和库文件列出来，让编译去获取。--cflags和--libs分别指定头文件和库文件。
         其实，pkg-config同其他命令一样，有很多选项，不过我们一般只会用到--libs和--cflags选项。
         https://blog.csdn.net/luotuo44/article/details/16970841
<!-- more -->
简单的说 pkg-config 维护了一个保存各个代码库的路径的数据库。当然这个”数据库” 非常的简单，其实就是一个特殊的目录，这个目录中有一系列的以 “.pc” 为后缀的文件。

因为pkg-config也只是一个命令，所以不是你安装了一个第三方的库，pkg-config就能知道第三方库的头文件和库文件所在的位置。pkg-config命令是通过查询XXX.pc文件而知道这些的。我们所需要做的是，写一个属于自己的库的.pc文件。

        但pkg-config又是如何找到所需的.pc文件呢？这就需要用到一个环境变量PKG_CONFIG_PATH了。这环境变量写明.pc文件的路径，pkg-config命令会读取这个环境变量的内容，这样就知道pc文件了。

        对于Ubuntu系统，可以用root权限打开/etc/bash.bashrc文件。在最后输入下面的内容。

        

        这样，pkg-config就会去/usr/local/lib/pkgconfig目录下，寻找.pc文件了。

.pc文件。只需写5个内容即可：Name、Description、Version、Cflags、Libs。

        比如简单的：

Name: opencv
Description:OpenCV pc file
Version: 2.4
Cflags:-I/usr/local/include
Libs:-L/usr/local/lib –lxxx –lxxx

其中，Cflags和Libs的写法，是使用了-I -L -l这些gcc的编译选项。

编译不同目录下的多个文件
1)  以相对路径的方式直接包含头文件
为了能够使用add函数，必须包含add所在的头文件。 最简单的方法是直接在main.cpp文件中，用相对路径包含head.h文件.即 #include”function/head.h”。

此时，编译命令为 ：$g++ main.cpp function/head.cpp -o main

        这种用相对路径包含头文件的方式有很多弊端。当function目录改成其它名字，或者head.h文件放到其它目录了，这时都要对main.cpp文件进行修改，如果head.h头文件被很多其它文件包含的话，这个工作量就大多了。



2)  用编译选项 –I(大写i)
 其实，可以想一下，为什么iostream文件不在当前目录下，就可以直接使用呢？这是因为，编译器会在一些默认的目录下(/usr/include，/usr/inlucde/c++/4.4.3等目录)搜索头文件。所以，iostream头文件不用添加。但我们不能每写一个头文件就放到那里。

        知道了原理，现在来用一下一个编译选项 –I（include的缩写）用来告诉编译器，还可以去哪里找头文件。

        使用这个编译命令，$g++ main.cpp function/head.cpp -Ifunction -o main

        此时main.cpp文件写成

#include  <iostream>
#include  <head.h>
可以看到head.h文件是用<>而不是””。想一下C语言书中，两者的区别。这说明，用-I选项，相当于说明了一条标准路径。



3)  使用.o文件
        此时，对于head.cpp文件，在编译命令中，还是要用到路径function/head.cpp。现在的想法是去掉这个。这时可以先根据head.cpp文件生成一个.o文件，然后就可以了去掉那个路径了。

 

        先cd 到function目录。

        输入命令：$g++ -c head.cpp -o head.o

        生成一个head.o目标文件，

 

        此时把生成的head.o文件复制到function的父目录，就是main.cpp所在的目录。

        然后回到function的父目录，输入命令$g++ main.cpp head.o -Ifunction -o main
  此时，直接使用head.o即可，无需head.cpp了。但头文件head.h还是要的。因为编译的时候要用到。链接的时候不用头文件。这个可以拆分成两条命令

        $g++ -c main.cpp -Ifunction -o main.o

        $g++ main.o head.o -o main

        第一条是编译命令，后一条是链接命令。
二、  静态库
        虽然上面说到的先生成.o目标文件，但如果function目录下有多个.cpp文件。那么就要为每一个.cpp文件都生成一个.o文件，这个工作量是会比较大。此时可以用静态库。静态库是把多个目标文件打包成一个文件。Anarchive(or static library) is simply a collection of object filesstored as a single file（摘自《Advanced Linux Programming》）。
首先，生成.o目标文件。

        cd进入function目录，输入命令$g++ -c sub.cpp add.cpp  这会分别为两个源文件目标文件，即生成sub.o和add.o文件。

        然后，打包生成静态库，即.a文件。

        输入命令$ar -cr libhead.a add.o sub.o



        注意：

        1)  命令中，.a文件要放到 .o文件的前面

        2)  .a文件的格式。要以lib作为前缀, .a作为后缀。

        选项 c是代表create，创建.a文件。

        r是代表replace，如果之前有创建过.a文件，现在为了提高性能而更改了add函数里面的代码，此时，就可以用r选项来代替之前.a文件里面的add.o

        可以用命令$ar -t libhead.a 查看libhead.a文件里面包含了哪些目标文件。其执行结果自然为add.o  sub.o
回到main.cpp文件所在的目录。

        输入命令:$g++ main.cpp -Ifunction -Lfunction -lhead -o main 生成可执行程序

        

        现在要解释一下使用静态库要用到的-L和-l（小写的L）选项。

        -L表示要使用的静态库的目录。这和前面所讲的-I(大写i)差不多，就是用来告诉编译器去哪里找静态库。因为可能-L所指明的目录下有很多静态库，所以除了要告诉去哪里找之外，还要告诉编译器，找哪一个静态库。此时，就要用到-l(小写L)了。它用来说明链接的时候要用到哪个静态库。

        注意：

        1. 注意是使用-lhead,而不是-llibhead

        命令中是使用-lhead,这是因为编译器会自动在库中添加lib前缀和.a后缀。

        2. 要把-l放到命令的尽可能后的位置，必须放到源文件的后面。

        如果使用命令中的顺序，将出现下面错误。
三、  动态库
        使用命令$g++ -c -fPIC add.cpp sub.cpp生成位置无关的目标文件。

        使用命令$g++ -shared -fPIC add.o sub.o -o libhead.so 生成.so动态链接库。

 

        利用动态库生成可执行文件 $g++ -Ifunction -Lfunction -lhead main.cpp -o main


        尝试运行. $./main 得到下面错误
这说明加载的时候没有找到libhead.so动态库。这是因为，Linux查找动态库的地方依次是

环境变量LD_LIBRARY_PATH指定的路径
缓存文件/etc/ld.so.cache指定的路径
默认的共享库目录，先是/lib，然后是/usr/lib
        运行./main时，明显这三个地方都没有找到。因为我们没有把libhead.so文件放到那里。



        其实，我们可以在生成可执行文件的时候指定要链接的动态库是在哪个地方的。

        $g++ -Ifunction ./libhead.so main.cpp -o main

        这个命令告诉可执行文件，在当前目录下查找libhead.so动态库。注意这个命令不再使用-L 和 -l了。

        另外，还可以使用选项-Wl,-rpath,XXXX.其中XXXX表示路径。

        

四、  打造自己的库目录和头文件目录
        三个要解决的东西：指定编译时的头文件路径、指定链接时的动态库路径、指定运行时Linux加载动态库的查找路径

 

 

1.指定运行时Linux加载动态库的查找路径
        利用前面所说的Linux程序运行时查找动态库的顺序，让Linux在运行程序的时候，去自己指定的路径搜索动态库。



        可以修改环境变量LD_LIBRARY_PATH或者修改/etc/ld.so.cache文件。这里选择修改/etc/ld.so.cache文件。

        1)  创建目录/mylib/so。这个目录可以用来存放自己以后写的所有库文件。由于是在系统目录下创建一个目录，所以需要root权限

        2)  创建并编辑一个mylib.conf文件。输入命令$sudo vim /etc/ld.so.conf.d/mylib.conf

        在mylib.conf文件中输入 /mylib/so 

        保存，退出。

        3)  重建/etc/ld.so.cache文件。输入命令$sudo ldconfig

 

        输入下面命令，生成main文件。注意，其链接的时候是用function目录下的libhead动态库。

        $g++  -Ifunction -Lfunction -lhead main.cpp 

        直接运行./main。并没有错误。可以运行。说明，Linux已经会在运行程序时自动在/mylib/so目录下找动态链接库了。

 

2. 指定编译时的头文件路径
        先弄清编译器搜索头文件的顺序。

        1.先搜索当前目录（使用include””时）

        2.然后搜索-I指定的目录

        3.再搜索环境变量CPLUS_INCLUDE_PATH、 C_INCLUDE_PATH。两者分别是g++、gcc使用的。

        4.最后搜索默认目录 /usr/include  和 /usr/local/include等

 

        知道这些就简单了。输入下面命令。编辑该文件。

        $sudo vim /etc/bash.bashrc  这个文件是作用了所有Linux用户的，如果不想影响其他用户，那么就编辑~/.bashrc文件。

        打开文件后，去到文件的最后一行。输入下面的语句。

        

        修改环境变量。然后保存并推出。

        输入命令$bash 或者直接打开一个新的命令行窗口，使得配置信息生效。原理可以参考:http://blog.csdn.net/luotuo44/article/details/8917764



        此时，可以看到 已经可以不用-I选项 下面编译命令能通过了。

        $g++ -Lfunction -lhead mian.cpp -o main



3.指定链接时的动态库路径
         需要注意的是，链接时使用动态库和运行时使用动态库是不同的。

        同样先搞清搜索顺序：

        1. 编译命令中-L指定的目录

        2. 环境变量LIBRARY_PATH所指定的目录

        3. 默认目录。/lib、/usr/lib等。

 

        接下来和指定头文件路径一样。输入命令$sudo vim /etc/bash.bashrc   在文件的最后一行输入

        

        保存，退出。

        同样，输入bash，使得配置信息生效。

        这是终极目标了。其中-lhead是不能省的。因为，编译器要知道，你要链接到哪一个动态库。当然，如果想像C运行库那样，链接时默认添加的动态库，那么应该也是可以通过设置，把libhead.so库作为默认库。但并不是所有的程序都会使用这个库。要是设置为默认添加的，反而不好。

一、编译和连接

        一般来说，如果库的头文件不在 /usr/include 目录中，那么在编译的时候需要用 -I 参数指定其路径。由于同一个库在不同系统上可能位于不同的目录下，用户安装库的时候也可以将库安装在不同的目录下，所以即使使用同一个库，由于库的路径的 不同，造成了用 -I 参数指定的头文件的路径也可能不同，其结果就是造成了编译命令界面的不统一。如果使用 -L 参数，也会造成连接界面的不统一。编译和连接界面不统一会为库的使用带来麻烦。

       为了解决编译和连接界面不统一的问题，人们找到了一些解决办法。其基本思想就是：事先把库的位置信息等保存起来，需要的时候再通过特定的工具将其中有用的 信息提取出来供编译和连接使用。这样，就可以做到编译和连接界面的一致性。其中，目前最为常用的库信息提取工具就是下面介绍的 pkg-config。

       pkg-config 是通过库提供的一个 .pc 文件获得库的各种必要信息的，包括版本信息、编译和连接需要的参数等。这些信息可以通过 pkg-config 提供的参数单独提取出来直接供编译器和连接器使用。

The pkgconfig package contains tools for passing the include path and/or library paths to build tools during the make file execution.

pkg-config is a function that returns meta information for the specified library.

The default setting for PKG_CONFIG_PATH is /usr/lib/pkgconfig because of the prefix we use to install pkgconfig. You may add to PKG_CONFIG_PATH by exporting additional paths on your system where pkgconfig files are installed. Note that PKG_CONFIG_PATH is only needed when compiling packages, not during run-time.

        在默认情况下，每个支持 pkg-config 的库对应的 .pc 文件在安装后都位于安装目录中的 lib/pkgconfig 目录下。例如，我们在上面已经将 Glib 安装在 /opt/gtk 目录下了，那么这个 Glib 库对应的 .pc 文件是 /opt/gtk/lib/pkgconfig 目录下一个叫 glib-2.0.pc 的文件：

prefix=/opt/gtk/
exec_prefix=${prefix}
libdir=${exec_prefix}/lib
includedir=${prefix}/include

glib_genmarshal=glib-genmarshal
gobject_query=gobject-query
glib_mkenums=glib-mkenums

Name: GLib
Description: C Utility Library
Version: 2.12.13
Libs: -L${libdir} -lglib-2.0 
Cflags: -I${includedir}/glib-2.0 -I${libdir}/glib-2.0/include 

        使用 pkg-config 的 --cflags 参数可以给出在编译时所需要的选项，而 --libs 参数可以给出连接时的选项。例如，假设一个 sample.c 的程序用到了 Glib 库，就可以这样编译：

$ gcc -c `pkg-config --cflags glib-2.0` sample.c

然后这样连接：

$ gcc sample.o -o sample `pkg-config --libs glib-2.0`

或者上面两步也可以合并为以下一步：

$ gcc sample.c -o sample `pkg-config --cflags --libs glib-2.0`

可以看到：由于使用了 pkg-config 工具来获得库的选项，所以不论库安装在什么目录下，都可以使用相同的编译和连接命令，带来了编译和连接界面的统一。

使用 pkg-config 工具提取库的编译和连接参数有两个基本的前提：

库本身在安装的时候必须提供一个相应的 .pc 文件。不这样做的库说明不支持 pkg-config 工具的使用。pkg-config 必须知道要到哪里去寻找此 .pc 文件。
GTK+ 及其依赖库支持使用 pkg-config 工具，所以剩下的问题就是如何告诉 pkg-config 到哪里去寻找库对应的 .pc 文件，这也是通过设置搜索路径来解决的。

       对于支持 pkg-config 工具的 GTK+ 及其依赖库来说，库的头文件的搜索路径的设置变成了对 .pc 文件搜索路径的设置。.pc 文件的搜索路径是通过环境变量 PKG_CONFIG_PATH 来设置的，pkg-config 将按照设置路径的先后顺序进行搜索，直到找到指定的 .pc 文件为止。

安装完 Glib 后，在 bash 中应该进行如下设置：
$ export PKG_CONFIG_PATH=/opt/gtk/lib/pkgconfig:$PKG_CONFIG_PATH

可以执行下面的命令检查是否 /opt/gtk/lib/pkgconfig 路径已经设置在 PKG_CONFIG_PATH 环境变量中：

$ echo $PKG_CONFIG_PATH

这样设置之后，使用 Glib 库的其它程序或库在编译的时候 pkg-config 就知道首先要到 /opt/gtk/lib/pkgconfig 这个目录中去寻找 glib-2.0.pc 了（GTK+ 和其它的依赖库的 .pc 文件也将拷贝到这里，也会首先到这里搜索它们对应的 .pc 文件）。之后，通过 pkg-config 就可以把其中库的编译和连接参数提取出来供程序在编译和连接时使用。

另外还需要注意的是：环境变量的设置只对当前的终端窗口有效。如果到了没有进行上述设置的终端窗口中，pkg-config 将找不到新安装的 glib-2.0.pc 文件、从而可能使后面进行的安装（如 Glib 之后的 Atk 的安装）无法进行。


       在我们采用的安装方案中，由于是使用环境变量对 GTK+ 及其依赖库进行的设置，所以当系统重新启动、或者新开一个终端窗口之后，如果想使用新安装的 GTK+ 库，需要如上面那样重新设置 PKG_CONFIG_PATH 和 LD_LIBRARY_PATH 环境变量。

这种使用 GTK+ 的方法，在使用之前多了一个对库进行设置的过程。虽然显得稍微繁琐了一些，但却是一种最安全的使用 GTK+ 库的方式，不会对系统上已经存在的使用了 GTK+ 库的程序（比如 GNOME 桌面）带来任何冲击。

为了使库的设置变得简单一些，可以把下面的这两句设置保存到一个文件中（比如 set_gtk-2.10 文件）:

export PKG_CONFIG_PATH=/opt/gtk/lib/pkgconfig:$PKG_CONFIG_PATH
export LD_LIBRARY_PATH=/opt/gtk/lib:$LD_LIBRARY_PATH
之后，就可以用下面的方法进行库的设置了（其中的 source 命令也可以用 . 代替）：

$ source set_gtk-2.10

只有在用新版的 GTK+ 库开发应用程序、或者运行使用了新版 GTK+ 库的程序的时候，才有必要进行上述设置。

           如果想避免使用 GTK+ 库之前上述设置的麻烦，可以把上面两个环境变量的设置在系统的配置文件中（如 /etc/profile）或者自己的用户配置文件中（如 ~/.bash_profile） ；库的搜索路径也可以设置在 /etc/ld.so.conf 文件中，等等。这种设置在系统启动时会生效，从而会导致使用 GTK+ 的程序使用新版的 GTK+ 运行库，这有可能会带来一些问题。当然，如果你发现用新版的 GTK+ 代替旧版没有什么问题的话，使用这种设置方式是比较方便的。加入到~/.bashrc中，例如：
PKG_CONFIG_PATH=/opt/gtk/lib/pkgconfig
重启之后：
[root@localhost ~]# echo $PKG_CONFIG_PATH
/opt/gtk/lib/pkgconfig 


二、运行时
        库文件在连接（静态库和共享库）和运行（仅限于使用共享库的程序）时被使用，其搜索路径是在系统中进行设置的。一般 Linux 系统把 /lib 和 /usr/lib 两个目录作为默认的库搜索路径，所以使用这两个目录中的库时不需要进行设置搜索路径即可直接使用。对于处于默认库搜索路径之外的库，需要将库的位置添加到 库的搜索路径之中。设置库文件的搜索路径有下列两种方式，可任选其一使用：

1、环境变量 LD_LIBRARY_PATH 中指明库的搜索路径。2、 /etc/ld.so.conf 文件中添加库的搜索路径。         将自己可能存放库文件的路径都加入到/etc/ld.so.conf中是明智的选择 ^_^
添加方法也极其简单，将库文件的绝对路径直接写进去就OK了，一行一个。例如：
/usr/X11R6/lib
/usr/local/lib
/opt/lib
        需要注意的是：第二种搜索路径的设置方式对于程序连接时的库（包括共享库和静态库）的定位已经足够了，但是对于使用了共享库的程序的执行还是不够的。这是 因为为了加快程序执行时对共享库的定位速度，避免使用搜索路径查找共享库的低效率，所以是直接读取库列表文件 /etc/ld.so.cache 从中进行搜索的。/etc/ld.so.cache 是一个非文本的数据文件，不能直接编辑，它是根据 /etc/ld.so.conf 中设置的搜索路径由 /sbin/ldconfig 命令将这些搜索路径下的共享库文件集中在一起而生成的（ldconfig 命令要以 root 权限执行）。因此，为了保证程序执行时对库的定位，在 /etc/ld.so.conf 中进行了库搜索路径的设置之后，还必须要运行 /sbin/ldconfig 命令更新 /etc/ld.so.cache 文件之后才可以。ldconfig ,简单的说，它的作用就是将/etc/ld.so.conf列出的路径下的库文件 缓存到/etc/ld.so.cache 以供使用。因此当安装完一些库文件，(例如刚安装好glib)，或者修改ld.so.conf增加新的库路径后，需要运行一下/sbin /ldconfig使所有的库文件都被缓存到ld.so.cache中，如果没做，即使库文件明明就在/usr/lib下的，也是不会被使用的，结果编译 过程中抱错，缺少xxx库，去查看发现明明就在那放着，搞的想大骂computer蠢猪一个。 ^_^
       在程序连接时，对于库文件（静态库和共享库）的搜索路径，除了上面的设置方式之外，还可以通过 -L 参数显式指定。因为用 -L 设置的路径将被优先搜索，所以在连接的时候通常都会以这种方式直接指定要连接的库的路径。
       
        前面已经说明过了，库搜索路径的设置有两种方式：在环境变量 LD_LIBRARY_PATH 中设置以及在 /etc/ld.so.conf 文件中设置。其中，第二种设置方式需要 root 权限，以改变 /etc/ld.so.conf 文件并执行 /sbin/ldconfig 命令。而且，当系统重新启动后，所有的基于 GTK2 的程序在运行时都将使用新安装的 GTK+ 库。不幸的是，由于 GTK+ 版本的改变，这有时会给应用程序带来兼容性的问题，造成某些程序运行不正常。为了避免出现上面的这些情况，在 GTK+ 及其依赖库的安装过程中对于库的搜索路径的设置将采用第一种方式进行。这种设置方式不需要 root 权限，设置也简单：

$ export LD_LIBRARY_PATH=/opt/gtk/lib:$LD_LIBRARY_PATH

可以用下面的命令查看 LD_LIBRAY_PATH 的设置内容：

$ echo $LD_LIBRARY_PATH

至此，库的两种设置就完成了。

https://www.cnblogs.com/youxin/p/4271978.html


