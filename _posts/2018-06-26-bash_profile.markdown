---
title: bash_profile等几个文件的执行过程
layout: post
category: linux
author: 夏泽民
---
利用~/.bash_profile: 每个用户都可使用该文件输入专用于自己使用的shell信息,当用户登录时,该文件仅仅执行一次可以实现登录公告
$~/.bash_profile
if [ -f ～/.greeting ]; then . ～/.greeting; fi
$vi ～/.greeting
echo -e "^[[1m^[[4;37;32m公    告:^[[0m"
<!-- more -->
1. 在登录Linux时要执行文件的过程如下：

在刚登录Linux时，首先启动 /etc/profile 文件，然后再启动用户目录下的 ~/.bash_profile、 ~/.bash_login或 ~/.profile文件中的其中一个，

 

执行的顺序为：~/.bash_profile、 ~/.bash_login、 ~/.profile。

如果 ~/.bash_profile文件存在的话，一般还会执行 ~/.bashrc文件。

 

因为在 ~/.bash_profile文件中一般会有下面的代码：
if [ -f ~/.bashrc ] ; then

. ./bashrc

fi

~/.bashrc中，一般还会有以下代码：
if [ -f /etc/bashrc ] ; then

. /bashrc

fi

所以，~/.bashrc会调用 /etc/bashrc文件。最后，在退出shell时，还会执行 ~/.bash_logout文件。

执行顺序为：/etc/profile -> (~/.bash_profile | ~/.bash_login | ~/.profile) -> ~/.bashrc -> /etc/bashrc -> ~/.bash_logout

关于各个文件的作用域，在网上找到了以下说明：

（1）/etc/profile： 此文件为系统的每个用户设置环境信息,当用户第一次登录时,该文件被执行. 并从/etc/profile.d目录的配置文件中搜集shell的设置。

（2）/etc/bashrc: 为每一个运行bash shell的用户执行此文件.当bash shell被打开时,该文件被读取。

（3）~/.bash_profile: 每个用户都可使用该文件输入专用于自己使用的shell信息,当用户登录时,该文件仅仅执行一次!默认情况下,他设置一些环境变量,执行用户的.bashrc文件。

（4）~/.bashrc: 该文件包含专用于你的bash shell的bash信息,当登录时以及每次打开新的shell时,该该文件被读取。

（5）~/.bash_logout:当每次退出系统(退出bash shell)时,执行该文件. 另外,/etc/profile中设定的变量(全局)的可以作用于任何用户,而~/.bashrc等中设定的变量(局部)只能继承 /etc/profile中的变量,他们是"父子"关系。

（6）~/.bash_profile 是交互式、login 方式进入 bash 运行的~/.bashrc 是交互式 non-login 方式进入 bash 运行的通常二者设置大致相同，所以通常前者会调用后者。



2. /etc/profile和/etc/environment等各种环境变量设置文件的用处


case: 先将export LANG=zh_CN加入/etc/profile ,退出系统重新登录，登录提示显示英文。


将/etc/profile 中的export LANG=zh_CN删除，将LNAG=zh_CN加入/etc/environment，退出系统重新登录，登录提示显示中文。


用户环境建立的过程中总是先执行/etc/profile然后在读取/etc/environment。为什么会有如上所叙的不同呢？

应该是先执行/etc/environment，后执行/etc/profile。

/etc/environment是设置整个系统的环境，而/etc/profile是设置所有用户的环境，前者与登录用户无关，后者与登录用户有关。

系统应用程序的执行与用户环境可以是无关的，但与系统环境是相关的，所以当你登录时，你看到的提示信息，象日期、时间信息的显示格式与系统环境的LANG是相关的，缺省LANG=en_US，如果系统环境LANG=zh_CN，则提示信息是中文的，否则是英文的。


对于用户的SHELL初始化而言是先执行/etc/profile,再读取文件/etc/environment.对整个系统而言是先执行/etc/environment。这样理解正确吗？

/etc/enviroment --> /etc/profile --> $HOME/.profile -->$HOME/.env (如果存在)

/etc/profile 是所有用户的环境变量


/etc/enviroment是系统的环境变量


3. 登陆系统时shell读取的顺序应该是
/etc/profile ->/etc/enviroment -->$HOME/.profile -->$HOME/.env
原因应该shi用户环境和系统环境的区别了

如果同一个变量在用户环境(/etc/profile)和系统环境(/etc/environment)有不同的值那应该是以用户环境为准了。


