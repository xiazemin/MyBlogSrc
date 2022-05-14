---
title: COMPOSER_HOME 环境变量
layout: post
category: web
author: 夏泽民
---
全局执行 global
global 命令允许你在 COMPOSER_HOME 目录下执行其它命令，像 install、require 或 update。

并且如果你将 $COMPOSER_HOME/vendor/bin 加入到了 $PATH 环境变量中，你就可以用它在命令行中安装全局应用，下面是一个例子：

php composer.phar global require fabpot/php-cs-fixer:dev-master
现在 php-cs-fixer 就可以在全局范围使用了（假设你已经设置了你的 PATH）。如果稍后你想更新它，你只需要运行 global update：

php composer.phar global update

https://docs.phpcomposer.com/03-cli.html#COMPOSER_HOME
<!-- more -->
COMPOSER_HOME
COMPOSER_HOME 环境变量允许你改变 Composer 的主目录。这是一个隐藏的、所有项目共享的全局目录（对本机的所有用户都可用）。

它在各个系统上的默认值分别为：

*nix /home/<user>/.composer。
OSX /Users/<user>/.composer。
Windows C:\Users\<user>\AppData\Roaming\Composer。

COMPOSER_HOME/config.json
你可以在 COMPOSER_HOME 目录中放置一个 config.json 文件。在你执行 install 和 update 命令时，Composer 会将它与你项目中的 composer.json 文件进行合并。

该文件允许你为用户的项目设置 配置信息 和 资源库。

若 全局 和 项目 存在相同配置项，那么项目中的 composer.json 文件拥有更高的优先级。


COMPOSER_CACHE_DIR
COMPOSER_CACHE_DIR 环境变量允许你设置 Composer 的缓存目录，这也可以通过 cache-dir 进行配置。

它在各个系统上的默认值分别为：

*nix and OSX $COMPOSER_HOME/cache。
Windows C:\Users\<user>\AppData\Local\Composer 或 %LOCALAPPDATA%/Composer。