---
title: ssh 原理
layout: post
category: linux
author: 夏泽民
---
SSH（22端口）是Secure Shell Protocol的简写，由IETF网络工作小组（Network Working Group）制定；在进行数据传输之前，SSH先对联机数据包通过加密技术进行加密处理，加密后在进行数据传输。确保了传递的数据安全。

SSH是专为远程登录会话和其他网络服务提供的安全性协议。利用SSH协议可以有效的防止远程管理过程中的信息泄露问题，在当前的生产环境运维工作中，绝大多数企业普通采用SSH协议服务来代替传统的不安全的远程联机服务软件，如telnet（23端口，非加密的）等。

telnet（23端口）实现远程控制管理，但是不对数据进行加密

1.2 SSH服务主要提供的服务功能
在默认状态下，SSH服务主要提供两个服务功能：

（1）一个是提供类似telnet远程联机服务器的服务，即上面提到的SSH服务；

（2）另一个是类似FTP服务的sftp-server，借助SSH协议来传输数据的，提供更安全的SFTP服务（vsftp.proftp）

 特别提醒：

 SSH客户端（ssh命令）还包含一个很有用的远程安全拷贝命令scp，也是通过ssh协议工作的。
 
 
1）. 基于口令的认证
1.远程Server收到Client端用户TopGun的登录请求，Server把自己的公钥发给用户。
2.Client使用这个公钥，将密码进行加密。
3.Client将加密的密码发送给Server端。
4.远程Server用自己的私钥，解密登录密码，然后验证其合法性。
5.若验证结果，给Client相应的响应。
私钥是Server端独有，这就保证了Client的登录信息即使在网络传输过程中被窃据，也没有私钥进行解密，保证了数据的安全性，这充分利用了非对称加密的特性。
<!-- more -->
2）基于公钥认证

在上面介绍的登录流程中可以发现，每次登录都需要输入密码，很麻烦。SSH提供了另外一种可以免去输入密码过程的登录方式：公钥登录。
1.Client将自己的公钥存放在Server上，追加在文件authorized_keys中。
2.Server端接收到Client的连接请求后，会在authorized_keys中匹配到Client的公钥pubKey，并生成随机数R，用Client的公钥对该随机数进行加密得到pubKey(R)
，然后将加密后信息发送给Client。
3.Client端通过私钥进行解密得到随机数R，然后对随机数R和本次会话的SessionKey利用MD5生成摘要Digest1，发送给Server端。
4.Server端会也会对R和SessionKey利用同样摘要算法生成Digest2。
5.Server端会最后比较Digest1和Digest2是否相同，完成认证过程。


ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa
$ cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
$ chmod 0600 ~/.ssh/authorized_keys
ssh-keygen是用于生产密钥的工具。

-t：指定生成密钥类型（rsa、dsa、ecdsa等）
-P：指定passphrase，用于确保私钥的安全
-f：指定存放密钥的文件（公钥文件默认和私钥同目录下，不同的是，存放公钥的文件名需要加上后缀.pub）

~/.ssh中的四个文件：
1.id_rsa：保存私钥
2.id_rsa.pub：保存公钥
3.authorized_keys：保存已授权的客户端公钥
4.known_hosts：保存已认证的远程主机ID

需要注意的是：一台主机可能既是Client，也是Server。所以会同时拥有authorized_keys和known_hosts

1. known_hosts中存储的内容是什么？
known_hosts中存储是已认证的远程主机host key，每个SSH Server都有一个secret, unique ID, called a host key。

2. host key何时加入known_hosts的？
当我们第一次通过SSH登录远程主机的时候，Client端会有如下提示：

Host key not found from the list of known hosts.
Are you sure you want to continue connecting (yes/no)?
此时，如果我们选择yes，那么该host key就会被加入到Client的known_hosts中，格式如下：

# domain name+encryption algorithm+host key
example.hostname.com ssh-rsa AAAAB4NzaC1yc2EAAAABIwAAAQEA。。。
3. 为什么需要known_hosts？
最后探讨下为什么需要known_hosts，这个文件主要是通过Client和Server的双向认证，从而避免中间人（man-in-the-middle attack）攻击，每次Client向Server发起连接的时候，不仅仅Server要验证Client的合法性，Client同样也需要验证Server的身份，SSH client就是通过known_hosts中的host key来验证Server的身份的。

 测试ssh服务是否开启
方法1：telnet 10.0.0.41 22

说明：正确标准的方式退出telnet连接，按ctrl+]键，切换到telnet命令行，输入quit命令进行退出

方法2： ss -lntup |grep 22

[root@m01 ~]# ss -lntup |grep 22

tcp    LISTEN     0      128             :::22               :::*      users:(("sshd",1188,4))

tcp    LISTEN     0      128             *:22               *:*      users:(("sshd",1188,3))

方法3：netstat -lntup

方法4：nmap -p 22 10.0.0.41

方法5： nc 10.0.0.41 22

1.3.3 一个服务启动始终无法启动起来
（1）查看日志

（2）检查服务端口有没有冲突，如果冲突则重新指定一个端口  例：rsync --daemon --port=874

 telnet服务部署过程（默认telnet远程管理服务不支持root用户登录）
2.2.1 第一个里程碑：安装telnet服务软件
yum install -y telnet telnet-server

2.2.2 第二个里程碑：配置telnet服务可以被xinetd服务管理
vim /etc/xinetd.d/telnet

   修改    disable         = no

[root@m01 ~]# vim /etc/xinetd.d/telnet

# default: on
# description: The telnet server serves telnet sessions; it uses \

#       unencrypted username/password pairs for authentication.

service telnet

{

        flags           = REUSE

        socket_type     = stream

        wait            = no

        user            = root

        server          = /usr/sbin/in.telnetd

        log_on_failure  += USERID

        disable         = no                 --------表示telnet愿意被xinetd管理

}

2.2.3 第三个里程碑：启动xinetd服务
/etc/init.d/xinetd start

[root@m01 ~]# /etc/init.d/xinetd start

Starting xinetd:                                           [  OK  ]

[root@m01 ~]# netstat -lntup |grep 23

tcp        0      0 :::23                       :::*                    LISTEN      1708/xinetd 

这时候服务名不叫telnet了，叫xineted

[root@m01 ~]# telnet 10.0.0.41

Trying 10.0.0.41...

Connected to 10.0.0.61.

Escape character is '^]'.

CentOS release 6.9 (Final)

Kernel 2.6.32-696.el6.x86_64 on an x86_64

backup login: wuhuang

Password:  

利用wireshark软件进行抓包

1. 实现确认ssh远程协议进行数据的加密，telnet远程协议未对数据进行加密处理

2. 实现观察tcp三次握手过程和四次挥手过程

2.3 SSH远程连接服务与telnet服务区别
（1）telnet是不安全的远程连接，连接内容是明文的； ssh是加密的远程连接，连接内容是加密的。

（2）SSH服务默认支持root用户登录，telnet服务默认不支持root用户登录

2.4 tcpdump抓包命令
tcpdump -i eth0 -nn -c 5 "port 53"

-i             ------指定抓取哪一个网卡上产生数据包流量

--nn        ------抓取数据包中端口信息以数字方式显示

-c            -----表示抓取的数据包数量

"port 53"  ----- 双引号里面表示搜索数据包的条件

-w             -----将抓取数据保存到指定文件中

-r              -----对保存的数据包文件进行读取

[root@m01 ~]# tcpdump -i eth0 -nn -c 8 "port 22"

tcpdump: verbose output suppressed, use -v or -vv for full protocol decode

listening on eth0, link-type EN10MB (Ethernet), capture size 65535 bytes

10:32:52.885672 IP 10.0.0.61.22 > 10.0.0.1.51651: Flags [P.], seq 1975289928:1975290136, ack 3602557944, win 634, length 208

SSH服务由服务端软件OpenSSH（openssl）和客户端（常见的有SSH），SecureCRT，Putty，xshell组成，SSH服务默认使用22端口提供服务，它有两个不兼容的SSH协议版本，分别是1.x和2.x

3.1.1 openssh
[root@m01 ~]# rpm -qf `which ssh`

openssh-clients-5.3p1-122.el6.x86_64

[root@m01 ~]# rpm -ql openssh-clients

/etc/ssh/ssh_config            -----SSH客户端配置文件

/usr/bin/.ssh.hmac             

/usr/bin/scp                  -----远程复制命令

/usr/bin/sftp                  ----远程文件传输服务-

/usr/bin/slogin                ----远程登录命令

/usr/bin/ssh                  ----实时远程登录管理主机命令

/usr/bin/ssh-add

/usr/bin/ssh-agent

/usr/bin/ssh-copy-id           ----ssh服务分发公钥命令

/usr/bin/ssh-keyscan


openssh-clients总结（客户端）

/etc/ssh/ssh_config       --- ssh客户端配置文件

    /usr/bin/.ssh.hmac       --- ssh服务的算法文件

    /usr/bin/scp             --- 基于ssh协议，实现远程拷贝数据命令

    /usr/bin/sftp            --- 基于ssh协议，实现数据传输密命令

    /usr/bin/slogin          --- 远程登录主机连接命令    

    /usr/bin/ssh             --- 远程登录主机连接命令

    /usr/bin/ssh-add  --- 此参数必须和ssh-agent命令结合使用，将秘钥信息注册到ssh-agent代理服务中

    /usr/bin/ssh-agent       --- 启动ssh认证代理服务命令

    /usr/bin/ssh-copy-id     --- 远程分发公钥命令（ok）

/usr/bin/ssh-keyscan     --- 显示本地主机上默认的ssh公钥信息（ok）

[root@m01 ~]# rpm -ql openssh-server

/etc/pam.d/ssh-keycat

/etc/pam.d/sshd

/etc/rc.d/init.d/sshd                -----ssh服务启动脚本文件

/etc/ssh/sshd_config               ----ssh服务配置文件

/etc/sysconfig/sshd

/usr/libexec/openssh/sftp-server

/usr/libexec/openssh/ssh-keycat

/usr/sbin/.sshd.hmac

/usr/sbin/sshd                    -----ssh服务进程启动命令

[root@m01 ~]# ps -ef |grep sshd

root       1188      1  0 08:36 ?        00:00:00 /usr/sbin/sshd

root       1226   1188  0 08:38 ?        00:00:01 sshd: root@pts/0

root       2378   1228  0 18:26 pts/0    00:00:00 grep sshd

[root@m01 ~]# /etc/init.d/sshd stop

Stopping sshd:                                             [  OK  ]

[root@m01 ~]# ps -ef |grep sshd

root       1226      1  0 08:38 ?        00:00:01 sshd: root@pts/0

root       2396   1228  0 18:27 pts/0    00:00:00 grep sshd

[root@m01 ~]# /usr/sbin/sshd

[root@m01 ~]# ps -ef |grep sshd

root       1226      1  0 08:38 ?        00:00:01 sshd: root@pts/0

root       2398      1  0 18:28 ?        00:00:00 /usr/sbin/sshd

root       2400   1228  0 18:28 pts/0    00:00:00 grep sshd

 

  openssh-server总结（服务端）：

/etc/rc.d/init.d/sshd      --- ssh服务端启动脚本命令

   /etc/ssh/sshd_config      --- ssh服务端配置文件

   /usr/sbin/sshd           --- 启动ssh服务进程命令   

  扩展说明：sshd服务主要的两个进程

   /usr/sbin/sshd         --- 此进程对客户端第一次连接ssh服务端有影响

   sshd: root@pts/0       --- 一旦ssh连接成功，是否可以始终保持连接有次进程决定

代理服务中

/usr/bin/ssh-agent            --- 启动ssh认证代理服务命令

/usr/bin/ssh-copy-id          --- 远程分发公钥命令

/usr/bin/ssh-keyscan          --- 显示本地主机上默认的ssh公钥信息

3.1.2 openssl：（https）
u OpenSSH同时支持SSH1.x和2.x。用SSH 2.x的客户端程序不能连接到SSH 1.x的服务程序上。

u SSH服务端是一个守护进程（daemon），它在后台运行并响应来自客户端的连接请求。SSH服务端的进程名为sshd，负责实时监听远程SSH客户端的远程连接请求，并进行处理，一般包括公共密钥认证，密钥交换，对称密钥加密和非安全连接等。这个SSH服务就是我们前面基础系统优化中保留开机自启动的服务之一。

u ssh客户端包含ssh以及像scp（远程拷贝），slogin（远程登录），sftp（安全FTP文件传输）等应用程序。

u ssh的工作机制大致是本地的ssh客户端先发送一个连接请求到远程的ssh服务端，服务端检查连接的客户端发送的数据包和IP地址，如果确认合法，就会发送密钥发回给服务端，自此连接建立

3.2 ssh协议实现加密技术原理
3.2.1  公钥和私钥
（1）利用了公钥和私钥实现数据加密和解密（公钥：锁头  私钥：钥匙）

（2）利用公钥和私钥实现了认证机制，公钥可以在网络中传输，私钥在本地主机保存

ssh_host_dsa_key

ssh_host_dsa_key.pub   

# ssh-keyscan -t dsa 172.16.1.61   根据IP查看对应主机的公钥

# 172.16.1.61 SSH-2.0-OpenSSH_5.3

172.16.1.61 ssh-dss AAAAB3NzaC1kc3MAAACBAKSV66UzxqEzt8TKEFcyQtYPMC3y7YeZh7YVsy+E4KaMQAEVzOwcp2b6IXFyMDGNrystP9jfV7cXKC+2S7LkayJnOr8l3NgmzY

3.2.2 sshv1与sshv2版本比较
v1版本钥匙和锁头默认不会变化，数据传输不安全

v2版本钥匙和锁头会经常变化，数据传输更安全

3.2.3 ssh加密算法v2.x
u 在SSH 1.x的联机过程中，当Server接受Client端的Private Key后，就不再针对该次联机的Key pair进行检验。此时若有恶意黑客针对该联机的Key pair对插入恶意的程序代码时，由于服务端你不会再检验联机的正确性，因此可能会接收该程序代码，从而造成系统被黑掉的问题。

u 为了改正这个缺点，SSH version 2 多加了一个确认联机正确性的Diffie-Hellman机制，在每次数据传输中，Server都会以该机制检查数据的来源是否正确，这样，可以避免联机过程中被插入恶意程序代码的问题。也就是说，SSH version 2 是比较安全的。

u 由于SSH1协议本身存在较大安全问题，因此，建议大家尽量都用SSH2的联机模式。而联机版本的设置则需要在SSH主机端与客户端均设置好才行。

3.3 SSH服务认证类型
3.3.1 基于密码的认证
基于口令的安全验证的方式就是大家现在一直在用的，只要知道服务器的SSH连接账号和口令（当然也要知道对应服务器的IP及开放的SSH端口，默认为22），就可以通过ssh客户端登录到这台远程主机。此时，联机过程中所有传输的数据都是加密的。

[root@m01 ~]# ssh -p22 root@10.0.0.7

root@10.0.0.7's password:

Last login: Mon Jan 29 08:38:11 2018 from 10.0.0.1

3.3.2 基于秘钥的认证（实现免密码管理）
u 基于密钥的安全验证方式是指，需要依靠密钥，也就是必须事先建立一对密钥对，然后把公用密钥（Public key）放在需要访问的目标服务器上，另外，还需要把私有密钥（Private key）放到SSH的客户端或对应的客户端服务器上。

   私钥不能在网络中传输-----私钥可以解密公钥

   公钥可以在网络中传输-----公钥不能解密私钥

u 此时，如果要想连接到这个带有公用密钥的SSH服务器，客户端SSH软件或客户端服务器就会向SSH服务器发出请求，请求用联机的用户密钥进行安全验证。SSH服务器收到请求之后，会先在该SSH服务器上连接的用户的家目录下寻找事先放上去的对应用户的公用密钥，然后把它和连接的SSH客户端发送过来的公用密钥进行比较。如果两个密钥一致，SSH服务器就用公用密钥加密“质询”并把它发送给SSH客户端。

u SSH客户端收到“质询”之后就可以用自己的私钥解密，再把它发送给SSH服务器。使用这种方式，需要知道联机用户的密钥文件。与第一种基于口令验证的方式相比，第二种方式不需要在网络上传送口令密码，所以安全性更高了，这时我们也要注意保护我们的密钥文件，特别是私钥文件，一旦被黑客获取，危险就很大了。

u 基于密钥的安全认证也有windows客户端和linux客户端的区别。在这里我们主要介绍的是linux客户端和linux服务端之间的密钥认证。


总结：

①. 在管理端创建出秘钥对（创建两个信物）

②. 管理端将公钥（锁头）传输给被管理端，锁头传输给被管理端要基于密码方式认证

③. 管理端向被管理端发出建立连接请求

④. 被管理端发出公钥质询

⑤. 管理端利用私钥解密公钥，进行公钥质询响应

⑥. 被管理端接收到质询响应，确认基于秘钥认证成功

3.4 基于秘钥认证配置部署过程
3.4.1 第一个里程碑：管理服务器上创建秘钥对
ssh-keygen -t rsa

[root@m01 ~]# ssh-keygen -t rsa

Generating public/private rsa key pair.

Enter file in which to save the key (/root/.ssh/id_rsa):               ---确认私钥保存路径（默认路径）

Created directory '/root/.ssh'.

Enter passphrase (empty for no passphrase):                  ---确认私钥文件是否设置密码（设置为空） 

Enter same passphrase again:                              ---确认私钥文件是否设置密码（设置为空）

Your identification has been saved in /root/.ssh/id_rsa.          ---私钥保存位置

Your public key has been saved in /root/.ssh/id_rsa.pub.         ---公钥保存位置

The key fingerprint is:

37:a5:23:85:89:5f:62:93:40:d6:9c:3d:04:d0:07:a8 root@m01

The key's randomart image is:

+--[ RSA 2048]----+

|     .=*o*.      |

|     ..o=++      |

|     .. B.o..    |

|    E  o = o     |

|        S =      |

|         o o     |

|                 |

|                 |

|                 |

+-----------------+

命令说明：

1）创建密钥对时，要你输入的密码，为进行密钥对验证时输入的密码（和linux角色登录的密码完全没有关系）；

2）如果我们要进行的是SSH免密码连接，那么这里密码为空跳过即可。

3）如果在这里你输入了密码，那么进行SSH密钥对匹配连接的时候，就需要输入这个密码了。（此密码为独立密码）

4）用户家目录下的.ssh隐藏目录下会生成：id_rsa id_rsa.pub 两个文件。id_rsa是用户的私钥；id_rsa.pub则是公钥

 创建完后，会在当前用户的宿主目录.ssh/下生成一个公钥与私钥

[root@m01 ~]#  ls .ssh/

id_rsa  id_rsa.pub

 

3.4.2 第二个里程碑：分发公钥给被管理端主机
ssh-copy-id -i /root/.ssh/id_rsa.pub root@172.16.1.31

[root@m01 ~]# ssh-copy-id -i /root/.ssh/id_rsa.pub root@172.16.1.31

The authenticity of host '172.16.1.31 (172.16.1.31)' can't be established.

RSA key fingerprint is 57:3f:64:68:95:4d:99:54:01:33:ab:47:a0:72:da:bf.

Are you sure you want to continue connecting (yes/no)? yes

Warning: Permanently added '172.16.1.31' (RSA) to the list of known hosts.

root@172.16.1.31's password: 输入密码

Now try logging into the machine, with "ssh 'root@172.16.1.31'", and check in:

 

  .ssh/authorized_keys

 

to make sure we haven't added extra keys that you weren't expecting.

命令说明：

-i   ---指定要进行分发的公钥文件

root表示客户端以什么样的身份登录到服务端

如果ssh端口不是默认的22端口，需要在公钥文件后面加参数：-p 3600

3.4.3 第三个里程碑：利于基于秘钥方式登录测试
[root@m01 ~]# ssh 172.16.1.31      可以直接免密登录

Last login: Mon Jan 29 08:37:44 2018 from 10.0.0.1

[root@nfs01 ~]# exit

3.5 实现多台主机之间基于秘钥的彼此相互访问
3.5.1 第一个里程碑：启动认证代理服务
[root@m01 ~]# eval `ssh-agent -s`

Agent pid 2602

说明：

eval这个命令，相当于执行俩次bash

3.5.2 第二个里程碑：向agent代理服务器注册本地服务器私钥信息
[root@m01 ~]# ssh-add

Identity added: /root/.ssh/id_rsa (/root/.ssh/id_rsa)

3.5.3 第三个里程碑：将注册凭证信息通过远程登录方式给被管理主机
[root@m01 ~]# ssh -A 172.16.1.31

Last login: Wed Jan 31 17:51:07 2018 from 172.16.1.61

[root@nfs01 ~]#

注意：

1.注册凭证不是私钥信息

2.被管理主机接受到后，会产生凭证信息/tmp/ssh-xxx/agent.12334

image.png 
3.6 ssh服务配置文件说明
修改SSH服务的运行参数，是通过修改配置文件/etc/ssh/sshd.config文件来实现的。
一般来说SSH服务使用默认的配置已经能够很好的工作了，如果对安全要求不高，仅仅提供SSH服务的情况，可以不需要修改任何配置。

[root@m01 ~]# vim /etc/ssh/sshd_config

#       $OpenBSD: sshd_config,v 1.80 2008/07/02 02:24:18 djm Exp $

 

13 #Port 22                                     --- 表示修改默认端口号信息，ssh连接默认端口22 

15 #ListenAddress 0.0.0.0                                      --- 指定监听本地主机网卡地址信息 

42 #PermitRootLogin yes            --- 是否允许root用户远程登录，默认允许，建议禁止root远程登录

65 #PermitEmptyPasswords no                                  --- 是否允许空密码

81 GSSAPIAuthentication yes      --- 默认此参数配置信息为yes，总是要对连接进行以下GSSAPI认证

122 UseDNS yes    --- 默认此参数配置信息为yes，要对访问过来主机信息做dns反向解析（建议设为no——服务端）

说明：

1）#号代表注释，#号后面有空格的表示对数据内容的描述信息，#号后面没空格的表示参数信息，注释的 参数信息是默认的参数配置

2）一旦修改了Port，那么ssh登录时就需要-p指定端口号，不然会登录失败，ssh默认登录22端口

3）一旦修改了ListenAddress，监听地址，那么不再地址范围内的所有客户端将无法远程连接服务器。

4） GSSAPIAuthentication yes 和 UseDNS yes   设为no可以提高连接速度

5）修改配置文件后，需要重启sshd服务才能生效

image.png 

第4章 SSH服务安全配置
4.1  SSH入侵案例说明
SSH入侵网友案例：http://phenixikki.blog.51cto.com/7572938/1546669

4.1.1 如何防止SSH登录入侵小结：
    1、用密钥登录，不用密码登陆。

    2、牤牛阵法：解决SSH安全问题

       a.防火墙封闭SSH,指定源IP限制(局域网、信任公网)

       b.开启SSH只监听本地内网IP（ListenAddress 172.16.1.61）。

    3、尽量不给服务器外网IP

    4、最小化（软件安装-授权）

    5、给系统的重要文件或命令做一个指纹

    6. 给他锁上 chattr +i   +a

    
 