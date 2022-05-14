---
title: Go unsafe Pointer
layout: post
category: golang
author: 夏泽民
---
强类型意味着一旦定义了，它的类型就不能改变了；静态意味着类型检查在运行前就做了。

同时为了安全的考虑，Go语言是不允许两个指针类型进行转换的。

指针类型转换
我们一般使用*T作为一个指针类型，表示一个指向类型T变量的指针。为了安全的考虑，两个不同的指针类型不能相互转换，比如*int不能转为*float64。

func main() {
	i:= 10
	ip:=&i

	var fp *float64 = (*float64)(ip)
	
	fmt.Println(fp)
}
以上代码我们在编译的时候，会提示cannot convert ip (type *int) to type *float64，也就是不能进行强制转型。那如果我们还是需要进行转换怎么做呢？这就需要我们使用unsafe包里的Pointer了，下面我们先看看unsafe.Pointer是什么，然后再介绍如何转换。

Pointer
unsafe.Pointer是一种特殊意义的指针，它可以包含任意类型的地址，有点类似于C语言里的void*指针，全能型的。

func main() {
	i:= 10
	ip:=&i

	var fp *float64 = (*float64)(unsafe.Pointer(ip))
	
	*fp = *fp * 3

	fmt.Println(i)
}
以上示例，我们可以把*int转为*float64,并且我们尝试了对新的*float64进行操作，打印输出i，就会发现i的址同样被改变。

以上这个例子没有任何实际的意义，但是我们说明了，通过unsafe.Pointer这个万能的指针，我们可以在*T之间做任何转换。

type ArbitraryType int

type Pointer *ArbitraryType
可以看到unsafe.Pointer其实就是一个*int,一个通用型的指针。

我们看下关于unsafe.Pointer的4个规则。

任何指针都可以转换为unsafe.Pointer
unsafe.Pointer可以转换为任何指针
uintptr可以转换为unsafe.Pointer
unsafe.Pointer可以转换为uintptr
前面两个规则我们刚刚已经演示了，主要用于*T1和*T2之间的转换，那么最后两个规则是做什么的呢？我们都知道*T是不能计算偏移量的，也不能进行计算，但是uintptr可以，所以我们可以把指针转为uintptr再进行偏移计算，这样我们就可以访问特定的内存了，达到对不同的内存读写的目的。

下面我们以通过指针偏移修改Struct结构体内的字段为例，来演示uintptr的用法。

func main() {
	u:=new(user)
	fmt.Println(*u)

	pName:=(*string)(unsafe.Pointer(u))
	*pName="张三"

	pAge:=(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)))
	*pAge = 20

	fmt.Println(*u)
}

type user struct {
	name string
	age int
}
以上我们通过内存偏移的方式，定位到我们需要操作的字段，然后改变他们的值。

第一个修改user的name值的时候，因为name是第一个字段，所以不用偏移，我们获取user的指针，然后通过unsafe.Pointer转为*string进行赋值操作即可。

第二个修改user的age值的时候，因为age不是第一个字段，所以我们需要内存偏移，内存偏移牵涉到的计算只能通过uintptr，所我们要先把user的指针地址转为uintptr，然后我们再通过unsafe.Offsetof(u.age)获取需要偏移的值，进行地址运算(+)偏移即可。


 
现在偏移后，地址已经是user的age字段了，如果要给它赋值，我们需要把uintptr转为*int才可以。所以我们通过把uintptr转为unsafe.Pointer,再转为*int就可以操作了。

这里我们可以看到，我们第二个偏移的表达式非常长，但是也千万不要把他们分段，不能像下面这样。

	temp:=uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)
	pAge:=(*int)(unsafe.Pointer(temp))
	*pAge = 20
逻辑上看，以上代码不会有什么问题，但是这里会牵涉到GC，如果我们的这些临时变量被GC，那么导致的内存操作就错了，我们最终操作的，就不知道是哪块内存了，会引起莫名其妙的问题。
<!-- more -->
