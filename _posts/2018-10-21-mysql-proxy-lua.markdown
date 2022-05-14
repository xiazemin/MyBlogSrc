---
title: mysql-proxy-lua
layout: post
category: storage
author: 夏泽民
---
MySQL Proxy处于客户端应用程序和MySQL服务器之间，通过截断、改变并转发客户端和后端数据库之间的通信来实现其功能，这和WinGate之类的网络代理服务器的基本思想是一样的。代理服务器是和TCP/IP协议打交道，而要理解MySQL Proxy的工作机制，同样要清楚MySQL客户端和服务器之间的通信协议，MySQL Protocol包括认证和查询两个基本过程：
　　认证过程包括：
　　客户端向服务器发起连接请求
　　服务器向客户端发送握手信息
　　客户端向服务器发送认证请求
　　服务器向客户端发送认证结果
　　如果认证通过，则进入查询过程：
　　客户端向服务器发起查询请求
　　服务器向客户端返回查询结果
　　当然，这只是一个粗略的描述，每个过程中发送的包都是有固定格式的，想详细了解MySQL Protocol的同学，可以去这里看看。MySQL Proxy要做的，就是介入协议的各个过程。首先MySQL Proxy以服务器的身份接受客户端请求，根据配置对这些请求进行分析处理，然后以客户端的身份转发给相应的后端数据库服务器，再接受服务器的信息，返回给客户端。所以MySQL Proxy需要同时实现客户端和服务器的协议。
　　由于要对客户端发送过来的SQL语句进行分析，还需要包含一个SQL解析器。可以说MySQL Proxy相当于一个轻量级的MySQL了，实际上，MySQL Proxy的admin server是可以接受SQL来查询状态信息的。
　　MySQL Proxy通过lua脚本来控制连接转发的机制。主要的函数都是配合MySQL Protocol各个过程的，这一点从函数名上就能看出来：
　　connect_server()   
　　read_handshake()   
　　read_auth()   
　　read_auth_result()   
　　read_query()   
　　read_query_result()  
　　至于为什么采用lua脚本语言，我想这是因为MySQL Proxy中采用了wormhole存储引擎的关系吧，这个虫洞存储引擎很有意思，数据的存储格式就是一段lua脚本
<!-- more -->
通过这几个入口函数我们可以控制mysql-proxy的一些行为。

connect_server()          当代理服务器接受到客户端连接请求时(tcp中的握手)会调用该函数
read_handshake()        当mysql服务器返回握手相应时会被调用
read_auth()　　           当客户端发送认证信息(username,password,port,database)时会被调用
read_auth_result(aut)  当mysql返回认证结果时会被调用
read_query(packet)      当客户端提交一个sql语句时会被调用
read_query_result(inj)　当mysql返回查询结果时会被调用

配置文件

mysql-proxy.cnf(权限设为660)

 [mysql-proxy]

    admin-username=root

    admin-password=123456

    admin-lua-script=/usr/local/lib/admin.lua

    proxy-read-only-backend-addresses=192.168.2.115

    proxy-backend-addresses=192.168.2.117

    proxy-lua-script=/usr/local/lib/rw-splitting.lua

    log-file=/var/log/mysql-proxy.log

    log-level=debug

    daemon=true

keepalive=true


proxy-lua-script，指定一个Lua脚本来控制mysql-proxy的运行和设置，这个脚本在每次新建连接和脚本发生修改的的时候将重新调用

keepalive，额外建立一个进程专门监控mysql_proxy进程，当mysql_proxy crash予以重新启动；



启动

/usr/local/mysql-proxy/bin/mysql-proxy -P 192.168.2.112:3306 --defaults-file=/etc/mysql-proxy.cnf



读写分离

当proxy-lua-script指定为rw-splitting.lua时，mysql_proxy会对客户端传入的sql执行读写分离；

同一个事务，DML传输给backend，select则被传到read-only-backend；

Lua脚本默认最小4个最大8个以上的客户端连接才会实现读写分离（这是因为mysql-proxy会检测客户端连接, 当连接没有超过min_idle_connections预设值时，不会进行读写分离，即查询操作会发生到Master上），现改为最小1个最大2个，我们用vim修改/usr/local/lib/rw-splitting.lua脚本，改动内容如下所示：

    if not proxy.global.config.rwsplit then

            proxy.global.config.rwsplit = {

                    min_idle_connections = 1,

                    max_idle_connections = 2,

   

                    is_debug = false

            }

    end

 

read_query()函数内有这么一个判断

if stmt.token_name == "TK_SQL_SELECT" then 

这个语句的作用就是判断sql语句是不是以SELECT开始的，如果是查询的话，接下来会有这么个语句

local backend_ndx = lb.idle_ro() 

lb.idle_ro() 是通过 local lb = require("proxy.balance") 引入的balance.lua文件

这个函数的作用就是选择使用哪个读服务器，并返回服务器的index:max_conns_ndx

如何选择服务器呢？ 它通过循环遍历所有服务器，然后选出一个客户端连接（s.connected_clients）最少的服务器，这样在一定程度上实现负载均衡

    function idle_ro()  

        local max_conns = -1 

        local max_conns_ndx = 0 

     

        for i = 1, #proxy.global.backends do 

            local s = proxy.global.backends[i] 

            local conns = s.pool.users[proxy.connection.client.username] 

            -- pick a slave which has some idling connections 

            if s.type == proxy.BACKEND_TYPE_RO and s.state ~= proxy.BACKEND_STATE_DOWN and conns.cur_idle_connections > 0 then 

                if max_conns == -1 or s.connected_clients < max_conns then 

                    max_conns = s.connected_clients 

                    max_conns_ndx = i 

                end 

            end

        end 

     

        return max_conns_ndx 

    end

http://blog.csdn.net/clh604/article/details/8906022




failover

利用mysql_proxy实现failover

缺点：mysql_proxy单点故障；

原理：默认连接A，如果A宕掉则连接B，A启动后再连接到A；


编写failover脚本

vi $mysql-proxy_path/share/doc/mysql-proxy/mysql_failover.lua

function connect_server()

    for i = 1, #proxy.backends do

        local s = proxy.backends[i]

        print ("s.state:" + s.state)

        if s.state ~= proxy.BACKEND_STATE_DOWN then

            proxy.connection.backend_ndx = i

            print ("connecting to " .. i)

            return

        end

    end

end


function read_query(packet)

    for i = 1, #proxy.backends do

        local s = proxy.backends[i]

        print ("s.state:" + s.state)

        if s.state ~= proxy.BACKEND_STATE_DOWN then

            proxy.connection.backend_ndx = i

            print ("connecting to " .. i)

            return

        end

    end

end


启动mysql-proxy

$mysql-proxy_path/bin/mysql-proxy --proxy-address=:4040 --proxy-lua-script=$mysql-proxy_path/share/doc/mysql-proxy/mysql_failover.lua --proxy-backend-addresses=$A:3306 --proxy-backend-addresses=$B:3306 --log-level=error  --log-file=$mysql-proxy_path/mysql-proxy.log --keepalive --proxy-fix-bug-25371

此时客户端直接连接mysql-proxy即可；
