---
title: http_proxy
layout: post
category: web
author: 夏泽民
---
支持http以及Sock4、Sock5类型的代理。

格式形如：http://user:pass@ip:port （socks4://user:pass@ip:port）如果没有用户名密码，那么格式形如http://:@ip:port，例如http://:@127.0.0.1:8888。
<!-- more -->
一、wget下的代理设置
1、临时生效
set "http_proxy=http://[user]:[pass]@host:port/"
或
export "http_proxy=http://[user]:[pass]@host:port/" 
执行完，就可以在当前shell 下使用wget程序了。

2、使用wget参数
wget -e "http_proxy=http://[user]:[pass]@host:port/" http://baidu.com
3、当前用户永久生效
创建$HOME/.wgetrc文件，加入以下内容：

http_proxy=代理主机IP:端口 
配置完后，就可以通过代理wget下载包了。

注：如果使用ftp代理，将http_proxy 改为ftp_proxy 即可。

二、lftp下代理设置
使lftp可以通过代理上网，可以做如下配置

echo "export http_proxy=proxy.361way.com:8888" > ~/.lftp
三、yum设置
编辑/etc/yum.conf文件，按如下配置

proxy=http://yourproxy:8080/      #匿名代理
proxy=http://username:password@yourproxy:8080/   #需验证代理
四、全局代理配置
编辑/etc/profile 或~/.bash_profile ，增加如下内容：

http_proxy=proxy.361way.com:8080
https_proxy=proxy.361way.com:8080
ftp_proxy=proxy.361way.com:8080
export http_proxy https_proxy ftp_proxy 
五、socket代理配置
这里以两个常见的socket代理软件socks5 和 tsocks 为例：

1、tsocks代理
在终端中:
sudo apt-get install tsocks
修改配置文件:
sudo nano /etc/tsocks.conf
将其内容改成以下几行并保存退出:

local = 192.168.1.0/255.255.255.0 #local表示本地的网络，也就是不使用socks代理的网络
server = 127.0.0.1 # SOCKS 服务器的 IP
server_type = 5 # SOCKS 服务版本
server_port = 9999 ＃SOCKS 服务使用的端口
运行软件：

用 tsocks 运行你的软件很简单，在终端中:tsocks 你的软件 ，如tsocks wget url

2、socks5代理
安装socks客户端工具runsocks(正常安装socks5后自带)。在libsocks5.conf文件里加入所要使用的代理服务器。配置完成，可以通过如下命令运行测试：

runsocks wget -m [http://site1 | ftp://site2]
