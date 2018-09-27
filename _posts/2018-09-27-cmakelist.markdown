---
title: cmakelist
layout: post
category: algorithm
author: 夏泽民
---
1.CMake编译原理

CMake是一种跨平台编译工具，比make更为高级，使用起来要方便得多。CMake主要是编写CMakeLists.txt文件，然后用cmake命令将CMakeLists.txt文件转化为make所需要的makefile文件，最后用make命令编译源码生成可执行程序或共享库（so(shared object)）。因此CMake的编译基本就两个步骤：

1. cmake
2. make
cmake  指向CMakeLists.txt所在的目录，例如cmake .. 表示CMakeLists.txt在当前目录的上一级目录。cmake后会生成很多编译的中间文件以及makefile文件，所以一般建议新建一个新的目录，专门用来编译，例如

mkdir build
cd build
cmake ..
make
make根据生成makefile文件，编译程序。

2.使用Cmake编译程序

我们编写一个关于开平方的C/C++程序项目，即b= sqrt(a)，以此理解整个CMake编译的过程。

a.准备程序文件

文件目录结构如下：

复制代码
.
├── build
├── CMakeLists.txt
├── include
│   └── b.h
└── src
    ├── b.c
    └── main.c
复制代码
头文件b.h，如下所示：

复制代码
#ifndef B_FILE_HEADER_INC
#define B_FIEL_HEADER_INC

#include<math.h>

double cal_sqrt(double value);

#endif
复制代码
 

头文件b.c，如下所示：

复制代码
#include "../include/b.h"

double cal_sqrt(double value)
{
    return sqrt(value);
}
复制代码
 

main.c主函数，如下所示：

复制代码
#include "../include/b.h"
#include <stdio.h>
int main(int argc, char** argv)
{
    double a = 49.0; 
    double b = 0.0;

    printf("input a:%f\n",a);
    b = cal_sqrt(a);
    printf("sqrt result:%f\n",b);
    return 0;
}
复制代码
 

b.编写CMakeLists.txt

接下来编写CMakeLists.txt文件，该文件放在和src，include的同级目录，实际方哪里都可以，只要里面编写的路径能够正确指向就好了。CMakeLists.txt文件，如下所示：

复制代码
 1 #1.cmake verson，指定cmake版本 
 2 cmake_minimum_required(VERSION 3.2)
 3 
 4 #2.project name，指定项目的名称，一般和项目的文件夹名称对应
 5 PROJECT(test_sqrt)
 6 
 7 #3.head file path，头文件目录
 8 INCLUDE_DIRECTORIES(
 9 include
10 )
11 
12 #4.source directory，源文件目录
13 AUX_SOURCE_DIRECTORY(src DIR_SRCS)
14 
15 #5.set environment variable，设置环境变量，编译用到的源文件全部都要放到这里，否则编译能够通过，但是执行的时候会出现各种问题，比如"symbol lookup error xxxxx , undefined symbol"
16 SET(TEST_MATH
17 ${DIR_SRCS}
18 )
19 
20 #6.add executable file，添加要编译的可执行文件
21 ADD_EXECUTABLE(${PROJECT_NAME} ${TEST_MATH})
22 
23 #7.add link library，添加可执行文件所需要的库，比如我们用到了libm.so（命名规则：lib+name+.so），就添加该库的名称
24 TARGET_LINK_LIBRARIES(${PROJECT_NAME} m)
复制代码
 CMakeLists.txt主要包含以上的7个步骤，具体的意义，请阅读相应的注释。

c.编译和运行程序

准备好了以上的所有材料，接下来，就可以编译了，由于编译中出现许多中间的文件，因此最好新建一个独立的目录build，在该目录下进行编译，编译步骤如下所示：

mkdir build
cd build
cmake ..
make
操作后，在build下生成的目录结构如下：

复制代码
├── build
│   ├── CMakeCache.txt
│   ├── CMakeFiles
│   │   ├── 3.2.2
│   │   │   ├── CMakeCCompiler.cmake
│   │   │   ├── CMakeCXXCompiler.cmake