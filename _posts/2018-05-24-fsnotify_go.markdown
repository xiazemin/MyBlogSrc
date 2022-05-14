---
title: golang fsnotify
layout: post
category: golang
author: 夏泽民
---
fsnotify是一个文件系统监控库, 它可以支持在如下系统上使用:

Windows
Linux
BSD
OSX
<!-- more -->
fsnotify的github地址是
https://github.com/howeyc/fsnotify

fsnotify是一个文件夹监控应用。可以使用创建一个watcher来对某个文件夹进行监控


文件目录很简单，实际就两个程序文件，fsnotify.go 和 各平台的fsnotify_XXX.go
后一个文件是各个不同平台的实现

example_test.go中给的是最简单的实际应用
先fsnotify.NewWatcher()
再开一个goroutine监听文件夹修改的事件
然后使用watcher.Watch()来监听一个文件夹

fsnotify中的几个public函数Watch，WatchFlags，RemoveWatch是对Watcher的具体封装，函数名一看就明白了什么意思。
这里的flag标志watcher要监听文件夹的哪些事件，Watch默认监听所有事件。

String函数能用string表示出事件。这里学了一招使用events = events[1:] 来达到trim同样的目的。

purgeEvents是将内部事件转成外部事件。这个内部事件指的是syscall包有的对事件的封装和标志位，外部事件指fsnotify对事件的再次封装

下面就到fsnotify_linux.go看linux平台下的实现。
FileEvent类型：
mask，代表事件的掩码，这里的事件码对应的实际上是syscall包中constants对应的一些位置码
cookie，每个事件会分配一个唯一的cookie，这个具体是什么也不理解
Name，触发事件的文件名


下面是一个watch类型
wd，syscall中对文件监控返回的watch id
flags，syscall中对文件的flag

watcher结构：
mu：互斥锁，控制并发，对watcher要进行互斥监控
fd：watcher的文件描述符，不要把这个理解成监控的文件的文件描述符。理解成通知watch消息的文件描述符
watches：要监控的文件夹路径和watch结构的映射
fsnFlags：要监控的事件标志位
paths：要监控watch id和文件夹路径的映射，上面三个其实和起来就能完成了path和watch的互相查找
Error：如果发生错误，从这个channel将错误通知主go routine
internalEvent：文件事件队列，内部的文件事件就放在这个队列中
Event：已经处理的文件事件队列
done：主goroutine监听是否已经结束的通知通道
isClose：是否已经结束的标志位，当然只能自身的goroutine使用

下面看NewWatcher这个函数
这里调用了syscall的InotifyInit来进行初始化
学了一点，当syscall出现错误的时候，可以使用os.NewSyscallError来抛出错误
里面起了两个goroutine
readEvents()和purgeEvents()

purgeEvents()上面已经有了，下面是readEvents
先从w.fd中获取出syscall.InotifyEvent，这个是syscall包的通知事件。这个事件是怎么被塞入这个fd的呢？是syscall的syscall.InotifyAddWatch之后如果文件有修改就会将event写入到这个fd中。这个fd就相当于是一个先进先出的队列了。

读出InitofifyEvent之后就需要将它变成我们这个包中定义的fileEvent。并将这个event放入到internalEvent中去。这里只是捕获消息，并没有对消息进行过滤之类的操作。考虑是否弹出和是否返回是在purgeEvent中进行过滤。

对readEvents读完之后其他的就很好理解了。
addWatch就是调用了一下syscall.InotifyAddWatch
removeWatch就是调用了一下syscall.InotifyRmWatch 

注意: 

当一个文件重命名并移到了另一个目录, 这个文件将不会继续被监控, 除非你监控了这个文件所属的目录.
当一个目录被监控时,如果想监控它的子目录需要自己添加子目录来监控他们
你需要自己来处理Error和Event channels
首先里面有几个核心方法:   NewWatcher, Watch, WatchFlags, RemoveWatch, readEvents和purgeEvents

NewWatcher就是通过调用syscall.InotifyInit()首先建立监控初始化.

然后根据返回文件描述符构造Watcher, 同时起两个goroutine, 分别运行readEvents和purgeEvents, readEvents负责读取新的事件并发送到internalEvent, purgeEvents负责将internalEvent的事件转换到Event channel供外部程序使用.

Watch方法就是通过syscall.InotifyAddWatch建立监控列表,并将路径添加到Watcher结构的paths中.

这里有一些深入学习内部实现的文章, 看完后你会发现fsnotify其实就是在外面加了个壳而已, 结构很简单.

{% highlight golang linenos %}
watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    // Process events
    go func() {
        for {
            select {
            case ev := <-watcher.Event:
                log.Println("event:", ev)
            case err := <-watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.Watch("/tmp")
    if err != nil {
        log.Fatal(err)
    }

    /* ... do stuff ... */
    watcher.Close()
    {% endhighlight %}
