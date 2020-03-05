---
title: sync.WaitGroup 实现逻辑和源码解析
layout: post
category: golang
author: 夏泽民
---
方便的并发，是Golang的一大特色优势，而使用并发，对sync包的WaitGroup不会陌生。WaitGroup主要用来做Golang并发实例即Goroutine的等待，当使用go启动多个并发程序，通过waitgroup可以等待所有go程序结束后再执行后面的代码逻辑，比如：
func Main() {
    wg := sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            time.Sleep(10 * time.Second)
        }()

    }
    wg.Wait() // 等待在此，等所有go func里都执行了Done()才会退出
}

https://juejin.im/post/5e5b62f86fb9a07cb1578fda
<!-- more -->
实现原理
WaitGroup对外提供三个方法，Add(int),Done()和Wait(), 其中Done()是调用了Add(-1)，一般使用方法是，先统一Add，在goroutine里并发的Done，然后Wait。
WaitGroup主要维护了2个计数器，一个是请求计数器 v，一个是等待计数器 w，二者组成一个64bit的值，请求计数器占高32bit，等待计数器占低32bit。
每次Add执行，请求计数器v加1，Done方法执行，请求计数器减1，v为0时通过信号量唤醒Wait()。
那么等待计数器拿来干嘛？是因为同一个实例的Wait()方法支持多处调用，每一次Wait()方法执行，等待计数器 w 就会加1，而当请求计数器v为0触发Wait()时，要根据w的数量发送w份的信号量，正确的触发所有的Wait()，这虽然不是常用的一个特性，但是在一些特殊场合是有用处的(比如多个并发都依赖于WaitGroup的实例的结束信号来进行下一个action)，演示代码如下：
func main() {
  wg := sync.WaitGroup{}
  for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
​    }()
  }
  time.Sleep(2 * time.Second)
  for j := 0; j < 3; j++ {
    go func(i int) {
      // 3个地方调用Wait()，通过等待j计时器，每个Wati都会被hu唤醒
      wg.Wait()
      fmt.Println("wait done now ", i)
    }(j)
  }
  time.Sleep(10 * time.Second)
  return
}
/*
输出如下，数字出现的顺序随机
wait done now  1
wait done now  0
wait done now  2
*/
复制代码同时，WaitGroup里还对使用逻辑进行了严格的检查，比如Wait()一旦开始不能Add().
下面是带注释的代码，去掉了不影响代码逻辑的trace部分：
func (wg *WaitGroup) Add(delta int) {
    statep := wg.state()
    // 更新statep，statep将在wait和add中通过原子操作一起使用
    state := atomic.AddUint64(statep, uint64(delta)<<32)
    v := int32(state >> 32)
    w := uint32(state)
        if v < 0 {
        panic("sync: negative WaitGroup counter")
    }
    if w != 0 && delta > 0 && v == int32(delta) {
        // wait不等于0说明已经执行了Wait，此时不容许Add
        panic("sync: WaitGroup misuse: Add called concurrently with Wait")
    }
    // 正常情况，Add会让v增加，Done会让v减少，如果没有全部Done掉，此处v总是会大于0的，直到v为0才往下走
    // 而w代表是有多少个goruntine在等待done的信号，wait中通过compareAndSwap对这个w进行加1
     if v > 0 || w == 0 {
        return
    }
    // This goroutine has set counter to 0 when waiters > 0.
    // Now there can't be concurrent mutations of state:
    // - Adds must not happen concurrently with Wait,
    // - Wait does not increment waiters if it sees counter == 0.
    // Still do a cheap sanity check to detect WaitGroup misuse.
    // 当v为0(Done掉了所有)或者w不为0(已经开始等待)才会到这里，但是在这个过程中又有一次Add，导致statep变化，panic
    if *statep != state {
        panic("sync: WaitGroup misuse: Add called concurrently with Wait")
    }
    // Reset waiters count to 0.
    // 将statep清0，在Wait中通过这个值来保护信号量发出后还对这个Waitgroup进行操作
    *statep = 0
    // 将信号量发出，触发wait结束
    for ; w != 0; w-- {
        runtime_Semrelease(&wg.sema, false)
    }
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
    wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
    statep := wg.state()
        for {
        state := atomic.LoadUint64(statep)
        v := int32(state >> 32)
        w := uint32(state)
        if v == 0 {
            // Counter is 0, no need to wait.
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }
        // Increment waiters count.
        // 如果statep和state相等，则增加等待计数，同时进入if等待信号量
        // 此处做CAS，主要是防止多个goroutine里进行Wait()操作，每有一个goroutine进行了wait，等待计数就加1
        // 如果这里不相等，说明statep，在 从读出来 到 CAS比较 的这个时间区间内，被别的goroutine改写了，那么不进入if，回去再读一次，这样写避免用锁，更高效些
        if atomic.CompareAndSwapUint64(statep, state, state+1) {
            if race.Enabled && w == 0 {
                // Wait must be synchronized with the first Add.
                // Need to model this is as a write to race with the read in Add.
                // As a consequence, can do the write only for the first waiter,
                // otherwise concurrent Waits will race with each other.
                race.Write(unsafe.Pointer(&wg.sema))
            }
            // 等待信号量
            runtime_Semacquire(&wg.sema)
            // 信号量来了，代表所有Add都已经Done
            if *statep != 0 {
                // 走到这里，说明在所有Add都已经Done后，触发信号量后，又被执行了Add
                panic("sync: WaitGroup is reused before previous Wait has returned")
            }
            return
        }
    }
}
