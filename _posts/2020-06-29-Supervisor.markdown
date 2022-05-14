---
title: Supervisor
layout: post
category: linux
author: 夏泽民
---
supervisor是一个Linux/Unix系统上的进程监控工具，supervisor是一个Python开发的通用的进程管理程序，可以管理和监控Linux上面的进程，能将一个普通的命令行进程变为后台daemon，并监控进程状态，异常退出时能自动重启。不过同daemontools一样，它不能监控daemon进程

注意：supervisor只能监控前台程序， 如果你的程序是通过fork方式实现的daemon服务，则不能用它监控，否则supervisor> status 会提示：BACKOFF  Exited too quickly (process log may have details)。 因此像apache、tomcat服务默认启动都是按daemon方式启动的，则不能通过supervisor直接运行启动脚本(service httpd start)，相反要通过一个包装过的启停脚本来完成，比如tomcat在supervisor下的启停脚本请参考：Controlling tomcat with supervisor或者supervisor-tomcat.conf。
另外，可以将supervisor随系统启动而启动，Linux 在启动的时候会执行 /etc/rc.local 里面的脚本，所以只要在这里添加执行命令即可：
# 如果是 Ubuntu 添加以下内容（这里要写全路径，因为此时PATH的环境变量未必设置）
/usr/local/bin/supervisord -c /etc/supervisord.conf

# 如果是 Centos 添加以下内容
/usr/bin/supervisord -c /etc/supervisord.conf

http://supervisord.org/
　Supervisor是一个客户端/服务器系统，采用 Python(2.4+) 开发的，它是一个允许用户管理，基于 Unix 系统进程的 Client/Server 系统，提供了大量功能来实现对进程的管理。
　
　supervisor是一个 Client/Server模式的系统，允许用户在类unix操作系统上监视和控制多个进程，或者可以说是多个程序。supervisor与launchd，daemontools，runit等程序有着相同的功能，与其中某些程序不同的是，它并不作为“id 为 1的进程”而替代init。相反，它用于控制应用程序，像启动其它程序一样，通俗理解就是，把Supervisor服务管理的进程程序，它们作为supervisor的子进程来运行，而supervisor是父进程。supervisor来监控管理子进程的启动关闭和异常退出后的自动启动。

至于为什么要用supervisor来管理进程，是因为相对于linux传统的进程管理(即系统自带的init 进程管理)方式来说，它有很多的优势：

1) 简单方便
通常管理linux进程的时候，一般来说都需要自己编写一个能够实现进程start/stop/restart/reload功能的脚本，然后丢到/etc/init.d/下面。其实这么做有很多不好的地方:
a) 编写这个脚本，耗时耗力。
b) 当这个进程挂掉的时候，linux不会自动重启它的，想要自动重启的话，还要自己另外写一个监控重启脚本。

supervisor则可以完美的解决上面这那两个问题! 那么supervisor怎么解决呢?
a) supervisor管理进程，就是通过fork/exec的方式把这些被管理的进程，当作supervisor的子进程来启动。这样的话，只要在supervisor的配置文件中，把要管理的进程的可执行文件的路径写进去就OK了。这样就省下了自己写脚本管理linux进程的麻烦了。
b) 被管理进程作为supervisor的子进程，当子进程挂掉的时候，父进程可以准确获取子进程挂掉的信息的，所以也就可以对挂掉的子进程进行自动重启了, 至于重启还是不重启，也要看配置文件里面有没有设置autostart=true。

2) 精确
linux对进程状态的反馈有时候不太准确, 也就是说linux进程通常很难获得准确的up/down状态, Pidfiles经常说谎!  而supervisor监控子进程，得到的子进程状态无疑是准确的。supervisord将进程作为子进程启动，所以它总是知道其子进程的正确的up/down状态，可以方便的对这些数据进行查询. 

3) 进程分组
进程支持分组启动和停止，也支持启动顺序，即‘优先级’，supervisor允许为进程分配优先级，并允许用户通过supervisorctl客户端发出命令，如“全部启动”和”重新启动所有“，它们以预先分配的优先级顺序启动。还可以将进程分为”进程组“，一组逻辑关联的进程可以作为一个单元停止或启动。进程组supervisor可以对进程组统一管理，也就是说我们可以把需要管理的进程写到一个组里面，然后把这个组作为一个对象进行管理，如启动，停止，重启等等操作。而linux系统则是没有这种功能的，想要停止一个进程，只能一个一个的去停止，要么就自己写个脚本去批量停止。

4) 集中式管理
supervisor管理的进程，进程组信息，全部都写在一个ini格式的文件里就OK了。管理supervisor时, 可以在本地进行管理，也可以远程管理，而且supervisor提供了一个web界面，可以在web界面上监控，管理进程。 当然了，本地，远程和web管理的时候，需要调用supervisor的xml_rpc接口。

5) 可扩展性
supervisor有一个简单的事件（event）通知协议，还有一个用于控制的XML-RPC接口，可以用Python开发人员来扩展构建。

6) 权限
总所周知, linux的进程特别是侦听在1024端口之下的进程，一般用户大多数情况下，是不能对其进行控制的。想要控制的话，必须要有root权限。然而supervisor提供了一个功能，可以为supervisord或者每个子进程，设置一个非root的user，这个user就可以管理它对应的进程了。

7) 兼容性，稳定性
supervisor由Python编写，在除Windows操作系统以外基本都支持，如linux，Mac OS x,solaris,FreeBSD系统

二、Supervisor组成部分
1)supervisord: 服务守护进程
supervisor服务器的进程名是supervisord。它主要负责在自己的调用中启动子程序，响应客户端的命令，重新启动崩溃或退出的进程，记录其子进程stdout和stderr的输出，以及生成和处理对应于子进程生命周期中的"event"服务器进程使用的配置文件，通常路径存放在/etc/supervisord.confa中。此配置文件是INI格式的配置文件。
2) supervisorctl：命令行客户端
supervisor命令行的客户端名称是supervisorctl。它为supervisord提供了一个类似于shell的交互界面。使用supervisorctl，用户可以查看不同的supervisord进程列表，获取控制子进程的状态，如停止和启动子进程
3) Web Server：提供与supervisorctl功能相当的WEB操作界面
一个可以通过Web界面来查看和控制进程的状态，默认监听在9091上。
4) XML-RPC Interface：XML-RPC接口
supervisor用于控制的XML-RPC接口


<!-- more -->
简单
Supervisor通过简单的INI样式（可以修改为.conf后缀）配置文件进行配置，该文件易于学习。它提供了许多每个进程选项，使您的生活更轻松，如重新启动失败的进程和自动日志轮换。
集中
主管为您提供一个启动，停止和监控流程的位置。流程可以单独控制，也可以成组控制。您可以将Supervisor配置为提供本地或远程命令行和Web界面。
高效
主管通过fork / exec启动其子进程，子进程不进行守护。当进程终止时，操作系统会立即向Supervisor发出信号，这与某些依赖麻烦的PID文件和定期轮询重新启动失败进程的解决方案不同。
扩展
Supervisor有一个简单的事件通知协议，用任何语言编写的程序都可以用它来监视它，以及一个用于控制的XML-RPC接口。它还使用可由Python开发人员利用的扩展点构建。
兼容
除了Windows之外，Supervisor几乎可以处理所有事情。它在Linux，Mac OS X，Solaris和FreeBSD上经过测试和支持。它完全用Python编写，因此安装不需要C编译器。
久经考验
虽然Supervisor今天非常活跃，但它并不是新软件。主管已存在多年，已在许多服务器上使用。
 

3、安装配置 Supervisor
3.1 各个平台安装Supervisor
（1）在 linux 中使用以下命令进行安装：

 centos
1
yum install supervisor
 ubuntu
1
sudo apt-get install supervisor
 python
1
pip install supervosor easy_install supervisor
 

（2）在 masOS 中直接使用brew工具进行安装即可：

1
brew install supervisor
 

3.2 Supervisor 配置
（1）linux 安装完后会有一个主配置文件/etc/supervisord.conf，和一个/etc/supervisord.d 自配置文件目录
$ ls /etc/supervisord.*
/etc/supervisord.conf
/etc/supervisor
　　

（2）修改主配置文件，设置自配置文件生效的后缀

$ vim /etc/supervisord.conf   在最后一行
[include]
files = supervisord.d/*.conf
　　

（3）为了方便管理，就在自配置文件目录下，创建项目的配置文件
$ cd /etc/supervisord.d/
$ vim ProjectName.conf
[program: ProjectName]
command=dotnet ProjectName.dll   ; 运行程序的命令
directory=/usr/local/ProjectName/   ; 命令执行的目录
autorestart=true   ; 程序意外退出是否自动重启
autostart=true   ; 是否自动启动
stderr_logfile=/var/log/ProjectName.err.log   ; 错误日志文件
stdout_logfile=/var/log/ProjectName.out.log   ; 输出日志文件
environment=ASPNETCORE_ENVIRONMENT=Production   ; 进程环境变量
user=root   ; 进程执行的用户身份
stopsignal=INT
startsecs=1   ; 自动重启间隔
　　

3.3 启动 Supervisor 服务
（1）开启服务，并设为开机自启

$ systemctl start supervisord.service
$ systemctl enable supervisord.service
Created symlink from /etc/systemd/system/multi-user.target.wants/supervisord.service to /usr/lib/systemd/system/supervisord.service.
 

（2）查询服务状态

$ systemctl status supervisord.service
● supervisord.service - Process Monitoring and Control Daemon
   Loaded: loaded (/usr/lib/systemd/system/supervisord.service; disabled; vendor preset: disabled)
   Active: active (running) since Fri 2019-01-11 15:00:23 CST; 57min ago
  Process: 910 ExecStart=/usr/bin/supervisord -c /etc/supervisord.conf (code=exited, status=0/SUCCESS)
 Main PID: 913 (supervisord)
   CGroup: /system.slice/supervisord.service
           ├─913 /usr/bin/python /usr/bin/supervisord -c /etc/supervisord.conf
           └─914 dotnet eXiu.OBD.Host.dll
 
Jan 11 15:00:23 iZe4iwiics91xjZ systemd[1]: Starting Process Monitoring and Control Daemon...
Jan 11 15:00:23 iZe4iwiics91xjZ systemd[1]: Started Process Monitoring and Control Daemon.
　　

（3）查看进程认证

1
2
3
$ ps -ef | grep dotnet ProjectName
root       914   913  0 15:00 ?        00:00:05 dotnet ProjectName.dll
root      3455  3058  0 15:58 pts/0    00:00:00 grep --color=auto dotnet
 

4、报错处理
（1）使用Supervisor 为服务创建守护进程失败

1
2
3
Error: Another program is already listening on a port that one of our HTTP servers is configured to use. 
Shut this program down first before starting supervisord.
For help, use /usr/bin/supervisord –h
　　是因为有一个使用supervisor配置的应用程序正在运行，需要执行supervisorctl shutdown命令终止它，或重新创建一个ProjectName.conf文件再执行第一条命令。

 

（2）如果运行supervisorctl出现以下错误

1
error: <class 'socket.error'>, [Errno 111] Connection refused: file: /usr/lib64/python2.6/socket.py line: 567
 

　　说明Supervisor 服务没有启动成功，或Supervisor 服务被关闭了，重启启动服务即可。

 

5、supervisorctl 常用命令
$ sudo service supervisor stop 停止supervisor服务
$ sudo service supervisor start 启动supervisor服务
$ supervisorctl shutdown #关闭所有任务
$ supervisorctl stop|start program_name #启动或停止服务
$ supervisorctl status #查看所有任务状态

https://www.cnblogs.com/kevingrace/p/7525200.html


http://liyangliang.me/posts/2015/06/using-supervisor/


https://juejin.im/post/5d80da83e51d45620c1c5471

https://www.jianshu.com/p/535c22ea6e28

http://wangshengzhuang.com/2017/05/26/%E8%BF%90%E7%BB%B4%E7%9B%B8%E5%85%B3/Supervisor/Supervisor%E7%AE%80%E4%BB%8B/

https://agvszwk.github.io/2019/06/23/supervisor%E7%9A%84%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86/


supervisor管理nginx进程

【nginx的这个启动命令默认是后台启动，supervisor不能监控后台程序，所以就会一直执行这个命令，不断的启动Nginx，就会报错。 加上-g ‘daemon off;’ 这个参数可解决这问题。】

[program: nginx]    
;管理的子进程。后面的是名字，最好写的具有代表性，避免日后”认错“
command=/usr/local/nginx/sbin/nginx  -g 'daemon off;'    
;我们的要启动进程的命令路径，可以带参数。
directory=/root ;   
;进程运行前，会先切换到这个目录
autorestart=true； 
;设置子进程挂掉后自动重启的情况，有三个选项，false,unexpected和true。false表示无论什么情况下，都不会被重新启动；unexpected表示只有当进程的退出码不在下面的exitcodes里面定义的退出码的时候，才会被自动重启。当为true的时候，只要子进程挂掉，将会被无条件的重启
autostart=true ;    
;如果是true的话，子进程将在supervisord启动后被自动启动，默认就是true
stderr_logfile=/home/work/super/nginx_error.log ;   
;日志，没什么好说的
stdout_logfile=/home/work/super/nginx_stdout.log ;  
;日志
environment=ASPNETCORE_ENVIRONMENT=Production ;  
;这个是子进程的环境变量,默认为空
user=nginx ;     
;可以用来管理该program的用户
stopsignal=INT   
;进程停止信号，可以为TERM, HUP, INT, QUIT, KILL, USR1等,默认为TERM
startsecs=10 ;    
;子进程启动多少秒之后,此时状态如果是running,我们认为启动成功了,默认值1
startretries=5 ;   当进程启动失败后，最大尝试的次数。当超过5次后，进程的状态变为FAIL
stopasgroup=true   
;这个东西主要用于，supervisord管理的子进程，这个子进程本身还有子进程。那么我们如果仅仅干掉supervisord的子进程的话，子进程的子进程有可能会变成孤儿进程。所以可以设置这个选项，把整个该子进程的整个进程组干掉。默认false

https://blog.csdn.net/weixin_44449055/article/details/88966342

https://www.oschina.net/translate/why-you-dont-need-to-run-sshd-in-docker


1.这是因为启动 tomcat的方式不对，在linux命令行模式下我们启动可以使用如下脚本
./apache-tomcat-7.0.70/bin/startup.sh

2.在supervisor的启动命令中不能使用这种方式了要使用如下方式：
command=/root/tools/apache-tomcat-7.0.70/bin/catalina.sh run

https://www.cnblogs.com/ExMan/p/12570394.html

supervisor是如何知道某项目挂掉的？

子进程异常退出的时候，作为父进程肯定是能收到信号的。

https://blog.csdn.net/jia_xiaoli/article/details/16983531

1，要了解他的实现原理

通过配置command这个参数，来运行子进程，并且获取到该进程的pid，然后再对该pid进行监控，如果，pid不存在就就拉起来

所以，重点来了！！！！ pid 一定要保证存在

下面是常犯的错误：command = /etc/init.d/nginx start

这样一定是会失败的，因为/etc/ini.t/nginx start 启动完之后，pid就消失了，这样supervisor就会认为他管理的进程挂了，于是执行拉起的动作，这样就会导致程序反复快速结束并拉起，最终导致监控失败



下面是常犯的错误：针对于php的问题，supervisor只能监控前端进程，不能监控守护进程，所以要进行更改





解决思路：

针对于/etc/init.d/nginx start的问题

可以写个while脚本，脚本的内容是个while true的检测脚本，这样保证pid一直存在的情况下，还能通过脚本拉起nginx 这个进程

https://zhuanlan.zhihu.com/p/48248668

https://www.jianshu.com/p/caa57a74ac7f

