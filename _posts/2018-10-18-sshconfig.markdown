---
title: ssh config
layout: post
category: web
author: 夏泽民
---
在 $HOME/.ssh/config 中加入以下内容：
Host *
ControlPersist yes
ControlMaster auto
ControlPath ~/.ssh/master-%r@%h:%p
这种方式第一次需要输入密码，然后一段时间内不需要输入密码了。
<!-- more -->
SSH 文件的结构及解释算法
本地系统的每个用户都可以维护一个客户端的 SSH 配置文件，这个配置文件可以包含你在命令行中使用 ssh 时参数，也可以存储公共连接选项并在连接时自动处理。你可以在命令上中使用 ssh 来指定 flag ，以覆盖配置文件中的选项。

SSH 客户端配置文件的位置
配置文件的文件名为 config ，位于用户 home 目录下的 .ssh 文件夹下。

~/.ssh/config
通常，该文件不是默认创建的，因此你可能要自己创建它。

配置文件的结构
配置文件通过 Host 来组织，每一个 Host 定义项为匹配的主机定义连接选项。通配符可以用，为了让选项有更大的范围。

配置文件看起来是这样的：

Host firsthost
    SSH_OPTIONS_1 custom_value
    SSH_OPTIONS_2 custom_value
    SSH_OPTIONS_3 custom_value

Host secondhost
    ANOTHER_OPTION custom_value

Host *host
    ANOTHER_OPTION custom_value

Host *
    CHANGE_DEFAULT custom_value
解释算法
只有理解 SSH 怎么解释配置文件，你才能写出合理的配置文件。

SSH 使命令行中给出的主机名与配置文件中定义的 Host 来匹配。它从文件顶部向下执行此操作，所以顺序非常重要。

现在是指出 Host 定义中的模式不必与您要连接的实际主机匹配的好时机。 实际上，您可以使用这些定义为主机设置别名，以替代实际的主机名。

看个例子：

Host dev1
    HostName dev1.example.com
    User tom
现在要连接到 tom@dev1.example.com，就可以通过在命令行中输入如下命令：

ssh dev1
记住这一点，我们现在继续讨论在由上而下的过程中，SSH 怎么应用每一个配置选中。它从顶部开始，检查每一个 Host 定义是否与命令行中给出的主机匹配。在上一个例子中，就是检查 dev1 。

当找到第一个匹配的主机定义时，每个关联的SSH选项都将应用于即将到来的连接（为了方便下边的讨论，这里我们称该连接接为“连接a”）。 尽管如此，解释并没有结束。

SSH 继续在文件中向下查找，检查是否有其他匹配的 Host 定义。如果有另一个 Host 定义匹配，SSH 将考虑该 Host 定义下的配置选项。如果新的配置选项中有 连接a 咱时没有使用的选项，就把这些选项也加入 连接a 中。
总结一下，SSH 将按顺序解释与命令行上给出主机名匹配的每个 Host 定义。在这个过程中，SSH 始终使用为每个选项给出的第一个值。没有办法覆盖之前已经匹配的 Host 定义给出的值。

HostName： 是目标主机的主机名，也就是平时我们使用ssh后面跟的地址名称。

Port：指定的端口号。

User：指定的登陆用户名。

IdentifyFile：指定的私钥地址。

当然不需要的时候 你也可以使用

ssh-add -D 删除所有管理的密钥

ssh-add -d 删除指定的

ssh-add -l 查看现在增加进去的指纹信息

ssh-add -L 查看现在增加进去的私钥

 

如果重启之后，会发现需要重新load一下ssh-agent

ssh-add -K 将指纹加到钥匙串里面去

ssh-add -A 可以把钥匙串里面的私钥密码，load进ssh-agent

SSH支持 ControlMaster 模式，可以复用之前已经建立的连接。所以开启这个功能之后，如果已经有一条到relay的链接，那么再连接的时候，就不需要再输入密码了。
而 ControlPersist 参数的含义就是在最后一个连接关闭之后也不真正的关掉连接，这样后面再连接的时候就还是不用输入密码。
启用这两个功能，就可以解决ssh登录时每次都需要重复输入密码的问题了。
在 $HOME/.ssh/config 中加入以下内容：（如果没有这个文件就touch一个，权限需要改成用户可访问才可以）
Host *
ControlPersist yes
ControlMaster auto
ControlPath ~/.ssh/master-%r@%h:%p

**Host**  
用于我们执行 SSH 命令的时候如何匹配到该配置。


* `*`，匹配所有主机名。
* `*.example.com`，匹配以 .example.com 结尾。
* `!*.dialup.example.com,*.example.com`，以 ! 开头是排除的意思。
* `192.168.0.?`，匹配 192.168.0.[0-9] 的 IP。


**AddKeysToAgent**  
是否自动将 key 加入到 `ssh-agent`，值可以为
 no(default)/confirm/ask/yes。


如果是 yes，key 和密码都将读取文件并以加入到 agent ，就像 `ssh-add`。其他分别是询问、确认、不加入的意思。添加到 ssh-agent 意味着将私钥和密码交给它管理，让它来进行身份认证。


**AddressFamily**  
指定连接的时候使用的地址族，值可以为 any(default)/inet(IPv4)/inet6(IPv6)。


**BindAddress**  
指定连接的时候使用的本地主机地址，只在系统有多个地址的时候有用。在 UsePrivilegedPort 值为 yes 的时候无效。


**ChallengeResponseAuthentication**  
是否响应支持的身份验证 chanllenge，yes(default)/no。


**Compression**  
是否压缩，值可以为 no(default)/yes。


**CompressionLevel**  
压缩等级，值可以为 1(fast)-9(slow)。6(default)，相当于 gzip。


**ConnectionAttempts**  
退出前尝试连接的次数，值必须为整数，1(default)。


**ConnectTimeout**  
连接 SSH 服务器超时时间，单位 s，默认系统 TCP 超时时间。


**ControlMaster**  
是否开启单一网络共享多个 session，值可以为 no(default)/yes/ask/auto。需要和 ControlPath 配合使用，当值为 yes 时，ssh 会监听该路径下的 control socket，多个 session 会去连接该 socket，它们会尽可能的复用该网络连接而不是重新建立新的。


**ControlPath**  
指定 control socket 的路径，值可以直接指定也可以用一下参数代替：


* %L 本地主机名的第一个组件
* %l 本地主机名（包括域名）
* %h 远程主机名（命令行输入）
* %n 远程原始主机名
* %p 远程主机端口
* %r 远程登录用户名
* %u 本地 ssh 正在使用的用户名
* %i 本地 ssh 正在使用 uid
* %C 值为 %l%h%p%r 的 hash


请最大限度的保持 ControlPath 的唯一。至少包含 %h，%p，%r（或者 %C）。


**ControlPersist**  
结合 ControlMaster 使用，指定连接打开后后台保持的时间。值可以为 no/yes/整数，单位 s。如果为 no，最初的客户端关闭就关闭。如果 yes/0，无限期的，直到杀死或通过其它机制，如：ssh -O exit。


**GatewayPorts**  
指定是否允许远程主机连接到本地转发端口，值可以为 no(default)/yes。默认情况，ssh 为本地回环地址绑定了端口转发器。


**HostName**  
真实的主机名，默认值为命令行输入的值（允许 IP）。你也可以使用 %h，它将自动替换，只要替换后的地址是完整的就 ok。


**IdentitiesOnly**  
指定 ssh 只能使用配置文件指定的 identity 和 certificate 文件或通过 ssh 命令行通过身份验证，即使 ssh-agent 或 PKCS11Provider 提供了多个 identities。值可以为 no(default)/yes。


**IdentityFile**  
指定读取的认证文件路径，允许 DSA，ECDSA，Ed25519 或 RSA。值可以直接指定也可以用一下参数代替：


* %d，本地用户目录 ~
* %u，本地用户
* %l，本地主机名
* %h，远程主机名
* %r，远程用户名


**LocalCommand**  
指定在连接成功后，本地主机执行的命令（单纯的本地命令）。可使用 %d，%h，%l，%n，%p，%r，%u，%C 替换部分参数。只在 PermitLocalCommand 开启的情况下有效。


**LocalForward**  
指定本地主机的端口通过 ssh 转发到指定远程主机。格式：LocalForward [bind_address:]post host:hostport，支持 IPv6。


**PasswordAuthentication**  
是否使用密码进行身份验证，yes(default)/no。


**PermitLocalCommand**  
是否允许指定 LocalCommand，值可以为 no(default)/yes。


**Port**  
指定连接远程主机的哪个端口，22(default)。


**ProxyCommand**  
指定连接的服务器需要执行的命令。%h，%p，%r


如：ProxyCommand /usr/bin/nc -X connect -x 192.0.2.0:8080 %h %p


**User**  
登录用户名




### 相关技巧


#### 管理多组密钥对
有时候你会针对多个服务器有不同的密钥对，每次通过指定 `-i` 参数也是非常的不方便。比如你使用 github 和 coding。那么你需要添加如下配置到 `~/.ssh/config`：
```
Host github
    HostName %h.com
    IdentityFile ~/.ssh/id_ecdsa_github
    User git
Host coding
    HostName git.coding.net
    IdentityFile ~/.ssh/id_rsa_coding
    User git
```
当你克隆 coding 上的某个仓库时：
```
# 原来
$ git clone git@git.coding.net:deepzz/test.git


# 现在
$ git clone coding:deepzz/test.git
```


#### vim 访问远程文件
vim 可以直接编辑远程服务器上的文件：
```
$ vim scp://example/docker-compose.yml
```


#### 远程服务当本地用
通过 LocalForward 将本地端口上的数据流量通过 ssh 转发到远程主机的指定端口。感觉你是使用的本地服务，其实你使用的远程服务。如远程服务器上运行着 Postgres，端口 5432（未暴露端口给外部）。那么，你可以：
```
Host db
    HostName db.example.com
    LocalForward 5433 localhost:5432
```
当你连接远程主机时，它会在本地打开一个 5433 端口，并将该端口的流量通过 ssh 转发到远程服务器上的 5432 端口。


首先，建立连接：
```
$ ssh db
```
之后，就可以通过 Postgres 客户端连接本地 5433 端口：
```
$ psql -h localhost -p 5433 orders
```


#### 多连接共享
什么是多连接共享？在你打开多个 shell 窗口时需要连接同一台服务器，如果你不想每次都输入用户名，密码，或是等待连接建立，那么你需要添加如下配置到 `~/.ssh/config`：
```
ControlMaster auto
ControlPath /tmp/%r@%h:%p
```


#### 禁用密码登录
如果你对服务器安全要求很高，那么禁用密码登录是必须的。因为使用密码登录服务器容易受到暴力破解的攻击，有一定的安全隐患。那么你需要编辑服务器的系统配置文件 `/etc/ssh/sshd_config`：
```
PasswordAuthentication no
ChallengeResponseAuthentication no
```


#### 关键词登录
为了更方便的登录服务器，我们也可以省略用户名和主机名，采用关键词登录。那么你需要添加如下配置到 `~/.ssh/config`：
```
Host deepzz                        # 别名
    HostName deepzz.com            # 主机地址
    User root                      # 用户名
    # IdentityFile ~/.ssh/id_ecdsa # 认证文件
    # Port 22                      # 指定端口
```
那么使用 `$ ssh deepzz` 就可以直接登录服务器了。


#### 代理登录
有的时候你可能没法直接登录到某台服务器，而需要使用一台中间服务器进行中转，如公司内网服务器。首先确保你已经为服务器配置了公钥访问，并开启了agent forwarding，那么你需要添加如下配置到 `~/.ssh/config`：
```
Host gateway
    HostName proxy.example.com
    User root
Host db
    HostName db.internal.example.com                  # 目标服务器地址
    User root                                         # 用户名
    # IdentityFile ~/.ssh/id_ecdsa                    # 认证文件
    ProxyCommand ssh gateway netcat -q 600 %h %p      # 代理命令
```
那么你现在可以使用 `$ ssh db` 连接了。

