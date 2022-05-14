---
title: AlgorithmVisualizer
layout: post
category: algorithm
author: 夏泽民
---
https://algorithm-visualizer.org/dynamic-programming/floyd-warshalls-shortest-path
https://github.com/algorithm-visualizer/algorithm-visualizer/wiki
https://github.com/xiazemin/algorithm-visualizer
https://github.com/xiazemin/fucking-algorithm/tree/master/%E9%AB%98%E9%A2%91%E9%9D%A2%E8%AF%95%E7%B3%BB%E5%88%97	
https://github.com/xiazemin/fucking-algorithm
<!-- more -->

{% raw %}
AlgorithmVisualizer项目运行环境搭建
在Ubuntu14.04 x64下搭建AlgorithmVisualizer项目运行环境

项目地址：https://github.com/parkjs814/AlgorithmVisualizer

演示项目：http://parkjs814.github.io/AlgorithmVisualizer

 

1. node+npm安装

默认下Ubuntu14.04是没有安装nodejs，需要用户自己安装[推荐：方式三]

方式一：

官网下载最新的nodejs[v7.5,2017.2.16] ,解压，建立软链接


# 下载并解压 node-v7.5.0-linux-x86.tar.xz
tar -xJf node-v7.5.0-linux-x86.tar.xz
# 移到通用的软件安装目录 /opt/
sudo mv node-v7.5.0-linux-x86 /opt/
 
# 安装 npm 和 node 命令到系统命令
sudo ln -s /opt/node-v7.5.0-linux-x86/bin/node /usr/local/bin/node
sudo ln -s /opt/node-v7.5.0-linux-x86/bin/npm /usr/local/bin/npm
 
# 验证：
node -v
v7.5.0
npm -v
4.1.2 　　　　
 

参考链接：http://www.linuxidc.com/Linux/2016-09/135487.htm

 

方式二：

使用Ubuntu提示的方式安装:


sudo apt-get install nodejs
sudo apt-get install npm
nodejs -v
npm -v　　
成功安装，但是版本很老[补充：Ubuntu里node命令无效解决方法 ]

终于发现了一个可以管理node版本的第三方库，n来自tj大神。
安装n有几种方式，最快捷的是用npm安装，前面的安装已经为这里打好了铺垫，现在只需要运行 sudo npm install -g n ，安装好后升级nodejs  sudo n latest  

Use or install the latest official release:
sudo n latest
 
Use or install the stable official release:
sudo n stable　　　
 参考链接：https://segmentfault.com/a/1190000007148749

方式三*：

下载到本地后解压，移动到/opt目录，配置/etc/profile全局环境变量 

sudo mv node-v7.5.0-linux-x86 /opt/
sudo gedit /etc/profile
# NODEJS ENV
export NODE=/opt/node-v7.5.0-linux-x86
export PATH=${NODE}/bin:$PATH
# 立即生效
. /etc/profile　　
 

2. 编译运行AlgorithmVisualizer项目 


# install gulp globally so you can run it from the command line
npm install -g gulp-cli
 
# install all dependencies
npm install
 
# run gulp to start the livereload server on http://localhost:8080
gulp 
注：如果使用方式一，在执行 npm install -g gulp-cli 之后需要额外执行 sudo ln -s /opt/node-v7.5.0-linux-x86/bin/gulp /usr/local/bin/gulp ，

否则在使用 gulp -v 时会报错：找不到gulp　
{% endraw %}

https://www.cnblogs.com/AbcFly/p/6405164.html

https://www.jianshu.com/p/9f693b7d0bd4
https://blog.csdn.net/m0_37577608/article/details/91680724
https://greyireland.gitbook.io/algorithm-pattern/suan-fa-si-wei/recursion
https://github.com/xiazemin/algorithm-pattern
