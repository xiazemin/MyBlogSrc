---
title: gracehttp
layout: post
category: golang
author: 夏泽民
---
经典平滑升级方案

服务器开发运维中，平滑升级是一个老生常谈的话题。拿一个http server来说，最常见的方案就是在http server前面加挂一个lvs负载，通过健康检查接口决定负载的导入与摘除。具体来说就是http server 提供一个/status 接口，服务器返回一个status文件，内容为ok，lvs负载定时访问这个接口，判断服务健康状况决定导入流量和切断流量。一般都会定一些策略，比如：访问间隔5秒，健康阈值2，异常阈值2之类的。意思就是每隔5秒访问一次/status接口，2次成功后，确认服务正常，开始导入流量，2次失败确认服务异常切断流量。当服务升级时，修改status文件内容为off，等待lvs健康检查确认服务为异常状态时主动切断流量，此时进行服务器的升级操作，服务重启完毕后，将status文件内容修改回ok，等待lvs健康检查确认服务正常后导入流量，以此步骤逐步完成剩余的机器的发布操作。将以上步骤完善成脚本，拆分为pre（预升级，ok修改为off）、post（发布代码，重启服务）、check（服务检查）、online（上线，off修改为ok）几个动作，与代码发布平台结合基本就实现了一般服务的自动化发版管理。360内部的代码发布平台Furion就是基于此原理工作的。

 
经典平滑升级方案的问题

一般的web服务使用上述平滑升级方案，基本上已经够用了。那这个方案还有什么问题吗？吹毛求疵的讲，还是有的。

发布过程中，正在发布的机器被摘除，其他机器承压增大。

发布过程仍然花费一些时间，按照上述策略指定的参数，发布一次至少需要20秒，当然我们可以调整参数，但是要面临浪费资源或者网络抖动误判导致切断流量的问题。

切断流量瞬间会导致未完成请求返回不完整。

这些问题一般来说都不算大问题，服务器资源做好冗余就够了，但是当服务器数量很大，服务器QPS很高的情况，小问题也会变大问题。所有寻求完美无缝重启的方案就是解决问题的关键了。
https://mp.weixin.qq.com/s?__biz=MzU4ODgyMDI0Mg==&mid=2247487071&idx=1&sn=c0098f0ea50f6b1fc5c94ea9e68e8bfb
<!-- more -->
优雅重启

golang语言http服务的优雅重启开源库也有一些，我们选择Facebook开源的库进行研究。代码地址https://github.com/facebookarchive/grace.git。网上的开源库的实现或简单或复杂，其实原理都差不多，执行优雅重启的过程基本如下：

发布新的bin文件去覆盖老的bin文件 

发送一个信号量，告诉正在运行的进程，进行重启 

正在运行的进程收到信号后，会以子进程的方式启动新的bin文件

新进程接受新请求，并处理 

老进程不再接受请求，但是要等正在处理的请求处理完成，所有在处理的请求处理完之后，便自动退出
其实我总结了一下，就两个关键点，一个是子进程继承端口监听启动，接受新请求处理；另一个是父进程优雅关闭。通过以上两个步骤基本上就实现了服务的无缝重启，发布过程中流量无损，发布消耗时间理论上最大也就是一个请求的超时时间，回滚服务也很简单，将旧版本服务重发一次就好了。

 
源码分析

1

使用方法

示例使用了流行的http库 gin，我们一般用法如下

func main
()
 
{

  engine 
:=
 gin
.
New
()

  engine
.
Use
(
httpserver
.
NewAccessLogger
(),
 gin
.
Recovery
())

  controller
.
Regist
(
engine
)

  srv 
:=
 
&
http
.
Server
{

        
Addr
:
         
":80"
,

        
Handler
:
      engine
,

        
ReadTimeout
:
  
30
 
*
 time
.
Second
,

        
WriteTimeout
:
 
30
 
*
 time
.
Second
,

    
}

  monitor
.
Init
()

  srvMonitor 
:=
 
&
http
.
Server
{

        
Addr
:
         
":9900"
,

        
Handler
:
      
nil
,

        
ReadTimeout
:
  
30
 
*
 time
.
Second
,

        
WriteTimeout
:
 
30
 
*
 time
.
Second
,

    
}

  grace
.
Serve
(
srv
,
 srvMonitor
)

}

grace.Serve函数参数是一个切片，可以处理多个server的端口监听继承与优雅关闭。此外还提供了关闭前的hook，使用方法如下：

gracehttp
.
ServeWithOptions
([]*
http
.
Server
{
srv
,
 srvMonitor
},
 gracehttp
.
PreStartProcess
(
func
()
 error 
{

        logger
.
Info
(
"do PreStartProcess\n"
)

        
return
 
nil

    
}))

在调研中我发现项目上有错误的使用方法，如下：

func startHttp
()
 
{

    engine 
:=
 gin
.
New
()

    engine
.
Use
(
httpserver
.
NewAccessLogger
(),
 gin
.
Recovery
())

    controller
.
Regist
(
engine
)

    srv 
:=
 
&
http
.
Server
{

        
Addr
:
         
":80"
,

        
Handler
:
      engine
,

        
ReadTimeout
:
  
30
 
*
 time
.
Second
,

        
WriteTimeout
:
 
30
 
*
 time
.
Second
,

    
}

    monitor
.
Init
()

    srvMonitor 
:=
 
&
http
.
Server
{

        
Addr
:
         
":9900"
,

        
Handler
:
      
nil
,

        
ReadTimeout
:
  
30
 
*
 time
.
Second
,

        
WriteTimeout
:
 
30
 
*
 time
.
Second
,

    
}

    grace
.
Serve
(
srv
,
 srvMonitor
)

}



func main
()
 
{

    go startHttp
()

    
//注册信号

    go signalHandler
()

    
<-
quiet

    logger
.
Info
(
"Close Server"
)

}



func signalHandler
()
 
{

    c 
:=
 make
(
chan os
.
Signal
)

    signal
.
Notify
(
c
,
 syscall
.
SIGHUP
,
 syscall
.
SIGINT
,
 syscall
.
SIGTERM
,
 syscall
.
SIGKILL
,
 syscall
.
SIGQUIT
)

    s 
:=
 
<-
c

    logger
.
Info
(
"get siginal  siginal=%v"
,
 s
)

    quiet 
<-
 
1

}

这里为什么出错了呢，是因为他将grace.Serve(srv,srvMonitor) 放在goroutine里面了，并且自己又监听了一遍信号，这样会导致旧进程优雅关闭前，父进程已经已经退出了，优雅关闭就失效了。

2

关键代码

我们按照程序启动的顺序逻辑来讲，大体如下：

执行启动端口监听，挂载server，判断当前进程如果是子进程就向父进程发送SIGTERM信号。

goroutine 执行wg.Add 和wg.Wait() ，等待所有挂载的server停止工作后执行退出进程。

goroutine 执行 signalHandler，等待SIGTERM和SIGUSR2信号。收到SIGTERM信号执行每个server的优雅关闭，关闭完后执行wg.Done()，wg全部Done之后在2中执行了退出进程操作；收到SIGUSR2信号时，执行启动子进程操作。

子进程启动执行1，会向父进程发送SIGTERM信号，父进程收到SIGTERM信号执行3，进行优雅关闭操作。

总结起来就是执行启动重启时，执行shell命令：

 pgrep 
(你的项目名)
 
|
xargs kill 
-
SIGUSR2

#(注意：要使用bash)。

你的项目会启动子进程，并继承父进程监听的端口，启动成功后再向父进程发送SIGTERM信号， 旧进程执行优雅关闭。我们看关键的struct

// gracehttp/http.go

type app 
struct
 
{

    servers         
[]*
http
.
Server

    http            
*
httpdown
.
HTTP

    net             
*
gracenet
.
Net

    listeners       
[]
net
.
Listener

    sds             
[]
httpdown
.
Server

    preStartProcess func
()
 error

    errors          chan error

}

// httpdown/httpdown.go

type HTTP 
struct
 
{

    
// StopTimeout is the duration before we begin force closing connections.

    
// Defaults to 1 minute.

    
StopTimeout
 time
.
Duration



    
// KillTimeout is the duration before which we completely give up and abort

    
// even though we still have connected clients. This is useful when a large

    
// number of client connections exist and closing them can take a long time.

    
// Note, this is in addition to the StopTimeout. Defaults to 1 minute.

    
KillTimeout
 time
.
Duration



    
// Stats is optional. If provided, it will be used to record various metrics.

    
Stats
 stats
.
Client



    
// Clock allows for testing timing related functionality. Do not specify this

    
// in production code.

    
Clock
 clock
.
Clock

}



// gracenet/net.go

type 
Net
 
struct
 
{

    inherited   
[]
net
.
Listener

    active      
[]
net
.
Listener

    mutex       sync
.
Mutex

    inheritOnce sync
.
Once

    
// used in tests to override the default behavior of starting from fd 3.

    fdStart 
int

}

我们知道函数调用是从grace.Serve(srv, srvMonitor)开始的,Serve函数会new一个app，一路执行下去关键函数如下：a.run()、a.listen()、a.serve()、 a.wait()、a.signalHandler()、 a.term()、a.net.StartProcess()。

a.run() 大体逻辑如下：

var
 
(

    didInherit 
=
 os
.
Getenv
(
"LISTEN_FDS"
)
 
!=
 
""

    ppid       
=
 os
.
Getppid
()

)



func 
(
a 
*
app
)
 run
()
 error 
{

    a
.
listen
()

    a
.
serve
()

    
if
 didInherit 
&&
 ppid 
!=
 
1
 
{

    syscall
.
Kill
(
ppid
,
 syscall
.
SIGTERM
)

    
}

    waitdone 
:=
 make
(
chan 
struct
{})

    go func
()
 
{

    defer close
(
waitdone
)

    a
.
wait
()

    
}()

    
select
 
{

    
case
 err 
:=
 
<-
a
.
errors
:

        
...

    
case
 
<-
waitdone
:

        logger
.
Printf
(
"Exiting pid %d."
,
 os
.
Getpid
())

        
return
 
nil

    
}

}

启动监听、挂载server，通过环境变量LISTEN_FDS判断当前进程是否为子进程，如果是就发送信号杀父进程。goroutine中执行wait()函数等待优雅关闭或者平滑启动子进程。

a.listen() 关键逻辑如下：

func 
(
a 
*
app
)
 listen
()
 error 
{

    
for
 _
,
 s 
:=
 range a
.
servers 
{

        l
,
 err 
:=
 a
.
net
.
Listen
(
"tcp"
,
 s
.
Addr
)

                
......

        a
.
listeners 
=
 append
(
a
.
listeners
,
 l
)

    
}

    
return
 
nil

}

这里看出app struct 中listeners用来存储监听的net.Listener的数组 ，net就是Net，封装了net.ListenTCP等逻辑（这里我只关注了TCP逻辑），inherited 和 active 两个数组分别用来存储继承自父进程的net.Listener 和 启动的net.Listener，这块父进程启动，即首次启动时逻辑很简单，略过，子进程启动，即非首次启动在介绍a.net.StartProccess时细讲。

a.serve() 关键逻辑如下：

func 
(
a 
*
app
)
 serve
()
 
{

    
for
 i
,
 s 
:=
 range a
.
servers 
{

        a
.
sds 
=
 append
(
a
.
sds
,
 a
.
http
.
Serve
(
s
,
 a
.
listeners
[
i
]))

    
}

}

这里涉及了app struct里面的两个字段，http和sds。http即 HTTP struct， 这里面封装了http server优雅关闭相关的逻辑，具体的细节很繁琐，我用一个简单的模型来说明一下吧。a.http.Serve(srv,l) 函数封装执行了srv.Serve(l)，即挂载srv， 并返回了一个httpdown.server的实例， 这个实例实现了httpdown.Server 接口，如下：

// httpdown/httpdown.go

type 
Server
 
interface
 
{

    
// Wait waits for the serving loop to finish. This will happen when Stop is

    
// called, at which point it returns no error, or if there is an error in the

    
// serving loop. You must call Wait after calling Serve or ListenAndServe.

    
Wait
()
 error



    
// Stop stops the listener. It will block until all connections have been

    
// closed.

    
Stop
()
 error

}

精简后实现的模型如下：

func 
(
s 
*
server
)
 serve
()
 
{

        
// 即前面提到的 srv.Serve(l)，被封装的挂载srv的代码

    s
.
serveErr 
<-
 s
.
server
.
Serve
(
s
.
listener
)

    close
(
s
.
serveDone
)

    close
(
s
.
serveErr
)

}



func 
(
s 
*
server
)
 
Wait
()
 error 
{

    
if
 err 
:=
 
<-
s
.
serveErr
;
 
!
isUseOfClosedError
(
err
)
 
{

        
return
 err

    
}

    
return
 
nil

}



func 
(
s 
*
server
)
 
Stop
()
 error 
{

    s
.
stopOnce
.
Do
(
func
()
 
{

        closeErr 
:=
 s
.
listener
.
Close
()

        
<-
s
.
serveDone

                
......

                
// 等待连接关闭或者超时后强杀连接等复杂逻辑

                
......

        
if
 closeErr 
!=
 
nil
 
&&
 
!
isUseOfClosedError
(
closeErr
)
 
{

            s
.
stopErr 
=
 closeErr

        
}

    
})

    
return
 s
.
stopErr

}

s.serveErr <- s.server.Serve(s.listener) 启动成功后会在这里挂住，失败直接返回错误，Wait() 函数提供给a.wait()调用，正常情况也是挂住，等Stop() 里面 closeErr := s.listener.Close() 执行后返回。这块的逻辑要结合 a.wait()、 a.signalHandler()、 a.term() 一起来分析

a.wait() 和 a.term() 的代码

func 
(
a 
*
app
)
 wait
()
 
{

    
var
 wg sync
.
WaitGroup

    wg
.
Add
(
len
(
a
.
sds
)
 
*
 
2
)
 
// Wait & Stop

    go a
.
signalHandler
(&
wg
)

    
for
 _
,
 s 
:=
 range a
.
sds 
{

        go func
(
s httpdown
.
Server
)
 
{

            defer wg
.
Done
()

            
if
 err 
:=
 s
.
Wait
();
 err 
!=
 
nil
 
{

                a
.
errors 
<-
 err

            
}

        
}(
s
)

    
}

    wg
.
Wait
()

}



func 
(
a 
*
app
)
 term
(
wg 
*
sync
.
WaitGroup
)
 
{

    
for
 _
,
 s 
:=
 range a
.
sds 
{

        go func
(
s httpdown
.
Server
)
 
{

            defer wg
.
Done
()

            
if
 err 
:=
 s
.
Stop
();
 err 
!=
 
nil
 
{

                a
.
errors 
<-
 err

            
}

        
}(
s
)

    
}

}

a.run() 函数里面会goroutine 执行 a.wait()，它会goroutine执行信号处理 a.signalHandler() 函数，创建一个WaitGroup 等待所有的httpdown.server执行s.Wait()函数返回。a.signalHandler() 函数基本上逻辑就是监听signal.Notify信号，收到SIGTERM信号执行a.term() ，收到SIGUSR2信号执行a.net.StartProcess()。a.term() 函数就是遍历执行所有httpdown.server的s.Stop()，进行优雅关闭，结合上面的代码来看，每一个s.Stop() 会导致s.Wait() 返回，即执行了两次wg.Done()， 所有httpdown.server 优雅关闭后导致a.wait()返回，进而waitdone关闭， 进程最后退出。下面是a.signalHandler()函数的代码

func 
(
a 
*
app
)
 signalHandler
(
wg 
*
sync
.
WaitGroup
)
 
{

    ch 
:=
 make
(
chan os
.
Signal
,
 
10
)

    signal
.
Notify
(
ch
,
 syscall
.
SIGINT
,
 syscall
.
SIGTERM
,
 syscall
.
SIGUSR2
)

    
for
 
{

        sig 
:=
 
<-
ch

        
switch
 sig 
{

        
case
 syscall
.
SIGINT
,
 syscall
.
SIGTERM
:

            
// this ensures a subsequent INT/TERM will trigger standard go behaviour of

            
// terminating.

            signal
.
Stop
(
ch
)

            a
.
term
(
wg
)

            
return

        
case
 syscall
.
SIGUSR2
:

            err 
:=
 a
.
preStartProcess
()

            
if
 err 
!=
 
nil
 
{

                a
.
errors 
<-
 err

            
}

            
// we only return here if there's an error, otherwise the new process

            
// will send us a TERM when it's ready to trigger the actual shutdown.

            
if
 _
,
 err 
:=
 a
.
net
.
StartProcess
();
 err 
!=
 
nil
 
{

                a
.
errors 
<-
 err

            
}

        
}

    
}

}

a.net.StartProcess() 函数是启动子进程的逻辑，这里需要详细介绍一下

const
 
(

    
// Used to indicate a graceful restart in the new process.

    envCountKey       
=
 
"LISTEN_FDS"

    envCountKeyPrefix 
=
 envCountKey 
+
 
"="

)



type filer 
interface
 
{

    
File
()
 
(*
os
.
File
,
 error
)

}



func 
(
n 
*
Net
)
 
StartProcess
()
 
(
int
,
 error
)
 
{

    listeners
,
 err 
:=
 n
.
activeListeners
()

    
if
 err 
!=
 
nil
 
{

        
return
 
0
,
 err

    
}

    
// Extract the fds from the listeners.

    files 
:=
 make
([]*
os
.
File
,
 len
(
listeners
))

    
for
 i
,
 l 
:=
 range listeners 
{

        files
[
i
],
 err 
=
 l
.(
filer
).
File
()

        
if
 err 
!=
 
nil
 
{

            
return
 
0
,
 err

        
}

        defer files
[
i
].
Close
()

    
}

    
// Use the original binary location. This works with symlinks such that if

    
// the file it points to has been changed we will use the updated symlink.

    argv0
,
 err 
:=
 
exec
.
LookPath
(
os
.
Args
[
0
])

    
if
 err 
!=
 
nil
 
{

        
return
 
0
,
 err

    
}



    
// Pass on the environment and replace the old count key with the new one.

    
var
 env 
[]
string

    
for
 _
,
 v 
:=
 range os
.
Environ
()
 
{

        
if
 
!
strings
.
HasPrefix
(
v
,
 envCountKeyPrefix
)
 
{

            env 
=
 append
(
env
,
 v
)

        
}

    
}

    env 
=
 append
(
env
,
 fmt
.
Sprintf
(
"%s%d"
,
 envCountKeyPrefix
,
 len
(
listeners
)))



    allFiles 
:=
 append
([]*
os
.
File
{
os
.
Stdin
,
 os
.
Stdout
,
 os
.
Stderr
},
 files
...)

    process
,
 err 
:=
 os
.
StartProcess
(
argv0
,
 os
.
Args
,
 
&
os
.
ProcAttr
{

        
Dir
:
   originalWD
,

        
Env
:
   env
,

        
Files
:
 allFiles
,

    
})

    
if
 err 
!=
 
nil
 
{

        
return
 
0
,
 err

    
}

    
return
 process
.
Pid
,
 
nil

}

n.activeListeners()返回 n.active中的net.Listener 数组的副本，files是从中提取出的fd列表。注意allFiles在files前面拼接了3个标准输入输出，记住这个数字。env 中修改了环境变量LISTEN_FDS等于listener的数量。这里的启动子进程的方法是os.StartProcess()，我看了其他的开源库都用syscall.ForkExec

fork
,
 err 
:=
 syscall
.
ForkExec
(
os
.
Args
[
0
],
 os
.
Args
,
 
&
os
.
ProcAttr
{

        
Dir
:
   originalWD
,

        
Env
:
   env
,

        
Files
:
 allFiles
,

    
})

两种的区别后续还有待研究。还记得前面没有展开的Net中的inherited 和 active么，这里我们细讲一下。

func 
(
n 
*
Net
)
 
Listen
(
nett
,
 laddr 
string
)
 
(
net
.
Listener
,
 error
)
 
{

    
......

       
// 仅关注tcp逻辑

       
return
 n
.
ListenTCP
(
nett
,
 addr
)

}

func 
(
n 
*
Net
)
 
ListenTCP
(
nett 
string
,
 laddr 
*
net
.
TCPAddr
)
 
(*
net
.
TCPListener
,
 error
)
 
{

    
if
 err 
:=
 n
.
inherit
();
 err 
!=
 
nil
 
{

        
return
 
nil
,
 err

    
}

    n
.
mutex
.
Lock
()

    defer n
.
mutex
.
Unlock
()

    
// look for an inherited listener

    
for
 i
,
 l 
:=
 range n
.
inherited 
{

        
if
 l 
==
 
nil
 
{
 
// we nil used inherited listeners

            
continue

        
}

        
if
 isSameAddr
(
l
.
Addr
(),
 laddr
)
 
{

            n
.
inherited
[
i
]
 
=
 
nil

            n
.
active 
=
 append
(
n
.
active
,
 l
)

            
return
 l
.(*
net
.
TCPListener
),
 
nil

        
}

    
}

    
// make a fresh listener

    l
,
 err 
:=
 net
.
ListenTCP
(
nett
,
 laddr
)

    
if
 err 
!=
 
nil
 
{

        
return
 
nil
,
 err

    
}

    n
.
active 
=
 append
(
n
.
active
,
 l
)

    
return
 l
,
 
nil

}



func 
(
n 
*
Net
)
 inherit
()
 error 
{

    
var
 retErr error

    n
.
inheritOnce
.
Do
(
func
()
 
{

        n
.
mutex
.
Lock
()

        defer n
.
mutex
.
Unlock
()

        countStr 
:=
 os
.
Getenv
(
envCountKey
)

        
if
 countStr 
==
 
""
 
{

            
return

        
}

        count
,
 err 
:=
 strconv
.
Atoi
(
countStr
)

        
// In tests this may be overridden.

        fdStart 
:=
 n
.
fdStart

        
if
 fdStart 
==
 
0
 
{

            fdStart 
=
 
3

        
}



        
for
 i 
:=
 fdStart
;
 i 
<
 fdStart
+
count
;
 i
++
 
{

            file 
:=
 os
.
NewFile
(
uintptr
(
i
),
 
"listener"
)

            l
,
 err 
:=
 net
.
FileListener
(
file
)

            
if
 err 
!=
 
nil
 
{

                file
.
Close
()

                retErr 
=
 fmt
.
Errorf
(
"error inheriting socket fd %d: %s"
,
 i
,
 err
)

                
return

            
}

            
if
 err 
:=
 file
.
Close
();
 err 
!=
 
nil
 
{

                retErr 
=
 fmt
.
Errorf
(
"error closing inherited socket fd %d: %s"
,
 i
,
 err
)

                
return

            
}

            n
.
inherited 
=
 append
(
n
.
inherited
,
 l
)

        
}

    
})

    
return
 retErr

}

这里ListenTCP 先执行inherit() 将继承来的net.Listener 保存在n.inherited里面，启动时判断是否是继承的listener，没有才 make a fresh listener呢，这里的fdStart 初始值设置为3，就是前面提到的那个数字3 （三个标准输入输出占了3位）。

总结起来启动子进程流程如下：

1、提取listener的fd，修改LISTENFDS环境变量为listener的数量，os.StartProcess启动子进程.

files
[
i
],
 err 
=
 l
.(
filer
).
File
()
 

2、子进程启动执行a.net.Listen()时，根据环境变量LISTENFDS和fdStart 变量取出listener

file 
:=
 os
.
NewFile
(
uintptr
(
i
),
 
"listener"
)

l
,
 err 
:=
 net
.
FileListener
(
file
)

file
.
Close
()

根据fd创建一个文件，通过文件拿到listener的副本，然后关闭文件。最终a.net.Listen()的逻辑是如果是继承端口就返回一个listener副本，如果不是就启动一个新的listener。3、后续执行a.serve() 挂载server，然后通知父进程优雅关闭等逻辑。