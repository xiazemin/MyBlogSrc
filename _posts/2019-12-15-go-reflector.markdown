---
title: go-reflector
layout: post
category: golang
author: 夏泽民
---
https://github.com/tkrajina/go-reflector
源码分析系列文章已经开源到github，地址如下：
github：
https://github.com/farmer-hutao/k8s-source-code-analysis
gitbook：
https://farmer-hutao.github.io/k8s-source-code-analysis
<!-- more -->
reflector
前面说到 Store 要给 Reflector 服务，我们看一下 Reflector 的定义：

tools/cache/reflector.go:47

type Reflector struct {
   name string
   metrics *reflectorMetrics
   expectedType reflect.Type
   // The destination to sync up with the watch source
   store Store
   // listerWatcher is used to perform lists and watches.
   listerWatcher ListerWatcher
  // ……
}
Copy
Reflector 要做的事情是 watch 一个指定的资源，然后将这个资源的变化反射到给定的store中。很明显这里的两个属性 listerWatcher 和 store 就是这些逻辑的关键。

我们简单看一下往 store 中添加数据的代码：

tools/cache/reflector.go:324

switch event.Type {
case watch.Added:
   err := r.store.Add(event.Object)
   // ……
case watch.Modified:
   err := r.store.Update(event.Object)
   // ……
case watch.Deleted:
   // ……
   err := r.store.Delete(event.Object)
Copy
这个 store 一般用的是 DeltaFIFO，到这里大概就知道 Refactor 从 API Server watch 资源，然后写入 DeltaFIFO 的过程了，大概长这个样子：

1555420426565

然后我们关注一下 DeltaFIFO 的 knownObjects 属性，在创建一个 DeltaFIFO 实例的时候有这样的逻辑：

tools/cache/delta_fifo.go:59

func NewDeltaFIFO(keyFunc KeyFunc, knownObjects KeyListerGetter) *DeltaFIFO {
   f := &DeltaFIFO{
      items:        map[string]Deltas{},
      queue:        []string{},
      keyFunc:      keyFunc,
      knownObjects: knownObjects,
   }
   f.cond.L = &f.lock
   return f
}
Copy
这里接收了 KeyListerGetter 类型的 knownObjects，继续往前跟可以看到我们前面提到的 SharedIndexInformer 的初始化逻辑中将 indexer 对象当作了这里的 knownObjects 的实参：

tools/cache/shared_informer.go:192

fifo := NewDeltaFIFO(MetaNamespaceKeyFunc, s.indexer)
Copy
s.indexer 来自于：NewSharedIndexInformer() 函数的逻辑：

func NewSharedIndexInformer(lw ListerWatcher, objType runtime.Object, defaultEventHandlerResyncPeriod time.Duration, indexers Indexers) SharedIndexInformer {
   realClock := &clock.RealClock{}
   sharedIndexInformer := &sharedIndexInformer{
      processor:                       &sharedProcessor{clock: realClock},
      indexer:                         NewIndexer(DeletionHandlingMetaNamespaceKeyFunc, indexers),
      listerWatcher:                   lw,
      objectType:                      objType,
      resyncCheckPeriod:               defaultEventHandlerResyncPeriod,
      defaultEventHandlerResyncPeriod: defaultEventHandlerResyncPeriod,
      cacheMutationDetector:           NewCacheMutationDetector(fmt.Sprintf("%T", objType)),
      clock:                           realClock,
   }
   return sharedIndexInformer
}
Copy
这里的 NewIndexer() 函数中就可以看到我们前面提到的 Indexer 接口的实现 cache 对象了：

!FILENMAE tools/cache/store.go:239

func NewIndexer(keyFunc KeyFunc, indexers Indexers) Indexer {
   return &cache{
      cacheStorage: NewThreadSafeStore(indexers, Indices{}),
      keyFunc:      keyFunc,
   }
}

