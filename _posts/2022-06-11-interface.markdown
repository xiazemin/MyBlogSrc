---
title: interface
layout: post
category: golang
author: 夏泽民
---
go语言接口有两种表现形式，对应到底层实现也是用两种不用的数据结构表示的。对没有定义方法的接口，底层用 用eface结构表示，对定义的有方法的接口，底层用iface结构表示。这里两种结构定义在 runtime/runtime2.go中 eface和iface都占16字节，eface只有2个字段，因为它代表的是没有方法的接口，只需要存储被赋值对象的类型和数据即可，正好对应到这里的_type 和 data字段。iface代表含有方法的接口，定义里面的 data字段也是表示被存储对象的值，注意这里的值是原始值的一个拷贝，如果原始值是一个值类型，这里的data是执行的数据时原始数据的一个副本。
<!-- more -->
// 空接口，表示不含有method的interface结构
type eface struct {
   // 赋值给空接口变量的类型型，_type是基础数据类型的抽象
   _type *_type
   // 指向数据的指针,对于值类型变量，指向的是值拷贝一份后的地址
   // 对于指针类型变量，指向的是原数据地址
   data unsafe.Pointer
}

/ 非空接口，含有method的interface结构 type iface struct { // itab描述信息有接口的类型和赋值给接口变量的类型，大小等 tab *itab // 指向数据的地址 data unsafe.Pointer }

type itab struct {
   // 描述接口的类型，接口有哪些方法，接口的包名
   inter *interfacetype
   // 描述赋值变量的类型
   _type *_type
   // hash值，用在接口断言时候
   hash  uint32 // copy of _type.hash. Used for type switches.
   _     [4]byte
   // 赋值变量，即接口实现者的方法地址，这里虽然定义了数组长度为1，并不表示只能有1个方法
   // fun是第一个方法的地址，因为方法的地址是一个指针地址，占用固定的8个字节，所以后面的
   // 方法的地址可以根据fun偏移计算得到
   fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
itab比较简单，只有5个字段， inter存储的是interface自己本身的类型，_type存储的是接口被赋值对象的类型，也就一个具体对象的类型，对应前面的Animal对象，inter存在的是Animal类型，inter含有 一个_type类型，同时包含package名pkgpath，mhdr描述Animal含有哪些方法信息。hash值在类型断言的时候用，这里的hash值与*type里面的hash值是一样的。fun表示接口被赋值对象的方法，这里虽然 定义了长度为1的数组，并不表示只有1个方法，uintptr是一个指针地址，占用固定的8个字节，方法地址从fun开始是连续存储的，所以后面的方法地址可以根据fun偏移计算得到。

// 接口的变量的类型
type interfacetype struct {
   // golang 基础类型，struct, array, slice,map...
   typ _type
   // 变量类型定义的结构所在的包位置信息
   pkgpath name
   // method信息
   mhdr []imethod
}

// Needs to be in sync with ../cmd/link/internal/ld/decodesym.go:/^func.commonsize,
// ../cmd/compile/internal/gc/reflect.go:/^func.dcommontype and
// ../reflect/type.go:/^type.rtype.
// ../internal/reflectlite/type.go:/^type.rtype.
type _type struct {
   // 类型占用的内存大小
   size uintptr
   // 包含所有指针的内存前缀大小
   ptrdata uintptr // size of memory prefix holding all pointers
   // 类型的hash值
   hash uint32
   // 标记值，在反射的时候会用到
   tflag tflag
   // 字节对齐方式
   align uint8
   // 结构体字段对齐的字节数
   fieldAlign uint8
   // 基础类型，见下面文中的说明
   kind uint8
   // function for comparing objects of this type
   // (ptr to object A, ptr to object B) -> ==?
   // 比较2个形参对象的类型是否相同 对象A和对象B的类型是否相同
   equal func(unsafe.Pointer, unsafe.Pointer) bool
   // gcdata stores the GC type data for the garbage collector.
   // If the KindGCProg bit is set in kind, gcdata is a GC program.
   // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.

   // gc信息
   gcdata *byte
   // 类型名称字符串在可执行二进制文件段中的偏移量
   str nameOff
   // 类型的元数据信息在可执行二进制文件段中的偏移量
   ptrToThis typeOff
}
_type类型定义如上，我们关注几个核心字段，size类型占用的内存大小，hash值与itab中的hash值一致，在类型断言的时候用，kind表示基础类型，描述类型的元素数据信息，是对具体类型的一种抽象。 kind具体类型如下，定义在 runtime/typekind.go文件中

https://zhuanlan.zhihu.com/p/372817535



根据 channel 中收发元素的类型和缓冲区的大小初始化 runtime.hchan 和缓冲区

   如果 channel 不存在缓冲区，分配 hchan 结构体空间，即无缓存 channel
   如果 channel 存储的类型不是指针类型，分配连续地址空间，包括 hchan 结构体 + 数据
   默认情况包括指针，为 hchan 和 buf 单独分配数据地址空间

https://blog.csdn.net/zhonglinzhang/article/details/83959974

type I interface{
    Get() int
    Put(int)
    A() int     //可以自由添加，只为检验是否增加后，会改变占用字节

}

func main(){
        var i I;
        fmt.Println(unsafe.Sizeof(i)) 
}

答案是16个字节
https://blog.csdn.net/carler_czj/article/details/79870080


