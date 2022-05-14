---
title: 使用Phabricator做为Code Review工具
layout: post
category: web
author: 夏泽民
---
目录

0x10 概述
0x20 我的应用环境
0x30 路线图
0x40 安装
0x41 LNMP环境的安装
0x42 Phabricator源码下载及运行
0x50 配置
0x51 解决基本的配置问题
0x52 设置用户登录认证方式
0x53 设置邮件发送服务参数
0x54 配置代码仓库访问方式：SSH/HTTP
0x60 使用Phabricator进行Code Review
0x61 Phabricator Code Review工作流
0x62 进行Code Review所用工具
0x63 配置进行Code Review
0x70 与GitHub集成
0x80 与Jenkins集成
0x90 结束语
0xA0 Q/A
0x10 概述
<!-- more -->
Phabricator是一套基于Web的软件开发协作工具，包括代码审查工具Differential，资源库浏览器Diffusion，变更监测工具Herald，Bug跟踪工具Maniphest和维基工具Phriction。Phabricator可与Git、Mercurial、Subversion集成使用。
Phabricator是开源软件，可在Apache许可证第2版下作为自由软件分发。
Phabricator最初是Facebook的一个内部工具，主要开发者为Evan Priestley。Evan Priestley离开Facebook后，在名为Phacility的新公司继续Phabricator的开发。
官网：https://www.phacility.com/

官网中的文档很多很全，但是如果对这个工具不太了解，或者对于Code Review不太了解的话，读起来可能会觉得没有头绪。这篇文章就自己的安装及使用过程做一个梳理，对于同样想用这个工具的读者，或许起到一些帮助作用。
对于本文中的一些章节，如果在官方文档有所提及，我会把官方文档地址附上，读者可以阅读一下官方文档，因为他们的语言和表达更优秀。

0x20 我的应用环境

我：一个做了很多年Android的程序员啊 （所以当这篇文章有幸被所涉及领域的专家看到，又发现有的部分有所不妥，或者有更好的想法时，请主动联系我改进，多谢~~~ 有任何问题，欢迎评论交流）
主机：Ubuntu 14.04 PC一台
网络：内网
代码库：GitHub
CI：Jenkins
0x30 路线图


安装和使用路线大致如此图，下面开始详细说明。

0x40 安装

Phabricator是一个基于Web的工具软件，使用PHP语言编写的，为了能让他运行起来，我们需要搭建一个LNMP（Linux，Nginx，MySQL，PHP）的Web Server环境。搭建完LNMP的环境后，下载Phabricator源码，配置后即可使用。

先看我
如果你不想读那么多文字，在Ubuntu环境下可以试试下面的脚本，这个脚本可以安装LNMP环境和下载Phabricator源码，执行完脚本并成功后，跳到 0x42 Phabricator源码下载及运行 查看如何让Phabricator跑起来。
{% highlight bash linenos %}
#!/bin/bash
confirm() {
  echo "Press RETURN to continue, or ^C to cancel.";
  read -e ignored
}
GIT='git'
LTS="Ubuntu 10.04"
ISSUE='cat /etc/issue'
if [[ $ISSUE != Ubuntu* ]]
then
  echo "This script is intended for use on Ubuntu, but this system appears";
  echo "to be something else. Your results may vary.";
  echo
  confirm
elif [[ 'expr match "$ISSUE" "$LTS"' -eq ${#LTS} ]]
then
  GIT='git-core'
fi
echo "PHABRICATOR UBUNTU INSTALL SCRIPT";
echo "This script will install Phabricator and all of its core dependencies.";
echo "Run it from the directory you want to install into.";
echo
ROOT='pwd'
echo "Phabricator will be installed to: ${ROOT}.";
confirm
echo "Testing sudo..."
sudo true
if [ $? -ne 0 ]
then
  echo "ERROR: You must be able to sudo to run this script.";
  exit 1;
fi;
echo "Installing dependencies: git, nginx, mysql, php...";
echo
set +x
sudo apt-get -qq update
sudo apt-get install \
  $GIT nginx mysql-server dpkg-dev \
  php5 php5-mysql php5-gd php5-dev php5-curl php-apc php5-cli php5-json
# Enable mod_rewrite
sudo a2enmod rewrite
HAVEPCNTL='php -r "echo extension_loaded('pcntl');"'
if [ $HAVEPCNTL != "1" ]
then
  echo "Installing pcntl...";
  echo
  apt-get source php5
  PHP5='ls -1F | grep '^php5-.*/$''
  (cd $PHP5/ext/pcntl && phpize && ./configure && make && sudo make install)
else
  echo "pcntl already installed";
fi
if [ ! -e libphutil ]
then
  git clone https://github.com/phacility/libphutil.git
else
  (cd libphutil && git pull --rebase)
fi
if [ ! -e arcanist ]
then
  git clone https://github.com/phacility/arcanist.git
else
  (cd arcanist && git pull --rebase)
fi
if [ ! -e phabricator ]
then
  git clone https://github.com/phacility/phabricator.git
else
  (cd phabricator && git pull --rebase)
fi
echo
echo
echo "Install probably worked mostly correctly. Continue with the 'Configuration Guide':";
echo
echo "    https://secure.phabricator.com/book/phabricator/article/configuration_guide/";
echo
echo "You can delete any php5-* stuff that's left over in this directory if you want.";
{% endhighlight %}
关于安装和配置，官方文档中有所提及（官方介绍的是 LAMP），请参考

https://secure.phabricator.com/book/phabricator/article/installation_guide/
0x41 LNMP环境的安装

安装Linux
关于Linux的安装，这里就不说了。

安装Nginx

sudo apt-get install nginx
安装完成后，Nginx 的配置文件存放在 /etc/nginx 目录下。使用下面的命令可以启动Nginx

sudo service nginx start
在安装完并启动后，可以使用浏览器访问 http://127.0.0.1 试试是否可以跳转到Nginx欢迎页面


安装MySQL

sudo apt-get install mysql-server
在安装过程中，会两次提示输入 root 用户密码。
在安装完成后，打开终端，使用以下命令登录MySQL

mysql -u root -p
安装PHP
Phabricator需要 PHP 5.2 或者更高版本，但是 不支持 PHP 7 。
可以使用以下命令安装PHP

sudo apt-get install -y php5 php5-fpm php5-mysql
安装完成后，可使用以下命令查看是否安装成功

php -v
安装成功后，输出类似以下信息


安装其它
如果你使用 git 来管理代码库的话，你还需要安装 git

sudo apt-get install git
一些必要的PHP扩展

mbstring, iconv, mysql (or mysqli), curl, pcntl（这些扩展一般会以 "php-mysql" 或 "php5-mysql" 方式使用）
一些可选的PHP扩展

gd, apc（官方文档中有详细的介绍和安装说明）, xhprof（如果你想自己开发Phabricator的话，你需要安装这个，官方文档中有详细的介绍说明）
0x42 Phabricator源码下载及运行

在成功安装LNMP环境后，需要下载Phabricator的源码并配置让它跑起来。

源码下载
在你想要存放Phabricator源码的位置（假设为 ./path_to_pha），执行这些命令

git clone https://github.com/phacility/libphutil.git
git clone https://github.com/phacility/arcanist.git
git clone https://github.com/phacility/phabricator.git
或者，你也可以直接点击上面的链接去GitHub下载压缩包，下载完成后解压。

Nginx配置
在下载完成后，我们需要配置Nginx，让Phabricator跑起来。假设你想为Phabricator分配这个域名：http://pha.example.com
在 /etc/nginx/conf.d 目录下创建文件 pha.example.com.conf，存放Phabricator代理配置信息，以下为我的文件内容（注意把 你存放Phabricator的路径 改为你的实际路径）

server {
    listen       80;
    server_name  pha.example.com; 
    location / {
        index index.php;
        rewrite ^/(.*)$ /index.php?__path__=/$1 last;
    }

    #error_page  404              /404.html;

    # redirect server error pages to the static page /50x.html
    #
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }

     # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
    #
    location ~ \.php$ {
        root           /你存放Phabricator的路径/phabricator/webroot;
        fastcgi_pass   127.0.0.1:9000;
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME  /你存放Phabricator的路径/phabricator/webroot$fastcgi_script_name;
        include        fastcgi_params;
    }
}
配置完成后，重启Nginx

sudo service nginx restart
然后在你的 hosts 文件中，加入 pha.example.com 对应的IP

127.0.0.1 pha.example.com
打开浏览器，访问 http://pha.example.com，会跳转到Phabricator用户注册界面，在这个界面注册的第一个用户，将会成为管理员用户。

0x50 配置

0x51 解决基本的配置问题

使用管理员账号登录，左上角会出现黄色感叹号图标，提示有一些配置问题未解决



这些问题基本都是关于一些参数的设置。点击每一个问题，显示的界面中会有很详细的关于这个问题的描述，和如何解决。

0x52 设置用户登录认证方式

使用管理员账号登录，在左侧的菜单中选择 Auth ，然后点击右上侧 Add Provider，在列表中选则你需要的认证方式。


我选择是 Username/Password 的方式，即用户自己注册Phabricator账号。为了保障安全，我设置了只允许公司邮箱地址注册：Config ---> Core Settings ---> Authentication ---> auth.email-domains。你还可以选择 auth.require-approval ，即新注册用户需要管理员批准。

0x53 设置邮件发送服务参数

首先，配置 mail-adapter （邮件发送方式）：Config ---> Core Settings ---> Mail ---> metamta.mail-adapter，我选择的是 PhabricatorMailImplementationPHPMailerAdapter ，通过SMTP的方式发送邮件。在选择完之后，需要设置SMTP服务器地址、账号和密码：Config ---> Core Settings ---> PHPMailer ---> metamta.mail-adapter，根据你自己邮箱的配置，相应的设置 phpmailer.smtp-host、phpmailer.smtp-port、phpmailer.smtp-protocol、phpmailer.smtp-user、phpmailer.smtp-password、phpmailer.smtp-encoding 。

0x54 配置代码仓库访问方式：SSH/HTTP

SSH
(如果你不打算允许使用SSH的方式访问代码仓库的话，请忽略这部分)
1）配置用户账号
Phabricator需要三个用户账号（三种用户身份）：两个用于基本运行，一个用于配置SSH访问。这些账号是指Phabricator所运行服务器系统的账号，不是Phabricator用户账号。
三个账号分别是：
www-user：Phabricator Web服务器运行身份。
daemon-user ：daemons （守护进程）运行身份。这个账号是唯一直接与代码仓库交互的账号，其它账号需要切换到这个账号身份（sudo）才能操作代码仓库。
vcs-user：我们需要以这个账号SSH连接Phabricator。
如果你的服务器系统中现在没有这三个账号，需要创建：
www-user：大部分情况下，这个账号已经存在了，我们不需要理这个账号。
daemon-user ：一般情况下，我们直接使用 root 账号，因为会需要很多权限（当然这可能不安全）。
vcs-user：可以使用系统中现有的一个用户账号，直接创建一个就叫 vcsuser。当用户克隆仓库的时候，需要使用类似 vcsuser@pha.example.com 的URI。
2）配置Phabricator
首先，设置 phd.user 为 daemon-user（root）

./path_to_pha/bin/config set phd.user root
重启 daemons 以确认这个配置工作正常

./path_to_pha/bin/phd restart
然后，配置SSH用户账号vcs-user（vcsuser 或其它你想用的用户）

./path_to_pha/bin/config set diffusion.ssh-user vcsuser
3）配置 Sudo
www-user 和 vcs-user 需要能够使用 sudo 切换到 daemon-user 用户身份才能与仓库交互，所以我们需要配置更改系统的 sudo 配置。
直接编辑 /etc/sudoers 或者在 /etc/sudoers.d 下创建一个新文件，然后把这些内容写到文件内容中

www-user ALL=(root) SETENV: NOPASSWD: /usr/lib/git-core/git, /usr/bin/git, /var/lib/git, /usr/lib/git-core/git-http-backend, /usr/bin/ssh, /etc/ssh, /etc/default/ssh, /etc/init.d/ssh
vcs-user ALL=(root) SETENV: NOPASSWD: /bin/sh, /usr/bin/git-upload-pack, /usr/bin/git-receive-pack
当然，别忘了把 www-user 和 vcs-user 替换为你实际对应的用户。
接下来，看看你文件中是不是有这行

Defaults requiretty
如果有的话，请用 # 注释掉。

4）其它SSH配置
我们还需要查看这两个文件 /etc/shadow 和 /etc/passwd 中 vcs-user 对应的配置是否正确。
打开 /etc/shadow 文件，找到 vcs-user 对应的那行，看一下第二个字段（密码），是不是 !! ，如果是，请改为 空值（什么都不写） 或者 NP 。
打开 /etc/passwd 文件，找到 vcs-user 对应的那行，如果有类似于这样的配置 /bin/false ，请修改为 /bin/sh，否则 sshd 无法执行命令。

5）配置SSHD端口
注意：Phabricator运行的服务器系统中 sshd 的版本 必须高于 6.2。
假设我们把Phabricator使用的sshd端口设置为 22，这样做的好处是我们不需要在仓库的URI中加入端口号，类似ssh://vcs-user@pha.example.com/xxx/xxx/xxx.git。当然，如果这样做需要我们更改系统已存在的sshd配置改为其它端口。下面来看一下配置的三个步骤：
i）创建脚本 phabricator-ssh-hook.sh，并且把这个脚本放到类似 /usr/libexec/phabricator-ssh-hook.sh 的目录中（我直接放在 /etc/ssh/ 中，后面会要求变更这个脚本和它的父文件夹所有者，所以这个脚本和它的父文件夹所在的文件夹的所有者不正确的话可能会导致这个脚本执行失败），脚本内容如下

#!/bin/sh

# NOTE: Replace this with the username that you expect users to connect with.
VCSUSER="vcs-user"

# NOTE: Replace this with the path to your Phabricator directory.
ROOT="/path_to_pha"

if [ "$1" != "$VCSUSER" ];
then
  exit 1
fi

exec "$ROOT/bin/ssh-auth" $@
注意把 VCSUSER 替换为你实际的用户，把 ROOT 值替换为你Phabricator源码路径。
创建完脚本后，需要把脚本和它的父文件夹所有者改为 root，并且赋予脚本 755 权限：

sudo chown root /path/to/somewhere/
sudo chown root /path/to/somewhere/phabricator-ssh-hook.sh
sudo chmod 755 /path/to/somewhere/phabricator-ssh-hook.sh
如果你不这么做，sshd 会拒绝执行 hook。

ii）为Phabricator创建 sshd_config
在 /etc/ssh 中创建文件名类似 sshd_config.phabricator 的文件，文件内容如下：

# NOTE: You must have OpenSSHD 6.2 or newer; support for AuthorizedKeysCommand
# was added in this version.

# NOTE: Edit these to the correct values for your setup.

AuthorizedKeysCommand /你的脚本路径/phabricator-ssh-hook.sh
AuthorizedKeysCommandUser vcs-user
AllowUsers vcs-user

# You may need to tweak these options, but mostly they just turn off everything
# dangerous.

Port 你配置的端口号
Protocol 2
PermitRootLogin no
AllowAgentForwarding no
AllowTcpForwarding no
PrintMotd no
PrintLastLog no
PasswordAuthentication no
AuthorizedKeysFile none

PidFile /var/run/sshd-phabricator.pid
注意把 AuthorizedKeysCommand 值替换为你在上一步中脚本实际路径，把 AuthorizedKeysCommandUser 和 AllowUsers 替换为你实际的用户，把 Port 替换为你想配置的端口号。如果你的 Port 值为 22，在你进行下面的操作之前，请查看当前系统中 22 端口是已否占用

sudo netstat -atlunp | grep ssh
如果已经被占用，请修改使用 22 端口的 sshd 配置，一般它们会在 /etc/ssh 下，名称类似 sshd_config，修改完成后，请重启 ssh 服务

sudo /etc/init.d/ssh restart
在完成上面的步骤后，我们来启动Phabricator的 ssh 服务

sudo /path/to/sshd -f /你的Phabricator sshd配置路径/sshd_config.phabricator
一般情况下，sshd 路径为 /usr/sbin。
在启动后，我们需要验证以下配置是否有效：
首先，请把你的公钥添加到Phabricator自己的账号中（你可以自己注册一个新的账号），注册完成后登录，然后 点击你的头像 ---> 左侧菜单面板 Manage ---> 右侧菜单面板 Edit Settings ---> 左侧菜单面板 SSH Public Keys ---> 右上角 SSH Key Actions ---> Upload Public Key


上传公钥后，执行下面的命令

echo {} | ssh vcs-user@phabricator.yourcompany.com conduit conduit.ping
如果出现类似下面的结果，说明配置有效

{"result":"phabricator.yourcompany.com","error_code":null,"error_info":null}
如果没有出现别的情况，请参考官方文档 Troubleshooting SSH 部分，官方文档地址如下

https://secure.phabricator.com/book/phabricator/article/diffusion_hosting/
接下来，看一下如何配置 HTTP

HTTP
首先，请确认Phabricator的配置项 diffusion.allow-http-auth 设置为 true。可以在 左侧菜单面板 All Setttings 中查找 diffusion.allow-http-auth ，点击之后可设置，请设置为 Allow HTTP Basic Auth。
然后，所有用户需要使用 HTTP 访问仓库之前，需要设置自己的密码：点击你的头像 ---> 左侧菜单面板 Manage ---> 右侧菜单面板 Edit Settings ---> 左侧菜单面板 VCS Password

强烈建议不要把这个密码设置为你的Phabricator登录密码，因为 vcs 密码很容易泄露。
一般来说，不需要其它配置就可以使用 HTTP 了，如果有问题，请参考官方文档 Troubleshooting HTTP 部分

https://secure.phabricator.com/book/phabricator/article/diffusion_hosting/
配置完仓库访问方式后，我们来看一下如何使用 Phabricator 进行 Code Review。

0x60 使用Phabricator进行Code Review

在进行 Code Review 实践前，先说一些理论方面的东西（开头和 0x61 ，不喜欢可绕过）
Code Review，有时候就像打架一样：我提交了变更，你说不行，要修改；我又提交了一次，你说还是不行，还要改。我不知道你究竟要怎样，你也不知道我感觉受到了打击有多不爽。所以，大家需要对Code Review这件事抱有开放的态度：

为什么我的代码需要其他人审查？
因为我不是神，我会制造Bug，我会当局者迷。
为什么我要审查其他人的代码？
因为我要对我们的团队负责，我要保证我们产品的质量，我可能会看到他人代码的Bug，在这些Bug显示出它们的"威力"前，把它们弄死。
Code Review这件事，旨在创造一个共进的团队氛围（交流和技术等），在产品交付给用户（包括我们的测试人员）前，保证产品的质量。

在了解如何使用Phabricator进行Code Review前，我们先了解一下Phabricator Code Review的流程，对其有一个整体上的了解。

0x61 Phabricator Code Review工作流

Phabricator提供两种Code Review的方式：pre-push，post-push
pre-push 是指审查发生在变更发布前；post-push 是指审查发生在变更已经被发布或者正在发布。
这里我们认为 pre-push 的方式更适合，所以接下来说一下 pre-push 的工作流：

Write, Review, Merge, Publish
从这篇文章，我了解到了这个流程

https://secure.phabricator.com/phame/post/view/766/write_review_merge_publish_phabricator_review_workflow/
如果你之前用过其它的Code Review工具，可能会对这样的流程感到不习惯。在其它工具中，变更（代码，资源文件或其它）会经历这样一个流程： Write, Publish, Review, Merge。首先，你做出一些变更（Write），然后把他们推送到远程仓库（Publish）等待审查者审查。一旦这些代码被审查（Review）并通过，变更会合并（Merge）到一个指定的功能分支。在这个流程中，被合并的变更恰好是被推送的变更（这句话有点模模糊糊，不痛不痒，接下来我们看一下Phabricator的流程，也许会清晰很多）。
接下来，我们看一下Phabricator略有不同的工作流：Write, Review, Merge, Publish。像上面一样，开始的时候，你做出一些变更。但是，接下来的流程就不一样了。
Phabricator认为在开发过程中审查（Review）是一个重要的步骤，对于那些没有审查过的变更，是不可以发布的。
理论上来说，没有审查过的变更不算数：这些变更可能只是临时的，易变的。可能方法上不对，可能缺少来龙去脉，可能根本就是解决错误的问题，等等。审查的参考基础是建立在开发人员和审查人员拥有一个共同认可的变更处理方式，并且这种处理方式是开发过程所有参与人员（项目管理、产品、开发）都期望的，而不仅仅是仅仅做到最终的产品看起来没问题。直到变更经过了这样的审查，我们才能得到稳固的版本。
这样的工作流跟其它工具的审查流程没有实质上的技术区别，但是存在明显的社交活动上的不同：由于变更必须经过审查才能被合并、发布，变更作者需要根据反馈对变更进行调整。另外，审查者根据粗略的草图（所有开发参与人员共同认可的变更方式）进行反馈，而不是简单的批判一件已经完成的变更工作。
Phabricator和其它工具的工作流都有着同样的目的：未审查的代码都只是临时的变更，没有长久或者明显的价值，直到通过审查。
Phabricator工作流的第二步是审查（Review）,审查的对象是还没有发布的变更。没有发布的变更被发送到Phabricator等待被审查（通常我们使用 arc diff命令发送审查请求），然后审查者做出反馈。变更作者根据反馈进行修改，在修改过程中，作者不必担心版本、解决方式这些事情。作者可以自由的复位、使用、移除或者舍弃老的变更。在从变更提交审核到审核者反馈，以及作者再次修改整个过程中，没有那种审核者把作者推入一个必须接受或者只能做少量改变的默认发布状态。
一旦通过审核，变更会被合并（Merge）和发布（Publish）（通常，这两个步骤由一个命令完成 arc land）。
这里，Phabricator也与那些先 Publish 的工具不同：默认情况下，Phabricator会舍弃到达最后变更前的所做的中间过程，把最后变更的整个过程压缩成一次提交。总体来说，这意味着舍弃checkpoint commits, rebases, squash-merges, 并且把整个变更过程做为一次 fast-forward commit 提交到目标分支。
Phabricator在一定程度上能做到这些，是因为：什么都没有被发布，所以这种工作流可以以任何想要的方式发布变更。
有了这些，我们可以以我们想要的版本自由的rebase，fast-commit，这些是Phabricator默认的行为。

0x62 进行Code Review所用工具

做为一般用户，常用的工具有两个 Differential 和 Arcanist 。

Differential-审查代码的工作台

我们在这里查看变更审查情况，对变更进行审查或评论等操作。
这是某次变更界面操作部分截图


做为 审查人，可进行的操作有：
Comment：说点什么。可以针对某行代码进行评论，直接点击行号即可
Accept Revision：接受变更，这哥们代码写得不错，不需要改
Request Changes：不行，还要改
Resign as Reviewer：重新指定审查代码的人
Commandeer Revision：字面意思是将这个Revision据为己有的意思，实际上这个时候Reviewer的身份已经变为Owner的身份了，不能再进行Review了，但是Comment还是可以的
Add Reviewer：添加审查人
Add Subscribers：添加订阅者，CC

做为 作者，可进行的操作有：
Comment：说点什么。可以针对某行代码进行评论，直接点击行号即可
Abandon Revision：废除版本。废除后，这个版本就不需要再审核了
Plan Changes：计划变更，我自己发现了一些问题或者需求有变，正在改
Add Reviewer：添加其它审查人（除当前审查人外）
Add Subscribers：添加订阅者，CC

Arcanist - 命令交互
我们用这个工具提交变更和审查请求，对变更做出更改，或者在通过审查后发布到远程仓库分支中。
常用的命令有：
arc diff：发送变更详情和审查请求
arc land：推送变更（Git and Mercurial），当通过审查后使用这个命令
arc list：显示变更处理的情况
arc cover：查找最有可能审查变更的人
arc patch：给版本打补丁
arc export：从Differential下载补丁
arc amend：更新Git commit
arc commit：提交变更（SVN）
arc branch：查看Git branches更加详细的信息

在配置了 lint 和 unit test intergration后，可以用这些命令：
arc lint：静态代码检查
arc unit：单元测试

与其它工具交互：
arc upload：上传文件
arc download：下载文件
arc paste：创建和查看剪贴

还有一些高级功能：
arc call-conduit：执行 Conduit 方法
arc liberate：创建或更新 libphutil 库
arc shell-complete：激活 tab 补全

0x63 配置进行Code Review

一些基本的配置和安装 ---> 写代码 ---> 提交审查请求(arc diff) ---> 审查（Differential） ---> （审查通过后）合并提交（arc land）
一些基本的配置和安装
包括：
配置代码仓库（Diffusion）
把你本地的Git远程URL设置为Phabricator上代码仓库地址
安装Arcanist
配置Project信息

配置代码仓库（Diffusion）
在开始进行代码审查后，我们的代码是由Phabricator直接托管的，所以我们需要配置代码仓库。
使用管理员账号登录Phabricator，点击左侧面板菜单 Diffusion ，然后点击右上侧 Create Repository ，选择你所使用的 Repository 类型，填写 Name 等信息，在创建完成后即可使用。如果没有什么特殊的需求，不需要进行特别的配置，这里列举两种你可能遇到的打算开始使用 Phabricator 时的场景：

1、代码之前由 GitHub 或其它托管，现在我需要把之前的代码导入
点击 Manage Repository ，点击左侧 URIs，点击 Add New URI，填写GitHub或其它托管系统对应仓库的 URI ， I/O Type 选择 Observe，点击 Create Repository URI 添加新的 URI 。



在添加完新的 URI 后，你还需要点击 Set Credential 设置访问新的 URI 的认证方式。
如果你打算此时就开始使用 Phabricator ，请务必通知你的团队，暂停一下，不要再向GitHub等提交代码。如果你的 GitHub 等也设置了代码审查，请督促相关人员完成代码审查流程。 
稍等片刻，待 Phabricator 同步完之前的代码后，编辑你添加的 GitHub 或其它代码托管系统的 URI ，务必修改 I/O Type：
1）如果你不再需要使用之前的托管系统，选择 No I/O
2）如果你想继续把代码备份到之前的代码托管系统，选择 Mirror，这时， Phabricator 代码仓库的变更会覆盖推送到之前的代码托管系统
如果你不修改 I/O Type，向 Phabricator 代码仓库提交代码会失败，因为是只读的。
当然，对于导入之前的代码，还有别的方式，例如直接把本地的代码再次向 Phabricator 代码仓库再提交一次。

2、开始一个新的项目，创建一个新的仓库
参考第1种场景，在 Phabricator 创建代码仓库。如果你希望把代码备份到其它的托管系统，只需要添加对应的 URI，并且把 I/O Type 选为 Mirror。

把你本地的Git远程URL设置为Phabricator上代码仓库地址

git remote set-url 远程名称 新的url
安装Arcanist

https://secure.phabricator.com/book/phabricator/article/arcanist_quick_start/
配置Project信息
在你项目代码的根目录下，创建 .arcconfig 文件，内容如下：

{
"phabricator.uri" : "你Phabricator系统访问URL"
}
Windows系统下，创建类似这种文件名的文件可能很麻烦，可以使用这条命令创建：

arc set-config phabricator.uri "你Phabricator系统访问URL"
Windows系统下，还需要配置 Editor ，详情参考：

https://secure.phabricator.com/book/phabricator/article/arcanist_windows/
Arcanist 使用可参考：

https://secure.phabricator.com/book/phabricator/article/arcanist_quick_start/
https://secure.phabricator.com/book/phabricator/article/arcanist/
在进行完基本的配置和安装后，可以开始 Code Review 了。

写代码
当然，不只是代码可以被审查，图标等资源文件的变更也可以被审查。

提交审查请求（arc diff）
一般情况下，我们直接使用 arc diff 即可，默认情况下，Arcanist 会把本地分支的 HEAD 与远程对应分支的 HEAD 进行对比，并生成差异对比发送到 Phabricator。当你所做的修改没有 commit 时，会提示你进行 commit。
在一些情况下，我们并不希望与本地分支的 HEAD 进行比较，假设想要与上次的 commit 比较，上次 commit id 是 8ffc88dc05d31fffd28e3ff1129d1b8c321dffff，那么我们需要在 arc diff 后把这个 id 加上：arc diff 8ffc88dc05d31fffd28e3ff1129d1b8c321dffff。
执行这条命令时我们需要按照模板填写title（必填），summary（必填），Test Plan（必填，没有可写 N/A 之类的标识），Reviewers（必填，且必须为真实有效的用户名），Subscribers（可选），填写完成后关闭编辑器，Arcanist会自动提交审查请求。

审查（Differential）
做为 审查人，需要在 Differential 工作台完成代码审查工作，上面已经介绍了 Differential，这里就不再多说了。

（审查通过后）合并提交（arc land）
做为 作者，在代码审查通过后，需要使用 arc land 把变更发布到远程分支。
注意，在首次执行这条命令前（不一定非要到这个步骤，可以是这个步骤前的任意时刻，例如开始写代码前），建议请使用 git branch -u 远程名称/远程分支名称 把本地的分支与远程分支相关联，否则，在执行完 arc land 后，本地分支会被删除。如果你不想这样做，又需要保留现在分支的话，请使用 arc land --keep-branch。
默认情况下，Arcanist 会把变更推送到与本地分支相关联的远程分支，你可以使用 --remote 和 --onto 参数推送到你想要的远程和远程分支。
关于 arc land 的详细说明，可使用 arc land --help 查看。

0x70 与GitHub集成

这里的“集成”其实说的很心虚，因为在使用 Phabricator 后，GitHub 已经变为一个文件存储服务器了。常见的使用情况已经在 0x63 配置进行Code Review 的 配置代码仓库（Diffusion） 中写出，所以你懂的。

0x80 与Jenkins集成

在很久很久以前，我已经搭建了 Jenkins 用于自动构建，所以这次把 Phabricator 与 Jenkins 做了集成。
在 Jenkins 中安装完插件： Phabricator Differential Plugin 后，请参考这篇文章：

https://github.com/uber/phabricator-jenkins-plugin#phabricator-jenkins-plugin--
0x90 结束语

从配置完到开始试用已经有一个月左右，期间遇到了各种问题。经历了这些问题的洗礼，算是对 Phabricator 使用入门了吧。
这篇文章时隔多日才完成，尽管我想把每个细节写的详尽，但是记忆总是像被虫蛀过的木头，难免有些疏漏。加之文笔水平有限，各位就凑合着看吧。
在安装和使用过程中遇到的问题，欢迎各位评论交流。

0xA0 Q/A

1、如何强制用户 Code Review？
再次强调一下前提：在开始 Code Review 流程前，请先确认团队成员的git remote url已经切换为Phabricator上对应仓库地址。
强制用户Code Review，需要创建Herald Rule。在创建时，New Rule for 选：Commit Hook: Commit Content.；Rule Type选：Global，或者根据自己需要选择；Conditions中是组合条件，可以根据自己需要指定一些条件；Action 指定当前情况符合你指定的条件组合时执行的动作。示例如图：


示例中定义了以下规则：在向develop分支提交代码时，所提交的代码必须是通过 Code Review 流程审查通过的，否则会被拒绝；除非 Commit Message 中包含字符 @bypass-review 。
指定 Commit Message 中包含字符 @bypass-review 这种例外情况，主要考虑到在紧急修复一些问题，没时间等待审查时使用。
当提交被拒绝时，如图：
