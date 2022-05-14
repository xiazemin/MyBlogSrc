---
title: expand_aliases Bash非交互模式下alias不能使用
layout: post
category: linux
author: 夏泽民
---
在Bash的非交互模式下(一般的脚本), alias是无效的.

例如, 我在.bashrc设置了:

alias code="/Applications/Visual\ Studio\ Code.app/Contents/Resources/app/bin/code"
export EDITOR=code
但是在脚本执行$EDITOR (默认编辑器) 时, 报错code 命令找不到. 即使在脚本里面source ~/.bashrc亦无效. 这是因为, code是使用了alias来实现的.

在非交互式模式下alias扩展功能默认是关闭的，此时仍然可以定义alias别名，但是shell不会将alias别名扩展成对应的命令，而是将alias别名本身当作命令执行，如果shell内置命令和PATH中均没有与alias别名同名的命令，则shell会“抱怨”找不到指定的命令。

如果想使得非交互模式下支持alias扩展, 就要使用shopt expand_aliases. 注意shopt这个命令是shell内置命令，可以控制shell功能选项的开启和关闭，从而控制shell的行为。shopt的使用方式如下：

shopt -s opt_name            # Enable (set) opt_name.
shopt -u opt_name            # Disable (unset) opt_name.
shopt opt_name                # Show current status of opt_name.
测试一下alias在非交互模式的表现和使用shopt的异同:

#!/bin/bash
 
alias echo_hello="echo Hello!"
shopt expand_aliases  
echo_hello
 
#> expand_aliases  off
#> ./test.sh: line 5: echo_hello: command not found

shopt -s  expand_aliases 
shopt expand_aliases  
echo_hello
#> expand_aliases  on
#> Hello!
另外，alias别名只在当前shell有效，不能被子shell继承，也不能像环境变量一样export。可以把alias别名定义写在.bashrc文件中，这样如果启动交互式的子shell，则子shell会读取.bashrc，从而得到alias别名定义。但是执行shell脚本时，启动的子shell处于非交互式模式，是不会读取.bashrc的。

不过，如果你一定要让执行shell脚本的子shell读取.bashrc的话，可以给shell脚本第一行的解释器加上参数：

#!/bin/bash --login
--login使得执行脚本的子shell成为一个login shell，login shell会读取系统和用户的profile及rc文件，因此用户自定义的.bashrc文件中的内容将在执行脚本的子shell中生效。

还有一个简单的办法让执行脚本的shell读取.bashrc，在脚本中主动source ~/.bashrc即可。
<!-- more -->
https://www.jianshu.com/p/1b103f6996a3

Linux shell有交互式与非交互式两种工作模式。我们日常使用shell输入命令得到结果的方式是交互式的方式，而shell脚本使用的是非交互式方式。 

 
shell提供了alias功能来简化我们的日常操作，使得我们可以为一个复杂的命令取一个简单的名字，从而提高我们的工作效率。在交互式模式下，shell的alias扩展功能是打开的，因此我们可以键入自己定义的alias别名来执行对应的命令。
 
但是，在非交互式模式下alias扩展功能默认是关闭的，此时仍然可以定义alias别名，但是shell不会将alias别名扩展成对应的命令，而是将alias别名本身当作命令执行，如果shell内置命令和PATH中均没有与alias别名同名的命令，则shell会“抱怨”找不到指定的命令。
 
那么，有没有办法在非交互式模式下启用alias扩展呢？答案是使用shell内置命令shopt命令来开启alias扩展选项。shopt是shell的内置命令，可以控制shell功能选项的开启和关闭，从而控制shell的行为。shopt的使用方式如下：
shopt -s opt_name                 Enable (set) opt_name.
shopt -u opt_name                 Disable (unset) opt_name.
shopt opt_name                    Show current status of opt_name.
alias扩展功能的选项名称是expand_aliases，我们可以在交互式模式下查看此选项是否开启：
sw@gentoo ~ $ shopt expand_aliases
expand_aliases  on

https://blog.csdn.net/liuxiangke0210/article/details/66476970

#!/bin/bash --login
–login使得执行脚本的子shell成为一个login shell，login shell会读取系统和用户的profile及rc文件，因此用户自定义的.bashrc文件中的内容将在执行脚本的子shell中生效。

还有一个简单的办法让执行脚本的shell读取.bashrc，在脚本中主动source ~/.bashrc即可。

还有一种解决办法是用source命令：
source script.sh

使当前shell读入路径为script.sh的文件并依次执行文件中的所有语句。

那么source和sh去执行脚本有什么不同呢？

sh script.sh 会重新建立一个子shell，在子shell中执行脚本里面的语句，该子shell继承父shell的环境变量，但子shell是新建的，其改变的变量不会被带回父shell，除非使用export。

source script.sh是简单地读取脚本里面的语句依次在当前shell里面执行，没有建立新的子shell。那么脚本里面所有新建、改变变量的语句都会保存在当前shell里面。

https://blog.csdn.net/qq_33709508/article/details/101822329

在根据教程安装composer之后，会生成个composer.phar的文件。但是，直接执行compser命令，提升找不到这个命令，怎么办？

答：这个composer.phar就是composer命令的文件，放到/usr/local/bin/的目录下，并改名成composer，就可以直接在系统中使用了


which 一般也是在usr/local/bin里查找，如果找不到就不返回
