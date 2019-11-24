---
title: DYLD_FORCE_FLAT_NAMESPACE
layout: post
category: linux
author: 夏泽民
---
1、gcc生成dylib。

 

gcc -dynamiclib -o mysharedlib.dylib mysharedlib.c
 

2、gcc生成dylib，指定flatnamespace。

gcc -flat_namespace -dynamiclib -o openhook.dylib openhook.c
3、如何Hook？

dani-2:test leedani$ export DYLD_FORCE_FLAT_NAMESPACE=1
dani-2:test leedani$ export DYLD_INSERT_LIBRARIES=openhook.dylib
dani-2:test leedani$ ./main 
--------zz------hello,dani
<!-- more -->

4、Mac offers a way to override functions in a shared library with DYLD_INSERT_LIBRARIES environment variable (which is similar to LD_PRELOAD on Linux). When you make a twin brother of a function that is defined in an existing shared library, put it in you a shared library, and you register your shared library name in DYLD_INSERT_LIBRARIES, your function is used instead of the original one. This is my simple test. Here I’ve replaced f() in mysharedlib.dylib with f() in openhook.dylib.

5、关于DYLD_INSERT_LIBRARIES & DYLD_FORCE_FLAT_NAMESPACE

参考：

1、http://www.h4ck.org.cn/2013/04/hooking-library-calls-on-mac-using-dyld_insert_libraries/

2、http://blog.sina.com.cn/s/blog_45e2b66c0101cde0.html


Mac可以通过设置DYLD_INSERT_LIBRARIES环境变量（linux上对应的环境变量是LD_PRELOAD ，效果实例可见 Android hook——LD_PRELOAD），重写动态链接库中的函数，实现hook功能。


以下是演示实例



一、替换动态链接库中的c函数

实例一：使用openhook.dylib中的f() 替换原始动态链接库mysharedlib.dylib中的f() 。实例来源



（一）、源文件

1. mysharedlib.h

 


void f();
2. mysharedlib.c
 

 


#include
#include "mysharedlib.h"

void f(){
printf("hello,dani \n");
}
3. main.c

 


#include
#include "mysharedlib.h"

int main(){
f();
return 0;
}
 

4. openhook.c
 

 


#include
#include
#include
#include "mysharedlib.h"

typedef void (*fType) ();
static void (*real_f)()=NULL;

void f(){
if (! real_f){
void * handle = dlopen("mysharedlib.dylib", RTLD_NOW);
real_f = (fType) dlsym(handle,"f");

if(! real_f) printf("NG");
}
printf("--------zz------");

real_f();
}
关键函数：
dlopen函数原型，void * dlopen( const char * pathname, int mode)，pathname是指定动态链接库地址，mode是打开模式
dlsym函数原型，void* dlsym(void* handle,const char* symbol)，handle是由dlopen打开动态链接库后返回的指针，symbol是指定获取的符号名，对c语言而言，符号名就是函数名，我们可以使用nm查看mysharedlib.dylib

dani-2:testC leedani$ nm mysharedlib.dylib
0000000000000f20 T _f
U _puts
U dyld_stub_binder

（二）、编译
1. 生成mysharedlib.dylib, 该动态链接库的功能就是f(),打印“hello，dani”
 


gcc -dynamiclib -o mysharedlib.dylib mysharedlib.c
dynamiclib选项是指生成动态链接库

2. 编译mysharedlib.dylib与main.c文件，生成最终的可执行文件
 

 


gcc mysharedlib.dylib main.c -o main
3. 生成openhook.dylib，该动态链接库的功能就是替换mysharedlib.dylib中的f()
 

 


gcc -flat_namespace -dynamiclib -o openhook.dylib openhook.c
flat_namespace选项指定了链接模式，有两种模式,flat-namespace与two-level
namespace,模式不一样生成的符号表也会不一样（具体区别）。

实例中mysharedlib.dylib没有采用该选项，而openhook.dylib采用了该选项，我们可以查看以下这两个文件的头结构，来对比一下

 


dani-2:testC leedani$ otool -hV mysharedlib.dylib
mysharedlib.dylib:
Mach header
magic cputype cpusubtype caps filetype ncmds sizeofcmds flags
MH_MAGIC_64 X86_64 ALL 0x00 DYLIB 13 1200 NOUNDEFS DYLDLINK TWOLEVEL NO_REEXPORTED_DYLIBS

dani-2:testC leedani$ otool -hV openhook.dylib
openhook.dylib:
Mach header
magic cputype cpusubtype caps filetype ncmds sizeofcmds flags
MH_MAGIC_64 X86_64 ALL 0x00 DYLIB 13 1272 DYLDLINK NO_REEXPORTED_DYLIBS
 

（三）、运行
 

1. 正常的运行结果

 


dani-2:test leedani$ ./main
hello,dani
2. Hook后的运行结果
通过设置环境变量DYLD_INSERT_LIBRARIES（linux上对应的环境变量是LD_PRELOAD ，效果实例可见 Android hook——LD_PRELOAD）
 


dani-2:test leedani$ export DYLD_FORCE_FLAT_NAMESPACE=1
dani-2:test leedani$ export DYLD_INSERT_LIBRARIES=openhook.dylib
dani-2:test leedani$ ./main
--------zz------hello,dani
DYLD_INSERT_LIBRARIES与DYLD_FORCE_FLAT_NAMESPACE环境变量在apple官方手册中有说明，如下所示：

Mac hook——DYLD_INSERT_LIBRARIES - 碳基体 - 碳基体
 
实例二：替换系统动态链接库中的函数，如下所示替换/usr/lib/libSystem.dylib中的time函数
实例来源
（一）、源码
time .c

#include

//This function will override the one in /usr/lib/libSystem.dylib

time_t time(time_t *tloc){
//January 1st,2013
struct tm timeStruct;
timeStruct.tm_year= 2013-1900;
timeStruct.tm_mon = 0;
timeStruct.tm_mday = 1;
timeStruct.tm_hour = 0;
timeStruct.tm_min = 0;
timeStruct.tm_sec = 0;
timeStruct.tm_isdst = -1;

*tloc = mktime(&timeStruct);

return *tloc;
}

（二）、编译

gcc -flat_namespace -dynamiclib -current_version 1.0 time.o -o libTime.dylib

（三）、运行
1. 正常的运行结果

dani-2:test leedani$ date
2013年 2月 1日 星期五 14时46分16秒 CST
2.替换系统函数后的运行结果

dani-2:test leedani$ export DYLD_FORCE_FLAT_NAMESPACE=1
dani-2:test leedani$ export DYLD_INSERT_LIBRARIES=libTime.dylib
dani-2:test leedani$ date
2013年 1月 1日 星期二 00时00分00秒 CST

 

二、替换动态链接库中的c++ 类方法
 

实例来源
（一）、源码
1. mysharedlib.h

class AAA
{
public:
int m;

AAA()
{
m = 1234;
}

void fff(int a);
};
2. mysharedlib.cpp

#include
#include "mysharedlib.h"

void AAA::fff(int a)
{

printf("-- Original: %d --", a);

}
3. main.cpp

#include
#include "mysharedlib.h"

int main()
{

AAA a;

printf("---------main1-------\n");

a.fff(50);

printf("\n---------main2-------\n");

return 0;
}
4. openhook.cpp
#include
#include
#include
#include "mysharedlib.h"

typedef void (*AAAfffType)(AAA*,int);
static void (*real_AAAfff)(AAA*,int);

extern "C"
{

void _ZN3AAA3fffEi(AAA* a, int b)
{

printf("---------AAA::fff------\n");
printf("%d,%d \n",b,a->m);

void * handle = dlopen("mysharedlib.dylib", RTLD_NOW);

real_AAAfff = (AAAfffType)dlsym(handle, "_ZN3AAA3fffEi");

if(real_AAAfff) printf("OK");

real_AAAfff(a,b);
}
}
关键函数：
dlopen函数原型，void * dlopen( const char * pathname, int mode)，pathname是指定动态链接库地址，mode是打开模式
dlsym函数原型，void* dlsym(void* handle,const char* symbol)，handle是由dlopen打开动态链接库后返回的指针，symbol是指定获取的符号名，对c++语言而言，由于存在name mangling，符号名不再是函数名了，编译器不同生成的符号名也会有所区别，我们可以使用nm查看mysharedlib.dylib

dani-2:testCPP leedani$ nm mysharedlib.dylib
0000000000000f0c T __ZN3AAA3fffEi
U _printf
U dyld_stub_binder
使用关键字extern "C"是为了防止符号名被mangle，使其可以像c一样被dlsym加载，具体的如何在unix环境下使用dlopen 动态加载c＋＋类函数可以看这篇文章《c++ dlopen mini HOWTO》

（二）、编译
1. 生成mysharedlib.dylib, 该动态链接库的功能就是f(),打印“hello，dani”
 


gcc -dynamiclib -lstdc++ -o mysharedlib.dylib mysharedlib.cpp
2. 编译mysharedlib.dylib与main.c文件，生成最终的可执行文件
 

 


gcc -lstdc++ mysharedlib.dylib main.cpp -o main
3. 生成openhook.dylib，该动态链接库的功能就是替换mysharedlib.dylib中的f()
 

 


gcc -flat_namespace -dynamiclib -lstdc++ -o openhook.dylib openhook.cpp
（三）、运行
1. 正常运行

dani-2:testCPP leedani$ ./main
---------main1-------
-- Original: 50 --
---------main2-------
2. hook后的结果，通过设置环境变量DYLD_INSERT_LIBRARIES
 

 




dani-2:test leedani$ export DYLD_FORCE_FLAT_NAMESPACE=1
dani-2:test leedani$ export DYLD_INSERT_LIBRARIES=openhook.dylib




dani-2:testCPP leedani$ ./main
---------main1-------
---------AAA::fff------
50,1234
OK
-- Original: 50 --
---------main2-------


三、小结
这种通过设置环境变量DYLD_INSERT_LIBRARIES，动态加载函数、类方法来实现使用自己编写的动态连接库dylib来patch运行中的应用的手段，是外挂、MobileSubstrate插件的主要原理，推广到PC windows平台（dll hook），Android平台（linux平台）（so hook），iOS平台（mac平台）（dylib hook），可以说动态加载技术奠定了软件patch的基础，需要深入了解。

 参考：
http://koichitamura.blogspot.com/2008/11/hooking-library-calls-on-mac.html
http://hactheplanet.com/blog/80
https://developer.apple.com/library/mac/#documentation/Darwin/Reference/Manpages/man1/dyld.1.html
https://developer.apple.com/library/mac/#documentation/Darwin/Reference/ManPages/man3/dlopen.3.html
https://developer.apple.com/library/mac/#documentation/Darwin/Reference/ManPages/man3/dlsym.3.html
https://developer.apple.com/library/mac/#documentation/developertools/conceptual/MachOTopics/1-Articles/executing_files.html

