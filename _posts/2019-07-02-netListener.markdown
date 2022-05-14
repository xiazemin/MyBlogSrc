---
title: netListener.File
layout: post
category: golang
author: 夏泽民
---
POSIX提供了fork和exec调用来启动一个新进程，fork复制父进程，然后通过exec来替换自己要执行的程序。在go中，我们使用exec.Command或者os.StartProcess来达到类似效果。
在启动子进程时，需要让子进程知道，我正处于热更新过程中。通常使用环境变量或者参数来实现，例子中使用了-graceful这个参数。
file := netListener.File() // this returns a Dup()
path := "/path/to/executable"
args := []string{
    "-graceful"}

cmd := exec.Command(path, args...)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.ExtraFiles = []*os.File{file}

err := cmd.Start()
if err != nil {
    log.Fatalf("gracefulRestart: Failed to launch, error: %v", err)
}
然后在子进程中使用net.FileListener来从fd创建一个Listener
func FileListener
func FileListener(f *os.File) (ln Listener, err error)
FileListener returns a copy of the network listener corresponding to the open file f. It is the caller's responsibility to close ln when finished. Closing ln does not affect f, and closing f does not affect ln.

flag.BoolVar(&gracefulChild, "graceful", false, "listen on fd open 3 (internal use only)")

if gracefulChild {
    log.Print("main: Listening to existing file descriptor 3.")
    f := os.NewFile(3, "") // 3就是我们传递的listening fd
    l, err = net.FileListener(f)
} else {
    log.Print("main: Listening on a new file descriptor.")
    l, err = net.Listen("tcp", server.Addr)
}
到这里，子进程就可以Accept并接受连接了，现在我们还需要立刻干掉父进程。使用getpid调用获取到父进程的id，然后kill它。

parent := syscall.Getppid()
syscall.Kill(parent, syscall.SIGTERM)
当然，更加完美的方式还需要父进程可以优雅退出，即不再接受新连接，并且处理完当前所有连接后再退出，如果一段时间内没能处理完，也可以选择直接退出。
<!-- more -->
上面 netListener.File() 与 dup 函数类似，返回的是一个拷贝的文件描述符。另外，该文件描述符不应该设置 FD_CLOEXEC 标识，这将会导致出现我们不想要的结果：子进程的该文件描述符被关闭。

你可能会想到可以使用命令行参数把该文件描述符的值传递给子进程，但相对来说，我使用的这种方式更为简单

最终， args 数组包含了一个 -graceful 选项，你的进程需要以某种方式通知子进程要复用父进程的描述符而不是新打开一个。

子进程初始化
server := &http.Server{Addr: "0.0.0.0:8888"}

var gracefulChild bool
var l net.Listever
var err error

flag.BoolVar(&gracefulChild, "graceful", false, "listen on fd open 3 (internal use only)")

if gracefulChild {
    log.Print("main: Listening to existing file descriptor 3.")
    f := os.NewFile(3, "")
    l, err = net.FileListener(f)
} else {
    log.Print("main: Listening on a new file descriptor.")
    l, err = net.Listen("tcp", server.Addr)
}
通知父进程停止
if gracefulChild {
    parent := syscall.Getppid()
    log.Printf("main: Killing parent pid: %v", parent)
    syscall.Kill(parent, syscall.SIGTERM)
}

server.Serve(l)
父进程停止接收请求并等待当前所处理的所有请求结束
为了做到这一点我们需要使用 sync.WaitGroup 来保证对当前打开的连接的追踪，基本上就是：每当接收一个新的请求时，给wait group做原子性加法，当请求结束时给wait group做原子性减法。也就是说wait group存储了当前正在处理的请求的数量

var httpWg sync.WaitGroup
匆匆一瞥，我发现go中的http标准库并没有为Accept()和Close()提供钩子函数，但这就到了 interface 展现其魔力的时候了(非常感谢 Jeff R. Allen 的这篇 文章 )

下面是一个例子，该例子实现了每当执行Accept()的时候会原子性增加wait group。首先我们先继承 net.Listener 实现一个结构体

type gracefulListener struct {
    net.Listener
    stop    chan error
    stopped bool
}

func (gl *gracefulListener) File() *os.File {
    tl := gl.Listener.(*net.TCPListener)
    fl, _ := tl.File()
    return fl
}
接下来我们覆盖Accept方法(暂时先忽略 gracefulConn )

func (gl *gracefulListener) Accept() (c net.Conn, err error) {
    c, err = gl.Listener.Accept()
    if err != nil {
        return
    }

    c = gracefulConn{Conn: c}

    httpWg.Add(1)
    return
}
我们还需要一个构造函数以及一个Close方法，构造函数中另起一个goroutine关闭，为什么要另起一个goroutine关闭，请看 refer^{[1]}

func newGracefulListener(l net.Listener) (gl *gracefulListener) {
    gl = &gracefulListener{Listener: l, stop: make(chan error)}
    // 这里为什么使用go 另起一个goroutine关闭请看文章末尾
    go func() {
        _ = <-gl.stop
        gl.stopped = true
        gl.stop <- gl.Listener.Close()
    }()
    return
}

func (gl *gracefulListener) Close() error {
    if gl.stopped {
        return syscall.EINVAL
    }
    gl.stop <- nil
    return <-gl.stop
}
我们的 Close 方法简单的向stop chan中发送了一个nil，让构造函数中的goroutine解除阻塞状态并执行Close操作。最终，goroutine执行的函数释放了 net.TCPListener 文件描述符。

接下来，我们还需要一个 net.Conn 的变种来原子性的对wait group做减法

type gracefulConn struct {
    net.Conn
}

func (w gracefulConn) Close() error {
    httpWg.Done()
    return w.Conn.Close()
}
为了让我们上面所写的优雅启动方案生效，我们需要替换 server.Serve(l) 行为:

netListener = newGracefulListener(l)
server.Serve(netListener)
最后补充：我们还需要避免客户端长时间不关闭连接的情况，所以我们创建server的时候可以指定超时时间：

server := &http.Server{
        Addr:           "0.0.0.0:8888",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 16}
