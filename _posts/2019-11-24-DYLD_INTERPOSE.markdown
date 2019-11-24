---
title: DYLD_INTERPOSE
layout: post
category: linux
author: 夏泽民
---
https://opensource.apple.com/source/dyld/dyld-353.2.1/include/mach-o/dyld-interposing.h
//  演示代码 
// #import <mach-o/dyld-interposing.h>
// from dyld-interposing.h
#define DYLD_INTERPOSE(_replacement,_replacee) __attribute__((used)) static struct{ const void* replacement; const void* replacee; } _interpose_##_replacee __attribute__ ((section ("__DATA,__interpose"))) = { (const void*)(unsigned long)&_replacement, (const void*)(unsigned long)&_replacee };

ssize_t hacked_write(int fildes, const void *buf, size_t nbyte)
{
    printf("[++++]into hacked_write－－－by piaoyun");
    return write(fildes, buf, nbyte);
}

DYLD_INTERPOSE(hacked_write, write);



// 再来个演示代码：

// 编译
// cc -dynamiclib main.c -o libHook.dylib -Wall
// 强行注入ls测试
// DYLD_INSERT_LIBRARIES=libHook.dylib ls
#include <malloc/malloc.h>

#define DYLD_INTERPOSE(_replacement,_replacee) \
__attribute__((used)) static struct{ const void* replacement; const void* replacee; } _interpose_##_replacee \
__attribute__ ((section ("__DATA,__interpose"))) = { (const void*)(unsigned long)&_replacement, (const void*)(unsigned long)&_replacee };


void *hacked_malloc(size_t size){
    void *ret = malloc(size);
    
    malloc_printf("+ %p %d\n", ret, size);
    return ret;
}

void hacked_free(void *freed){
    malloc_printf("- %p\n", freed);
    free(freed);
}

DYLD_INTERPOSE(hacked_malloc, malloc)
DYLD_INTERPOSE(hacked_free, free);


<!-- more -->
/*
 * Copyright (c) 2005-2008 Apple Computer, Inc. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 * 
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 * 
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 * 
 * @APPLE_LICENSE_HEADER_END@
 */

#if !defined(_DYLD_INTERPOSING_H_)
#define _DYLD_INTERPOSING_H_

/*
 *  Example:
 *
 *  static
 *  int
 *  my_open(const char* path, int flags, mode_t mode)
 *  {
 *    int value;
 *    // do stuff before open (including changing the arguments)
 *    value = open(path, flags, mode);
 *    // do stuff after open (including changing the return value(s))
 *    return value;
 *  }
 *  DYLD_INTERPOSE(my_open, open)
 */

#define DYLD_INTERPOSE(_replacement,_replacee) \
   __attribute__((used)) static struct{ const void* replacement; const void* replacee; } _interpose_##_replacee \
            __attribute__ ((section ("__DATA,__interpose"))) = { (const void*)(unsigned long)&_replacement, (const void*)(unsigned long)&_replacee };

#endif

https://www.cnblogs.com/cobbliu/p/7347923.html?utm_source=itdadao&utm_medium=referral

There are times when you want to wrap a library function in order to provide some additional functionality. A common example of this is wrapping the standard library’s malloc() and free() so that you can easily track memory allocations in your program. While there are several techniques for wrapping library functions, one well-known method is using dlsym() with RTLD_NEXT to locate the wrapped function’s address so that you can correctly forward calls to it.

Problem
So what can go wrong? Let’s look at an example:

LibWrap.h

void* memAlloc(size_t s);
// Allocate a memory block of size 's' bytes.
void memDel(void* p);
// Free the block of memory pointed to by 'p'.
LibWrap.c

#define _GNU_SOURCE
#include <dlfcn.h>
#include "LibWrap.h"

static void* malloc(size_t s) {
   // Wrapper for standard library's 'malloc'.
   // The 'static' keyword forces all calls to malloc() in this file to resolve
   // to this functions.
   void* (*origMalloc)(size_t) = dlsym(RTLD_NEXT,"malloc");
   return origMalloc(s);
}

static void free(void* p) {
   // Wrapper for standard library's 'free'.
   // The 'static' keyword forces all calls to free() in this file to resolve
   // to this functions.
   void (*origFree)(void*) = dlsym(RTLD_NEXT,"free");
   origFree(p);
}

void* memAlloc(size_t s) {
   return malloc(s);
   // Call the malloc() wrapper.
}

void memDel(void* p) {
   free(p);
   // Call the free() wrapper.
}
Main.c

#include <malloc.h>
#include "LibWrap.h"

int main() {
   struct mallinfo beforeMalloc = mallinfo();
   printf("Bytes allocated before malloc: %d\n",beforeMalloc.uordblks);

   void* p = memAlloc(57);
   struct mallinfo afterMalloc = mallinfo();
   printf("Bytes allocated after malloc: %d\n",afterMalloc.uordblks);

   memDel(p);
   struct mallinfo afterFree = mallinfo();
   printf("Bytes allocated after free: %d\n",afterFree.uordblks);

   return 0;
}
First compile LibWrap.c into a shared library:

$ gcc -Wall -Werror -fPIC -shared -o libWrap.so LibWrap.c
Next compile Main.c and link it against the libWrap.so that we just created:

$ gcc -Wall -Werror -o Main Main.c ./libWrap.so -ldl
Time to run the program!

$ ./Main
Bytes allocated before malloc: 0
Bytes allocated after malloc: 80
Bytes allocated after free: 0
So far, so good. No surprises. We allocated a bunch of memory and then freed it. The statistics returned by mallinfo() confirm this.

Out of curiosity, let’s look at ldd output for the application binary we created.


$ ldd Main
       linux-vdso.so.1 =>  (0x00007fff1b1fe000)
       ./libWrap.so (0x00007fe7d2755000)
       libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007fe7d2542000)
       libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fe7d217c000)
       /lib64/ld-linux-x86-64.so.2 (0x00007fe7d2959000)

Take note of the relative placement of libWrap.so with respect to libc.so.6: libWrap.socomes before libc.so.6. Remember this. It will be important later.

Now for fun, let’s re-compile Main.c with libc.so.6 explicitly specified on the command-line and coming before libWrap.so:

$ gcc -Wall -Werror -o Main Main.c /lib/x86_64-linux-gnu/libc.so.6 ./libWrap.so -ldl
Re-run:

$ ./Main
Bytes allocated before malloc: 0
Bytes allocated after malloc: 80
Bytes allocated after free: 80
Uh oh, why are we leaking memory all of a sudden? We de-allocate everything we allocate, so why the memory leak?

It turns out that the leak is occurring because we are not actually forwarding malloc() and free() calls to libc.so.6‘s implementations. Instead, we are forwarding them to malloc() and free() inside ld-linux-x86-64.so.2!

“What are you talking about?!” you might be asking.

Well, it just so happens that ld-linux-x86-64.so.2, which is the dynamic linker/loader, has its own copy of malloc() and free(). Why? Because ld-linux has to allocate memory from the heap before it loads libc.so.6. But the version of malloc/free that ld-linuxhas does not actually free memory!


[RTLD_NEXT] will find the next occurrence of a function in the search order after the current library. This allows one to provide a wrapper around a function in another shared library.But why does libWrap.so forward calls to ld-linux instead of libc? The answer comes down to how dlsym() searches for symbols when RTLD_NEXT is specified. Here’s the relevant excerpt from the dlsym(3) man page:— dlsym(3)

To understand this better, take a look at ldd output for the new Main binary:

$ ldd Main
        linux-vdso.so.1 =>  (0x00007fffe1da0000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f32c2e91000)
        ./libWrap.so (0x00007f32c2c8f000)
        libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f32c2a8a000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f32c3267000)
Unlike earlier, libWrap.so comes after libc.so.6. So when dlsym() is called inside libWrap.so to search for functions, it skips libc.so.6 since it precedes libWrap.so in the search order list. That means the searches continue through to ld-linux-x86-64.so.2where they find linker/loader’s malloc/free and return pointers to those functions. And so, libWrap.so ends up forwading calls to ld-linux instead of libc!


The answer is unfortunately no. At OptumSoft, we recently encountered this very same memory leak with a binary compiled using the standard ./configure && make on x86-64 Ubuntu 14.04.1 LTS. For reasons we don’t understand, the linking order for the binary was such that using dlsym() with RTLD_NEXT to lookup malloc/free resulted in pointers to implementations inside ld-linux. It took a ton of effort and invaluable help from Mozilla’s rr tool to root-cause the issue. After the whole ordeal, we decided to write a blog post about this strange behavior in case someone else encounters it in the future.At this point you might be wondering: We ran a somewhat funky command to build our application and then encountered a memory leak due to weird library linking order caused by said command. Isn’t this whole thing a silly contrived scenario?

Solution
If you find dlsym() with RTLD_NEXT returning pointers to malloc/free inside ld-linux, what can you do?

For starters, you need to detect that a function address indeed does belong to ld-linuxusing dladdr():

void* func = dlsym(RTLD_NEXT,"malloc");
Dl_info dlInfo;
if(!dladdr(func,&dlInfo)) {
   // dladdr() failed.
}
if(strstr(dlInfo.dli_fname,"ld-linux")) {
   // 'malloc' is inside linker/loader.
}
Once you have figured out that a function is inside ld-linux, you need to decide what to do next. Unfortunately, there is no straightforward way to continue searching for the same function name in all other libraries. But if you know the name of a specific library in which the function exists (e.g. libc), you can use dlopen() and dlsym() to fetch the desired pointer:

void* handle = dlopen("libc.so.6",RTLD_LAZY);
// NOTE: libc.so.6 may *not* exist on Alpha and IA-64 architectures.
if(!handle) {
   // dlopen() failed.
}
void* func = dlsym(handle,"free");
if(!func) {
   // Bad! 'free' was not found inside libc.
}
 

Summary
One can use dlsym() with RTLD_NEXT to implement wrappers around malloc() and free().
Due to unexpected linking behavior, dlsym() when using RTLD_NEXT can return pointers to malloc/free implementations inside ld-linux (dynamic linker/loader). Using ld-linux‘s malloc/free for general heap allocations leads to memory leaks because that particular version of free() doesn’t actually release memory.
You can check if an address returned by dlsym() belongs to ld-linux via dladdr(). You can also lookup a function in a specific library using dlopen() and dlsym().

dlsym用法
1. 包含头文件 #include<dlfcn.h>

2. 函数定义 void *dlsym(void *handle, const char* symbol);

handle是使用dlopen函数之后返回的句柄，symbol是要求获取的函数的名称，函数，返回值是void*,指向函数的地址，供调用使用

dlsym与dlopen的以如下例子解释：

#include<dlfcn.h>

void * handle = dlopen("./testListDB.so",RTLD_LAZY);

如果createListDB函数定义为int32_t createListDB(std::string);

那么dlsym的用法则为：int32_t  (*create_listDB)(std::string) = reinterpret_cast<int32_t (*)(std::string)>(dlsym(handle, "createListDB"))

createListDB库函数的定义要用extern来声明，这样在主函数中才能通过createListDB来查找函数，

 相比于已知函数的所在动态库，函数dlsym的参数RTLD_NEXT可以在对函数实现所在动态库名称未知的情况下完成对库函数的替代。这提供了巨大的便利。但是凡是有一利必有一弊，在使用该参数时，需要注意一些问题。

使用的函数文件
main函数.c

#include <stdio.h>
#include <malloc.h>

int main (void) {
    struct mallinfo bfmalloc = mallinfo();
    printf("bfmalloc: %d\n",bfmalloc.uordblks);
    int *p = malloc(32);
    struct mallinfo afmalloc  = mallinfo();
    printf("afmalloc: %d\n",afmalloc.uordblks);
    free(p);
    struct mallinfo fmalloc = mallinfo();
    printf("fmalloc: %d\n",fmalloc.uordblks);
    return 0;
}

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
包装函数实现文件.h

void *memAlloc(size_t size);
void memFree(void *ptr);
1
2
包装实现文件.c

#ifdef RUNTIME
#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <malloc.h>
#include <assert.h>
#include "mymalloc.h"
void *malloc(size_t size) {
    void*(*mallocp)(size_t size);
    char *error;
    mallocp = dlsym(RTLD_NEXT,"malloc");
    assert(dlerror() == NULL);
    char *ptr = mallocp(size);
    return ptr;
}
void free(void *ptr) {
    assert(ptr != NULL);
    char *error;
    void (*freep)(void *ptr);
    freep = dlsym(RTLD_NEXT,"free");
    assert(dlerror() == NULL);
    freep(ptr);
}
void *memAlloc (size_t size) {
    void *ptr = malloc(size);
    printf("ret:%d : %p\n",size,ptr);
    return ptr;
}
void memFree(void *ptr) {
    free(ptr);
    printf("free ret: %p\n",ptr);
}
#endif
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
我们首先编译并运行一下程序，看看运行结果
编译步骤：

gcc -DRUNTIME -shared -fPIC -o mymalloc.so mymalloc.c
gcc -o t main.c ./mymalloc.so -ldl
1
2
运行./t

我们可以很明显的看到，我们并没有释放掉我们所申请的空间，虽然我们显示的调用了内存分配函数free。那么为什么会造成内存泄露呢？简单的来说原因是因为在使用了RTLD_NEXT参数后，编译器在查询函数的实现库时，对库的访问顺序不同造成的，详细的原因大家可以参考下面这篇博文：
Dangers of using dlsym() with RTLD_NEXT

其次在使用该参数时还碰到了一个问题：不能在malloc的实现文件里面加入任何带有缓冲的输出语句，否则便或造成段错误：有兴趣的读者可以试一下。这个是加了printf语句之后，函数的调用栈

至于程序为什么会在__vfprintf_internal()子函数处宕掉，我有一个大胆的想法：

当加了带缓冲的输出语句：以printf为例，程序的编译是没有问题的，但是当程序运行时便会产生错误。原因是因为：printf函数首先会将格式控制串语句写入到自己的缓冲区里面，此时修改了内存中的只读文本段，所以会导致程序触发一般保护故障处理流程。那么我们如何印证这个猜想呢？可以利用stderr标准错误输出来打印我们想要输出的语句，此时你便会发现程序得以正常运行，并没有引发任何的故障。
如果有强迫症就想在malloc的实现文件中加入printf语句，那么可以像上面的memAlloc一样加一层包装函数，在该函数中可以执行替代原来库函数功能的执行逻辑。

当然了你也可以采用在编译时期或者链接时期完成对函数的 “打桩”，在此便不在赘述。

I reduced my problem using below test codes,

main.cc

#include <iostream>

int main(int argc, const char** argv) {
  void init2();
  init2();
  return 0;
}
2.cc

#include <iostream>

int init2() {
  void init1();
  init1();
  std::cout<<"init2 called\n";
  return 0;
}
1.cc

#include <dlfcn.h>
#include <pthread.h>
#include <stdio.h>
#include <iostream>

typedef FILE* (*FopenFunction)(const char* path, const char* mode);

static FopenFunction g_libc_fopen = NULL;

void init1() {
 g_libc_fopen = reinterpret_cast<FopenFunction>(
          dlsym(RTLD_NEXT, "fopen"));

 std::cout<<"init1: fopen addr:"<<(void*)g_libc_fopen<<"\n";
}

__attribute__ ((__visibility__("default")))
FILE* fopen_override(const char* path, const char* mode)  __asm__ ("fopen");

__attribute__ ((__visibility__("default")))
FILE* fopen_override(const char* path, const char* mode) {
  return g_libc_fopen(path, mode);
}
Compiled 1.cc into a lib1.so and 2.cc into lib2.so like below,

g++ 1.cc -shared -ldl -fvisibility=default -fPIC -o lib1.so -L.
g++ 2.cc -shared -ldl -fvisibility=default -fPIC -o lib2.so -l1 -L.
g++ main.cc -l2 -l1 -L.
Above steps will produce lib1.so, lib2.so and a.out. The problem here is while running the executable a.out, it is unable to lookup the original "fread" symbol when using dlsym(RTLD_NEXT).

The output is,

arunprasadr@demo:~/works/myex/c++/rtdl_next$ LD_LIBRARY_PATH=./ ./a.out
init1: fopen addr:0
init2 called
But if the change the link process of lib2.so(like below), it seems to be working

g++ 2.cc -shared -ldl -fvisibility=default -fPIC -o lib2.so -L.
g++ main.cc -l2 -l1 -L.
LD_LIBRARY_PATH=./ ./a.out
output:

arunprasadr@demo:~/works/myex/c++/rtdl_next$ LD_LIBRARY_PATH=./ ./a.out
init1: fopen addr:0x7f9e84a9e2c0
init2 called
Can anyone please explain what is happening in the background? Thanks in advance.

gcc glibc dlopen dlsym uclibc
shareimprove this question
asked Apr 8 '14 at 8:57

Arunprasad Rajkumar
1,14299 silver badges2626 bronze badges
add a comment
1 Answer
activeoldestvotes

2

This is an interesting (and unexpected for me) result.

First, using your original commands, I observe:

LD_DEBUG=symbols,bindings LD_LIBRARY_PATH=./ ./a.out |& grep fopen
     10204: symbol=fopen;  lookup in file=/lib/x86_64-linux-gnu/libm.so.6 [0]
     10204: symbol=fopen;  lookup in file=/lib64/ld-linux-x86-64.so.2 [0]
     10204: symbol=fopen;  lookup in file=/lib/x86_64-linux-gnu/libgcc_s.so.1 [0]
     10204: symbol=fopen;  lookup in file=/lib/x86_64-linux-gnu/libdl.so.2 [0]
init1: fopen addr:0
Compare this with the with the same output, but removing -l1 from lib2.sos link line:

LD_DEBUG=symbols,bindings LD_LIBRARY_PATH=./ ./a.out |& grep fopen
     10314: symbol=fopen;  lookup in file=/usr/lib/x86_64-linux-gnu/libstdc++.so.6 [0]
     10314: symbol=fopen;  lookup in file=/lib/x86_64-linux-gnu/libc.so.6 [0]
     10314: binding file ./lib1.so [0] to /lib/x86_64-linux-gnu/libc.so.6 [0]: normal symbol `fopen'
init1: fopen addr:0x7f03692352c0
The question then is: why isn't the loader searching libc.so.6 for fopen in the first case?

The answer: the loader has a linear list of libraries in the _r_debug.r_map link chain, and for RTLD_NEXT will search libraries after the one that is calling dlopen.

Is the order of libraries different for case 1 and case 2? You bet:

case 1:

LD_LIBRARY_PATH=./ ldd ./a.out
    linux-vdso.so.1 =>  (0x00007fff2f1ff000)
    lib2.so => ./lib2.so (0x00007f54a2b12000)
    libstdc++.so.6 => /usr/lib/x86_64-linux-gnu/libstdc++.so.6 (0x00007f54a27f1000)
    libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f54a2430000)
    lib1.so => ./lib1.so (0x00007f54a222e000)
    libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f54a1f32000)
    /lib64/ld-linux-x86-64.so.2 (0x00007f54a2d16000)
    libgcc_s.so.1 => /lib/x86_64-linux-gnu/libgcc_s.so.1 (0x00007f54a1d1b000)
    libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f54a1b17000)
case 2:

LD_LIBRARY_PATH=./ ldd ./a.out
    linux-vdso.so.1 =>  (0x00007fff39fff000)
    lib2.so => ./lib2.so (0x00007f8502329000)
    lib1.so => ./lib1.so (0x00007f8502127000)
    libstdc++.so.6 => /usr/lib/x86_64-linux-gnu/libstdc++.so.6 (0x00007f8501e05000)
    libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f8501a45000)
    libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f8501841000)
    libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f8501544000)
    /lib64/ld-linux-x86-64.so.2 (0x00007f850252d000)
    libgcc_s.so.1 => /lib/x86_64-linux-gnu/libgcc_s.so.1 (0x00007f850132e000)
It should now be clear, that for case 2 libc.so.6 follows lib1.so, but for case 1 it does not.

I do not yet understand what causes this particular ordering though. I'll have to think some more about it.


OS X系统中，仅有很少的进程只需要内核加载器就可以完成，几乎所有的程序都是动态连接的，通常采用/usr/lib/dyld作为动态链接器。
作为一个私有的加载器，dyld提供了一些独有的特性，如函数拦截等。DYLD_INTERPOSE宏定义允许一个库将其函数实现替换为另一个函数实现。以下代码取自dyld的源代码，演示了这个功能。

#if !defined(_DYLD_INTERPOSING_H_)
#define _DYLD_INTERPOSING_H_

#define DYLD_INTERPOSE(_replacment,_replacee) \ __attribute__((used)) static strut{const void* replacment;const void* replacee;}
_interpose_##_replace \ __attribute__((section ("__DATA,__interpose"))) = { (const void*)(unsigned long)&_replacement, (const void*)(unsigned long)&_replacee};

#endif
dyld的函数拦截功能提供了一个新的__DATA区，名为__interpose,在这个区中依次列出了替换函数和被替换的函数，其他事情就交由dyld处理了。
下面通过一个实验展示如何通过函数拦截机制跟踪malloc（）函数

#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdlib.h>
#include <malloc/malloc.h>

//标准的interpose数据结构
typedef struct interpose_s{
    void *new_func;
    void *orig_fnc;
}interpose_t;

//我们的原型
void *my_malloc(int size);//对应真实的malloc函数
void my_free(void *);//对应真实的free函数

static const interpose_t interposing_functions[] \
__attribute__ ((section ("__DATA,__interpose"))) = {{(void *)my_free,(void *)free},{(void *)my_malloc,(void *)malloc}};

void *my_malloc (int size){
    //在我们的函数中，要访问真正的malloc()函数，因为不想自己管理整个堆，所以就调用了原来的malloc()
    void *returned = malloc(size);
    //调用malloc_printf是因为printf中会调用malloc(),产生无限递归调用。
    malloc_printf("+ %p %d\n",returned,size);
    return (returned);
}

void my_free(void *freed){
    malloc_printf("- %p\n",freed);
    free(freed);
}

int main(int argc, const char * argv[]) {
    // 释放内存——打印出地址，然后调用真正的free()
    printf("Hello, World!\n");
    return 0;
}
在终端中执行以下代码编译为dylib并强制插入ls中

cc -dynamiclib 1.c -o libMTrace.dylib -Wall
DYLD_INSERT_LIBRARIES=libMTrace.dylib ls

最终会发现调用malloc的地方会打印这样的信息

ls(24346) malloc: + 0x100100020 88
ls(24346) malloc: + 0x100800000 4096
ls(24346) malloc: + 0x100801000 2160
ls(24346) malloc: - 0x100800000
ls(24346) malloc: + 0x100801a00 3312