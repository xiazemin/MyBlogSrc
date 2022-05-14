---
title: configure.ac
layout: post
category: linux
author: 夏泽民
---
编译源码目录下只有configure.ac文件和Makefile.am文件的工程
<!-- more -->
aclocal

autoconf

autoheader

automake --add-missing

./configure

make

sudo make install

<img src="{{site.url}}{{site.baseurl}}/img/configure_auto.jpg"/>
	
autoscan: 扫描源代码以搜寻普通的可移植性问题，比如检查编译器，库，头文件等，生成文件configure.scan,它是configure.ac的一个雏形。

aclocal:根据已经安装的宏，用户定义宏和acinclude.m4文件中的宏将configure.ac文件所需要的宏集中定义到文件 aclocal.m4中。aclocal是一个perl 脚本程序，它的定义是：“aclocal - create aclocal.m4 by scanning configure.ac”

automake:将Makefile.am中定义的结构建立Makefile.in，然后configure脚本将生成的Makefile.in文件转换 为Makefile。如果在configure.ac中定义了一些特殊的宏，比如AC_PROG_LIBTOOL，它会调用libtoolize，否则它 会自己产生config.guess和config.sub

autoconf:将configure.ac中的宏展开，生成configure脚本。这个过程可能要用到aclocal.m4中定义的宏。

三、实例
1.测试代码(定义两个文件hello.h和hello.c)
复制代码
/*hello.c*/
#include <iostream>
#include "hello.h"

using namespace std;

int main()
{
    CHello a;
    return 0;
}
复制代码
复制代码
/*hello.h*/
#ifndef __HELLO_H__
#define __HELLO_H__

#include<iostream>
using namespace std;

class CHello
{
public:
    CHello(){ cout<<"Hello!"<<endl;}
    ~CHello(){ cout<<"Bye!"<<endl;}
};

#endif
复制代码
2.操作步骤
(1)安装依赖的包
[root@bogon autoconfig]# yum -y install automake autoconf
automake包括:aclocal、automake等

autoconf包括：autoscan、autoconf等

(2)autoscan
复制代码
[root@bogon autoconfig]# ll
-rw-r--r-- 1 root root 105 Jun  4 hello.cpp
-rw-r--r-- 1 root root 189 Jun  4 hello.h
[root@bogon autoconfig]# autoscan
[root@bogon autoconfig]# ll
total 12
-rw-r--r-- 1 root root   0 Jun  4 autoscan.log
-rw-r--r-- 1 root root 481 Jun  4 configure.scan
-rw-r--r-- 1 root root 105 Jun  4 hello.cpp
-rw-r--r-- 1 root root 189 Jun  4 hello.h
复制代码
(3)aclocal
复制代码
[root@bogon autoconfig]# mv configure.scan configure.ac
[root@bogon autoconfig]# vim configure.ac  /*将下面红色加粗的部分修改掉*/
#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.
AC_PREREQ([2.69])
AC_INIT(hello, 1.0, admin@163.com)
AM_INIT_AUTOMAKE(hello, 1.0)
AC_CONFIG_SRCDIR([hello.cpp])
AC_CONFIG_HEADERS([config.h])

# Checks for programs.
AC_PROG_CXX
AC_PROG_CC

# Checks for libraries.
# Checks for header files.
# Checks for typedefs, structures, and compiler characteristics.
# Checks for library functions.
AC_OUTPUT(Makefile)
[root@bogon autoconfig]# aclocal
[root@bogon autoconfig]# ll
total 52
-rw-r--r-- 1 root root 37794 Jun  4 aclocal.m4
drwxr-xr-x 2 root root    51 Jun  4 autom4te.cache
-rw-r--r-- 1 root root     0 Jun  4 autoscan.log
-rw-r--r-- 1 root root   492 Jun  4 configure.ac
-rw-r--r-- 1 root root   105 Jun  4 hello.cpp
-rw-r--r-- 1 root root   189 Jun  4 hello.h
复制代码
下面给出本文件的简要说明（所有以”#”号开始的行为注释）：
· AC_PREREQ 宏声明本文件要求的autoconf版本，本例使用的版本为2.59。
· AC_INIT 宏用来定义软件的名称和版本等信息，”FULL-PACKAGE-NAME”为软件包名称，”VERSION”为软件版本号，”BUG-REPORT-ADDRESS”为BUG报告地址（一般为软件作者邮件地址）。
·AC_CONFIG_SRCDIR 宏用来侦测所指定的源码文件是否存在，来确定源码目录的有效性。此处为当前目录下的hello.c。
·AC_CONFIG_HEADER 宏用于生成config.h文件，以便autoheader使用。
·AC_PROG_CC 用来指定编译器，如果不指定，选用默认gcc。
·AC_OUTPUT 用来设定 configure 所要产生的文件，如果是makefile，configure会把它检查出来的结果带入makefile.in文件产生合适的makefile。使用Automake时，还需要一些其他的参数，这些额外的宏用aclocal工具产生。

(3)autoconf
复制代码
[root@bogon autoconfig]# autoconf
[root@bogon autoconfig]# ll
total 204
-rw-r--r-- 1 root root  37794 Jun  4 aclocal.m4
drwxr-xr-x 2 root root     81 Jun  4 autom4te.cache
-rw-r--r-- 1 root root      0 Jun  4 autoscan.log
-rwxr-xr-x 1 root root 154727 Jun  4 configure
-rw-r--r-- 1 root root    492 Jun  4 configure.ac
-rw-r--r-- 1 root root    105 Jun  4 hello.cpp
-rw-r--r-- 1 root root    189 Jun  4 hello.h
复制代码
此时可以看到已经生成了configure

(4)autoheader
复制代码
[root@bogon autoconfig]# autoheader
[root@bogon autoconfig]# ll
total 208
-rw-r--r-- 1 root root  37794 Jun  4 aclocal.m4
drwxr-xr-x 2 root root     81 Jun  4 autom4te.cache
-rw-r--r-- 1 root root      0 Jun  4 autoscan.log
-rw-r--r-- 1 root root    625 Jun  4 config.h.in
-rwxr-xr-x 1 root root 154727 Jun  4 configure
-rw-r--r-- 1 root root    492 Jun  4 configure.ac
-rw-r--r-- 1 root root    105 Jun  4 hello.cpp
-rw-r--r-- 1 root root    189 Jun  4 hello.h
复制代码
autoheader生成了configure.h.in如果在configure.ac中定义了AC_CONFIG_HEADER，那么此文件就需要；

(5)Makefile.am
[root@bogon autoconfig]# vim Makefile.am
[root@bogon autoconfig]# cat Makefile.am 
AUTOMAKE_OPTIONS=foreign 
bin_PROGRAMS=hello 
hello_SOURCES=hello.cpp hello.h
· AUTOMAKE_OPTIONS 为设置Automake的选项。由于GNU对自己发布的软件有严格的规范，比如必须附带许可证声明文件COPYING等，否则Automake执行时会报错。Automake提供了3种软件等级：foreign、gnu和gnits，供用户选择，默认等级为gnu。本例使需用foreign等级，它只检测必须的文件。
· bin_PROGRAMS 定义要产生的执行文件名。如果要产生多个执行文件，每个文件名用空格隔开。
· hello_SOURCES 定义”hello”这个执行程序所需要的原始文件。如果”hello”这个程序是由多个原始文件所产生的，则必须把它所用到的所有原始文件都列出来，并用空格隔开。例如：若目标体”hello”需要”hello.c”、”hello.h”两个依赖文件，则定义hello_SOURCES=hello.c hello.h。

(6)automake
复制代码
[root@bogon autoconfig]# automake --add-missing
configure.ac:6: warning: AM_INIT_AUTOMAKE: two- and three-arguments forms are deprecated.  For more info, see:
configure.ac:6: http://www.gnu.org/software/automake/manual/automake.html#Modernize-AM_005fINIT_005fAUTOMAKE-invocation
configure.ac:6: installing './install-sh'
configure.ac:6: installing './missing'
[root@bogon autoconfig]# ll
total 236
-rw-r--r-- 1 root root  37794 Jun  4  aclocal.m4
drwxr-xr-x 2 root root     81 Jun  4  autom4te.cache
-rw-r--r-- 1 root root      0 Jun  4  autoscan.log
-rw-r--r-- 1 root root    625 Jun  4  config.h.in
-rwxr-xr-x 1 root root 154727 Jun  4  configure
-rw-r--r-- 1 root root    492 Jun  4  configure.ac
-rw-r--r-- 1 root root    105 Jun  4  hello.cpp
-rw-r--r-- 1 root root    189 Jun  4  hello.h
lrwxrwxrwx 1 root root     35 Jun  4  install-sh -> /usr/share/automake-1.13/install-sh
-rw-r--r-- 1 root root     79 Jun  4  Makefile.am
-rw-r--r-- 1 root root  22227 Jun  4  Makefile.in
lrwxrwxrwx 1 root root     32 Jun  4  missing -> /usr/share/automake-1.13/missing
复制代码
此步主要是为了生成Makefile.in，加上--add-missing参数后，会补全缺少的脚本；

(6)测试
[root@bogon autoconfig]# ./configure
[root@bogon autoconfig]# make
[root@bogon autoconfig]# ./hello
Hello!
Bye!
和平时安装许多开源软件一样操作

(7)打包
复制代码
[root@bogon autoconfig]# make dist
[root@bogon autoconfig]# ll
total 436
-rw-r--r-- 1 root root  37794 Jun  4 aclocal.m4
drwxr-xr-x 2 root root     81 Jun  4 autom4te.cache
-rw-r--r-- 1 root root      0 Jun  4 autoscan.log
-rw-r--r-- 1 root root    758 Jun  4 config.h
-rw-r--r-- 1 root root    625 Jun  4 config.h.in
-rw-r--r-- 1 root root  11031 Jun  4 config.log
-rwxr-xr-x 1 root root  32557 Jun  4 config.status
-rwxr-xr-x 1 root root 154727 Jun  4 configure
-rw-r--r-- 1 root root    492 Jun  4 configure.ac
lrwxrwxrwx 1 root root     32 Jun  4 depcomp -> /usr/share/automake-1.13/depcomp
-rwxr-xr-x 1 root root  22250 Jun  4 hello
-rw-r--r-- 1 root root  72021 Jun  4 hello-1.0.tar.gz
-rw-r--r-- 1 root root    105 Jun  4 hello.cpp
-rw-r--r-- 1 root root    189 Jun  4 hello.h
-rw-r--r-- 1 root root  26008 Jun  4 hello.o
lrwxrwxrwx 1 root root     35 Jun  4 install-sh -> /usr/share/automake-1.13/install-sh
-rw-r--r-- 1 root root  23564 Jun  4 Makefile
-rw-r--r-- 1 root root     79 Jun  4 Makefile.am
-rw-r--r-- 1 root root  23869 Jun  4 Makefile.in
lrwxrwxrwx 1 root root     32 Jun  4 missing -> /usr/share/automake-1.13/missing
-rw-r--r-- 1 root root     23 Jun  4 stamp-h1
复制代码
如果细心的话可以发现，人群中已经出现了hello-1.0.tar.gz就是将已经编译好的文件进行了打包。