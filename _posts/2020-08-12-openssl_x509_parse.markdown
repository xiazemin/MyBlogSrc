---
title: openssl_x509_parse
layout: post
category: php
author: 夏泽民
---
https://www.php.net/manual/en/function.openssl-x509-parse.php

https://www.seebug.org/vuldb/ssvid-61173

  [Symfony\Component\Debug\Exception\FatalThrowableError]
  Call to undefined function Composer\CaBundle\openssl_x509_parse()
  
  这个函数依赖openssl扩展
<!-- more -->

PHP Warning:  PHP Startup: Unable to load dynamic library '/usr/local/lib/php/extensions/no-debug-non-zts-20151012/openssl.so' - dlopen(/usr/local/lib/php/extensions/no-debug-non-zts-20151012/openssl.so, 9): initializer function 0x7fff9ac360a0 not in mapped image for /usr/local/lib/php/extensions/no-debug-non-zts-20151012/openssl.so
 in Unknown on line 0
 
 版本冲突问题
 
  ./Configure darwin64-x86_64-cc -shared --prefix=/Users/xcl/Desktop/openssl/myopenssl --openssldir=/usr/local/ssl
  
  https://blog.csdn.net/xcltapestry/article/details/48580735
  
  
  
  编译OpenSSL:
./Configure darwin64-x86_64-cc --prefix=/Users/xcl/Desktop/openssl/myopenssl --openssldir=/usr/local/ssl
make
会编译出静态库:
 -- libssl.a libcrypto.a
 XCLiMac:openssl-0.9.8zg xcl$ ls lib*
libcrypto.a	libcrypto.pc	libssl.a	libssl.pc

编译动态库加上shared参数即可
-- 生成 libssl.dylib  libcrypto.dylib
 ./Configure darwin64-x86_64-cc -shared --prefix=/Users/xcl/Desktop/openssl/myopenssl --openssldir=/usr/local/ssl
 XCLiMac:openssl-0.9.8zg xcl$ openssl version -a
OpenSSL 1.0.2a 19 Mar 2015
built on: reproducible build, date unspecified
platform: darwin64-x86_64-cc
options:  bn(64,64) rc4(ptr,int) des(idx,cisc,16,int) idea(int) blowfish(idx) 
compiler: /usr/bin/clang -I. -I.. -I../include  -fPIC -fno-common -DOPENSSL_PIC -DZLIB -DOPENSSL_THREADS -D_REENTRANT -DDSO_DLFCN -DHAVE_DLFCN_H -arch x86_64 -O3 -DL_ENDIAN -Wall -DOPENSSL_IA32_SSE2 -DOPENSSL_BN_ASM_MONT -DOPENSSL_BN_ASM_MONT5 -DOPENSSL_BN_ASM_GF2m -DSHA1_ASM -DSHA256_ASM -DSHA512_ASM -DMD5_ASM -DAES_ASM -DVPAES_ASM -DBSAES_ASM -DWHIRLPOOL_ASM -DGHASH_ASM -DECP_NISTZ256_ASM


https://blog.csdn.net/xcltapestry/article/details/48580735

链接最新的openssl版本
我们可以直接用一句指令

$ brew link openssl --force
执行过后，重新打开终端，输入openssl version，即可看到就是新的版本了。

https://www.jianshu.com/p/1cad918f34d6
https://www.jianshu.com/p/32f068922baf



brew是mac机上面程序猿非常常用的软件包安装方式，其中有两组命令是需要大家知晓的。分别是：

第一组：brew install和brew uninstall。
第二组，brew link和brew unlink。
不过关于第一组brew install命令，比较常用，所以大家可能会比较熟悉。后面的这组brew link命令才是本文要讲述的重点。苏南大叔将以前不久刚刚降级安装的php71为例，说明一下brew link命令。



正常情况下来说，brew link php71并不是需要主动执行的，因为在brew install php71的过程中，就已经默认执行了brew link php71。但是，由于各种各样的权限之类的问题，导致brew link php71操作是失败的。在brew install php71的过程中，就会体现为一个警告信息。

而对于实际的应用上来说，可能表现为：不能识别php命令，不能识别phpize命令，或者不能识别php-config命令。这些问题实际上是很fatal的，会导致一系列的后续错误。比如安装扩展插件识别，或者编译扩展插件失败，composer命令不能使用等问题

http://caibaojian.com/a-programmer/software/mac/softwares/brew.html



The symlinks can be seen with ls as normal links. ls -lh /usr/local/bin/python => /usr/local/bin/python -> ../Cellar/python/3.6.4_3/bin/python. For a complete reference of all the symlinks homebrew manages I am curious too. Cellar is simply where all the Homebrew packages reside. It is under /usr/local/Cellar

https://stackoverflow.com/questions/33268389/what-does-brew-link-do


ln 创建软链接时覆盖

创建软链接：ln -s /path/to/file link-file-name

创建软链接并覆盖： ln -sfn /path/to/file link-file-name

-f 强制覆盖原文件

-n 覆盖目录link文件，若不加则会在link目录中创建链接
https://www.cnblogs.com/flashBoxer/p/9790509.html

ln建立时符号链接时出现同名文件或目录
给ln命令加上-s选项，则建立软链接。

格式：ln -s [真正的文件或者目录] [链接名]

 

[链接名]可以是任何一个文件名或者目录名，并且允许它与原文件不在同一个文件系统中。

如果[链接名]是一个已经存在的文件，将不做链接。

如果[链接名]是一个已经存在的目录，linux系统会分两种情况自行进行处理：

若链接指向的是一个文件名，系统将在已经存在的目录下建立一个与源文件名同名的符号链接文件

若链接指向的是一个目录名，系统将在已经存在的目录下建立一个与源目录名同名的符号链接文件

总之，建立软链接就是建立了一个新文件。当访问链接文件时，系统就会发现它是个链接文件，系统读取链接文件找到真正要访问的文件然后打开。


这是linux中一个非常重要命令，请大家一定要熟悉。它的功能是为某一个文件在另外一个位置建立一个同不的链接，这个命令最常用的参数是-s,具体用法是：ln -s 源文件 目标文件
 
这是linux中一个非常重要命令，请大家一定要熟悉。它的功能是为某一个文件在另外一个位置建立一个同不的链接，这个命令最常用的参数是-s,具体用法是：ln -s 源文件 目标文件。 
当 我们需要在不同的目录，用到相同的文件时，我们不需要在每一个需要的目录下都放一个必须相同的文件，我们只要在某个固定的目录，放上该文件，然后在其它的 目录下用ln命令链接（link）它就可以，不必重复的占用磁盘空间。

例如：ln -s /bin/less /usr/local/bin/less 

-s 是代号（symbolic）的意思。 
这 里有两点要注意：第一，ln命令会保持每一处链接文件的同步性，也就是说，不论你改动了哪一处，其它的文件都会发生相同的变化；第二，ln的链接又软链接 和硬链接两种，软链接就是ln -s ** **,它只会在你选定的位置上生成一个文件的镜像，不会占用磁盘空间，硬链接ln ** **,没有参数-s, 它会在你选定的位置上生成一个和源文件大小相同的文件，无论是软链接还是硬链接，文件都保持同步变化。 
如果你用ls察看一个目录时，发现有的文件后面有一个@的符号，那就是一个用ln命令生成的文件，用ls -l命令去察看，就可以看到显示的link的路径了。 

ln是linux中又一个非常重要命令，它的功能是为某一个文件在另外一个位置建立一个同步的链接.当我们需要在不同的目录，用到相同的文件时，我们不需要在每一个需要的目录下都放一个必须相同的文件，我们只要在某个固定的目录，放上该文件，然后在 其它的目录下用ln命令链接（link）它就可以，不必重复的占用磁盘空间。

1．命令格式：

 ln [参数][源文件或目录][目标文件或目录]

2．命令功能：

Linux文件系统中，有所谓的链接(link)，我们可以将其视为档案的别名，而链接又可分为两种 : 硬链接(hard link)与软链接(symbolic link)，硬链接的意思是一个档案可以有多个名称，而软链接的方式则是产生一个特殊的档案，该档案的内容是指向另一个档案的位置。硬链接是存在同一个文件系统中，而软链接却可以跨越不同的文件系统。

软链接：

1.软链接，以路径的形式存在。类似于Windows操作系统中的快捷方式
2.软链接可以 跨文件系统 ，硬链接不可以
3.软链接可以对一个不存在的文件名进行链接
4.软链接可以对目录进行链接

硬链接:

1.硬链接，以文件副本的形式存在。但不占用实际空间。
2.不允许给目录创建硬链接
3.硬链接只有在同一个文件系统中才能创建

  这里有两点要注意：

第一，ln命令会保持每一处链接文件的同步性，也就是说，不论你改动了哪一处，其它的文件都会发生相同的变化；
第二，ln的链接又分软链接和硬链接两种，软链接就是ln –s 源文件 目标文件，它只会在你选定的位置上生成一个文件的镜像，不会占用磁盘空间，硬链接 ln 源文件 目标文件，没有参数-s， 它会在你选定的位置上生成一个和源文件大小相同的文件，无论是软链接还是硬链接，文件都保持同步变化。

ln指令用在链接文件或目录，如同时指定两个以上的文件或目录，且最后的目的地是一个已经存在的目录，则会把前面指定的所有文件或目录复制到该目录中。若同时指定多个文件或目录，且最后的目的地并非是一个已存在的目录，则会出现错误信息。


3．命令参数：

必要参数:

-b 删除，覆盖以前建立的链接
-d 允许超级用户制作目录的硬链接
-f 强制执行
-i 交互模式，文件存在则提示用户是否覆盖
-n 把符号链接视为一般目录
-s 软链接(符号链接)
-v 显示详细的处理过程


选择参数:

-S “-S<字尾备份字符串> ”或 “--suffix=<字尾备份字符串>”

-V “-V<备份方式>”或“--version-control=<备份方式>”

--help 显示帮助信息

--version 显示版本信息

https://www.cnblogs.com/zhuyeshen/p/11693406.html
https://blog.csdn.net/guojin08/article/details/38702919/

https://blog.csdn.net/m0_37450089/article/details/80297361

https://www.cnblogs.com/gaoBlog/p/12264197.html


$sudo ln -sfn /usr/local/etc/openssl\@1.0/bin/openssl /usr/bin/openssl
Password:
ln: /usr/bin/openssl: Operation not permitted

原因

这是因为苹果在OS X 10.11中引入的SIP特性使得即使加了sudo（也就是具有root权限）也无法修改系统级的目录，其中就包括了/usr/bin。要解决这个问题有两种做法：

一种是比较不安全的就是关闭SIP，也就是rootless特性；

另一种是将本要链接到/usr/bin下的改链接到/usr/local/bin下就好了。

三、解决办法

sudo ln -s /usr/local/mysql/bin/mysql /usr/local/bin

$sudo ln -sfn /usr/local/etc/openssl\@1.0/bin/openssl /usr/local/bin/openssl
15:34:10-didi@localhost:~/PhpstormProjects/c/php-src/ext/openssl$which openssl
/usr/local/bin/openssl

$ln -snf /usr/local/etc/openssl\@1.0/ /usr/local/Cellar/openssl\@1.0



./config --prefix=/usr/local/openssl

https://blog.csdn.net/focusjava/article/details/51179297


$./Configure darwin64-x86_64-cc -shared --prefix=/usr/local/openssl

$./Configure darwin64-x86_64-cc -shared --prefix=/usr/local

  557  sudo make -j4
  558
  559   sudo make install
  560  which openssl
  
 $which openssl
/usr/local/bin/openssl

 563  phpize
  564  make clean
  565  ./configure
  566  make -je
  567  make -j4
  568  make install
  
  问题解决
  
  总结，mac 默认在/usr/local/目录查找bin和lib，如果用户安装软件，推荐安装在／usr／local/这个目录，不会出现找不到的问题
  
  如果／usr／local
  目录没有，就会去／usr目录去查找预先安装的版本
  会出现版本冲突
  brew 安装的原理是，安装在／usr/local/Celler目录
  
  然后软连接到/usr/local目录
  
  which 命令和非交互式shell都是如此的
  
