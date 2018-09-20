---
title:nginx php-fpm unix-socket  通信
layout: post
category: web
author: 夏泽民
---
首先就是找php-fpm的配置文件修改配置，我用的php是7.1.4

在php-fpm.conf里面找不到lister=127.0.0.1:9000的配置，然后找到末尾发现有个inclue

进入php-fpm.d找到www.conf，在里面修改listen = 127.0.0.1:9000为

listen = /tmp/php-fpm.sock

然后到nginx修改配置文件

把fastcgi_pass 127.0.0.1:9000;改为fastcgi_pass unix:/tmp/php-fpm.sock;

然后重启php-fpm和nginx，发现无效，这个时候是/tmp/php-fpm.sock没有权限，直接chmod 777 /tmp/php-fpm.sock

然后就能访问了，用sock文件会比访问9000端口快些

PS:后记

今天要安装一个拓展，安装拓展之后重启php-fpm发现网页报错了

具体我猜是重新启动了php-fpm，新产生一个php-fpm.sock，所以把这个新的php-fpm.sock设置权限777就可以了
<!-- more -->
nginx和fastcgi的通信方式有两种，一种是TCP的方式，一种是unix socke方式。两种方式各有优缺点，这里先给出两种的配置方法，然后再对性能、安全性等做出总结。

TCP是使用TCP端口连接127.0.0.1:9000
Socket是使用unix domain socket连接套接字/dev/shm/php-cgi.sock（很多教程使用路径/tmp，而路径/dev/shm是个tmpfs，速度比磁盘快得多）,在服务器压力不大的情况下，tcp和socket差别不大，但在压力比较满的时候，用套接字方式，效果确实比较好。

配置指南

TCP配置方式

TCP通信配置起来很简单，三步即可搞定

第一步，编辑 /etc/nginx/conf.d/你的站点配置文件（如果使用的默认配置文件，修改/etc/nginx/sites-available/default）

将fastcgi_pass参数修改为127.0.0.1:9000，像这样：
location ~ \.php$ {
      index index.php index.html index.htm;
      include /etc/nginx/fastcgi_params;
      fastcgi_pass 127.0.0.1:9000;
      fastcgi_index index.php;
      include fastcgi_params;
 }
 第二步，编辑php-fpm配置文件 /etc/php5/fpm/pool.d/www.conf

将listen参数修改为127.0.0.1:9000，像这样：
listen=127.0.0.1:9000
 第三步，重启php-fpm，重启nginx
unix socket配置方式

unix socket其实严格意义上应该叫unix domain socket，它是*nix系统进程间通信（IPC）的一种被广泛采用方式，以文件（一般是.sock）作为socket的唯一标识（描述符），需要通信的两个进程引用同一个socket描述符文件就可以建立通道进行通信了。

Unix domain socket 或者 IPC socket是一种终端，可以使同一台操作系统上的两个或多个进程进行数据通信。与管道相比，Unix domain sockets 既可以使用字节流和数据队列，而管道通信则只能通过字节流。Unix domain sockets的接口和Internet socket很像，但它不使用网络底层协议来通信。Unix domain socket 的功能是POSIX操作系统里的一种组件。Unix domain sockets 使用系统文件的地址来作为自己的身份。它可以被系统进程引用。所以两个进程可以同时打开一个Unix domain sockets来进行通信。不过这种通信方式是发生在系统内核里而不会在网络里传播。

配置需要五步

第一步，决定你的socket描述符文件的存储位置。

可以放在系统的任意位置，如果想要更快的通信速度，可以放在/dev/shm下面，这个目录是所谓的tmpfs，是RAM可以直接使用的区域，所以，读写速度都会很快。

决定了文件位置，就要修改文件的权限了，要让nginx和php-fpm对它都有读写的权限，可以这样：
sudo touch /dev/shm/fpm-cgi.sock
sudo chown www-data:www-data /dev/shm/fpm-cgi.sock
sudo chmod 666 /dev/shm/fpm-cgi.sock
 第二步，修改php-fpm配置文件/etc/php5/fpm/pool.d/www.conf
将listen参数修改为/dev/shm/fpm-cgi.sock，像这样：
listen=/dev/shm/fpm-cgi.sock
 将listen.backlog参数改为-1，内存积压无限大，默认是128，并发高了之后就会报错
 ; Set listen(2) backlog. A value of '-1' means unlimited.
 ; Default Value: 128 (-1 on FreeBSD and OpenBSD)
 listen.backlog = -1
 第三步，修改nginx站点配置文件

将fastcgi_pass参数修改为unix:/dev/shm/fpm-cgi.sock，像这样：
location~\.php${
      indexindex.phpindex.htmlindex.htm;
      include/etc/nginx/fastcgi_params;
      fastcgi_passunix:/dev/shm/fpm-cgi.sock;
      fastcgi_indexindex.php;
      includefastcgi_params;
}
第四步，修改/etc/sysctl.conf 文件，提高内核级别的并发连接数（这个系统级的配置文件我也不是特别熟悉，参考的是这篇博客：《Php-fpm TcpSocket vs UnixSocket》）
sudoecho'net.core.somaxconn = 2048'>>/etc/sysctl.conf
sudosysctl-p
第五步， 重启nginx和php-fpm服务（最好先重启php-fpm再重启nginx）
两种通信方式的分析和总结

从原理上来说，unix socket方式肯定要比tcp的方式快而且消耗资源少，因为socket之间在nginx和php-fpm的进程之间通信，而tcp需要经过本地回环驱动，还要申请临时端口和tcp相关资源。

当然还是从原理上来说，unix socket会显得不是那么稳定，当并发连接数爆发时，会产生大量的长时缓存，在没有面向连接协议支撑的情况下，大数据包很有可能就直接出错并不会返回异常。而TCP这样的面向连接的协议，多少可以保证通信的正确性和完整性。

当然以上主要是半懂不懂的理论分析加主观臆测，具体的差别还是要通过测试数据来说话，以后有空，会进行这方面的测试。从网上别人博客的测试数据，我的理论分析差不多是对的。至于你选择哪种方式，我只能说“鱼和熊掌不可兼得也”，通过高超的运维和配置技巧，在性能和稳定性上做一个平衡吧。
其实，如果nginx做要做负载均衡的话，根本也不要考虑unix socket的方式了，只能采用TCP的方式。现在我的小站没有那么高的并发量，所以就用unix socket了，以后如果有了高并发业务，再进行一些参数调整即可应付，如果真要是无法支撑，那只能做负载均衡了，到时候自然会选择TCP方式。
