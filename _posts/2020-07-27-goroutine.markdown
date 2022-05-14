---
title: goroutine 的同步和协作
layout: post
category: golang
author: 夏泽民
---
https://www.cnblogs.com/win-for-life/p/13372984.html
竞争条件
一份数据被多个线程共享，可能会产生争用和冲突的情况。这种情况被称为竞态条件，竞态条件会破坏共享数据的一致性，影响一些线程中代码和流程的正确执行。

同步可以解决竞态问题。它本质上是在控制多个线程对共享资源的访问。这种控制主要包含两点：

避免多个线程在同一时刻操作同一个数据块。
协调多个线程，以避免它们在同一时刻执行同一个代码块。
在同步控制下，多个并发运行的线程对这个共享资源的访问是完全串行的。对这个共享资源进行操作的代码片段可以视为一个临界区。

一个互斥锁可以被用来保护一个临界区或者一组相关临界区。它可以保证，在同一时刻只有一个 goroutine 处于该临界区之内。
每当有 goroutine 想进入临界区时，都需要先加锁，每个 goroutine 离开临界区时，都要及时解锁。

var mutex sync.Mutex

func updatePublicResource() {
    mutex.Lock()
    doUpdate()
    mutex.Unlock()
}
使用互斥锁的注意事项：

不要重复锁定互斥锁。
不要忘记解锁互斥锁，推荐使用defer。
不要对尚未锁定或者已解锁的互斥锁解锁。
不要在多个函数之间直接传递互斥锁。(即，不要复制锁)
对一个已经被锁定的互斥锁进行锁定，会阻塞当前的 goroutine 。如果其他的用户级 goroutine 也处于等待状态，整个程序就停止执行了，Go 语言运行时系统会抛出一个死锁的 panic 错误，程序就会崩溃。因此，切记，每一个锁定操作，都要有且只有一个对应的解锁操作。
<!-- more -->

{% raw %}

读写锁是读 / 写互斥锁的简称，读写锁是互斥锁的一种扩展。一个读写锁中包含了两个锁，即：读锁和写锁。
读写锁可以对共享资源的“读操作”和“写操作”进行区别，实现更加细腻的访问控制。
对于某个受到读写锁保护的共享资源，多个写操作不能同时进行，写操作和读操作也不能同时进行，多个读操作可以同时进行。

var mutex sync.RWMutex

func updatePublicResource() {
    mutex.Lock()
    doUpdate()
    mutex.Unlock()
}

func readPublicResource() {
    mutex.RLock()
    read()
    mutex.RUnlock()
}
对写锁进行解锁，会唤醒“所有因试图锁定读锁，而被阻塞的 goroutine”，通常它们都能成功完成对读锁的锁定。
对读锁进行解锁，会在没有其他锁定中读锁的前提下，唤醒“因试图锁定写锁，而被阻塞的 goroutine”；只有一个等待时间最长的被唤醒的 goroutine 能够成功完成对写锁的锁定。
读写锁是互斥锁的扩展，因此有些方面它还是沿用了互斥锁的行为模式。比如，解锁未被锁定的写锁或读锁，会立刻引发 panic。

条件变量是基于互斥锁的，它不用于保护临界区和共享资源，而是用于协调想要访问共享资源的那些线程的。当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程。
io.Pipe 的实现就基于 sync.Cond。
sync.Cond 需要 sync.Locker 类型的参数用于初始化。

type Locker interface {
	Lock()
	Unlock()
}
大多数同步工具禁止在使用后进行复制。Golang 使用两个内嵌字段实现 coCopy 功能：noCopy 和 checker。noCopy 字段用于代码检查工具，checker 字段用于保证运行时不发生复制。

type Cond struct {
    // 用于标识当前结构体在第一次使用后不应该再复制
    // 用于 go vet 编译检查
	noCopy noCopy
	// Cond 基于的锁
	L Locker
    // 一个基于ticket的通知列表
    // 保存了 goroutine 信息的双向链表
	notify  notifyList
	// 保证运行时发生拷贝抛出  panic
	// 在第一次生成时，初始化为 Cond 地址，如果发生复制，复制对象的地址和当前地址将会不同
	checker copyChecker
}
sync.Cond 提供 3 个方法：

Broadcast()：唤醒所有等待 Cond 的 goroutine。不需要在锁的保护下进行。
Signal()：唤醒一个等待 Cond 的 goroutine。不需要在锁的保护下进行。
Wait():解锁互斥锁，挂起当前 goroutine。当 Broadcast 或 Signal 唤醒这个 goroutine，Wait 在返回前会再锁定互斥锁。因此 Wait() 需要在锁的保护下进行。
var lock sync.RWMutex
var sendCond, recvCond *sync.Cond

func init() {
    sendCond = sync.NewCond(&lock)
    recvCond = sync.NewCond(&lock) // 获取读写锁中的读锁
}

func send() {
    lock.Lock()
    for !writeCondition() {
        sendCond.Wait()
    }
    writeResource()
    lock.Unlock()
    recvCond.Signal()// 如果有多个接收的 goroutine 就使用 recvCond.Broadcast()
}

func receive() {
    lock.Lock()
    for !readCondition() {
        recvCond.Wait()
    }
    receiveResource()
    lock.Unlock()
    sendCond.Signal()// 如果有多个发送的 goroutine 就使用
}
有时 sync.Cond 的功能用 channel 也能实现，不过 channel 的意义更多地在于传递数据，而 sync.Cond 的意义在于协程的协作；并且 sync.Cond 更为底层，效率更高。

Cond 在第一次使用后不能复制。
条件变量的通知具有即时性。如果发送通知的时候没有 goroutine 为此等待，该通知就会被直接丢弃。
Signal() 和 Broadcast() 需要在非锁定的情况下调用，因为 Wait() 的调用方处于阻塞状态，可能错过通知。
Wait() 的调用需要基于锁定状态。
func (c *Cond) Wait() {
    // 检查是否发生复制
	c.checker.check()
	// 将当前 gorouitne 加入当前条件变量的通知队列
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	// 阻塞当前的 goroutine，直至收到通知
	runtime_notifyListWait(&c.notify, t)
	// 收到通知后，加锁，进入临界区
	c.L.Lock()
}
为什么要由调用方先加锁，再由Wait()解锁？
调用方在对共享资源的条件进行判断时，保证共享资源的状态不被修改，因此进行加锁。
而当共享资源不满足当前goroutine的条件时，需要让出共享资源的执行权，以便其他 goroutine 对其进行修改，因此进行解锁。

为什么使用for循环多次多次检查共享资源条件？

如果存在多个 goroutine 同时等待通知，最终只有一个 goroutine 可以成功获得执行权限。那么其他的 goroutine 应该在检查不满足执行条件后继续等待。
共享资源存在多种状态，状态改变通知是基于锁的，无法实现更细腻的判断。这时需要每个 goroutine 对自己所需的状态反复检查。
即使共享资源的状态只有两个，并且每种状态都只有一个 goroutine 在关注，如上文展示，也应当使用 for 循环。因为一个 gorouinte 即使没有收到条件通知，也可能被唤醒。这是多核 CPU 计算机硬件层面的调度机制。
条件变量适合保护那些可执行两个对立操作的共享资源。比如，一个既可读又可写的共享文件。又比如，既有生产者又有消费者的产品池。
对于有着对立操作的共享资源（比如一个共享文件），我们通常需要基于同一个读写锁的两个条件变量（比如 rcond 和 wcond）分别保护读操作和写操作（比如 rcond 保护读，wcond 保护写）。读操作在操作完成后要向 wcond 发通知；写操作在操作完成后要向 rcond 发通知。
// 针对读写操作的控制只在初始化时有所变化
var lock sync.RWMutex
var sendCond, recvCond *sync.Cond

func init() {
    sendCond = sync.NewCond(&lock)
    recvCond = sync.NewCond(&lock.RLocker())
}
互斥锁可以保证临界区中代码的串行执行，但却不能保证这些代码执行的原子性（atomicity）。
只有原子操作才能保证代码片段的原子性，原子操作由底层的 CPU 提供了芯片级别的支持。
针对同一共享资源的原子操作不能同时进行，针对不同共享资源的原子操作可以同时进行。
因为原子操作不能被中断，所以它需要足够简单和快速。
sync/atomic 提供了以下操作：

加法（add）
比较并交换（compare and swap，简称 CAS）
加载（load）
存储（store）
交换（swap）
支持的数据类型有：

int32
int64
uint32
uint64
uintptr
unsafe.Pointer
CAS 包含2步操作，但 Load、Store 这类操作只有一步，不具原子性吗？
即使像 a = 1 这种简单的赋值操作也并不一定能够一次完成。如果右边的值的存储宽度超出了计算机的字宽，那么实际的步骤就会多于一个（或者说底层指令多于一个）。比如，你计算机是32位的，但是你要把一个Int64类型的数赋给变量a，那么底层指令就肯定多于一个。在这种情况下，多个底层指令的执行期间是可以被打断的，也就是说CPU在这时可以被切换到别的任务上。如果新任务恰巧要读写这个变量a，那么就会出现值不完整的问题。况且，就算是 a = 1，操作系统和CPU也都不保证这个操作一定不会被打断。只要被打断，就很有可能出现并发访问上的问题，并发安全性也就被破坏了。
所以，当有多个goroutine在并发的读写同一变量时，它们之间就可能会造成干扰。这种操作不是原子性，并发安全性也无法得到保障。

// 法一
var num uint32
num = 100
delta := int32(-3)
atomic.AddUint32(&num, uint32(delta))
fmt.Println(num) // 97

// 法二
var num uint32
num = 100
delta := -3
atomic.AddUint32(&num, ^uint32(-delta-1))
fmt.Println(num) // 97
自旋锁（spinlock）是指当一个线程在获取锁的时候，如果锁已经被其它线程获取，那么该线程将循环等待，然后不断的判断锁是否能够被成功获取，直到获取到锁才会退出循环。
获取锁的线程一直处于活跃状态，但是并没有执行任何有效的任务，使用这种锁会造成busy-waiting。
自旋锁利用了 CPU 层面的指令，因此性能比互斥锁高很多。适合简单对象的操作以及冲突较少的场景。

var num int32 = 10
for {
 if atomic.CompareAndSwapInt32(&num, 10, 0) {
  fmt.Println("The second number has gone to zero.")
  break
 }
 time.Sleep(time.Millisecond * 500)
}
这在效果上与互斥锁有些类似。我们在使用互斥锁的时候，总是假设共享资源的状态会被其他的 goroutine 频繁地改变。而for语句加 CAS 操作的假设往往是：共享资源状态的改变并不频繁，或者，它的状态总会变成期望的那样。这是一种更加乐观，或者说更加宽松的做法。

当真正使用了一个 atomic.Value 变量（第一次赋值）后，就不应该再进行复制操作了。
不能存储 nil 值。不过对于接口类型的变量，它的动态值是 nil，动态类型不是 nil，它就不是 nil。
对于一个原子变量，向它存储的第一个值决定了它的可存储类型。即使是同一接口的不同类型，也是禁止更换的。对于暴露给外部的存储函数，应当先判断其存储值的合法性。
存储引用类型时，注意不要把指针暴露给外部。
sync.Pool 是一个临时对象池。初次使用后禁止复制。它存储的对象应该满足以下特征：

不需要持久使用，对程序来说可有可无，对象的创建和销毁不会影响程序功能。因为 Go 语言的 GC 每次执行时都会将临时对象池清空。
池子中的每一个对象都可以相互替代。
因此，sync.Pool 很适合作为缓存池。
GC 是如何清理临时对象池的？
sync 初始化时，向运行时系统注册一个函数，这个函数用于清除所有已创建的临时对象池中的值。这个函数在每次 GC 运行时被调用。sync 包中有一个全局变量 allPools 负责保存使用中的池列表，供池清理函数使用。

type Pool struct {
	noCopy noCopy

	local     unsafe.Pointer // per-P pool, 实际类型是 [P]poolLocal
	localSize uintptr        // size of the local array

	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // size of victims array

	// 创建一个临时对象
	New func() interface{}
}

// Local per-P Pool
type poolLocalInternal struct {
	private interface{} // 只能由当前 P 使用
	shared  poolChain   // 双向队列，Local P can pushHead/popHead; any P can popTail.
}
Pool 提供了 Put 和 Get 方法用于存取临时对象。存取临时对象时，优先操作private，其次是 poolLocal 的共享临时对象列表 shared （先访问 goroutine 关联的 P 对应的 poolLocal，再访问非关联的 poolLocal ）。当 Get 无法找到可用的临时对象，就会调用 New 创建以一个新的临时对象。

sync.Map 是一个并发安全的字典。

// 可自定义键类型和值类型的并发安全字典

type ConcurrentMap struct {
	m         sync.Map
	keyType   reflect.Type
	valueType reflect.Type
}

func NewConcurrentMap(keyType, valueType reflect.Type) (*ConcurrentMap, error) {
	if keyType == nil {
		return nil, errors.New("nil key type")
	}
	if !keyType.Comparable() {
		return nil, fmt.Errorf("incomparable key type: %s", keyType)
	}
	if valueType == nil {
		return nil, errors.New("nil value type")
	}
	cMap := &ConcurrentMap{
		keyType:   keyType,
		valueType: valueType,
	}
	return cMap, nil
}

func (cMap *ConcurrentMap) Delete(key interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	cMap.m.Delete(key)
}

func (cMap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	return cMap.m.Load(key)
}

func (cMap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	actual, loaded = cMap.m.LoadOrStore(key, value)
	return
}

func (cMap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	cMap.m.Range(f)
}

func (cMap *ConcurrentMap) Store(key, value interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cMap.m.Store(key, value)
}
type Map struct {
	mu Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	//
	// The read field itself is always safe to load, but must only be stored with
	// mu held.
	//
	// Entries stored in read may be updated concurrently without mu, but updating
	// a previously-expunged entry requires that the entry be copied to the dirty
	// map and unexpunged with mu held.
	read atomic.Value // readOnly

	// dirty contains the portion of the map's contents that require mu to be
	// held. To ensure that the dirty map can be promoted to the read map quickly,
	// it also includes all of the non-expunged entries in the read map.
	//
	// Expunged entries are not stored in the dirty map. An expunged entry in the
	// clean map must be unexpunged and added to the dirty map before a new value
	// can be stored to it.
	//
	// If the dirty map is nil, the next write to the map will initialize it by
	// making a shallow copy of the clean map, omitting stale entries.
	dirty map[interface{}]*entry

	// misses counts the number of loads since the read map was last updated that
	// needed to lock mu to determine whether the key was present.
	//
	// Once enough misses have occurred to cover the cost of copying the dirty
	// map, the dirty map will be promoted to the read map (in the unamended
	// state) and the next store to the map will make a new dirty copy.
	misses int
}
Map.read 相当于字典的快照，支持更新和查询操作，原子操作，不需要持有锁。Map.dirty 是原生字典，支持增删改查操作，所有操作需要持有锁 mu 。
Map.read 和 Map.dirty 中存储的键值都是指针，而不是基本值。
查找键值对时，首先去 read 字典查找，如果没找到，再加锁去 dirty 字典查找。
存储键值对时，如果 read 字典中存在这个键，就直接更新。如果这个键被标记为“已删除”，则保存到 dirty 字典，清除“已删除”的标记。
删除键值时，如果只读字典中不存在该键值对，就直接在 dirty 字典中进行删除。如果只读字典中存在该键值对，还要对其进行逻辑删除（标记为“已删除”）。
在脏字典中查找键值对次数足够多的时候，sync.Map 会把脏字典直接作为只读字典，保存在它的 read 字段中，然后把代表脏字典的 dirty 字段的值置为 nil。在这之后，一旦再有新的键值对存入，它就会依据只读字典去重建脏字典。这个时候，它会把只读字典中已被逻辑删除的键值对过滤掉。
总的来说，只读字典可能只包含部分键值对（含逻辑删除键值对），而脏字典中始终包含全量的键值对（不含逻辑删除键值对）。
sync.Map 适用于读多写少的情况，如果写数据比较频繁可以参考：https://github.com/orcaman/concurrent-map

用于同步 goroutine 的协作流程。它可以使一个 goroutine 在其协程完成后才继续执行后续任务。
开始使用后禁止复制。

var wg sync.WaitGroup
func main() {
    wg.Add(3)
    for i := 0; i < 3; i++ {
        go doSomething()
    }
    wg.Wait()
}
func doSomething() {
    defer wg.Done()
}
禁止同时调用 WaitGroup 的 Add() 和 Wait()，即杜绝并发执行用 WaitGroup 的方法。原因是在 Wait() 执行时更改其计数器的值会引发 panic。
执行首次被调用时的入参函数，并且只执行一次。

func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.

	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
由于 Once.Do() 保证在返回前 f() 已经执行完成，如果存在多个 goroutine 并发调用 Do()，会导致除了获胜者，其余 goroutine 都被阻塞在 o.m.Lock() 上。如果 f() 阻塞，可能会导致死锁。
Once.Do() 不保证 f() 执行成功。

func coordinateWithContext() {
    cxt, cancelFunc := context.WithCancel(context.Background())
    // 启动 3 个具有相关任务的协程
    // 如果有一个协程出现问题，取消其他协程
    for i := 1; i < 3; i++ {
        go func() {
            r, e := fn(ctx)
            if e != nil {
                cancelFunc()
            }
        }
    }
    time.Sleep(10 * time.Second)
    fmt.Println("End.")
}

func fn(ctx context.Context) string, error {
    resp := make(chan string)
    err := make(chan error)
    go func(){
        responseString, e := doSomething()
        if e != nil {
            err <- e
        } else {
            resp <- responseString
        }
    }()
    select {
        case <- ctx.Done():
            return "", ctx.Err()
        case r:= <- resp
            return r, nil
        case e := <- err
            return "", e
    }
}
Context.Done() 返回一个 <-chan struct{} 类型的值，这是一个接收通道。调用 cancelFunc() 时，该通道会关闭，阻塞的接收操作会立刻返回。
Context 类型值的撤销操作会联动它的子值。

Context 类型还提供了 WithDeadline() 和 WithTimeout() 方法，生成拥有生命周期的 Context 类型。
此外，Context.WithValue() 可以提供协程间的数据传输功能。在 Context 中查询数据时，先在当前 Context 中查找，如果没找到，再去父值中查找。不过 Context 不提供数据更新的方法，只能通过 在子值中覆盖同名数据、或撤销 Context 丢弃数据 间接实现。


{% endraw %}

https://www.cnblogs.com/Tassdar/p/13373289.html