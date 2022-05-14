---
title: LD_PRELOAD/DYLD_INSERT_LIBRARIES libc hook
layout: post
category: linux
author: 夏泽民
---
一、LD_PRELOAD是什么
LD_PRELOAD是Linux系统的一个环境变量，它可以影响程序的运行时的链接（Runtime linker），它允许你定义在程序运行前优先加载的动态链接库。这个功能主要就是用来有选择性的载入不同动态链接库中的相同函数。通过这个环境变量，我们可以在主程序和其动态链接库的中间加载别的动态链接库，甚至覆盖正常的函数库。一方面，我们可以以此功能来使用自己的或是更好的函数（无需别人的源码），而另一方面，我们也可以以向别人的程序注入程序，从而达到特定的目的。
编译、设置指令
gcc mystrcmp.c -fPIC -shared -o libmystrcmp.so      #编译动态链接库
gcc myverifypasswd.c -L. -lmystrcmp -o myverifypasswd      #编译主程序
export LD_LIBRARY_PATH=/home/LD_PRELOAD      #指定动态链接库所在目录位置
ldd myverifypasswd      #显示、确认依赖关系
./myverifypasswd      #运行主程序myverifypasswd
二、LD_PRELOAD运用总结
定义与目标函数完全一样的函数，包括名称、变量及类型、返回值及类型等
将包含替换函数的源码编译为动态链接库
通过命令 export LD_PRELOAD="库文件路径"，设置要优先替换动态链接库
如果找不替换库，可以通过 export LD_LIBRARY_PATH=库文件所在目录路径，设置系统查找库的目录
替换结束，要还原函数调用关系，用命令unset LD_PRELOAD 解除
想查询依赖关系，可以用ldd 程序名称
<!-- more -->
LD_PRELOAD：
在Unix操作系统的动态链接库的世界中，LD_PRELOAD就是这样一个环境变量
用以指定预先装载的一些共享库或目标文件，且无论程序是否依赖这些共享库或者文件，LD_PRELOAD指定的这些文件都会被装载
其优先级比LD_LIBRARY_PATH自定义的进程的共享库查找路径的执行还要早
全局符号介入
指在不同的共享库（对象）中存在同名符号时，一个共享对象中的全局符号被另一个共享对象的同名全局符号覆盖
因为LD_PRELOAD指定的共享库或者目标文件的装载顺序十分靠前，几乎是程序运行最先装载的，所以其中的全局符号如果和后面的库中的全局符号重名的话，就会覆盖后面装载的共享库或者目标文件中的全局符号。
因为装载顺序和全局符号介入的原理
它可以影响程序的运行时的链接（Runtime linker），它允许你定义在程序运行前优先加载的动态链接库。
这个功能主要就是用来有选择性的载入Unix操作系统不同动态链接库中的相同函数。通过这个环境变量，我们可以在主程序和其动态链接库的中间加载别的动态链接库，甚至覆盖正常的函数库。

这个环境变量带来了很大的安全隐患，解决方法有两个：
使用静态链接
使用SETUID/SETGID，有SUID权限的程序执行加载时会忽略这个环境变量
LD_PRELOAD实际上和dlsym函数族相关。具体看dlsym的使用
void *dlsym(void *handle, const char *symbol)
handle表示dlopen的返回值或者预定义的两个值(RTLD_NEXT RTLD_DEFAULT)
symbol表示符号名称
当找到符号时返回符号地址，没有找到返回空指针。
RTLD_DEFAULT和RTLD_NEXT提供了在多个动态库中查询相同函数或者变量的方法。
当多个动态库存在相同的符号时，RTLD_DEFAULT查找的是这些动态库的第一个。并将第一个设置为当前的动态库，RTLD_NEXT查找的时当前动态库的下一个动态库。
因此，LD_PRELOAD影响的实际上就是RTLD_DEFAULT，它将指定的动态库放在查找的第一个位置。

定义与目标函数完全一样的函数，包括名称、变量及类型、返回值及类型等
将包含替换函数的源码编译为动态链接库
通过命令 export LD_PRELOAD="库文件路径"，设置要优先替换动态链接库
如果找不替换库，可以通过 export LD_LIBRARY_PATH=库文件所在目录路径，设置系统查找库的目录
替换结束，要还原函数调用关系，用命令unset LD_PRELOAD 解除
想查询依赖关系，可以用ldd 程序名称

应用
获得 root 权限。你想多了！你不会通过这种方法绕过安全机制的。（一个专业的解释是：如果 ruid != euid，库不会通过这种方法预加载的。）
欺骗游戏：取消随机化。这是我演示的第一个示例。对于一个完整的工作案例，你将需要去实现一个定制的 random() 、rand_r()、random_r()，也有一些应用程序是从 /dev/urandom 之类的读取，你可以通过使用一个修改过的文件路径来运行原始的 open() 来把它们重定向到 /dev/null。而且，一些应用程序可能有它们自己的随机数生成算法，这种情况下你似乎是没有办法的（除非，按下面的第 10 点去操作）。但是对于一个新手来说，它看起来很容易上手。
欺骗游戏：让子弹飞一会 。实现所有的与时间有关的标准函数，让假冒的时间变慢两倍，或者十倍。如果你为时间测量和与时间相关的 sleep 或其它函数正确地计算了新的值，那么受影响的应用程序将认为时间变慢了（你想的话，也可以变快），并且，你可以体验可怕的 “子弹时间” 的动作。或者 甚至更进一步，你的共享库也可以成为一个 DBus 客户端，因此你可以使用它进行实时的通讯。绑定一些快捷方式到定制的命令，并且在你的假冒的时间函数上使用一些额外的计算，让你可以有能力按你的意愿去启用和禁用慢进或快进任何时间。
研究应用程序：列出访问的文件。它是我演示的第二个示例，但是这也可以进一步去深化，通过记录和监视所有应用程序的文件 I/O。
研究应用程序：监视因特网访问。你可以使用 Wireshark 或者类似软件达到这一目的，但是，使用这个诀窍你可以真实地控制基于 web 的应用程序发送了什么，不仅是看看，而是也能影响到交换的数据。这里有很多的可能性，从检测间谍软件到欺骗多用户游戏，或者分析和逆向工程使用闭源协议的应用程序。
研究应用程序：检查 GTK 结构 。为什么只局限于标准库？让我们在所有的 GTK 调用中注入一些代码，因此我们就可以知道一个应用程序使用了哪些组件，并且，知道它们的构成。然后这可以渲染出一个图像或者甚至是一个 gtkbuilder 文件！如果你想去学习一些应用程序是怎么管理其界面的，这个方法超级有用！
在沙盒中运行不安全的应用程序。如果你不信任一些应用程序，并且你可能担心它会做一些如 rm -rf / 或者一些其它不希望的文件活动，你可以通过修改传递到文件相关的函数（不仅是 open ，也包括删除目录等）的参数，来重定向所有的文件 I/O 操作到诸如 /tmp 这样地方。还有更难的诀窍，如 chroot，但是它也给你提供更多的控制。它可以更安全地完全 “封装”，但除非你真的知道你在做什么，不要以这种方式真的运行任何恶意软件。
实现特性 。zlibc 是明确以这种方法运行的一个真实的库；它可以在访问文件时解压文件，因此，任何应用程序都可以在无需实现解压功能的情况下访问压缩数据。
修复 bug。另一个现实中的示例是：不久前（我不确定现在是否仍然如此）Skype（它是闭源的软件）从某些网络摄像头中捕获视频有问题。因为 Skype 并不是自由软件，源文件不能被修改，这就可以通过使用预加载一个解决了这个问题的库的方式来修复这个 bug。
手工方式 访问应用程序拥有的内存。请注意，你可以通过这种方式去访问所有应用程序的数据。如果你有类似的软件，如 CheatEngine/scanmem/GameConqueror 这可能并不会让人惊讶，但是，它们都要求 root 权限才能工作，而 LD_PRELOAD 则不需要。事实上，通过一些巧妙的诀窍，你注入的代码可以访问所有的应用程序内存，从本质上看，是因为它是通过应用程序自身得以运行的。你可以修改这个应用程序能修改的任何东西。你可以想像一下，它允许你做许多的底层的侵入


简单的例子
 #include <stdio.h>
 #include <stdlib.h>
 #include <unistd.h>
int main(int argc, const char *argv[]) {
  char buffer[1000];
  int amount_read;
  int fd;
  fd = fileno(stdin);
  if ((amount_read = read(fd, buffer, sizeof buffer)) == -1) {
    perror("error reading");
    return EXIT_FAILURE;
  }
  if (fwrite(buffer, sizeof(char), amount_read, stdout) == -1) {
    perror("error writing");
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}
这个程序很简单就是从stdin读取，输出到stdout
然后我们编译得到二进制执行文件

gcc main.c -o out
然后我们再写一个read函数
inject.c

#include <string.h>
ssize_t read(int fd, void *data, size_t size) {
  strcpy(data, "I love cats");
  return 12;
}
编译获得.so库文件

gcc -shared -fPIC -o inject.so inject.c
然后我们开始注入，用LD_PRELOAD环境变量指定预先加载的库为上面的库

LD_PRELOAD=$PWD/inject.so ./out
现在运行out，它的输出将会是打印 "I love cats" 而不是我们的输入值。
OS X 略有不同，但是类似

gcc -shared -fPIC -o inject.dylib inject.c
DYLD_INSERT_LIBRARIES=$PWD/inject.dylib DYLD_FORCE_FLAT_NAMESPACE=1 ./out
用这种办法，就可以在不改源码的情况下改写指定函数的逻辑。然而这还不是我们要的。

链接:
编译过程只是把源文件翻译成二进制而已，这个二进制还不能直接执行，这个时候就需要做一个动作，将翻译成的二进制与需要用到库绑定在一块。
系统里有个头文件 <dlfcn.h>，提供了dlsym方法从链接器里取符号。接下去我们要做的事情就发生在 链接 这一步。

#define _GNU_SOURCE
#include <string.h>
#include <dlfcn.h>
#include <stdio.h>
typedef ssize_t (*real_read_t)(int, void *, size_t);
// 找到原来的read
ssize_t real_read(int fd, void *data, size_t size) {
  return ((real_read_t)dlsym(RTLD_NEXT, "read"))(fd, data, size);
}
// 自己的read
ssize_t read(int fd, void *data, size_t size) {
  strcpy(data, "I love cats");
  return 12;
}
RTLD_NEXT是一个c的宏，作用是查找符号表里下一个名字为参数值的符号，在这个例子里，就是查找下一个read函数，也就是最开始那个demo里的read函数。
最后整合起来

#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdio.h>
#include <string.h>
typedef ssize_t (*real_read_t)(int, void *, size_t);
ssize_t real_read(int fd, void *data, size_t size) {
  return ((real_read_t)dlsym(RTLD_NEXT, "read"))(fd, data, size);
}
ssize_t read(int fd, void *data, size_t size) {
  ssize_t amount_read;
  // Perform the actual system call
  amount_read = real_read(fd, data, size);
  // Our evil code, for example, log, record, stat. whatever.
strcpy(data, "I love cats");
  // Behave just like the regular syscall would
  return amount_read;
}
最后重新编译这个带其它功能的代码

gcc -shared -fPIC -ldl -o inject.so inject.c
LD_PRELOAD=$PWD/inject.so gcc -shared -fPIC -ldl -o inject.so inject.c

参考：
http://www.goldsborough.me/c/low-level/kernel/2016/08/29/16-48-53-the_-ld_preload-_trick/
https://github.com/ccurtsinger/interpose
