---
title: composer-source
layout: post
category: php
author: 夏泽民
---
<!-- more -->
1、安装
curl -sS https://getcomposer.org/installer | php mv composer.phar /usr/local/bin/composer
如果上面出现问题
可以这样
curl -sS https://getcomposer.org/installer | php 
mv composer.phar /usr/local/bin/composer 
如果curl下载较慢,一直卡在downloading
可以这样
wget https://getcomposer.org/installer
php installer
mv composer.phar /usr/local/bin/composer
2、配置国内源

方法一： 修改 composer 的全局配置文件（推荐方式）
 打开命令行窗口（windows用户）或控制台（Linux、Mac 用户）并执行如下命令：
composer config -g repo.packagist composer https://packagist.phpcomposer.com
方法二： 修改当前项目的 composer.json 配置文件：
打开命令行窗口（windows用户）或控制台（Linux、Mac 用户），进入你的项目的根目录（也就是 composer.json文件所在目录），执行如下命令：
composer config repo.packagist composer https://packagist.phpcomposer.com
上述命令将会在当前项目中的 composer.json 文件的末尾自动添加镜像的配置信息

$  composer config repo.packagist composer https://packagist.phpcomposer.com

$vi composer.json

  "repositories": {
        "packagist": {
            "type": "composer",
            "url": "https://packagist.phpcomposer.com"
        }
    }
}