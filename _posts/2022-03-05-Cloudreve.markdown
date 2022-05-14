---
title: Cloudreve 搭建私人网盘
layout: post
category: golang
author: 夏泽民
---
https://github.com/cloudreve/Cloudreve

<!-- more -->
Cloudreve是个公有网盘程序，你可以用它快速搭建起自己的网盘服务，公有云/私有云都可。Cloudreve底层支持 本机存储、从机存储、阿里云OSS、又拍云、腾讯云COS、七牛云存储、OneDrive（国际版/世纪互联版），每种存储方式的上传下载都是客户端直传。具有以下特性：
支持本机、从机、七牛、阿里云OSS、腾讯云COS、又拍云、OneDrive (包括世纪互联版) 作为存储端

上传/下载 支持客户端直传，支持下载限速

可对接Aria2离线下载（支持所有存储策略，下载完成后自动中转）

在线压缩/解压缩、多文件打包下载（支持所有存储策略）

覆盖全部存储策略的WebDAV协议支持

文件拖拽管理，拖拽上传、目录上传、流式上传处理

多用户、用户组

创建文件、目录的分享链接，可设定自动过期

视频、图像、音频、文本、Office文档在线预览

自定义配色、黑暗模式、PWA应用、全站单页应用


启动 Cloudreve

Linux下，直接解压并执行主程序即可：


#解压获取到的主程序
tar -zxvf cloudreve_VERSION_OS_ARCH.tar.gz
# 赋予执行权限
chmod +x ./cloudreve
# 启动 Cloudreve
./cloudreve


Cloudreve 在首次启动时，会创建初始管理员账号，请注意保管管理员密码，此密码只会在首次启动时出现。如果您忘记初始管理员密码，需要删除同级目录下的cloudreve.db，重新启动主程序以初始化新的管理员账户。

Cloudreve 默认会监听5212端口。你可以在浏览器中访问http://服务器IP:5212进入 Cloudreve。以上步骤操作完后，最简单的部署就完成了。

https://mp.weixin.qq.com/s/xNYU8Zg5jGSUDzU4g2Ozfw