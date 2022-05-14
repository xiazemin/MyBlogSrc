---
title: completion bash自动补全
layout: post
category: linux
author: 夏泽民
---
要让可编程命令补全功能在你的终端起作用 ，你只需要如下执行/etc/bash_completion即可：

# . /etc/bash_completion

你也可以取消/etc/bash.bashrc
3. 定义一个命令名补全

通过 -c 选项可以将所有的可用命令作为一个命令的补全参数。在下面的例子里面，为which命令定义了一个补全(LCTT译注：在按两下TAB时，可以列出所有命令名作为可补全的参数)。

$ complete -c which

$ which [TAB][TAB]

Display all 2116 possibilities? (y or n)

如上，如果按下 ‘y’，就会列出所有的命令名。

4. 定义一个目录补全

通过选项 -d，可以定义一个仅包含目录名的补全参数。在下面的例子中，为ls命令定义了补全。

$ ls

countfiles.sh  dir1/          dir2/          dir3/

$ complete -d ls

$ ls [TAB][TAB]

dir1/          dir2/          dir3/

如上，连按下 TAB 仅会显示目录名。
https://blog.csdn.net/weixin_36294922/article/details/116831695
<!-- more -->
https://www.cnblogs.com/327999487heyu/articles/5621546.html
https://blog.csdn.net/weixin_36021459/article/details/116588849
https://www.cnblogs.com/xulei13140106/p/5946359.html

