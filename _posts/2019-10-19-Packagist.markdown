---
title: Packagist
layout: post
category: lang
author: 夏泽民
---
https://packagist.org/
国内镜像：
https://pkg.phpcomposer.com/
国内仓库
http://packagist.p2hp.com/

镜像用法
有两种方式启用本镜像服务：

系统全局配置： 即将配置信息添加到 Composer 的全局配置文件 config.json 中。见“方法一”
单个项目配置： 将配置信息添加到某个项目的 composer.json 文件中

<!-- more -->
方法一： 修改 composer 的全局配置文件（推荐方式）
打开命令行窗口（windows用户）或控制台（Linux、Mac 用户）并执行如下命令：

复制
composer config -g repo.packagist composer https://packagist.phpcomposer.com

方法二： 修改当前项目的 composer.json 配置文件：
打开命令行窗口（windows用户）或控制台（Linux、Mac 用户），进入你的项目的根目录（也就是 composer.json 文件所在目录），执行如下命令：

复制
composer config repo.packagist composer https://packagist.phpcomposer.com
上述命令将会在当前项目中的 composer.json 文件的末尾自动添加镜像的配置信息（你也可以自己手工添加）：

复制
"repositories": {
    "packagist": {
        "type": "composer",
        "url": "https://packagist.phpcomposer.com"
    }
}
以 laravel 项目的 composer.json 配置文件为例，执行上述命令后如下所示（注意最后几行）：

复制
{
    "name": "laravel/laravel",
    "description": "The Laravel Framework.",
    "keywords": ["framework", "laravel"],
    "license": "MIT",
    "type": "project",
    "require": {
        "php": ">=5.5.9",
        "laravel/framework": "5.2.*"
    },
    "config": {
        "preferred-install": "dist"
    },
    "repositories": {
        "packagist": {
            "type": "composer",
            "url": "https://packagist.phpcomposer.com"
        }
    }
}
OK，一切搞定！试一下 composer install 来体验飞一般的速度吧！

镜像原理：
一般情况下，安装包的数据（主要是 zip 文件）一般是从 github.com 上下载的，安装包的元数据是从 packagist.org 上下载的。

然而，由于众所周知的原因，国外的网站连接速度很慢，并且随时可能被“墙”甚至“不存在”。

“Packagist 中国全量镜像”所做的就是缓存所有安装包和元数据到国内的机房并通过国内的 CDN 进行加速，这样就不必再去向国外的网站发起请求，从而达到加速 composer install 以及 composer update 的过程，并且更加快速、稳定。因此，即使 packagist.org、github.com 发生故障（主要是连接速度太慢和被墙），你仍然可以下载、更新安装包。


解除镜象：
如果需要解除镜像并恢复到 packagist 官方源，请执行以下命令：

复制
composer config -g --unset repos.packagist
执行之后，composer 会利用默认值（也就是官方源）重置源地址。

将来如果还需要使用镜像的话，只需要根据前面的“镜像用法”中介绍的方法再次设置镜像地址即可。

项目配置：
$ composer config repo.packagist composer  https://packagist.org/

composer.json 新增内容
"repositories": {
    "packagist": {
        "type": "composer",
        "url": "https://packagist.org/"
    }
}

$ composer require init/lib
 The "https://packagist.org/packages.json" file could not be downloaded: SSL operation failed with code 1. OpenSSL Error messages:  
  
改成国内镜像
composer config repo.packagist composer https://packagist.phpcomposer.com



composer require xxxx/xxx               # 这时候会报错， Could not find package xxxx/xxx at any version for your minimum-stability (stable). Check the package。。。猜测是我的composer使用的国内镜像，可能是没有同步的原因，使用这个命令把“源”改回去还是不行。
composer config repo.packagist composer https://packagist.org  # 继续猜测，原来我的组件还没有在github上发布正式，这个时候还是开发版本dev-master.应该加上dev-master版本。
composer require xxxx/xxx:dev-master     # 成功

但是一般无法成功

使用国内仓库

搭建私有仓库
使用 Satis 搭建私有仓库
1. 建立项目
使用 Composer 自带的建项目功能，这个相当于git clone+composer install+ 运行 post-install 脚本。

$ composer create-project composer/satis my-satis --stability=dev --keep-vcs
2. 建立配置文件
在/path/to/my-satis目录下建立satis.json文件
name: 项目名称
homepage : 私有包主页，后续会用到。
repositories : 资源包来源，里面配置私有仓库url，就是上前面创建的私有Git仓库地址。
require : 配置 git仓库中存在的包。
3. 生成仓库列表
执行：
php bin/satis build satis.json ./web
就可以在path/to/my-satis/web/里生成仓库列表了。
可能会报协议错误，默认是禁止 http 方式获取代码。需要单独配置开启。
执行完毕后。会在项目根目录生成 web 目录。
4. 配置 webServer
将 web 目录配置 webServer 访问。虚拟域名就是之前我们配置的 homepage : packagist.example.com
