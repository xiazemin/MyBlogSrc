---
title: golang 实现 LD_PRELOAD 拦截 libc
layout: post
category: golang
author: 夏泽民
---
https://blog.gopheracademy.com/advent-2015/libc-hooking-go-shared-libraries/
先把 socket() 的 libc 调用拦截下来
package main

// #cgo LDFLAGS: -ldl
// #include <stddef.h>
// #include <netinet/in.h>
// #include "network_hook.h"
import "C"
import (
	"fmt"
	"syscall"
)

func init() {
	C.libc_hook_init()
}

//export socket
func socket(domain C.int, type_ C.int, protocol C.int) C.int {
	fmt.Println(fmt.Sprintf("open socket from thread id: %v", syscall.Gettid()))
	return C.orig_socket(domain, type_, protocol)
}

func main() {
}
把这个编译成 so

go build -buildmode=c-shared -o libmotrix.so main.go
然后调用

LD_PRELOAD=libmotrix.so curl http://www.baidu.com
就会把 socket 的 lib 调用给拦截下来

把 socket() 透明转交给 libc
network_hook.h

#ifndef __DLSYM_WRAPPER_H__
#define __DLSYM_WRAPPER_H__

void libc_hook_init();
int orig_socket(int, int, int);

#endif
network_hook.c

#include <dlfcn.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include <netdb.h>
#include <math.h>
#include "network_hook.h"
#include "_cgo_export.h"

#define RTLD_NEXT	((void *) -1l)

#define HOOK_SYS_FUNC(name) if( !orig_##name##_func ) { orig_##name##_func = (name##_pfn_t)dlsym(RTLD_NEXT,#name); }

typedef int (*socket_pfn_t)(int, int, int);
static socket_pfn_t orig_socket_func;

void libc_hook_init() {
    HOOK_SYS_FUNC( socket );
}

int orig_socket(int domain, int type, int protocol) {
    return orig_socket_func(domain, type, protocol);
}

其中

orig_socket 是函数暴露给 golang 调用
orig_socket_func 是函数指针，指向了libc的原来的实现
socket_pfn_t 是socket()这个函数指针的类型定义
原来文章里的实现在我的机器上报错，原因未知。感觉这样用 cgo 包装一下更简单直接一些，还少了一次反射。
<!-- more -->
Golang调用C分两个步骤：1 写一个C的wrapper，这个很简单；2 对wrapper做编译，这个步骤有点复杂，而且涉及众多中间文件。应该是有办法用自动化的工具简化这个过程的。

先来展示一下C程序。为了将描述集中在如何调用上，C的程序很简单：

prints.h
#ifndef PRINTS_HEAD
void prints(char* str);
#endif

prints.c
#include "prints.h"
#include <stdio.h>

void prints(char* str)
{
  printf("%s\n", str);
}
之后是Golang对C的一个wrapper：

prints.go
package prints

//#include "prints.h"
// // some comment
import "C"

func Prints(s string) {
  p := C.CString(s);
  C.prints(p);
}
需要注意的是红色高亮的几行。在编译过程中，go会根据import "C"之前的几行注释生成一个c程序，并将这个c程序里的符号导入到模块C里，最后由import "C"再导入到go程序里。如果需要在其他go程序里调用api，需要参照prints.go里的Prints函数（要导出的go模块需要首字母大写）写一个wrapper func。其中对c程序里符号的引用都需要通过C来引用，包括一些c的类型定义，比如传给c api的int需要通过C.int来定义，而字符串则是C.CString。

有了这几个文件，就可以编译一个可以在go里加载的库了。以下都是在x86 linux下操作过程，如果是其他环境，请替换相应的编译命令。

cgo prints.go

编译wrapper，生成文件：
cgo_defun.c：根据prints.go里标红的注释，生成用于在go里调用的c符号和函数
_cgo_gotypes.go：_cgo_defun.c里的符号在go里对应的定义
_cgo.o
prints.cgo1.go：根据prints.go生成的go wrapper func
prints.cgo2.c：根据prints.go生成的c wrapper func
8g -o go.8 prints.cgo1.go cgo_gotypes.go
编译go wrapper相关的文件，生成_go.8
8c -FVw -I"/home/lizh/go/src/pkg/runtime/" _cgo_defun.c
编译c wrapper的通用部分，生成_cgo_defun.8
gopack grc prints.a _go_.8 _cgo_defun.8

对上面两个编译好的wrapper打包，生成prints.a

cp prints.a $GOROOT/pkg/linux_386/

将生成的prints.a放到go的包目录下

之后是对c部分的编译：

gcc -m32 -fPIC -O2 -o prints.cgo2.o -c prints.cgo2.c 
gcc -m32 -fPIC -O2 -o prints.o -c prints.c
gcc -m32 -o prints.so prints.o prints.cgo2.o -shared
根据prints.c和prints.cgo2.c生成prints.so，是一个可供go程序引入的动态库。通过objdump查看prints.so的符号，可以发现： prints：需要引入的c api符号 _cgo_prints：由go生成的对c api的wrapper，具体可以查看prints.cgo2.c

cp prints.so /home/lizh/go/pkg/linux_386/

将编译好的动态库放入go的包目录下

之后就可以在go里调用prints这个c函数了：

package main

import "prints"

func main() {
  s := "Hello world!";
  prints.Prints(s);
}
查看生成的调用程序，可以看到对$GOROOT/pkg/linux_386/libcgo.so和$GOROOT/pkg/linux_386/prints.so两个动态库的引用。发布时需要将这两个库放到发布环境里。


