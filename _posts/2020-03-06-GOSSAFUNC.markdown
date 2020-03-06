---
title: GOSSAFUNC Dig101-Go 之 interface 调用的一个优化点
layout: post
category: golang
author: 夏泽民
---
https://gocn.vip/topics/9946
为什么对于以下 interface Stringer 和构造类型 Binary

下面代码conversion会调用转换函数convT64，而devirt不会调用？

func conversion() {
  var b Stringer
  var i Binary = 1
  b = i //convT64
  _=b.String()
}

func devirt() {
  var b Stringer = Binary(1)
  _ = b.String() //static call Binary.String
}
<!-- more -->
这里可以使用 ssa 可视化工具查看，更容易了解每行代码的编译过程 如 GOSSAFUNC=main go1.14 build types/interface/interface.go 生成ssa.html
<img src="{{site.url}}{{site.baseurl}}/img/ssa_html.png"/>

搜索发现相关 issue Devirtualize calls when concrete type behind interface is statically known 和提交 De-virtualize interface calls

原来这个是为了优化如果 interface 内部的构造类型如果可以内联后被静态推断出来的话，就将其直接重写为静态调用

最初主要希望避免一些 interface 调用的 gc 压力（interface 调用在逃逸分析时，会使函数的接受者 (receiver) 和参数 (argument) 逃逸到堆上（而不是留在栈上），增加 gc 压力。不过这一点目前还未实现，参见Use devirtualization in escape analysis）

暂时先优化为静态调用避免转换调用（convXXX），减少代码大小和提升细微的性能

摘录主要处理点如下：

// 对iface=类指针（pointer-shaped）构造类型 记录itab
// 用于后续优化掉 OCONVIFACE
cmd/compile/internal/gc/subr.go:implements
  if isdirectiface(t0) && !iface.IsEmptyInterface() {
    itabname(t0, iface)
  }
cmd/compile/internal/gc/reflect.go:itabname
  itabs = append(itabs, itabEntry{t: t, itype: itype, lsym: s.Linksym()})
// 编译前，获取itabs
cmd/compile/internal/gc/reflect.go:peekitabs
// ssa时利用函数内联和itabs推断可重写为静态调用，避免convXXX
cmd/compile/internal/ssa/rewrite.go:devirt
Go 编译步骤相关参见 Go compiler
这种优化对于常见的返回 interface 的构造函数还是有帮助的。

func New() Interface { return &impl{...} }
要注意返回构造类型需为类指针才可以。

我们可以利用这一点来应用此 interface 调用优化

https://github.com/NewbMiao/
https://mp.weixin.qq.com/s/81mLETTbbNmA86qKHCGOZQ
https://bitfieldconsulting.com/blog/building-a-golang-docker-image
