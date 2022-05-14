---
title: systemctl init.d
layout: post
category: linux
author: 夏泽民
---
docker容器不支持 systemctl

情况一：正常情况（系统有service命令）

重启服务命令：[root@centos6 /]# service crond restart



启动服务命令：[root@centos6 /]# service crond start

停止服务命令：[root@centos6 /]# service crond stop



情况二：当linux发行的版本没有service这个命令时候，用如下命令进行停止启动：

停止服务：[root@centos6 /]# /etc/init.d/cron空格stop

启动服务：[root@centos6 /]# /etc/init.d/cron空格start
<!-- more -->
ubuntu 在执行crond restart 时提示cron: can’t lock /var/run/crond.pid, otherpid may be 2699: Resource temporarily unavailable

解决方案： rm -rf /var/run/crond.pid 重启即可 
重新加载
处理方法： /etc/init.d/cron reload 

重启服务

处理方法：/etc/init.d/crond restart  


sudo systemctl status apache2.service
sudo /bin/systemctl status apache2.service
sudo /etc/init.d/apache2 status
sudo service apache2 status
All the above commands work.

Should I prefer one command over the other?
If yes then why?
Are there any other commands I need to be aware of?
Using init.d in Monit caused issues when I wanted to use the status option (status will be that the service is offline when it was actually online -- restarted by Monit). Change the code in Monit from inid.d to /bin/systemctl fixed it.

It seems that using init.d provides more information on what happened that the others. If I should be using one of the other commands, is it possible to have them display more information on what was done?


To start, there's a whole history and struggle between going from SysVInit to SystemD. Rather than trying to break that all down in one answer though, I'll refer you to some google venturing for more details on the history as well as one particular article on the topic:

http://www.tecmint.com/systemd-replaces-init-in-linux/

In summary though, it's been a slow and arduous transition. Some legacy features were kept intact (such as init.d to some degree). If you have the option to use systemctl for your service control I recommend using that one. It's the foreseeable future for Linux and eventually older SysVInit methods will be considered deprecated entirely and removed.

To cover each one you listed specifically:

sudo systemctl status apache2.service
This is the new SystemD approach to handling services. Moving forward, applications on Linux are designed to uses the systemd method, not any other.

sudo /bin/systemctl status apache2.service
This is the same thing as the previous command. The only difference in this case is that it's not depending on the shell's $PATH environment variable to find the command, it's listing the command explicitly by including the path to the command.

sudo /etc/init.d/apache2 status
This is the original SysVInit method of calling on a service. Init scripts would be written for a service and placed into this directory. While this method is still used by many, service was the command that replaced this method of calling on services in SysVInit. There's some legacy functionality for this on newer systems with SystemD, but most newer programs don't include this, and not all older application init scripts work with it.

sudo service apache2 status
This was the primary tool used on SysVInit systems for services. In some cases it just linked to the /etc/init.d/ scripts, but in other cases it went to an init script stored elsewhere. It was intended to provide a smoother transition into service dependency handling.


Lastly, you mention wanting to know how to get more information out of the commands, since some provide more information than others. This is almost always determined by the application and how they designed their init or service file. As a general rule though, if it completed silently it was successful. However, to verify a start, stop, or restart, you can use the status sub-command to see how it is doing. You mentioned a status command being incorrect on an old init script. That is a bug that the application developers would have to look at. However, since init scripts are becoming the deprecated method of handling services, they may just ignore the bug until they remove the init script option entirely. The systemctl status should always work correctly otherwise a bug should be logged with the application developers.

init,service和systemctl的区别
、service是一个脚本命令，分析service可知是去/etc/init.d目录下执行相关程序。service和chkconfig结合使用。
服务配置文件存放目录/etc/init.d/

例如

# 启动sshd服务
service sshd start
# 设置sshd服务开机启动
chkconfig sshd start
2、systemd
centos7版本中使用了systemd，systemd同时兼容service,对应的命令就是systemctl
Systemd 是 Linux 系统中最新的初始化系统（init），它主要的设计目标是克服 sysvinit 固有的缺点，提高系统的启动速度
使用systemd的目的是获取更快的启动速度。
为了减少系统启动时间，systemd的目标是
尽可能启动较少的进程
尽可能将更多进程并发启动
可以去查看系统进程的pid，initd的pid是0，如果支持systemd的系统的systemd进程pid为1

systemd把不同的资源称为Unit
每一个 Unit 都有一个配置文件，告诉 Systemd 怎么启动这个 Unit
存放目录：/etc/systemd/system和/usr/lib/systemd/system

对于有先后依赖关系的任务
systemctl融合service和chkconfig功能
systemctl的使用例如

# 开启服务
systemctl start sshd.service
# 设置开机启动
systemctl enable sshd.service
# 本质上是建立一个软链接 ln -s /usr/lib/systemd/system/sshd.service /etc/systemd/system/multi-user.target.wants/sshd.service

Systemd 的使用
下面针对技术人员的不同角色来简单地介绍一下 systemd 的使用。本文只打算给出简单的描述，让您对 systemd 的使用有一个大概的理解。具体的细节内容太多，即无法在一篇短文内写全，本人也没有那么强大的能力。还需要读者自己去进一步查阅 systemd 的文档。
系统软件开发人员
开发人员需要了解 systemd 的更多细节。比如您打算开发一个新的系统服务，就必须了解如何让这个服务能够被 systemd 管理。这需要您注意以下这些要点：
后台服务进程代码不需要执行两次派生来实现后台精灵进程，只需要实现服务本身的主循环即可。
不要调用 setsid()，交给 systemd 处理
不再需要维护 pid 文件。
Systemd 提供了日志功能，服务进程只需要输出到 stderr 即可，无需使用 syslog。
处理信号 SIGTERM，这个信号的唯一正确作用就是停止当前服务，不要做其他的事情。
SIGHUP 信号的作用是重启服务。
需要套接字的服务，不要自己创建套接字，让 systemd 传入套接字。
使用 sd_notify()函数通知 systemd 服务自己的状态改变。一般地，当服务初始化结束，进入服务就绪状态时，可以调用它。
Unit 文件的编写
对于开发者来说，工作量最大的部分应该是编写配置单元文件，定义所需要的单元。
举例来说，开发人员开发了一个新的服务程序，比如 httpd，就需要为其编写一个配置单元文件以便该服务可以被 systemd 管理，类似 UpStart 的工作配置文件。在该文件中定义服务启动的命令行语法，以及和其他服务的依赖关系等。
此外我们之前已经了解到，systemd 的功能繁多，不仅用来管理服务，还可以管理挂载点，定义定时任务等。这些工作都是由编辑相应的配置单元文件完成的。我在这里给出几个配置单元文件的例子。
下面是 SSH 服务的配置单元文件，服务配置单元文件以.service 为文件名后缀。
#cat /etc/system/system/sshd.service
[Unit]
Description=OpenSSH server daemon
[Service]
EnvironmentFile=/etc/sysconfig/sshd
ExecStartPre=/usr/sbin/sshd-keygen
ExecStart=/usrsbin/sshd –D $OPTIONS
ExecReload=/bin/kill –HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=42s
[Install]
WantedBy=multi-user.target
文件分为三个小节。第一个是[Unit]部分，这里仅仅有一个描述信息。第二部分是 Service 定义，其中，ExecStartPre 定义启动服务之前应该运行的命令；ExecStart 定义启动服务的具体命令行语法。第三部分是[Install]，WangtedBy 表明这个服务是在多用户模式下所需要的。
那我们就来看下 multi-user.target 吧：
#cat multi-user.target
[Unit]
Description=Multi-User System
Documentation=man.systemd.special(7)
Requires=basic.target
Conflicts=rescue.service rescure.target
After=basic.target rescue.service rescue.target
AllowIsolate=yes
[Install]
Alias=default.target
第一部分中的 Requires 定义表明 multi-user.target 启动的时候 basic.target 也必须被启动；另外 basic.target 停止的时候，multi-user.target 也必须停止。如果您接着查看 basic.target 文件，会发现它又指定了 sysinit.target 等其他的单元必须随之启动。同样 sysinit.target 也会包含其他的单元。采用这样的层层链接的结构，最终所有需要支持多用户模式的组件服务都会被初始化启动好。
在[Install]小节中有 Alias 定义，即定义本单元的别名，这样在运行 systemctl 的时候就可以使用这个别名来引用本单元。这里的别名是 default.target，比 multi-user.target 要简单一些。。。
此外在/etc/systemd/system 目录下还可以看到诸如*.wants 的目录，放在该目录下的配置单元文件等同于在[Unit]小节中的 wants 关键字，即本单元启动时，还需要启动这些单元。比如您可以简单地把您自己写的 foo.service 文件放入 multi-user.target.wants 目录下，这样每次都会被默认启动了。
最后，让我们来看看 sys-kernel-debug.mout 文件，这个文件定义了一个文件挂载点：

#cat sys-kernel-debug.mount
[Unit]
Description=Debug File Syste
DefaultDependencies=no
ConditionPathExists=/sys/kernel/debug
Before=sysinit.target
[Mount]
What=debugfs
Where=/sys/kernel/debug
Type=debugfs
这个配置单元文件定义了一个挂载点。挂载配置单元文件有一个[Mount]配置小节，里面配置了 What，Where 和 Type 三个数据项。这都是挂载命令所必须的，例子中的配置等同于下面这个挂载命令：
mount –t debugfs /sys/kernel/debug debugfs
配置单元文件的编写需要很多的学习，必须参考 systemd 附带的 man 等文档进行深入学习。希望通过上面几个小例子，大家已经了解配置单元文件的作用和一般写法了。

我的总结：systemd功能更加强大，是某些linux发行版用来代替系统启动init和service的一个新的系统工具，同时还没有失去兼容对init的兼容性。

/etc/init.d 是sysVinit服务的启动方式,对于一些古老的系统或者服务 使用这个.
service 也是sysVinit, 比/etc/init.d先进一点,底层还是调用/etc/init.d
systemctl 是systemD命令的主要方式, 尽管一些老的系统或者命令不支持systemctl, 但是systemctl最后会逐渐的替代其他的命令方式的, 能用这个就优先用这个,是最时尚/方便的

systemctl的config文件存在 /lib/systemd/system/xxxxxx.serive

我们知道，Docker运行一个容器起来的时候，只是为你提供特定的文件系统层和进程隔离，它给你一个VM的感觉却并不是VM，所以你可能偶尔会想要像在物理机那样使用systemctl start|status|stop来管理服务进程，然后你通常会看到
Failed to get D-Bus connection: Operation not permitted
这个错误。
原因很简单：

你需要启动systemd进程
你需要特权
所以你如果想要一个可以使用Systemd的容器，你可以尝试这样启动容器:

1.解决办法一：给权限
cat /etc/redhat-release 
//CentOS Linux release 7.2.1511 (Core) 
docker run -tdi --privileged centos init
在容器中，你可以使用systemd管理服务进程了:

yum install -y vsftpd
systemctl start vsftpd
systemctl status vsftpd
2.解决办法二：init.d
除此之外你也可以通过其他启动软件的方式: init.d来达到 systemctl 的效果：
比如要启动 mysql 服务：/etc/init.d/mysql start

Linux 服务管理两种方式service和systemctl

1.service命令
service命令其实是去/etc/init.d目录下，去执行相关程序

# service命令启动redis脚本
service redis start
# 直接启动redis脚本
/etc/init.d/redis start
# 开机自启动
update-rc.d redis defaults
其中脚本需要我们自己编写

2.systemctl命令
systemd是Linux系统最新的初始化系统(init),作用是提高系统的启动速度，尽可能启动较少的进程，尽可能更多进程并发启动。
systemd对应的进程管理命令是systemctl

1)systemctl命令兼容了service
即systemctl也会去/etc/init.d目录下，查看，执行相关程序

systemctl redis start
systemctl redis stop
# 开机自启动
systemctl enable redis
2)systemctl命令管理systemd的资源Unit
systemd的Unit放在目录/usr/lib/systemd/system(Centos)或/etc/systemd/system(Ubuntu)


主要有四种类型文件.mount,.service,.target,.wants

.mount文件


.mount文件定义了一个挂载点，[Mount]节点里配置了What,Where,Type三个数据项
等同于以下命令：

mount -t hugetlbfs /dev/hugepages hugetlbfs
.service文件


.service文件定义了一个服务，分为[Unit]，[Service]，[Install]三个小节
[Unit]
Description:描述，
After：在network.target,auditd.service启动后才启动
ConditionPathExists: 执行条件

[Service]
EnvironmentFile:变量所在文件
ExecStart: 执行启动脚本
Restart: fail时重启

[Install]
Alias:服务别名
WangtedBy: 多用户模式下需要的

.target文件


.target定义了一些基础的组件，供.service文件调用

.wants文件


.wants文件定义了要执行的文件集合，每次执行，.wants文件夹里面的文件都会执行

1、service是一个脚本命令，分析service可知是去/etc/init.d目录下执行相关程序。service和chkconfig结合使用。
服务配置文件存放目录/etc/init.d/

例如

# 启动sshd服务
service sshd start
# 设置sshd服务开机启动
chkconfig sshd start
2、systemd
centos7版本中使用了systemd，systemd同时兼容service,对应的命令就是systemctl
Systemd 是 Linux 系统中最新的初始化系统（init），它主要的设计目标是克服 sysvinit 固有的缺点，提高系统的启动速度
使用systemd的目的是获取更快的启动速度。
为了减少系统启动时间，systemd的目标是
尽可能启动较少的进程
尽可能将更多进程并发启动
可以去查看系统进程的pid，initd的pid是0，如果支持systemd的系统的systemd进程pid为1

systemd把不同的资源称为Unit
每一个 Unit 都有一个配置文件，告诉 Systemd 怎么启动这个 Unit
存放目录：/etc/systemd/system和/usr/lib/systemd/system

对于有先后依赖关系的任务
systemctl融合service和chkconfig功能
systemctl的使用例如

# 开启服务
systemctl start sshd.service
# 设置开机启动
systemctl enable sshd.service
# 本质上是建立一个软链接 ln -s /usr/lib/systemd/system/sshd.service /etc/systemd/system/multi-user.target.wants/sshd.service
1
2
3
4
5
转载自
http://www.ibm.com/developerworks/cn/linux/1407_liuming_init3/
Systemd 的使用
下面针对技术人员的不同角色来简单地介绍一下 systemd 的使用。本文只打算给出简单的描述，让您对 systemd 的使用有一个大概的理解。具体的细节内容太多，即无法在一篇短文内写全，本人也没有那么强大的能力。还需要读者自己去进一步查阅 systemd 的文档。
系统软件开发人员
开发人员需要了解 systemd 的更多细节。比如您打算开发一个新的系统服务，就必须了解如何让这个服务能够被 systemd 管理。这需要您注意以下这些要点：
后台服务进程代码不需要执行两次派生来实现后台精灵进程，只需要实现服务本身的主循环即可。
不要调用 setsid()，交给 systemd 处理
不再需要维护 pid 文件。
Systemd 提供了日志功能，服务进程只需要输出到 stderr 即可，无需使用 syslog。
处理信号 SIGTERM，这个信号的唯一正确作用就是停止当前服务，不要做其他的事情。
SIGHUP 信号的作用是重启服务。
需要套接字的服务，不要自己创建套接字，让 systemd 传入套接字。
使用 sd_notify()函数通知 systemd 服务自己的状态改变。一般地，当服务初始化结束，进入服务就绪状态时，可以调用它。
Unit 文件的编写
对于开发者来说，工作量最大的部分应该是编写配置单元文件，定义所需要的单元。
举例来说，开发人员开发了一个新的服务程序，比如 httpd，就需要为其编写一个配置单元文件以便该服务可以被 systemd 管理，类似 UpStart 的工作配置文件。在该文件中定义服务启动的命令行语法，以及和其他服务的依赖关系等。
此外我们之前已经了解到，systemd 的功能繁多，不仅用来管理服务，还可以管理挂载点，定义定时任务等。这些工作都是由编辑相应的配置单元文件完成的。我在这里给出几个配置单元文件的例子。
下面是 SSH 服务的配置单元文件，服务配置单元文件以.service 为文件名后缀。
#cat /etc/system/system/sshd.service
[Unit]
Description=OpenSSH server daemon
[Service]
EnvironmentFile=/etc/sysconfig/sshd
ExecStartPre=/usr/sbin/sshd-keygen
ExecStart=/usrsbin/sshd –D $OPTIONS
ExecReload=/bin/kill –HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=42s
[Install]
WantedBy=multi-user.target
文件分为三个小节。第一个是[Unit]部分，这里仅仅有一个描述信息。第二部分是 Service 定义，其中，ExecStartPre 定义启动服务之前应该运行的命令；ExecStart 定义启动服务的具体命令行语法。第三部分是[Install]，WangtedBy 表明这个服务是在多用户模式下所需要的。
那我们就来看下 multi-user.target 吧：
#cat multi-user.target
[Unit]
Description=Multi-User System
Documentation=man.systemd.special(7)
Requires=basic.target
Conflicts=rescue.service rescure.target
After=basic.target rescue.service rescue.target
AllowIsolate=yes
[Install]
Alias=default.target
第一部分中的 Requires 定义表明 multi-user.target 启动的时候 basic.target 也必须被启动；另外 basic.target 停止的时候，multi-user.target 也必须停止。如果您接着查看 basic.target 文件，会发现它又指定了 sysinit.target 等其他的单元必须随之启动。同样 sysinit.target 也会包含其他的单元。采用这样的层层链接的结构，最终所有需要支持多用户模式的组件服务都会被初始化启动好。
在[Install]小节中有 Alias 定义，即定义本单元的别名，这样在运行 systemctl 的时候就可以使用这个别名来引用本单元。这里的别名是 default.target，比 multi-user.target 要简单一些。。。
此外在/etc/systemd/system 目录下还可以看到诸如*.wants 的目录，放在该目录下的配置单元文件等同于在[Unit]小节中的 wants 关键字，即本单元启动时，还需要启动这些单元。比如您可以简单地把您自己写的 foo.service 文件放入 multi-user.target.wants 目录下，这样每次都会被默认启动了。
最后，让我们来看看 sys-kernel-debug.mout 文件，这个文件定义了一个文件挂载点：

#cat sys-kernel-debug.mount
[Unit]
Description=Debug File Syste
DefaultDependencies=no
ConditionPathExists=/sys/kernel/debug
Before=sysinit.target
[Mount]
What=debugfs
Where=/sys/kernel/debug
Type=debugfs
这个配置单元文件定义了一个挂载点。挂载配置单元文件有一个[Mount]配置小节，里面配置了 What，Where 和 Type 三个数据项。这都是挂载命令所必须的，例子中的配置等同于下面这个挂载命令：
mount –t debugfs /sys/kernel/debug debugfs
配置单元文件的编写需要很多的学习，必须参考 systemd 附带的 man 等文档进行深入学习。希望通过上面几个小例子，大家已经了解配置单元文件的作用和一般写法了。

我的总结：systemd功能更加强大，是某些linux发行版用来代替系统启动init和service的一个新的系统工具，同时还没有失去兼容对init的兼容性
