---
title: workerpool
layout: post
category: golang
author: 夏泽民
---
fasthttp 的源码，其中读到了 workpool
<!-- more -->
ackage fasthttp

import (
    "net"
    "runtime"
    "strings"
    "sync"
    "time"
)

// workerPool serves incoming connections via a pool of workers
// in FILO order, i.e. the most recently stopped worker will serve the next
// incoming connection.
//
// Such a scheme keeps CPU caches hot (in theory).
type workerPool struct {
    // Function for serving server connections.
    // It must leave c unclosed.
    WorkerFunc func(c net.Conn) error //注册的conn 处理函数

    MaxWorkersCount int //最大的工作协程数

    LogAllErrors bool

    MaxIdleWorkerDuration time.Duration //协程最大的空闲时间，超过了就清理掉，其实就是退出协程函数 ，退出 go

    Logger Logger

    lock         sync.Mutex
    workersCount int  //当前的工作协程数
    mustStop     bool //workpool 停止标记

    ready []*workerChan //准备工作的协程，记当时还在空闲的协程

    stopCh chan struct{} //workpool 停止信号

    workerChanPool sync.Pool //避免每次频繁分配workerChan，使用pool
}

type workerChan struct { //工作协程
    lastUseTime time.Time
    ch          chan net.Conn // 带缓冲区 chan 处理完了一个conn 通过for range 再处理下一个，都在一个协程里面
}

func (wp *workerPool) Start() {
    if wp.stopCh != nil {
        panic("BUG: workerPool already started")
    }
    wp.stopCh = make(chan struct{})
    stopCh := wp.stopCh
    go func() {
        var scratch []*workerChan
        for {
            wp.clean(&scratch) //定时清理掉协程  （workerChan）
            select {
            case <-stopCh:
                return
            default:
                time.Sleep(wp.getMaxIdleWorkerDuration())
            }
        }
    }()
}

func (wp *workerPool) Stop() {
    if wp.stopCh == nil {
        panic("BUG: workerPool wasn't started")
    }
    close(wp.stopCh) //停止
    wp.stopCh = nil

    // Stop all the workers waiting for incoming connections.
    // Do not wait for busy workers - they will stop after
    // serving the connection and noticing wp.mustStop = true.
    wp.lock.Lock()
    ready := wp.ready
    for i, ch := range ready { //清空
        ch.ch <- nil  //使用nil 值来关闭，而不是close,因为 ch 是池化了的，会循环使用，所以不能close
        ready[i] = nil
    }
    wp.ready = ready[:0]
    wp.mustStop = true
    wp.lock.Unlock()
}

func (wp *workerPool) getMaxIdleWorkerDuration() time.Duration {
    if wp.MaxIdleWorkerDuration <= 0 {
        return 10 * time.Second
    }
    return wp.MaxIdleWorkerDuration
}

func (wp *workerPool) clean(scratch *[]*workerChan) {
    // 传入scratch ，要淘汰的ch, 避免每次分配
    maxIdleWorkerDuration := wp.getMaxIdleWorkerDuration()

    // Clean least recently used workers if they didn't serve connections
    // for more than maxIdleWorkerDuration.
    currentTime := time.Now()

    wp.lock.Lock()
    ready := wp.ready
    n := len(ready)
    i := 0
    for i < n && currentTime.Sub(ready[i].lastUseTime) > maxIdleWorkerDuration {
        i++ //过期的ch 个数
    }
    *scratch = append((*scratch)[:0], ready[:i]...) //淘汰的ch,放到scratch
    if i > 0 {
        m := copy(ready, ready[i:]) //把需要保留的ch，平移到前面，并且几下要保留的数量 m
        for i = m; i < n; i++ {
            ready[i] = nil //把ready 后面的ch淘汰 赋值nil
        }
        wp.ready = ready[:m] //保留的ch到ready
    }
    wp.lock.Unlock()

    // Notify obsolete workers to stop.
    // This notification must be outside the wp.lock, since ch.ch
    // may be blocking and may consume a lot of time if many workers
    // are located on non-local CPUs.
    tmp := *scratch
    for i, ch := range tmp { //淘汰的ch 赋值nil
        ch.ch <- nil
        tmp[i] = nil
    }
}

func (wp *workerPool) Serve(c net.Conn) bool {
    ch := wp.getCh() //获取一个协程
    if ch == nil {
        return false
    }
    ch.ch <- c //传入 conn 到协程
    return true
}

var workerChanCap = func() int {
    // Use blocking workerChan if GOMAXPROCS=1.
    // This immediately switches Serve to WorkerFunc, which results
    // in higher performance (under go1.5 at least).
    if runtime.GOMAXPROCS(0) == 1 {
        return 0
    }

    // Use non-blocking workerChan if GOMAXPROCS>1,
    // since otherwise the Serve caller (Acceptor) may lag accepting
    // new connections if WorkerFunc is CPU-bound.
    return 1
}()

func (wp *workerPool) getCh() *workerChan {
    var ch *workerChan //ch 是一个conn chan 阻塞的，通过for range 不停的处理不同的conn,可以看做是一个协程，不停的处理不同的链接
    createWorker := false

    wp.lock.Lock()
    ready := wp.ready
    n := len(ready) - 1
    if n < 0 {
        if wp.workersCount < wp.MaxWorkersCount {
            createWorker = true
            wp.workersCount++ //没有可用的了需要 new
        }
    } else {
        ch = ready[n]
        ready[n] = nil
        wp.ready = ready[:n] //获取ch 并且ready - 1
    }
    wp.lock.Unlock()

    if ch == nil {
        if !createWorker {
            return nil
        }
        vch := wp.workerChanPool.Get() //new 一个，这里的new 其实是在 pool 拿一个workerChan，从这里可以看出基本上只要是频繁要分配的变量，都使用pool
        if vch == nil {
            vch = &workerChan{
                ch: make(chan net.Conn, workerChanCap),
            }
        }
        ch = vch.(*workerChan)
        go func() { //新建一个协程处理
            wp.workerFunc(ch)
            wp.workerChanPool.Put(vch) //归还 workerChan
        }()
    }
    return ch
}

func (wp *workerPool) release(ch *workerChan) bool {
    ch.lastUseTime = CoarseTimeNow() //更新最后使用这个协程的时间
    wp.lock.Lock()
    if wp.mustStop {
        wp.lock.Unlock()
        return false //如果停止了，则上层 停止协程
    }
    wp.ready = append(wp.ready, ch) //归还 ch 到ready,这里很巧妙，这样 getch 的时候就又可以把新的conn放到这个协程处理
    wp.lock.Unlock()
    return true
}

func (wp *workerPool) workerFunc(ch *workerChan) {
    var c net.Conn

    var err error
    for c = range ch.ch { //不停的获取 ch（阻塞的chan） ，处理不同的conn,在一个协程里面
        if c == nil { //接受到nil 值 就break
            break
        }

        if err = wp.WorkerFunc(c); err != nil && err != errHijacked { //注册的conn链接处理函数
            errStr := err.Error()
            if wp.LogAllErrors || !(strings.Contains(errStr, "broken pipe") ||
                strings.Contains(errStr, "reset by peer") ||
                strings.Contains(errStr, "i/o timeout")) {
                wp.Logger.Printf("error when serving connection %q<->%q: %s", c.LocalAddr(), c.RemoteAddr(), err)
            }
        }
        if err != errHijacked {
            c.Close()
        }
        c = nil //释放conn

        if !wp.release(ch) {
            break
        } //如果stop 了就 break
    }

    wp.lock.Lock()
    wp.workersCount-- //释放此协程
    wp.lock.Unlock()
}
